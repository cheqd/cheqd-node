package utils_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cheqd/cheqd-node/x/cheqd/utils"
)

func TestValidateURI(t *testing.T) {
	cases := []struct {
		name  string
		valid bool
		URI   string
	}{
		// Path: all the possible symbols
		{"Valid: General http URI path", true, "http://a.com/a/b/c/d/?query=123#fragment=another_part"},
		{"Valid: General https URI path", true, "https://a.com/a/b/c/d/?query=123#fragment=another_part"},
		{"Valid: only alphabet symbols", true, "SomeAnotherPath"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err_ := utils.ValidateURI(tc.URI)

			if tc.valid {
				require.NoError(t, err_)
			} else {
				require.Error(t, err_)
			}
		})
	}
}
