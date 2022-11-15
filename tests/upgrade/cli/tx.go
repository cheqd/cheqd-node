package cli

import (
	"fmt"
	"strconv"
	"strings"

	integrationhelpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func Tx(container string, binary string, module, tx, from string, txArgs ...string) (sdk.TxResponse, error) {
	args := []string{
		binary,
		"tx",
		module,
		tx,
	}

	// Common params
	args = append(args, TX_PARAMS...)

	// Cosmos account
	args = append(args, "--from", from)

	// Other args
	args = append(args, txArgs...)

	output, err := LocalnetExecExec(container, args...)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	// Skip 'gas estimate: xxx' string
	perLineOutput := strings.Split(output, "\n")
	if len(perLineOutput) > 3 {
		output = TrimExtraLineOffset(output, 2)
	} else {
		output = TrimExtraLineOffset(output, 1)
	}

	var resp sdk.TxResponse

	err = integrationhelpers.Codec.UnmarshalJSON([]byte(output), &resp)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	return resp, nil
}

func SubmitUpgradeProposal(upgradeHeight int64, container string) (sdk.TxResponse, error) {
	fmt.Println("Submitting upgrade proposal from", container)
	args := append([]string{
		CLI_BINARY_NAME,
		"tx", "gov", "submit-proposal", "software-upgrade",
		UPGRADE_NAME,
		"--title", "Upgrade Title",
		"--description", "Upgrade Description",
		"--upgrade-height", strconv.FormatInt(upgradeHeight, 10),
		"--upgrade-info", "Upgrade Info",
		"--from", OperatorAccounts[container],
	}, TX_PARAMS...)

	out, err := LocalnetExecExec(container, args...)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	out = TrimExtraLineOffset(out, 1)

	var resp sdk.TxResponse

	err = integrationhelpers.Codec.UnmarshalJSON([]byte(out), &resp)
	if err != nil {
		fmt.Println("JSON unmarshal error: output:", out)
		return sdk.TxResponse{}, err
	}
	return resp, nil
}

func DepositGov(container string) (sdk.TxResponse, error) {
	fmt.Println("Depositing from", container)
	args := append([]string{
		CLI_BINARY_NAME,
		"tx", "gov", "deposit", "1", DEPOSIT_AMOUNT,
		"--from", OperatorAccounts[container],
	}, TX_PARAMS...)

	out, err := LocalnetExecExec(container, args...)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	out = TrimExtraLineOffset(out, 1)

	var resp sdk.TxResponse

	err = integrationhelpers.Codec.UnmarshalJSON([]byte(out), &resp)
	if err != nil {
		return sdk.TxResponse{}, err
	}
	return resp, nil
}

func VoteUpgradeProposal(container string) (sdk.TxResponse, error) {
	fmt.Println("Voting from", container)
	args := append([]string{
		CLI_BINARY_NAME,
		"tx", "gov", "vote", "1", "yes",
		"--from", OperatorAccounts[container],
	}, TX_PARAMS...)

	out, err := LocalnetExecExec(container, args...)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	out = TrimExtraLineOffset(out, 1)

	var resp sdk.TxResponse

	err = integrationhelpers.Codec.UnmarshalJSON([]byte(out), &resp)
	if err != nil {
		return sdk.TxResponse{}, err
	}
	return resp, nil
}
