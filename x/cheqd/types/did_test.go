package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestVerificationMethodValidation(t *testing.T) {
	cases := []struct {
		name    string
		struct_ *VerificationMethod
		valid   bool
		errMsg  string
	}{
		{
			"Id is required",
			&VerificationMethod{
				Id:                 "",
				Type:               "",
				Controller:         "",
				PublicKeyJwk:       nil,
				PublicKeyMultibase: "key_content",
			},
			false,
			"Key: 'VerificationMethod.Id' Error:Field validation for 'Id' failed on the 'required' tag",
		},
		//{
		//	"Id must be a DID",
		//	&VerificationMethod{
		//		Id:                 "abba",
		//		Type:               "",
		//		Controller:         "",
		//		PublicKeyJwk:       nil,
		//		PublicKeyMultibase: "key_content",
		//	},
		//	false,
		//	"",
		//},
		//{
		//	"Valid verification method",
		//	&VerificationMethod{
		//		Id:                 "did:cheqd:alternet:TG9yZW0gaXBzdW0g",
		//		Type:               "",
		//		Controller:         "",
		//		PublicKeyJwk:       nil,
		//		PublicKeyMultibase: "key_content",
		//	},
		//	true,
		//	"",
		//},
	}

	validate, err := BuildValidator("", nil)
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {

			err := validate.Struct(testCase.struct_)

			if testCase.valid {
				require.Nil(t, err)
			} else {
				require.Error(t, err)
				require.Equal(t, testCase.errMsg, err.Error())
			}
		})
	}
}
