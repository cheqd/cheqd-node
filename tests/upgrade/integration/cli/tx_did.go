package cli

import (
	"encoding/base64"
	"encoding/json"

	integrationhelpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
	didtypesv2 "github.com/cheqd/cheqd-node/x/did/types"
	didtypesv1 "github.com/cheqd/cheqd-node/x/did/types/v1"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PayloadWithSignInputs struct {
	Payload    json.RawMessage
	SignInputs []SignInput
}

type SignInput struct {
	VerificationMethodId string `json:"verificationMethodId"`
	PrivateKey           []byte `json:"privateKey"`
}

func CreateDidLegacy(payload didtypesv1.MsgCreateDidPayload, signInputs []SignInput, container string) (sdk.TxResponse, error) {
	payloadJson, err := integrationhelpers.Codec.MarshalJSON(&payload)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	args := []string{string(payloadJson)}

	for _, signInput := range signInputs {
		args = append(args, signInput.VerificationMethodId)
		args = append(args, base64.StdEncoding.EncodeToString(signInput.PrivateKey))
	}

	return Tx(container, CLI_BINARY_NAME, "cheqd", "create-did", OperatorAccounts[container], args...)
}

func CreateDid(payload didtypesv2.MsgCreateDidDocPayload, signInputs []SignInput, container string) (sdk.TxResponse, error) {
	payloadJson, err := integrationhelpers.Codec.MarshalJSON(&payload)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	payloadWithSignInput := PayloadWithSignInputs{
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

func UpdateDidLegacy(payload didtypesv1.MsgUpdateDidPayload, signInputs []SignInput, container string) (sdk.TxResponse, error) {
	payloadJson, err := integrationhelpers.Codec.MarshalJSON(&payload)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	args := []string{string(payloadJson)}

	for _, signInput := range signInputs {
		args = append(args, signInput.VerificationMethodId)
		args = append(args, base64.StdEncoding.EncodeToString(signInput.PrivateKey))
	}

	return Tx(container, CLI_BINARY_NAME, "cheqd", "update-did", OperatorAccounts[container], args...)
}

func UpdateDid(payload didtypesv2.MsgUpdateDidDocPayload, signInputs []SignInput, container string) (sdk.TxResponse, error) {
	payloadJson, err := integrationhelpers.Codec.MarshalJSON(&payload)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	args := []string{string(payloadJson)}

	for _, signInput := range signInputs {
		args = append(args, signInput.VerificationMethodId)
		args = append(args, base64.StdEncoding.EncodeToString(signInput.PrivateKey))
	}

	return Tx(container, CLI_BINARY_NAME, "cheqd", "update-did", OperatorAccounts[container], args...)
}

func DeactivateDidLegacy(payload didtypesv1.MsgDeactivateDidPayload, signInputs []SignInput, container string) (sdk.TxResponse, error) {
	payloadJson, err := integrationhelpers.Codec.MarshalJSON(&payload)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	args := []string{string(payloadJson)}

	for _, signInput := range signInputs {
		args = append(args, signInput.VerificationMethodId)
		args = append(args, base64.StdEncoding.EncodeToString(signInput.PrivateKey))
	}

	return Tx(container, CLI_BINARY_NAME, "cheqd", "deactivate-did", OperatorAccounts[container], args...)
}

func DeactivateDid(payload didtypesv1.MsgDeactivateDidPayload, signInputs []SignInput, container string) (sdk.TxResponse, error) {
	payloadJson, err := integrationhelpers.Codec.MarshalJSON(&payload)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	args := []string{string(payloadJson)}

	for _, signInput := range signInputs {
		args = append(args, signInput.VerificationMethodId)
		args = append(args, base64.StdEncoding.EncodeToString(signInput.PrivateKey))
	}

	return Tx(container, CLI_BINARY_NAME, "cheqd", "deactivate-did", OperatorAccounts[container], args...)
}