package cli

import (
	"fmt"

	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
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
	args = append(args, QueryParamsConst...)

	return LocalnetExecExec(container, args...)
}

func QueryModuleVersionMap(container string) (upgradetypes.QueryModuleVersionsResponse, error) {
	fmt.Println("Querying module version map from", container)
	args := append([]string{
		CLIBinaryName,
		"query", "upgrade", "module_versions",
	}, QueryParamsConst...)

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
		CLIBinaryName,
		"query", "params", "subspace", subspace, key,
	}, QueryParamsConst...)

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

func QueryResourceFeeParams(container, subspace, key string) (resourcetypes.FeeParams, error) {
	params, err := QueryParams(container, subspace, key)
	if err != nil {
		return resourcetypes.FeeParams{}, err
	}

	var feeParams resourcetypes.FeeParams

	err = MakeCodecWithExtendedRegistry().UnmarshalJSON([]byte(params.Value), &feeParams)
	if err != nil {
		return resourcetypes.FeeParams{}, err
	}

	return feeParams, nil
}

func QueryProposalLegacy(container, id string) (govtypesv1beta1.Proposal, error) {
	fmt.Println("Querying proposal from", container)
	args := append([]string{
		CLIBinaryName,
		"query", "gov", "proposal", id,
	}, QueryParamsConst...)

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

func QueryProposal(container, id string) (govtypesv1.Proposal, error) {
	fmt.Println("Querying proposal from", container)
	args := append([]string{
		CLIBinaryName,
		"query", "gov", "proposal", id,
	}, QueryParamsConst...)

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
