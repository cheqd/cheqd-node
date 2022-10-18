package app

import (
	cheqdapp "github.com/cheqd/cheqd-node/app"
	sdk "github.com/cosmos/cosmos-sdk/types"
	module "github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func Upgrade0_6to1_0(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
	return cheqdapp.App
}
