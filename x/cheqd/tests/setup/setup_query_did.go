package setup

import "github.com/cheqd/cheqd-node/x/cheqd/types"

func (s *TestSetup) QueryDid(did string) (*types.QueryGetDidResponse, error) {
	req := &types.QueryGetDidRequest{
		Id: did,
	}

	return s.QueryServer.Did(s.StdCtx, req)
}
