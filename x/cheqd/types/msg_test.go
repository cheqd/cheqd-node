package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const Prefix = "did:cheqd:test:"

func TestNewMsgCreateDidValidation(t *testing.T) {
	cases := []struct {
		valid  bool
		name   string
		msg    *MsgCreateDid
		errMsg string
	}{
		{true, "Valid Create Did Msg", NewMsgCreateDid(&MsgCreateDidPayload{Id: "1"}, []*SignInfo{{VerificationMethodId: "foo", Signature: "bar"}}), ""},
		{false, "Payload is missed", NewMsgCreateDid(nil, nil), "Payload: is required"},
		{false, "Signatures is missed", NewMsgCreateDid(&MsgCreateDidPayload{Id: "1"}, nil), "Signatures: is required"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()

			if tc.valid {
				require.Nil(t, err)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.errMsg, err.Error())
			}
		})
	}
}

func TestNewMsgUpdateDidValidation(t *testing.T) {
	cases := []struct {
		valid  bool
		name   string
		msg    *MsgUpdateDid
		errMsg string
	}{
		{true, "Valid Update Did Msg", NewMsgUpdateDid(&MsgUpdateDidPayload{Id: "1"}, []*SignInfo{{VerificationMethodId: "foo", Signature: "bar"}}), ""},
		{false, "Payload is missed", NewMsgUpdateDid(nil, nil), "Payload: is required"},
		{false, "Signatures is missed", NewMsgUpdateDid(&MsgUpdateDidPayload{Id: "1"}, nil), "Signatures: is required"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()

			if tc.valid {
				require.Nil(t, err)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.errMsg, err.Error())
			}
		})
	}
}

func TestMsgCreateDidPayloadPayload(t *testing.T) {
	cases := []struct {
		valid  bool
		msg    *MsgCreateDidPayload
		errMsg string
	}{
		{
			false,
			&MsgCreateDidPayload{},
			"Id: is not DID",
		},
		{
			false,
			&MsgCreateDidPayload{Id: ""},
			"Id: is not DID",
		},
		{
			false,
			&MsgCreateDidPayload{Id: "did:ch:test:alice"},
			"Id: is not DID",
		},
		{
			false,
			&MsgCreateDidPayload{Id: "did:cheqd:test:alice"},
			"The message must contain either a Controller or a Authentication: bad request",
		},
		{
			false,
			&MsgCreateDidPayload{Id: "did:cheqd:test:alice",
				VerificationMethod: []*VerificationMethod{},
			},
			"The message must contain either a Controller or a Authentication: bad request",
		},
		{
			false,
			&MsgCreateDidPayload{Id: "did:cheqd:test:alice",
				VerificationMethod: []*VerificationMethod{{}},
			},
			"index 0, value : : is not DID fragment: invalid verification method",
		},
		{
			false,
			&MsgCreateDidPayload{Id: "did:cheqd:test:alice", Authentication: []string{"dd"}},
			"Authentication item dd: is not DID fragment",
		},
		{
			false,
			&MsgCreateDidPayload{Id: "did:cheqd:test:alice", Authentication: []string{""}},
			"Authentication item : is not DID fragment",
		},
		{
			false,
			&MsgCreateDidPayload{Id: "did:cheqd:test:alice", Authentication: []string{"did:cheqd:test:alice"}},
			"Authentication item did:cheqd:test:alice: is not DID fragment",
		},
		{
			false,
			&MsgCreateDidPayload{
				Id:             "did:cheqd:test:alice",
				Authentication: []string{"did:cheqd:test:alice#key-1"},
				Controller:     []string{"did:cheqd:test:alice"},
			},
			"did:cheqd:test:alice#key-1: verification method not found",
		},
		{
			false,
			&MsgCreateDidPayload{Id: "did:cheqd:test:alice", CapabilityInvocation: []string{"dd"}},
			"CapabilityInvocation item dd: is not DID fragment",
		},
		{
			false,
			&MsgCreateDidPayload{Id: "did:cheqd:test:alice", CapabilityInvocation: []string{""}},
			"CapabilityInvocation item : is not DID fragment",
		},
		{
			false,
			&MsgCreateDidPayload{Id: "did:cheqd:test:alice", CapabilityInvocation: []string{"did:cheqd:test:alice"}},
			"CapabilityInvocation item did:cheqd:test:alice: is not DID fragment",
		},
		{
			false,
			&MsgCreateDidPayload{
				Id:                   "did:cheqd:test:alice",
				CapabilityInvocation: []string{"did:cheqd:test:alice#key-1"},
				Controller:           []string{"did:cheqd:test:alice"},
			},
			"did:cheqd:test:alice#key-1: verification method not found",
		},
		{
			false,
			&MsgCreateDidPayload{Id: "did:cheqd:test:alice", CapabilityDelegation: []string{"dd"}},
			"CapabilityDelegation item dd: is not DID fragment",
		},
		{
			false,
			&MsgCreateDidPayload{Id: "did:cheqd:test:alice", CapabilityDelegation: []string{""}},
			"CapabilityDelegation item : is not DID fragment",
		},
		{
			false,
			&MsgCreateDidPayload{Id: "did:cheqd:test:alice", CapabilityDelegation: []string{"did:cheqd:test:alice"}},
			"CapabilityDelegation item did:cheqd:test:alice: is not DID fragment",
		},
		{
			false,
			&MsgCreateDidPayload{
				Id:                   "did:cheqd:test:alice",
				CapabilityDelegation: []string{"did:cheqd:test:alice#key-1"},
				Controller:           []string{"did:cheqd:test:alice"},
			},
			"did:cheqd:test:alice#key-1: verification method not found",
		},
		{
			false,
			&MsgCreateDidPayload{Id: "did:cheqd:test:alice", KeyAgreement: []string{"dd"}},
			"KeyAgreement item dd: is not DID fragment",
		},
		{
			false,
			&MsgCreateDidPayload{Id: "did:cheqd:test:alice", KeyAgreement: []string{""}},
			"KeyAgreement item : is not DID fragment",
		},
		{
			false,
			&MsgCreateDidPayload{Id: "did:cheqd:test:alice", KeyAgreement: []string{"did:cheqd:test:alice"}},
			"KeyAgreement item did:cheqd:test:alice: is not DID fragment",
		},
		{
			false,
			&MsgCreateDidPayload{
				Id:           "did:cheqd:test:alice",
				KeyAgreement: []string{"did:cheqd:test:alice#key-1"},
				Controller:   []string{"did:cheqd:test:alice"},
			},
			"did:cheqd:test:alice#key-1: verification method not found",
		},
		{
			false,
			&MsgCreateDidPayload{
				Id:           "did:cheqd:test:alice",
				KeyAgreement: []string{"did:cheqd:test:alice#key-1"},
				Controller:   []string{"did:cheqd::alice"},
			},
			"Controller item did:cheqd::alice at position 0: is not DID",
		},
		{
			false,
			&MsgCreateDidPayload{
				Id: "did:cheqd:test:alice",
				VerificationMethod: []*VerificationMethod{
					{Id: "dasda"},
				},
				Controller: []string{"did:cheqd:test:alice"},
			},
			"index 0, value dasda: dasda: is not DID fragment: invalid verification method",
		},
		{
			false,
			&MsgCreateDidPayload{
				Id: "did:cheqd:test:alice",
				VerificationMethod: []*VerificationMethod{
					{Id: "did:cheqd:test:alice#key-1"},
				},
				Controller: []string{"did:cheqd:test:alice"},
			},
			"index 0, value did:cheqd:test:alice#key-1: : unsupported verification method type: bad request: invalid verification method",
		},
		{
			false,
			&MsgCreateDidPayload{
				Id: "did:cheqd:test:alice",
				VerificationMethod: []*VerificationMethod{
					{Id: "did:cheqd:test:alice#key-1", Type: "YES"},
				},
				Controller: []string{"did:cheqd:test:alice"},
			},
			"index 0, value did:cheqd:test:alice#key-1: YES: unsupported verification method type: bad request: invalid verification method",
		},
		{
			false,
			&MsgCreateDidPayload{
				Id: "did:cheqd:test:alice",
				VerificationMethod: []*VerificationMethod{
					{
						Id:   "did:cheqd:test:alice#key-1",
						Type: "JsonWebKey2020",
					},
				},
				Controller: []string{"did:cheqd:test:alice"},
			},
			"index 0, value did:cheqd:test:alice#key-1: JsonWebKey2020: should contain `PublicKeyJwk` verification material property: bad request: invalid verification method",
		},
		{
			false,
			&MsgCreateDidPayload{
				Id: "did:cheqd:test:alice",
				VerificationMethod: []*VerificationMethod{
					{
						Id:                 "did:cheqd:test:alice#key-1",
						Type:               "JsonWebKey2020",
						PublicKeyMultibase: "tetetet",
					},
				},
				Controller: []string{"did:cheqd:test:alice"},
			},
			"index 0, value did:cheqd:test:alice#key-1: JsonWebKey2020: should contain `PublicKeyJwk` verification material property: bad request: invalid verification method",
		},
		{
			false,
			&MsgCreateDidPayload{
				Id: "did:cheqd:test:alice",
				VerificationMethod: []*VerificationMethod{
					{
						Id:                 "did:cheqd:test:alice#key-1",
						Type:               "Ed25519VerificationKey2020",
						PublicKeyMultibase: "tetetet",
						Controller:         "did:cheqd:test:alice",
					},
					{
						Id:                 "did:cheqd:test:alice#key-1",
						Type:               "Ed25519VerificationKey2020",
						PublicKeyMultibase: "tetetet",
						Controller:         "did:cheqd:test:alice",
					},
				},
				Controller: []string{"did:cheqd:test:alice"},
			},
			"did:cheqd:test:alice#key-1 is duplicated: invalid verification method",
		},
		{
			false,
			&MsgCreateDidPayload{
				Id: "did:cheqd:test:alice",
				VerificationMethod: []*VerificationMethod{
					{
						Id:                 "did:cheqd:test:alice#key-1",
						Type:               "Ed25519VerificationKey2020",
						PublicKeyMultibase: "tetetet",
						Controller:         "did:cheqd:test:alice",
					},
					{
						Id:   "did:cheqd:test:alice#key-2",
						Type: "JsonWebKey2020",
						PublicKeyJwk: []*KeyValuePair{
							{
								Key:   "x",
								Value: "sadad",
							},
						},
						Controller: "did:cheqd:test:alice",
					},
					{
						Id:                 "did:cheqd:test:alice#key-3",
						Type:               "JsonWebKey20212",
						PublicKeyMultibase: "tetetet",
						Controller:         "did:cheqd:test:alice",
					},
				},
				Controller: []string{"did:cheqd:test:alice"},
			},
			"index 2, value did:cheqd:test:alice#key-3: JsonWebKey20212: unsupported verification method type: bad request: invalid verification method",
		},
		{
			false,
			&MsgCreateDidPayload{
				Id:         "did:cheqd:test:alice",
				Controller: []string{"did:cheqd:test:alice"},
				Service: []*Service{
					{},
				},
			},
			"index 0, value : : is not DID fragment: invalid service",
		},
		{
			false,
			&MsgCreateDidPayload{
				Id:         "did:cheqd:test:alice",
				Controller: []string{"did:cheqd:test:alice"},
				Service: []*Service{
					{
						Id: "weqweqw",
					},
				},
			},
			"index 0, value weqweqw: weqweqw: is not DID fragment: invalid service",
		},
		{
			false,
			&MsgCreateDidPayload{
				Id:         "did:cheqd:test:alice",
				Controller: []string{"did:cheqd:test:alice"},
				Service: []*Service{
					{
						Id: "#service-1",
					},
				},
			},
			"index 0, value #service-1: : unsupported service type: bad request: invalid service",
		},
		{
			false,
			&MsgCreateDidPayload{
				Id:         "did:cheqd:test:alice",
				Controller: []string{"did:cheqd:test:alice"},
				Service: []*Service{
					{
						Id:   "#service-1",
						Type: "DIDCommMessaging",
					},
					{
						Id:   "#service-1",
						Type: "DIDCommMessaging",
					},
				},
			},
			"#service-1 is duplicated: invalid service",
		},
		{
			false,
			&MsgCreateDidPayload{
				Id:         "did:cheqd:test:alice",
				Controller: []string{"did:cheqd:test:alice"},
				Service: []*Service{
					{
						Id:   "did:cheqd:test:alice#service-1",
						Type: "DIDCommMessaging",
					},
					{
						Id:   "#service-1",
						Type: "DIDCommMessaging",
					},
				},
			},
			"did:cheqd:test:alice#service-1 is duplicated: invalid service",
		},
		{
			true,
			&MsgCreateDidPayload{
				Id:                 "did:cheqd:test:alice",
				Controller:         []string{"did:cheqd:test:alice"},
				VerificationMethod: []*VerificationMethod{},
			},
			"",
		},
		{
			true,
			&MsgCreateDidPayload{
				Id:         "did:cheqd:test:alice",
				Controller: []string{"did:cheqd:test:alice"},
				Service:    []*Service{},
			},
			"",
		},
		{
			true,
			&MsgCreateDidPayload{
				Id:         "did:cheqd:test:alice",
				Controller: []string{"did:cheqd:test:alice"}},
			"",
		},
		{
			true,
			&MsgCreateDidPayload{
				Id:             "did:cheqd:test:alice",
				Controller:     []string{"did:cheqd:test:alice", "did:cheqd:test:bob"},
				Authentication: []string{"#key-1", "did:cheqd:test:alice#key-2"},
				VerificationMethod: []*VerificationMethod{
					{
						Id:   "did:cheqd:test:alice#key-1",
						Type: "JsonWebKey2020",
						PublicKeyJwk: []*KeyValuePair{
							{
								Key:   "x",
								Value: "sadad",
							},
						},
						Controller: "did:cheqd:test:alice",
					},
					{
						Id:                 "did:cheqd:test:alice#key-2",
						Type:               "Ed25519VerificationKey2020",
						PublicKeyMultibase: "tetetet",
						Controller:         "did:cheqd:test:alice",
					},
				},
			},
			"",
		},
	}

	for _, tc := range cases {
		err := tc.msg.Validate(Prefix)

		if tc.valid {
			require.Nil(t, err)
		} else {
			require.Error(t, err)
			require.Equal(t, tc.errMsg, err.Error())
		}
	}
}

func TestNewMsgUpdateDidPayload(t *testing.T) {
	cases := []struct {
		valid  bool
		msg    *MsgUpdateDidPayload
		errMsg string
	}{
		{
			false,
			&MsgUpdateDidPayload{},
			"Id: is not DID",
		},
		{
			false,
			&MsgUpdateDidPayload{Id: ""},
			"Id: is not DID",
		},
		{
			false,
			&MsgUpdateDidPayload{Id: "did:ch:test:alice"},
			"Id: is not DID",
		},
		{
			false,
			&MsgUpdateDidPayload{Id: "did:cheqd:test:alice"},
			"The message must contain either a Controller or a Authentication: bad request",
		},
		{
			false,
			&MsgUpdateDidPayload{Id: "did:cheqd:test:alice",
				VerificationMethod: []*VerificationMethod{},
			},
			"The message must contain either a Controller or a Authentication: bad request",
		},
		{
			false,
			&MsgUpdateDidPayload{Id: "did:cheqd:test:alice",
				VerificationMethod: []*VerificationMethod{{}},
			},
			"index 0, value : : is not DID fragment: invalid verification method",
		},
		{
			false,
			&MsgUpdateDidPayload{Id: "did:cheqd:test:alice", Authentication: []string{"dd"}},
			"Authentication item dd: is not DID fragment",
		},
		{
			false,
			&MsgUpdateDidPayload{Id: "did:cheqd:test:alice", Authentication: []string{""}},
			"Authentication item : is not DID fragment",
		},
		{
			false,
			&MsgUpdateDidPayload{Id: "did:cheqd:test:alice", Authentication: []string{"did:cheqd:test:alice"}},
			"Authentication item did:cheqd:test:alice: is not DID fragment",
		},
		{
			false,
			&MsgUpdateDidPayload{
				Id:             "did:cheqd:test:alice",
				Authentication: []string{"did:cheqd:test:alice#key-1"},
				Controller:     []string{"did:cheqd:test:alice"},
			},
			"did:cheqd:test:alice#key-1: verification method not found",
		},
		{
			false,
			&MsgUpdateDidPayload{Id: "did:cheqd:test:alice", CapabilityInvocation: []string{"dd"}},
			"CapabilityInvocation item dd: is not DID fragment",
		},
		{
			false,
			&MsgUpdateDidPayload{Id: "did:cheqd:test:alice", CapabilityInvocation: []string{""}},
			"CapabilityInvocation item : is not DID fragment",
		},
		{
			false,
			&MsgUpdateDidPayload{Id: "did:cheqd:test:alice", CapabilityInvocation: []string{"did:cheqd:test:alice"}},
			"CapabilityInvocation item did:cheqd:test:alice: is not DID fragment",
		},
		{
			false,
			&MsgUpdateDidPayload{
				Id:                   "did:cheqd:test:alice",
				CapabilityInvocation: []string{"did:cheqd:test:alice#key-1"},
				Controller:           []string{"did:cheqd:test:alice"},
			},
			"did:cheqd:test:alice#key-1: verification method not found",
		},
		{
			false,
			&MsgUpdateDidPayload{Id: "did:cheqd:test:alice", CapabilityDelegation: []string{"dd"}},
			"CapabilityDelegation item dd: is not DID fragment",
		},
		{
			false,
			&MsgUpdateDidPayload{Id: "did:cheqd:test:alice", CapabilityDelegation: []string{""}},
			"CapabilityDelegation item : is not DID fragment",
		},
		{
			false,
			&MsgUpdateDidPayload{Id: "did:cheqd:test:alice", CapabilityDelegation: []string{"did:cheqd:test:alice"}},
			"CapabilityDelegation item did:cheqd:test:alice: is not DID fragment",
		},
		{
			false,
			&MsgUpdateDidPayload{
				Id:                   "did:cheqd:test:alice",
				CapabilityDelegation: []string{"did:cheqd:test:alice#key-1"},
				Controller:           []string{"did:cheqd:test:alice"},
			},
			"did:cheqd:test:alice#key-1: verification method not found",
		},
		{
			false,
			&MsgUpdateDidPayload{Id: "did:cheqd:test:alice", KeyAgreement: []string{"dd"}},
			"KeyAgreement item dd: is not DID fragment",
		},
		{
			false,
			&MsgUpdateDidPayload{Id: "did:cheqd:test:alice", KeyAgreement: []string{""}},
			"KeyAgreement item : is not DID fragment",
		},
		{
			false,
			&MsgUpdateDidPayload{Id: "did:cheqd:test:alice", KeyAgreement: []string{"did:cheqd:test:alice"}},
			"KeyAgreement item did:cheqd:test:alice: is not DID fragment",
		},
		{
			false,
			&MsgUpdateDidPayload{
				Id:           "did:cheqd:test:alice",
				KeyAgreement: []string{"did:cheqd:test:alice#key-1"},
				Controller:   []string{"did:cheqd:test:alice"},
			},
			"did:cheqd:test:alice#key-1: verification method not found",
		},
		{
			false,
			&MsgUpdateDidPayload{
				Id:           "did:cheqd:test:alice",
				KeyAgreement: []string{"did:cheqd:test:alice#key-1"},
				Controller:   []string{"did:cheqd::alice"},
			},
			"Controller item did:cheqd::alice at position 0: is not DID",
		},
		{
			false,
			&MsgUpdateDidPayload{
				Id: "did:cheqd:test:alice",
				VerificationMethod: []*VerificationMethod{
					{Id: "dasda"},
				},
				Controller: []string{"did:cheqd:test:alice"},
			},
			"index 0, value dasda: dasda: is not DID fragment: invalid verification method",
		},
		{
			false,
			&MsgUpdateDidPayload{
				Id: "did:cheqd:test:alice",
				VerificationMethod: []*VerificationMethod{
					{Id: "did:cheqd:test:alice#key-1"},
				},
				Controller: []string{"did:cheqd:test:alice"},
			},
			"index 0, value did:cheqd:test:alice#key-1: : unsupported verification method type: bad request: invalid verification method",
		},
		{
			false,
			&MsgUpdateDidPayload{
				Id: "did:cheqd:test:alice",
				VerificationMethod: []*VerificationMethod{
					{Id: "did:cheqd:test:alice#key-1", Type: "YES"},
				},
				Controller: []string{"did:cheqd:test:alice"},
			},
			"index 0, value did:cheqd:test:alice#key-1: YES: unsupported verification method type: bad request: invalid verification method",
		},
		{
			false,
			&MsgUpdateDidPayload{
				Id: "did:cheqd:test:alice",
				VerificationMethod: []*VerificationMethod{
					{
						Id:   "did:cheqd:test:alice#key-1",
						Type: "Ed25519VerificationKey2020",
					},
				},
				Controller: []string{"did:cheqd:test:alice"},
			},
			"index 0, value did:cheqd:test:alice#key-1: Ed25519VerificationKey2020: should contain `PublicKeyMultibase` verification material property: bad request: invalid verification method",
		},
		{
			false,
			&MsgUpdateDidPayload{
				Id: "did:cheqd:test:alice",
				VerificationMethod: []*VerificationMethod{
					{
						Id:   "did:cheqd:test:alice#key-1",
						Type: "JsonWebKey2020",
						PublicKeyJwk: []*KeyValuePair{
							{Key: "x", Value: "y"},
						},
					},
				},
				Controller: []string{"did:cheqd:test:alice"},
			},
			"index 0, value did:cheqd:test:alice#key-1: Controller: is required: invalid verification method",
		},
		{
			false,
			&MsgUpdateDidPayload{
				Id: "did:cheqd:test:alice",
				VerificationMethod: []*VerificationMethod{
					{
						Id:                 "did:cheqd:test:alice#key-1",
						Type:               "Ed25519VerificationKey2020",
						PublicKeyMultibase: "tetetet",
						Controller:         "did:cheqd:test:alice",
					},
					{
						Id:                 "did:cheqd:test:alice#key-1",
						Type:               "Ed25519VerificationKey2020",
						PublicKeyMultibase: "tetetet",
						Controller:         "did:cheqd:test:alice",
					},
				},
				Controller: []string{"did:cheqd:test:alice"},
			},
			"did:cheqd:test:alice#key-1 is duplicated: invalid verification method",
		},
		{
			false,
			&MsgUpdateDidPayload{
				Id: "did:cheqd:test:alice",
				VerificationMethod: []*VerificationMethod{
					{
						Id:                 "did:cheqd:test:alice#key-1",
						Type:               "Ed25519VerificationKey2020",
						PublicKeyMultibase: "tetetet",
						Controller:         "did:cheqd:test:alice",
					},
					{
						Id:                 "did:cheqd:test:alice#key-2",
						Type:               "Ed25519VerificationKey2020",
						PublicKeyMultibase: "tetetet",
						Controller:         "did:cheqd:test:alice",
					},
					{
						Id:                 "did:cheqd:test:alice#key-3",
						Type:               "JsonWebKey20212",
						PublicKeyMultibase: "tetetet",
						Controller:         "did:cheqd:test:alice",
					},
				},
				Controller: []string{"did:cheqd:test:alice"},
			},
			"index 2, value did:cheqd:test:alice#key-3: JsonWebKey20212: unsupported verification method type: bad request: invalid verification method",
		},
		{
			false,
			&MsgUpdateDidPayload{
				Id:         "did:cheqd:test:alice",
				Controller: []string{"did:cheqd:test:alice"},
				Service: []*Service{
					{},
				},
			},
			"index 0, value : : is not DID fragment: invalid service",
		},
		{
			false,
			&MsgUpdateDidPayload{
				Id:         "did:cheqd:test:alice",
				Controller: []string{"did:cheqd:test:alice"},
				Service: []*Service{
					{
						Id: "weqweqw",
					},
				},
			},
			"index 0, value weqweqw: weqweqw: is not DID fragment: invalid service",
		},
		{
			false,
			&MsgUpdateDidPayload{
				Id:         "did:cheqd:test:alice",
				Controller: []string{"did:cheqd:test:alice"},
				Service: []*Service{
					{
						Id: "#service-1",
					},
				},
			},
			"index 0, value #service-1: : unsupported service type: bad request: invalid service",
		},
		{
			false,
			&MsgUpdateDidPayload{
				Id:         "did:cheqd:test:alice",
				Controller: []string{"did:cheqd:test:alice"},
				Service: []*Service{
					{
						Id:   "#service-1",
						Type: "DIDCommMessaging",
					},
					{
						Id:   "#service-1",
						Type: "DIDCommMessaging",
					},
				},
			},
			"#service-1 is duplicated: invalid service",
		},
		{
			false,
			&MsgUpdateDidPayload{
				Id:         "did:cheqd:test:alice",
				Controller: []string{"did:cheqd:test:alice"},
				Service: []*Service{
					{
						Id:   "did:cheqd:test:alice#service-1",
						Type: "DIDCommMessaging",
					},
					{
						Id:   "#service-1",
						Type: "DIDCommMessaging",
					},
				},
			},
			"did:cheqd:test:alice#service-1 is duplicated: invalid service",
		},
		{
			true,
			&MsgUpdateDidPayload{
				Id:                 "did:cheqd:test:alice",
				Controller:         []string{"did:cheqd:test:alice"},
				VerificationMethod: []*VerificationMethod{},
			},
			"",
		},
		{
			true,
			&MsgUpdateDidPayload{
				Id:         "did:cheqd:test:alice",
				Controller: []string{"did:cheqd:test:alice"},
				Service:    []*Service{},
			},
			"",
		},
		{
			true,
			&MsgUpdateDidPayload{
				Id:         "did:cheqd:test:alice",
				Controller: []string{"did:cheqd:test:alice"}},
			"",
		},
		{
			true,
			&MsgUpdateDidPayload{
				Id:             "did:cheqd:test:alice",
				Controller:     []string{"did:cheqd:test:alice", "did:cheqd:test:bob"},
				Authentication: []string{"#key-1", "did:cheqd:test:alice#key-2"},
				VerificationMethod: []*VerificationMethod{
					{
						Id:                 "did:cheqd:test:alice#key-1",
						Type:               "Ed25519VerificationKey2020",
						PublicKeyMultibase: "tetetet",
						Controller:         "did:cheqd:test:alice",
					},
					{
						Id:                 "did:cheqd:test:alice#key-2",
						Type:               "Ed25519VerificationKey2020",
						PublicKeyMultibase: "tetetet",
						Controller:         "did:cheqd:test:alice",
					},
				},
			},
			"",
		},
	}

	for _, tc := range cases {
		err := tc.msg.Validate(Prefix)

		if tc.valid {
			require.Nil(t, err)
		} else {
			require.Error(t, err)
			require.Equal(t, tc.errMsg, err.Error())
		}
	}
}
