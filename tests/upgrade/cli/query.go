package cli

import (
	"fmt"

	integrationhelpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
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

func QueryUpgradeProposal(container string) (govtypesv1beta1.Proposal, error) {
	args := append([]string{
		CLI_BINARY_NAME,
		"query", "gov", "proposal", "1",
	}, QUERY_PARAMS...)

	out, err := LocalnetExecExec(container, args...)
	if err != nil {
		return govtypesv1beta1.Proposal{}, err
	}

	fmt.Println("QueryUpgradeProposal", out)

	fmt.Println("Registry Implementations", integrationhelpers.Registry.ListImplementations("cosmos.gov.v1beta1.Content"))

	var resp govtypesv1beta1.Proposal

	err = integrationhelpers.Codec.UnmarshalJSON([]byte(out), &resp)
	if err != nil {
		return govtypesv1beta1.Proposal{}, err
	}
	return resp, nil
}
