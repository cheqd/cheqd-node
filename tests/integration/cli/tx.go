package cli

import (
	"encoding/base64"
	"strings"

	"github.com/cheqd/cheqd-node/tests/integration/helpers"
	"github.com/cheqd/cheqd-node/x/cheqd/client/cli"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var CLI_TX_PARAMS = []string{
	"--chain-id",
	CHAIN_ID,
	"--keyring-backend",
	KEYRING_BACKEND,
	"--output",
	OUTPUT_FORMAT,
	"--gas",
	GAS,
	"--gas-adjustment",
	GAS_ADJUSTMENT,
	"--gas-prices",
	GAS_PRICES,
	"--yes",
}

func SendTx(module, txName, from string, otherArgs ...string) (sdk.TxResponse, error) {
	args := []string{"tx", module, txName}

	// Common params
	args = append(args, CLI_TX_PARAMS...)

	// Cosmos account
	args = append(args, "--from", from)

	// Other args
	args = append(args, otherArgs...)

	output, err := Exec(args...)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	// Skip 'gas estimate: xxx' string
	output = strings.Split(output, "\n")[1]

	var resp sdk.TxResponse

	err = helpers.Codec.UnmarshalJSON([]byte(output), &resp)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	return resp, nil
}

func CreateDid(payload types.MsgCreateDidPayload, signInputs []cli.SignInput, from string) (sdk.TxResponse, error) {
	// Payload
	payloadJson, err := helpers.Codec.MarshalJSON(&payload)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	args := []string{string(payloadJson)}

	// Sign inputs
	for _, signInput := range signInputs {
		args = append(args, signInput.VerificationMethodId)
		args = append(args, base64.StdEncoding.EncodeToString(signInput.PrivKey))
	}

	return SendTx("cheqd", "create-did", from, args...)
}
