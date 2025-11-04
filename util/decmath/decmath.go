package decmath

import (
	"fmt"
	"sort"
	"strconv"

	"cosmossdk.io/math"
)

var ErrEmptyList = fmt.Errorf("empty price list passed in")

// Median returns the median of a list of math.LegacyDec. Returns error
// if ds is empty list.
func Median(ds []math.LegacyDec) (math.LegacyDec, error) {
	if len(ds) == 0 {
		return math.LegacyZeroDec(), ErrEmptyList
	}

	sort.Slice(ds, func(i, j int) bool {
		return ds[i].BigInt().
			Cmp(ds[j].BigInt()) < 0
	})

	if len(ds)%2 == 0 {
		return ds[len(ds)/2-1].
			Add(ds[len(ds)/2]).
			QuoInt64(2), nil
	}
	return ds[len(ds)/2], nil
}

// MedianDeviation returns the standard deviation around the
// median of a list of math.LegacyDec. Returns error if ds is empty list.
// MedianDeviation = ∑((d - median)^2 / len(ds))
func MedianDeviation(median math.LegacyDec, ds []math.LegacyDec) (math.LegacyDec, error) {
	if len(ds) == 0 {
		return math.LegacyZeroDec(), ErrEmptyList
	}

	variance := math.LegacyZeroDec()
	for _, d := range ds {
		variance = variance.Add(
			d.Sub(median).Abs().Power(2).QuoInt64(int64(len(ds))))
	}

	medianDeviation, err := variance.ApproxSqrt()
	if err != nil {
		return math.LegacyZeroDec(), err
	}

	return medianDeviation, nil
}

// Average returns the average value of a list of math.LegacyDec. Returns error
// if ds is empty list.
func Average(ds []math.LegacyDec) (math.LegacyDec, error) {
	if len(ds) == 0 {
		return math.LegacyZeroDec(), ErrEmptyList
	}

	sumPrices := math.LegacyZeroDec()
	for _, d := range ds {
		sumPrices = sumPrices.Add(d)
	}

	return sumPrices.QuoInt64(int64(len(ds))), nil
}

// Max returns the max value of a list of math.LegacyDec. Returns error
// if ds is empty list.
func Max(ds []math.LegacyDec) (math.LegacyDec, error) {
	if len(ds) == 0 {
		return math.LegacyZeroDec(), ErrEmptyList
	}

	max := ds[0]
	for _, d := range ds[1:] {
		if d.GT(max) {
			max = d
		}
	}

	return max, nil
}

// Min returns the min value of a list of math.LegacyDec. Returns error
// if ds is empty list.
func Min(ds []math.LegacyDec) (math.LegacyDec, error) {
	if len(ds) == 0 {
		return math.LegacyZeroDec(), ErrEmptyList
	}

	min := ds[0]
	for _, d := range ds[1:] {
		if d.LT(min) {
			min = d
		}
	}

	return min, nil
}

// NewDecFromFloat converts a float64 into a math.LegacyDec. Returns error
// if float64 cannot be converted into a string, or into a subsequent
// math.LegacyDec.
func NewDecFromFloat(f float64) (math.LegacyDec, error) {
	return math.LegacyNewDecFromStr(strconv.FormatFloat(f, 'f', -1, 64))
}
