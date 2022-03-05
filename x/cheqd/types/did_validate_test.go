package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsDid(t *testing.T) {
	cases := []struct {
		valid bool
		did   string
		method string
		allowedNS []string
	}{
		{true, "did:cheqd:testnet:123456789abcdefg", "cheqd", []string{"testnet"}},
		{true, "did:cheqd:testnet:123456789abcdefg", "cheqd", []string{"testnet", "mainnet"}},
		{true, "did:cheqd:testnet:123456789abcdefg", "", []string{"testnet"}},
		{true, "did:cheqd:testnet:123456789abcdefg", "", []string{}},
		{true, "did:cheqd:123456789abcdefg", "", []string{}},
		// Generic method validation
		{true, "did:NOTcheqd:123456789abcdefg", "", []string{}},
		{true, "did:NOTcheqd:123456789abcdefg123456789abcdefg", "", []string{}},

		{true, "did:cheqd:testnet:123456789abcdefg", "cheqd", []string{"mainnet", "testnet"}},
		// Wrong splitting (^did:([^:]+?)(:([^:]+?))?:([^:]+)$)
		{false, "did123:cheqd:::123456789abcdefg", "cheqd", []string{"testnet"}},
		{false, "did:cheqd::123456789abcdefg", "cheqd", []string{"testnet"}},
		{false, "did:cheqd:::123456789abcdefg", "cheqd", []string{"testnet"}},
		{false, "did:cheqd:testnet:123456789abcdefg:did:cheqd:testnet:123456789abcdefg", "cheqd", []string{"testnet"}},
		// Wrong method
		{false, "did:NOTcheqd:testnet:123456789abcdefg", "cheqd", []string{"mainnet", "testnet"}},
		{false, "did:cheqd:testnet:123456789abcdefg", "NOTcheqd", []string{"mainnet", "testnet"}},
		// Wrong namespace (^[a-zA-Z0-9]*)
		{false, "did:cheqd:testnet/:123456789abcdefg", "cheqd", []string{}},
		{false, "did:cheqd:testnet_:123456789abcdefg", "cheqd", []string{}},
		{false, "did:cheqd:testnet%:123456789abcdefg", "cheqd", []string{}},
		{false, "did:cheqd:testnet*:123456789abcdefg", "cheqd", []string{}},
		{false, "did:cheqd:testnet&:123456789abcdefg", "cheqd", []string{}},
		{false, "did:cheqd:testnet@/:123456789abcdefg", "cheqd", []string{}},
		{false, "did:cheqd:testnet:123456789abcdefg", "cheqd", []string{"not_testnet"}},
		// Base58 checks (^[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]*$)
		{false, "did:cheqd:testnet:123456789abcdefO", "cheqd", []string{}},
		{false, "did:cheqd:testnet:123456789abcdefI", "cheqd", []string{}},
		{false, "did:cheqd:testnet:123456789abcdefl", "cheqd", []string{}},
		{false, "did:cheqd:testnet:123456789abcdef0", "cheqd", []string{}},
		// Length checks (should be exactly 16 or 32)
		{false, "did:cheqd:testnet:123", "cheqd", []string{}},
		{false, "did:cheqd:testnet:123456789abcdefgABCDEF", "cheqd", []string{}},
		{false, "did:cheqd:testnet:123456789abcdefg123456789abcdefgABCDEF", "cheqd", []string{}},

	}

	for _, tc := range cases {
		isDid := IsValidDID(tc.did, tc.method, tc.allowedNS)

		if tc.valid {
			require.True(t, isDid)
		} else {
			require.False(t, isDid)
		}
	}
}


func TestIsDidURL(t *testing.T) {
	cases := []struct {
		valid bool
		didUrl   string
	}{
		// Path: all the possible symbols
		{true, "did:cheqd:testnet:123456789abcdefg/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff"},
		{true, "did:cheqd:testnet:123456789abcdefg/path/to/some/other/place/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff/"},
		{true, "did:cheqd:testnet:123456789abcdefg/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?query"},
		{true, "did:cheqd:testnet:123456789abcdefg/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?query#fragment"},
		{true, "did:cheqd:testnet:123456789abcdefg/12/ab/AB/-/./_/~/!/$/&/'/(/)/*/+/,/;/=/:/@/%20/%ff"},
		{true, "did:cheqd:testnet:123456789abcdefg/"},
		// Query: all the possible variants
		{true, "did:cheqd:testnet:123456789abcdefg/path?abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?"},
		{true, "did:cheqd:testnet:123456789abcdefg/path?abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?query=2?query=3/query=%A4"},
		// Fragment:
		{true, "did:cheqd:testnet:123456789abcdefg/path?query#abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff"},
		{true, "did:cheqd:testnet:123456789abcdefg/path?query#abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~!$&'()*+,;=:@%20%ff?"},
		// Wrong cases
		{false, "did:cheqd:testnet:123456789abcdefg/path%20%zz"},
		{false, "did:cheqd:testnet:123456789abcdefg/path?query%20%zz"},
		{false, "did:cheqd:testnet:123456789abcdefg/path?query#fragment%20%zz"},
		// Wrong splitting (^did:([^:]+?)(:([^:]+?))?:([^:]+)$)
		{false, "did123:cheqd:::123456789abcdefg/path?query#fragment"},
		{false, "did:cheqd::123456789abcdefg/path?query#fragment"},
		{false, "did:cheqd:::123456789abcdefg/path?query#fragment"},
		{false, "did:cheqd:testnet:123456789abcdefg:did:cheqd:testnet:123456789abcdefg/path?query#fragment"},
		// Wrong namespace (^[a-zA-Z0-9]*)
		{false, "did:cheqd:testnet/:123456789abcdefg/path?query#fragment"},
		{false, "did:cheqd:testnet_:123456789abcdefg/path?query#fragment"},
		{false, "did:cheqd:testnet%:123456789abcdefg/path?query#fragment"},
		{false, "did:cheqd:testnet*:123456789abcdefg/path?query#fragment"},
		{false, "did:cheqd:testnet&:123456789abcdefg/path?query#fragment"},
		{false, "did:cheqd:testnet@/:123456789abcdefg/path?query#fragment"},
		// Base58 checks (^[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]*$)
		{false, "did:cheqd:testnet:123456789abcdefO/path?query#fragment"},
		{false, "did:cheqd:testnet:123456789abcdefI/path?query#fragment"},
		{false, "did:cheqd:testnet:123456789abcdefl/path?query#fragment"},
		{false, "did:cheqd:testnet:123456789abcdef0/path?query#fragment"},
		// Length checks (should be exactly 16 or 32)
		{false, "did:cheqd:testnet:123/path?query#fragment"},
		{false, "did:cheqd:testnet:123456789abcdefgABCDEF/path?query#fragment"},
		{false, "did:cheqd:testnet:123456789abcdefg123456789abcdefgABCDEF/path?query#fragment"},
	}

	for _, tc := range cases {
		isDid := IsValidDIDUrl(tc.didUrl, "", []string{})

		if tc.valid {
			require.True(t, isDid)
		} else {
			require.False(t, isDid)
		}
	}
}


//func TestDIDURLErrors(t *testing.T) {
//	cases := []struct {
//		didUrl   string
//		error_code int
//		error_string string
//	}{
//		{"did:cheqd:testnet:123456789abcdefg/", ErrStaticDIDURLPathAbemptyNotValid},
//	}
//
//	for _, tc := range cases {
//		err := ValidateDIDUrl(tc.didUrl, "", []string{})
//		require.NotNil(t, err)
//		require.Exactly(t, err.code, tc.error_code)
//
//}

func TestVerificationMethod(t *testing.T) {
	cases := []struct {
		name   string
		struct_ VerificationMethod
		isValid bool
		errorMsg string
	}{
		{
			name: "test case 1",
			struct_: VerificationMethod{
				Id:                 "did1:cheqd:testnet:123456789abcdefg#sdfsdf",
				Type:               "jwk",
				Controller:         "",
				PublicKeyJwk:       nil,
				PublicKeyMultibase: "multibase",
			},
			isValid: true,
			errorMsg: "",
		},
	}

	validator, err := BuildValidator("", nil)
	require.NoError(t, err)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.Struct(tc.struct_)

			if tc.isValid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Equal(t, err.Error(), tc.errorMsg)
			}
		})
	}
}

