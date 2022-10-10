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

func QueryDid(did string) (cheqd_types.QueryGetDidResponse, error) {
	res, err := Query("cheqd", "did", did)
	if err != nil {
		return cheqd_types.QueryGetDidResponse{}, err
	}

	var resp cheqd_types.QueryGetDidResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return cheqd_types.QueryGetDidResponse{}, err
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
