package posthandler

import (
	"bytes"
	"fmt"

	errors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	cheqdante "github.com/cheqd/cheqd-node/ante"
	"github.com/cheqd/cheqd-node/pricefeeder"
	"github.com/cheqd/cheqd-node/util"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	oraclekeeper "github.com/cheqd/cheqd-node/x/oracle/keeper"
	oracletypes "github.com/cheqd/cheqd-node/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	globalfeekeeper "github.com/noble-assets/globalfee/keeper"
	feeabskeeper "github.com/osmosis-labs/fee-abstraction/v8/x/feeabs/keeper"
	feeabstypes "github.com/osmosis-labs/fee-abstraction/v8/x/feeabs/types"
	feemarkettypes "github.com/skip-mev/feemarket/x/feemarket/types"
)

// TaxDecorator will handle tax for all taxable messages
type TaxDecorator struct {
	accountKeeper     ante.AccountKeeper
	bankKeeper        BankKeeper
	feegrantKeeper    ante.FeegrantKeeper
	didKeeper         cheqdante.DidKeeper
	resourceKeeper    cheqdante.ResourceKeeper
	feemarketKeeper   FeeMarketKeeper
	oracleKeeper      cheqdante.OracleKeeper
	feeabsKeeper      feeabskeeper.Keeper
	oraclePricefeeder *pricefeeder.PriceFeeder
	globalFeeKeeper   *globalfeekeeper.Keeper
}

// NewTaxDecorator returns a new taxDecorator
func NewTaxDecorator(ak ante.AccountKeeper, bk BankKeeper, fk ante.FeegrantKeeper, dk cheqdante.DidKeeper, rk cheqdante.ResourceKeeper, fmk FeeMarketKeeper, ok cheqdante.OracleKeeper, fak feeabskeeper.Keeper, pf *pricefeeder.PriceFeeder, gfk *globalfeekeeper.Keeper) TaxDecorator {
	return TaxDecorator{
		accountKeeper:     ak,
		bankKeeper:        bk,
		feegrantKeeper:    fk,
		didKeeper:         dk,
		resourceKeeper:    rk,
		feemarketKeeper:   fmk,
		oracleKeeper:      ok,
		feeabsKeeper:      fak,
		oraclePricefeeder: pf,
		globalFeeKeeper:   gfk,
	}
}

// AnteHandle handles tax for all taxable messages
func (td TaxDecorator) PostHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, success bool, next sdk.PostHandler) (sdk.Context, error) {
	params, err := td.feemarketKeeper.GetParams(ctx)
	if err != nil {
		return ctx, err
	}

	nativeDenom := params.FeeDenom
	// must implement FeeTx
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, errors.Wrapf(sdkerrors.ErrTxDecode, "invalid transaction type: %T, must implement FeeTx", tx)
	}
	// if simulate, perform no-op
	if simulate {
		return next(ctx, tx, simulate, success)
	}
	// get metrics for tax
	rewards, burn, taxable, err := td.isTaxable(ctx, feeTx)
	if err != nil {
		return ctx, err
	}
	// if bypassable, perform no-op
	if cheqdante.ShouldBypassFeeMarket(ctx, td.globalFeeKeeper, tx) {
		return next(ctx, tx, simulate, success)
	}

	if taxable {
		err := td.handleTaxableTransaction(ctx, feeTx, simulate, rewards, burn, tx)
		if err != nil {
			return ctx, err
		}
		return next(ctx, tx, simulate, success)
	}
	params, err = td.feemarketKeeper.GetParams(ctx)
	if err != nil {
		return ctx, errors.Wrapf(err, "unable to get fee market params")
	}
	// return if disabled
	if !params.Enabled {
		return next(ctx, tx, simulate, success)
	}

	enabledHeight, err := td.feemarketKeeper.GetEnabledHeight(ctx)
	if err != nil {
		return ctx, errors.Wrapf(err, "unable to get fee market enabled height")
	}

	// if the current height is that which enabled the feemarket or lower, skip deduction
	if ctx.BlockHeight() <= enabledHeight {
		return next(ctx, tx, simulate, success)
	}

	// update fee market state
	state, err := td.feemarketKeeper.GetState(ctx)
	if err != nil {
		return ctx, errors.Wrapf(err, "unable to get fee market state")
	}

	onlyNativeDenom := true
	for _, fee := range feeTx.GetFee() {
		if fee.Denom != nativeDenom {
			// If any other token besides the native denom is present, set the flag to false
			onlyNativeDenom = false
			break
		}
	}
	var feeCoins sdk.Coins

	// if IBC Denom fetch the balance of did module.
	if onlyNativeDenom {
		feeCoins = feeTx.GetFee()
	} else {
		addr := td.accountKeeper.GetModuleAddress(feemarkettypes.FeeCollectorName)
		feeBal := td.bankKeeper.GetBalance(ctx, addr, nativeDenom)
		feeCoins = sdk.NewCoins(feeBal)
	}
	gas := ctx.GasMeter().GasConsumed() // use context gas consumed

	if len(feeCoins) == 0 && !simulate {
		return ctx, errors.Wrapf(feemarkettypes.ErrNoFeeCoins, "got length %d", len(feeCoins))
	}
	if len(feeCoins) > 1 {
		return ctx, errors.Wrapf(feemarkettypes.ErrTooManyFeeCoins, "got length %d", len(feeCoins))
	}

	var feeCoin sdk.Coin
	if simulate && len(feeCoins) == 0 {
		// if simulating and user did not provider a fee - create a dummy value for them
		feeCoin = sdk.NewCoin(params.FeeDenom, math.OneInt())
	} else {
		feeCoin = feeCoins[0]
	}

	feeGas := int64(feeTx.GetGas())

	var (
		tip     = sdk.NewCoin(feeCoin.Denom, math.ZeroInt())
		payCoin = feeCoin
	)

	minGasPrice, err := td.feemarketKeeper.GetMinGasPrice(ctx, feeCoin.GetDenom())
	if err != nil {
		return ctx, errors.Wrapf(err, "unable to get min gas price for denom %s", feeCoins[0].GetDenom())
	}

	ctx.Logger().Info("fee deduct post handle",
		"min gas prices", minGasPrice,
		"gas consumed", gas,
	)

	if !simulate {
		payCoin, tip, err = cheqdante.CheckTxFee(ctx, minGasPrice, feeCoin, feeGas, false)
		if err != nil {
			return ctx, err
		}
	}

	ctx.Logger().Info("fee deduct post handle",
		"fee", payCoin,
		"tip", tip,
	)
	if err := td.PayOutFeeAndTip(ctx, payCoin, tip); err != nil {
		return ctx, err
	}
	err = state.Update(gas, params)
	if err != nil {
		return ctx, errors.Wrapf(err, "unable to update fee market state")
	}

	err = td.feemarketKeeper.SetState(ctx, state)
	if err != nil {
		return ctx, errors.Wrapf(err, "unable to set fee market state")
	}

	return next(ctx, tx, simulate, success)
}

// PayOutFeeAndTip deducts the provided fee and tip from the fee payer.
// If the tx uses a feegranter, the fee granter address will pay the fee instead of the tx signer.
func (td TaxDecorator) PayOutFeeAndTip(ctx sdk.Context, fee, tip sdk.Coin) error {
	params, err := td.feemarketKeeper.GetParams(ctx)
	if err != nil {
		return fmt.Errorf("error getting feemarket params: %v", err)
	}

	var events sdk.Events
	// deduct the fees and tip
	if !fee.IsNil() {
		err := DeductCoins(td.bankKeeper, ctx, sdk.NewCoins(fee), params.DistributeFees)
		if err != nil {
			return err
		}

		events = append(events, sdk.NewEvent(
			feemarkettypes.EventTypeFeePay,
			sdk.NewAttribute(sdk.AttributeKeyFee, fee.String()),
		))
	}

	proposer := sdk.AccAddress(ctx.BlockHeader().ProposerAddress)
	if !tip.IsNil() {
		err := SendTip(td.bankKeeper, ctx, proposer, sdk.NewCoins(tip))
		if err != nil {
			return err
		}

		events = append(events, sdk.NewEvent(
			feemarkettypes.EventTypeTipPay,
			sdk.NewAttribute(feemarkettypes.AttributeKeyTip, tip.String()),
			sdk.NewAttribute(feemarkettypes.AttributeKeyTipPayee, proposer.String()),
		))
	}

	ctx.EventManager().EmitEvents(events)
	return nil
}

// DeductCoins deducts coins from the given account.
// Coins can be sent to the default fee collector (
// causes coins to be distributed to stakers) or kept in the fee collector account (soft burn).
func DeductCoins(bankKeeper BankKeeper, ctx sdk.Context, coins sdk.Coins, distributeFees bool) error {
	if distributeFees {
		err := bankKeeper.SendCoinsFromModuleToModule(ctx, feemarkettypes.FeeCollectorName, types.FeeCollectorName, coins)
		if err != nil {
			return err
		}
	}
	return nil
}

// SendTip sends a tip to the current block proposer.
func SendTip(bankKeeper BankKeeper, ctx sdk.Context, proposer sdk.AccAddress, coins sdk.Coins) error {
	err := bankKeeper.SendCoinsFromModuleToAccount(ctx, feemarkettypes.FeeCollectorName, proposer, coins)
	if err != nil {
		return err
	}

	return nil
}

// isTaxable returns true if the message is taxable and returns
func (td TaxDecorator) isTaxable(ctx sdk.Context, sdkTx sdk.Tx) (rewards sdk.Coins, burn sdk.Coins, taxable bool, err error) {
	feeTx, ok := sdkTx.(sdk.FeeTx)
	if !ok {
		return sdk.Coins{}, sdk.Coins{}, false, errors.Wrapf(sdkerrors.ErrTxDecode, "invalid transaction type: %T, must implement FeeTx", sdkTx)
	}
	// run lite validation
	taxable = cheqdante.IsTaxableTxLite(feeTx)
	if taxable {
		// run full validation
		_, rewards, burn, err = cheqdante.IsTaxableTx(ctx, td.didKeeper, td.resourceKeeper, feeTx, td.oracleKeeper, td.feeabsKeeper, td.oraclePricefeeder)
		return rewards, burn, taxable, err
	}

	return rewards, burn, taxable, err
}

// getFeePayer returns the fee payer and checks if a fee grant exists
func (td TaxDecorator) getFeePayer(ctx sdk.Context, feeTx sdk.FeeTx, tax sdk.Coins, msgs []sdk.Msg) (sdk.AccountI, error) {
	feePayer := feeTx.FeePayer()
	feeGranter := feeTx.FeeGranter()
	deductFrom := feePayer
	if feeGranter != nil {
		// check if fee grant is supported
		if td.feegrantKeeper == nil {
			return nil, sdkerrors.ErrInvalidRequest.Wrap("fee grants are not enabled")
		} else if !bytes.Equal(feeGranter, feePayer) {
			// check if fee grant exists
			err := td.feegrantKeeper.UseGrantedFees(ctx, feeGranter, feePayer, tax, msgs)
			if err != nil {
				return nil, errors.Wrapf(err, "%s does not not allow to pay fees for %s", feeGranter, feePayer)
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

func (td TaxDecorator) validateTax(tax sdk.Coins, simulate bool) error {
	// no-op if simulate
	if simulate {
		return nil
	}
	// check if denom is accepted
	if !tax.DenomsSubsetOf(sdk.NewCoins(sdk.NewCoin(didtypes.BaseMinimalDenom, math.NewInt(1)))) {
		return errors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid denom: %s", tax)
	}
	// check if tax is positive
	if !tax.IsAllPositive() {
		return errors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid tax: %s", tax)
	}
	return nil
}

// deductTaxFromFeePayer deducts fees from the account
func (td TaxDecorator) deductTaxFromFeePayer(ctx sdk.Context, acc sdk.AccountI, fees sdk.Coins) error {
	// ensure module account has been set
	if addr := td.accountKeeper.GetModuleAddress(didtypes.ModuleName); addr == nil {
		return fmt.Errorf("cheqd fee collector module account (%s) has not been set", didtypes.ModuleName)
	}
	// deduct fees to did module account
	err := td.bankKeeper.SendCoinsFromAccountToModule(ctx, acc.GetAddress(), didtypes.ModuleName, fees)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInsufficientFunds, "failed to deduct fees from %s: %s", acc.GetAddress(), err)
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
func (td TaxDecorator) distributeRewards(ctx sdk.Context, rewards sdk.Coins) error {
	if rewards.IsZero() {
		return nil
	}

	oracleShareRate := math.LegacyNewDecFromIntWithPrec(math.NewInt(5), 3) // 0.005 = 0.5%
	oracleRewards, feeCollectorRewards := SplitRewardsByRatio(rewards, oracleShareRate)

	if !oracleRewards.IsZero() {
		if err := td.bankKeeper.SendCoinsFromModuleToModule(ctx, didtypes.ModuleName, oracletypes.ModuleName, oracleRewards); err != nil {
			return err
		}
	}
	if !feeCollectorRewards.IsZero() {
		if err := td.bankKeeper.SendCoinsFromModuleToModule(ctx, didtypes.ModuleName, types.FeeCollectorName, feeCollectorRewards); err != nil {
			return err
		}
	}
	return nil
}

// SplitRewardsByRatio splits the input rewards by a share ratio.
// It returns two sdk.Coins:  and remainingShare.
func SplitRewardsByRatio(rewards sdk.Coins, ratio math.LegacyDec) (oracleRewards sdk.Coins, feeCollectorRewards sdk.Coins) {
	oracleRewards = sdk.NewCoins()
	feeCollectorRewards = sdk.NewCoins()

	for _, coin := range rewards {
		portion := ratio.MulInt(coin.Amount).TruncateInt()
		rest := coin.Amount.Sub(portion)

		if portion.IsPositive() {
			oracleRewards = oracleRewards.Add(sdk.NewCoin(coin.Denom, portion))
		}
		if rest.IsPositive() {
			feeCollectorRewards = feeCollectorRewards.Add(sdk.NewCoin(coin.Denom, rest))
		}
	}
	return
}

// burnFees burns fees from the module account
func (td TaxDecorator) burnFees(ctx sdk.Context, fees sdk.Coins) error {
	// burn fees
	err := td.bankKeeper.BurnCoins(ctx, didtypes.ModuleName, fees)
	if err != nil {
		return err
	}
	return nil
}

func (td *TaxDecorator) handleTaxableTransaction(
	ctx sdk.Context,
	feeTx sdk.FeeTx,
	simulate bool,
	rewards, burn sdk.Coins,
	tx sdk.Tx,
) error {
	var (
		convertedRewards sdk.Coins
		convertedBurn    sdk.Coins
	)

	params, err := td.feemarketKeeper.GetParams(ctx)
	if err != nil {
		return err
	}
	nativeDenom := params.FeeDenom
	onlyNativeDenom := td.isOnlyNativeDenom(feeTx.GetFee(), nativeDenom)

	cheqPrice, found := td.oracleKeeper.GetWMA(ctx, oracletypes.CheqdSymbol, string(oraclekeeper.WmaStrategyBalanced))
	if !found {
		return err
	}
	if onlyNativeDenom {
		err := td.processNativeDenomTax(ctx, feeTx, simulate, rewards, burn, tx, &convertedRewards, &convertedBurn, cheqPrice)
		if err != nil {
			return err
		}
	} else {
		if err := td.handleNonNativeTax(ctx, feeTx, simulate, rewards, burn, tx, cheqPrice, &convertedRewards, &convertedBurn); err != nil {
			return err
		}
	}

	if err := td.distributeRewards(ctx, convertedRewards); err != nil {
		return err
	}
	if err := td.burnFees(ctx, convertedBurn); err != nil {
		return err
	}

	return nil
}

func (td *TaxDecorator) handleNonNativeTax(
	ctx sdk.Context,
	feeTx sdk.FeeTx,
	simulate bool,
	rewards, burn sdk.Coins,
	tx sdk.Tx,
	cheqPrice math.LegacyDec,
	convertedRewards, convertedBurn *sdk.Coins,
) error {
	var err error

	*convertedRewards, err = ConvertToCheq(rewards, cheqPrice)
	if err != nil {
		return fmt.Errorf("failed to convert rewards to ncheq: %w", err)
	}

	*convertedBurn, err = ConvertToCheq(burn, cheqPrice)
	if err != nil {
		return fmt.Errorf("failed to convert burn to ncheq: %w", err)
	}

	userFee := feeTx.GetFee()
	if userFee.IsZero() {
		return fmt.Errorf("user fee is zero")
	}
	denom := userFee[0].Denom

	hostChainConfig, found := td.feeabsKeeper.GetHostZoneConfig(ctx, denom)
	if !found {
		return fmt.Errorf("host chain config not found for denom: %s", denom)
	}

	totalTax := convertedRewards.Add(*convertedBurn...)

	return td.processNonNativeDenomTax(ctx, tx, simulate, feeTx, hostChainConfig, totalTax)
}

func (td *TaxDecorator) isOnlyNativeDenom(fees sdk.Coins, nativeDenom string) bool {
	for _, fee := range fees {
		if fee.Denom != nativeDenom {
			return false
		}
	}
	return true
}

func (td *TaxDecorator) processNativeDenomTax(
	ctx sdk.Context,
	feeTx sdk.FeeTx,
	simulate bool,
	rewards, burn sdk.Coins,
	tx sdk.Tx,
	convertedRewards, convertedBurn *sdk.Coins,
	cheqPrice math.LegacyDec,
) error {
	if err := td.validateTax(feeTx.GetFee(), simulate); err != nil {
		return err
	}
	var err error
	*convertedRewards, err = ConvertToCheq(rewards, cheqPrice)
	if err != nil {
		return err
	}
	*convertedBurn, err = ConvertToCheq(sdk.NewCoins(burn...), cheqPrice)
	if err != nil {
		return err
	}
	tax := convertedRewards.Add(*convertedBurn...)
	feePayer, err := td.getFeePayer(ctx, feeTx, tax, tx.GetMsgs())
	if err != nil {
		return err
	}
	return td.deductTaxFromFeePayer(ctx, feePayer, tax)
}

func (td *TaxDecorator) processNonNativeDenomTax(ctx sdk.Context, tx sdk.Tx, simulate bool, feeTx sdk.FeeTx, hostChainConfig feeabstypes.HostChainFeeAbsConfig, nativeFeeTax sdk.Coins) error {
	if hostChainConfig.Status == feeabstypes.HostChainFeeAbsStatus_FROZEN {
		return errors.Wrap(feeabstypes.ErrHostZoneFrozen, "cannot deduct fee as host zone is frozen")
	}

	// if hostChainConfig.Status == feeabstypes.HostChainFeeAbsStatus_OUTDATED {
	// 	return ctx, sdkerrors.Wrap(feeabstypes.ErrHostZoneOutdated, "cannot deduct fee as host zone is outdated")
	// }
	fee := feeTx.GetFee()
	feePayer := feeTx.FeePayer()
	feeGranter := feeTx.FeeGranter()

	feeAbstractionPayer := feePayer
	// if feegranter set deduct fee from feegranter account.
	// this works with only when feegrant enabled.
	if feeGranter != nil {
		if td.feegrantKeeper == nil {
			return errors.Wrap(sdkerrors.ErrInvalidRequest, "fee grants are not enabled")
		} else if !bytes.Equal(feeGranter, feePayer) {
			err := td.feegrantKeeper.UseGrantedFees(ctx, feeGranter, feePayer, fee, tx.GetMsgs())
			if err != nil {
				return errors.Wrapf(err, "%s not allowed to pay fees from %s", feeGranter, feePayer)
			}
		}

		feeAbstractionPayer = feeGranter
	}

	deductFeesFrom := td.feeabsKeeper.GetFeeAbsModuleAddress()
	deductFeesFromAcc := td.accountKeeper.GetAccount(ctx, deductFeesFrom)
	if deductFeesFromAcc == nil {
		return errors.Wrapf(sdkerrors.ErrUnknownAddress, "fee abstraction didn't set : %s does not exist", deductFeesFrom)
	}

	// calculate the native token can be swapped from ibc token
	if len(fee) != 1 {
		return errors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid ibc token: %s", fee)
	}

	ibcFeetax, err := CalculateIBCCoinsFromNative(ctx, nativeFeeTax, hostChainConfig, td.feeabsKeeper)
	if err != nil {
		return err
	}

	// deduct the fees
	if !feeTx.GetFee().IsZero() {
		err := td.bankKeeper.SendCoinsFromAccountToModule(ctx, sdk.AccAddress(feeAbstractionPayer), feeabstypes.ModuleName, ibcFeetax)
		if err != nil {
			return err
		}
		err = DeductFees(td.bankKeeper, ctx, deductFeesFrom, nativeFeeTax)
		if err != nil {
			return err
		}
	}

	return nil
}

func CalculateIBCCoinsFromNative(ctx sdk.Context, nativeCoins sdk.Coins, chainConfig feeabstypes.HostChainFeeAbsConfig, feeabskeeper feeabskeeper.Keeper) (sdk.Coins, error) {
	// Support only 1 native denom at a time
	if len(nativeCoins) != 1 {
		return nil, errors.Wrapf(sdkerrors.ErrInvalidCoins, "expected single native coin, got: %s", nativeCoins)
	}

	nativeCoin := nativeCoins[0]

	// Get TWAP rate: ibcDenom price in native denom
	twapRate, err := feeabskeeper.GetTwapRate(ctx, chainConfig.IbcDenom)
	if err != nil {
		return nil, err
	}

	if twapRate.IsZero() {
		return nil, errors.Wrapf(sdkerrors.ErrInvalidRequest, "zero twap rate for denom %s", chainConfig.IbcDenom)
	}

	// inverse: ibcAmount = nativeAmount / twapRate
	ibcAmount := math.LegacyNewDecFromInt(nativeCoin.Amount).Quo(twapRate).TruncateInt()

	// Create and return the IBC coin
	ibcCoin := sdk.NewCoin(chainConfig.IbcDenom, ibcAmount)

	ibcCoins := sdk.NewCoins(ibcCoin)
	if ibcCoins.Len() != 1 {
		return nil, feeabstypes.ErrInvalidIBCFees
	}

	ibcDenom := ibcCoins[0].Denom
	if !feeabskeeper.HasHostZoneConfig(ctx, ibcDenom) {
		return nil, feeabstypes.ErrHostZoneConfigNotFound
	}
	return sdk.NewCoins(ibcCoin), nil
}

// DeductFees deducts fees from the given account.
func DeductFees(bankKeeper types.BankKeeper, ctx sdk.Context, accAddress sdk.AccAddress, fees sdk.Coins) error {
	if err := fees.Validate(); err != nil {
		return errors.Wrapf(sdkerrors.ErrInsufficientFee, "invalid fee amount: %s", fees)
	}

	if err := bankKeeper.SendCoinsFromAccountToModule(ctx, accAddress, didtypes.ModuleName, fees); err != nil {
		return errors.Wrapf(sdkerrors.ErrInsufficientFunds, "failed to deduct fees: %s", err)
	}

	return nil
}

func ConvertToCheq(coins sdk.Coins, cheqPrice math.LegacyDec) (sdk.Coins, error) {
	if coins.DenomsSubsetOf(sdk.NewCoins(sdk.NewCoin(oracletypes.CheqdDenom, math.ZeroInt()))) {
		return coins, nil
	}

	converted := sdk.NewCoins()

	for _, coin := range coins {
		switch coin.Denom {
		case oracletypes.UsdDenom:
			if cheqPrice.IsZero() {
				return nil, fmt.Errorf("cannot convert USD to ncheq: CHEQ price unavailable")
			}

			// Convert: USD (1e18) → CHEQ → ncheq (1e9)
			usdAmount := coin.Amount.ToLegacyDec().Quo(math.LegacyNewDecFromInt(util.UsdExponent))
			ncheqAmount := usdAmount.Quo(cheqPrice).MulInt64(util.CheqScale.Int64()).TruncateInt()

			converted = converted.Add(sdk.NewCoin(oracletypes.CheqdDenom, ncheqAmount))

		case oracletypes.CheqdDenom:
			converted = converted.Add(coin)

		default:
			return nil, fmt.Errorf("unexpected denom: %s", coin.Denom)
		}
	}

	return converted, nil
}
