package keeper_test

import (
	"fmt"
	"strings"
	"testing"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/secp256k1"
	"github.com/cosmos/cosmos-sdk/baseapp"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtestutil "github.com/cosmos/cosmos-sdk/x/staking/testutil"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/suite"

	cheqdapp "github.com/cheqd/cheqd-node/app"
	"github.com/cheqd/cheqd-node/x/oracle/keeper"
	"github.com/cheqd/cheqd-node/x/oracle/types"
	cmtrand "github.com/cometbft/cometbft/libs/rand"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
)

const (
	displayDenom string = "CHEQ"
	bondDenom    string = "ncheq" // should match the bond denom in appparams
)

type IntegrationTestSuite struct {
	suite.Suite

	ctx         sdk.Context
	app         *cheqdapp.TestApp
	queryClient types.QueryClient
	msgServer   types.MsgServer
}

const (
	initialPower = int64(10000000000)
)

func (s *IntegrationTestSuite) SetupTest() {
	require := s.Require()
	app, err := cheqdapp.Setup(false)
	require.NoError(err)
	ctx := app.BaseApp.NewContextLegacy(false, cmtproto.Header{
		ChainID: fmt.Sprintf("test-chain-%s", cmtrand.Str(4)),
		Height:  9,
	},
	)
	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, keeper.NewQuerier(app.OracleKeeper))
	newParams := stakingtypes.DefaultParams()
	newParams.BondDenom = bondDenom
	err = app.StakingKeeper.SetParams(ctx, newParams)
	require.NoError(err)
	sh := stakingtestutil.NewHelper(s.T(), ctx, app.StakingKeeper)
	amt := sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction)
	sh.Denom = bondDenom
	// mint and send coins to validators
	require.NoError(app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, initCoins))
	require.NoError(app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, initCoins))
	require.NoError(app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, initCoins))
	require.NoError(app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr2, initCoins))

	sh.CreateValidator(valAddr, valPubKey, amt, true)
	sh.CreateValidator(valAddr2, valPubKey2, amt, true)

	_, err = app.StakingKeeper.EndBlocker(ctx)
	require.NoError(err)

	s.app = app
	s.ctx = ctx
	s.queryClient = types.NewQueryClient(queryHelper)
	s.msgServer = keeper.NewMsgServerImpl(app.OracleKeeper)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

var (
	valPubKeys = simtestutil.CreateTestPubKeys(2)

	valPubKey = valPubKeys[0]
	pubKey    = secp256k1.GenPrivKey().PubKey()
	addr      = sdk.AccAddress(pubKey.Address())
	valAddr   = sdk.ValAddress(pubKey.Address())

	valPubKey2 = valPubKeys[1]
	pubKey2    = secp256k1.GenPrivKey().PubKey()
	addr2      = sdk.AccAddress(pubKey2.Address())
	valAddr2   = sdk.ValAddress(pubKey2.Address())

	initTokens = sdk.TokensFromConsensusPower(initialPower, sdk.DefaultPowerReduction)
	initCoins  = sdk.NewCoins(sdk.NewCoin(bondDenom, initTokens))
)

// NewTestMsgCreateValidator test msg creator
func NewTestMsgCreateValidator(address sdk.ValAddress, pubKey cryptotypes.PubKey, amt math.Int) *stakingtypes.MsgCreateValidator {
	commission := stakingtypes.NewCommissionRates(math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec())
	msg, _ := stakingtypes.NewMsgCreateValidator(
		address.String(), pubKey, sdk.NewCoin(bondDenom, amt),
		stakingtypes.Description{}, commission, math.OneInt(),
	)
	return msg
}

func (s *IntegrationTestSuite) TestSetFeederDelegation() {
	app, ctx := s.app, s.ctx
	feederAddr := sdk.AccAddress([]byte("addr________________"))
	feederAcc := app.AccountKeeper.NewAccountWithAddress(ctx, feederAddr)
	app.AccountKeeper.SetAccount(ctx, feederAcc)

	err := s.app.OracleKeeper.ValidateFeeder(ctx, addr, valAddr)
	s.Require().NoError(err)
	err = s.app.OracleKeeper.ValidateFeeder(ctx, feederAddr, valAddr)
	s.Require().Error(err)

	s.app.OracleKeeper.SetFeederDelegation(ctx, valAddr, feederAddr)

	err = s.app.OracleKeeper.ValidateFeeder(ctx, addr, valAddr)
	s.Require().Error(err)
	err = s.app.OracleKeeper.ValidateFeeder(ctx, feederAddr, valAddr)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TestGetFeederDelegation() {
	app, ctx := s.app, s.ctx

	feederAddr := sdk.AccAddress([]byte("addr________________"))
	feederAcc := app.AccountKeeper.NewAccountWithAddress(ctx, feederAddr)
	app.AccountKeeper.SetAccount(ctx, feederAcc)

	s.app.OracleKeeper.SetFeederDelegation(ctx, valAddr, feederAddr)
	resp, err := app.OracleKeeper.GetFeederDelegation(ctx, valAddr)
	s.Require().NoError(err)
	s.Require().Equal(resp, feederAddr)
}

func (s *IntegrationTestSuite) TestAggregateExchangeRatePrevote() {
	app, ctx := s.app, s.ctx

	prevote := types.AggregateExchangeRatePrevote{
		Hash:        "hash",
		Voter:       addr.String(),
		SubmitBlock: 0,
	}
	app.OracleKeeper.SetAggregateExchangeRatePrevote(ctx, valAddr, prevote)

	_, err := app.OracleKeeper.GetAggregateExchangeRatePrevote(ctx, valAddr)
	s.Require().NoError(err)

	app.OracleKeeper.DeleteAggregateExchangeRatePrevote(ctx, valAddr)

	_, err = app.OracleKeeper.GetAggregateExchangeRatePrevote(ctx, valAddr)
	s.Require().Error(err)
}

func (s *IntegrationTestSuite) TestAggregateExchangeRatePrevoteError() {
	app, ctx := s.app, s.ctx

	_, err := app.OracleKeeper.GetAggregateExchangeRatePrevote(ctx, valAddr)
	s.Require().Errorf(err, types.ErrNoAggregatePrevote.Error())
}

func (s *IntegrationTestSuite) TestAggregateExchangeRateVote() {
	app, ctx := s.app, s.ctx

	var decCoins sdk.DecCoins
	decCoins = append(decCoins, sdk.DecCoin{
		Denom:  displayDenom,
		Amount: math.LegacyZeroDec(),
	})

	vote := types.AggregateExchangeRateVote{
		ExchangeRates: decCoins,
		Voter:         addr.String(),
	}
	app.OracleKeeper.SetAggregateExchangeRateVote(ctx, valAddr, vote)

	_, err := app.OracleKeeper.GetAggregateExchangeRateVote(ctx, valAddr)
	s.Require().NoError(err)

	app.OracleKeeper.DeleteAggregateExchangeRateVote(ctx, valAddr)

	_, err = app.OracleKeeper.GetAggregateExchangeRateVote(ctx, valAddr)
	s.Require().Error(err)
}

func (s *IntegrationTestSuite) TestAggregateExchangeRateVoteError() {
	app, ctx := s.app, s.ctx

	_, err := app.OracleKeeper.GetAggregateExchangeRateVote(ctx, valAddr)
	s.Require().Errorf(err, types.ErrNoAggregateVote.Error())
}

func (s *IntegrationTestSuite) TestSetExchangeRateWithEvent() {
	app, ctx := s.app, s.ctx
	err := app.OracleKeeper.SetExchangeRateWithEvent(ctx, displayDenom, math.LegacyOneDec())
	s.Require().NoError(err)
	rate, err := app.OracleKeeper.GetExchangeRate(ctx, displayDenom)
	s.Require().NoError(err)
	s.Require().Equal(rate, math.LegacyOneDec())
}

func (s *IntegrationTestSuite) TestGetExchangeRate_InvalidDenom() {
	app, ctx := s.app, s.ctx

	_, err := app.OracleKeeper.GetExchangeRate(ctx, "uxyz")
	s.Require().Error(err)
}

func (s *IntegrationTestSuite) TestGetExchangeRate_NotSet() {
	app, ctx := s.app, s.ctx

	_, err := app.OracleKeeper.GetExchangeRate(ctx, displayDenom)
	s.Require().Error(err)
}

func (s *IntegrationTestSuite) TestGetExchangeRate_Valid() {
	app, ctx := s.app, s.ctx

	app.OracleKeeper.SetExchangeRate(ctx, displayDenom, math.LegacyOneDec())
	rate, err := app.OracleKeeper.GetExchangeRate(ctx, displayDenom)
	s.Require().NoError(err)
	s.Require().Equal(rate, math.LegacyOneDec())

	app.OracleKeeper.SetExchangeRate(ctx, strings.ToLower(displayDenom), math.LegacyOneDec())
	rate, err = app.OracleKeeper.GetExchangeRate(ctx, displayDenom)
	s.Require().NoError(err)
	s.Require().Equal(rate, math.LegacyOneDec())
}

func (s *IntegrationTestSuite) TestGetExchangeRateBase() {
	oracleParams := s.app.OracleKeeper.GetParams(s.ctx)
	var exponent uint64
	for _, denom := range oracleParams.AcceptList {
		if denom.BaseDenom == bondDenom {
			exponent = uint64(denom.Exponent)
		}
	}

	power := math.LegacyMustNewDecFromStr("10").Power(exponent)

	s.app.OracleKeeper.SetExchangeRate(s.ctx, displayDenom, math.LegacyOneDec())
	rate, err := s.app.OracleKeeper.GetExchangeRateBase(s.ctx, bondDenom)
	s.Require().NoError(err)
	s.Require().Equal(rate.Mul(power), math.LegacyOneDec())

	s.app.OracleKeeper.SetExchangeRate(s.ctx, strings.ToLower(displayDenom), math.LegacyOneDec())
	rate, err = s.app.OracleKeeper.GetExchangeRateBase(s.ctx, strings.ToLower(bondDenom))
	s.Require().NoError(err)
	s.Require().Equal(rate.Mul(power), math.LegacyOneDec())
}

func (s *IntegrationTestSuite) TestClearExchangeRate() {
	app, ctx := s.app, s.ctx

	app.OracleKeeper.SetExchangeRate(ctx, displayDenom, math.LegacyOneDec())
	app.OracleKeeper.ClearExchangeRates(ctx)
	_, err := app.OracleKeeper.GetExchangeRate(ctx, displayDenom)
	s.Require().Error(err)
}
