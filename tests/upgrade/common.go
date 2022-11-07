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
var (
	CURRENT_HEIGHT    int64
	VOTING_END_HEIGHT int64
	UPGRADE_HEIGHT    int64
	HEIGHT_ERROR	  error
)

var DidDoc cheqdtypes.MsgCreateDidPayload

var (
	SignInputs []cheqdcli.SignInput
	Err        error
)

var (
	ResourcePayload resourcetypes.MsgCreateResourcePayload
	ResourceFile    string
	ResourceFileErr error
	ResourceErr     error
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

func GenerateDidDocWithSignInputs() (cheqdtypes.MsgCreateDidPayload, []cheqdcli.SignInput, error) {
	did := "did:cheqd:" + network.DID_NAMESPACE + ":" + uuid.NewString()
	keyId := did + "#key1"

	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return cheqdtypes.MsgCreateDidPayload{}, []cheqdcli.SignInput{}, err
	}

	pubKeyMultibase58, err := multibase.Encode(multibase.Base58BTC, pubKey)
	if err != nil {
		return cheqdtypes.MsgCreateDidPayload{}, []cheqdcli.SignInput{}, err
	}

	payload := cheqdtypes.MsgCreateDidPayload{
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

	input := []cheqdcli.SignInput{
		{
			VerificationMethodId: keyId,
			PrivKey:              privKey,
		},
	}
	return payload, input, nil
}

func GenerateRotatedKeysDidDocWithSignInputs(payload cheqdtypes.MsgCreateDidPayload, input []cheqdcli.SignInput, versionId string) (cheqdtypes.MsgUpdateDidPayload, []cheqdcli.SignInput, error) {
	// Specifically, we want to update the DID doc by rotating keys.

	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return cheqdtypes.MsgUpdateDidPayload{}, []cheqdcli.SignInput{}, err
	}

	pubKeyMultibase58, err := multibase.Encode(multibase.Base58BTC, pubKey)
	if err != nil {
		return cheqdtypes.MsgUpdateDidPayload{}, []cheqdcli.SignInput{}, err
	}

	updatedPayload := cheqdtypes.MsgUpdateDidPayload{
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

	updatedInput := []cheqdcli.SignInput{
		input[0],
		{
			VerificationMethodId: input[0].VerificationMethodId,
			PrivKey:              privKey,
		},
	}

	return updatedPayload, updatedInput, nil
}

func GenerateResource(didDoc cheqdtypes.MsgCreateDidPayload) (resourcetypes.MsgCreateResourcePayload, error) {
	collectionId := strings.Replace(didDoc.Id, "did:cheqd:"+network.DID_NAMESPACE+":", "", 1)
	payload := resourcetypes.MsgCreateResourcePayload{
		CollectionId: collectionId,
		Id:           uuid.NewString(),
		Name:         "TestResource",
		ResourceType: "TestType",
		Data:         []byte(integrationtestdata.JSON_FILE_CONTENT),
	}

	return payload, nil
}
