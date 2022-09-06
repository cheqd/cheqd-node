package utils_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cheqd/cheqd-node/x/cheqd/utils"
)

func TestValidateUUID(t *testing.T) {
	cases := []struct {
		uuid  string
		valid bool
	}{
		{"42d9c704-ecb0-11ec-8ea0-0242ac120002", true},
		{"not uuid", false},
		{"e1cdbc10-858c-4d7d-8a4b-5d45e90a81b3", true},
		{"{42d9c704-ecb0-11ec-8ea0-0242ac120002}", false},
		{"urn:uuid:42d9c704-ecb0-11ec-8ea0-0242ac120002", false},
		{"42d9c704ecb011ec8ea00242ac120002", false},
	}

	for _, tc := range cases {
		t.Run(tc.uuid, func(t *testing.T) {
			err_ := utils.ValidateUUID(tc.uuid)

			if tc.valid {
				require.NoError(t, err_)
			} else {
				require.Error(t, err_)
			}
		})
	}
}
