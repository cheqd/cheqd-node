package migration

import (
	cheqdtypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Migration struct {
	name     string
	expected interface{}
	handler  func(ctx sdk.Context) error
}

func NewMigration(name string, expected interface{}, handler func(ctx sdk.Context) error) Migration {
	return Migration{
		name:     name,
		expected: expected,
		handler:  handler,
	}
}

func (m Migration) Name() string {
	return m.name
}

func (m Migration) Expected() interface{} {
	return m.expected
}

func (m Migration) Handler() func(ctx sdk.Context) error {
	return m.handler
}

type CheqdMigration struct {
	name     string
	expected cheqdtypes.MsgCreateDidPayload
	handler  func(ctx sdk.Context) error
}

func NewCheqdMigration(name string, expected cheqdtypes.MsgCreateDidPayload, handler func(ctx sdk.Context) error) CheqdMigration {
	return CheqdMigration{
		name:     name,
		expected: expected,
		handler:  handler,
	}
}

func (m CheqdMigration) Name() string {
	return m.name
}

func (m CheqdMigration) Expected() cheqdtypes.MsgCreateDidPayload {
	return m.expected
}

func (m CheqdMigration) Handler() func(ctx sdk.Context) error {
	return m.handler
}

type ResourceMigration struct {
	name     string
	expected resourcetypes.MsgCreateResourcePayload
	handler  func(ctx sdk.Context) error
}

func NewResourceMigration(name string, expected resourcetypes.MsgCreateResourcePayload, handler func(ctx sdk.Context) error) ResourceMigration {
	return ResourceMigration{
		name:     name,
		expected: expected,
		handler:  handler,
	}
}

func (m ResourceMigration) Handler() func(ctx sdk.Context) error {
	return m.handler
}

func (m ResourceMigration) Name() string {
	return m.name
}

func (m ResourceMigration) Expected() resourcetypes.MsgCreateResourcePayload {
	return m.expected
}

type Migrator interface {
	Migrate(ctx sdk.Context) error
}

type CheqdMigrator struct {
	migrations []CheqdMigration
}

func NewCheqdMigrator(migrations []CheqdMigration) CheqdMigrator {
	return CheqdMigrator{
		migrations: migrations,
	}
}

func (m CheqdMigrator) Migrate(ctx sdk.Context) error {
	for _, migration := range m.migrations {
		err := migration.Handler()(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

type ResourceMigrator struct {
	migrations []ResourceMigration
}

func NewResourceMigrator(migrations []ResourceMigration) ResourceMigrator {
	return ResourceMigrator{
		migrations: migrations,
	}
}

func (m ResourceMigrator) Migrate(ctx sdk.Context) error {
	for _, migration := range m.migrations {
		err := migration.Handler()(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
