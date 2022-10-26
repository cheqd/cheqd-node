package cli

import (
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	integrationhelpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
)

func QueryResource(collectionId string, resourceId string, container string) (resourcetypes.QueryGetResourceResponse, error) {
	res, err := Query(container, CLI_BINARY_NAME, "resource", "resource", collectionId, resourceId)
	if err != nil {
		return resourcetypes.QueryGetResourceResponse{}, err
	}

	var resp resourcetypes.QueryGetResourceResponse
	err = integrationhelpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return resourcetypes.QueryGetResourceResponse{}, err
	}

	return resp, nil
}