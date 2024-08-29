package ante

import (
	"bytes"
	"math"

	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	feemarkettypes "github.com/skip-mev/feemarket/x/feemarket/types"
)

type TxFeeChecker func(ctx sdk.Context, tx sdk.Tx) (sdk.Coins, int64, error)

type DeductFeeDecorator struct {
	accountKeeper   AccountKeeper
	bankKeeper      BankKeeper
	feegrantKeeper  ante.FeegrantKeeper
	txFeeChecker    TxFeeChecker
	feeMarketKeeper FeeMarketKeeper
}

func NewDeductFeeDecorator(ak AccountKeeper, bk BankKeeper, fk ante.FeegrantKeeper, tfc TxFeeChecker, fmk FeeMarketKeeper) DeductFeeDecorator {
	if tfc == nil {
		tfc = checkTxFeeWithValidatorMinGasPrices
	}

	return DeductFeeDecorator{
		accountKeeper:   ak,
		bankKeeper:      bk,
		feegrantKeeper:  fk,
		txFeeChecker:    tfc,
		feeMarketKeeper: fmk,
	}
}

func (dfd DeductFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
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
	params, err := dfd.feeMarketKeeper.GetParams(ctx)
	if err != nil {
		return ctx, err
	}
	if params.Enabled {

		if ctx.BlockHeight() == 0 {
			return next(ctx, tx, simulate)
		}

		feeTx, ok := tx.(sdk.FeeTx)
		if !ok {
			return ctx, errorsmod.Wrap(errors.ErrTxDecode, "Tx must be a FeeTx")
		}

		if !simulate && ctx.BlockHeight() > 0 && feeTx.GetGas() == 0 {
			return ctx, errors.ErrInvalidGasLimit.Wrapf("must provide positive gas")
		}

		params, err := dfd.feeMarketKeeper.GetParams(ctx)
		if err != nil {
			return ctx, errorsmod.Wrapf(err, "unable to get fee market params")
		}

		// return if disabled
		if !params.Enabled {
			return next(ctx, tx, simulate)
		}

		feeCoins := feeTx.GetFee()
		gas := feeTx.GetGas() // use provided gas limit

		if len(feeCoins) == 0 && !simulate {
			return ctx, errorsmod.Wrapf(feemarkettypes.ErrNoFeeCoins, "got length %d", len(feeCoins))
		}
		if len(feeCoins) > 1 {
			return ctx, errorsmod.Wrapf(feemarkettypes.ErrTooManyFeeCoins, "got length %d", len(feeCoins))
		}

		var feeCoin sdk.Coin
		if simulate && len(feeCoins) == 0 {
			// if simulating and user did not provider a fee - create a dummy value for them
			feeCoin = sdk.NewCoin(params.FeeDenom, sdkmath.OneInt())
		} else {
			feeCoin = feeCoins[0]
		}
		feeGas := int64(feeTx.GetGas())

		minGasPrice, err := dfd.feeMarketKeeper.GetMinGasPrice(ctx, feeCoin.GetDenom())
		if err != nil {
			return ctx, errorsmod.Wrapf(err, "unable to get min gas price for denom %s", feeCoin.GetDenom())
		}

		ctx.Logger().Info("fee deduct ante handle",
			"min gas prices", minGasPrice,
			"fee", feeCoins,
			"gas limit", gas,
		)

		ctx = ctx.WithMinGasPrices(sdk.NewDecCoins(minGasPrice))

		if !simulate {
			_, _, err := CheckTxFee(ctx, minGasPrice, feeCoin, feeGas, true)
			if err != nil {
				return ctx, errorsmod.Wrapf(err, "error checking fee")
			}
		}

		// escrow the entire amount that the account provided as fee (feeCoin)
		err = dfd.EscrowFunds(ctx, tx, feeCoin)
		if err != nil {
			return ctx, errorsmod.Wrapf(err, "error escrowing funds")
		}

		priorityFee, err := dfd.resolveTxPriorityCoins(ctx, feeCoin, params.FeeDenom)
		if err != nil {
			return ctx, errorsmod.Wrapf(err, "error resolving fee priority")
		}

		baseGasPrice, err := dfd.feeMarketKeeper.GetMinGasPrice(ctx, params.FeeDenom)
		if err != nil {
			return ctx, err
		}

		ctx = ctx.WithPriority(GetTxPriority(priorityFee, int64(gas), baseGasPrice))

		return next(ctx, tx, simulate)
	}

	return next(ctx, tx, simulate)
}

// The actual fee is deducted in the post handler along with the tip.
func (dfd DeductFeeDecorator) EscrowFunds(ctx sdk.Context, sdkTx sdk.Tx, providedFee sdk.Coin) error {
	feeTx, ok := sdkTx.(sdk.FeeTx)
	if !ok {
		return errorsmod.Wrap(errors.ErrTxDecode, "Tx must be a FeeTx")
	}

	feePayer := feeTx.FeePayer()
	feeGranter := feeTx.FeeGranter()
	deductFeesFrom := feePayer

	// if feegranter set deduct fee from feegranter account.
	// this works with only when feegrant enabled.
	if feeGranter != nil {
		if dfd.feegrantKeeper == nil {
			return errors.ErrInvalidRequest.Wrap("fee grants are not enabled")
		} else if !bytes.Equal(feeGranter, feePayer) {
			if !providedFee.IsNil() {
				err := dfd.feegrantKeeper.UseGrantedFees(ctx, feeGranter, feePayer, sdk.NewCoins(providedFee), sdkTx.GetMsgs())
				if err != nil {
					return errorsmod.Wrapf(err, "%s does not allow to pay fees for %s", feeGranter, feePayer)
				}
			}
		}

		deductFeesFrom = feeGranter
	}

	deductFeesFromAcc := dfd.accountKeeper.GetAccount(ctx, deductFeesFrom)
	if deductFeesFromAcc == nil {
		return errors.ErrUnknownAddress.Wrapf("fee payer address: %s does not exist", deductFeesFrom)
	}

	return escrow(dfd.bankKeeper, ctx, deductFeesFromAcc, sdk.NewCoins(providedFee))
}

// escrow deducts coins to the escrow.
func escrow(bankKeeper BankKeeper, ctx sdk.Context, acc authtypes.AccountI, coins sdk.Coins) error {
	targetModuleAcc := feemarkettypes.FeeCollectorName
	err := bankKeeper.SendCoinsFromAccountToModule(ctx, acc.GetAddress(), targetModuleAcc, coins)
	if err != nil {
		return err
	}
	return nil
}

// resolveTxPriorityCoins converts the coins to the proper denom used for tx prioritization calculation.
func (dfd DeductFeeDecorator) resolveTxPriorityCoins(ctx sdk.Context, fee sdk.Coin, baseDenom string) (sdk.Coin, error) {
	if fee.Denom == baseDenom {
		return fee, nil
	}

	feeDec := sdk.NewDecCoinFromCoin(fee)
	convertedDec, err := dfd.feeMarketKeeper.ResolveToDenom(ctx, feeDec, baseDenom)
	if err != nil {
		return sdk.Coin{}, err
	}
	// truncate down
	return sdk.NewCoin(baseDenom, convertedDec.Amount.TruncateInt()), nil
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

const (
	// gasPricePrecision is the amount of digit precision to scale the gas prices to.
	gasPricePrecision = 6
)

// GetTxPriority returns a naive tx priority based on the amount of gas price provided in a transaction.
//
// The fee amount is divided by the gasLimit to calculate "Effective Gas Price".
// This value is then normalized and scaled into an integer, so it can be used as a priority.
//
//	effectiveGasPrice = feeAmount / gas limit (denominated in fee per gas)
//	normalizedGasPrice = effectiveGasPrice / currentGasPrice (floor is 1.  The minimum effective gas price can ever be is current gas price)
//	scaledGasPrice = normalizedGasPrice * 10 ^ gasPricePrecision (amount of decimal places in the normalized gas price to consider when converting to int64).
func GetTxPriority(fee sdk.Coin, gasLimit int64, currentGasPrice sdk.DecCoin) int64 {
	// protections from dividing by 0
	if gasLimit == 0 {
		return 0
	}

	// if the gas price is 0, just use a raw amount
	if currentGasPrice.IsZero() {
		return fee.Amount.Int64()
	}

	effectiveGasPrice := fee.Amount.ToLegacyDec().QuoInt64(gasLimit)
	normalizedGasPrice := effectiveGasPrice.Quo(currentGasPrice.Amount)
	scaledGasPrice := normalizedGasPrice.MulInt64(int64(math.Pow10(gasPricePrecision)))

	// overflow panic protection
	if scaledGasPrice.GTE(sdkmath.LegacyNewDec(math.MaxInt64)) {
		return math.MaxInt64
	} else if scaledGasPrice.LTE(sdkmath.LegacyOneDec()) {
		return 0
	}

	return scaledGasPrice.TruncateInt64()
}
