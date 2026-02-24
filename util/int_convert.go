package util

import (
	"math"
)

func SafeIntToUint32(i int) uint32 {
	if i < 0 {
		return 0
	}
	if i > math.MaxUint32 {
		return math.MaxUint32
	}
	return uint32(i)
}

func SafeIntToUint64(i int) uint64 {
	if i < 0 {
		return 0
	}

	return uint64(i)
}

func SafeInt64ToUint64(i int64) uint64 {
	if i < 0 {
		return 0
	}

	return uint64(i)
}

func SafeUint64ToUint32(i uint64) uint32 {
	if i > math.MaxUint32 {
		return math.MaxUint32
	}
	return uint32(i)
}

func SafeUint64ToInt64(i uint64) int64 {
	if i > math.MaxInt64 {
		return math.MaxInt64
	}
	return int64(i)
}
