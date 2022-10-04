package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsId(t *testing.T) {
	cases := []struct {
		valid bool
		id    string
	}{
		{true, "123456789abcdefg"},
		{true, "123456789abcdefg123456789abcdefg"},
		{true, "3b9b8eec-5b5d-4382-86d8-9185126ff130"},
		{false, "sdf"},
		{false, "sdf:sdf"},
		{false, "12345"},
	}

	for _, tc := range cases {
		t.Run(tc.id, func(t *testing.T) {
			isDid := IsValidID(tc.id)

			if tc.valid {
				require.True(t, isDid)
			} else {
				require.False(t, isDid)
			}
		})
	}
}

func UUIDTestCases() []struct {
	name        string
	did         string
	expectedDid string
} {
	return []struct {
		name        string
		did         string
		expectedDid string
	}{
		{
			name:        "base58 identifier - not changed",
			did:         "did:cheqd:aaaaaaaaaaaaaaaa",
			expectedDid: "did:cheqd:aaaaaaaaaaaaaaaa",
		},
		{
			name:        "Mixed case UUID",
			did:         "did:cheqd:test:BAbbba14-f294-458a-9b9c-474d188680fd",
			expectedDid: "did:cheqd:test:babbba14-f294-458a-9b9c-474d188680fd",
		},
		{
			name:        "Low case UUID",
			did:         "did:cheqd:test:babbba14-f294-458a-9b9c-474d188680fd",
			expectedDid: "did:cheqd:test:babbba14-f294-458a-9b9c-474d188680fd",
		},
		{
			name:        "Upper case UUID",
			did:         "did:cheqd:test:A86F9CAE-0902-4a7c-a144-96b60ced2FC9",
			expectedDid: "did:cheqd:test:a86f9cae-0902-4a7c-a144-96b60ced2fc9",
		},
	}
}

func TestUpdateUUIDForDID(t *testing.T) {
	for _, tc := range UUIDTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			result := NormalizeIdentifier(tc.did)

			require.Equal(t, tc.expectedDid, result)
		})
	}
}
