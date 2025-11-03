package posthandler

import (
	cheqdante "github.com/cheqd/cheqd-node/ante"
	"github.com/cheqd/cheqd-node/pricefeeder"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	globalfeekeeper "github.com/noble-assets/globalfee/keeper"
	feeabskeeper "github.com/osmosis-labs/fee-abstraction/v8/x/feeabs/keeper"
)

// HandlerOptions are the options required for constructing a default post handler
type HandlerOptions struct {
	AccountKeeper   ante.AccountKeeper
	BankKeeper      BankKeeper
	FeegrantKeeper  ante.FeegrantKeeper
	DidKeeper       cheqdante.DidKeeper
	ResourceKeeper  cheqdante.ResourceKeeper
	FeeMarketKeeper FeeMarketKeeper
	OracleKeeper    cheqdante.OracleKeeper
	FeeabsKeeper    feeabskeeper.Keeper
	PriceFeeder     *pricefeeder.PriceFeeder
	GlobalFeeKeeper *globalfeekeeper.Keeper
}

// NewPostHandler returns a default post handler
func NewPostHandler(options HandlerOptions) (sdk.PostHandler, error) {
	postDecorators := []sdk.PostDecorator{
		NewTaxDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper, options.DidKeeper, options.ResourceKeeper, options.FeeMarketKeeper, options.OracleKeeper, options.FeeabsKeeper, options.PriceFeeder, options.GlobalFeeKeeper),
	}
	return sdk.ChainPostDecorators(postDecorators...), nil
}
