package cli

import (
	"github.com/cheqd/cheqd-node/tests/integration/helpers"
	"github.com/cheqd/cheqd-node/tests/integration/network"

	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
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

func QueryBalance(address, denom string) (banktypes.QueryBalanceResponse, error) {
	res, err := Query("bank", "balances", address)
	if err != nil {
		return banktypes.QueryBalanceResponse{}, err
	}

	var resp banktypes.QueryBalanceResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return banktypes.QueryBalanceResponse{}, err
	}

	return resp, nil
}

func QuerySupplyOf(denom string) (banktypes.QuerySupplyOfResponse, error) {
	res, err := Query("bank", "total", denom)
	if err != nil {
		return banktypes.QuerySupplyOfResponse{}, err
	}

	var resp banktypes.QuerySupplyOfResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return banktypes.QuerySupplyOfResponse{}, err
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

func QueryDidDoc(did string) (didtypes.QueryGetDidDocResponse, error) {
	res, err := Query("cheqd", "diddoc", did)
	if err != nil {
		return didtypes.QueryGetDidDocResponse{}, err
	}

	var resp didtypes.QueryGetDidDocResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return didtypes.QueryGetDidDocResponse{}, err
	}

	return resp, nil
}

func QueryResource(collectionId string, resourceId string) (resourcetypes.QueryGetResourceResponse, error) {
	res, err := Query("resource", "resource", collectionId, resourceId)
	if err != nil {
		return resourcetypes.QueryGetResourceResponse{}, err
	}

	var resp resourcetypes.QueryGetResourceResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return resourcetypes.QueryGetResourceResponse{}, err
	}

	return resp, nil
}

func QueryResourceMetadata(collectionId string, resourceId string) (resourcetypes.QueryGetResourceMetadataResponse, error) {
	res, err := Query("resource", "resource-metadata", collectionId, resourceId)
	if err != nil {
		return resourcetypes.QueryGetResourceMetadataResponse{}, err
	}

	var resp resourcetypes.QueryGetResourceMetadataResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return resourcetypes.QueryGetResourceMetadataResponse{}, err
	}

	return resp, nil
}

func QueryResourceCollection(collectionId string) (resourcetypes.QueryGetCollectionResourcesResponse, error) {
	res, err := Query("resource", "collection-resources", collectionId)
	if err != nil {
		return resourcetypes.QueryGetCollectionResourcesResponse{}, err
	}

	var resp resourcetypes.QueryGetCollectionResourcesResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return resourcetypes.QueryGetCollectionResourcesResponse{}, err
	}

	return resp, nil
}
