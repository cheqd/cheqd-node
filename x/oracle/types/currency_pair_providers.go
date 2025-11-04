package types

import (
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

// String implements fmt.Stringer interface
func (cpp CurrencyPairProviders) String() string {
	out, _ := yaml.Marshal(cpp)
	return string(out)
}

func (cpp CurrencyPairProviders) Equal(cpp2 *CurrencyPairProviders) bool {
	if !strings.EqualFold(cpp.BaseDenom, cpp2.BaseDenom) || !strings.EqualFold(cpp.QuoteDenom, cpp2.QuoteDenom) {
		return false
	}

	if !reflect.DeepEqual(cpp.PairAddress, cpp2.PairAddress) {
		return false
	}

	return reflect.DeepEqual(cpp.Providers, cpp2.Providers)
}

// CurrencyPairProvidersList is array of CurrencyPairProviders
type CurrencyPairProvidersList []CurrencyPairProviders

func (cppl CurrencyPairProvidersList) String() (out string) {
	for _, v := range cppl {
		out += v.String() + "\n"
	}

	return strings.TrimSpace(out)
}

func (cppl CurrencyPairProvidersList) RemovePair(pair CurrencyPairProviders) CurrencyPairProvidersList {
	for i := 0; i < len(cppl); i++ {
		if cppl[i].BaseDenom == pair.BaseDenom && cppl[i].QuoteDenom == pair.QuoteDenom {
			cppl = append(cppl[:i], cppl[i+1:]...)
			i-- // decrement i so the next iteration will check the next element
		}
	}

	return cppl
}
