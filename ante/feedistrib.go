package ante

import (
	cheqdtypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

const (
	BurnFeePortion int = iota
	RewardsFeePortion

	FeePortionCount
)

type DistributionFeeAllocation = [FeePortionCount]sdk.Coins

func GetDistributionFee(fee sdk.Coins) (DistributionFeeAllocation, error) {
	distrFeeAlloc := DistributionFeeAllocation{
		BurnFeePortion:       sdk.NewCoins(fee...).QuoInt(sdk.NewInt(int64(BurnFeeDivisor))),
		RewardsFeePortion:    sdk.NewCoins(fee...).Sub(sdk.NewCoins(fee...).QuoInt(sdk.NewInt(int64(BurnFeeDivisor)))...),
	}

	if ValidateDistributionFee(fee, distrFeeAlloc) != nil {
		return distrFeeAlloc, sdkerrors.Wrap(sdkerrors.ErrLogic, "fee distribution is invalid")
	}

	return distrFeeAlloc, nil
}

func SumDistributionFee(distrFeeAlloc DistributionFeeAllocation) sdk.Coins {
	sum := sdk.NewCoins()

	for _, fee := range distrFeeAlloc {
		sum = sum.Add(fee...)
	}

	return sum
}

func ValidateDistributionFee(fee sdk.Coins, distrFeeAlloc DistributionFeeAllocation) error {
	if fee.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInsufficientFee, "fee cannot be zero")
	}

	if !fee.IsEqual(SumDistributionFee(distrFeeAlloc)) {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "fee distribution is invalid")
	}

	return nil
}

//* Redundant, but useful for custom distribution later. Will be removed eventually.
// func DistributeFeeToAccount(bankKeeper BankKeeper, ctx sdk.Context, fee sdk.Coins) error {
// 	if fee.IsZero() {
// 		return sdkerrors.Wrap(sdkerrors.ErrLogic, "fee to be distributed is zero")
// 	}

// 	return bankKeeper.SendCoinsFromModuleToAccount(ctx, cheqdtypes.ModuleName, sdk.AccAddress(FoundationAccAddr), fee)
// }

func DistributeFeeToModule(bankKeeper BankKeeper, ctx sdk.Context, fee sdk.Coins) error {
	if fee.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrLogic, "fee to be distributed is zero")
	}

	return bankKeeper.SendCoinsFromModuleToModule(ctx, cheqdtypes.ModuleName, types.FeeCollectorName, fee)
}
