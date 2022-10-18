package migrations

import (
	cheqdkeeper "github.com/cheqd/cheqd-node/x/cheqd/keeper"
	resourcekeeper "github.com/cheqd/cheqd-node/x/resource/keeper"
)

type Migrator struct{}

type CheqdMigrator struct {
	Migrator
	cheqdKeeper cheqdkeeper.Keeper
}

func NewCheqdMigrator(cheqdKeeper cheqdkeeper.Keeper) CheqdMigrator {
	return CheqdMigrator{cheqdKeeper: cheqdKeeper}
}

type ResourceMigrator struct {
	CheqdMigrator
	resourceKeeper resourcekeeper.Keeper
}

func NewResourceMigrator(cheqdKeeper cheqdkeeper.Keeper, resourceKeeper resourcekeeper.Keeper) ResourceMigrator {
	return ResourceMigrator{CheqdMigrator: NewCheqdMigrator(cheqdKeeper), resourceKeeper: resourceKeeper}
}
