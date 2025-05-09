package cli

import (
	"encoding/json"
	"fmt"

	integrationhelpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
	"github.com/cheqd/cheqd-node/x/did/client/cli"
	"github.com/cheqd/cheqd-node/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CreateDid(payload cli.DIDDocument, signInputs []cli.SignInput, container, versionID, fees string) (sdk.TxResponse, error) {
	innerPayloadJSON, err := json.Marshal(&payload)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	outerPayload := cli.PayloadWithSignInputs{
		Payload:    innerPayloadJSON,
		SignInputs: signInputs,
	}

	outerPayloadJSON, err := json.Marshal(&outerPayload)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	out, err := LocalnetExecExec(container, "/bin/bash", "-c", "echo '"+string(outerPayloadJSON)+"' > payload.json")
	if err != nil {
		return sdk.TxResponse{}, err
	}

	fmt.Println(out)

	args := []string{string("payload.json")}

	args = append(args, integrationhelpers.GenerateFees(fees)...)

	if versionID != "" {
		args = append(args, versionID)
	}

	return Tx(container, CliBinaryName, "cheqd", "create-did", OperatorAccounts[container], args...)
}

func UpdateDid(payload cli.DIDDocument, signInputs []cli.SignInput, container, versionID, fees string) (sdk.TxResponse, error) {
	innerPayloadJSON, err := json.Marshal(&payload)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	outerPayload := cli.PayloadWithSignInputs{
		Payload:    innerPayloadJSON,
		SignInputs: signInputs,
	}

	outerPayloadJSON, err := json.Marshal(&outerPayload)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	out, err := LocalnetExecExec(container, "/bin/bash", "-c", "echo '"+string(outerPayloadJSON)+"' > payload.json")
	if err != nil {
		return sdk.TxResponse{}, err
	}

	fmt.Println(out)

	args := []string{string("payload.json")}

	args = append(args, integrationhelpers.GenerateFees(fees)...)

	if versionID != "" {
		args = append(args, versionID)
	}

	return Tx(container, CliBinaryName, "cheqd", "update-did", OperatorAccounts[container], args...)
}

func DeactivateDid(payload types.MsgDeactivateDidDocPayload, signInputs []cli.SignInput, container string, fees string) (sdk.TxResponse, error) {
	innerPayloadJSON := integrationhelpers.Codec.MustMarshalJSON(&payload)

	outerPayload := cli.PayloadWithSignInputs{
		Payload:    innerPayloadJSON,
		SignInputs: signInputs,
	}

	outerPayloadJSON, err := json.Marshal(&outerPayload)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	out, err := LocalnetExecExec(container, "/bin/bash", "-c", "echo '"+string(outerPayloadJSON)+"' > payload.json")
	if err != nil {
		return sdk.TxResponse{}, err
	}

	fmt.Println(out)

	args := []string{string("payload.json")}

	args = append(args, integrationhelpers.GenerateFees(fees)...)

	return Tx(container, CliBinaryName, "cheqd", "deactivate-did", OperatorAccounts[container], args...)
}
