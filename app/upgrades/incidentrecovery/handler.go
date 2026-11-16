// Package incidentrecovery implements the upgrade handler for the
// 2026-11-13 compromised-wallet incident. Integrated against cheqd-node's
// real app.go and its actual vendored cosmos-sdk keepers (not mocks) — see
// RegisterUpgradeHandler in upgrade.go for the wiring.
//
// Status: integrated, not yet run against any devnet/testnet. Do not run
// against mainnet until the full TESTING_PLAN.md rollout has passed.
package incidentrecovery

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// The minimal keeper surface this handler needs. Kept as local interfaces
// (rather than depending on the app's concrete keeper types directly) so
// handler_test.go can exercise this logic with hand-written mocks instead
// of standing up a full app — Go satisfies these implicitly, so passing
// the app's real *stakingkeeper.Keeper / bankkeeper.Keeper / distrkeeper.Keeper
// from upgrade.go works with no adapter code, PROVIDED the signatures here
// exactly match the real keeper methods.
//
// Every signature below has been checked directly against cheqd-node's
// vendored cosmos-sdk fork (github.com/cheqd/cosmos-sdk v0.50.14-height-
// mismatch-iavl...), not assumed from general cosmos-sdk familiarity — an
// earlier draft of this file had three signatures that were wrong in ways
// that would have failed to compile against the real app:
//   - DistrKeeper's real method is WithdrawDelegationRewards (plural
//     "Rewards"), not WithdrawDelegatorReward (singular) — that name only
//     exists on the internal message-server wrapper, taking a Msg struct.
//   - StakingKeeper.GetDelegatorUnbondingDelegations doesn't exist at all;
//     the real method is GetAllUnbondingDelegations.
//   - StakingKeeper.Unbond returns (math.Int, error), not
//     (math.LegacyDec, error) — it's already the final token amount, no
//     TruncateInt() needed or possible (math.Int has no such method).

type BankKeeper interface {
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
	UndelegateCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

type StakingKeeper interface {
	BondDenom(ctx context.Context) (string, error)
	GetDelegatorDelegations(ctx context.Context, delegator sdk.AccAddress, maxRetrieve uint16) ([]stakingtypes.Delegation, error)
	GetValidator(ctx context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error)
	// Unbond converts shares to tokens and removes the delegation record.
	// It does NOT move the underlying tokens to the delegator's account —
	// the caller is responsible for that via UndelegateCoinsFromModuleToAccount,
	// same as the real Undelegate message handler does internally.
	Unbond(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares math.LegacyDec) (math.Int, error)
	GetUnbondingDelegation(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.UnbondingDelegation, error)
	GetAllUnbondingDelegations(ctx context.Context, delegator sdk.AccAddress) ([]stakingtypes.UnbondingDelegation, error)
	RemoveUnbondingDelegation(ctx context.Context, ubd stakingtypes.UnbondingDelegation) error
	SetUnbondingDelegation(ctx context.Context, ubd stakingtypes.UnbondingDelegation) error
}

type DistrKeeper interface {
	WithdrawDelegationRewards(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdk.Coins, error)
}

// SweepResult records what actually happened for one affected address, for
// the on-chain event log and for post-execution auditing by anyone
// reviewing the upgrade afterwards.
type SweepResult struct {
	Address          string
	RewardsWithdrawn sdk.Coins
	UnbondingMoved   sdk.Coins
	DelegationsMoved sdk.Coins
	TotalSwept       sdk.Coins
	Skipped          bool
	SkipReason       string
	// Warnings records every per-step failure that was *not* fatal to the
	// overall sweep for this address (e.g. one validator's reward
	// withdrawal or delegation conversion failing while others succeed).
	// Deliberately never silently dropped: a partial failure here means
	// shares/tokens could be in an inconsistent state for that one
	// validator, and an auditor comparing swept totals against expected
	// amounts needs to know exactly where and why a shortfall happened,
	// not just that the overall address "succeeded".
	Warnings []string
}

func (r *SweepResult) warn(format string, args ...any) {
	r.Warnings = append(r.Warnings, fmt.Sprintf(format, args...))
}

// RunIncidentRecovery is the core logic, called from the registered
// UpgradeHandler (see RegisterUpgradeHandler below). It is a plain function,
// not itself the upgradetypes.UpgradeHandler, so it can be unit tested
// directly without the full module-manager machinery.
//
// For each address in cfg.AffectedAddresses, in order:
//  1. Withdraw any pending distribution rewards into the liquid balance.
//  2. Force-convert any active delegations to liquid balance immediately
//     (bypassing the normal 14-day unbonding period entirely).
//  3. Force-complete any already-in-progress unbonding delegation entries
//     immediately (bypassing both the remaining wait AND, critically, the
//     jailed-validator check that MsgCancelUnbondingDelegation enforces —
//     this is a direct state operation, not that message handler, so it
//     is not subject to that check).
//  4. Sweep the resulting full liquid balance to cfg.DestinationAddress.
//
// Reads live on-chain state at execution time throughout — it does not
// trust or replay any amounts recorded when this config was drafted, so it
// is correct even if balances have moved (partially or fully) by the time
// this actually runs.
func RunIncidentRecovery(
	ctx context.Context,
	cfg *Config,
	bankKeeper BankKeeper,
	stakingKeeper StakingKeeper,
	distrKeeper DistrKeeper,
) ([]SweepResult, error) {
	// Re-validate even though LoadConfig already did — this is the point of
	// actual execution, and a second independent check here costs nothing.
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("refusing to run: config failed validation at execution time: %w", err)
	}

	destAddr, err := sdk.AccAddressFromBech32(cfg.DestinationAddress)
	if err != nil {
		return nil, fmt.Errorf("refusing to run: destination_address did not parse as a valid address at execution time: %w", err)
	}

	bondDenom, err := stakingKeeper.BondDenom(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting bond denom: %w", err)
	}

	results := make([]SweepResult, 0, len(cfg.AffectedAddresses))

	for _, addrStr := range cfg.AffectedAddresses {
		result := SweepResult{Address: addrStr}

		delAddr, err := sdk.AccAddressFromBech32(addrStr)
		if err != nil {
			result.Skipped = true
			result.SkipReason = fmt.Sprintf("address did not parse: %v", err)
			results = append(results, result)
			continue
		}

		// 1. Withdraw pending rewards from every active delegation.
		delegations, err := stakingKeeper.GetDelegatorDelegations(ctx, delAddr, 0)
		if err != nil {
			result.Skipped = true
			result.SkipReason = fmt.Sprintf("listing delegations failed: %v", err)
			results = append(results, result)
			continue
		}

		for _, d := range delegations {
			valAddr, err := sdk.ValAddressFromBech32(d.ValidatorAddress)
			if err != nil {
				result.warn("reward withdrawal: validator address %q did not parse: %v", d.ValidatorAddress, err)
				continue
			}
			rewards, err := distrKeeper.WithdrawDelegationRewards(ctx, delAddr, valAddr)
			if err != nil {
				// Deliberately not failing the whole sweep for this address on
				// a reward-withdrawal error — the bulk of the value is
				// principal, not rewards — but it IS recorded, not dropped.
				result.warn("reward withdrawal from validator %s failed: %v", d.ValidatorAddress, err)
				continue
			}
			result.RewardsWithdrawn = result.RewardsWithdrawn.Add(rewards...)
		}

		// 2. Force-convert active delegations to liquid balance immediately.
		for _, d := range delegations {
			valAddr, err := sdk.ValAddressFromBech32(d.ValidatorAddress)
			if err != nil {
				result.warn("delegation conversion: validator address %q did not parse: %v", d.ValidatorAddress, err)
				continue
			}
			validator, err := stakingKeeper.GetValidator(ctx, valAddr)
			if err != nil {
				result.warn("delegation conversion: could not load validator %s: %v", d.ValidatorAddress, err)
				continue
			}
			amount, err := stakingKeeper.Unbond(ctx, delAddr, valAddr, d.Shares)
			if err != nil {
				result.warn("delegation conversion: Unbond from validator %s failed: %v", d.ValidatorAddress, err)
				continue
			}
			if !amount.IsPositive() {
				// Not a warning case: shares converted to zero tokens is a
				// legitimate (if unusual) outcome, not a failure.
				continue
			}
			coins := sdk.NewCoins(sdk.NewCoin(bondDenom, amount))
			sourceModule := stakingtypes.BondedPoolName
			if !validator.IsBonded() {
				sourceModule = stakingtypes.NotBondedPoolName
			}
			if err := bankKeeper.UndelegateCoinsFromModuleToAccount(ctx, sourceModule, delAddr, coins); err != nil {
				// This is the dangerous case: Unbond already removed the
				// delegation's shares above, but the tokens never reached
				// the account. That value is not accounted for anywhere
				// else — recording it loudly is the only way an auditor
				// can know it happened at all.
				result.warn("delegation conversion: Unbond succeeded (removed %s shares, %s tokens) from validator %s but crediting the account failed — TOKENS MAY BE UNACCOUNTED FOR: %v", d.Shares.String(), coins.String(), d.ValidatorAddress, err)
				continue
			}
			result.DelegationsMoved = result.DelegationsMoved.Add(coins...)
		}

		// 3. Force-complete in-progress unbonding delegation entries
		//    immediately. This is the step that recovers stake stuck behind
		//    a jailed validator — it never calls MsgCancelUnbondingDelegation
		//    or its jailed-validator check at all.
		unbondings, err := stakingKeeper.GetAllUnbondingDelegations(ctx, delAddr)
		if err != nil {
			result.warn("listing unbonding delegations failed: %v", err)
		} else {
			for _, ubd := range unbondings {
				var moved sdk.Coins
				for _, entry := range ubd.Entries {
					if !entry.Balance.IsPositive() {
						continue
					}
					coins := sdk.NewCoins(sdk.NewCoin(bondDenom, entry.Balance))
					if err := bankKeeper.UndelegateCoinsFromModuleToAccount(ctx, stakingtypes.NotBondedPoolName, delAddr, coins); err != nil {
						result.warn("unbonding completion: crediting account from validator %s (creation height %d) failed — TOKENS MAY BE UNACCOUNTED FOR: %v", ubd.ValidatorAddress, entry.CreationHeight, err)
						continue
					}
					moved = moved.Add(coins...)
				}
				if moved != nil {
					result.UnbondingMoved = result.UnbondingMoved.Add(moved...)
					if err := stakingKeeper.RemoveUnbondingDelegation(ctx, ubd); err != nil {
						// Tokens have already moved to the account at this point.
						// Not removing the (now-empty) unbonding record is a
						// bookkeeping inconsistency, not a fund-safety issue,
						// but it should be investigated if it happens.
						result.warn("unbonding completion: tokens moved from validator %s but RemoveUnbondingDelegation failed: %v", ubd.ValidatorAddress, err)
					}
				}
			}
		}

		// 4. Sweep the full resulting liquid balance.
		balance := bankKeeper.GetAllBalances(ctx, delAddr)
		if balance.IsZero() {
			result.Skipped = true
			if result.SkipReason == "" {
				result.SkipReason = "zero balance at execution time — nothing to sweep"
			}
			results = append(results, result)
			continue
		}

		if err := bankKeeper.SendCoins(ctx, delAddr, destAddr, balance); err != nil {
			result.Skipped = true
			result.SkipReason = fmt.Sprintf("final sweep SendCoins failed: %v", err)
			results = append(results, result)
			continue
		}

		result.TotalSwept = balance
		results = append(results, result)
	}

	return results, nil
}
