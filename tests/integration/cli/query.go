package cli

import (
	"github.com/cheqd/cheqd-node/tests/integration/helpers"
	"github.com/cheqd/cheqd-node/tests/integration/network"

	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
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
