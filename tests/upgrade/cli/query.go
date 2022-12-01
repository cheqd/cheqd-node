package cli

import (
	"fmt"

	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
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

func QueryModuleVersionMap(container string) (upgradetypes.QueryModuleVersionsResponse, error) {
	fmt.Println("Querying module version map from", container)
	args := append([]string{
		CLI_BINARY_NAME,
		"query", "upgrade", "module_versions",
	}, QUERY_PARAMS...)

	out, err := LocalnetExecExec(container, args...)
	if err != nil {
		return upgradetypes.QueryModuleVersionsResponse{}, err
	}

	fmt.Println("Module version map", out)

	var resp upgradetypes.QueryModuleVersionsResponse

	err = MakeCodecWithExtendedRegistry().UnmarshalJSON([]byte(out), &resp)
	if err != nil {
		return upgradetypes.QueryModuleVersionsResponse{}, err
	}

	return resp, nil
}

func QueryParams(container, subspace, key string) (paramproposal.ParamChange, error) {
	fmt.Println("Querying params from", container)
	args := append([]string{
		CLI_BINARY_NAME,
		"query", "params", "subspace", subspace, key,
	}, QUERY_PARAMS...)

	out, err := LocalnetExecExec(container, args...)
	if err != nil {
		return paramproposal.ParamChange{}, err
	}

	fmt.Println("Params", out)

	var resp paramproposal.ParamChange

	err = MakeCodecWithExtendedRegistry().UnmarshalJSON([]byte(out), &resp)
	if err != nil {
		return paramproposal.ParamChange{}, err
	}

	return resp, nil
}

func QueryDidFeeParams(container, subspace, key string) (didtypes.FeeParams, error) {
	params, err := QueryParams(container, subspace, key)
	if err != nil {
		return didtypes.FeeParams{}, err
	}

	var feeParams didtypes.FeeParams

	err = MakeCodecWithExtendedRegistry().UnmarshalJSON([]byte(params.Value), &feeParams)
	if err != nil {
		return didtypes.FeeParams{}, err
	}

	return feeParams, nil
}

func QueryProposalLegacy(container, id string) (govtypesv1beta1.Proposal, error) {
	fmt.Println("Querying proposal from", container)
	args := append([]string{
		CLI_BINARY_NAME,
		"query", "gov", "proposal", id,
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
