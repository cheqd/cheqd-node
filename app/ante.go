package app

import (
	errorsmod "cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"
	circuitante "cosmossdk.io/x/circuit/ante"
	txsigning "cosmossdk.io/x/tx/signing"
	cheqdante "github.com/cheqd/cheqd-node/ante"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	ibcante "github.com/cosmos/ibc-go/v8/modules/core/ante"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	globalfeekeeper "github.com/noble-assets/globalfee/keeper"
	feeabsante "github.com/osmosis-labs/fee-abstraction/v8/x/feeabs/ante"
	feeabskeeper "github.com/osmosis-labs/fee-abstraction/v8/x/feeabs/keeper"
	feemarketante "github.com/skip-mev/feemarket/x/feemarket/ante"
	feemarketkeeper "github.com/skip-mev/feemarket/x/feemarket/keeper"
)

// HandlerOptions are the options required for constructing a default SDK AnteHandler.
// Here we add the cheqd ante decorators, which extend default SDK AnteHandler.
type HandlerOptions struct {
	AccountKeeper          cheqdante.AccountKeeper
	BankKeeper             bankkeeper.Keeper
	ExtensionOptionChecker ante.ExtensionOptionChecker
	FeegrantKeeper         ante.FeegrantKeeper
	SignModeHandler        *txsigning.HandlerMap
	SigGasConsumer         func(meter storetypes.GasMeter, sig signing.SignatureV2, params types.Params) error
	TxFeeChecker           ante.TxFeeChecker
	IBCKeeper              *ibckeeper.Keeper
	DidKeeper              cheqdante.DidKeeper
	ResourceKeeper         cheqdante.ResourceKeeper
	FeeAbskeeper           feeabskeeper.Keeper
	FeeMarketKeeper        *feemarketkeeper.Keeper
	CircuitKeeper          circuitante.CircuitBreaker
	GlobalFeeKeeper        *globalfeekeeper.Keeper
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
		circuitante.NewCircuitBreakerDecorator(options.CircuitKeeper),
		ante.NewExtensionOptionsDecorator(options.ExtensionOptionChecker),
		cheqdante.NewFeeAbsBypassDecorator(options.GlobalFeeKeeper),
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
	fallbackDecorator := ante.NewDeductFeeDecorator(
		options.AccountKeeper,
		options.BankKeeper,
		options.FeegrantKeeper,
		options.TxFeeChecker,
	)

	feeMarketDecorator := feemarketante.NewFeeMarketCheckDecorator( // fee market check replaces fee deduct decorator
		options.AccountKeeper,
		options.BankKeeper,
		options.FeegrantKeeper,
		options.FeeMarketKeeper,
		fallbackDecorator,
	)

	return []sdk.AnteDecorator{
		cheqdante.NewFeeMarketBypassDecorator(options.GlobalFeeKeeper, feeMarketDecorator, fallbackDecorator),
	}
}
