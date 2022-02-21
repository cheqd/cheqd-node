package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsDid(t *testing.T) {
	cases := []struct {
		valid bool
		did   string
	}{
		{true,  "did:cheqd:test:000000wyywywywyw"},
		{false,  "did:cheqd:test:wyywywywyw:sdadasda"},
		{false, "did1:cheqd:test:wyywywywyw:sdadasda"},
		{false, "did:cheqd2:test:wyywywywyw:sdadasda"},
		{false, "did:cheqd:test4:wyywywywyw:sdadasda"},
		{false, ""},
		{false, "did:cheqd"},
		{false, "did:cheqd:test"},
		{false, "did:cheqd:test:dsdasdad#weqweqwew"},
		{false, "did:cheqd:test:sdasdasdasd/qeweqweqwee"},
		{false, "did:cheqd:test:sdasdasdasd?=qeweqweqwee"},
		{false, "did:cheqd:test:sdasdasdasd&qeweqweqwee"},
	}

	for _, tc := range cases {
		isDid := IsValidDid("did:cheqd:test:", tc.did)

		if tc.valid {
			require.True(t, isDid)
		} else {
			require.False(t, isDid)
		}
	}
}
