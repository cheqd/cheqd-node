package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	integrationhelpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func Tx(container string, binary string, module, tx, from string, txArgs ...string) (sdk.TxResponse, error) {
	var output string
	var err error

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

	// ToDo: refactor
	if RUN_INSIDE_DOCKER {
		output, err = LocalnetExecExec(container, args...)
	} else {
		output, err = ExecDirectWithHome(container, args...)
	}
	if err != nil {
		return sdk.TxResponse{}, err
	}

	// Skip 'gas estimate: xxx' string, trim 'Successfully migrated key' string
	output = integrationhelpers.TrimImportedStdout(output)

	var resp sdk.TxResponse

	err = integrationhelpers.Codec.UnmarshalJSON([]byte(output), &resp)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	return resp, nil
}

func SubmitParamChangeProposal(container string, pathToDir ...string) (sdk.TxResponse, error) {
	var err error
	var out string

	fmt.Println("Submitting param change proposal from", container)
	args := append([]string{
		CLI_BINARY_NAME,
		"tx", "gov", "submit-legacy-proposal", "param-change", filepath.Join(pathToDir...),
		"--from", OperatorAccounts[container],
	}, TX_PARAMS...)

	// ToDo: refactor
	if RUN_INSIDE_DOCKER {
		out, err = LocalnetExecExec(container, args...)
	} else {
		out, err = ExecDirectWithHome(container, args...)
	}

	if err != nil {
		fmt.Println("Error on submitting ParamChangeProposal", err)
		fmt.Println("Output:", out)
		return sdk.TxResponse{}, err
	}

	// Skip 'gas estimate: xxx' string, trim 'Successfully migrated key' string
	out = integrationhelpers.TrimImportedStdout(out)

	fmt.Println("Output:", out)

	var resp sdk.TxResponse

	err = integrationhelpers.Codec.UnmarshalJSON([]byte(out), &resp)
	if err != nil {
		return sdk.TxResponse{}, err
	}
	return resp, nil
}

func SubmitUpgradeProposal(upgradeHeight int64, container string) (sdk.TxResponse, error) {
	var err error
	var out string
	
	fmt.Println("Submitting upgrade proposal from", container)
	var upgrade_name = os.Getenv("UPGRADE_NAME")
	if upgrade_name == "" {
		upgrade_name = UPGRADE_NAME
	}
	fmt.Println("Upgrade name:", upgrade_name)

	args := append([]string{
		CLI_BINARY_NAME,
		"tx", "gov", "submit-proposal", "software-upgrade",
		upgrade_name,
		"--title", "Upgrade Title",
		"--description", "Upgrade Description",
		"--upgrade-height", strconv.FormatInt(upgradeHeight, 10),
		"--upgrade-info", "Upgrade Info",
		"--from", OperatorAccounts[container],
	}, TX_PARAMS...)

	// ToDo: refactor
	if RUN_INSIDE_DOCKER {
		out, err = LocalnetExecExec(container, args...)
	} else {
		out, err = ExecDirectWithHome(container, args...)
	}

	if err != nil {
		return sdk.TxResponse{}, err
	}

	// Skip 'gas estimate: xxx' string, trim 'Successfully migrated key' string
	out = integrationhelpers.TrimImportedStdout(out)

	var resp sdk.TxResponse

	err = integrationhelpers.Codec.UnmarshalJSON([]byte(out), &resp)
	if err != nil {
		fmt.Println("JSON unmarshal error: output:", out)
		return sdk.TxResponse{}, err
	}
	return resp, nil
}

func DepositGov(container string, proposal_id string) (sdk.TxResponse, error) {
	var err error
	var out string

	fmt.Println("Depositing from", container)
	args := append([]string{
		CLI_BINARY_NAME,
		"tx", "gov", "deposit", proposal_id, DEPOSIT_AMOUNT,
		"--from", OperatorAccounts[container],
	}, TX_PARAMS...)

	// ToDo: refactor
	if RUN_INSIDE_DOCKER {
		out, err = LocalnetExecExec(container, args...)
	} else {
		out, err = ExecDirectWithHome(container, args...)
	}

	if err != nil {
		return sdk.TxResponse{}, err
	}

	// Skip 'gas estimate: xxx' string, trim 'Successfully migrated key' string
	out = integrationhelpers.TrimImportedStdout(out)

	var resp sdk.TxResponse

	err = integrationhelpers.Codec.UnmarshalJSON([]byte(out), &resp)
	if err != nil {
		return sdk.TxResponse{}, err
	}
	return resp, nil
}

func VoteProposal(container, id, option string) (sdk.TxResponse, error) {
	var err error
	var out string

	fmt.Println("Voting from", container)
	args := append([]string{
		CLI_BINARY_NAME,
		"tx", "gov", "vote", id, option,
		"--from", OperatorAccounts[container],
	}, TX_PARAMS...)

	// ToDo: refactor
	if RUN_INSIDE_DOCKER {
		out, err = LocalnetExecExec(container, args...)
	} else {
		out, err = ExecDirectWithHome(container, args...)
	}

	if err != nil {
		return sdk.TxResponse{}, err
	}

	// Skip 'gas estimate: xxx' string, trim 'Successfully migrated key' string
	out = integrationhelpers.TrimImportedStdout(out)

	var resp sdk.TxResponse

	err = integrationhelpers.Codec.UnmarshalJSON([]byte(out), &resp)
	if err != nil {
		return sdk.TxResponse{}, err
	}
	return resp, nil
}
