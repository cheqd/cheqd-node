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
	CollectionID    string                  `json:"collection_id"`
	ResourceID      string                  `json:"resource_id"`
	ResourceName    string                  `json:"resource_name"`
	ResourceVersion string                  `json:"resource_version"`
	ResourceType    string                  `json:"resource_type"`
	ResourceFile    string                  `json:"resource_file"`
	AlsoKnownAs     []*types.AlternativeUri `json:"also_known_as"`
}

func CreateResourceLegacy(collectionID string, resourceID string, resourceName string, resourceType string, resourceFile string, signInputs []cli.SignInput, container string) (sdk.TxResponse, error) {
	args := []string{
		"--collection-id", collectionID,
		"--resource-id", resourceID,
		"--resource-name", resourceName,
		"--resource-type", resourceType,
		"--resource-file", resourceFile,
	}

	for _, signInput := range signInputs {
		args = append(args, signInput.VerificationMethodID)
		args = append(args, base64.StdEncoding.EncodeToString(signInput.PrivKey))
	}

	return Tx(container, CLIBinaryName, "resource", "create-resource", OperatorAccounts[container], args...)
}

func CreateResource(msg types.MsgCreateResourcePayload, resourceFile string, signInputs []cli.SignInput, container string) (sdk.TxResponse, error) {
	resourceFileName := filepath.Base(resourceFile)
	payloadFileName := "payload.json"

	resourceOptions := CreateResourceOptions{
		CollectionID:    msg.CollectionId,
		ResourceID:      msg.Id,
		ResourceName:    msg.Name,
		ResourceVersion: msg.Version,
		ResourceType:    msg.ResourceType,
		ResourceFile:    resourceFileName,
		AlsoKnownAs:     msg.AlsoKnownAs,
	}

	payloadJSON, err := json.Marshal(&resourceOptions)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	payloadWithSignInputs := cli.PayloadWithSignInputs{
		Payload:    payloadJSON,
		SignInputs: signInputs,
	}

	payloadWithSignInputsJSON, err := json.Marshal(&payloadWithSignInputs)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	_, err = LocalnetExecExec(container, "/bin/bash", "-c", "echo '"+string(payloadWithSignInputsJSON)+"' > "+payloadFileName)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	_, err = LocalnetExecExec(container, "/bin/bash", "-c", "echo '"+string(msg.Data)+"' > "+resourceFileName)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	return Tx(container, CLIBinaryName, "resource", "create", OperatorAccounts[container], payloadFileName)
}
