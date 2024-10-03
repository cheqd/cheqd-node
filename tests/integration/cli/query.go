package cli

import (
	"fmt"

	"github.com/cheqd/cheqd-node/tests/integration/helpers"
	"github.com/cheqd/cheqd-node/tests/integration/network"

	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
)

var CLIQueryParams = []string{
	"--chain-id", network.ChainID,
	"--output", OutputFormat,
}

func Query(module, query string, queryArgs ...string) (string, error) {
	args := []string{"query", module, query}

	// Common params
	args = append(args, CLIQueryParams...)

	// Other args
	args = append(args, queryArgs...)

	return Exec(args...)
}

func QueryBalance(address, denom string) (sdk.Coin, error) {
	res, err := Query("bank", "balances", address, "--denom", denom)
	if err != nil {
		return sdk.Coin{}, err
	}

	var resp sdk.Coin
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return sdk.Coin{}, err
	}

	return resp, nil
}

func QueryParams(subspace, key string) (paramproposal.ParamChange, error) {
	res, err := Query("params", "subspace", subspace, key)
	if err != nil {
		return paramproposal.ParamChange{}, err
	}

	var resp paramproposal.ParamChange
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return paramproposal.ParamChange{}, err
	}

	return resp, nil
}

func QueryDidDoc(did string) (didtypes.QueryDidDocResponse, error) {
	res, err := Query("cheqd", "did-document", did)
	if err != nil {
		return didtypes.QueryDidDocResponse{}, err
	}

	var resp didtypes.QueryDidDocResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return didtypes.QueryDidDocResponse{}, err
	}

	return resp, nil
}

func QueryResource(collectionID string, resourceID string) (resourcetypes.QueryResourceResponse, error) {
	res, err := Query("resource", "specific-resource", collectionID, resourceID)
	if err != nil {
		return resourcetypes.QueryResourceResponse{}, err
	}

	var resp resourcetypes.QueryResourceResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return resourcetypes.QueryResourceResponse{}, err
	}

	return resp, nil
}

func QueryResourceMetadata(collectionID string, resourceID string) (resourcetypes.QueryResourceMetadataResponse, error) {
	res, err := Query("resource", "metadata", collectionID, resourceID)
	if err != nil {
		return resourcetypes.QueryResourceMetadataResponse{}, err
	}

	var resp resourcetypes.QueryResourceMetadataResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return resourcetypes.QueryResourceMetadataResponse{}, err
	}

	return resp, nil
}

func QueryResourceCollection(collectionID string) (resourcetypes.QueryCollectionResourcesResponse, error) {
	res, err := Query("resource", "collection-metadata", collectionID)
	if err != nil {
		return resourcetypes.QueryCollectionResourcesResponse{}, err
	}

	var resp resourcetypes.QueryCollectionResourcesResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return resourcetypes.QueryCollectionResourcesResponse{}, err
	}

	return resp, nil
}

func QueryTxn(hash string) (sdk.TxResponse, error) {
	res, err := Query("tx", hash)
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

func QueryKeys(container string) (string, error) {
	fmt.Println("Querying proposal from", container)
	args := append([]string{
		CliBinaryName,
		"keys", "list",
	}, KeyringParams...)

	out, err := LocalnetExecExec(container, args...)

	return out, err
}

func FindFile(container, filename string) (string, error) {
	args := []string{
		"find", "/home", "-name", filename,
	}
	out, err := LocalnetExecExec(container, args...)
	return out, err
}
