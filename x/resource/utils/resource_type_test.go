package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateResourceType(t *testing.T) {
	cases := []struct {
		rt  string
		valid bool
	}{
		{"CL-Schema", true},
		{"JSONSchema2020", true},
		{"My-Schema", false},
		{"Not schema", false},
	}

	for _, tc := range cases {
		t.Run(tc.rt, func(t *testing.T) {
			err_ := ValidateResourceType(tc.rt)

			if tc.valid {
				require.NoError(t, err_)
			} else {
				require.Error(t, err_)
			}
		})
	}
}
