package cli

import (
	"encoding/json"
	"fmt"
	"time"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cheqd/cheqd-node/tests/integration/helpers"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
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
		CliBinaryName,
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
		CliBinaryName,
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

func QueryProposal(container, id string) (govtypesv1.Proposal, error) {
	fmt.Println("Querying proposal from", container)
	args := append([]string{
		CliBinaryName,
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

func QueryTxn(container, hash string) (sdk.TxResponse, error) {
	time.Sleep(2500 * time.Millisecond)
	res, err := Query(container, CliBinaryName, "tx", hash)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	var resp sdk.TxResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	return resp, nil
}

func GetProposalID(rawLog string) (string, error) {
	var logs []sdk.ABCIMessageLog
	err := json.Unmarshal([]byte(rawLog), &logs)
	if err != nil {
		return "", err
	}

	// Iterate over logs and their events
	for _, log := range logs {
		for _, event := range log.Events {
			// Look for the "submit_proposal" event type
			if event.Type == "submit_proposal" {
				for _, attr := range event.Attributes {
					// Look for the "proposal_id" attribute
					if attr.Key == "proposal_id" {
						return attr.Value, nil
					}
				}
			}
		}
	}

	return "", fmt.Errorf("proposal_id not found")
}
