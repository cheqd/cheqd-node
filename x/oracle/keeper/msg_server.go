package keeper

import (
	"context"
	"fmt"
	"strings"

	"cosmossdk.io/errors"
	ojoutils "github.com/cheqd/cheqd-node/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/cheqd/cheqd-node/x/oracle/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the oracle MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (ms msgServer) AggregateExchangeRatePrevote(
	goCtx context.Context,
	msg *types.MsgAggregateExchangeRatePrevote,
) (*types.MsgAggregateExchangeRatePrevoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	valAddr, err := sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return nil, err
	}
	feederAddr, err := sdk.AccAddressFromBech32(msg.Feeder)
	if err != nil {
		return nil, err
	}

	if err := ms.ValidateFeeder(ctx, feederAddr, valAddr); err != nil {
		return nil, err
	}

	// Ensure prevote wasn't already submitted
	if ms.HasAggregateExchangeRatePrevote(ctx, valAddr) {
		return nil, types.ErrExistingPrevote
	}

	// Convert hex string to votehash
	voteHash, err := types.AggregateVoteHashFromHexString(msg.Hash)
	if err != nil {
		return nil, types.ErrInvalidHash.Wrap(err.Error())
	}

	aggregatePrevote := types.NewAggregateExchangeRatePrevote(
		voteHash,
		valAddr,
		ojoutils.SafeInt64ToUint64(ctx.BlockHeight()),
	)
	ms.SetAggregateExchangeRatePrevote(ctx, valAddr, aggregatePrevote)

	return &types.MsgAggregateExchangeRatePrevoteResponse{}, nil
}

func (ms msgServer) AggregateExchangeRateVote(
	goCtx context.Context,
	msg *types.MsgAggregateExchangeRateVote,
) (*types.MsgAggregateExchangeRateVoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	valAddr, err := sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return nil, err
	}
	feederAddr, err := sdk.AccAddressFromBech32(msg.Feeder)
	if err != nil {
		return nil, err
	}
	if err := ms.ValidateFeeder(ctx, feederAddr, valAddr); err != nil {
		return nil, err
	}

	params := ms.GetParams(ctx)
	aggregatePrevote, err := ms.GetAggregateExchangeRatePrevote(ctx, valAddr)
	if err != nil {
		return nil, types.ErrNoAggregatePrevote.Wrap(msg.Validator)
	}

	// Check if the message is submitted in the proper period
	currentPeriod := ojoutils.SafeInt64ToUint64(ctx.BlockHeight()) / params.VotePeriod
	prevotePeriod := aggregatePrevote.SubmitBlock / params.VotePeriod
	if currentPeriod-prevotePeriod != 1 {
		return nil, types.ErrRevealPeriodMissMatch
	}

	exchangeRates, err := types.ParseExchangeRateDecCoins(msg.ExchangeRates)
	if err != nil {
		return nil, types.ErrInvalidExchangeRate.Wrap(err.Error())
	}

	// Verify a exchange rate with aggregate prevote hash
	hash := types.GetAggregateVoteHash(msg.Salt, msg.ExchangeRates, valAddr)
	if aggregatePrevote.Hash != hash.String() {
		return nil, types.ErrVerificationFailed.Wrapf("must be given %s not %s", aggregatePrevote.Hash, hash)
	}

	// Filter out rates which aren't included in the AcceptList
	filteredDecCoins := sdk.DecCoins{}
	for _, decCoin := range exchangeRates {
		if params.AcceptList.Contains(decCoin.Denom) {
			filteredDecCoins = append(filteredDecCoins, decCoin)
		}
	}

	// Move aggregate prevote to aggregate vote with given exchange rates
	ms.SetAggregateExchangeRateVote(ctx, valAddr, types.NewAggregateExchangeRateVote(filteredDecCoins, valAddr))
	ms.DeleteAggregateExchangeRatePrevote(ctx, valAddr)

	return &types.MsgAggregateExchangeRateVoteResponse{}, nil
}

func (ms msgServer) DelegateFeedConsent(
	goCtx context.Context,
	msg *types.MsgDelegateFeedConsent,
) (*types.MsgDelegateFeedConsentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	operatorAddr, err := sdk.ValAddressFromBech32(msg.Operator)
	if err != nil {
		return nil, err
	}

	delegateAddr, err := sdk.AccAddressFromBech32(msg.Delegate)
	if err != nil {
		return nil, err
	}

	val, err := ms.StakingKeeper.Validator(ctx, operatorAddr)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, stakingtypes.ErrNoValidatorFound.Wrap(msg.Operator)
	}

	ms.SetFeederDelegation(ctx, operatorAddr, delegateAddr)
	err = ctx.EventManager().EmitTypedEvent(&types.EventDelegateFeedConsent{
		Operator: msg.Operator, Delegate: msg.Delegate,
	})

	return &types.MsgDelegateFeedConsentResponse{}, err
}

func (ms msgServer) LegacyGovUpdateParams(
	goCtx context.Context,
	msg *types.MsgLegacyGovUpdateParams,
) (*types.MsgLegacyGovUpdateParamsResponse, error) {
	if msg.Authority != ms.authority {
		err := errors.Wrapf(
			types.ErrNoGovAuthority,
			"invalid authority; expected %s, got %s",
			ms.authority,
			msg.Authority,
		)
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	for _, key := range msg.Keys {
		switch key {
		case string(types.KeyVotePeriod):
			ms.SetVotePeriod(ctx, msg.Changes.VotePeriod)

		case string(types.KeyVoteThreshold):
			ms.SetVoteThreshold(ctx, msg.Changes.VoteThreshold)

		case string(types.KeyRewardBands):
			ms.SetRewardBand(ctx, msg.Changes.RewardBands)

		case string(types.KeyRewardDistributionWindow):
			if msg.Changes.RewardDistributionWindow < ms.Keeper.VotePeriod(ctx) {
				return nil, fmt.Errorf("oracle parameter RewardDistributionWindow must be greater than or equal with VotePeriod")
			}
			ms.SetRewardDistributionWindow(ctx, msg.Changes.RewardDistributionWindow)

		case string(types.KeyAcceptList):
			accept := msg.Changes.AcceptList.Normalize()
			mandatory := ms.Keeper.MandatoryList(ctx).Normalize()
			if !accept.ContainDenoms(mandatory) {
				return nil, fmt.Errorf("denom in MandatoryList not present in AcceptList")
			}
			ms.SetAcceptList(ctx, accept)

		case string(types.KeyMandatoryList):
			mandatory := msg.Changes.MandatoryList.Normalize()
			accept := ms.Keeper.AcceptList(ctx).Normalize()
			if !accept.ContainDenoms(mandatory) {
				return nil, fmt.Errorf("denom in MandatoryList not present in AcceptList")
			}
			ms.SetMandatoryList(ctx, mandatory)

		case string(types.KeySlashFraction):
			ms.SetSlashFraction(ctx, msg.Changes.SlashFraction)

		case string(types.KeySlashWindow):
			if msg.Changes.SlashWindow < ms.Keeper.VotePeriod(ctx) {
				return nil, fmt.Errorf("oracle parameter SlashWindow must be greater than or equal with VotePeriod")
			}
			ms.SetSlashWindow(ctx, msg.Changes.SlashWindow)

		case string(types.KeyMinValidPerWindow):
			ms.SetMinValidPerWindow(ctx, msg.Changes.MinValidPerWindow)

		case string(types.KeyHistoricStampPeriod):
			if msg.Changes.HistoricStampPeriod < 1 {
				return nil, fmt.Errorf("oracle parameters HistoricStampPeriod must be greater than 0")
			}
			if msg.Changes.HistoricStampPeriod > ms.Keeper.MedianStampPeriod(ctx) {
				return nil, fmt.Errorf("oracle parameter HistoricStampPeriod must be less than or equal with MedianStampPeriod")
			}
			if msg.Changes.HistoricStampPeriod%ms.Keeper.VotePeriod(ctx) != 0 {
				return nil, fmt.Errorf("oracle parameters HistoricStampPeriod must be exact multiples of VotePeriod")
			}
			ms.SetHistoricStampPeriod(ctx, msg.Changes.HistoricStampPeriod)

		case string(types.KeyMedianStampPeriod):
			if msg.Changes.MedianStampPeriod < 1 {
				return nil, fmt.Errorf("oracle parameters MedianStampPeriod must be greater than 0")
			}
			if msg.Changes.MedianStampPeriod < ms.Keeper.HistoricStampPeriod(ctx) {
				return nil, fmt.Errorf("oracle parameter MedianStampPeriod must be greater than or equal with HistoricStampPeriod")
			}
			if msg.Changes.MedianStampPeriod%ms.Keeper.VotePeriod(ctx) != 0 {
				return nil, fmt.Errorf("oracle parameters MedianStampPeriod must be exact multiples of VotePeriod")
			}
			ms.SetMedianStampPeriod(ctx, msg.Changes.MedianStampPeriod)

		case string(types.KeyMaximumPriceStamps):
			if msg.Changes.MaximumPriceStamps < 1 {
				return nil, fmt.Errorf("oracle parameters MaximumPriceStamps must be greater than 0")
			}
			ms.SetMaximumPriceStamps(ctx, msg.Changes.MaximumPriceStamps)

		case string(types.KeyMaximumMedianStamps):
			if msg.Changes.MaximumMedianStamps < 1 {
				return nil, fmt.Errorf("oracle parameters MaximumMedianStamps must be greater than 0")
			}
			ms.SetMaximumMedianStamps(ctx, msg.Changes.MaximumMedianStamps)

		case string(types.KeyUsdcIbcDenom):
			ms.SetUsdcIbcDenom(ctx, msg.Changes.UsdcIbcDenom)

		case string(types.KeySlashingEnabled):
			ms.SetSlashingEnabled(ctx, msg.Changes.SlashingEnabled)

		case string(types.KeyAveragingWindow):
			ms.SetAveragingWindow(ctx, msg.Changes.AveragingWindow)

		default:
			return nil, fmt.Errorf("%s is not an existing oracle param key", key)
		}
	}

	return &types.MsgLegacyGovUpdateParamsResponse{}, nil
}

func (ms msgServer) GovUpdateParams(
	goCtx context.Context,
	msg *types.MsgGovUpdateParams,
) (*types.MsgGovUpdateParamsResponse, error) {
	if msg.Authority != ms.authority {
		err := errors.Wrapf(
			types.ErrNoGovAuthority,
			"invalid authority; expected %s, got %s",
			ms.authority,
			msg.Authority,
		)
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	err := ms.ScheduleParamUpdatePlan(ctx, msg.Plan)
	if err != nil {
		return nil, err
	}

	return &types.MsgGovUpdateParamsResponse{}, nil
}

// GovAddDenoms adds new assets to the AcceptList, and adds
// it to the MandatoryList if specified.
func (ms msgServer) GovAddDenoms(
	goCtx context.Context,
	msg *types.MsgGovAddDenoms,
) (*types.MsgGovAddDenomsResponse, error) {
	if msg.Authority != ms.authority {
		err := errors.Wrapf(
			types.ErrNoGovAuthority,
			"invalid authority; expected %s, got %s",
			ms.authority,
			msg.Authority,
		)
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := ms.GetParams(ctx)

	plan := types.ParamUpdatePlan{
		Keys:    []string{},
		Height:  msg.Height,
		Changes: params,
	}
	for _, denom := range msg.DenomList {
		// if the AcceptList already contains this denom, and we're not
		// adding it to the "mandatory" list, error out.
		if plan.Changes.AcceptList.Contains(denom.SymbolDenom) && !msg.Mandatory {
			err := errors.Wrapf(
				types.ErrInvalidParamValue,
				"denom already exists in acceptList: %s",
				denom.SymbolDenom,
			)
			return nil, err
			// if the MandatoryList already contains this denom, and we're trying to
			// add it to the "mandatory" list, error out.
		} else if plan.Changes.MandatoryList.Contains(denom.SymbolDenom) && msg.Mandatory {
			err := errors.Wrapf(
				types.ErrInvalidParamValue,
				"denom already exists in mandatoryList: %s",
				denom.SymbolDenom,
			)
			return nil, err
		}

		// add to AcceptList & MandatoryList if necessary
		if !plan.Changes.AcceptList.Contains(denom.SymbolDenom) {
			plan.Changes.AcceptList = append(plan.Changes.AcceptList, denom)
			plan.Keys = ojoutils.AppendUniqueString(plan.Keys, string(types.KeyAcceptList))
		}
		if msg.Mandatory {
			plan.Changes.MandatoryList = append(plan.Changes.MandatoryList, denom)
			plan.Keys = ojoutils.AppendUniqueString(plan.Keys, string(types.KeyMandatoryList))
		}

		// add a RewardBand
		_, err := plan.Changes.RewardBands.GetBandFromDenom(denom.SymbolDenom)
		if err == types.ErrNoRewardBand {
			if msg.RewardBand != nil {
				plan.Changes.RewardBands.Add(denom.SymbolDenom, *msg.RewardBand)
			}
			plan.Changes.RewardBands.AddDefault(denom.SymbolDenom)
		} else if err != nil {
			return nil, err
		}
	}

	// append new currency pair providers
	for _, cpp := range msg.CurrencyPairProviders {
		plan.Keys = ojoutils.AppendUniqueString(plan.Keys, string(types.KeyCurrencyPairProviders))
		plan.Changes.CurrencyPairProviders = append(plan.Changes.CurrencyPairProviders, cpp)
	}

	// append new currency deviation thresholds
	for _, cdt := range msg.CurrencyDeviationThresholds {
		plan.Keys = ojoutils.AppendUniqueString(plan.Keys, string(types.KeyCurrencyDeviationThresholds))
		plan.Changes.CurrencyDeviationThresholds = append(plan.Changes.CurrencyDeviationThresholds, cdt)
	}

	// also update RewardBand key if new denoms are getting added
	if len(msg.DenomList) != 0 {
		plan.Keys = append(plan.Keys, string(types.KeyRewardBands))
	}

	// validate plan construction before scheduling
	err := plan.ValidateBasic()
	if err != nil {
		return nil, err
	}

	err = ms.ScheduleParamUpdatePlan(ctx, plan)
	if err != nil {
		return nil, err
	}

	return &types.MsgGovAddDenomsResponse{}, nil
}

// GovRemoveCurrencyPairProviders removes the specified currency pair
// providers in MsgGovRemoveCurrencyPairProviders if they exist in
// the current CurrencyPairProviders list.
func (ms msgServer) GovRemoveCurrencyPairProviders(
	goCtx context.Context,
	msg *types.MsgGovRemoveCurrencyPairProviders,
) (*types.MsgGovRemoveCurrencyPairProvidersResponse, error) {
	if msg.Authority != ms.authority {
		err := errors.Wrapf(
			types.ErrNoGovAuthority,
			"invalid authority; expected %s, got %s",
			ms.authority,
			msg.Authority,
		)
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := ms.GetParams(ctx)

	plan := types.ParamUpdatePlan{
		Keys:    []string{string(types.KeyCurrencyPairProviders)},
		Height:  msg.Height,
		Changes: params,
	}

	for _, cpp := range msg.CurrencyPairProviders {
		plan.Changes.CurrencyPairProviders = plan.Changes.CurrencyPairProviders.RemovePair(cpp)
	}

	// validate plan construction before scheduling
	err := plan.ValidateBasic()
	if err != nil {
		return nil, err
	}

	err = ms.ScheduleParamUpdatePlan(ctx, plan)
	if err != nil {
		return nil, err
	}

	return &types.MsgGovRemoveCurrencyPairProvidersResponse{}, nil
}

// GovRemoveCurrencyDeviationThresholds removes the specified currency
// deviation thresholds in MsgGovRemoveCurrencyDeviationThresholdsResponse
// if they exist in the current CurrencyDeviationThresholds list.
func (ms msgServer) GovRemoveCurrencyDeviationThresholds(
	goCtx context.Context,
	msg *types.MsgGovRemoveCurrencyDeviationThresholds,
) (*types.MsgGovRemoveCurrencyDeviationThresholdsResponse, error) {
	if msg.Authority != ms.authority {
		err := errors.Wrapf(
			types.ErrNoGovAuthority,
			"invalid authority; expected %s, got %s",
			ms.authority,
			msg.Authority,
		)
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := ms.GetParams(ctx)

	plan := types.ParamUpdatePlan{
		Keys:    []string{string(types.KeyCurrencyDeviationThresholds)},
		Height:  msg.Height,
		Changes: params,
	}

	for _, curr := range msg.Currencies {
		plan.Changes.CurrencyDeviationThresholds = plan.Changes.CurrencyDeviationThresholds.RemovePair(
			strings.ToUpper(curr),
		)
	}

	// validate plan construction before scheduling
	err := plan.ValidateBasic()
	if err != nil {
		return nil, err
	}

	err = ms.ScheduleParamUpdatePlan(ctx, plan)
	if err != nil {
		return nil, err
	}

	return &types.MsgGovRemoveCurrencyDeviationThresholdsResponse{}, nil
}

func (ms msgServer) GovCancelUpdateParamPlan(
	goCtx context.Context,
	msg *types.MsgGovCancelUpdateParamPlan,
) (*types.MsgGovCancelUpdateParamPlanResponse, error) {
	if msg.Authority != ms.authority {
		err := errors.Wrapf(
			types.ErrNoGovAuthority,
			"invalid authority; expected %s, got %s",
			ms.authority,
			msg.Authority,
		)
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	err := ms.ClearParamUpdatePlan(ctx, ojoutils.SafeInt64ToUint64(msg.Height))
	if err != nil {
		return nil, err
	}

	return &types.MsgGovCancelUpdateParamPlanResponse{}, nil
}
