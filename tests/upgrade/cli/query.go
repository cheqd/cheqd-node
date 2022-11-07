package cli

import (
	integrationhelpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

func Query(container string, binary string, module, query string, queryArgs ...string) (string, error) {
	args := []string{
		binary,
		"query",
		module,
		query,
	}

	args = append(args, queryArgs...)
	args = append(args, QUERY_PARAMS...)

	return LocalnetExecExec(container, args...)
}

func QueryUpgradeProposal(container string) (govtypesv1.QueryProposalResponse, error) {
	args := append([]string{
		CLI_BINARY_NAME,
		"query", "gov", "proposal", "1",
	}, QUERY_PARAMS...)

	out, err := LocalnetExecExec(container, args...)
	if err != nil {
		return govtypesv1.QueryProposalResponse{}, err
	}

	var resp govtypesv1.QueryProposalResponse

	err = integrationhelpers.Codec.UnmarshalJSON([]byte(out), &resp)
	if err != nil {
		return govtypesv1.QueryProposalResponse{}, err
	}
	return resp, nil
}
