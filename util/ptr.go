package util

import (
	sdkmath "cosmossdk.io/math"
)

const (
	CheqExponent     = 9
	UsdScaleExponent = 6
)

var (
	CheqScale    = sdkmath.NewIntWithDecimal(1, CheqExponent)
	UsdScale     = sdkmath.NewIntWithDecimal(1, UsdScaleExponent)
	UsdFrom18To6 = sdkmath.NewInt(1_000_000_000_000)
	UsdExponent  = sdkmath.NewIntWithDecimal(1, 18)
)

func PtrInt(val int64) *sdkmath.Int {
	i := sdkmath.NewInt(val)
	return &i
}
