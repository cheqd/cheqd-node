package ante

import (
	"github.com/cosmos/cosmos-sdk/types/errors"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type OverAllDecorator struct {
	decorators []sdk.AnteDecorator
}

func NewOverAllDecorator(decorators ...sdk.AnteDecorator) OverAllDecorator {
	return OverAllDecorator{
		decorators: decorators,
	}
}

func (dfd OverAllDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, errorsmod.Wrap(errors.ErrTxDecode, "Tx must be a FeeTx")
	}

	if !simulate && ctx.BlockHeight() > 0 && feeTx.GetGas() == 0 {
		return ctx, errorsmod.Wrap(errors.ErrInvalidGasLimit, "must provide positive gas")
	}

	var (
		priority int64
		err      error
	)
	// check if the tx is a taxable tx
	// CONTRACT: Taxable tx is a tx that has at least 1 taxable related Msg.
	taxable := IsTaxableTxLite(tx)
	// if taxable, include in the mempool
	if taxable {
		// default priority of tx
		newCtx := ctx.WithPriority(priority)
		// posthandler will deduct the fee from the fee payer
		return next(newCtx, tx, simulate)
	}

	handler := sdk.ChainAnteDecorators(dfd.decorators...)
	newCtx, err := handler(ctx, tx, simulate)
	if err != nil {
		return newCtx, err
	}
	return next(newCtx, tx, simulate)
}

// CheckTxFee implements the logic for the fee market to check if a Tx has provided sufficient
// fees given the current state of the fee market. Returns an error if insufficient fees.
func CheckTxFee(ctx sdk.Context, gasPrice sdk.DecCoin, feeCoin sdk.Coin, feeGas int64, isAnte bool) (payCoin sdk.Coin, tip sdk.Coin, err error) {
	payCoin = feeCoin
	// Ensure that the provided fees meet the minimum
	if !gasPrice.IsZero() {
		var (
			requiredFee sdk.Coin
			consumedFee sdk.Coin
		)

		// Determine the required fees by multiplying each required minimum gas
		// price by the gas, where fee = ceil(minGasPrice * gas).
		gasConsumed := int64(ctx.GasMeter().GasConsumed())
		gcDec := sdkmath.LegacyNewDec(gasConsumed)
		glDec := sdkmath.LegacyNewDec(feeGas)

		consumedFeeAmount := gasPrice.Amount.Mul(gcDec)
		limitFee := gasPrice.Amount.Mul(glDec)

		consumedFee = sdk.NewCoin(gasPrice.Denom, consumedFeeAmount.Ceil().RoundInt())
		requiredFee = sdk.NewCoin(gasPrice.Denom, limitFee.Ceil().RoundInt())

		if !payCoin.IsGTE(requiredFee) {
			return sdk.Coin{}, sdk.Coin{}, errors.ErrInsufficientFee.Wrapf(
				"got: %s required: %s, minGasPrice: %s, gas: %d",
				payCoin,
				requiredFee,
				gasPrice,
				gasConsumed,
			)
		}

		if isAnte {
			tip = payCoin.Sub(requiredFee)
			payCoin = requiredFee
		} else {
			tip = payCoin.Sub(consumedFee)
			payCoin = consumedFee
		}
	}

	return payCoin, tip, nil
}
