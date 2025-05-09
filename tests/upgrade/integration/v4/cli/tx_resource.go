package cli

import (
	"encoding/json"
	"path/filepath"

	integrationhelpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
	"github.com/cheqd/cheqd-node/x/did/client/cli"
	types "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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
	args = append(args, integrationhelpers.GenerateFees(fees)...)

	return Tx(container, CliBinaryName, "resource", "create", OperatorAccounts[container], args...)
}
