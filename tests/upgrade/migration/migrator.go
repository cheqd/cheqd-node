package migration

import (
	cheqdtestssetup "github.com/cheqd/cheqd-node/x/cheqd/tests/setup"
	cheqdtypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	resourcetestssetup "github.com/cheqd/cheqd-node/x/resource/tests/setup"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MigrationScenario struct {
	name     string
	setup    func()
	expected interface{}
	handler  func(ctx sdk.Context) error
	validate func(actual interface{}) error
}

func NewMigrationScenario(name string, setup func(), expected interface{}, handler func(ctx sdk.Context) error, validate func(actual interface{}) error) MigrationScenario {
	return MigrationScenario{
		name:     name,
		setup:    setup,
		expected: expected,
		handler:  handler,
		validate: validate,
	}
}

func (m MigrationScenario) Name() string {
	return m.name
}

func (m MigrationScenario) Setup() {
	m.setup()
}

func (m MigrationScenario) Expected() interface{} {
	return m.expected
}

func (m MigrationScenario) Handler() func(ctx sdk.Context) error {
	return m.handler
}

func (m MigrationScenario) Validate(actual interface{}) error {
	return m.validate(actual)
}

type CheqdMigrationScenario struct {
	name     string
	setup    func() cheqdtestssetup.TestSetup
	expected cheqdtypes.Did
	handler  func(ctx sdk.Context) error
	validate func(actual cheqdtypes.Did) error
}

func NewCheqdMigrationScenario(name string, setup func() cheqdtestssetup.TestSetup, expected cheqdtypes.Did, handler func(ctx sdk.Context) error, validate func(actual cheqdtypes.Did) error) CheqdMigrationScenario {
	return CheqdMigrationScenario{
		name:     name,
		setup:    setup,
		expected: expected,
		handler:  handler,
		validate: validate,
	}
}

func (m CheqdMigrationScenario) Name() string {
	return m.name
}

func (m CheqdMigrationScenario) Setup() cheqdtestssetup.TestSetup {
	return m.setup()
}

func (m CheqdMigrationScenario) Expected() cheqdtypes.Did {
	return m.expected
}

func (m CheqdMigrationScenario) Handler() func(ctx sdk.Context) error {
	return m.handler
}

func (m CheqdMigrationScenario) Validate(actual cheqdtypes.Did) error {
	return m.validate(actual)
}

type ResourceMigrationScenario struct {
	name     string
	setup    func() resourcetestssetup.TestSetup
	expected resourcetypes.Resource
	handler  func(ctx sdk.Context) error
	validate func(actual resourcetypes.Resource) error
}

func NewResourceMigrationScenario(name string, setup func() resourcetestssetup.TestSetup, expected resourcetypes.Resource, handler func(ctx sdk.Context) error, validate func(actual resourcetypes.Resource) error) ResourceMigrationScenario {
	return ResourceMigrationScenario{
		name:     name,
		setup:    setup,
		expected: expected,
		handler:  handler,
		validate: validate,
	}
}

func (m ResourceMigrationScenario) Name() string {
	return m.name
}

func (m ResourceMigrationScenario) Setup() resourcetestssetup.TestSetup {
	return m.setup()
}

func (m ResourceMigrationScenario) Expected() resourcetypes.Resource {
	return m.expected
}

func (m ResourceMigrationScenario) Handler() func(ctx sdk.Context) error {
	return m.handler
}

func (m ResourceMigrationScenario) Validate(actual resourcetypes.Resource) error {
	return m.validate(actual)
}

type Migrator interface {
	Migrate(ctx sdk.Context) error
}

type CheqdMigrator struct {
	migrations []CheqdMigrationScenario
}

func NewCheqdMigrator(migrations []CheqdMigrationScenario) CheqdMigrator {
	return CheqdMigrator{
		migrations: migrations,
	}
}

func (m CheqdMigrator) Migrate(ctx sdk.Context) error {
	for _, migration := range m.migrations {
		setup := migration.Setup()
		err := migration.Handler()(setup.SdkCtx)
		if err != nil {
			return err
		}
		err = migration.Validate(migration.Expected())
		if err != nil {
			return err
		}
	}
	return nil
}

type ResourceMigrator struct {
	migrations []ResourceMigrationScenario
}

func NewResourceMigrator(migrations []ResourceMigrationScenario) ResourceMigrator {
	return ResourceMigrator{
		migrations: migrations,
	}
}

func (m ResourceMigrator) Migrate() error {
	for _, migration := range m.migrations {
		setup := migration.Setup()
		err := migration.Handler()(setup.SdkCtx)
		if err != nil {
			return err
		}
		err = migration.Validate(migration.Expected())
		if err != nil {
			return err
		}
	}
	return nil
}
