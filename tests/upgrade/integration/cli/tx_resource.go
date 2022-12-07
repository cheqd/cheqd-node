package cli

import (
	"encoding/base64"
	"encoding/json"
	"path/filepath"

	// integrationhelpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
	"github.com/cheqd/cheqd-node/x/did/client/cli"
	types "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	// . "github.com/onsi/ginkgo/v2"
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

func CreateResource(msg types.MsgCreateResourcePayload, resourceFile string, signInputs []cli.SignInput, container string) (sdk.TxResponse, error) {
	resourceFileName := filepath.Base(resourceFile)
	payloadFileName := "payload.json"

	resourceOptions := CreateResourceOptions{
		CollectionId:    msg.CollectionId,
		ResourceId:      msg.Id,
		ResourceName:    msg.Name,
		ResourceVersion: msg.Version,
		ResourceType:    msg.ResourceType,
		ResourceFile:    resourceFileName,
		AlsoKnownAs:     msg.AlsoKnownAs,
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

	_, err = LocalnetExecExec(container, "/bin/bash", "-c", "echo '"+string(payloadWithSignInputsJson)+"' > "+payloadFileName)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	_, err = LocalnetExecExec(container, "/bin/bash", "-c", "echo '"+string(msg.Data)+"' > "+resourceFileName)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	return Tx(container, CLI_BINARY_NAME, "resource", "create", OperatorAccounts[container], payloadFileName)
}
