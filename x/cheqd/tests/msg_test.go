package tests

import (
	"fmt"
	"testing"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/stretchr/testify/require"
)

const Prefix = "did:cheqd:test:"

func TestNewMsgCreateDidValidation(t *testing.T) {
	cases := []struct {
		valid  bool
		name   string
		msg    *types.MsgCreateDid
		errMsg string
	}{
		{true, "Valid Create Did Msg", types.NewMsgCreateDid(&types.MsgCreateDidPayload{Id: "1"}, []*types.SignInfo{{VerificationMethodId: "foo", Signature: "bar"}}), ""},
		{false, "Payload is missed", types.NewMsgCreateDid(nil, nil), "Payload: is required"},
		{false, "Signature are missed", types.NewMsgCreateDid(&types.MsgCreateDidPayload{Id: "1"}, nil), "Signatures: are required"},
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
		msg    *types.MsgUpdateDid
		errMsg string
	}{
		{true, "Valid Update Did Msg", types.NewMsgUpdateDid(&types.MsgUpdateDidPayload{Id: "1"}, []*types.SignInfo{{VerificationMethodId: "foo", Signature: "bar"}}), ""},
		{false, "Payload is missed", types.NewMsgUpdateDid(nil, nil), "Payload: is required"},
		{false, "Signatures are missed", types.NewMsgUpdateDid(&types.MsgUpdateDidPayload{Id: "1"}, nil), "Signatures: are required"},
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
		msg    *types.MsgCreateDidPayload
		errMsg string
	}{
		{
			false,
			&types.MsgCreateDidPayload{},
			"Id: is not DID",
		},
		{
			false,
			&types.MsgCreateDidPayload{Id: ""},
			"Id: is not DID",
		},
		{
			false,
			&types.MsgCreateDidPayload{Id: "did:ch:test:alice"},
			"Id: is not DID",
		},
		{
			false,
			&types.MsgCreateDidPayload{Id: AliceDID},
			"The message must contain either a Controller or a Authentication: bad request",
		},
		{
			false,
			&types.MsgCreateDidPayload{Id: AliceDID,
				VerificationMethod: []*types.VerificationMethod{},
			},
			"The message must contain either a Controller or a Authentication: bad request",
		},
		{
			false,
			&types.MsgCreateDidPayload{Id: AliceDID,
				VerificationMethod: []*types.VerificationMethod{{}},
			},
			"index 0, value : : is not DID fragment: invalid verification method",
		},
		{
			false,
			&types.MsgCreateDidPayload{Id: AliceDID, Authentication: []string{"dd"}},
			"Authentication item dd: is not DID fragment",
		},
		{
			false,
			&types.MsgCreateDidPayload{Id: AliceDID, Authentication: []string{""}},
			"Authentication item : is not DID fragment",
		},
		{
			false,
			&types.MsgCreateDidPayload{Id: AliceDID, Authentication: []string{AliceDID}},
			fmt.Sprintf("Authentication item %v: is not DID fragment", AliceDID),
		},
		{
			false,
			&types.MsgCreateDidPayload{
				Id:             AliceDID,
				Authentication: []string{AliceKey1},
				Controller:     []string{AliceDID},
			},
			fmt.Sprintf("%v: verification method not found", AliceKey1),
		},
		{
			false,
			&types.MsgCreateDidPayload{Id: AliceDID, CapabilityInvocation: []string{"dd"}},
			"CapabilityInvocation item dd: is not DID fragment",
		},
		{
			false,
			&types.MsgCreateDidPayload{Id: AliceDID, CapabilityInvocation: []string{""}},
			"CapabilityInvocation item : is not DID fragment",
		},
		{
			false,
			&types.MsgCreateDidPayload{Id: AliceDID, CapabilityInvocation: []string{AliceDID}},
			fmt.Sprintf("CapabilityInvocation item %v: is not DID fragment", AliceDID),
		},
		{
			false,
			&types.MsgCreateDidPayload{
				Id:                   AliceDID,
				CapabilityInvocation: []string{AliceKey1},
				Controller:           []string{AliceDID},
			},
			fmt.Sprintf("%v: verification method not found", AliceKey1),
		},
		{
			false,
			&types.MsgCreateDidPayload{Id: AliceDID, CapabilityDelegation: []string{"dd"}},
			"CapabilityDelegation item dd: is not DID fragment",
		},
		{
			false,
			&types.MsgCreateDidPayload{Id: AliceDID, CapabilityDelegation: []string{""}},
			"CapabilityDelegation item : is not DID fragment",
		},
		{
			false,
			&types.MsgCreateDidPayload{Id: AliceDID, CapabilityDelegation: []string{AliceDID}},
			fmt.Sprintf("CapabilityDelegation item %v: is not DID fragment", AliceDID),
		},
		{
			false,
			&types.MsgCreateDidPayload{
				Id:                   AliceDID,
				CapabilityDelegation: []string{AliceKey1},
				Controller:           []string{AliceDID},
			},
			fmt.Sprintf("%v: verification method not found", AliceKey1),
		},
		{
			false,
			&types.MsgCreateDidPayload{Id: AliceDID, KeyAgreement: []string{"dd"}},
			"KeyAgreement item dd: is not DID fragment",
		},
		{
			false,
			&types.MsgCreateDidPayload{Id: AliceDID, KeyAgreement: []string{""}},
			"KeyAgreement item : is not DID fragment",
		},
		{
			false,
			&types.MsgCreateDidPayload{Id: AliceDID, KeyAgreement: []string{AliceDID}},
			fmt.Sprintf("KeyAgreement item %v: is not DID fragment", AliceDID),
		},
		{
			false,
			&types.MsgCreateDidPayload{
				Id:           AliceDID,
				KeyAgreement: []string{AliceKey1},
				Controller:   []string{AliceDID},
			},
			fmt.Sprintf("%v: verification method not found", AliceKey1),
		},
		{
			false,
			&types.MsgCreateDidPayload{
				Id:           AliceDID,
				KeyAgreement: []string{AliceKey1},
				Controller:   []string{"did:cheqd::alice"},
			},
			"Controller item did:cheqd::alice at position 0: is not DID",
		},
		{
			false,
			&types.MsgCreateDidPayload{
				Id: AliceDID,
				VerificationMethod: []*types.VerificationMethod{
					{Id: "dasda"},
				},
				Controller: []string{AliceDID},
			},
			"index 0, value dasda: dasda: is not DID fragment: invalid verification method",
		},
		{
			false,
			&types.MsgCreateDidPayload{
				Id: AliceDID,
				VerificationMethod: []*types.VerificationMethod{
					{Id: AliceKey1},
				},
				Controller: []string{AliceDID},
			},
			fmt.Sprintf("index 0, value %v: : unsupported verification method type: bad request: invalid verification method", AliceKey1),
		},
		{
			false,
			&types.MsgCreateDidPayload{
				Id: AliceDID,
				VerificationMethod: []*types.VerificationMethod{
					{Id: AliceKey1, Type: "YES"},
				},
				Controller: []string{AliceDID},
			},
			fmt.Sprintf("index 0, value %v: YES: unsupported verification method type: bad request: invalid verification method", AliceKey1),
		},
		{
			false,
			&types.MsgCreateDidPayload{
				Id: AliceDID,
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:   AliceKey1,
						Type: "JsonWebKey2020",
					},
				},
				Controller: []string{AliceDID},
			},
			fmt.Sprintf("index 0, value %v: JsonWebKey2020: should contain `PublicKeyJwk` verification material property: bad request: invalid verification method", AliceKey1),
		},
		{
			false,
			&types.MsgCreateDidPayload{
				Id: AliceDID,
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:                 AliceKey1,
						Type:               "JsonWebKey2020",
						PublicKeyMultibase: "tetetet",
					},
				},
				Controller: []string{AliceDID},
			},
			fmt.Sprintf("index 0, value %v: JsonWebKey2020: should contain `PublicKeyJwk` verification material property: bad request: invalid verification method", AliceKey1),
		},
		{
			false,
			&types.MsgCreateDidPayload{
				Id: AliceDID,
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:                 AliceKey1,
						Type:               "Ed25519VerificationKey2020",
						PublicKeyMultibase: "tetetet",
						Controller:         AliceDID,
					},
					{
						Id:                 AliceKey1,
						Type:               "Ed25519VerificationKey2020",
						PublicKeyMultibase: "tetetet",
						Controller:         AliceDID,
					},
				},
				Controller: []string{AliceDID},
			},
			fmt.Sprintf("%v is duplicated: invalid verification method", AliceKey1),
		},
		{
			false,
			&types.MsgCreateDidPayload{
				Id: AliceDID,
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:                 AliceKey1,
						Type:               "Ed25519VerificationKey2020",
						PublicKeyMultibase: "tetetet",
						Controller:         AliceDID,
					},
					{
						Id:   AliceKey2,
						Type: "JsonWebKey2020",
						PublicKeyJwk: []*types.KeyValuePair{
							{
								Key:   "x",
								Value: "sadad",
							},
						},
						Controller: AliceDID,
					},
					{
						Id:                 AliceKey3,
						Type:               "JsonWebKey20212",
						PublicKeyMultibase: "tetetet",
						Controller:         AliceDID,
					},
				},
				Controller: []string{AliceDID},
			},
			fmt.Sprintf("index 2, value %v: JsonWebKey20212: unsupported verification method type: bad request: invalid verification method", AliceKey3),
		},
		{
			false,
			&types.MsgCreateDidPayload{
				Id:         AliceDID,
				Controller: []string{AliceDID},
				Service: []*types.Service{
					{},
				},
			},
			"index 0, value : : is not DID fragment: invalid service",
		},
		{
			false,
			&types.MsgCreateDidPayload{
				Id:         AliceDID,
				Controller: []string{AliceDID},
				Service: []*types.Service{
					{
						Id: "weqweqw",
					},
				},
			},
			"index 0, value weqweqw: weqweqw: is not DID fragment: invalid service",
		},
		{
			false,
			&types.MsgCreateDidPayload{
				Id:         AliceDID,
				Controller: []string{AliceDID},
				Service: []*types.Service{
					{
						Id: "#service-1",
					},
				},
			},
			"index 0, value #service-1: : unsupported service type: bad request: invalid service",
		},
		{
			false,
			&types.MsgCreateDidPayload{
				Id:         AliceDID,
				Controller: []string{AliceDID},
				Service: []*types.Service{
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
			&types.MsgCreateDidPayload{
				Id:         AliceDID,
				Controller: []string{AliceDID},
				Service: []*types.Service{
					{
						Id:   AliceService1,
						Type: "DIDCommMessaging",
					},
					{
						Id:   "#service-1",
						Type: "DIDCommMessaging",
					},
				},
			},
			fmt.Sprintf("%v is duplicated: invalid service", AliceService1),
		},
		{
			true,
			&types.MsgCreateDidPayload{
				Id:                 AliceDID,
				Controller:         []string{AliceDID},
				VerificationMethod: []*types.VerificationMethod{},
			},
			"",
		},
		{
			true,
			&types.MsgCreateDidPayload{
				Id:         AliceDID,
				Controller: []string{AliceDID},
				Service:    []*types.Service{},
			},
			"",
		},
		{
			true,
			&types.MsgCreateDidPayload{
				Id:         AliceDID,
				Controller: []string{AliceDID}},
			"",
		},
		{
			true,
			&types.MsgCreateDidPayload{
				Id:             AliceDID,
				Controller:     []string{AliceDID, BobDID},
				Authentication: []string{"#key-1", AliceKey2},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:   AliceKey1,
						Type: "JsonWebKey2020",
						PublicKeyJwk: []*types.KeyValuePair{
							{
								Key:   "x",
								Value: "sadad",
							},
						},
						Controller: AliceDID,
					},
					{
						Id:                 AliceKey2,
						Type:               "Ed25519VerificationKey2020",
						PublicKeyMultibase: "tetetet",
						Controller:         AliceDID,
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
		msg    *types.MsgUpdateDidPayload
		errMsg string
	}{
		{
			false,
			&types.MsgUpdateDidPayload{},
			"Id: is not DID",
		},
		{
			false,
			&types.MsgUpdateDidPayload{Id: ""},
			"Id: is not DID",
		},
		{
			false,
			&types.MsgUpdateDidPayload{Id: "did:ch:test:alice"},
			"Id: is not DID",
		},
		{
			false,
			&types.MsgUpdateDidPayload{Id: AliceDID},
			"The message must contain either a Controller or a Authentication: bad request",
		},
		{
			false,
			&types.MsgUpdateDidPayload{Id: AliceDID,
				VerificationMethod: []*types.VerificationMethod{},
			},
			"The message must contain either a Controller or a Authentication: bad request",
		},
		{
			false,
			&types.MsgUpdateDidPayload{Id: AliceDID,
				VerificationMethod: []*types.VerificationMethod{{}},
			},
			"index 0, value : : is not DID fragment: invalid verification method",
		},
		{
			false,
			&types.MsgUpdateDidPayload{Id: AliceDID, Authentication: []string{"dd"}},
			"Authentication item dd: is not DID fragment",
		},
		{
			false,
			&types.MsgUpdateDidPayload{Id: AliceDID, Authentication: []string{""}},
			"Authentication item : is not DID fragment",
		},
		{
			false,
			&types.MsgUpdateDidPayload{Id: AliceDID, Authentication: []string{AliceDID}},
			fmt.Sprintf("Authentication item %v: is not DID fragment", AliceDID),
		},
		{
			false,
			&types.MsgUpdateDidPayload{
				Id:             AliceDID,
				Authentication: []string{AliceKey1},
				Controller:     []string{AliceDID},
			},
			fmt.Sprintf("%v: verification method not found", AliceKey1),
		},
		{
			false,
			&types.MsgUpdateDidPayload{Id: AliceDID, CapabilityInvocation: []string{"dd"}},
			"CapabilityInvocation item dd: is not DID fragment",
		},
		{
			false,
			&types.MsgUpdateDidPayload{Id: AliceDID, CapabilityInvocation: []string{""}},
			"CapabilityInvocation item : is not DID fragment",
		},
		{
			false,
			&types.MsgUpdateDidPayload{Id: AliceDID, CapabilityInvocation: []string{AliceDID}},
			fmt.Sprintf("CapabilityInvocation item %v: is not DID fragment", AliceDID),
		},
		{
			false,
			&types.MsgUpdateDidPayload{
				Id:                   AliceDID,
				CapabilityInvocation: []string{AliceKey1},
				Controller:           []string{AliceDID},
			},
			fmt.Sprintf("%v: verification method not found", AliceKey1),
		},
		{
			false,
			&types.MsgUpdateDidPayload{Id: AliceDID, CapabilityDelegation: []string{"dd"}},
			"CapabilityDelegation item dd: is not DID fragment",
		},
		{
			false,
			&types.MsgUpdateDidPayload{Id: AliceDID, CapabilityDelegation: []string{""}},
			"CapabilityDelegation item : is not DID fragment",
		},
		{
			false,
			&types.MsgUpdateDidPayload{Id: AliceDID, CapabilityDelegation: []string{AliceDID}},
			fmt.Sprintf("CapabilityDelegation item %v: is not DID fragment", AliceDID),
		},
		{
			false,
			&types.MsgUpdateDidPayload{
				Id:                   AliceDID,
				CapabilityDelegation: []string{AliceKey1},
				Controller:           []string{AliceDID},
			},
			fmt.Sprintf("%v: verification method not found", AliceKey1),
		},
		{
			false,
			&types.MsgUpdateDidPayload{Id: AliceDID, KeyAgreement: []string{"dd"}},
			"KeyAgreement item dd: is not DID fragment",
		},
		{
			false,
			&types.MsgUpdateDidPayload{Id: AliceDID, KeyAgreement: []string{""}},
			"KeyAgreement item : is not DID fragment",
		},
		{
			false,
			&types.MsgUpdateDidPayload{Id: AliceDID, KeyAgreement: []string{AliceDID}},
			fmt.Sprintf("KeyAgreement item %v: is not DID fragment", AliceDID),
		},
		{
			false,
			&types.MsgUpdateDidPayload{
				Id:           AliceDID,
				KeyAgreement: []string{AliceKey1},
				Controller:   []string{AliceDID},
			},
			fmt.Sprintf("%v: verification method not found", AliceKey1),
		},
		{
			false,
			&types.MsgUpdateDidPayload{
				Id:           AliceDID,
				KeyAgreement: []string{AliceKey1},
				Controller:   []string{"did:cheqd::alice"},
			},
			"Controller item did:cheqd::alice at position 0: is not DID",
		},
		{
			false,
			&types.MsgUpdateDidPayload{
				Id: AliceDID,
				VerificationMethod: []*types.VerificationMethod{
					{Id: "dasda"},
				},
				Controller: []string{AliceDID},
			},
			"index 0, value dasda: dasda: is not DID fragment: invalid verification method",
		},
		{
			false,
			&types.MsgUpdateDidPayload{
				Id: AliceDID,
				VerificationMethod: []*types.VerificationMethod{
					{Id: AliceKey1},
				},
				Controller: []string{AliceDID},
			},
			fmt.Sprintf("index 0, value %v: : unsupported verification method type: bad request: invalid verification method", AliceKey1),
		},
		{
			false,
			&types.MsgUpdateDidPayload{
				Id: AliceDID,
				VerificationMethod: []*types.VerificationMethod{
					{Id: AliceKey1, Type: "YES"},
				},
				Controller: []string{AliceDID},
			},
			fmt.Sprintf("index 0, value %v: YES: unsupported verification method type: bad request: invalid verification method", AliceKey1),
		},
		{
			false,
			&types.MsgUpdateDidPayload{
				Id: AliceDID,
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:   AliceKey1,
						Type: "Ed25519VerificationKey2020",
					},
				},
				Controller: []string{AliceDID},
			},
			fmt.Sprintf("index 0, value %v: Ed25519VerificationKey2020: should contain `PublicKeyMultibase` verification material property: bad request: invalid verification method", AliceKey1),
		},
		{
			false,
			&types.MsgUpdateDidPayload{
				Id: AliceDID,
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:   AliceKey1,
						Type: "JsonWebKey2020",
						PublicKeyJwk: []*types.KeyValuePair{
							{Key: "x", Value: "y"},
						},
					},
				},
				Controller: []string{AliceDID},
			},
			fmt.Sprintf("index 0, value %v: Controller: is required: invalid verification method", AliceKey1),
		},
		{
			false,
			&types.MsgUpdateDidPayload{
				Id: AliceDID,
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:                 AliceKey1,
						Type:               "Ed25519VerificationKey2020",
						PublicKeyMultibase: "tetetet",
						Controller:         AliceDID,
					},
					{
						Id:                 AliceKey1,
						Type:               "Ed25519VerificationKey2020",
						PublicKeyMultibase: "tetetet",
						Controller:         AliceDID,
					},
				},
				Controller: []string{AliceDID},
			},
			fmt.Sprintf("%v is duplicated: invalid verification method", AliceKey1),
		},
		{
			false,
			&types.MsgUpdateDidPayload{
				Id: AliceDID,
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:                 AliceKey1,
						Type:               "Ed25519VerificationKey2020",
						PublicKeyMultibase: "tetetet",
						Controller:         AliceDID,
					},
					{
						Id:                 AliceKey2,
						Type:               "Ed25519VerificationKey2020",
						PublicKeyMultibase: "tetetet",
						Controller:         AliceDID,
					},
					{
						Id:                 AliceKey3,
						Type:               "JsonWebKey20212",
						PublicKeyMultibase: "tetetet",
						Controller:         AliceDID,
					},
				},
				Controller: []string{AliceDID},
			},
			fmt.Sprintf("index 2, value %v: JsonWebKey20212: unsupported verification method type: bad request: invalid verification method", AliceKey3),
		},
		{
			false,
			&types.MsgUpdateDidPayload{
				Id:         AliceDID,
				Controller: []string{AliceDID},
				Service: []*types.Service{
					{},
				},
			},
			"index 0, value : : is not DID fragment: invalid service",
		},
		{
			false,
			&types.MsgUpdateDidPayload{
				Id:         AliceDID,
				Controller: []string{AliceDID},
				Service: []*types.Service{
					{
						Id: "weqweqw",
					},
				},
			},
			"index 0, value weqweqw: weqweqw: is not DID fragment: invalid service",
		},
		{
			false,
			&types.MsgUpdateDidPayload{
				Id:         AliceDID,
				Controller: []string{AliceDID},
				Service: []*types.Service{
					{
						Id: "#service-1",
					},
				},
			},
			"index 0, value #service-1: : unsupported service type: bad request: invalid service",
		},
		{
			false,
			&types.MsgUpdateDidPayload{
				Id:         AliceDID,
				Controller: []string{AliceDID},
				Service: []*types.Service{
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
			&types.MsgUpdateDidPayload{
				Id:         AliceDID,
				Controller: []string{AliceDID},
				Service: []*types.Service{
					{
						Id:   AliceService1,
						Type: "DIDCommMessaging",
					},
					{
						Id:   "#service-1",
						Type: "DIDCommMessaging",
					},
				},
			},
			fmt.Sprintf("%v is duplicated: invalid service", AliceService1),
		},
		{
			true,
			&types.MsgUpdateDidPayload{
				Id:                 AliceDID,
				Controller:         []string{AliceDID},
				VerificationMethod: []*types.VerificationMethod{},
			},
			"",
		},
		{
			true,
			&types.MsgUpdateDidPayload{
				Id:         AliceDID,
				Controller: []string{AliceDID},
				Service:    []*types.Service{},
			},
			"",
		},
		{
			true,
			&types.MsgUpdateDidPayload{
				Id:         AliceDID,
				Controller: []string{AliceDID}},
			"",
		},
		{
			true,
			&types.MsgUpdateDidPayload{
				Id:             AliceDID,
				Controller:     []string{AliceDID, BobDID},
				Authentication: []string{"#key-1", AliceKey2},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:                 AliceKey1,
						Type:               "Ed25519VerificationKey2020",
						PublicKeyMultibase: "tetetet",
						Controller:         AliceDID,
					},
					{
						Id:                 AliceKey2,
						Type:               "Ed25519VerificationKey2020",
						PublicKeyMultibase: "tetetet",
						Controller:         AliceDID,
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
