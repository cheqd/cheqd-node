package setup

import "github.com/cheqd/cheqd-node/x/did/types"

func (s *TestSetup) QueryAllDidDocVersionsMetadata(did string) (*types.QueryGetAllDidDocVersionsMetadataResponse, error) {
	req := &types.QueryGetAllDidDocVersionsMetadataRequest{
		Id: did,
	}

	return s.QueryServer.AllDidDocVersionsMetadata(s.StdCtx, req)
}
