package keeper

import (
	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cheqd/cheqd-node/util"
	"github.com/cheqd/cheqd-node/util/decmath"
	"github.com/cheqd/cheqd-node/x/oracle/types"
)

const denomErr = "denom: "

// HistoricMedians returns a list of a given denom's last numStamps medians.
func (k Keeper) HistoricMedians(
	ctx sdk.Context,
	denom string,
	numStamps uint64,
) types.PriceStamps {
	medians := types.PriceStamps{}

	k.IterateHistoricMedians(ctx, denom, uint(numStamps), func(median types.PriceStamp) bool {
		medians = append(medians, median)
		return false
	})

	return medians
}

func (k Keeper) HistoricDeviations(
	ctx sdk.Context,
	denom string,
	numStamps uint64,
) types.PriceStamps {
	deviations := types.PriceStamps{}

	k.IterateHistoricDeviations(ctx, denom, uint(numStamps), func(median types.PriceStamp) bool {
		deviations = append(deviations, median)
		return false
	})

	return deviations
}

// CalcAndSetHistoricMedian uses all the historic prices of a given denom to
// calculate its median price at the current block and set it to the store.
// It will also call setMedianDeviation with the calculated median.
func (k Keeper) CalcAndSetHistoricMedian(
	ctx sdk.Context,
	denom string,
) error {
	historicPrices := k.historicPrices(ctx, denom, k.MaximumPriceStamps(ctx))
	median, err := decmath.Median(historicPrices)
	if err != nil {
		return errors.Wrap(err, denomErr+denom)
	}

	block := util.SafeInt64ToUint64(ctx.BlockHeight())
	k.SetHistoricMedian(ctx, denom, block, median)

	return k.calcAndSetHistoricMedianDeviation(ctx, denom, median, historicPrices)
}

func (k Keeper) SetHistoricMedian(
	ctx sdk.Context,
	denom string,
	blockNum uint64,
	median math.LegacyDec,
) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&sdk.DecProto{Dec: median})
	store.Set(types.KeyMedian(denom, blockNum), bz)
}

// HistoricMedianDeviation returns a given denom's most recently stamped
// standard deviation around its median price at a given block.
func (k Keeper) HistoricMedianDeviation(
	ctx sdk.Context,
	denom string,
) (*types.PriceStamp, error) {
	store := ctx.KVStore(k.storeKey)
	blockDiff := util.SafeInt64ToUint64(ctx.BlockHeight())%k.MedianStampPeriod(ctx) + 1
	blockNum := util.SafeInt64ToUint64(ctx.BlockHeight()) - blockDiff
	bz := store.Get(types.KeyMedianDeviation(denom, blockNum))
	if bz == nil {
		return &types.PriceStamp{}, types.ErrNoMedianDeviation.Wrap(denomErr + denom)
	}

	decProto := sdk.DecProto{}
	k.cdc.MustUnmarshal(bz, &decProto)

	return types.NewPriceStamp(decProto.Dec, denom, blockNum), nil
}

// WithinHistoricMedianDeviation returns whether or not the current price of a
// given denom is within the latest stamped Standard Deviation around
// the Median.
func (k Keeper) WithinHistoricMedianDeviation(
	ctx sdk.Context,
	denom string,
) (bool, error) {
	// get latest median
	medians := k.HistoricMedians(ctx, denom, 1)
	if len(medians) == 0 {
		return false, types.ErrNoMedian.Wrap(denomErr + denom)
	}
	median := medians[0].ExchangeRate.Amount

	// get latest historic price
	prices := k.historicPrices(ctx, denom, 1)
	if len(prices) == 0 {
		return false, types.ErrNoHistoricPrice.Wrap(denomErr + denom)
	}
	price := prices[0]

	medianDeviation, err := k.HistoricMedianDeviation(ctx, denom)
	if err != nil {
		return false, err
	}

	return price.Sub(median).Abs().LTE(medianDeviation.ExchangeRate.Amount), nil
}

// calcAndSetHistoricMedianDeviation calculates and sets a given denom's standard
// deviation around its median price in the current block.
func (k Keeper) calcAndSetHistoricMedianDeviation(
	ctx sdk.Context,
	denom string,
	median math.LegacyDec,
	prices []math.LegacyDec,
) error {
	medianDeviation, err := decmath.MedianDeviation(median, prices)
	if err != nil {
		return errors.Wrap(err, denomErr+denom)
	}

	block := util.SafeInt64ToUint64(ctx.BlockHeight())
	k.SetHistoricMedianDeviation(ctx, denom, block, medianDeviation)
	return nil
}

func (k Keeper) SetHistoricMedianDeviation(
	ctx sdk.Context,
	denom string,
	blockNum uint64,
	medianDeviation math.LegacyDec,
) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&sdk.DecProto{Dec: medianDeviation})
	store.Set(types.KeyMedianDeviation(denom, blockNum), bz)
}

// MedianOfHistoricMedians calculates and returns the median of the last stampNum
// historic medians as well as the amount of medians used to calculate that median.
// If no medians are available, all returns are zero and error is nil.
func (k Keeper) MedianOfHistoricMedians(
	ctx sdk.Context,
	denom string,
	numStamps uint64,
) (math.LegacyDec, uint32, error) {
	medians := k.HistoricMedians(ctx, denom, numStamps)
	if len(medians) == 0 {
		return math.LegacyZeroDec(), 0, nil
	}
	median, err := decmath.Median(medians.Decs())
	if err != nil {
		return math.LegacyZeroDec(), 0, errors.Wrap(err, denomErr+denom)
	}

	return median, util.SafeIntToUint32(len(medians)), nil
}

// AverageOfHistoricMedians calculates and returns the average of the last stampNum
// historic medians as well as the amount of medians used to calculate that average.
// If no medians are available, all returns are zero and error is nil.
func (k Keeper) AverageOfHistoricMedians(
	ctx sdk.Context,
	denom string,
	numStamps uint64,
) (math.LegacyDec, uint32, error) {
	medians := k.HistoricMedians(ctx, denom, numStamps)
	if len(medians) == 0 {
		return math.LegacyZeroDec(), 0, nil
	}
	average, err := decmath.Average(medians.Decs())
	if err != nil {
		return math.LegacyZeroDec(), 0, errors.Wrap(err, denomErr+denom)
	}

	return average, util.SafeIntToUint32(len(medians)), nil
}

// MaxOfHistoricMedians calculates and returns the maximum value of the last stampNum
// historic medians as well as the amount of medians used to calculate that maximum.
// If no medians are available, all returns are zero and error is nil.
func (k Keeper) MaxOfHistoricMedians(
	ctx sdk.Context,
	denom string,
	numStamps uint64,
) (math.LegacyDec, uint32, error) {
	medians := k.HistoricMedians(ctx, denom, numStamps)
	if len(medians) == 0 {
		return math.LegacyZeroDec(), 0, nil
	}
	max, err := decmath.Max(medians.Decs())
	if err != nil {
		return math.LegacyZeroDec(), 0, errors.Wrap(err, denomErr+denom)
	}

	return max, util.SafeIntToUint32(len(medians)), nil
}

// MinOfHistoricMedians calculates and returns the minimum value of the last stampNum
// historic medians as well as the amount of medians used to calculate that minimum.
// If no medians are available, all returns are zero and error is nil.
func (k Keeper) MinOfHistoricMedians(
	ctx sdk.Context,
	denom string,
	numStamps uint64,
) (math.LegacyDec, uint32, error) {
	medians := k.HistoricMedians(ctx, denom, numStamps)
	if len(medians) == 0 {
		return math.LegacyZeroDec(), 0, nil
	}
	min, err := decmath.Min(medians.Decs())
	if err != nil {
		return math.LegacyZeroDec(), 0, errors.Wrap(err, denomErr+denom)
	}

	return min, util.SafeIntToUint32(len(medians)), nil
}

// historicPrices returns all the historic prices of a given denom.
func (k Keeper) historicPrices(
	ctx sdk.Context,
	denom string,
	numStamps uint64,
) []math.LegacyDec {
	// calculate start block to iterate from
	historicPrices := []math.LegacyDec{}

	k.IterateHistoricPrices(ctx, denom, uint(numStamps), func(exchangeRate math.LegacyDec) bool {
		historicPrices = append(historicPrices, exchangeRate)
		return false
	})

	return historicPrices
}

// IterateHistoricPrices iterates over historic prices of a given
// denom in the store in reverse.
// Iterator stops when exhausting the source, or when the handler returns `true`.
func (k Keeper) IterateHistoricPrices(
	ctx sdk.Context,
	denom string,
	numStamps uint,
	handler func(math.LegacyDec) bool,
) {
	store := ctx.KVStore(k.storeKey)

	// make sure we have one zero byte to correctly separate denoms
	prefix := util.ConcatBytes(1, types.KeyPrefixHistoricPrice, []byte(denom))
	iter := storetypes.KVStoreReversePrefixIteratorPaginated(store, prefix, 1, numStamps)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		decProto := sdk.DecProto{}
		k.cdc.MustUnmarshal(iter.Value(), &decProto)
		if handler(decProto.Dec) {
			break
		}
	}
}

// IterateHistoricMedians iterates over medians of a given
// denom in the store in reverse.
// Iterator stops when exhausting the source, or when the handler returns `true`.
func (k Keeper) IterateHistoricMedians(
	ctx sdk.Context,
	denom string,
	numStamps uint,
	handler func(types.PriceStamp) bool,
) {
	store := ctx.KVStore(k.storeKey)

	// make sure we have one zero byte to correctly separate denoms
	prefix := util.ConcatBytes(1, types.KeyPrefixMedian, []byte(denom))
	iter := storetypes.KVStoreReversePrefixIteratorPaginated(store, prefix, 1, numStamps)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		denom, block := types.ParseDenomAndBlockFromKey(iter.Key(), types.KeyPrefixMedian)
		decProto := sdk.DecProto{}
		k.cdc.MustUnmarshal(iter.Value(), &decProto)
		price := types.NewPriceStamp(decProto.Dec, denom, block)
		if handler(*price) {
			break
		}
	}
}

// IterateHistoricDeviations iterates over medians of a given
// denom in the store in reverse.
// Iterator stops when exhausting the source, or when the handler returns `true`.
func (k Keeper) IterateHistoricDeviations(
	ctx sdk.Context,
	denom string,
	numStamps uint,
	handler func(types.PriceStamp) bool,
) {
	store := ctx.KVStore(k.storeKey)

	// make sure we have one zero byte to correctly separate denoms
	prefix := util.ConcatBytes(1, types.KeyPrefixMedianDeviation, []byte(denom))
	iter := storetypes.KVStoreReversePrefixIteratorPaginated(store, prefix, 1, numStamps)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		denom, block := types.ParseDenomAndBlockFromKey(iter.Key(), types.KeyPrefixMedian)
		decProto := sdk.DecProto{}
		k.cdc.MustUnmarshal(iter.Value(), &decProto)
		price := types.NewPriceStamp(decProto.Dec, denom, block)
		if handler(*price) {
			break
		}
	}
}

// AddHistoricPrice adds the historic price of a denom at the current
// block height.
func (k Keeper) AddHistoricPrice(
	ctx sdk.Context,
	denom string,
	exchangeRate math.LegacyDec,
) {
	block := util.SafeInt64ToUint64(ctx.BlockHeight())
	k.SetHistoricPrice(ctx, denom, block, exchangeRate)
}

func (k Keeper) SetHistoricPrice(
	ctx sdk.Context,
	denom string,
	blockNum uint64,
	exchangeRate math.LegacyDec,
) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&sdk.DecProto{Dec: exchangeRate})
	store.Set(types.KeyHistoricPrice(denom, blockNum), bz)

	// Store the latest block number for this denom
	lastBlockBz := sdk.Uint64ToBigEndian(blockNum)
	store.Set(types.KeyLastHistoricPriceBlock(denom), lastBlockBz)
}

// DeleteHistoricPrice deletes the historic price of a denom at a
// given block.
func (k Keeper) DeleteHistoricPrice(
	ctx sdk.Context,
	denom string,
	blockNum uint64,
) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyHistoricPrice(denom, blockNum))
}

// DeleteHistoricMedian deletes a given denom's median price at a given block.
func (k Keeper) DeleteHistoricMedian(
	ctx sdk.Context,
	denom string,
	blockNum uint64,
) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyMedian(denom, blockNum))
}

// DeleteHistoricMedianDeviation deletes a given denom's standard deviation
// around its median price at a given block.
func (k Keeper) DeleteHistoricMedianDeviation(
	ctx sdk.Context,
	denom string,
	blockNum uint64,
) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyMedianDeviation(denom, blockNum))
}

func (k Keeper) PruneHistoricPricesBeforeBlock(ctx sdk.Context, blockNum uint64) {
	k.IterateAllHistoricPrices(ctx, func(price types.PriceStamp) (stop bool) {
		if price.BlockNum <= blockNum {
			k.DeleteHistoricPrice(ctx, price.ExchangeRate.Denom, price.BlockNum)
		}
		return false
	})
}

func (k Keeper) PruneMediansBeforeBlock(ctx sdk.Context, blockNum uint64) {
	k.IterateAllMedianPrices(ctx, func(price types.PriceStamp) (stop bool) {
		if price.BlockNum <= blockNum {
			k.DeleteHistoricMedian(ctx, price.ExchangeRate.Denom, price.BlockNum)
		}
		return false
	})
}

func (k Keeper) PruneMedianDeviationsBeforeBlock(ctx sdk.Context, blockNum uint64) {
	k.IterateAllMedianDeviationPrices(ctx, func(price types.PriceStamp) (stop bool) {
		if price.BlockNum <= blockNum {
			k.DeleteHistoricMedianDeviation(ctx, price.ExchangeRate.Denom, price.BlockNum)
		}
		return false
	})
}
