package posthandler

import (
	cheqdante "github.com/cheqd/cheqd-node/ante"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
)

// HandlerOptions are the options required for constructing a default post handler
type HandlerOptions struct {
	AccountKeeper  ante.AccountKeeper
	BankKeeper     cheqdante.BankKeeper
	FeegrantKeeper ante.FeegrantKeeper
	DidKeeper      cheqdante.DidKeeper
	ResourceKeeper cheqdante.ResourceKeeper
}

// NewPostHandler returns a default post handler
func NewPostHandler(options HandlerOptions) (sdk.AnteHandler, error) {
	postDecorators := []sdk.AnteDecorator{
		NewTaxDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper, options.DidKeeper, options.ResourceKeeper),
	}
	return sdk.ChainAnteDecorators(postDecorators...), nil
}
