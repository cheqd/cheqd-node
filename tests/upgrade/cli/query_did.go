package cli

import (
	integrationhelpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
	didtypesv2 "github.com/cheqd/cheqd-node/x/did/types"
	didtypesv1 "github.com/cheqd/cheqd-node/x/did/types/v1"
)

func QueryDidLegacy(did string, container string) (didtypesv1.QueryGetDidResponse, error) {
	res, err := Query(container, CLI_BINARY_NAME, "cheqd", "did", did)
	if err != nil {
		return didtypesv1.QueryGetDidResponse{}, err
	}

	var resp didtypesv1.QueryGetDidResponse
	err = integrationhelpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return didtypesv1.QueryGetDidResponse{}, err
	}

	return resp, nil
}

func QueryDid(did string, container string) (didtypesv2.QueryGetDidDocResponse, error) {
	res, err := Query(container, CLI_BINARY_NAME, "cheqd", "did", did)
	if err != nil {
		return didtypesv2.QueryGetDidDocResponse{}, err
	}

	var resp didtypesv2.QueryGetDidDocResponse
	err = integrationhelpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return didtypesv2.QueryGetDidDocResponse{}, err
	}

	return resp, nil
}
