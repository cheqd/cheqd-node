//go:build upgrade

package upgrade

import (
	"crypto/ed25519"
	"strings"

	integrationtestdata "github.com/cheqd/cheqd-node/tests/integration/testdata"
	network "github.com/cheqd/cheqd-node/tests/upgrade/network"
	cheqdcli "github.com/cheqd/cheqd-node/x/cheqd/client/cli"
	cheqdtypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/google/uuid"
	"github.com/multiformats/go-multibase"
)

// Pre
var DidDoc cheqdtypes.MsgCreateDidPayload

var (
	SignInputs []cheqdcli.SignInput
	Err        = GenerateDidDocWithSignInputs(&DidDoc, &SignInputs)
)

var (
	ResourcePayload resourcetypes.MsgCreateResourcePayload
	ResourceFile    string
	ResourceFileErr error
	ResourceErr     = GenerateResource(&ResourcePayload)
)

var (
	RotatedKeysDidDoc     cheqdtypes.MsgUpdateDidPayload
	RotatedKeysSignInputs []cheqdcli.SignInput
	RotatedKeysErr        error
)

// Post
var PostDidDoc cheqdtypes.MsgCreateDidPayload

var (
	PostSignInputs []cheqdcli.SignInput
	PostErr        error
)

var (
	PostResourcePayload resourcetypes.MsgCreateResourcePayload
	PostResourceFile    string
	PostResourceFileErr error
	PostResourceErr     error
)

var (
	PostRotatedKeysDidDoc     cheqdtypes.MsgUpdateDidPayload
	PostRotatedKeysSignInputs []cheqdcli.SignInput
	PostRotatedKeysErr        error
)

// Migration
var (
	QueriedDidDoc   cheqdtypes.Did
	QueriedResource resourcetypes.Resource
)

func GenerateDidDocWithSignInputs(payload *cheqdtypes.MsgCreateDidPayload, input *[]cheqdcli.SignInput) error {
	did := "did:cheqd:" + network.DID_NAMESPACE + ":" + uuid.NewString()
	keyId := did + "#key1"

	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return err
	}

	pubKeyMultibase58, err := multibase.Encode(multibase.Base58BTC, pubKey)
	if err != nil {
		return err
	}

	payload = &cheqdtypes.MsgCreateDidPayload{
		Id:         did,
		Controller: []string{did},
		VerificationMethod: []*cheqdtypes.VerificationMethod{
			{
				Id:                 keyId,
				Type:               "Ed25519VerificationKey2020",
				Controller:         did,
				PublicKeyMultibase: string(pubKeyMultibase58),
			},
		},
		Authentication: []string{keyId},
	}

	input = &[]cheqdcli.SignInput{
		{
			VerificationMethodId: keyId,
			PrivKey:              privKey,
		},
	}
	return nil
}

func GenerateRotatedKeysDidDocWithSignInputs(payload *cheqdtypes.MsgCreateDidPayload, updatedPayload *cheqdtypes.MsgUpdateDidPayload, input *[]cheqdcli.SignInput, updatedInput *[]cheqdcli.SignInput, versionId string) error {
	// Specifically, we want to update the DID doc by rotating keys.

	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return err
	}

	pubKeyMultibase58, err := multibase.Encode(multibase.Base58BTC, pubKey)
	if err != nil {
		return err
	}

	updatedPayload = &cheqdtypes.MsgUpdateDidPayload{
		Id:         payload.Id,
		Controller: []string{payload.Id},
		VerificationMethod: []*cheqdtypes.VerificationMethod{
			{
				Id:                 payload.VerificationMethod[0].Id,
				Type:               "Ed25519VerificationKey2020",
				Controller:         payload.Id,
				PublicKeyMultibase: string(pubKeyMultibase58),
			},
		},
		Authentication: []string{payload.VerificationMethod[0].Id},
		VersionId:      versionId,
	}

	updatedInput = &[]cheqdcli.SignInput{
		(*input)[0],
		{
			VerificationMethodId: (*input)[0].VerificationMethodId,
			PrivKey:              privKey,
		},
	}

	return nil
}

func GenerateResource(payload *resourcetypes.MsgCreateResourcePayload) error {
	collectionId := strings.Replace(DidDoc.Id, "did:cheqd:"+network.DID_NAMESPACE, "", 1)
	payload = &resourcetypes.MsgCreateResourcePayload{
		CollectionId: collectionId,
		Id:           uuid.NewString(),
		Name:         "TestResource",
		ResourceType: "TestType",
		Data:         []byte(integrationtestdata.JSON_FILE_CONTENT),
	}

	return nil
}

func ResetDidDocInMem(payload *cheqdtypes.MsgCreateDidPayload, input *[]cheqdcli.SignInput) error {
	payload = &cheqdtypes.MsgCreateDidPayload{}
	input = &[]cheqdcli.SignInput{}
	Err = GenerateDidDocWithSignInputs(payload, input)
	return Err
}

func ResetRotatedKeysDidDocInMem(updatedPayload *cheqdtypes.MsgUpdateDidPayload, updatedInput *[]cheqdcli.SignInput) error {
	updatedPayload = &cheqdtypes.MsgUpdateDidPayload{}
	updatedInput = &[]cheqdcli.SignInput{}
	return nil
}

func ResetResourceInMem(payload *resourcetypes.MsgCreateResourcePayload) error {
	payload = &resourcetypes.MsgCreateResourcePayload{}
	return nil
}
