package types

import (
	ctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewMsgCreateCredDef(t *testing.T) {
	cases := []struct {
		valid  bool
		msg    *MsgCreateCredDef
		errMsg string
	}{
		{
			true,
			NewMsgCreateCredDef(
				"did:cheqd:test:aaaaa/credDef",
				"schema",
				"",
				"CL-Sig-Cred_def",
				[]string{"did:cheqd:test:alice"},
				&MsgCreateCredDef_ClType{ClType: &CredDefValue{Primary: nil, Revocation: nil}}),
			"",
		},
		{
			true,
			NewMsgCreateCredDef(
				"did:cheqd:test:aaaaa/credDef",
				"schema",
				"tag",
				"CL-Sig-Cred_def",
				[]string{"did:cheqd:test:alice"},
				&MsgCreateCredDef_ClType{ClType: &CredDefValue{Primary: nil, Revocation: nil}}),
			"",
		},
		{
			false,
			NewMsgCreateCredDef(
				"/credDef",
				"",
				"",
				"",
				nil,
				&MsgCreateCredDef_ClType{ClType: nil}),
			"Id: is not DID",
		},
		{
			false,
			NewMsgCreateCredDef(
				"did:cheqd:test:alice#key-1/credDef",
				"",
				"",
				"",
				nil,
				&MsgCreateCredDef_ClType{ClType: nil}),
			"Id: is not DID",
		},
		{
			false,
			NewMsgCreateCredDef(
				"did:cheqd:test:alice",
				"",
				"",
				"",
				nil,
				&MsgCreateCredDef_ClType{ClType: &CredDefValue{Primary: nil, Revocation: nil}}),
			"Id must end with resource type '/credDef': bad request",
		},
		{
			false,
			NewMsgCreateCredDef(
				"did:cheqd:test:alice/credDef",
				"",
				"",
				"",
				nil,
				&MsgCreateCredDef_ClType{ClType: &CredDefValue{Primary: nil, Revocation: nil}}),
			"SchemaId: is required",
		},
		{
			false,
			NewMsgCreateCredDef(
				"did:cheqd:test:alice/credDef",
				"schema-1",
				"",
				"",
				nil,
				&MsgCreateCredDef_ClType{ClType: &CredDefValue{Primary: nil, Revocation: nil}}),
			"SignatureType: is required",
		},
		{
			false,
			NewMsgCreateCredDef(
				"did:cheqd:test:alice/credDef",
				"schema-1",
				"",
				"ss",
				[]string{"did:cheqd:test:alice"},
				&MsgCreateCredDef_ClType{ClType: &CredDefValue{Primary: nil, Revocation: nil}}),
			"ss is not allowed signature type: bad request",
		},
		{
			false,
			NewMsgCreateCredDef(
				"did:cheqd:test:alice/credDef",
				"schema-1",
				"",
				"CL-Sig-Cred_def",
				nil,
				&MsgCreateCredDef_ClType{ClType: &CredDefValue{Primary: nil, Revocation: nil}}),
			"Controller: is required",
		},
		{
			false,
			NewMsgCreateCredDef(
				"did:cheqd:test:alice/credDef",
				"schema-1",
				"",
				"CL-Sig-Cred_def",
				[]string{"1"},
				&MsgCreateCredDef_ClType{ClType: &CredDefValue{Primary: nil, Revocation: nil}}),
			"Controller item 1: is not DID",
		},
		{
			false,
			NewMsgCreateCredDef(
				"did:cheqd:test:alice/credDef",
				"schema-1",
				"",
				"CL-Sig-Cred_def",
				[]string{"did:cheqd:test:alice"},
				&MsgCreateCredDef_ClType{ClType: nil}),
			"Value: is required",
		},
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

func TestNewMsgWriteRequestValidation(t *testing.T) {
	cases := []struct {
		valid  bool
		msg    *MsgWriteRequest
		errMsg string
	}{
		{true, NewMsgWriteRequest(&ctypes.Any{TypeUrl: "1", Value: []byte{1}}, nil, map[string]string{"foo": "bar"}), ""},
		{false, NewMsgWriteRequest(&ctypes.Any{TypeUrl: "1"}, nil, map[string]string{"foo": "bar"}), "Invalid Data: it cannot be empty: bad request"},
		{false, NewMsgWriteRequest(&ctypes.Any{Value: []byte{1}}, nil, map[string]string{"foo": "bar"}), "Invalid Data: it cannot be empty: bad request"},
		{false, NewMsgWriteRequest(nil, nil, nil), "Data: is required"},
		{false, NewMsgWriteRequest(&ctypes.Any{TypeUrl: "1", Value: []byte{1}}, nil, nil), "Signatures: is required"},
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

func TestNewMsgCreateSchema(t *testing.T) {
	cases := []struct {
		valid  bool
		msg    *MsgCreateSchema
		errMsg string
	}{
		{
			true,
			NewMsgCreateSchema(
				"did:cheqd:test:aaaaa/schema",
				"CL-Schema",
				"schema",
				"version1",
				[]string{"did:cheqd:test:alice"},
				[]string{"did:cheqd:test:alice"}),
			"",
		},
		{
			false,
			NewMsgCreateSchema(
				"/schema",
				"",
				"",
				"",
				nil,
				nil),
			"Id: is not DID",
		},
		{
			false,
			NewMsgCreateSchema(
				"did:cheqd:test:alice#key-1/schema",
				"",
				"",
				"",
				nil,
				nil),
			"Id: is not DID",
		},
		{
			false,
			NewMsgCreateSchema(
				"did:cheqd:test:alice",
				"",
				"",
				"",
				nil,
				nil),
			"Id must end with resource type '/schema': bad request",
		},
		{
			false,
			NewMsgCreateSchema(
				"did:cheqd:test:alice/schema",
				"",
				"",
				"",
				nil,
				nil),
			"Type: is required",
		},
		{
			false,
			NewMsgCreateSchema(
				"did:cheqd:test:alice/schema",
				"schema-1",
				"",
				"",
				nil,
				nil),
			"schema-1 is not allowed type: bad request",
		},
		{
			false,
			NewMsgCreateSchema(
				"did:cheqd:test:alice/schema",
				"schema-1",
				"",
				"ss",
				[]string{"did:cheqd:test:alice"},
				nil),
			"schema-1 is not allowed type: bad request",
		},
		{
			false,
			NewMsgCreateSchema(
				"did:cheqd:test:alice/schema",
				"CL-Schema",
				"",
				"",
				nil,
				nil),
			"AttrNames: is required",
		},
		{
			false,
			NewMsgCreateSchema(
				"did:cheqd:test:alice/schema",
				"CL-Schema",
				"",
				"",
				[]string{},
				nil),
			"AttrNames: is required",
		},
		{
			false,
			NewMsgCreateSchema(
				"did:cheqd:test:alice/schema",
				"CL-Schema",
				"",
				"",
				[]string{"1", "2"},
				nil),
			"Name: is required",
		},
		{
			false,
			NewMsgCreateSchema(
				"did:cheqd:test:alice/schema",
				"CL-Schema",
				"",
				"",
				make([]string, 126),
				nil),
			"AttrNames: Expected max length 125, got: 126: bad request",
		},
		{
			false,
			NewMsgCreateSchema(
				"did:cheqd:test:alice/schema",
				"CL-Schema",
				"schema",
				"",
				[]string{"1"},
				nil),
			"Version: is required",
		},
		{
			false,
			NewMsgCreateSchema(
				"did:cheqd:test:alice/schema",
				"CL-Schema",
				"schema",
				"version",
				[]string{"1"},
				nil),
			"Controller: is required",
		},
		{
			false,
			NewMsgCreateSchema(
				"did:cheqd:test:alice/schema",
				"CL-Schema",
				"schema",
				"version",
				[]string{"1"},
				[]string{"1"}),
			"Controller item 1: is not DID",
		},
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
