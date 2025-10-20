package ante

import (
	"math"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"github.com/noble-assets/globalfee/keeper"
)

// TxFeeChecker returns a custom ante.TxFeeChecker that ensures the fees for a
// given transaction respect the gas prices set in the GlobalFee module.
//
//nolint:nestif
func TxFeeChecker(keeper *keeper.Keeper) ante.TxFeeChecker {
	return func(ctx sdk.Context, tx sdk.Tx) (sdk.Coins, int64, error) {
		feeTx, ok := tx.(sdk.FeeTx)
		if !ok {
			return nil, 0, errors.Wrap(errorstypes.ErrTxDecode, "Tx must be a FeeTx")
		}
		fees := feeTx.GetFee()

		if ctx.IsCheckTx() {
			allBypassMessages := true
			for _, msg := range feeTx.GetMsgs() {
				if has, _ := keeper.BypassMessages.Has(ctx, sdk.MsgTypeURL(msg)); !has {
					allBypassMessages = false
					break
				}
			}
			if allBypassMessages && fees.IsZero() {
				return sdk.Coins{}, 0, nil
			}

			requiredFees, err := keeper.GetRequiredFees(ctx, feeTx)
			if err != nil {
				return nil, 0, err
			}
			if len(requiredFees) == 0 {
				return sdk.Coins{}, 0, nil
			}

			sufficientFees := false
			for _, fee := range fees {
				found, requiredFee := requiredFees.Find(fee.Denom)
				if found && fee.Amount.GTE(requiredFee.Amount) {
					sufficientFees = true
					break
				}
			}

			if !sufficientFees {
				return nil, 0, errors.Wrapf(errorstypes.ErrInsufficientFee, "expected at least one of %s", requiredFees)
			}
		}

		return fees, getTxPriority(fees, int64(feeTx.GetGas())), nil
	}
}

// getTxPriority is copied from the Cosmos SDK as it's not exported for reuse.
// https://github.com/cosmos/cosmos-sdk/blob/v0.50.10/x/auth/ante/validator_tx_fee.go#L50-L68
func getTxPriority(fee sdk.Coins, gas int64) int64 {
	var priority int64
	for _, c := range fee {
		p := int64(math.MaxInt64)
		gasPrice := c.Amount.QuoRaw(gas)
		if gasPrice.IsInt64() {
			p = gasPrice.Int64()
		}
		if priority == 0 || p < priority {
			priority = p
		}
	}

	return priority
}

type feeMarketBypassDecorator struct {
	globalFeeKeeper *keeper.Keeper
	feeMarket       sdk.AnteDecorator
	fallback        sdk.AnteDecorator
}

// NewFeeMarketBypassDecorator wraps the fee market decorator so transactions
// that only contain bypass messages (either configured in the globalfee module
// or hard-coded such as MsgAcknowledgement) skip the fee market checks and rely
// on the globalfee fallback logic instead.
func NewFeeMarketBypassDecorator(
	globalFeeKeeper *keeper.Keeper,
	feeMarket sdk.AnteDecorator,
	fallback sdk.AnteDecorator,
) sdk.AnteDecorator {
	return feeMarketBypassDecorator{
		globalFeeKeeper: globalFeeKeeper,
		feeMarket:       feeMarket,
		fallback:        fallback,
	}
}

func (d feeMarketBypassDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	if ShouldBypassFeeMarket(ctx, d.globalFeeKeeper, tx) {
		return d.fallback.AnteHandle(ctx, tx, simulate, next)
	}

	return d.feeMarket.AnteHandle(ctx, tx, simulate, next)
}

func ShouldBypassFeeMarket(ctx sdk.Context, k *keeper.Keeper, tx sdk.Tx) bool {
	if tx == nil {
		return false
	}

	msgs := tx.GetMsgs()
	if len(msgs) == 0 {
		return false
	}

	for _, msg := range msgs {
		if IsDefaultFeeMarketBypassMsg(msg) {
			continue
		}

		if k == nil {
			return false
		}

		has, err := k.BypassMessages.Has(ctx, sdk.MsgTypeURL(msg))
		if err != nil || !has {
			return false
		}
	}

	return true
}

func IsDefaultFeeMarketBypassMsg(msg sdk.Msg) bool {
	switch msg.(type) {
	case *channeltypes.MsgAcknowledgement:
		return true
	default:
		return false
	}
}
