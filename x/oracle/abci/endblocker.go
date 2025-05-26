package abci

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cheqd/cheqd-node/util"
	"github.com/cheqd/cheqd-node/x/oracle/keeper"
	"github.com/cheqd/cheqd-node/x/oracle/types"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx context.Context, k keeper.Keeper) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Check for Oracle parameter update plans and execute plans that are
	// at their plan height.
	plans := k.GetParamUpdatePlans(sdkCtx)
	for _, plan := range plans {
		if plan.ShouldExecute(sdkCtx) {
			if err := k.ExecuteParamUpdatePlan(sdkCtx, plan); err != nil {
				return err
			}
		}
	}

	params := k.GetParams(sdkCtx)

	// Start price feeder if it hasn't been started, and it is enabled.
	if k.PriceFeeder.Oracle == nil && k.PriceFeeder.AppConfig.Enable {
		go func() {
			err := k.PriceFeeder.Start(sdkCtx.BlockHeight(), params)
			if err != nil {
				sdkCtx.Logger().Error("Error starting Oracle Keeper price feeder", "err", err)
			}
		}()
	}

	// Set all current active validators into the ValidatorRewardSet at
	// the beginning of a new Slash Window.
	if k.IsPeriodLastBlock(sdkCtx, params.SlashWindow+1) {
		if err := k.SetValidatorRewardSet(sdkCtx); err != nil {
			return err
		}
	}

	if k.IsPeriodLastBlock(sdkCtx, params.VotePeriod) {
		RunPriceFeederIfEnabled(ctx, sdkCtx, params, k)

		if err := UpdateOraclePrices(sdkCtx, params, k); err != nil {
			return err
		}

		if k.IsPeriodLastBlock(sdkCtx, params.HistoricStampPeriod*keeper.AveragingWindow) {
			if err := ComputeAllAverages(sdkCtx, params, k); err != nil {
				return err
			}
		}
	}

	// Slash oracle providers who missed voting over the threshold and reset
	// miss counters of all validators at the last block of slash window.
	if k.IsPeriodLastBlock(sdkCtx, params.SlashWindow) {
		k.SlashAndResetMissCounters(sdkCtx)
	}
	k.PruneAllPrices(sdkCtx)
	return nil
}

func RunPriceFeederIfEnabled(ctx context.Context, sdkCtx sdk.Context, params types.Params, k keeper.Keeper) {
	if k.PriceFeeder.Oracle != nil && k.PriceFeeder.AppConfig.Enable {
		k.PriceFeeder.Oracle.ParamCache.UpdateParamCache(
			sdkCtx.BlockHeight(),
			k.GetParams(sdkCtx),
			nil,
		)

		if err := k.PriceFeeder.Oracle.TickClientless(ctx); err != nil {
			sdkCtx.Logger().Error("Error in Oracle Keeper price feeder clientless tick", "err", err)
		}
	}
}

func UpdateOraclePrices(sdkCtx sdk.Context, params types.Params, k keeper.Keeper) error {
	return CalcPrices(sdkCtx, params, k)
}

func ComputeAllAverages(sdkCtx sdk.Context, params types.Params, k keeper.Keeper) error {
	for _, v := range params.AcceptList {
		if err := k.ComputeAverages(sdkCtx, v.SymbolDenom); err != nil {
			return err
		}
	}
	return nil
}

func CalcPrices(ctx sdk.Context, params types.Params, k keeper.Keeper) error {
	// Build claim map over all validators in active set
	validatorClaimMap := make(map[string]types.Claim)
	powerReduction := k.StakingKeeper.PowerReduction(ctx)
	// Calculate total validator power
	var totalBondedPower int64
	vals, err := k.StakingKeeper.GetBondedValidatorsByPower(ctx)
	if err != nil {
		return err
	}
	for _, v := range vals {
		addrString := v.GetOperator()
		addr, err := sdk.ValAddressFromBech32(addrString)
		if err != nil {
			return err
		}
		power := v.GetConsensusPower(powerReduction)
		totalBondedPower += power
		validatorClaimMap[addrString] = types.NewClaim(power, 0, 0, addr)
	}

	// voteTargets defines the symbol (ticker) denoms that we require votes on
	voteTargetDenoms := make([]string, 0)
	for _, v := range params.AcceptList {
		voteTargetDenoms = append(voteTargetDenoms, v.BaseDenom)
	}

	k.ClearExchangeRates(ctx)

	// NOTE: it filters out inactive or jailed validators
	ballotDenomSlice := k.OrganizeBallotByDenom(ctx, validatorClaimMap)
	threshold := k.VoteThreshold(ctx).MulInt64(types.MaxVoteThresholdMultiplier).TruncateInt64()

	// Iterate through ballots and update exchange rates; drop if not enough votes have been achieved.
	for _, ballotDenom := range ballotDenomSlice {
		// Increment Mandatory Win count if Denom in Mandatory list
		incrementWin := params.MandatoryList.Contains(ballotDenom.Denom)

		// If the asset is not in the mandatory or accept list, continue
		if !incrementWin && !params.AcceptList.Contains(ballotDenom.Denom) {
			ctx.Logger().Info("Unsupported denom, dropping ballot", "denom", ballotDenom)
			continue
		}

		// Calculate the portion of votes received as an integer, scaled up using the
		// same multiplier as the `threshold` computed above
		support := ballotDenom.Ballot.Power() * types.MaxVoteThresholdMultiplier / totalBondedPower
		if support < threshold {
			ctx.Logger().Info("Ballot voting power is under vote threshold, dropping ballot", "denom", ballotDenom)
			continue
		}

		// Get the current denom's reward band
		rewardBand, err := params.RewardBands.GetBandFromDenom(ballotDenom.Denom)
		if err != nil {
			return err
		}

		// Get weighted median of exchange rates
		exchangeRate, err := Tally(ballotDenom.Ballot, rewardBand, validatorClaimMap, incrementWin)
		if err != nil {
			return err
		}

		// Set the exchange rate, emit ABCI event
		if err = k.SetExchangeRateWithEvent(ctx, ballotDenom.Denom, exchangeRate); err != nil {
			return err
		}

		if k.IsPeriodLastBlock(ctx, params.HistoricStampPeriod) {
			fmt.Println(">>>>>>>>>>>>>>>>>>>.ballotDenom", ballotDenom.Denom, exchangeRate)
			k.AddHistoricPrice(ctx, ballotDenom.Denom, exchangeRate)
		}
		// Calculate and stamp median/median deviation if median stamp period has passed
		if k.IsPeriodLastBlock(ctx, params.MedianStampPeriod) {
			if err = k.CalcAndSetHistoricMedian(ctx, ballotDenom.Denom); err != nil {
				return err
			}
		}
	}

	// Get the validators which can earn rewards in this Slash Window.
	validatorRewardSet := k.GetValidatorRewardSet(ctx)

	// update miss counting & slashing
	voteTargetsLen := len(params.MandatoryList)
	claimSlice, rewardSlice := types.ClaimMapToSlices(validatorClaimMap, validatorRewardSet.ValidatorSet)
	for _, claim := range claimSlice {
		misses := util.SafeIntToUint64(voteTargetsLen - int(claim.MandatoryWinCount))
		if misses == 0 {
			continue
		}

		// Increase miss counter
		k.SetMissCounter(ctx, claim.Recipient, k.GetMissCounter(ctx, claim.Recipient)+misses)
	}

	// Distribute rewards to ballot winners
	k.RewardBallotWinners(
		ctx,
		util.SafeUint64ToInt64(params.VotePeriod),
		util.SafeUint64ToInt64(params.RewardDistributionWindow),
		voteTargetDenoms,
		rewardSlice,
	)

	// Clear the ballot
	k.ClearBallots(ctx, params.VotePeriod)
	return nil
}

// Tally calculates and returns the median. It sets the set of voters to be
// rewarded, i.e. voted within a reasonable spread from the weighted median to
// the store. Note, the ballot is sorted by ExchangeRate.
func Tally(
	ballot types.ExchangeRateBallot,
	rewardBand math.LegacyDec,
	validatorClaimMap map[string]types.Claim,
	incrementWin bool,
) (math.LegacyDec, error) {
	weightedMedian, err := ballot.WeightedMedian()
	if err != nil {
		return math.LegacyZeroDec(), err
	}
	standardDeviation, err := ballot.StandardDeviation()
	if err != nil {
		return math.LegacyZeroDec(), err
	}

	// rewardSpread is the MAX((weightedMedian * (rewardBand/2)), standardDeviation)
	rewardSpread := weightedMedian.Mul(rewardBand.QuoInt64(2))
	rewardSpread = math.LegacyMaxDec(rewardSpread, standardDeviation)

	for _, tallyVote := range ballot {
		// Filter ballot winners. For voters, we filter out the tally vote iff:
		// (weightedMedian - rewardSpread) <= ExchangeRate <= (weightedMedian + rewardSpread)
		if (tallyVote.ExchangeRate.GTE(weightedMedian.Sub(rewardSpread)) &&
			tallyVote.ExchangeRate.LTE(weightedMedian.Add(rewardSpread))) ||
			!tallyVote.ExchangeRate.IsPositive() {

			key := tallyVote.Voter.String()
			claim := validatorClaimMap[key]

			if incrementWin {
				claim.MandatoryWinCount++
			}
			claim.Weight += tallyVote.Power
			validatorClaimMap[key] = claim
		}
	}

	return weightedMedian, nil
}
