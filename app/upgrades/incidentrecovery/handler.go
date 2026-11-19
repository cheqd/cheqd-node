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
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	vestingexported "github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
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
	GetAllUnbondingDelegations(ctx context.Context, delegator sdk.AccAddress) ([]stakingtypes.UnbondingDelegation, error)
	RemoveUnbondingDelegation(ctx context.Context, ubd stakingtypes.UnbondingDelegation) error
}

type DistrKeeper interface {
	// SetDelegatorWithdrawAddr (NOT the message-server-facing SetWithdrawAddr)
	// is called for each affected address before withdrawing any rewards or
	// calling Unbond, resetting its withdraw address back to itself. Unbond
	// itself withdraws rewards internally via the BeforeDelegationSharesModified
	// staking hook, paying out to whatever withdraw address is currently on
	// record — so this reset has to happen, and has to actually succeed,
	// before step 2 (delegation conversion) runs too, not just before the
	// explicit reward-withdrawal loop. SetWithdrawAddr (the public, checked
	// setter) refuses to run at all if the withdraw_addr_enabled chain
	// parameter happens to be false, or if the address is bank-blocked —
	// either of which would silently leave an attacker-set withdraw address
	// in place while this handler proceeded anyway. SetDelegatorWithdrawAddr
	// is the unchecked setter SetWithdrawAddr itself calls internally; using
	// it directly bypasses both of those gates, since what we're doing here
	// is an administrative recovery action, not a user-initiated message.
	SetDelegatorWithdrawAddr(ctx context.Context, delegatorAddr, withdrawAddr sdk.AccAddress) error
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

// blockTime returns the block time if ctx carries an sdk.Context. Callers
// must treat !ok as "unknown" and fail safe rather than assume a time.
func blockTime(ctx context.Context) (time.Time, bool) {
	sdkCtx, ok := ctx.Value(sdk.SdkContextKey).(sdk.Context)
	if !ok {
		return time.Time{}, false
	}
	return sdkCtx.BlockTime(), true
}

// RunIncidentRecovery is the core logic, called from the registered
// UpgradeHandler (see RegisterUpgradeHandler below). It is a plain function,
// not itself the upgradetypes.UpgradeHandler, so it can be unit tested
// directly without the full module-manager machinery.
//
// For each address in cfg.AffectedAddresses, in order:
//  0. Reset its distribution withdraw address back to itself, so step 1
//     can't pay rewards out to an address the attacker redirected them to.
//  1. Withdraw any pending distribution rewards into the liquid balance
//     (skipped if step 0 failed).
//  2. Force-convert any active delegations to liquid balance immediately
//     (bypassing the normal 14-day unbonding period entirely) — skipped
//     for a vesting account with coins still vesting, see the isVesting
//     comment below.
//  3. Force-complete any already-in-progress unbonding delegation entries
//     immediately (bypassing both the remaining wait AND, critically, the
//     jailed-validator check that MsgCancelUnbondingDelegation enforces —
//     this is a direct state operation, not that message handler, so it
//     is not subject to that check) — also skipped for a still-vesting account.
//  4. Sweep the resulting spendable balance to cfg.DestinationAddress.
//
// Any failure in step 2 or 3 after the destructive Unbond/pool-debit call
// has already run aborts the ENTIRE run (not just that one address) with an
// error, so the upgrade block fails and nothing commits — see the comments
// at each such call site for why a "skip and continue" there would have
// been actively unsafe.
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
	// Re-validate even though ParsePlanInfo already did — this is the point of
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

		// Dynamic module-account check, same as the destination check
		// above: config.go's static affected-address check only knows the
		// standard SDK module names, so this catches any module account
		// (including any cheqd-specific one) a typo'd or malicious config
		// entry might have slipped past it.
		if acc := accountKeeper.GetAccount(ctx, delAddr); acc != nil {
			if _, ok := acc.(sdk.ModuleAccountI); ok {
				result.Skipped = true
				result.SkipReason = "address resolves to a module account on-chain, refusing to treat it as a compromised wallet"
				results = append(results, result)
				continue
			}
		}

		// Reset the withdraw address to the delegator itself before
		// withdrawing anything OR calling Unbond. If the attacker already
		// called MsgSetWithdrawAddress to point it somewhere they control,
		// withdrawing rewards without this first would hand those rewards
		// straight to them while the on-chain event still logs it as
		// "recovered" — and Unbond (step 2 below) withdraws rewards
		// internally via a staking hook, paying out to whatever withdraw
		// address is on record, regardless of whether this handler also
		// withdraws rewards explicitly. Not fatal if it fails — nothing
		// destructive has happened yet — but BOTH explicit reward
		// withdrawal AND delegation conversion are skipped for this address
		// rather than risk paying rewards out to an attacker-controlled
		// address via either path.
		canWithdrawRewardsSafely := true
		if err := distrKeeper.SetDelegatorWithdrawAddr(ctx, delAddr, delAddr); err != nil {
			canWithdrawRewardsSafely = false
			result.warn("could not reset withdraw address to the delegator itself, skipping reward withdrawal and delegation conversion for this address to avoid paying rewards out to a possibly attacker-controlled address: %v", err)
		}

		// A vesting account's delegated-vesting coins, once force-converted
		// by steps 2/3 below, come back as LOCKED balance (per the vesting
		// accounting the bank keeper runs on undelegation), not spendable —
		// and locked balance unlocks and becomes spendable by the
		// compromised account itself over the vesting schedule, with none
		// of the normal 21-day unbonding wait. Force-converting a vesting
		// account's stake would make this handler actively help the
		// attacker rather than recover funds, so active delegations and
		// unbonding entries are deliberately left untouched for one; only
		// its already-spendable balance (step 4) and rewards (above) are
		// swept.
		//
		// The check is on coins still vesting at block time, not the account
		// type: a vesting account whose schedule has ended keeps its type
		// forever but has nothing left to release, so it is handled like any
		// other account. LockedCoins is deliberately not used here — it nets
		// off DelegatedVesting, so a mid-schedule account with its vesting
		// coins delegated would report zero locked and be wrongly converted.
		isVesting := false
		if acc := accountKeeper.GetAccount(ctx, delAddr); acc != nil {
			if va, ok := acc.(vestingexported.VestingAccount); ok {
				now, known := blockTime(ctx)
				if !known || !va.GetVestingCoins(now).IsZero() {
					isVesting = true
					result.warn("address is a vesting account with coins still vesting: active delegations and unbonding entries were left untouched and require separate, manual handling — force-converting them would release locked-but-unvested coins back to the compromised account with no unbonding wait")
				} else {
					result.warn("address is a vesting account whose schedule has fully ended (nothing still vesting at block time %s); handled as an ordinary account", now.UTC().Format(time.RFC3339))
				}
			}
		}

		// 1. Withdraw pending rewards from every active delegation.
		delegations, err := stakingKeeper.GetAllDelegatorDelegations(ctx, delAddr)
		if err != nil {
			result.Skipped = true
			result.SkipReason = fmt.Sprintf("listing delegations failed: %v", err)
			results = append(results, result)
			continue
		}

		if canWithdrawRewardsSafely {
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
		}

		//nolint:nestif // fund-moving steps are kept inline, in order, so each error path is auditable in one place
		if !isVesting && canWithdrawRewardsSafely {
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
					// Unbond is not atomic: internally it can reduce or
					// remove the delegation's shares (SetDelegation /
					// RemoveDelegation) and only THEN call
					// RemoveValidatorTokensAndShares, possibly
					// RemoveValidator, and the AfterDelegationModified hook —
					// any of which can fail after the shares have already
					// been written to the store. An error here does not mean
					// "nothing happened," the same reasoning as the credit
					// failure below. Aborting the whole run, rather than
					// warning and continuing, means the block fails and
					// nothing commits either way.
					return nil, fmt.Errorf("%s: delegation conversion: Unbond from validator %s returned an error, but Unbond is not atomic and may have already mutated delegation/validator state before failing; aborting entire recovery run rather than risk committing that partial mutation: %w", addrStr, d.ValidatorAddress, err)
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
		}

		//nolint:nestif // fund-moving steps are kept inline, in order, so each error path is auditable in one place
		if !isVesting {
			// 3. Force-complete in-progress unbonding delegation entries
			//    immediately. This is the step that recovers stake stuck behind
			//    a jailed validator — it never calls MsgCancelUnbondingDelegation
			//    or its jailed-validator check at all.
			unbondings, err := stakingKeeper.GetAllUnbondingDelegations(ctx, delAddr)
			if err != nil {
				result.warn("listing unbonding delegations failed: %v", err)
			} else {
				for _, ubd := range unbondings {
					for _, entry := range ubd.Entries {
						if !entry.Balance.IsPositive() {
							continue
						}
						coins := sdk.NewCoins(sdk.NewCoin(bondDenom, entry.Balance))
						if err := bankKeeper.UndelegateCoinsFromModuleToAccount(ctx, stakingtypes.NotBondedPoolName, delAddr, coins); err != nil {
							// UndelegateCoinsFromModuleToAccount debits the
							// NotBonded pool BEFORE crediting the account
							// internally — a failure here does NOT mean
							// "nothing moved". Treating it as a skippable
							// warning (the original bug) would let the pool
							// debit commit while the matching record is kept
							// as if nothing happened, letting a later
							// CompleteUnbonding pay this entry out of OTHER
							// delegators' funds in that pool. Abort the whole
							// run instead, same as the Unbond case above.
							return nil, fmt.Errorf("%s: unbonding completion: crediting account from validator %s (creation height %d) failed, aborting entire recovery run to avoid committing a pool debit with no matching credit: %w", addrStr, ubd.ValidatorAddress, entry.CreationHeight, err)
						}
						result.UnbondingMoved = result.UnbondingMoved.Add(coins...)
					}
					if err := stakingKeeper.RemoveUnbondingDelegation(ctx, ubd); err != nil {
						// Tokens have already moved to the account at this
						// point, but the record still holds entries with a
						// positive balance. If it isn't removed, a later
						// EndBlocker CompleteUnbonding call for this exact
						// record will see those same entries as still mature
						// and pay their balance out AGAIN from the
						// NotBonded pool — a real double payout funded by
						// other delegators' unbonding funds, not just a
						// bookkeeping inconsistency. Abort the whole run
						// instead, same reasoning as the credit failure above.
						return nil, fmt.Errorf("%s: unbonding completion: tokens already credited from validator %s but RemoveUnbondingDelegation failed, aborting entire recovery run to avoid leaving a stale unbonding record that a later EndBlocker would pay out a second time: %w", addrStr, ubd.ValidatorAddress, err)
					}
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
