package cli

import (
	integrationhelpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func Query(container string, binary string, module, query string, queryArgs ...string) (string, error) {
	args := []string{
		binary,
		"query",
		module,
		query,
	}

	args = append(args, QUERY_PARAMS)
	args = append(args, queryArgs...)

	return LocalnetExecExec(container, args...)
}

func QueryUpgradeProposal(container string) (govtypes.QueryProposalResponse, error) {
	args := append([]string{
		CLI_BINARY_NAME,
		"query", "gov", "proposal", "1",
	}, QUERY_PARAMS)

	out, err := LocalnetExecExec(container, args...)
	if err != nil {
		return govtypes.QueryProposalResponse{}, err
	}

	var resp govtypes.QueryProposalResponse

	err = integrationhelpers.Codec.UnmarshalJSON([]byte(out), &resp)
	if err != nil {
		return govtypes.QueryProposalResponse{}, err
	}
	return resp, nil
}
