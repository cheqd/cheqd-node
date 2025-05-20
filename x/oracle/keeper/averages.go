package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	"github.com/cheqd/cheqd-node/util"
	"github.com/cheqd/cheqd-node/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func KeySMA(denom string) []byte { return []byte(fmt.Sprintf("sma:%s", denom)) }
func KeyWMA(denom string) []byte { return []byte(fmt.Sprintf("wma:%s", denom)) }
func KeyEMA(denom string) []byte { return []byte(fmt.Sprintf("ema:%s", denom)) }

const AveragingWindow = 10 // N-period for SMA/WMA/EMA // TODO: make it a param

func (k Keeper) ComputeAverages(ctx sdk.Context, denom string) error {
	store := ctx.KVStore(k.storeKey)
	currentBlock := util.SafeInt64ToUint64(ctx.BlockHeight())

	// Collect last N prices
	var prices []math.LegacyDec
	for i := uint64(0); i < AveragingWindow; i++ {
		bz := store.Get(types.KeyHistoricPrice(denom, currentBlock-i))
		if bz == nil {
			continue
		}
		var proto sdk.DecProto
		if err := k.cdc.Unmarshal(bz, &proto); err != nil {
			return err
		}
		prices = append(prices, proto.Dec)
	}
	if len(prices) == 0 {
		return nil
	}

	// 1. SMA
	sum := math.LegacyZeroDec()
	for _, p := range prices {
		sum = sum.Add(p)
	}
	sma := sum.QuoInt64(int64(len(prices)))

	// 2. WMA (weight = 1, 2, ..., N)
	weightedSum := math.LegacyZeroDec()
	weightTotal := int64(0)
	for i := 0; i < len(prices); i++ {
		weight := int64(i + 1)
		weightedSum = weightedSum.Add(prices[i].MulInt64(weight))
		weightTotal += weight
	}
	wma := weightedSum.QuoInt64(weightTotal)

	// 3. EMA (smoothing factor α = 2 / (N + 1))
	ema, present := k.GetAverage(ctx, KeyEMA(denom))

	if !present {
		ema = sma
	}
	alpha := math.LegacyNewDecWithPrec(2, 0).QuoInt64(int64(AveragingWindow + 1))
	for i := 1; i < len(prices); i++ {
		ema = prices[i].Mul(alpha).Add(ema.Mul(math.LegacyOneDec().Sub(alpha)))
	}

	// Store averages
	k.setAverage(ctx, KeySMA(denom), sma)
	k.setAverage(ctx, KeyWMA(denom), wma)
	k.setAverage(ctx, KeyEMA(denom), ema)
	return nil
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

func (k Keeper) GetAverage(ctx sdk.Context, key []byte) (math.LegacyDec, bool) {
	bz := ctx.KVStore(k.storeKey).Get(key)
	if bz == nil {
		return math.LegacyZeroDec(), false
	}
	var proto sdk.DecProto
	k.cdc.MustUnmarshal(bz, &proto)
	return proto.Dec, true
}
