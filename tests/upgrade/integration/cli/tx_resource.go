package cli

import (
	"encoding/base64"
	"encoding/json"

	"github.com/cheqd/cheqd-node/x/did/client/cli"
	types "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CreateResourceOptions struct {
	CollectionId    string                  `json:"collection_id"`
	ResourceId      string                  `json:"resource_id"`
	ResourceName    string                  `json:"resource_name"`
	ResourceVersion string                  `json:"resource_version"`
	ResourceType    string                  `json:"resource_type"`
	ResourceFile    string                  `json:"resource_file"`
	AlsoKnownAs     []*types.AlternativeUri `json:"also_known_as"`
}

func CreateResourceLegacy(collectionId string, resourceId string, resourceName string, resourceType string, resourceFile string, signInputs []cli.SignInput, container string) (sdk.TxResponse, error) {
	args := []string{
		"--collection-id", collectionId,
		"--resource-id", resourceId,
		"--resource-name", resourceName,
		"--resource-type", resourceType,
		"--resource-file", resourceFile,
	}

	for _, signInput := range signInputs {
		args = append(args, signInput.VerificationMethodId)
		args = append(args, base64.StdEncoding.EncodeToString(signInput.PrivKey))
	}

	return Tx(container, CLI_BINARY_NAME, "resource", "create-resource", OperatorAccounts[container], args...)
}

func CreateResource(collectionId string, resourceId string, resourceName string, resourceVersion string, resourceType string, resourceFile string, resourceAsKnownAs []*types.AlternativeUri, signInputs []cli.SignInput, container string) (sdk.TxResponse, error) {
	resourceOptions := CreateResourceOptions{
		CollectionId:    collectionId,
		ResourceId:      resourceId,
		ResourceName:    resourceName,
		ResourceVersion: resourceVersion,
		ResourceType:    resourceType,
		ResourceFile:    resourceFile,
		AlsoKnownAs:     resourceAsKnownAs,
	}

	payloadJson, err := json.Marshal(&resourceOptions)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	payloadWithSignInputs := cli.PayloadWithSignInputs{
		Payload:    payloadJson,
		SignInputs: signInputs,
	}

	payloadWithSignInputsJson, err := json.Marshal(&payloadWithSignInputs)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	args := []string{string(payloadWithSignInputsJson)}

	return Tx(container, CLI_BINARY_NAME, "resource", "create-resource", OperatorAccounts[container], args...)
}
