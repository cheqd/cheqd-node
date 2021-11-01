package types

import (
	"testing"

	ctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/stretchr/testify/require"
)

const Prefix = "did:cheqd:test:"

func TestNewMsgWriteRequestValidation(t *testing.T) {
	cases := []struct {
		valid  bool
		msg    *MsgWriteRequest
		errMsg string
	}{
		{true, NewMsgWriteRequest(&ctypes.Any{TypeUrl: "1", Value: []byte{1}}, []*SignInfo{{VerificationMethodId: "foo", Signature: "bar"}}), ""},
		{false, NewMsgWriteRequest(&ctypes.Any{TypeUrl: "1"}, []*SignInfo{{VerificationMethodId: "foo", Signature: "bar"}}), "Invalid Data: it cannot be empty: bad request"},
		{false, NewMsgWriteRequest(&ctypes.Any{Value: []byte{1}}, []*SignInfo{{VerificationMethodId: "foo", Signature: "bar"}}), "Invalid Data: it cannot be empty: bad request"},
		{false, NewMsgWriteRequest(nil, nil), "Data: is required"},
		{false, NewMsgWriteRequest(&ctypes.Any{TypeUrl: "1", Value: []byte{1}}, nil), "Signatures: is required"},
	}

	for _, tc := range cases {
		err := tc.msg.ValidateBasic()

		if tc.valid {
			require.Nil(t, err)
		} else {
			require.Error(t, err)
			require.Equal(t, tc.errMsg, err.Error())
		}
	}
}

func TestNewMsgCreateDid(t *testing.T) {
	cases := []struct {
		valid  bool
		msg    *MsgCreateDid
		errMsg string
	}{
		{
			false,
			&MsgCreateDid{},
			"Id: is not DID",
		},
		{
			false,
			&MsgCreateDid{Id: ""},
			"Id: is not DID",
		},
		{
			false,
			&MsgCreateDid{Id: "did:ch:test:alice"},
			"Id: is not DID",
		},
		{
			false,
			&MsgCreateDid{Id: "did:cheqd:test:alice"},
			"The message must contain either a Controller or a Authentication: bad request",
		},
		{
			false,
			&MsgCreateDid{Id: "did:cheqd:test:alice",
				VerificationMethod: []*VerificationMethod{},
			},
			"The message must contain either a Controller or a Authentication: bad request",
		},
		{
			false,
			&MsgCreateDid{Id: "did:cheqd:test:alice",
				VerificationMethod: []*VerificationMethod{{}},
			},
			"index 0, value : : is not DID fragment: invalid verification method",
		},
		{
			false,
			&MsgCreateDid{Id: "did:cheqd:test:alice", Authentication: []string{"dd"}},
			"Authentication item dd: is not DID fragment",
		},
		{
			false,
			&MsgCreateDid{Id: "did:cheqd:test:alice", Authentication: []string{""}},
			"Authentication item : is not DID fragment",
		},
		{
			false,
			&MsgCreateDid{Id: "did:cheqd:test:alice", Authentication: []string{"did:cheqd:test:alice"}},
			"Authentication item did:cheqd:test:alice: is not DID fragment",
		},
		{
			false,
			&MsgCreateDid{
				Id:             "did:cheqd:test:alice",
				Authentication: []string{"did:cheqd:test:alice#key-1"},
				Controller:     []string{"did:cheqd:test:alice"},
			},
			"did:cheqd:test:alice#key-1: verification method not found",
		},
		{
			false,
			&MsgCreateDid{Id: "did:cheqd:test:alice", CapabilityInvocation: []string{"dd"}},
			"CapabilityInvocation item dd: is not DID fragment",
		},
		{
			false,
			&MsgCreateDid{Id: "did:cheqd:test:alice", CapabilityInvocation: []string{""}},
			"CapabilityInvocation item : is not DID fragment",
		},
		{
			false,
			&MsgCreateDid{Id: "did:cheqd:test:alice", CapabilityInvocation: []string{"did:cheqd:test:alice"}},
			"CapabilityInvocation item did:cheqd:test:alice: is not DID fragment",
		},
		{
			false,
			&MsgCreateDid{
				Id:                   "did:cheqd:test:alice",
				CapabilityInvocation: []string{"did:cheqd:test:alice#key-1"},
				Controller:           []string{"did:cheqd:test:alice"},
			},
			"did:cheqd:test:alice#key-1: verification method not found",
		},
		{
			false,
			&MsgCreateDid{Id: "did:cheqd:test:alice", CapabilityDelegation: []string{"dd"}},
			"CapabilityDelegation item dd: is not DID fragment",
		},
		{
			false,
			&MsgCreateDid{Id: "did:cheqd:test:alice", CapabilityDelegation: []string{""}},
			"CapabilityDelegation item : is not DID fragment",
		},
		{
			false,
			&MsgCreateDid{Id: "did:cheqd:test:alice", CapabilityDelegation: []string{"did:cheqd:test:alice"}},
			"CapabilityDelegation item did:cheqd:test:alice: is not DID fragment",
		},
		{
			false,
			&MsgCreateDid{
				Id:                   "did:cheqd:test:alice",
				CapabilityDelegation: []string{"did:cheqd:test:alice#key-1"},
				Controller:           []string{"did:cheqd:test:alice"},
			},
			"did:cheqd:test:alice#key-1: verification method not found",
		},
		{
			false,
			&MsgCreateDid{Id: "did:cheqd:test:alice", KeyAgreement: []string{"dd"}},
			"KeyAgreement item dd: is not DID fragment",
		},
		{
			false,
			&MsgCreateDid{Id: "did:cheqd:test:alice", KeyAgreement: []string{""}},
			"KeyAgreement item : is not DID fragment",
		},
		{
			false,
			&MsgCreateDid{Id: "did:cheqd:test:alice", KeyAgreement: []string{"did:cheqd:test:alice"}},
			"KeyAgreement item did:cheqd:test:alice: is not DID fragment",
		},
		{
			false,
			&MsgCreateDid{
				Id:           "did:cheqd:test:alice",
				KeyAgreement: []string{"did:cheqd:test:alice#key-1"},
				Controller:   []string{"did:cheqd:test:alice"},
			},
			"did:cheqd:test:alice#key-1: verification method not found",
		},
		{
			false,
			&MsgCreateDid{
				Id:           "did:cheqd:test:alice",
				KeyAgreement: []string{"did:cheqd:test:alice#key-1"},
				Controller:   []string{"did:cheqd::alice"},
			},
			"Controller item did:cheqd::alice at position 0: is not DID",
		},
		{
			false,
			&MsgCreateDid{
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
			&MsgCreateDid{
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
			&MsgCreateDid{
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
			&MsgCreateDid{
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
			&MsgCreateDid{
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
			&MsgCreateDid{
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
			&MsgCreateDid{
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
			&MsgCreateDid{
				Id:         "did:cheqd:test:alice",
				Controller: []string{"did:cheqd:test:alice"},
				Service: []*ServiceEndpoint{
					{},
				},
			},
			"index 0, value : : is not DID fragment: invalid service",
		},
		{
			false,
			&MsgCreateDid{
				Id:         "did:cheqd:test:alice",
				Controller: []string{"did:cheqd:test:alice"},
				Service: []*ServiceEndpoint{
					{
						Id: "weqweqw",
					},
				},
			},
			"index 0, value weqweqw: weqweqw: is not DID fragment: invalid service",
		},
		{
			false,
			&MsgCreateDid{
				Id:         "did:cheqd:test:alice",
				Controller: []string{"did:cheqd:test:alice"},
				Service: []*ServiceEndpoint{
					{
						Id: "#service-1",
					},
				},
			},
			"index 0, value #service-1: : unsupported service type: bad request: invalid service",
		},
		{
			false,
			&MsgCreateDid{
				Id:         "did:cheqd:test:alice",
				Controller: []string{"did:cheqd:test:alice"},
				Service: []*ServiceEndpoint{
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
			&MsgCreateDid{
				Id:         "did:cheqd:test:alice",
				Controller: []string{"did:cheqd:test:alice"},
				Service: []*ServiceEndpoint{
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
			&MsgCreateDid{
				Id:                 "did:cheqd:test:alice",
				Controller:         []string{"did:cheqd:test:alice"},
				VerificationMethod: []*VerificationMethod{},
			},
			"",
		},
		{
			true,
			&MsgCreateDid{
				Id:         "did:cheqd:test:alice",
				Controller: []string{"did:cheqd:test:alice"},
				Service:    []*ServiceEndpoint{},
			},
			"",
		},
		{
			true,
			&MsgCreateDid{
				Id:         "did:cheqd:test:alice",
				Controller: []string{"did:cheqd:test:alice"}},
			"",
		},
		{
			true,
			&MsgCreateDid{
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

func TestNewMsgUpdateDid(t *testing.T) {
	cases := []struct {
		valid  bool
		msg    *MsgUpdateDid
		errMsg string
	}{
		{
			false,
			&MsgUpdateDid{},
			"Id: is not DID",
		},
		{
			false,
			&MsgUpdateDid{Id: ""},
			"Id: is not DID",
		},
		{
			false,
			&MsgUpdateDid{Id: "did:ch:test:alice"},
			"Id: is not DID",
		},
		{
			false,
			&MsgUpdateDid{Id: "did:cheqd:test:alice"},
			"The message must contain either a Controller or a Authentication: bad request",
		},
		{
			false,
			&MsgUpdateDid{Id: "did:cheqd:test:alice",
				VerificationMethod: []*VerificationMethod{},
			},
			"The message must contain either a Controller or a Authentication: bad request",
		},
		{
			false,
			&MsgUpdateDid{Id: "did:cheqd:test:alice",
				VerificationMethod: []*VerificationMethod{{}},
			},
			"index 0, value : : is not DID fragment: invalid verification method",
		},
		{
			false,
			&MsgUpdateDid{Id: "did:cheqd:test:alice", Authentication: []string{"dd"}},
			"Authentication item dd: is not DID fragment",
		},
		{
			false,
			&MsgUpdateDid{Id: "did:cheqd:test:alice", Authentication: []string{""}},
			"Authentication item : is not DID fragment",
		},
		{
			false,
			&MsgUpdateDid{Id: "did:cheqd:test:alice", Authentication: []string{"did:cheqd:test:alice"}},
			"Authentication item did:cheqd:test:alice: is not DID fragment",
		},
		{
			false,
			&MsgUpdateDid{
				Id:             "did:cheqd:test:alice",
				Authentication: []string{"did:cheqd:test:alice#key-1"},
				Controller:     []string{"did:cheqd:test:alice"},
			},
			"did:cheqd:test:alice#key-1: verification method not found",
		},
		{
			false,
			&MsgUpdateDid{Id: "did:cheqd:test:alice", CapabilityInvocation: []string{"dd"}},
			"CapabilityInvocation item dd: is not DID fragment",
		},
		{
			false,
			&MsgUpdateDid{Id: "did:cheqd:test:alice", CapabilityInvocation: []string{""}},
			"CapabilityInvocation item : is not DID fragment",
		},
		{
			false,
			&MsgUpdateDid{Id: "did:cheqd:test:alice", CapabilityInvocation: []string{"did:cheqd:test:alice"}},
			"CapabilityInvocation item did:cheqd:test:alice: is not DID fragment",
		},
		{
			false,
			&MsgUpdateDid{
				Id:                   "did:cheqd:test:alice",
				CapabilityInvocation: []string{"did:cheqd:test:alice#key-1"},
				Controller:           []string{"did:cheqd:test:alice"},
			},
			"did:cheqd:test:alice#key-1: verification method not found",
		},
		{
			false,
			&MsgUpdateDid{Id: "did:cheqd:test:alice", CapabilityDelegation: []string{"dd"}},
			"CapabilityDelegation item dd: is not DID fragment",
		},
		{
			false,
			&MsgUpdateDid{Id: "did:cheqd:test:alice", CapabilityDelegation: []string{""}},
			"CapabilityDelegation item : is not DID fragment",
		},
		{
			false,
			&MsgUpdateDid{Id: "did:cheqd:test:alice", CapabilityDelegation: []string{"did:cheqd:test:alice"}},
			"CapabilityDelegation item did:cheqd:test:alice: is not DID fragment",
		},
		{
			false,
			&MsgUpdateDid{
				Id:                   "did:cheqd:test:alice",
				CapabilityDelegation: []string{"did:cheqd:test:alice#key-1"},
				Controller:           []string{"did:cheqd:test:alice"},
			},
			"did:cheqd:test:alice#key-1: verification method not found",
		},
		{
			false,
			&MsgUpdateDid{Id: "did:cheqd:test:alice", KeyAgreement: []string{"dd"}},
			"KeyAgreement item dd: is not DID fragment",
		},
		{
			false,
			&MsgUpdateDid{Id: "did:cheqd:test:alice", KeyAgreement: []string{""}},
			"KeyAgreement item : is not DID fragment",
		},
		{
			false,
			&MsgUpdateDid{Id: "did:cheqd:test:alice", KeyAgreement: []string{"did:cheqd:test:alice"}},
			"KeyAgreement item did:cheqd:test:alice: is not DID fragment",
		},
		{
			false,
			&MsgUpdateDid{
				Id:           "did:cheqd:test:alice",
				KeyAgreement: []string{"did:cheqd:test:alice#key-1"},
				Controller:   []string{"did:cheqd:test:alice"},
			},
			"did:cheqd:test:alice#key-1: verification method not found",
		},
		{
			false,
			&MsgUpdateDid{
				Id:           "did:cheqd:test:alice",
				KeyAgreement: []string{"did:cheqd:test:alice#key-1"},
				Controller:   []string{"did:cheqd::alice"},
			},
			"Controller item did:cheqd::alice at position 0: is not DID",
		},
		{
			false,
			&MsgUpdateDid{
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
			&MsgUpdateDid{
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
			&MsgUpdateDid{
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
			&MsgUpdateDid{
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
			&MsgUpdateDid{
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
			&MsgUpdateDid{
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
			&MsgUpdateDid{
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
			&MsgUpdateDid{
				Id:         "did:cheqd:test:alice",
				Controller: []string{"did:cheqd:test:alice"},
				Service: []*ServiceEndpoint{
					{},
				},
			},
			"index 0, value : : is not DID fragment: invalid service",
		},
		{
			false,
			&MsgUpdateDid{
				Id:         "did:cheqd:test:alice",
				Controller: []string{"did:cheqd:test:alice"},
				Service: []*ServiceEndpoint{
					{
						Id: "weqweqw",
					},
				},
			},
			"index 0, value weqweqw: weqweqw: is not DID fragment: invalid service",
		},
		{
			false,
			&MsgUpdateDid{
				Id:         "did:cheqd:test:alice",
				Controller: []string{"did:cheqd:test:alice"},
				Service: []*ServiceEndpoint{
					{
						Id: "#service-1",
					},
				},
			},
			"index 0, value #service-1: : unsupported service type: bad request: invalid service",
		},
		{
			false,
			&MsgUpdateDid{
				Id:         "did:cheqd:test:alice",
				Controller: []string{"did:cheqd:test:alice"},
				Service: []*ServiceEndpoint{
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
			&MsgUpdateDid{
				Id:         "did:cheqd:test:alice",
				Controller: []string{"did:cheqd:test:alice"},
				Service: []*ServiceEndpoint{
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
			&MsgUpdateDid{
				Id:                 "did:cheqd:test:alice",
				Controller:         []string{"did:cheqd:test:alice"},
				VerificationMethod: []*VerificationMethod{},
			},
			"",
		},
		{
			true,
			&MsgUpdateDid{
				Id:         "did:cheqd:test:alice",
				Controller: []string{"did:cheqd:test:alice"},
				Service:    []*ServiceEndpoint{},
			},
			"",
		},
		{
			true,
			&MsgUpdateDid{
				Id:         "did:cheqd:test:alice",
				Controller: []string{"did:cheqd:test:alice"}},
			"",
		},
		{
			true,
			&MsgUpdateDid{
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
