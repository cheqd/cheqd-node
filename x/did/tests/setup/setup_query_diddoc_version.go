package setup

import "github.com/cheqd/cheqd-node/x/did/types"

func (s *TestSetup) QueryDidDocVersion(did, version string) (*types.QueryDidDocVersionResponse, error) {
	req := &types.QueryDidDocVersionRequest{
		Id:      did,
		Version: version,
	}

	return s.QueryServer.DidDocVersion(s.StdCtx, req)
}
