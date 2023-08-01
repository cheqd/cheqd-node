package setup

import (
	"github.com/cheqd/cheqd-node/x/resource"
	"github.com/cheqd/cheqd-node/x/resource/types"
)

func (s *TestSetup) StorePort() {
	s.ResourceKeeper.SetPort(s.SdkCtx, types.ResourcePortId)
}

func (s *TestSetup) StorePortWithGenesis() {
	genesis := types.DefaultGenesis()
	resource.InitGenesis(s.SdkCtx, s.ResourceKeeper, genesis)
}
