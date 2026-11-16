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

// A step that has already irreversibly mutated chain state (Unbond removing
// shares, or a final sweep leaving liquidated funds spendable) must abort
// the entire run with an error, not just skip or warn, if its matching
// credit/send step then fails. Returned as a plain error from
// RunIncidentRecovery, which propagates up through the registered upgrade
// handler and causes the node to halt without committing the block — the
// same effect as panicking, but via a normal error return. Everything
// already applied in this block is discarded along with it, so no
// inconsistent state (shares destroyed but tokens unaccounted for) can ever
// be committed on chain; a corrected binary has to ship before the upgrade
// can proceed.

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
	// SpendableCoins excludes locked/vesting coins. The final sweep must use
	// this, not GetAllBalances, since SendCoins itself refuses to move
	// locked coins — using the full balance there would make the final
	// sweep fail outright for any affected address that happens to be a
	// vesting account.
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
	UndelegateCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

type StakingKeeper interface {
	BondDenom(ctx context.Context) (string, error)
	// GetAllDelegatorDelegations, not GetDelegatorDelegations with a
	// maxRetrieve of 0 — the real keeper treats maxRetrieve as a cap on how
	// many delegations to return, so 0 always returns an empty list. An
	// earlier draft of this handler passed 0 here, which silently skipped
	// reward withdrawal and active-delegation sweeping entirely.
	GetAllDelegatorDelegations(ctx context.Context, delegator sdk.AccAddress) ([]stakingtypes.Delegation, error)
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

// AccountKeeper is used only to check, at execution time, whether the
// configured destination address actually resolves to a module account on
// this specific chain — a dynamic, always-complete check, rather than
// relying solely on the hardcoded name list in config.go's
// knownModuleAccountNames (which only covers the standard SDK modules and
// would miss any cheqd-specific or newer module this binary registers).
type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
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
	accountKeeper AccountKeeper,
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

	// Dynamic check against the real account state, on top of config.go's
	// static name-list check — catches any module account the hardcoded
	// list doesn't know about.
	if destAcc := accountKeeper.GetAccount(ctx, destAddr); destAcc != nil {
		if _, ok := destAcc.(sdk.ModuleAccountI); ok {
			return nil, fmt.Errorf("refusing to run: destination_address %q resolves to a module account on-chain", cfg.DestinationAddress)
		}
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
		delegations, err := stakingKeeper.GetAllDelegatorDelegations(ctx, delAddr)
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
				// Unbond has ALREADY removed the delegation's shares at
				// this point — there is no way to put them back. Letting
				// the loop continue would commit that destructive change
				// while the matching tokens never reach any account.
				// Aborting the whole run here means the block fails and
				// nothing in it commits, so the shares are never actually
				// lost: the chain halts, instead of silently losing funds.
				return nil, fmt.Errorf("%s: delegation conversion: Unbond succeeded (removed %s shares, %s tokens) from validator %s but crediting the account failed, aborting entire recovery run to avoid committing unaccounted-for tokens: %w", addrStr, d.Shares.String(), coins.String(), d.ValidatorAddress, err)
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
				// failedEntries keeps every entry whose credit transfer
				// failed, so the unbonding record below can be trimmed down
				// to just those rather than deleted wholesale. The transfer
				// itself is a single bank call per entry — if it errors,
				// nothing has moved for that entry yet, so leaving it in
				// place (instead of treating it as moved) is safe, not just
				// a cosmetic bookkeeping choice.
				var failedEntries []stakingtypes.UnbondingDelegationEntry
				for _, entry := range ubd.Entries {
					if !entry.Balance.IsPositive() {
						continue
					}
					coins := sdk.NewCoins(sdk.NewCoin(bondDenom, entry.Balance))
					if err := bankKeeper.UndelegateCoinsFromModuleToAccount(ctx, stakingtypes.NotBondedPoolName, delAddr, coins); err != nil {
						result.warn("unbonding completion: crediting account from validator %s (creation height %d) failed, leaving this entry in place for a later run: %v", ubd.ValidatorAddress, entry.CreationHeight, err)
						failedEntries = append(failedEntries, entry)
						continue
					}
					moved = moved.Add(coins...)
				}
				if moved == nil {
					// Nothing moved for this record at all — leave it
					// completely untouched.
					continue
				}
				result.UnbondingMoved = result.UnbondingMoved.Add(moved...)
				if len(failedEntries) == 0 {
					if err := stakingKeeper.RemoveUnbondingDelegation(ctx, ubd); err != nil {
						// Tokens have already moved to the account at this point.
						// Not removing the (now-empty) unbonding record is a
						// bookkeeping inconsistency, not a fund-safety issue,
						// but it should be investigated if it happens.
						result.warn("unbonding completion: tokens moved from validator %s but RemoveUnbondingDelegation failed: %v", ubd.ValidatorAddress, err)
					}
					continue
				}
				// Some, but not all, entries succeeded: keep the record but
				// shrink it to only the entries that are still genuinely
				// unbonding. Removing the whole record here (the original
				// bug) would permanently strand the failed entries' tokens
				// in the NotBonded pool with no record pointing to them.
				ubd.Entries = failedEntries
				if err := stakingKeeper.SetUnbondingDelegation(ctx, ubd); err != nil {
					result.warn("unbonding completion: partially moved tokens from validator %s but failed to persist the trimmed unbonding record (%d entries remaining): %v", ubd.ValidatorAddress, len(failedEntries), err)
				}
			}
		}

		// 4. Sweep the resulting spendable balance. Deliberately
		// SpendableCoins, not GetAllBalances: SendCoins itself refuses to
		// move locked coins, so sweeping the full balance would make this
		// final step fail outright for any address holding a vesting
		// account. Any locked remainder is flagged below, not silently
		// dropped — it is a known limitation (vesting lockup requires
		// separate handling, out of scope here), not a bug to hide.
		fullBalance := bankKeeper.GetAllBalances(ctx, delAddr)
		spendable := bankKeeper.SpendableCoins(ctx, delAddr)
		if locked := fullBalance.Sub(spendable...); !locked.IsZero() {
			result.warn("address holds %s in locked/vesting balance that this handler cannot move via SendCoins and did NOT sweep — needs separate handling", locked.String())
		}

		if spendable.IsZero() {
			result.Skipped = true
			if result.SkipReason == "" {
				result.SkipReason = "zero spendable balance at execution time — nothing to sweep"
			}
			results = append(results, result)
			continue
		}

		if err := bankKeeper.SendCoins(ctx, delAddr, destAddr, spendable); err != nil {
			// By this point steps 1-3 may already have liquidated rewards,
			// delegations and unbonding entries into this address's
			// spendable balance. Skipping and moving on, as the original
			// code did, would leave that freshly-liquidated balance sitting
			// immediately spendable in the compromised account — worse
			// than doing nothing. Abort the whole run instead, so the
			// block fails and none of this commits.
			return nil, fmt.Errorf("%s: final sweep SendCoins of %s failed, aborting entire recovery run to avoid leaving liquidated funds spendable in the compromised account: %w", addrStr, spendable.String(), err)
		}

		result.TotalSwept = spendable
		results = append(results, result)
	}

	return results, nil
}
