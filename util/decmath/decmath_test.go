package decmath

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestMedian(t *testing.T) {
	require := require.New(t)
	prices := []math.LegacyDec{
		math.LegacyMustNewDecFromStr("1.12"),
		math.LegacyMustNewDecFromStr("1.07"),
		math.LegacyMustNewDecFromStr("1.11"),
		math.LegacyMustNewDecFromStr("1.2"),
	}

	median, err := Median(prices)
	require.NoError(err)
	require.Equal(math.LegacyMustNewDecFromStr("1.115"), median)

	// test empty prices list
	_, err = Median([]math.LegacyDec{})
	require.ErrorIs(err, ErrEmptyList)
}

func TestMedianDeviation(t *testing.T) {
	require := require.New(t)
	prices := []math.LegacyDec{
		math.LegacyMustNewDecFromStr("1.12"),
		math.LegacyMustNewDecFromStr("1.07"),
		math.LegacyMustNewDecFromStr("1.11"),
		math.LegacyMustNewDecFromStr("1.2"),
	}
	median := math.LegacyMustNewDecFromStr("1.115")

	medianDeviation, err := MedianDeviation(median, prices)
	require.NoError(err)
	require.Equal(math.LegacyMustNewDecFromStr("0.048218253804964775"), medianDeviation)

	// test empty prices list
	_, err = MedianDeviation(median, []math.LegacyDec{})
	require.ErrorIs(err, ErrEmptyList)
}

func TestAverage(t *testing.T) {
	require := require.New(t)
	prices := []math.LegacyDec{
		math.LegacyMustNewDecFromStr("1.12"),
		math.LegacyMustNewDecFromStr("1.07"),
		math.LegacyMustNewDecFromStr("1.11"),
		math.LegacyMustNewDecFromStr("1.2"),
	}

	average, err := Average(prices)
	require.NoError(err)
	require.Equal(math.LegacyMustNewDecFromStr("1.125"), average)

	// test empty prices list
	_, err = Average([]math.LegacyDec{})
	require.ErrorIs(err, ErrEmptyList)
}

func TestMin(t *testing.T) {
	require := require.New(t)
	prices := []math.LegacyDec{
		math.LegacyMustNewDecFromStr("1.12"),
		math.LegacyMustNewDecFromStr("1.07"),
		math.LegacyMustNewDecFromStr("1.11"),
		math.LegacyMustNewDecFromStr("1.2"),
	}

	min, err := Min(prices)
	require.NoError(err)
	require.Equal(math.LegacyMustNewDecFromStr("1.07"), min)

	// test empty prices list
	_, err = Min([]math.LegacyDec{})
	require.ErrorIs(err, ErrEmptyList)
}

func TestMax(t *testing.T) {
	require := require.New(t)
	prices := []math.LegacyDec{
		math.LegacyMustNewDecFromStr("1.12"),
		math.LegacyMustNewDecFromStr("1.07"),
		math.LegacyMustNewDecFromStr("1.11"),
		math.LegacyMustNewDecFromStr("1.2"),
	}

	max, err := Max(prices)
	require.NoError(err)
	require.Equal(math.LegacyMustNewDecFromStr("1.2"), max)

	// test empty prices list
	_, err = Max([]math.LegacyDec{})
	require.ErrorIs(err, ErrEmptyList)
}

func TestNewDecFromFloat(t *testing.T) {
	testCases := []struct {
		name       string
		float      float64
		dec        math.LegacyDec
		expectPass bool
	}{
		{
			name:       "max float64 precision",
			float:      1.000_000_000_000_001,
			dec:        math.LegacyMustNewDecFromStr("1.000000000000001"),
			expectPass: true,
		},
		{
			name:       "over max float64 precision",
			float:      1.000_000_000_000_000_1,
			dec:        math.LegacyMustNewDecFromStr("1"),
			expectPass: true,
		},
		{
			name:       "simple float",
			float:      2999999.9,
			dec:        math.LegacyMustNewDecFromStr("2999999.9"),
			expectPass: true,
		},
		{
			name:       "negative float",
			float:      -10.598,
			dec:        math.LegacyMustNewDecFromStr("-10.598"),
			expectPass: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dec, err := NewDecFromFloat(tc.float)
			if tc.expectPass {
				require.NoError(t, err)
				require.Equal(t, tc.dec, dec)
			} else {
				require.Error(t, err)
			}
		})
	}
}
