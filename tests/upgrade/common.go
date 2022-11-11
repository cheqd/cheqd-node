//go:build upgrade

package upgrade

import (
	"crypto/ed25519"
	"strings"

	integrationtestdata "github.com/cheqd/cheqd-node/tests/integration/testdata"
	cli "github.com/cheqd/cheqd-node/tests/upgrade/cli"
	network "github.com/cheqd/cheqd-node/tests/upgrade/network"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/google/uuid"
	"github.com/multiformats/go-multibase"
)

var (
	ExistingDidDocCreatePayloads, ExistingDidDocUpdatePayloads, ExistingDidDocDeactivatePayloads          []string
	ExistingSignInputCreatePayloads, ExistingSignInputUpdatePayloads, ExistingSignInputDeactivatePayloads []string
	ExistingResourceCreatePayloads                                                                        []string
)

// Pre
var (
	CURRENT_HEIGHT    int64
	VOTING_END_HEIGHT int64
	UPGRADE_HEIGHT    int64
	HEIGHT_ERROR      error
)

var DidDoc didtypes.MsgCreateDidDocPayload

var (
	SignInputs []cli.SignInput
	Err        error
)

var (
	ResourcePayload resourcetypes.MsgCreateResourcePayload
	ResourceFile    string
	ResourceFileErr error
	ResourceErr     error
)

var (
	RotatedKeysDidDoc     didtypes.MsgUpdateDidDocPayload
	RotatedKeysSignInputs []cli.SignInput
	RotatedKeysErr        error
)

// Post
var PostDidDoc didtypes.MsgCreateDidDocPayload

var (
	PostSignInputs []cli.SignInput
	PostErr        error
)

var (
	PostResourcePayload resourcetypes.MsgCreateResourcePayload
	PostResourceFile    string
	PostResourceFileErr error
	PostResourceErr     error
)

var (
	PostRotatedKeysDidDoc     didtypes.MsgUpdateDidDocPayload
	PostRotatedKeysSignInputs []cli.SignInput
	PostRotatedKeysErr        error
)

// Migration
var (
	QueriedDidDoc   didtypes.DidDoc
	QueriedResource resourcetypes.Resource
)

func GenerateDidDocWithSignInputs() (didtypes.MsgCreateDidDocPayload, []cli.SignInput, error) {
	did := "did:cheqd:" + network.DID_NAMESPACE + ":" + uuid.NewString()
	keyId := did + "#key1"

	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return didtypes.MsgCreateDidDocPayload{}, []cli.SignInput{}, err
	}

	pubKeyMultibase58, err := multibase.Encode(multibase.Base58BTC, pubKey)
	if err != nil {
		return didtypes.MsgCreateDidDocPayload{}, []cli.SignInput{}, err
	}

	payload := didtypes.MsgCreateDidDocPayload{
		Id:         did,
		Controller: []string{did},
		VerificationMethod: []*didtypes.VerificationMethod{
			{
				Id:                   keyId,
				Type:                 "Ed25519VerificationKey2020",
				Controller:           did,
				VerificationMaterial: string(pubKeyMultibase58),
			},
		},
		Authentication: []string{keyId},
	}

	input := []cli.SignInput{
		{
			VerificationMethodId: keyId,
			PrivateKey:           privKey,
		},
	}
	return payload, input, nil
}

func GenerateRotatedKeysDidDocWithSignInputs(payload didtypes.MsgCreateDidDocPayload, input []cli.SignInput, versionId string) (didtypes.MsgUpdateDidDocPayload, []cli.SignInput, error) {
	// Specifically, we want to update the DID doc by rotating keys.

	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return didtypes.MsgUpdateDidDocPayload{}, []cli.SignInput{}, err
	}

	pubKeyMultibase58, err := multibase.Encode(multibase.Base58BTC, pubKey)
	if err != nil {
		return didtypes.MsgUpdateDidDocPayload{}, []cli.SignInput{}, err
	}

	updatedPayload := didtypes.MsgUpdateDidDocPayload{
		Id:         payload.Id,
		Controller: []string{payload.Id},
		VerificationMethod: []*didtypes.VerificationMethod{
			{
				Id:                   payload.VerificationMethod[0].Id,
				Type:                 "Ed25519VerificationKey2020",
				Controller:           payload.Id,
				VerificationMaterial: string(pubKeyMultibase58),
			},
		},
		Authentication: []string{payload.VerificationMethod[0].Id},
		VersionId:      versionId,
	}

	updatedInput := []cli.SignInput{
		input[0],
		{
			VerificationMethodId: input[0].VerificationMethodId,
			PrivateKey:           privKey,
		},
	}

	return updatedPayload, updatedInput, nil
}

func GenerateResource(didDoc didtypes.MsgCreateDidDocPayload) (resourcetypes.MsgCreateResourcePayload, error) {
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
