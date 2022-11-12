package cli

import (
	"fmt"

	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
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

func QueryUpgradeProposalLegacy(container string) (govtypesv1beta1.Proposal, error) {
	fmt.Println("Querying upgrade proposal from", container)
	args := append([]string{
		CLI_BINARY_NAME,
		"query", "gov", "proposal", "1",
	}, QUERY_PARAMS...)

	out, err := LocalnetExecExec(container, args...)
	if err != nil {
		return govtypesv1beta1.Proposal{}, err
	}

	fmt.Println("Proposal", out)

	var resp govtypesv1beta1.Proposal

	err = MakeCodecWithExtendedRegistry().UnmarshalJSON([]byte(out), &resp)
	if err != nil {
		return govtypesv1beta1.Proposal{}, err
	}
	return resp, nil
}

func QueryUpgradeProposal(container string) (govtypesv1.Proposal, error) {
	fmt.Println("Querying upgrade proposal from", container)
	args := append([]string{
		CLI_BINARY_NAME,
		"query", "gov", "proposal", "1",
	}, QUERY_PARAMS...)

	out, err := LocalnetExecExec(container, args...)
	if err != nil {
		return govtypesv1.Proposal{}, err
	}

	fmt.Println("Proposal", out)

	var resp govtypesv1.Proposal

	err = MakeCodecWithExtendedRegistry().UnmarshalJSON([]byte(out), &resp)
	if err != nil {
		return govtypesv1.Proposal{}, err
	}
	return resp, nil
}
