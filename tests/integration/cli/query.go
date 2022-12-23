package cli

import (
	"github.com/cheqd/cheqd-node/tests/integration/helpers"
	"github.com/cheqd/cheqd-node/tests/integration/network"

	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
)

var CLI_QUERY_PARAMS = []string{
	"--chain-id", network.CHAIN_ID,
	"--output", OUTPUT_FORMAT,
}

func Query(module, query string, queryArgs ...string) (string, error) {
	args := []string{"query", module, query}

	// Common params
	args = append(args, CLI_QUERY_PARAMS...)

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

func QueryResource(collectionId string, resourceId string) (resourcetypes.QueryResourceResponse, error) {
	res, err := Query("resource", "resource", collectionId, resourceId)
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

func QueryResourceMetadata(collectionId string, resourceId string) (resourcetypes.QueryResourceMetadataResponse, error) {
	res, err := Query("resource", "resource-metadata", collectionId, resourceId)
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

func QueryResourceCollection(collectionId string) (resourcetypes.QueryCollectionResourcesResponse, error) {
	res, err := Query("resource", "collection-resources", collectionId)
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
