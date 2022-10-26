package cli

import (
	cheqdtypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	integrationhelpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
)

func QueryDid(did string, container string) (cheqdtypes.QueryGetDidResponse, error) {
	res, err := Query(container, CLI_BINARY_NAME, "cheqd", "did", did)
	if err != nil {
		return cheqdtypes.QueryGetDidResponse{}, err
	}

	var resp cheqdtypes.QueryGetDidResponse
	err = integrationhelpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return cheqdtypes.QueryGetDidResponse{}, err
	}

	return resp, nil
}
