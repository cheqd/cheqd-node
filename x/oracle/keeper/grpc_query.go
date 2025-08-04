package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cheqd/cheqd-node/util"
	"github.com/cheqd/cheqd-node/x/oracle/types"
)

var _ types.QueryServer = Querier{}

// Querier implements a QueryServer for the x/oracle module.
type Querier struct {
	Keeper
}

// NewQuerier returns an implementation of the oracle QueryServer interface
// for the provided Keeper.
func NewQuerier(keeper Keeper) types.QueryServer {
	return &Querier{Keeper: keeper}
}

// Params queries params of x/oracle module.
func (q Querier) Params(
	goCtx context.Context,
	req *types.QueryParams,
) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	params := q.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

// ExchangeRates queries exchange rates of all denoms, or, if specified, returns
// a single denom.
func (q Querier) ExchangeRates(
	goCtx context.Context,
	req *types.QueryExchangeRates,
) (*types.QueryExchangeRatesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	var exchangeRates sdk.DecCoins

	if len(req.Denom) > 0 {
		exchangeRate, err := q.GetExchangeRate(ctx, req.Denom)
		if err != nil {
			return nil, err
		}

		exchangeRates = exchangeRates.Add(sdk.NewDecCoinFromDec(req.Denom, exchangeRate))
	} else {
		q.IterateExchangeRates(ctx, func(denom string, rate math.LegacyDec) (stop bool) {
			exchangeRates = exchangeRates.Add(sdk.NewDecCoinFromDec(denom, rate))
			return false
		})
	}

	return &types.QueryExchangeRatesResponse{ExchangeRates: exchangeRates}, nil
}

// ActiveExchangeRates queries all denoms for which exchange rates exist.
func (q Querier) ActiveExchangeRates(
	goCtx context.Context,
	req *types.QueryActiveExchangeRates,
) (*types.QueryActiveExchangeRatesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	denoms := []string{}
	q.IterateExchangeRates(ctx, func(denom string, _ math.LegacyDec) (stop bool) {
		denoms = append(denoms, denom)
		return false
	})

	return &types.QueryActiveExchangeRatesResponse{ActiveRates: denoms}, nil
}

// FeederDelegation queries the account address to which the validator operator
// delegated oracle vote rights.
func (q Querier) FeederDelegation(
	goCtx context.Context,
	req *types.QueryFeederDelegation,
) (*types.QueryFeederDelegationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	feederAddr, err := q.GetFeederDelegation(ctx, valAddr)
	if err != nil {
		return nil, err
	}

	return &types.QueryFeederDelegationResponse{
		FeederAddr: feederAddr.String(),
	}, nil
}

// MissCounter queries oracle miss counter of a validator.
func (q Querier) MissCounter(
	goCtx context.Context,
	req *types.QueryMissCounter,
) (*types.QueryMissCounterResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryMissCounterResponse{
		MissCounter: q.GetMissCounter(ctx, valAddr),
	}, nil
}

// SlashWindow queries the current slash window progress of the oracle.
func (q Querier) SlashWindow(
	goCtx context.Context,
	req *types.QuerySlashWindow,
) (*types.QuerySlashWindowResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	params := q.GetParams(ctx)

	return &types.QuerySlashWindowResponse{
		WindowProgress: (util.SafeInt64ToUint64(ctx.BlockHeight()) % params.SlashWindow) /
			params.VotePeriod,
	}, nil
}

// AggregatePrevote queries an aggregate prevote of a validator.
func (q Querier) AggregatePrevote(
	goCtx context.Context,
	req *types.QueryAggregatePrevote,
) (*types.QueryAggregatePrevoteResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	prevote, err := q.GetAggregateExchangeRatePrevote(ctx, valAddr)
	if err != nil {
		return nil, err
	}

	return &types.QueryAggregatePrevoteResponse{
		AggregatePrevote: prevote,
	}, nil
}

// AggregatePrevotes queries aggregate prevotes of all validators
func (q Querier) AggregatePrevotes(
	goCtx context.Context,
	req *types.QueryAggregatePrevotes,
) (*types.QueryAggregatePrevotesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	var prevotes []types.AggregateExchangeRatePrevote
	q.IterateAggregateExchangeRatePrevotes(ctx, func(_ sdk.ValAddress, prevote types.AggregateExchangeRatePrevote) bool {
		prevotes = append(prevotes, prevote)
		return false
	})

	return &types.QueryAggregatePrevotesResponse{
		AggregatePrevotes: prevotes,
	}, nil
}

// AggregateVote queries an aggregate vote of a validator
func (q Querier) AggregateVote(
	goCtx context.Context,
	req *types.QueryAggregateVote,
) (*types.QueryAggregateVoteResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	vote, err := q.GetAggregateExchangeRateVote(ctx, valAddr)
	if err != nil {
		return nil, err
	}

	return &types.QueryAggregateVoteResponse{
		AggregateVote: vote,
	}, nil
}

// AggregateVotes queries aggregate votes of all validators
func (q Querier) AggregateVotes(
	goCtx context.Context,
	req *types.QueryAggregateVotes,
) (*types.QueryAggregateVotesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	var votes []types.AggregateExchangeRateVote
	q.IterateAggregateExchangeRateVotes(ctx, func(_ sdk.ValAddress, vote types.AggregateExchangeRateVote) bool {
		votes = append(votes, vote)
		return false
	})

	return &types.QueryAggregateVotesResponse{
		AggregateVotes: votes,
	}, nil
}

// Medians queries medians of all denoms, or, if specified, returns
// a single median.
func (q Querier) Medians(
	goCtx context.Context,
	req *types.QueryMedians,
) (*types.QueryMediansResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	medians := types.PriceStamps{}

	if len(req.Denom) > 0 {
		if req.NumStamps == 0 {
			return nil, status.Error(codes.InvalidArgument, "parameter NumStamps must be greater than 0")
		}

		if req.NumStamps > util.SafeUint64ToUint32(q.MaximumMedianStamps(ctx)) {
			req.NumStamps = util.SafeUint64ToUint32(q.MaximumMedianStamps(ctx))
		}

		medians = q.HistoricMedians(ctx, req.Denom, uint64(req.NumStamps))
	} else {
		medians = q.AllMedianPrices(ctx)
	}

	return &types.QueryMediansResponse{Medians: *medians.Sort()}, nil
}

// MedianDeviations queries median deviations of all denoms, or, if specified, returns
// a single median deviation.
func (q Querier) MedianDeviations(
	goCtx context.Context,
	req *types.QueryMedianDeviations,
) (*types.QueryMedianDeviationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	medianDeviations := types.PriceStamps{}

	if len(req.Denom) > 0 {
		price, err := q.HistoricMedianDeviation(ctx, req.Denom)
		if err != nil {
			return nil, err
		}
		medianDeviations = append(medianDeviations, *price)
	} else {
		medianDeviations = q.AllMedianDeviationPrices(ctx)
	}

	return &types.QueryMedianDeviationsResponse{MedianDeviations: *medianDeviations.Sort()}, nil
}

// ValidatorRewardSet queries the list of validators that can earn rewards in
// the current Slash Window.
func (q Querier) ValidatorRewardSet(
	goCtx context.Context,
	req *types.QueryValidatorRewardSet,
) (*types.QueryValidatorRewardSetResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	validatorRewardSet := q.GetValidatorRewardSet(ctx)

	return &types.QueryValidatorRewardSetResponse{
		Validators: validatorRewardSet,
	}, nil
}

func (q Querier) EMA(ctx context.Context, req *types.QueryEMARequest) (*types.QueryEMAResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	price, present := q.GetEMA(sdkCtx, req.Denom)
	if !present {
		return nil, errors.New("ema not present")
	}
	return &types.QueryEMAResponse{
		Price: price,
	}, nil
}

func (q Querier) WMA(goCtx context.Context, req *types.QueryWMARequest) (*types.QueryWMAResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	if !IsValidWmaStrategy(req.Strategy) {
		return nil, types.ErrInvalidWmaStrategy.Wrapf("invalid strategy: %s", req.Strategy)
	}

	// Handle custom strategy separately (needs live price data + weights)
	if req.Strategy == string(WmaStrategyCustom) {
		prices, err := CalculateHistoricPrices(ctx, ctx.KVStore(q.storeKey), req.Denom, uint64(ctx.BlockHeight()), q.Keeper)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get historic prices: %v", err)
		}
		if len(prices) == 0 {
			return nil, status.Errorf(codes.NotFound, "no historic prices found for denom: %s", req.Denom)
		}
		if len(req.CustomWeights) != AveragingWindow {
			return nil, status.Errorf(codes.InvalidArgument,
				"custom_weights must have exactly %d elements (one per period)", AveragingWindow)
		}

		result := CalculateWMA(prices, req.Strategy, req.CustomWeights)
		return &types.QueryWMAResponse{Price: result}, nil
	}
	result, found := q.Keeper.GetWMA(ctx, req.Denom, req.Strategy)
	if !found {
		return nil, status.Errorf(codes.NotFound, "no WMA found for denom: %s with strategy: %s", req.Denom, req.Strategy)
	}

	return &types.QueryWMAResponse{Price: result}, nil
}

func (q Querier) SMA(ctx context.Context, req *types.QuerySMARequest) (*types.QuerySMAResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	price, present := q.GetSMA(sdkCtx, req.Denom)
	if !present {
		return nil, errors.New("sma not present")
	}
	return &types.QuerySMAResponse{
		Price: price,
	}, nil
}
