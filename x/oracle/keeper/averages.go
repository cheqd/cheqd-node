package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	cosmosstore "cosmossdk.io/store/types"
	"github.com/cheqd/cheqd-node/util"
	"github.com/cheqd/cheqd-node/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type WmaStrategy string

const (
	WmaStrategyOldest   WmaStrategy = "OLDEST"
	WmaStrategyRecent   WmaStrategy = "RECENT"
	WmaStrategyBalanced WmaStrategy = "BALANCED"
	WmaStrategyCustom   WmaStrategy = "CUSTOM"
)

func KeySMA(denom string) []byte { return []byte(fmt.Sprintf("sma:%s", denom)) }
func KeyEMA(denom string) []byte { return []byte(fmt.Sprintf("ema:%s", denom)) }
func KeyWMAWithStrategy(denom, strategy string) []byte {
	return []byte(fmt.Sprintf("wma:%s:%s", denom, strategy))
}

const AveragingWindow = 3 // N-period for SMA/WMA/EMA // TODO: make it a param

func (k Keeper) ComputeAverages(ctx sdk.Context, denom string) error {
	store := ctx.KVStore(k.storeKey)
	currentBlock := util.SafeInt64ToUint64(ctx.BlockHeight())

	// Collect last N prices
	prices, err := CalculateHistoricPrices(ctx, store, denom, currentBlock, k)
	if err != nil {
		ctx.Logger().Error("Failed to calculate historic prices", "denom", denom, "error", err)
		return nil
	}
	if len(prices) == 0 {
		ctx.Logger().Error("No historic prices for denom", "denom", denom)
		return nil
	}
	// calculate sma and store it
	sma := CalculateSMA(prices)
	k.setAverage(ctx, KeySMA(denom), sma)

	// Calculate WMA for all strategies
	strategies := []string{string(WmaStrategyBalanced), string(WmaStrategyOldest), string(WmaStrategyRecent)}

	for _, strategy := range strategies {
		if !IsValidWmaStrategy(strategy) {
			return types.ErrInvalidWmaStrategy.Wrapf("invalid WMA strategy: %s", strategy)
		}
		wma := CalculateWMA(prices, strategy, nil)
		k.setAverage(ctx, KeyWMAWithStrategy(denom, strategy), wma)
	}

	// 3. EMA (smoothing factor α = 2 / (N + 1))
	prevEMA, present := k.GetAverage(ctx, KeyEMA(denom))

	ema := CalculateEMA(prevEMA, present, prices)
	k.setAverage(ctx, KeyEMA(denom), ema)
	return nil
}

func CalculateHistoricPrices(ctx sdk.Context, store cosmosstore.KVStore, denom string, currentBlock uint64, k Keeper) ([]math.LegacyDec, error) {
	var prices []math.LegacyDec
	// Get the last recorded block for this denom
	lastBlockBz := store.Get(types.KeyLastHistoricPriceBlock(denom))
	if lastBlockBz == nil {
		ctx.Logger().Error("No historic prices recorded", "denom", denom)
		return []math.LegacyDec{}, nil

	}
	lastBlock := sdk.BigEndianToUint64(lastBlockBz)
	for i := uint64(0); i < AveragingWindow; i++ {
		bz := store.Get(types.KeyHistoricPrice(denom, lastBlock-i*types.DefaultParams().HistoricStampPeriod))
		if bz == nil {
			continue
		}
		var proto sdk.DecProto
		if err := k.cdc.Unmarshal(bz, &proto); err != nil {
			return nil, err
		}
		prices = append(prices, proto.Dec)
	}
	return prices, nil
}

func CalculateEMA(previousEMA math.LegacyDec, present bool, prices []math.LegacyDec) math.LegacyDec {
	if len(prices) == 0 {
		return math.LegacyZeroDec()
	}

	// Initialize EMA with previous value or first price
	ema := previousEMA
	if !present {
		ema = prices[0]
	}

	alpha := math.LegacyNewDecWithPrec(2, 0).QuoInt64(int64(len(prices) + 1))

	for i := 1; i < len(prices); i++ {
		ema = prices[i].Mul(alpha).Add(ema.Mul(math.LegacyOneDec().Sub(alpha)))
	}

	return ema
}

func (k Keeper) setAverage(ctx sdk.Context, key []byte, value math.LegacyDec) {
	bz := k.cdc.MustMarshal(&sdk.DecProto{Dec: value})
	ctx.KVStore(k.storeKey).Set(key, bz)
}

func (k Keeper) GetSMA(ctx sdk.Context, denom string) (math.LegacyDec, bool) {
	return k.GetAverage(ctx, KeySMA(denom))
}

func (k Keeper) GetEMA(ctx sdk.Context, denom string) (math.LegacyDec, bool) {
	return k.GetAverage(ctx, KeyEMA(denom))
}

func (k Keeper) GetWMA(ctx sdk.Context, denom string, strategy string) (math.LegacyDec, bool) {
	return k.GetAverage(ctx, KeyWMAWithStrategy(denom, strategy))
}

func (k Keeper) GetAverage(ctx sdk.Context, key []byte) (math.LegacyDec, bool) {
	bz := ctx.KVStore(k.storeKey).Get(key)
	if bz == nil {
		return math.LegacyZeroDec(), false
	}
	var proto sdk.DecProto
	k.cdc.MustUnmarshal(bz, &proto)
	return proto.Dec, true
}

func CalculateSMA(prices []math.LegacyDec) math.LegacyDec {
	sum := math.LegacyZeroDec()
	for _, p := range prices {
		sum = sum.Add(p)
	}
	sma := sum.QuoInt64(int64(len(prices)))
	return sma
}

func CalculateWMA(prices []math.LegacyDec, strategy string, customWeights []int64) math.LegacyDec {
	n := len(prices)
	if n == 0 {
		return math.LegacyZeroDec()
	}

	weightedSum := math.LegacyZeroDec()
	weightTotal := int64(0)

	switch strategy {
	case "OLDEST":
		// Weights: [N, N-1, ..., 1]
		for i := 0; i < n; i++ {
			weight := int64(n - i)
			weightedSum = weightedSum.Add(prices[i].MulInt64(weight))
			weightTotal += weight
		}

	case "RECENT":
		// Weights: [1, 2, ..., N]
		for i := 0; i < n; i++ {
			weight := int64(i + 1)
			weightedSum = weightedSum.Add(prices[i].MulInt64(weight))
			weightTotal += weight
		}

	case "BALANCED":
		// Weights: [1–10], then ten × 10s, then [9–1] to make 30 entries
		// Adapt to whatever len(prices) is, but assume 30 ideal entries
		weights := make([]int64, n)
		for i := 0; i < n; i++ {
			switch {
			case i < 10:
				weights[i] = int64(i + 1)
			case i < 20:
				weights[i] = 10
			default:
				weights[i] = int64(30 - i)
			}
		}

		for i := 0; i < n; i++ {
			weightedSum = weightedSum.Add(prices[i].MulInt64(weights[i]))
			weightTotal += weights[i]
		}

	case "CUSTOM":
		// Use customWeights array provided by governance param or config
		if len(customWeights) != n {
			panic(fmt.Sprintf("custom weight length %d does not match price list length %d", len(customWeights), n))
		}

		for i := 0; i < n; i++ {
			weight := customWeights[i]
			weightedSum = weightedSum.Add(prices[i].MulInt64(weight))
			weightTotal += weight
		}

	default:
		panic(fmt.Sprintf("unsupported WMA strategy: %s", strategy))
	}

	return weightedSum.QuoInt64(weightTotal)
}

func IsValidWmaStrategy(s string) bool {
	switch WmaStrategy(s) {
	case WmaStrategyOldest, WmaStrategyRecent, WmaStrategyBalanced, WmaStrategyCustom:
		return true
	default:
		return false
	}
}
