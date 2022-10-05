package ante

// 1. Implement `txFeeChecker` to wrap the default fee checker along with the custom fee logic
//    a. `IsIdentityTx` returns true if the tx is an identity tx (at least 1 create, update)
//        	i. 	If true, calculate the number of identity related Msgs and the total fee
//    		ii.	If false, handle the fee logic as default with `checkTxFeeWithValidatorMinGasPrices`
//    b. `checkTxFeeWithCustomFixedFee` is called second to check the custom fee logic
// 2. Override `checkDeductFee` to wrap `DeductFees` along with the custom fee logic
//    a. `DeductFees` is called first as default fee deduction from the fee payer to the corresponding module account (where the Msg is routed -> `cheqd` or `resource` currently)
//    b. `DeductFeesEvent` is emitted for fee deduction from fee payer
//    c. `BurnFees` is called to burn portion of the fee from the module account
//    c. `BurnFeesEvent` is emitted for a predefined portion of identity fees (deducted from module account)
//    d. `DistributeToFoundation` is called to distribute the predefined amount to the foundation account
//    e. `DistributeToFoundationEvent` is emitted for the predefined amount of identity fees (deducted from module account)
//    f. `DistributeToFeeCollector` is called to distribute the predefined amount to the fee collector account
//    g. `DistributeToFeeCollectorEvent` is emitted for the predefined amount of identity fees (deducted from module account)

import (
	"fmt"

	cheqdtypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// TxFeeChecker check if the provided fee is enough and returns the effective fee and tx priority,
// the effective fee should be deducted later, and the priority should be returned in abci response.
// Here the default fee checker type is augmented to return a boolean flag to indicate if the tx is an identity tx,
// hence a custom fee logic is applied and the total custom fee is returned as well.
type TxFeeChecker func(ctx sdk.Context, cheqdKeeper CheqdKeeper, resourceKeeper ResourceKeeper, tx sdk.Tx) (sdk.Coins, sdk.Coins, int64, bool, error)

// DeductFeeDecorator deducts fees from the first signer of the tx
// If the first signer does not have the funds to pay for the fees, return with InsufficientFunds error
// Call next AnteHandler if fees successfully deducted
// CONTRACT: Tx must implement FeeTx interface to use DeductFeeDecorator
type DeductFeeDecorator struct {
	accountKeeper  ante.AccountKeeper
	bankKeeper     BankKeeper
	feegrantKeeper ante.FeegrantKeeper
	cheqdKeeper    CheqdKeeper
	resourceKeeper ResourceKeeper
	txFeeChecker   TxFeeChecker
}

func NewDeductFeeDecorator(ak ante.AccountKeeper, bk BankKeeper, fk ante.FeegrantKeeper, ck CheqdKeeper, rk ResourceKeeper, tfc TxFeeChecker) DeductFeeDecorator {
	if tfc == nil {
		tfc = checkTxFeeWithCustomFixedFee
	}

	return DeductFeeDecorator{
		accountKeeper:  ak,
		bankKeeper:     bk,
		feegrantKeeper: fk,
		cheqdKeeper:    ck,
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
		priority    int64
		isCustomFee bool
		err         error
	)

	fee := feeTx.GetFee()
	burn := (sdk.Coins)(nil)
	// CONTRACT: simulate=true means the fee is ignored at this point, as we don't support a fixed fee for simulation mode for `TaxableMsg` Tx.
	// As a result, a lite version of `IsIdentityTx` is called to check if the tx is an identity tx, as the most minimal check.
	// If true, an *error* bubbles up to indicate that the tx is an identity tx and a fixed fee is not supported in simulation mode.
	if !simulate {
		fee, burn, priority, isCustomFee, err = dfd.txFeeChecker(ctx, dfd.cheqdKeeper, dfd.resourceKeeper, tx)
		if err != nil {
			return ctx, err
		}
	} else {
		isCustomFee = IsIdentityTxLite(feeTx)
	}

	if err := dfd.checkDeductFeeWithFixedFee(ctx, tx, isCustomFee, fee, burn, simulate); err != nil {
		return ctx, err
	}

	newCtx := ctx.WithPriority(priority)

	return next(newCtx, tx, simulate)
}

func (dfd DeductFeeDecorator) checkDeductFeeWithFixedFee(ctx sdk.Context, sdkTx sdk.Tx, isCustomFee bool, fee sdk.Coins, burnFeePortion sdk.Coins, simulate bool) error {
	feeTx, ok := sdkTx.(sdk.FeeTx)
	if !ok {
		return sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	if !isCustomFee {
		return dfd.checkDeductFee(ctx, feeTx, fee)
	}

	if simulate {
		return sdkerrors.Wrap(sdkerrors.ErrNotSupported, "simulation of fees is not supported for network specific transactions")
	}

	if addr := dfd.accountKeeper.GetModuleAddress(cheqdtypes.ModuleName); addr == nil {
		return fmt.Errorf("cheqd fee collector module account (%s) has not been set", cheqdtypes.ModuleName)
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
		err := DeductFeesToModule(dfd.bankKeeper, ctx, deductFeesFromAcc, fee, cheqdtypes.ModuleName)
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

	distrFeeAlloc, err := GetDistributionFee(ctx, fee, burnFeePortion)
	if err != nil {
		return err
	}

	// 1. Fixed fee resides in the cheqd module account now
	//		a. Predefined amount of fixed fee is burnt
	//		b. Predefined amount of fixed fee is distributed to the validators as rewards

	// 1a. Burn fixed fee
	if err := BurnFee(dfd.bankKeeper, ctx, distrFeeAlloc[BurnFeePortion]); err != nil {
		return err
	}

	// 1b. Distribute fixed fee to the validators as rewards
	if err := DistributeFeeToModule(dfd.bankKeeper, ctx, distrFeeAlloc[RewardsFeePortion]); err != nil {
		return err
	}

	return nil
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

// DeductFees deducts fees from the given account.
// Used for `cheqd` fee module, but can accept any valid module with an instantiated account.
// NOTE: Validation of the instantiated module account is done before calling this function on `checkDeductFeeWithFixedFee`.
func DeductFeesToModule(bankKeeper types.BankKeeper, ctx sdk.Context, acc types.AccountI, fees sdk.Coins, recipientModule string) error {
	if !fees.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "invalid fee amount: %s", fees)
	}

	err := bankKeeper.SendCoinsFromAccountToModule(ctx, acc.GetAddress(), recipientModule, fees)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, err.Error())
	}

	return nil
}

// checkTxWithCustomMinGasPrices checks if the provided fee is enough for `cheqd`, `resource` module specific Msg and returns the effective fee and tx priority,
// otherwise it returns the default fee and priority.
func checkTxFeeWithCustomFixedFee(ctx sdk.Context, cheqdKeeper CheqdKeeper, resourceKeeper ResourceKeeper, tx sdk.Tx) (sdk.Coins, sdk.Coins, int64, bool, error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return nil, nil, 0, false, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	// check if the provided fee is enough for `cheqd`, `resource` module specific Msg
	isIdentityTx, fee, burn := IsIdentityTx(ctx, cheqdKeeper, resourceKeeper, feeTx)
	if !isIdentityTx {
		return checkTxFeeWithValidatorMinGasPrices(ctx, feeTx)
	}

	feeCoins := feeTx.GetFee()
	gas := feeTx.GetGas()

	// Ensure that the provided fees meet a minimum threshold for the validator,
	// if this is a CheckTx. This is only for local mempool purposes, and thus
	// is only ran on check tx.
	if ctx.IsCheckTx() || ctx.IsReCheckTx() {
		_, err := IsSufficientCustomFee(ctx, fee, feeCoins, burn, int64(gas))
		if err != nil {
			return nil, nil, 0, false, err
		}
	}

	priority := getTxPriority(feeCoins, int64(gas))
	return fee, burn, priority, true, nil
}

func IsSufficientCustomFee(ctx sdk.Context, feeRequired sdk.Coins, fee sdk.Coins, burnFeePortion sdk.Coins, gasRequested int64) (DistributionFeeAllocation, error) {
	// 1. Check if the provided fee is enough for `cheqd`, `resource` module specific Msg
	// 2. Calculate further distributions
	// 3. Check with the default validator min gas prices based on rewards distribution
	// 4. If the fee is not sufficient, return error

	// 1. Check if the provided fee is enough for `cheqd`, `resource` module specific Msg
	if !fee.IsAnyGTE(feeRequired) {
		return DistributionFeeAllocation{}, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "insufficient fees; got: %s required: %s", fee, feeRequired)
	}

	// 2. Calculate further distributions
	// TODO: At any case, we need to also decide on handling reverting txs if the dynamic validation fails (e.g. if the identity signatures are invalid, verification method is not supported yet, etc).
	// TODO: This is achieved with defining a `postHandler` in the `anteHandler` and reverting the tx if the dynamic validation fails.
	//* NOTE: Here we are accepting the total fee provided by the user, if it is greater than or equal to the required fee.
	distrFeeAlloc, err := GetDistributionFee(ctx, feeRequired, burnFeePortion)
	if err != nil {
		return distrFeeAlloc, err
	}

	// 3. Check with the default validator min gas prices based on rewards distribution
	rewardsFee := distrFeeAlloc[RewardsFeePortion]
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

		// 4. If the fee is not sufficient, return error
		if !rewardsFee.IsAnyGTE(requiredFees) {
			return distrFeeAlloc, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "insufficient fees; got: %s required: %s", rewardsFee, requiredFees)
		}
	}

	return distrFeeAlloc, nil
}
