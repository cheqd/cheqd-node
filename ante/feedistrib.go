package ante

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	BurnFeePortion int = iota
	RewardsFeePortion

	FeePortionCount
)

type DistributionFeeAllocation = [FeePortionCount]sdk.Coins

func GetDistributionFee(ctx sdk.Context, fee sdk.Coins, burnFeePortion sdk.Coins) (DistributionFeeAllocation, error) {
	rewardsFeePortion := sdk.NewCoins(fee...).Sub(burnFeePortion...)
	distrFeeAlloc := DistributionFeeAllocation{
		BurnFeePortion:    burnFeePortion,
		RewardsFeePortion: rewardsFeePortion,
	}

	if ValidateDistributionFee(fee, distrFeeAlloc) != nil {
		return distrFeeAlloc, errorsmod.Wrap(sdkerrors.ErrLogic, "fee distribution is invalid")
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
		return errorsmod.Wrap(sdkerrors.ErrInsufficientFee, "fee cannot be zero")
	}

	if !fee.IsEqual(SumDistributionFee(distrFeeAlloc)) {
		return errorsmod.Wrap(sdkerrors.ErrLogic, "fee distribution is invalid")
	}

	return nil
}
