package types

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

var CorrectTestDID = "did:cheqd:testnet:123456789abcdefg"
var CorrectTestDID2 = "did:cheqd:testnet:gfedcba987654321"
var BadTestDID = "badDid"

func TestMsgCreateDidPayloadValidation(t *testing.T) {
	cases := []struct {
		name     string
		struct_  *MsgCreateDidPayload
		isValid  bool
		errorMsg string
	}{
		{
			name: "Valid: Verification Method: all is fine with type Ed25519VerificationKey2020",
			struct_: &MsgCreateDidPayload{
				Context:    nil,
				Id:         CorrectTestDID,
				Controller: nil,
				VerificationMethod: []*VerificationMethod{
					{
						Id:                 fmt.Sprintf("%s#fragment", CorrectTestDID),
						Type:               "Ed25519VerificationKey2020",
						Controller:         CorrectTestDID,
						PublicKeyJwk:       nil,
						PublicKeyMultibase: "zABCDEFG12345678",
					},
				},
				Authentication:       nil,
				AssertionMethod:      nil,
				CapabilityInvocation: nil,
				CapabilityDelegation: nil,
				KeyAgreement:         nil,
				AlsoKnownAs:          nil,
				Service:              nil,
			},
			isValid:  true,
			errorMsg: "",
		},
		{
			name: "Valid: Verification Method: all is fine with type jwk",
			struct_: &MsgCreateDidPayload{
				Context:    nil,
				Id:         CorrectTestDID,
				Controller: nil,
				VerificationMethod: []*VerificationMethod{
					{
						Id:                  fmt.Sprintf("%s#fragment", CorrectTestDID),
						Type:               "jwk",
						Controller:         CorrectTestDID,
						PublicKeyJwk:       []*KeyValuePair{&KeyValuePair{"key", "value"}},
						PublicKeyMultibase: "",
					},
				},
				Authentication:       nil,
				AssertionMethod:      nil,
				CapabilityInvocation: nil,
				CapabilityDelegation: nil,
				KeyAgreement:         nil,
				AlsoKnownAs:          nil,
				Service:              nil,
			},
			isValid:  true,
			errorMsg: "",
		},
		{
			name: "Not valid: Verification Method: Wrong id",
			struct_: &MsgCreateDidPayload{
				Context:    nil,
				Id:         CorrectTestDID,
				Controller: nil,
				VerificationMethod: []*VerificationMethod{
					{
						Id:                 BadTestDID,
						Type:               "jwk",
						Controller:         CorrectTestDID,
						PublicKeyJwk:       []*KeyValuePair{&KeyValuePair{"key", "value"}},
						PublicKeyMultibase: "",
					},
				},
				Authentication:       nil,
				AssertionMethod:      nil,
				CapabilityInvocation: nil,
				CapabilityDelegation: nil,
				KeyAgreement:         nil,
				AlsoKnownAs:          nil,
				Service:              nil,
			},
			isValid:  false,
			errorMsg: "id: did must match the following regex exactly one time",
		},
		{
			name: "Not valid: Verification Method: Wrong controller",
			struct_: &MsgCreateDidPayload{
				Context:    nil,
				Id:         CorrectTestDID,
				Controller: nil,
				VerificationMethod: []*VerificationMethod{
					{
						Id:                  fmt.Sprintf("%s#fragment", CorrectTestDID),
						Type:               "jwk",
						Controller:         BadTestDID,
						PublicKeyJwk:       []*KeyValuePair{&KeyValuePair{"key", "value"}},
						PublicKeyMultibase: "",
					},
				},
				Authentication:       nil,
				AssertionMethod:      nil,
				CapabilityInvocation: nil,
				CapabilityDelegation: nil,
				KeyAgreement:         nil,
				AlsoKnownAs:          nil,
				Service:              nil,
			},
			isValid:  false,
			errorMsg: "controller: did must match the following regex exactly one time",
		},
		{
			name: "Not valid: Verification Method: Wrong id and controller",
			struct_: &MsgCreateDidPayload{
				Context:    nil,
				Id:         CorrectTestDID,
				Controller: nil,
				VerificationMethod: []*VerificationMethod{
					{
						Id:                  fmt.Sprintf("%s#fragment", CorrectTestDID),
						Type:               "jwk",
						Controller:         BadTestDID,
						PublicKeyJwk:       []*KeyValuePair{&KeyValuePair{"key", "value"}},
						PublicKeyMultibase: "",
					},
				},
				Authentication:       nil,
				AssertionMethod:      nil,
				CapabilityInvocation: nil,
				CapabilityDelegation: nil,
				KeyAgreement:         nil,
				AlsoKnownAs:          nil,
				Service:              nil,
			},
			isValid:  false,
			errorMsg: "controller: did must match the following regex exactly one time",
		},
		{
			name: "Not valid: Verification Method: type - jwk but PublicKeyJwk is nil",
			struct_: &MsgCreateDidPayload{
				Context:    nil,
				Id:         CorrectTestDID,
				Controller: nil,
				VerificationMethod: []*VerificationMethod{
					{
						Id:                  fmt.Sprintf("%s#fragment", CorrectTestDID),
						Type:               "jwk",
						Controller:         CorrectTestDID,
						PublicKeyMultibase: "",
					},
				},
				Authentication:       nil,
				AssertionMethod:      nil,
				CapabilityInvocation: nil,
				CapabilityDelegation: nil,
				KeyAgreement:         nil,
				AlsoKnownAs:          nil,
				Service:              nil,
			},
			isValid:  false,
			errorMsg: "public_key_jwk: must be set when type is jwk",
		},
		{
			name: "Not valid: Verification Method: type - Ed25519VerificationKey2020 but PublicKeyMultibase is empty",
			struct_: &MsgCreateDidPayload{
				Context:    nil,
				Id:         CorrectTestDID,
				Controller: nil,
				VerificationMethod: []*VerificationMethod{
					{
						Id:                  fmt.Sprintf("%s#fragment", CorrectTestDID),
						Type:               "Ed25519VerificationKey2020",
						Controller:         CorrectTestDID,
						PublicKeyMultibase: "",
					},
				},
				Authentication:       nil,
				AssertionMethod:      nil,
				CapabilityInvocation: nil,
				CapabilityDelegation: nil,
				KeyAgreement:         nil,
				AlsoKnownAs:          nil,
				Service:              nil,
			},
			isValid:  false,
			errorMsg: "public_key_multibase: PublicKeyMultibase cannot be empty string",
		},
		{
			name: "Valid: Controller: List of DIDs allowed",
			struct_: &MsgCreateDidPayload{
				Context:    nil,
				Id:         CorrectTestDID,
				Controller: []string{CorrectTestDID, CorrectTestDID2},
				VerificationMethod: []*VerificationMethod{
					{
						Id:                  fmt.Sprintf("%s#fragment", CorrectTestDID),
						Type:               "Ed25519VerificationKey2020",
						Controller:         CorrectTestDID,
						PublicKeyMultibase: "zABCDEFG12345678",
					},
				},
				Authentication:       nil,
				AssertionMethod:      nil,
				CapabilityInvocation: nil,
				CapabilityDelegation: nil,
				KeyAgreement:         nil,
				AlsoKnownAs:          nil,
				Service:              nil,
			},
			isValid:  true,
			errorMsg: "",
		},
		{
			name: "Not valid: Controller: List of DIDs is not allowed",
			struct_: &MsgCreateDidPayload{
				Context:    nil,
				Id:         CorrectTestDID,
				Controller: []string{CorrectTestDID, BadTestDID},
				VerificationMethod: []*VerificationMethod{
					{
						Id:                  fmt.Sprintf("%s#fragment", CorrectTestDID),
						Type:               "Ed25519VerificationKey2020",
						Controller:         CorrectTestDID,
						PublicKeyMultibase: "zABCDEFG12345678",
					},
				},
				Authentication:       nil,
				AssertionMethod:      nil,
				CapabilityInvocation: nil,
				CapabilityDelegation: nil,
				KeyAgreement:         nil,
				AlsoKnownAs:          nil,
				Service:              nil,
			},
			isValid:  false,
			errorMsg: "Errors after the validation process of DID's list: did must match the following regex exactly one time",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.struct_.Validate()

			if tc.isValid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.True(t, strings.Contains(err.Error(), tc.errorMsg))
			}
		})
	}
}
