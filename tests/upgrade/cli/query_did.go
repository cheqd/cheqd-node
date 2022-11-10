package cli

import (
	integrationhelpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
	didtypes "github.com/cheqd/cheqd-node/x/did/types/v1"
)

func QueryDid(did string, container string) (didtypes.QueryGetDidResponse, error) {
	res, err := Query(container, CLI_BINARY_NAME, "cheqd", "did", did)
	if err != nil {
		return didtypes.QueryGetDidResponse{}, err
	}

	var resp didtypes.QueryGetDidResponse
	err = integrationhelpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return didtypes.QueryGetDidResponse{}, err
	}

	return resp, nil
}
