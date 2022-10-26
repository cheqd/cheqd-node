package cli

import (
	"encoding/base64"

	integrationhelpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
	cheqdcli "github.com/cheqd/cheqd-node/x/cheqd/client/cli"
	cheqdtypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CreateDid(payload cheqdtypes.MsgCreateDidPayload, signInputs []cheqdcli.SignInput, container string) (sdk.TxResponse, error) {
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

func UpdateDid(payload cheqdtypes.MsgUpdateDidPayload, signInputs []cheqdcli.SignInput, container string) (sdk.TxResponse, error) {
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
