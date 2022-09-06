package utils_test

import (
	"testing"

	"github.com/multiformats/go-multibase"
	"github.com/stretchr/testify/require"

	"github.com/cheqd/cheqd-node/x/cheqd/utils"
)

func TestValidateEd25519PubKey(t *testing.T) {
	cases := []struct {
		name     string
		key      string
		valid    bool
		errorMsg string
	}{
		{"Valid: General Ed25519 public key", "zF1hVGXXK9rmx5HhMTpGnGQJiab9qrFJbQXBRhSmYjQWX", true, ""},
		{"Valid: General Ed25519 public key", "zF1hVGXXK9rmx5HhMTpGnGQJiab9qr1111111111111", false, "ed25519: bad public key length: 31"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, keyBytes, _ := multibase.Decode(tc.key)
			err := utils.ValidateEd25519PubKey(keyBytes)

			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errorMsg)
			}
		})
	}
}

func TestValidateJwk(t *testing.T) {
	cases := []struct {
		name     string
		key      string
		valid    bool
		errorMsg string
	}{
		{"positive ed25519", "{\"crv\":\"Ed25519\",\"kty\":\"OKP\",\"x\":\"9Ov80OqMlNrILAUG8DBBlYQ1rUhp7wDomr2I5muzpTc\"}", true, ""},
		{"positive ecdsa", "{\"crv\":\"P-256\",\"kty\":\"EC\",\"x\":\"tcEgxIPyYMiyR2_Vh_YMYG6Grg7axhK2N8JjWta5C0g\",\"y\":\"imiXD9ahVA_MKY066TrNA9r6l35lRrerP6JRey5SryQ\"}", true, ""},
		{"positive rsa", "{\"e\":\"AQAB\",\"kty\":\"RSA\",\"n\":\"skKXRn44WN2DpXDwm4Ip25kIAGRA8y3iXlaoAhPmFiuSDkx97lXcJYrjxX0wSfehgCiSoZOBv6mFzgSVv0_pXQ6zI35xi2dsbexrc87m7Q24q2chpG33ttnVwQkoXrrm0zDzSX32EVxYQyTu9aWp-zxUdAWcrWUarT24RmgjU78v8JmUzkLmwbzsEImnIZ8Hce2ruisAmuAQBVVA4bWwQm_x1KPoQW-TP5_UR3gGugvf0XrQfMJaVpcxcJ9tduMUw6ffZOsqgbvAiZYnrezxSIjnd5lFTFBIEYdGR6ZgjYZoWvQB7U72o_TJoka-zfSODOUbxNBvxvFhA3uhoo3ZKw\"}", true, ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := utils.ValidateJWK(tc.key)

			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errorMsg)
			}
		})
	}
}
