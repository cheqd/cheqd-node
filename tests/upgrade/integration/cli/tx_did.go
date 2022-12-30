package cli

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	integrationhelpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
	"github.com/cheqd/cheqd-node/x/did/client/cli"
	didtypesv2 "github.com/cheqd/cheqd-node/x/did/types"
	didtypesv1 "github.com/cheqd/cheqd-node/x/did/types/v1"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CreateDidLegacy(payload didtypesv1.MsgCreateDidPayload, signInputs []cli.SignInput, container string) (sdk.TxResponse, error) {
	payloadJSON, err := integrationhelpers.Codec.MarshalJSON(&payload)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	args := []string{string(payloadJSON)}

	for _, signInput := range signInputs {
		args = append(args, signInput.VerificationMethodID)
		args = append(args, base64.StdEncoding.EncodeToString(signInput.PrivKey))
	}

	return Tx(container, CLIBinaryName, "cheqd", "create-did", OperatorAccounts[container], args...)
}

func CreateDid(payload didtypesv2.MsgCreateDidDocPayload, signInputs []cli.SignInput, container string) (sdk.TxResponse, error) {
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

	return Tx(container, CLIBinaryName, "cheqd", "create-did", OperatorAccounts[container], args...)
}

func UpdateDidLegacy(payload didtypesv1.MsgUpdateDidPayload, signInputs []cli.SignInput, container string) (sdk.TxResponse, error) {
	payloadJSON, err := integrationhelpers.Codec.MarshalJSON(&payload)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	args := []string{string(payloadJSON)}

	for _, signInput := range signInputs {
		args = append(args, signInput.VerificationMethodID)
		args = append(args, base64.StdEncoding.EncodeToString(signInput.PrivKey))
	}

	return Tx(container, CLIBinaryName, "cheqd", "update-did", OperatorAccounts[container], args...)
}

func UpdateDid(payload didtypesv2.MsgUpdateDidDocPayload, signInputs []cli.SignInput, container string) (sdk.TxResponse, error) {
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

	return Tx(container, CLIBinaryName, "cheqd", "update-did", OperatorAccounts[container], args...)
}

func DeactivateDidLegacy(payload didtypesv1.MsgDeactivateDidPayload, signInputs []cli.SignInput, container string) (sdk.TxResponse, error) {
	payloadJSON, err := integrationhelpers.Codec.MarshalJSON(&payload)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	args := []string{string(payloadJSON)}

	for _, signInput := range signInputs {
		args = append(args, signInput.VerificationMethodID)
		args = append(args, base64.StdEncoding.EncodeToString(signInput.PrivKey))
	}

	return Tx(container, CLIBinaryName, "cheqd", "deactivate-did", OperatorAccounts[container], args...)
}

func DeactivateDid(payload didtypesv2.MsgDeactivateDidDocPayload, signInputs []cli.SignInput, container string) (sdk.TxResponse, error) {
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

	return Tx(container, CLIBinaryName, "cheqd", "deactivate-did", OperatorAccounts[container], args...)
}
