package cli

import (
	"github.com/cheqd/cheqd-node/tests/integration/helpers"
	"github.com/cheqd/cheqd-node/tests/integration/network"

	cheqd_types "github.com/cheqd/cheqd-node/x/cheqd/types"
	resource_types "github.com/cheqd/cheqd-node/x/resource/types"
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

func QueryDidDoc(did string) (cheqd_types.QueryGetDidDocResponse, error) {
	res, err := Query("cheqd", "diddoc", did)
	if err != nil {
		return cheqd_types.QueryGetDidDocResponse{}, err
	}

	var resp cheqd_types.QueryGetDidDocResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return cheqd_types.QueryGetDidDocResponse{}, err
	}

	return resp, nil
}

func QueryResource(collectionId string, resourceId string) (resource_types.QueryGetResourceResponse, error) {
	res, err := Query("resource", "resource", collectionId, resourceId)
	if err != nil {
		return resource_types.QueryGetResourceResponse{}, err
	}

	var resp resource_types.QueryGetResourceResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return resource_types.QueryGetResourceResponse{}, err
	}

	return resp, nil
}

func QueryResourceMetadata(collectionId string, resourceId string) (resource_types.QueryGetResourceMetadataResponse, error) {
	res, err := Query("resource", "resource-metadata", collectionId, resourceId)
	if err != nil {
		return resource_types.QueryGetResourceMetadataResponse{}, err
	}

	var resp resource_types.QueryGetResourceMetadataResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return resource_types.QueryGetResourceMetadataResponse{}, err
	}

	return resp, nil
}

func QueryResourceCollection(collectionId string) (resource_types.QueryGetCollectionResourcesResponse, error) {
	res, err := Query("resource", "collection-resources", collectionId)
	if err != nil {
		return resource_types.QueryGetCollectionResourcesResponse{}, err
	}

	var resp resource_types.QueryGetCollectionResourcesResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return resource_types.QueryGetCollectionResourcesResponse{}, err
	}

	return resp, nil
}
