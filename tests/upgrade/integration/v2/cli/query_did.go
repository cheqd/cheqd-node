package cli

import (
	integrationhelpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
	"github.com/cheqd/cheqd-node/x/did/types"
)

func QueryDid(did string, container string) (types.QueryDidDocResponse, error) {
	res, err := Query(container, CliBinaryName, "cheqd", "did-document", did)
	if err != nil {
		return types.QueryDidDocResponse{}, err
	}

	var resp types.QueryDidDocResponse
	err = integrationhelpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return types.QueryDidDocResponse{}, err
	}

	return resp, nil
}
