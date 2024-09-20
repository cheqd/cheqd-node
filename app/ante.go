package app

import (
	errorsmod "cosmossdk.io/errors"
	cheqdante "github.com/cheqd/cheqd-node/ante"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	ibcante "github.com/cosmos/ibc-go/v7/modules/core/ante"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	feeabsante "github.com/osmosis-labs/fee-abstraction/v7/x/feeabs/ante"
	feeabskeeper "github.com/osmosis-labs/fee-abstraction/v7/x/feeabs/keeper"
	feemarketante "github.com/skip-mev/feemarket/x/feemarket/ante"
	feemarketkeeper "github.com/skip-mev/feemarket/x/feemarket/keeper"
)

// HandlerOptions are the options required for constructing a default SDK AnteHandler.
// Here we add the cheqd ante decorators, which extend default SDK AnteHandler.
type HandlerOptions struct {
	AccountKeeper          cheqdante.AccountKeeper
	AccountKeeper          cheqdante.AccountKeeper
	BankKeeper             cheqdante.BankKeeper
	ExtensionOptionChecker ante.ExtensionOptionChecker
	FeegrantKeeper         ante.FeegrantKeeper
	SignModeHandler        authsigning.SignModeHandler
	SigGasConsumer         func(meter sdk.GasMeter, sig signing.SignatureV2, params types.Params) error
	TxFeeChecker           ante.TxFeeChecker
	TxFeeChecker           ante.TxFeeChecker
	IBCKeeper              *ibckeeper.Keeper
	DidKeeper              cheqdante.DidKeeper
	ResourceKeeper         cheqdante.ResourceKeeper
	FeeAbskeeper           feeabskeeper.Keeper
	FeeMarketKeeper        *feemarketkeeper.Keeper
}

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func NewAnteHandler(options HandlerOptions) (sdk.AnteHandler, error) {
	if options.AccountKeeper == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "account keeper is required for ante builder")
	}

	if options.BankKeeper == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "bank keeper is required for ante builder")
	}

	if options.SignModeHandler == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "sign mode handler is required for ante builder")
	}

	if options.IBCKeeper == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "IBC keeper is required for ante builder")
	}

	anteDecorators := []sdk.AnteDecorator{
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		ante.NewExtensionOptionsDecorator(options.ExtensionOptionChecker),
		feeabsante.NewFeeAbstrationMempoolFeeDecorator(options.FeeAbskeeper),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		feeabsante.NewFeeAbstractionDeductFeeDecorate(options.AccountKeeper, options.BankKeeper, options.FeeAbskeeper, options.FeegrantKeeper),
		cheqdante.NewOverAllDecorator( // fee market check replaces fee deduct decorator
			feeDecorators(options)...,
		),
		ante.NewSetPubKeyDecorator(options.AccountKeeper), // SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewValidateSigCountDecorator(options.AccountKeeper),
		ante.NewSigGasConsumeDecorator(options.AccountKeeper, options.SigGasConsumer),
		ante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		ante.NewIncrementSequenceDecorator(options.AccountKeeper),
		// v4 -> v5 ibc-go migration, v6 does not need migration
		ibcante.NewRedundantRelayDecorator(options.IBCKeeper),
	}

	return sdk.ChainAnteDecorators(anteDecorators...), nil
}

func feeDecorators(options HandlerOptions) []sdk.AnteDecorator {
	return []sdk.AnteDecorator{
		feemarketante.NewFeeMarketCheckDecorator( // fee market check replaces fee deduct decorator
			options.AccountKeeper,
			options.BankKeeper,
			options.FeegrantKeeper,
			options.FeeMarketKeeper,
			ante.NewDeductFeeDecorator(
				options.AccountKeeper,
				options.BankKeeper,
				options.FeegrantKeeper,
				options.TxFeeChecker,
			),
		),
	}
}

func feeDecorators(options HandlerOptions) []sdk.AnteDecorator {
	return []sdk.AnteDecorator{
		feemarketante.NewFeeMarketCheckDecorator( // fee market check replaces fee deduct decorator
			options.AccountKeeper,
			options.BankKeeper,
			options.FeegrantKeeper,
			options.FeeMarketKeeper,
			ante.NewDeductFeeDecorator(
				options.AccountKeeper,
				options.BankKeeper,
				options.FeegrantKeeper,
				options.TxFeeChecker,
			),
		),
	}
}
