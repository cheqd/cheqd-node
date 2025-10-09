package keeper

import (
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cheqd/cheqd-node/util"
	"github.com/cheqd/cheqd-node/x/oracle/types"
)

// ScheduleParamUpdatePlan schedules a param update plan.
func (k Keeper) ScheduleParamUpdatePlan(ctx sdk.Context, plan types.ParamUpdatePlan) error {
	if plan.Height < ctx.BlockHeight() {
		return types.ErrInvalidRequest.Wrap("param update cannot be scheduled in the past")
	}
	if err := k.ValidateParamChanges(ctx, plan.Keys, plan.Changes); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&plan)
	store.Set(types.KeyParamUpdatePlan(util.SafeInt64ToUint64(plan.Height)), bz)

	return nil
}

// ClearParamUpdatePlan will clear an upcoming param update plan if one exists and return
// an error if one isn't found.
func (k Keeper) ClearParamUpdatePlan(ctx sdk.Context, planHeight uint64) error {
	if !k.haveParamUpdatePlan(ctx, planHeight) {
		return types.ErrInvalidRequest.Wrapf("No param update plan found at block height %d", planHeight)
	}

	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyParamUpdatePlan(planHeight))
	return nil
}

// haveParamUpdatePlan will return whether a param update plan exists and the specified
// plan height.
func (k Keeper) haveParamUpdatePlan(ctx sdk.Context, planHeight uint64) bool {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyParamUpdatePlan(planHeight))
	return bz != nil
}

// GetParamUpdatePlans returns all the param update plans in the store.
func (k Keeper) GetParamUpdatePlans(ctx sdk.Context) (plans []types.ParamUpdatePlan) {
	k.IterateParamUpdatePlans(ctx, func(plan types.ParamUpdatePlan) bool {
		plans = append(plans, plan)
		return false
	})

	return plans
}

// IterateParamUpdatePlans iterates rate over param update plans in the store
func (k Keeper) IterateParamUpdatePlans(
	ctx sdk.Context,
	handler func(types.ParamUpdatePlan) bool,
) {
	store := ctx.KVStore(k.storeKey)

	iter := storetypes.KVStorePrefixIterator(store, types.KeyPrefixParamUpdatePlan)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var paramUpdatePlan types.ParamUpdatePlan
		k.cdc.MustUnmarshal(iter.Value(), &paramUpdatePlan)

		if handler(paramUpdatePlan) {
			break
		}
	}
}

// ValidateParamChanges validates parameter changes against the existing oracle parameters.
func (k Keeper) ValidateParamChanges(ctx sdk.Context, keys []string, changes types.Params) error {
	params := k.GetParams(ctx)

	for _, key := range keys {
		switch key {
		case string(types.KeyVotePeriod):
			params.VotePeriod = changes.VotePeriod

		case string(types.KeyVoteThreshold):
			params.VoteThreshold = changes.VoteThreshold

		case string(types.KeyRewardBands):
			params.RewardBands = changes.RewardBands

		case string(types.KeyRewardDistributionWindow):
			params.RewardDistributionWindow = changes.RewardDistributionWindow

		case string(types.KeyAcceptList):
			params.AcceptList = changes.AcceptList.Normalize()

		case string(types.KeyMandatoryList):
			params.MandatoryList = changes.MandatoryList.Normalize()

		case string(types.KeySlashFraction):
			params.SlashFraction = changes.SlashFraction

		case string(types.KeySlashWindow):
			params.SlashWindow = changes.SlashWindow

		case string(types.KeyMinValidPerWindow):
			params.MinValidPerWindow = changes.MinValidPerWindow

		case string(types.KeyHistoricStampPeriod):
			params.HistoricStampPeriod = changes.HistoricStampPeriod

		case string(types.KeyMedianStampPeriod):
			params.MedianStampPeriod = changes.MedianStampPeriod

		case string(types.KeyMaximumPriceStamps):
			params.MaximumPriceStamps = changes.MaximumPriceStamps

		case string(types.KeyMaximumMedianStamps):
			params.MaximumMedianStamps = changes.MaximumMedianStamps

		case string(types.KeyCurrencyPairProviders):
			params.CurrencyPairProviders = changes.CurrencyPairProviders

		case string(types.KeyCurrencyDeviationThresholds):
			params.CurrencyDeviationThresholds = changes.CurrencyDeviationThresholds

		case string(types.KeyUsdcIbcDenom):
			params.UsdcIbcDenom = changes.UsdcIbcDenom

		case string(types.KeySlashingEnabled):
			params.SlashingEnabled = changes.SlashingEnabled
		}
	}

	return params.Validate()
}

// ExecuteParamUpdatePlan will execute a given param update plan and emit a param
// update event.
func (k Keeper) ExecuteParamUpdatePlan(ctx sdk.Context, plan types.ParamUpdatePlan) error {
	for _, key := range plan.Keys {
		switch key {
		case string(types.KeyVotePeriod):
			k.SetVotePeriod(ctx, plan.Changes.VotePeriod)

		case string(types.KeyVoteThreshold):
			k.SetVoteThreshold(ctx, plan.Changes.VoteThreshold)

		case string(types.KeyRewardBands):
			k.SetRewardBand(ctx, plan.Changes.RewardBands)

		case string(types.KeyRewardDistributionWindow):
			k.SetRewardDistributionWindow(ctx, plan.Changes.RewardDistributionWindow)

		case string(types.KeyAcceptList):
			k.SetAcceptList(ctx, plan.Changes.AcceptList.Normalize())

		case string(types.KeyMandatoryList):
			k.SetMandatoryList(ctx, plan.Changes.MandatoryList.Normalize())

		case string(types.KeySlashFraction):
			k.SetSlashFraction(ctx, plan.Changes.SlashFraction)

		case string(types.KeySlashWindow):
			k.SetSlashWindow(ctx, plan.Changes.SlashWindow)

		case string(types.KeyMinValidPerWindow):
			k.SetMinValidPerWindow(ctx, plan.Changes.MinValidPerWindow)

		case string(types.KeyHistoricStampPeriod):
			k.SetHistoricStampPeriod(ctx, plan.Changes.HistoricStampPeriod)

		case string(types.KeyMedianStampPeriod):
			k.SetMedianStampPeriod(ctx, plan.Changes.MedianStampPeriod)

		case string(types.KeyMaximumPriceStamps):
			k.SetMaximumPriceStamps(ctx, plan.Changes.MaximumPriceStamps)

		case string(types.KeyMaximumMedianStamps):
			k.SetMaximumMedianStamps(ctx, plan.Changes.MaximumMedianStamps)

		case string(types.KeyCurrencyPairProviders):
			k.SetCurrencyPairProviders(ctx, plan.Changes.CurrencyPairProviders)

		case string(types.KeyCurrencyDeviationThresholds):
			k.SetCurrencyDeviationThresholds(ctx, plan.Changes.CurrencyDeviationThresholds)

		case string(types.KeyUsdcIbcDenom):
			k.SetUsdcIbcDenom(ctx, plan.Changes.UsdcIbcDenom)

		case string(types.KeySlashingEnabled):
			k.SetSlashingEnabled(ctx, plan.Changes.SlashingEnabled)
		}
	}

	event := sdk.NewEvent(
		types.EventParamUpdate,
		sdk.NewAttribute(types.AttributeKeyNotifyPriceFeeder, "1"),
	)
	ctx.EventManager().EmitEvent(event)

	// clear plan from store after executing it
	return k.ClearParamUpdatePlan(ctx, util.SafeInt64ToUint64(plan.Height))
}
