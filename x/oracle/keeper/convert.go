package keeper

import (
	"context"
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cheqd/cheqd-node/util"
	"github.com/cheqd/cheqd-node/x/oracle/types"
)

func (q Querier) ConvertUSDCtoCHEQ(
	ctx context.Context,
	req *types.ConvertUSDCtoCHEQRequest,
) (*types.ConvertUSDCtoCHEQResponse, error) {
	// Step 1: Parse the string coin
	coin, err := sdk.ParseCoinNormalized(req.Amount)
	if err != nil {
		return nil, fmt.Errorf("invalid amount format: %w", err)
	}

	Sdkctx := sdk.UnwrapSDKContext(ctx)

	// Step 2: Get price using selected moving average
	var price sdkmath.LegacyDec
	switch req.MaType {
	case "sma":
		var found bool
		price, found = q.GetSMA(Sdkctx, types.CheqdSymbol)
		if !found {
			return nil, fmt.Errorf("no SMA found for %s", types.CheqdSymbol)
		}
	case "ema":
		var found bool
		price, found = q.GetEMA(Sdkctx, types.CheqdSymbol)
		if !found {
			return nil, fmt.Errorf("no EMA found for %s", types.CheqdSymbol)
		}
	case "wma":
		switch req.WmaStrategy {
		case string(WmaStrategyRecent), string(WmaStrategyOldest), string(WmaStrategyBalanced):
			var found bool
			price, found = q.GetWMA(Sdkctx, types.CheqdSymbol, req.WmaStrategy)
			if !found {
				return nil, fmt.Errorf("missing WMA for strategy: %s", req.WmaStrategy)
			}
		case string(WmaStrategyCustom):
			price, err = q.GetCustomWMA(Sdkctx, types.CheqdSymbol, req.CustomWeights)
			if err != nil {
				return nil, fmt.Errorf("failed to compute custom WMA: %w", err)
			}
		default:
			return nil, fmt.Errorf("invalid WMA strategy: %s", req.WmaStrategy)
		}
	default:
		return nil, fmt.Errorf("invalid MA type: %s", req.MaType)
	}

	// Step 3: Convert USD to ncheq using helper
	ncheq, err := q.ConvertUSDToNcheq(coin, price)
	if err != nil {
		return nil, fmt.Errorf("conversion error: %w", err)
	}

	// Step 4: Return response
	return &types.ConvertUSDCtoCHEQResponse{
		Amount: ncheq.String(), // returns something like "7424625ncheq"
	}, nil
}

func (q Querier) ConvertUSDToNcheq(usd sdk.Coin, price sdkmath.LegacyDec) (sdk.Coin, error) {
	if usd.Denom != types.UsdDenom {
		return sdk.Coin{}, fmt.Errorf("expected denom to be 'usd', got: %s", usd.Denom)
	}

	if price.IsZero() {
		return sdk.Coin{}, fmt.Errorf("cannot convert: price is zero")
	}

	// 1. Convert USD (1e18) → float USD (1.0, 2.5, etc.)
	usdRaw := usd.Amount.ToLegacyDec().Quo(sdkmath.LegacyNewDecFromInt(util.UsdExponent))

	// 2. Divide USD / CHEQ/USD → CHEQ
	cheq := usdRaw.Quo(price)

	// 3. Scale to ncheq (1e9)
	ncheq := cheq.Mul(sdkmath.LegacyNewDecFromInt(util.CheqScale)).TruncateInt()

	if ncheq.IsZero() {
		return sdk.Coin{}, fmt.Errorf("converted amount is zero")
	}

	return sdk.NewCoin(types.CheqdDenom, ncheq), nil
}
