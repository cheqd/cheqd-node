package posthandler

import (
	"fmt"

	cheqdante "github.com/cheqd/cheqd-node/ante"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// taxDecorator will handle tax for all taxable messages
type taxDecorator struct {
	accountKeeper  ante.AccountKeeper
	bankKeeper     cheqdante.BankKeeper
	feegrantKeeper ante.FeegrantKeeper
	didKeeper      cheqdante.DidKeeper
	resourceKeeper cheqdante.ResourceKeeper
}

// NewTaxDecorator returns a new taxDecorator
func NewTaxDecorator(ak ante.AccountKeeper, bk cheqdante.BankKeeper, fk ante.FeegrantKeeper, dk cheqdante.DidKeeper, rk cheqdante.ResourceKeeper) taxDecorator {
	return taxDecorator{
		accountKeeper:  ak,
		bankKeeper:     bk,
		feegrantKeeper: fk,
		didKeeper:      dk,
		resourceKeeper: rk,
	}
}

// AnteHandle handles tax for all taxable messages
func (td taxDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// must implement FeeTx
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrTxDecode, "invalid transaction type: %T, must implement FeeTx", tx)
	}
	// if simulate, perform no-op
	if simulate {
		return next(ctx, tx, simulate)
	}
	// get metrics for tax
	rewards, burn, taxable, err := td.isTaxable(ctx, feeTx)
	if err != nil {
		return ctx, err
	}
	// if not taxable, skip
	if !taxable {
		return next(ctx, tx, simulate)
	}
	// get fee payer and check if fee grant exists
	tax := rewards.Add(burn...)
	feePayer, err := td.getFeePayer(ctx, feeTx, tax, tx.GetMsgs())
	if err != nil {
		return ctx, err
	}
	// deduct tax (rewards + burn) from fee payer to did module account
	if err := td.deductTaxFromFeePayer(ctx, feePayer, tax); err != nil {
		return ctx, err
	}
	// move rewards to fee collector to follow the default proposer logic
	if err := td.distributeRewards(ctx, rewards); err != nil {
		return ctx, err
	}
	// finally, burn tax portion from did module account
	if err := td.burnFees(ctx, burn); err != nil {
		return ctx, err
	}

	return next(ctx, tx, simulate)
}

// isTaxable returns true if the message is taxable and returns
func (td taxDecorator) isTaxable(ctx sdk.Context, sdkTx sdk.Tx) (rewards sdk.Coins, burn sdk.Coins, taxable bool, err error) {
	feeTx, ok := sdkTx.(sdk.FeeTx)
	if !ok {
		return sdk.Coins{}, sdk.Coins{}, false, sdkerrors.Wrapf(sdkerrors.ErrTxDecode, "invalid transaction type: %T, must implement FeeTx", sdkTx)
	}
	// run lite validation
	taxable = cheqdante.IsTaxableTxLite(feeTx)
	if taxable {
		// run full validation
		_, rewards, burn = cheqdante.IsTaxableTx(ctx, td.didKeeper, td.resourceKeeper, feeTx)
		return rewards, burn, taxable, nil
	}

	return rewards, burn, taxable, err
}

// getFeePayer returns the fee payer and checks if a fee grant exists
func (td taxDecorator) getFeePayer(ctx sdk.Context, feeTx sdk.FeeTx, tax sdk.Coins, msgs []sdk.Msg) (types.AccountI, error) {
	feePayer := feeTx.FeePayer()
	feeGranter := feeTx.FeeGranter()
	deductFrom := feePayer
	if feeGranter != nil {
		// check if fee grant is supported
		if td.feegrantKeeper == nil {
			return nil, sdkerrors.ErrInvalidRequest.Wrap("fee grants are not enabled")
		} else if !feeGranter.Equals(feePayer) {
			// check if fee grant exists
			err := td.feegrantKeeper.UseGrantedFees(ctx, feeGranter, feePayer, tax, msgs)
			if err != nil {
				return nil, sdkerrors.Wrapf(err, "%s does not not allow to pay fees for %s", feeGranter, feePayer)
			}
		}
		deductFrom = feeGranter
	}

	deductFromAcc := td.accountKeeper.GetAccount(ctx, deductFrom)
	if deductFromAcc == nil {
		return nil, sdkerrors.ErrUnknownAddress.Wrapf("fee payer address: %s does not exist", deductFrom)
	}

	return deductFromAcc, nil
}

// deductTaxFromFeePayer deducts fees from the account
func (td taxDecorator) deductTaxFromFeePayer(ctx sdk.Context, acc types.AccountI, fees sdk.Coins) error {
	// ensure module account has been set
	if addr := td.accountKeeper.GetModuleAddress(didtypes.ModuleName); addr == nil {
		return fmt.Errorf("cheqd fee collector module account (%s) has not been set", didtypes.ModuleName)
	}
	// deduct fees to did module account
	err := td.bankKeeper.SendCoinsFromAccountToModule(ctx, acc.GetAddress(), didtypes.ModuleName, fees)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "failed to deduct fees from %s: %s", acc.GetAddress(), err)
	}
	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeTx,
				sdk.NewAttribute(sdk.AttributeKeyFee, fees.String()),
				sdk.NewAttribute(sdk.AttributeKeyFeePayer, acc.GetAddress().String()),
			),
		},
	)
	return nil
}

// distributeRewards distributes rewards to the fee collector
func (td taxDecorator) distributeRewards(ctx sdk.Context, rewards sdk.Coins) error {
	// move rewards to fee collector
	err := td.bankKeeper.SendCoinsFromModuleToModule(ctx, didtypes.ModuleName, types.FeeCollectorName, rewards)
	if err != nil {
		return err
	}
	return nil
}

// burnFees burns fees from the module account
func (td taxDecorator) burnFees(ctx sdk.Context, fees sdk.Coins) error {
	// burn fees
	err := td.bankKeeper.BurnCoins(ctx, didtypes.ModuleName, fees)
	if err != nil {
		return err
	}
	return nil
}
