package ante

import (
	"fmt"

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

	// fee := feeTx.GetFee()

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

// // The actual fee is deducted in the post handler along with the tip.
// func (dfd DeductFeeDecorator) EscrowFunds(ctx sdk.Context, sdkTx sdk.Tx, providedFee sdk.Coin) error {
// 	feeTx, ok := sdkTx.(sdk.FeeTx)
// 	if !ok {
// 		return errorsmod.Wrap(errors.ErrTxDecode, "Tx must be a FeeTx")
// 	}

// 	feePayer := feeTx.FeePayer()
// 	feeGranter := feeTx.FeeGranter()
// 	deductFeesFrom := feePayer

// 	// if feegranter set deduct fee from feegranter account.
// 	// this works with only when feegrant enabled.
// 	if feeGranter != nil {
// 		if dfd.feegrantKeeper == nil {
// 			return errors.ErrInvalidRequest.Wrap("fee grants are not enabled")
// 		} else if !bytes.Equal(feeGranter, feePayer) {
// 			if !providedFee.IsNil() {
// 				err := dfd.feegrantKeeper.UseGrantedFees(ctx, feeGranter, feePayer, sdk.NewCoins(providedFee), sdkTx.GetMsgs())
// 				if err != nil {
// 					return errorsmod.Wrapf(err, "%s does not allow to pay fees for %s", feeGranter, feePayer)
// 				}
// 			}
// 		}

// 		deductFeesFrom = feeGranter
// 	}

// 	deductFeesFromAcc := dfd.accountKeeper.GetAccount(ctx, deductFeesFrom)
// 	if deductFeesFromAcc == nil {
// 		return errors.ErrUnknownAddress.Wrapf("fee payer address: %s does not exist", deductFeesFrom)
// 	}

// 	return escrow(dfd.bankKeeper, ctx, deductFeesFromAcc, sdk.NewCoins(providedFee))
// }

// // escrow deducts coins to the escrow.
// func escrow(bankKeeper BankKeeper, ctx sdk.Context, acc authtypes.AccountI, coins sdk.Coins) error {
// 	targetModuleAcc := feemarkettypes.FeeCollectorName
// 	err := bankKeeper.SendCoinsFromAccountToModule(ctx, acc.GetAddress(), targetModuleAcc, coins)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // resolveTxPriorityCoins converts the coins to the proper denom used for tx prioritization calculation.
// func (dfd DeductFeeDecorator) resolveTxPriorityCoins(ctx sdk.Context, fee sdk.Coin, baseDenom string) (sdk.Coin, error) {
// 	if fee.Denom == baseDenom {
// 		return fee, nil
// 	}

// 	feeDec := sdk.NewDecCoinFromCoin(fee)
// 	convertedDec, err := dfd.feeMarketKeeper.ResolveToDenom(ctx, feeDec, baseDenom)
// 	if err != nil {
// 		return sdk.Coin{}, err
// 	}
// 	// truncate down
// 	return sdk.NewCoin(baseDenom, convertedDec.Amount.TruncateInt()), nil
// }

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

// const (
// 	// gasPricePrecision is the amount of digit precision to scale the gas prices to.
// 	gasPricePrecision = 6
// )

// // GetTxPriority returns a naive tx priority based on the amount of gas price provided in a transaction.
// //
// // The fee amount is divided by the gasLimit to calculate "Effective Gas Price".
// // This value is then normalized and scaled into an integer, so it can be used as a priority.
// //
// //	effectiveGasPrice = feeAmount / gas limit (denominated in fee per gas)
// //	normalizedGasPrice = effectiveGasPrice / currentGasPrice (floor is 1.  The minimum effective gas price can ever be is current gas price)
// //	scaledGasPrice = normalizedGasPrice * 10 ^ gasPricePrecision (amount of decimal places in the normalized gas price to consider when converting to int64).
// func GetTxPriority(fee sdk.Coin, gasLimit int64, currentGasPrice sdk.DecCoin) int64 {
// 	// protections from dividing by 0
// 	if gasLimit == 0 {
// 		return 0
// 	}

// 	// if the gas price is 0, just use a raw amount
// 	if currentGasPrice.IsZero() {
// 		return fee.Amount.Int64()
// 	}

// 	effectiveGasPrice := fee.Amount.ToLegacyDec().QuoInt64(gasLimit)
// 	normalizedGasPrice := effectiveGasPrice.Quo(currentGasPrice.Amount)
// 	scaledGasPrice := normalizedGasPrice.MulInt64(int64(math.Pow10(gasPricePrecision)))

// 	// overflow panic protection
// 	if scaledGasPrice.GTE(sdkmath.LegacyNewDec(math.MaxInt64)) {
// 		return math.MaxInt64
// 	} else if scaledGasPrice.LTE(sdkmath.LegacyOneDec()) {
// 		return 0
// 	}

// 	return scaledGasPrice.TruncateInt64()
// }
