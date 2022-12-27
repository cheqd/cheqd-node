package cli

import (
	integrationhelpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
	didtypesv2 "github.com/cheqd/cheqd-node/x/did/types"
	didtypesv1 "github.com/cheqd/cheqd-node/x/did/types/v1"
)

func QueryDidLegacy(did string, container string) (didtypesv1.QueryGetDidResponse, error) {
	res, err := Query(container, CLIBinaryName, "cheqd", "did", did)
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

func QueryDid(did string, container string) (didtypesv2.QueryDidDocResponse, error) {
	res, err := Query(container, CLIBinaryName, "cheqd", "did-document", did)
	if err != nil {
		return didtypesv2.QueryDidDocResponse{}, err
	}

	var resp didtypesv2.QueryDidDocResponse
	err = integrationhelpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return didtypesv2.QueryDidDocResponse{}, err
	}

	return resp, nil
}
