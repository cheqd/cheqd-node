package ante

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/noble-assets/globalfee/keeper"
	feeabstypes "github.com/osmosis-labs/fee-abstraction/v8/x/feeabs/types"
)

// feeAbsBypassDecorator annotates the context so the fee abstraction mempool
// decorator skips fee checks for transactions that only contain bypass
// messages.
type feeAbsBypassDecorator struct {
	globalFeeKeeper *keeper.Keeper
}

func NewFeeAbsBypassDecorator(globalFeeKeeper *keeper.Keeper) sdk.AnteDecorator {
	return feeAbsBypassDecorator{
		globalFeeKeeper: globalFeeKeeper,
	}
}

func (d feeAbsBypassDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	if ShouldBypassFeeMarket(ctx, d.globalFeeKeeper, tx) {
		ctx = ctx.WithContext(context.WithValue(ctx.Context(), feeabstypes.ByPassMsgKey{}, true))
	}

	return next(ctx, tx, simulate)
}
