package ante

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

type TxFeeChecker func(ctx sdk.Context, tx sdk.Tx) (sdk.Coins, int64, error)

// DeductFeeDecorator deducts fees from the first signer of the tx
// If the first signer does not have the funds to pay for the fees, return with InsufficientFunds error
// Call next AnteHandler if fees successfully deducted
// CONTRACT: Tx must implement FeeTx interface to use DeductFeeDecorator
type DeductFeeDecorator struct {
	accountKeeper  ante.AccountKeeper
	bankKeeper     BankKeeper
	feegrantKeeper ante.FeegrantKeeper
	didKeeper      DidKeeper
	resourceKeeper ResourceKeeper
	txFeeChecker   TxFeeChecker
}

func NewDeductFeeDecorator(ak ante.AccountKeeper, bk BankKeeper, fk ante.FeegrantKeeper, dk DidKeeper, rk ResourceKeeper, tfc TxFeeChecker) DeductFeeDecorator {
	if tfc == nil {
		tfc = checkTxFeeWithValidatorMinGasPrices
	}

	return DeductFeeDecorator{
		accountKeeper:  ak,
		bankKeeper:     bk,
		feegrantKeeper: fk,
		didKeeper:      dk,
		resourceKeeper: rk,
		txFeeChecker:   tfc,
	}
}

func (dfd DeductFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	if !simulate && ctx.BlockHeight() > 0 && feeTx.GetGas() == 0 {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidGasLimit, "must provide positive gas")
	}

	var (
		priority int64
		err      error
	)

	fee := feeTx.GetFee()
	gas := feeTx.GetGas()

	// check if the tx is a taxable tx
	// CONTRACT: Taxable tx is a tx that has at least 1 taxable related Msg.
	taxable := IsTaxableTxLite(tx)
	// if taxable, check if the fee is enough to be included in the mempool
	if taxable {
		// check if the fee is enough to be included in the mempool
		if !simulate {
			// calculate spread
			_, reward, burn := IsTaxableTx(ctx, dfd.didKeeper, dfd.resourceKeeper, tx)
			// define tax as the sum of reward and burn
			tax := reward.Add(burn...)
			_, priority, err = IsSufficientFee(ctx, tax, reward, burn, fee, int64(gas))
			if err != nil {
				return ctx, err
			}
		}
		newCtx := ctx.WithPriority(priority)
		return next(newCtx, tx, simulate)
	}

	// if not taxable, follow default
	if !simulate {
		fee, priority, err = dfd.txFeeChecker(ctx, tx)
		if err != nil {
			return ctx, err
		}
	}
	if err := dfd.checkDeductFee(ctx, tx, fee); err != nil {
		return ctx, err
	}

	newCtx := ctx.WithPriority(priority)

	return next(newCtx, tx, simulate)
}

func (dfd DeductFeeDecorator) checkDeductFee(ctx sdk.Context, sdkTx sdk.Tx, fee sdk.Coins) error {
	feeTx, ok := sdkTx.(sdk.FeeTx)
	if !ok {
		return sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	if addr := dfd.accountKeeper.GetModuleAddress(types.FeeCollectorName); addr == nil {
		return fmt.Errorf("fee collector module account (%s) has not been set", types.FeeCollectorName)
	}

	feePayer := feeTx.FeePayer()
	feeGranter := feeTx.FeeGranter()
	deductFeesFrom := feePayer

	// if feegranter set deduct fee from feegranter account.
	// this works with only when feegrant enabled.
	if feeGranter != nil {
		if dfd.feegrantKeeper == nil {
			return sdkerrors.ErrInvalidRequest.Wrap("fee grants are not enabled")
		} else if !feeGranter.Equals(feePayer) {
			err := dfd.feegrantKeeper.UseGrantedFees(ctx, feeGranter, feePayer, fee, sdkTx.GetMsgs())
			if err != nil {
				return sdkerrors.Wrapf(err, "%s does not not allow to pay fees for %s", feeGranter, feePayer)
			}
		}

		deductFeesFrom = feeGranter
	}

	deductFeesFromAcc := dfd.accountKeeper.GetAccount(ctx, deductFeesFrom)
	if deductFeesFromAcc == nil {
		return sdkerrors.ErrUnknownAddress.Wrapf("fee payer address: %s does not exist", deductFeesFrom)
	}

	// deduct the fees
	if !fee.IsZero() {
		err := DeductFees(dfd.bankKeeper, ctx, deductFeesFromAcc, fee)
		if err != nil {
			return err
		}
	}

	events := sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeTx,
			sdk.NewAttribute(sdk.AttributeKeyFee, fee.String()),
			sdk.NewAttribute(sdk.AttributeKeyFeePayer, deductFeesFrom.String()),
		),
	}
	ctx.EventManager().EmitEvents(events)

	return nil
}

// DeductFees deducts fees from the given account.
func DeductFees(bankKeeper types.BankKeeper, ctx sdk.Context, acc types.AccountI, fees sdk.Coins) error {
	if !fees.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "invalid fee amount: %s", fees)
	}

	err := bankKeeper.SendCoinsFromAccountToModule(ctx, acc.GetAddress(), types.FeeCollectorName, fees)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, err.Error())
	}

	return nil
}

func IsSufficientFee(ctx sdk.Context, tax, reward, burn, feeProvided sdk.Coins, gasRequested int64) (bool, int64, error) {
	// check if the provided fee is enough for `did`, `resource` module specific Msg
	if !feeProvided.IsAnyGTE(tax) {
		return false, 0, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "insufficient fees; got: %s required: %s", feeProvided, tax)
	}

	// check with the default validator min gas prices based on rewards distribution
	if ctx.IsCheckTx() {
		minGasPrices := ctx.MinGasPrices()
		if !minGasPrices.IsZero() {
			requiredFees := make(sdk.Coins, len(minGasPrices))

			// Determine the required fees by multiplying each required minimum gas
			// price by the gas limit, where fee = ceil(minGasPrice * gasLimit).
			glDec := sdk.NewDec(int64(gasRequested))
			for i, gp := range minGasPrices {
				fee := gp.Amount.Mul(glDec)
				requiredFees[i] = sdk.NewCoin(gp.Denom, fee.Ceil().RoundInt())
			}

			// check if the fee is sufficient
			if !reward.IsAnyGTE(requiredFees) {
				return false, 0, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "insufficient fees; got: %s required: %s", reward, requiredFees)
			}
		}
	}

	priority := getTxPriority(tax, gasRequested)

	return true, priority, nil
}
