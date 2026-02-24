package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCurrencyDeviationThresholdsString(t *testing.T) {
	cdt := CurrencyDeviationThreshold{
		BaseDenom: "OJO",
		Threshold: "1.5",
	}
	require.Equal(
		t,
		cdt.String(),
		"base_denom: OJO\nthreshold: \"1.5\"\n",
	)

	cdtl := CurrencyDeviationThresholdList{cdt}
	require.Equal(
		t,
		cdtl.String(),
		"base_denom: OJO\nthreshold: \"1.5\"",
	)
}

func TestCurrencyDeviationThresholdsEqual(t *testing.T) {
	cdt1 := CurrencyDeviationThreshold{
		BaseDenom: "OJO",
		Threshold: "1.5",
	}
	cdt2 := CurrencyDeviationThreshold{
		BaseDenom: "OJO",
		Threshold: "1.5",
	}
	cdt3 := CurrencyDeviationThreshold{
		BaseDenom: "OJO",
		Threshold: "1.6",
	}
	cdt4 := CurrencyDeviationThreshold{
		BaseDenom: "UMEE",
		Threshold: "1.5",
	}

	require.True(t, cdt1.Equal(&cdt2))
	require.False(t, cdt1.Equal(&cdt3))
	require.False(t, cdt1.Equal(&cdt4))
}
