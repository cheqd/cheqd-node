package cli

import (
	"encoding/base64"
	"encoding/json"
	"path/filepath"

	integrationhelpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
	"github.com/cheqd/cheqd-node/x/did/client/cli"
	types "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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

	args = append(args, GasParams...)

	return Tx(container, CliBinaryName, "resource", "create-resource", OperatorAccounts[container], args...)
}

func CreateResource(payload types.MsgCreateResourcePayload, resourceFile string, signInputs []cli.SignInput, container, fees string) (sdk.TxResponse, error) {
	resourceFileName := filepath.Base(resourceFile)
	payloadFileName := "payload.json"

	payloadJSON, err := integrationhelpers.Codec.MarshalJSON(&payload)
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

	_, err = LocalnetExecExec(container, "/bin/bash", "-c", "echo '"+string(payload.Data)+"' > "+resourceFileName)
	if err != nil {
		return sdk.TxResponse{}, err
	}
	args := []string{payloadFileName}
	args = append(args, resourceFileName)
	args = append(args, payload.Id)
	args = append(args, integrationhelpers.GenerateFees(fees)...)

	return Tx(container, CliBinaryName, "resource", "create", OperatorAccounts[container], args...)
}
