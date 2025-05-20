package cli

import (
	"fmt"
	"time"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	integrationhelpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

func QueryDidFeeParams(container string) (didtypes.QueryParamsResponse, error) {
	res, err := Query(container, CliBinaryName, "cheqd", "params")
	if err != nil {
		return didtypes.QueryParamsResponse{}, err
	}

	var resp didtypes.QueryParamsResponse
	err = MakeCodecWithExtendedRegistry().UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return didtypes.QueryParamsResponse{}, err
	}

	return resp, nil
}

func QueryResourceFeeParams(container string) (resourcetypes.QueryParamsResponse, error) {
	res, err := Query(container, CliBinaryName, "resource", "params")
	if err != nil {
		return resourcetypes.QueryParamsResponse{}, err
	}

	var resp resourcetypes.QueryParamsResponse
	err = MakeCodecWithExtendedRegistry().UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return resourcetypes.QueryParamsResponse{}, err
	}

	return resp, nil
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

	// FIX: getting type instead of @type in messages struct when querying proposal via cli
	convertedJSON, err := convertProposalJSON(out)
	if err != nil {
		return govtypesv1.Proposal{}, err
	}

	var resp govtypesv1.QueryProposalResponse
	err = integrationhelpers.Codec.UnmarshalJSON(convertedJSON, &resp)
	if err != nil {
		return govtypesv1.Proposal{}, err
	}
	return *resp.Proposal, nil
}

func GetProposalID(events []abcitypes.Event) (string, error) {
	// Iterate over events
	for _, event := range events {
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
	return "", fmt.Errorf("proposal_id not found")
}

func QueryTxn(container, hash string) (sdk.TxResponse, error) {
	time.Sleep(2000 * time.Millisecond)
	res, err := Query(container, CliBinaryName, "tx", hash)
	if err != nil {
		fmt.Println("Error querying tx", res)
		return sdk.TxResponse{}, err
	}

	var resp sdk.TxResponse
	err = integrationhelpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		fmt.Println("Error unmarshalling tx", res)
		return sdk.TxResponse{}, err
	}

	return resp, nil
}

func QueryResource(collectionID string, resourceID string, container string) (resourcetypes.QueryResourceResponse, error) {
	res, err := Query(container, CliBinaryName, "resource", "specific-resource", collectionID, resourceID)
	if err != nil {
		return resourcetypes.QueryResourceResponse{}, err
	}

	var resp resourcetypes.QueryResourceResponse
	err = integrationhelpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return resourcetypes.QueryResourceResponse{}, err
	}

	return resp, nil
}

func QueryDid(did string, container string) (didtypes.QueryDidDocResponse, error) {
	res, err := Query(container, CliBinaryName, "cheqd", "did-document", did)
	if err != nil {
		return didtypes.QueryDidDocResponse{}, err
	}

	var resp didtypes.QueryDidDocResponse
	err = integrationhelpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return didtypes.QueryDidDocResponse{}, err
	}

	return resp, nil
}
