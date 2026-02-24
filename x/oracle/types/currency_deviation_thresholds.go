package types

import (
	"strings"

	"gopkg.in/yaml.v3"
)

// String implements fmt.Stringer interface
func (cdt CurrencyDeviationThreshold) String() string {
	out, _ := yaml.Marshal(cdt)
	return string(out)
}

func (cdt CurrencyDeviationThreshold) Equal(cdt2 *CurrencyDeviationThreshold) bool {
	return strings.EqualFold(cdt.BaseDenom, cdt2.BaseDenom) && strings.EqualFold(cdt.Threshold, cdt2.Threshold)
}

// CurrencyDeviationThresholdList is array of CurrencyDeviationThresholds
type CurrencyDeviationThresholdList []CurrencyDeviationThreshold

func (cdtl CurrencyDeviationThresholdList) String() (out string) {
	for _, v := range cdtl {
		out += v.String() + "\n"
	}

	return strings.TrimSpace(out)
}

func (cdtl CurrencyDeviationThresholdList) RemovePair(curr string) CurrencyDeviationThresholdList {
	for i := 0; i < len(cdtl); i++ {
		if cdtl[i].BaseDenom == curr {
			cdtl = append(cdtl[:i], cdtl[i+1:]...)
			i-- // decrement i so the next iteration will check the next element
		}
	}

	return cdtl
}
