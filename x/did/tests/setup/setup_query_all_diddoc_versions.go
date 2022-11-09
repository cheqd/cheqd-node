package setup

import "github.com/cheqd/cheqd-node/x/did/types"

func (s *TestSetup) QueryAllDidDocVersions(did string) (*types.QueryGetAllDidDocVersionsResponse, error) {
	req := &types.QueryGetAllDidDocVersionsRequest{
		Id: did,
	}

	return s.QueryServer.AllDidDocVersions(s.StdCtx, req)
}
