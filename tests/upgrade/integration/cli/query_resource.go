package cli

import (
	integrationhelpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
	resourcetypesv2 "github.com/cheqd/cheqd-node/x/resource/types"
	resourcetypesv1 "github.com/cheqd/cheqd-node/x/resource/types/v1"
)

func QueryResourceLegacy(collectionID string, resourceID string, container string) (resourcetypesv1.QueryResourceResponse, error) {
	res, err := Query(container, CLIBinaryName, "resource", "resource", collectionID, resourceID)
	if err != nil {
		return resourcetypesv1.QueryResourceResponse{}, err
	}

	var resp resourcetypesv1.QueryResourceResponse
	err = integrationhelpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return resourcetypesv1.QueryResourceResponse{}, err
	}

	return resp, nil
}

func QueryResource(collectionID string, resourceID string, container string) (resourcetypesv2.QueryResourceResponse, error) {
	res, err := Query(container, CLIBinaryName, "resource", "resource", collectionID, resourceID)
	if err != nil {
		return resourcetypesv2.QueryResourceResponse{}, err
	}

	var resp resourcetypesv2.QueryResourceResponse
	err = integrationhelpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return resourcetypesv2.QueryResourceResponse{}, err
	}

	return resp, nil
}
