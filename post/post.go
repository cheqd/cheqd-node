package posthandler

import (
	// cheqdante "github.com/cheqd/cheqd-node/ante"
	sdk "github.com/cosmos/cosmos-sdk/types"
	feemarketpost "github.com/skip-mev/feemarket/x/feemarket/post"
)

// HandlerOptions are the options required for constructing a default post handler
type HandlerOptions struct {
	AccountKeeper  feemarketpost.AccountKeeper
	BankKeeper     feemarketpost.BankKeeper
	FeegrantKeeper feemarketpost.FeeGrantKeeper
	// DidKeeper       cheqdante.DidKeeper
	// ResourceKeeper  cheqdante.ResourceKeeper
	FeeMarketKeeper feemarketpost.FeeMarketKeeper
}

// NewPostHandler returns a default post handler
func NewPostHandler(options HandlerOptions) (sdk.PostHandler, error) {
	postDecorators := []sdk.PostDecorator{
		feemarketpost.NewFeeMarketDeductDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper, options.FeeMarketKeeper),
	}
	return sdk.ChainPostDecorators(postDecorators...), nil
}
