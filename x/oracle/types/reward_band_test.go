package types

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestRewardBandString(t *testing.T) {
	rb := RewardBand{
		SymbolDenom: "cheq",
		RewardBand:  math.LegacyOneDec(),
	}
	require.Equal(t, rb.String(), "symbol_denom: cheq\nreward_band: \"1.000000000000000000\"\n")

	rbl := RewardBandList{rb}
	require.Equal(t, rbl.String(), "symbol_denom: cheq\nreward_band: \"1.000000000000000000\"")
}

func TestRewardBandEqual(t *testing.T) {
	rb := RewardBand{
		SymbolDenom: "cheq",
		RewardBand:  math.LegacyOneDec(),
	}
	rb2 := RewardBand{
		SymbolDenom: "cheq",
		RewardBand:  math.LegacyOneDec(),
	}
	rb3 := RewardBand{
		SymbolDenom: "inequal",
		RewardBand:  math.LegacyOneDec(),
	}

	require.True(t, rb.Equal(&rb2))
	require.False(t, rb.Equal(&rb3))
	require.False(t, rb2.Equal(&rb3))
}

func TestRewardBandDenomFinder(t *testing.T) {
	rbl := RewardBandList{
		{
			SymbolDenom: "foo",
			RewardBand:  math.LegacyOneDec(),
		},
		{
			SymbolDenom: "bar",
			RewardBand:  math.LegacyZeroDec(),
		},
	}

	band, err := rbl.GetBandFromDenom("foo")
	require.NoError(t, err)
	require.Equal(t, band, math.LegacyOneDec())

	band, err = rbl.GetBandFromDenom("bar")
	require.NoError(t, err)
	require.Equal(t, band, math.LegacyZeroDec())

	band, err = rbl.GetBandFromDenom("baz")
	require.Error(t, err)
	require.Equal(t, band, math.LegacyZeroDec())
}
