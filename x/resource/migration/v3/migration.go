package v3

import (
	"github.com/cheqd/cheqd-node/x/resource/keeper"
	"github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Migrator struct {
	k keeper.Keeper
}

func NewMigrator(k keeper.Keeper) Migrator {
	return Migrator{k: k}
}

func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	m.k.SetPort(ctx, types.ResourcePortId)
	m.k.BindPort(ctx, types.ResourcePortId)
	return nil
}
