package incidentrecovery_test

import (
	"os"
	"path/filepath"
	"testing"

	"cosmossdk.io/core/header"
	sdkmath "cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"

	"github.com/cheqd/cheqd-node/app/testsuite"
	"github.com/cheqd/cheqd-node/app/upgrades/incidentrecovery"
)

// This is the real integration test: it runs the actual registered
// app.UpgradeKeeper upgrade handler (wired in app.go's
// RegisterUpgradeHandlers, not a hand-written mock) against a real, fully
// constructed app instance with real staking/bank/distribution keepers.
//
// It specifically exercises TESTING_PLAN.md stage 1's "Account C" scenario
// — the single most important case this whole mechanism exists for: an
// unbonding delegation behind a validator that is jailed, which the normal
// MsgCancelUnbondingDelegation path rejects outright with no workaround.
type IncidentRecoveryUpgradeTestSuite struct {
	testsuite.UpgradeTestSuite
}

func TestIncidentRecoveryUpgrade(t *testing.T) {
	suite.Run(t, new(IncidentRecoveryUpgradeTestSuite))
}

func (s *IncidentRecoveryUpgradeTestSuite) TestJailedValidatorUnbondingIsRecovered() {
	s.Setup()
	ctx := s.Ctx

	affectedAddr := s.TestAccs[0]
	destAddr := s.TestAccs[1]

	bondDenom, err := s.App.StakingKeeper.BondDenom(ctx)
	s.Require().NoError(err)

	// --- Fund the affected address with both a liquid balance and enough
	// to delegate, via the mint module (the standard way to conjure test
	// funds — no dedicated FundAccount helper exists in this repo). ---
	liquidAmt := sdkmath.NewInt(1_000_000)
	delegateAmt := sdkmath.NewInt(5_000_000)
	totalMint := sdk.NewCoins(sdk.NewCoin(bondDenom, liquidAmt.Add(delegateAmt)))
	s.Require().NoError(s.App.BankKeeper.MintCoins(ctx, minttypes.ModuleName, totalMint))
	s.Require().NoError(s.App.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, affectedAddr, totalMint))

	// --- Create and bond a validator directly via the staking keeper
	// (bypassing the full MsgCreateValidator flow — we only need a real,
	// queryable validator + consensus address pair, not full consensus
	// power-index bookkeeping, for this test). ---
	valPriv := ed25519.GenPrivKey()
	valPub := valPriv.PubKey()
	valAddr := sdk.ValAddress(valPub.Address())
	consAddr := sdk.ConsAddress(valPub.Address())

	validator, err := stakingtypes.NewValidator(valAddr.String(), valPub, stakingtypes.Description{Moniker: "jailed-test-validator"})
	s.Require().NoError(err)
	validator.Status = stakingtypes.Bonded
	s.Require().NoError(s.App.StakingKeeper.SetValidator(ctx, validator))
	s.Require().NoError(s.App.StakingKeeper.SetValidatorByConsAddr(ctx, validator))
	// MsgCreateValidator normally triggers this hook automatically, which
	// initializes the validator's distribution historical-rewards record.
	// Without it, the very first Delegate() to this validator panics inside
	// the distribution hooks (BeforeDelegationCreated -> IncrementValidatorPeriod
	// -> decrementReferenceCount) because no record exists yet to decrement.
	s.Require().NoError(s.App.DistrKeeper.Hooks().AfterValidatorCreated(ctx, valAddr))

	// --- Delegate the full delegateAmt to it. ---
	storedValidator, err := s.App.StakingKeeper.GetValidator(ctx, valAddr)
	s.Require().NoError(err)
	shares, err := s.App.StakingKeeper.Delegate(ctx, affectedAddr, delegateAmt, stakingtypes.Unbonded, storedValidator, true)
	s.Require().NoError(err)
	s.Require().True(shares.IsPositive())

	// --- Undelegate it, creating a real unbonding delegation entry with a
	// completion time 14 days out — i.e. still fully pending, not mature,
	// so CompleteUnbonding (the SDK's own maturity-gated path) would
	// correctly refuse to touch it. ---
	_, undelegatedAmt, err := s.App.StakingKeeper.Undelegate(ctx, affectedAddr, valAddr, shares)
	s.Require().NoError(err)
	s.Require().True(undelegatedAmt.Equal(delegateAmt))

	// --- NOW jail the validator. This is the crux of the whole test: the
	// unbonding entry above already exists, behind a validator that is
	// now jailed. MsgCancelUnbondingDelegation would reject this outright.
	s.Require().NoError(s.App.StakingKeeper.Jail(ctx, consAddr))
	jailedValidator, err := s.App.StakingKeeper.GetValidator(ctx, valAddr)
	s.Require().NoError(err)
	s.Require().True(jailedValidator.Jailed, "validator must actually be jailed for this test to mean anything")

	// --- Write the incident-recovery config pointing at this address, and
	// point the handler at it via the env var it reads lazily at
	// upgrade-execution time (see upgrade.go's ConfigPath). ---
	cfg := map[string]interface{}{
		"upgrade_name":                incidentrecovery.UpgradeName,
		"destination_address":         destAddr.String(),
		"destination_address_confirm": destAddr.String(),
		"affected_addresses":          []string{affectedAddr.String()},
	}
	cfgBytes, err := yaml.Marshal(cfg)
	s.Require().NoError(err)
	cfgPath := filepath.Join(s.T().TempDir(), "config.yaml")
	s.Require().NoError(os.WriteFile(cfgPath, cfgBytes, 0o600))
	s.T().Setenv(incidentrecovery.ConfigPathEnvVar, cfgPath)

	// --- Schedule the upgrade and advance to its height, same pattern as
	// app/upgrades/v4/upgrade_test.go uses for the real app's other
	// upgrade handlers. ---
	const upgradeHeight = 10
	ctx = ctx.WithBlockHeight(upgradeHeight - 1).WithHeaderInfo(header.Info{Height: upgradeHeight - 1})
	plan := upgradetypes.Plan{Name: incidentrecovery.UpgradeName, Height: upgradeHeight}
	s.Require().NoError(s.App.UpgradeKeeper.ScheduleUpgrade(ctx, plan))

	ctx = s.Ctx.WithBlockHeight(upgradeHeight).WithHeaderInfo(header.Info{Height: upgradeHeight})
	s.Require().NotPanics(func() {
		_, err = s.App.PreBlocker(ctx, &types.RequestFinalizeBlock{})
	})
	s.Require().NoError(err)

	// --- Verify: the destination received the full amount (liquid +
	// recovered-from-unbonding), and it happened despite the validator
	// still being jailed throughout — proving the bypass is real, not
	// just asserted in a comment. ---
	destBalance := s.App.BankKeeper.GetBalance(ctx, destAddr, bondDenom)
	expectedTotal := liquidAmt.Add(delegateAmt)
	s.Require().True(destBalance.Amount.Equal(expectedTotal),
		"expected destination to receive %s%s (liquid %s + recovered unbonding %s), got %s",
		expectedTotal, bondDenom, liquidAmt, delegateAmt, destBalance)

	sourceBalance := s.App.BankKeeper.GetBalance(ctx, affectedAddr, bondDenom)
	s.Require().True(sourceBalance.IsZero(), "expected affected address to be fully swept, got %s", sourceBalance)

	stillJailed, err := s.App.StakingKeeper.GetValidator(ctx, valAddr)
	s.Require().NoError(err)
	s.Require().True(stillJailed.Jailed, "the validator was never unjailed during this test — the recovery did not depend on, or trigger, an unjail")
}
