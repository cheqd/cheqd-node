package util

import (
	sdkmath "cosmossdk.io/math"
)

func PtrInt(val int64) *sdkmath.Int {
	i := sdkmath.NewInt(val)
	return &i
}
