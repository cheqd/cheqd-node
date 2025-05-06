package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cheqd/cheqd-node/util"
	"github.com/cheqd/cheqd-node/x/oracle/types"
)

// SlashAndResetMissCounters iterates over all the current missed counters and
// calculates the "valid vote rate" as:
// (possibleWinsPerSlashWindow - missCounter)/possibleWinsPerSlashWindow.
//
// If the valid vote rate is below the minValidPerWindow, the validator will be
// slashed and jailed.
func (k Keeper) SlashAndResetMissCounters(ctx sdk.Context) {
	var (
		possibleWinsPerSlashWindow = k.PossibleWinsPerSlashWindow(ctx)
		minValidPerWindow          = k.MinValidPerWindow(ctx)

		distributionHeight = ctx.BlockHeight() - sdk.ValidatorUpdateDelay - 1
		slashFraction      = k.SlashFraction(ctx)
		powerReduction     = k.StakingKeeper.PowerReduction(ctx)
	)

	k.IterateMissCounters(ctx, func(operator sdk.ValAddress, missCounter uint64) bool {
		validVotes := math.NewInt(possibleWinsPerSlashWindow - util.SafeUint64ToInt64(missCounter))
		validVoteRate := math.LegacyNewDecFromInt(validVotes).QuoInt64(possibleWinsPerSlashWindow)

		// Slash and jail the validator if their valid vote rate is smaller than the
		// minimum threshold.
		if validVoteRate.LT(minValidPerWindow) {
			validator, err := k.StakingKeeper.Validator(ctx, operator)
			if validator.IsBonded() && !validator.IsJailed() && err == nil {
				consAddr, err := validator.GetConsAddr()
				if err != nil {
					panic(err)
				}

				_, err = k.StakingKeeper.Slash(
					ctx,
					consAddr,
					distributionHeight,
					validator.GetConsensusPower(powerReduction), slashFraction,
				)
				if err != nil {
					panic(err)
				}

				err = k.StakingKeeper.Jail(ctx, consAddr)
				if err != nil {
					panic(err)
				}
			}
		}

		k.DeleteMissCounter(ctx, operator)
		return false
	})
}

// PossibleWinsPerSlashWindow returns the total number of possible correct votes
// that a validator can have per asset multiplied by the number of vote
// periods in the slash window
func (k Keeper) PossibleWinsPerSlashWindow(ctx sdk.Context) int64 {
	slashWindow := util.SafeUint64ToInt64(k.SlashWindow(ctx))
	votePeriod := util.SafeUint64ToInt64(k.VotePeriod(ctx))

	votePeriodsPerWindow := math.LegacyNewDec(slashWindow).QuoInt64(votePeriod).TruncateInt64()
	numberOfAssets := int64(len(k.GetParams(ctx).MandatoryList))

	return (votePeriodsPerWindow * numberOfAssets)
}

// SetValidatorRewardSet will take all the current validators and store them
// in the ValidatorRewardSet to earn rewards in the current Slash Window.
func (k Keeper) SetValidatorRewardSet(ctx sdk.Context) error {
	validatorRewardSet := types.ValidatorRewardSet{
		ValidatorSet: []string{},
	}
	vals, err := k.StakingKeeper.GetBondedValidatorsByPower(ctx)
	if err != nil {
		return err
	}
	for _, v := range vals {
		addr := v.GetOperator()
		validatorRewardSet.ValidatorSet = append(validatorRewardSet.ValidatorSet, addr)
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&validatorRewardSet)
	store.Set(types.KeyValidatorRewardSet(), bz)
	return nil
}

// CurrentValidatorRewardSet returns the latest ValidatorRewardSet in the store.
func (k Keeper) GetValidatorRewardSet(ctx sdk.Context) types.ValidatorRewardSet {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyValidatorRewardSet())
	if bz == nil {
		return types.ValidatorRewardSet{}
	}

	var rewardSet types.ValidatorRewardSet
	k.cdc.MustUnmarshal(bz, &rewardSet)

	return rewardSet
}
