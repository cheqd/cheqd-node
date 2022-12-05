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
	payloadJson, err := integrationhelpers.Codec.MarshalJSON(&payload)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	args := []string{string(payloadJson)}

	for _, signInput := range signInputs {
		args = append(args, signInput.VerificationMethodId)
		args = append(args, base64.StdEncoding.EncodeToString(signInput.PrivKey))
	}

	return Tx(container, CLI_BINARY_NAME, "cheqd", "create-did", OperatorAccounts[container], args...)
}

func CreateDid(payload didtypesv2.MsgCreateDidDocPayload, signInputs []cli.SignInput, container string) (sdk.TxResponse, error) {
	payloadJson, err := integrationhelpers.Codec.MarshalJSON(&payload)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	payloadWithSignInput := cli.PayloadWithSignInputs{
		Payload:    payloadJson,
		SignInputs: signInputs,
	}

	payloadWithSignInputJson, err := json.Marshal(&payloadWithSignInput)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	args := []string{string(payloadWithSignInputJson)}

	return Tx(container, CLI_BINARY_NAME, "cheqd", "create-did", OperatorAccounts[container], args...)
}

func UpdateDidLegacy(payload didtypesv1.MsgUpdateDidPayload, signInputs []cli.SignInput, container string) (sdk.TxResponse, error) {
	payloadJson, err := integrationhelpers.Codec.MarshalJSON(&payload)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	args := []string{string(payloadJson)}

	for _, signInput := range signInputs {
		args = append(args, signInput.VerificationMethodId)
		args = append(args, base64.StdEncoding.EncodeToString(signInput.PrivKey))
	}

	return Tx(container, CLI_BINARY_NAME, "cheqd", "update-did", OperatorAccounts[container], args...)
}

func UpdateDid(payload didtypesv2.MsgUpdateDidDocPayload, signInputs []cli.SignInput, container string) (sdk.TxResponse, error) {
	innerPayloadJson := integrationhelpers.Codec.MustMarshalJSON(&payload)

	outerPayload := cli.PayloadWithSignInputs{
		Payload:    innerPayloadJson,
		SignInputs: signInputs,
	}

	outerPayloadJson, err := json.Marshal(&outerPayload)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	out, err := LocalnetExecExec(container, "/bin/bash", "-c", "echo '"+string(outerPayloadJson)+"' > payload.json")
	if err != nil {
		return sdk.TxResponse{}, err
	}

	fmt.Println(out)

	args := []string{string("payload.json")}

	return Tx(container, CLI_BINARY_NAME, "cheqd", "update-did", OperatorAccounts[container], args...)
}

func DeactivateDidLegacy(payload didtypesv1.MsgDeactivateDidPayload, signInputs []cli.SignInput, container string) (sdk.TxResponse, error) {
	payloadJson, err := integrationhelpers.Codec.MarshalJSON(&payload)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	args := []string{string(payloadJson)}

	for _, signInput := range signInputs {
		args = append(args, signInput.VerificationMethodId)
		args = append(args, base64.StdEncoding.EncodeToString(signInput.PrivKey))
	}

	return Tx(container, CLI_BINARY_NAME, "cheqd", "deactivate-did", OperatorAccounts[container], args...)
}

func DeactivateDid(payload didtypesv1.MsgDeactivateDidPayload, signInputs []cli.SignInput, container string) (sdk.TxResponse, error) {
	payloadJson, err := integrationhelpers.Codec.MarshalJSON(&payload)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	args := []string{string(payloadJson)}

	for _, signInput := range signInputs {
		args = append(args, signInput.VerificationMethodId)
		args = append(args, base64.StdEncoding.EncodeToString(signInput.PrivKey))
	}

	return Tx(container, CLI_BINARY_NAME, "cheqd", "deactivate-did", OperatorAccounts[container], args...)
}
