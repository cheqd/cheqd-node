package scenarios

import (
	appmigrations "github.com/cheqd/cheqd-node/app/migrations"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MigrationScenario struct {
	name    string
	handler func(ctx sdk.Context, migrationCtx appmigrations.MigrationContext) error
}

func NewMigrationScenario(
	name string,
	handler func(ctx sdk.Context, migrationCtx appmigrations.MigrationContext) error,
) MigrationScenario {
	return MigrationScenario{
		name:    name,
		handler: handler,
	}
}

func (m MigrationScenario) Name() string {
	return m.name
}

func (m MigrationScenario) Handler() func(ctx sdk.Context, migrationCtx appmigrations.MigrationContext) error {
	return m.handler
}