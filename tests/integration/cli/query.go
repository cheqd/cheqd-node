package cli

import (
	"github.com/cheqd/cheqd-node/tests/integration/helpers"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
)

var CLI_QUERY_PARAMS = []string{
	"--chain-id",
	CHAIN_ID,
	"--output",
	OUTPUT_FORMAT,
}

func Query(module, query string, queryArgs ...string) (string, error) {
	args := []string{"query", module, query}

	// Common params
	args = append(args, CLI_QUERY_PARAMS...)

	// Other args
	args = append(args, queryArgs...)

	return Exec(args...)
}

func QueryDid(did string) (types.QueryGetDidResponse, error) {
	res, err := Query("cheqd", "did", did)
	if err != nil {
		return types.QueryGetDidResponse{}, err
	}

	var resp types.QueryGetDidResponse
	err = helpers.Codec.UnmarshalJSON([]byte(res), &resp)
	if err != nil {
		return types.QueryGetDidResponse{}, err
	}

	return resp, nil
}
