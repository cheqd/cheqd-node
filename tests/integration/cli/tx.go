package cli

import (
	"encoding/json"

	"github.com/cheqd/cheqd-node/tests/integration/helpers"
	"github.com/cheqd/cheqd-node/tests/integration/network"
	"github.com/cheqd/cheqd-node/x/did/client/cli"
	"github.com/cheqd/cheqd-node/x/did/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	FlagVersionID = "--version-id"
)

var CLITxParams = []string{
	"--chain-id", network.ChainID,
	"--keyring-backend", KeyringBackend,
	"--output", OutputFormat,
	"--yes",
}

var CliGasParams = []string{
	"--gas", Gas,
	"--gas-adjustment", GasAdjustment,
	"--gas-prices", GasPrices,
}

func Tx(module, tx, from string, feeParams []string, txArgs ...string) (sdk.TxResponse, error) {
	args := []string{"tx", module, tx}

	// Common params
	args = append(args, CLITxParams...)

	// Fee params
	args = append(args, feeParams...)

	// Cosmos account
	args = append(args, "--from", from)

	// Other args
	args = append(args, txArgs...)

	output, err := Exec(args...)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	// Skip 'gas estimate: xxx' string, trim 'Successfully migrated key' string
	output = helpers.TrimImportedStdout(output)

	var resp sdk.TxResponse

	err = helpers.Codec.UnmarshalJSON([]byte(output), &resp)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	return resp, nil
}

func GrantFees(granter, grantee string, feeParams []string) (sdk.TxResponse, error) {
	return Tx("feegrant", "grant", granter, feeParams, granter, grantee)
}

func RevokeFeeGrant(granter, grantee string, feeParams []string) (sdk.TxResponse, error) {
	return Tx("feegrant", "revoke", granter, feeParams, granter, grantee)
}

func CreateDidDoc(tmpDir string, payload cli.DIDDocument, signInputs []cli.SignInput, versionID, from string, feeParams []string) (sdk.TxResponse, error) {
	payloadJSON, err := json.Marshal(&payload)
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

	payloadFile := helpers.MustWriteTmpFile(tmpDir, payloadWithSignInputsJSON)

	if versionID != "" {
		return Tx("cheqd", "create-did", from, feeParams, payloadFile, FlagVersionID, versionID)
	}

	return Tx("cheqd", "create-did", from, feeParams, payloadFile)
}

func UpdateDidDoc(tmpDir string, payload cli.DIDDocument, signInputs []cli.SignInput, versionID, from string, feeParams []string) (sdk.TxResponse, error) {
	payloadJSON, err := json.Marshal(&payload)
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

	payloadFile := helpers.MustWriteTmpFile(tmpDir, payloadWithSignInputsJSON)

	if versionID != "" {
		return Tx("cheqd", "update-did", from, feeParams, payloadFile, FlagVersionID, versionID)
	}

	return Tx("cheqd", "update-did", from, feeParams, payloadFile)
}

func DeactivateDidDoc(tmpDir string, payload types.MsgDeactivateDidDocPayload, signInputs []cli.SignInput, versionID, from string, feeParams []string) (sdk.TxResponse, error) {
	payloadJSON, err := helpers.Codec.MarshalJSON(&payload)
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

	payloadFile := helpers.MustWriteTmpFile(tmpDir, payloadWithSignInputsJSON)

	if versionID != "" {
		return Tx("cheqd", "deactivate-did", from, feeParams, payloadFile, FlagVersionID, versionID)
	}

	return Tx("cheqd", "deactivate-did", from, feeParams, payloadFile)
}

func CreateResource(tmpDir string, payload resourcetypes.MsgCreateResourcePayload, signInputs []cli.SignInput, dataFile, from string, feeParams []string) (sdk.TxResponse, error) {
	payloadJSON, err := helpers.Codec.MarshalJSON(&payload)
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

	payloadFile := helpers.MustWriteTmpFile("", payloadWithSignInputsJSON)

	return Tx("resource", "create", from, feeParams, payloadFile, dataFile)
}

func BurnMsg(from string, coins string, feeParams []string) (sdk.TxResponse, error) {
	return Tx("cheqd", "burn", from, feeParams, coins)
}

func SubmitProposalTx(from, pathToDir string, feeParams []string) (sdk.TxResponse, error) {
	return Tx("gov", "submit-proposal", from, feeParams, pathToDir)
}

func VoteProposalTx(from, option, id string, feeParams []string) (sdk.TxResponse, error) {
	return Tx("gov", "vote", from, feeParams, option, id)
}
