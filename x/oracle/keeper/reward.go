package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cheqd/cheqd-node/util"
	"github.com/cheqd/cheqd-node/util/decmath"
	"github.com/cheqd/cheqd-node/util/genmap"
	"github.com/cheqd/cheqd-node/util/reward"
	"github.com/cheqd/cheqd-node/x/oracle/types"
)

// prependOjoIfUnique pushes `uojo` denom to the front of the list, if it is not yet included.
func prependOjoIfUnique(voteTargets []string) []string {
	if genmap.Contains(types.CheqdDenom, voteTargets) {
		return voteTargets
	}
	rewardDenoms := make([]string, len(voteTargets)+1)
	rewardDenoms[0] = types.CheqdDenom
	copy(rewardDenoms[1:], voteTargets)
	return rewardDenoms
}

// smallestMissCountInBallot iterates through a given list of Claims and returns the smallest
// misscount in that list
func (k Keeper) smallestMissCountInBallot(ctx sdk.Context, ballotWinners []types.Claim) int64 {
	missCount := k.GetMissCounter(ctx, ballotWinners[0].Recipient)
	for _, winner := range ballotWinners[1:] {
		count := k.GetMissCounter(ctx, winner.Recipient)
		if count < missCount {
			missCount = count
		}
	}

	return util.SafeUint64ToInt64(missCount)
}

// RewardBallotWinners is executed at the end of every voting period, where we
// give out a portion of seigniorage reward(reward-weight) to the oracle voters
// that voted correctly.
func (k Keeper) RewardBallotWinners(
	ctx sdk.Context,
	votePeriod int64,
	rewardDistributionWindow int64,
	voteTargets []string,
	ballotWinners []types.Claim,
) {
	if len(ballotWinners) == 0 {
		return
	}

	distributionRatio := math.LegacyNewDec(votePeriod).QuoInt64(rewardDistributionWindow)
	var periodRewards sdk.DecCoins
	rewardDenoms := prependOjoIfUnique(voteTargets)
	for _, denom := range rewardDenoms {
		rewardPool := k.GetRewardPool(ctx, denom)

		// return if there's no rewards to give out
		if rewardPool.IsZero() {
			continue
		}

		periodRewards = periodRewards.Add(sdk.NewDecCoinFromDec(
			denom,
			math.LegacyNewDecFromInt(rewardPool.Amount).Mul(distributionRatio),
		))
	}

	// distribute rewards
	var distributedReward sdk.Coins

	smallestMissCount := k.smallestMissCountInBallot(ctx, ballotWinners)
	for _, winner := range ballotWinners {
		receiverVal, err := k.StakingKeeper.Validator(ctx, winner.Recipient)
		// in case absence of the validator, we just skip distribution
		if receiverVal == nil || err != nil {
			continue
		}

		missCount := util.SafeUint64ToInt64(k.GetMissCounter(ctx, winner.Recipient))
		maxMissCount := int64(len(voteTargets)) * (util.SafeUint64ToInt64((k.SlashWindow(ctx) / k.VotePeriod(ctx))))
		rewardFactor := reward.CalculateRewardFactor(
			missCount,
			maxMissCount,
			smallestMissCount,
		)
		rewardDec, err := decmath.NewDecFromFloat(rewardFactor)
		if err != nil {
			k.Logger(ctx).With(err).Error("unable to calculate validator reward factor!")
			return
		}
		ballotLength := int64(len(ballotWinners))

		rewardCoins, _ := periodRewards.MulDec(rewardDec.QuoInt64(
			ballotLength)).TruncateDecimal()
		if rewardCoins.IsZero() {
			continue
		}

		err = k.distrKeeper.AllocateTokensToValidator(ctx, receiverVal, sdk.NewDecCoinsFromCoins(rewardCoins...))
		if err != nil {
			k.Logger(ctx).With(err).Error("Failed to allocate tokens to validator!")
			return
		}
		distributedReward = distributedReward.Add(rewardCoins...)
	}

	// move distributed reward to distribution module
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.distrName, distributedReward)
	if err != nil {
		panic(fmt.Errorf("failed to send coins to distribution module %w", err))
	}
}
