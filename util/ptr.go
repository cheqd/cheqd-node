package util

import (
	sdkmath "cosmossdk.io/math"
)

const UsdExponent = 1e18

func PtrInt(val int64) *sdkmath.Int {
	i := sdkmath.NewInt(val)
	return &i
}
