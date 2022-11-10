package cli

import (
	integrationhelpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
)

func QueryDid(did string, container string) (didtypes.QueryGetDidDocResponse, error) {
	res, err := Query(container, CLI_BINARY_NAME, "cheqd", "did", did)
	if err != nil {
		return didtypes.QueryGetDidDocResponse{}, err
	}

	var resp didtypes.QueryGetDidDocResponse
	err = integrationhelpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return didtypes.QueryGetDidDocResponse{}, err
	}

	return resp, nil
}
