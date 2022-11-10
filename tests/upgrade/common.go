//go:build upgrade

package upgrade

import (
	"crypto/ed25519"
	"strings"

	integrationtestdata "github.com/cheqd/cheqd-node/tests/integration/testdata"
	network "github.com/cheqd/cheqd-node/tests/upgrade/network"
	didcli "github.com/cheqd/cheqd-node/x/did/client/cli"
	didtypes "github.com/cheqd/cheqd-node/x/did/types/v1"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/google/uuid"
	"github.com/multiformats/go-multibase"
)

// Pre
var DidDoc didtypes.MsgCreateDidPayload

var (
	SignInputs []didcli.SignInput
	Err        = GenerateDidDocWithSignInputs(&DidDoc, &SignInputs)
)

var (
	ResourcePayload resourcetypes.MsgCreateResourcePayload
	ResourceFile    string
	ResourceFileErr error
	ResourceErr     = GenerateResource(&ResourcePayload)
)

var (
	RotatedKeysDidDoc     didtypes.MsgUpdateDidPayload
	RotatedKeysSignInputs []didcli.SignInput
	RotatedKeysErr        error
)

// Post
var PostDidDoc didtypes.MsgCreateDidPayload

var (
	PostSignInputs []didcli.SignInput
	PostErr        error
)

var (
	PostResourcePayload resourcetypes.MsgCreateResourcePayload
	PostResourceFile    string
	PostResourceFileErr error
	PostResourceErr     error
)

var (
	PostRotatedKeysDidDoc     didtypes.MsgUpdateDidPayload
	PostRotatedKeysSignInputs []didcli.SignInput
	PostRotatedKeysErr        error
)

// Migration
var (
	QueriedDidDoc   didtypes.Did
	QueriedResource resourcetypes.Resource
)

func GenerateDidDocWithSignInputs(payload *didtypes.MsgCreateDidPayload, input *[]didcli.SignInput) error {
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

	payload = &didtypes.MsgCreateDidPayload{
		Id:         did,
		Controller: []string{did},
		VerificationMethod: []*didtypes.VerificationMethod{
			{
				Id:                 keyId,
				Type:               "Ed25519VerificationKey2020",
				Controller:         did,
				PublicKeyMultibase: string(pubKeyMultibase58),
			},
		},
		Authentication: []string{keyId},
	}

	input = &[]didcli.SignInput{
		{
			VerificationMethodId: keyId,
			PrivKey:              privKey,
		},
	}
	return nil
}

func GenerateRotatedKeysDidDocWithSignInputs(payload *didtypes.MsgCreateDidPayload, updatedPayload *didtypes.MsgUpdateDidPayload, input *[]didcli.SignInput, updatedInput *[]didcli.SignInput, versionId string) error {
	// Specifically, we want to update the DID doc by rotating keys.

	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return err
	}

	pubKeyMultibase58, err := multibase.Encode(multibase.Base58BTC, pubKey)
	if err != nil {
		return err
	}

	updatedPayload = &didtypes.MsgUpdateDidPayload{
		Id:         payload.Id,
		Controller: []string{payload.Id},
		VerificationMethod: []*didtypes.VerificationMethod{
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

	updatedInput = &[]didcli.SignInput{
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

func ResetDidDocInMem(payload *didtypes.MsgCreateDidPayload, input *[]didcli.SignInput) error {
	payload = &didtypes.MsgCreateDidPayload{}
	input = &[]didcli.SignInput{}
	Err = GenerateDidDocWithSignInputs(payload, input)
	return Err
}

func ResetRotatedKeysDidDocInMem(updatedPayload *didtypes.MsgUpdateDidPayload, updatedInput *[]didcli.SignInput) error {
	updatedPayload = &didtypes.MsgUpdateDidPayload{}
	updatedInput = &[]didcli.SignInput{}
	return nil
}

func ResetResourceInMem(payload *resourcetypes.MsgCreateResourcePayload) error {
	payload = &resourcetypes.MsgCreateResourcePayload{}
	return nil
}
