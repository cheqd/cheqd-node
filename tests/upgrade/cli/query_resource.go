package cli

import (
	integrationhelpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
	resourcetypesv2 "github.com/cheqd/cheqd-node/x/resource/types"
	resourcetypesv1 "github.com/cheqd/cheqd-node/x/resource/types/v1"
)

func QueryResourceLegacy(collectionId string, resourceId string, container string) (resourcetypesv1.QueryGetResourceResponse, error) {
	res, err := Query(container, CLI_BINARY_NAME, "resource", "resource", collectionId, resourceId)
	if err != nil {
		return resourcetypesv1.QueryGetResourceResponse{}, err
	}

	var resp resourcetypesv1.QueryGetResourceResponse
	err = integrationhelpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return resourcetypesv1.QueryGetResourceResponse{}, err
	}

	return resp, nil
}

func QueryResource(collectionId string, resourceId string, container string) (resourcetypesv2.QueryGetResourceResponse, error) {
	res, err := Query(container, CLI_BINARY_NAME, "resource", "resource", collectionId, resourceId)
	if err != nil {
		return resourcetypesv2.QueryGetResourceResponse{}, err
	}

	var resp resourcetypesv2.QueryGetResourceResponse
	err = integrationhelpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return resourcetypesv2.QueryGetResourceResponse{}, err
	}

	return resp, nil
}
