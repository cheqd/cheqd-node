package migration

import (
	didkeeper "github.com/cheqd/cheqd-node/x/did/keeper"
	didtestssetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	resourcekeeper "github.com/cheqd/cheqd-node/x/resource/keeper"
	resourcetestssetup "github.com/cheqd/cheqd-node/x/resource/tests/setup"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MigrationScenario struct {
	name     string
	setup    func()
	existing interface{}
	expected interface{}
	handler  func(ctx sdk.Context, didKeeper didkeeper.Keeper) error
	validate func(actual interface{}) error
}

func NewMigrationScenario(name string, setup func(), existing interface{}, expected interface{}, handler func(ctx sdk.Context, didKeeper didkeeper.Keeper) error, validate func(actual interface{}) error) MigrationScenario {
	return MigrationScenario{
		name:     name,
		setup:    setup,
		existing: existing,
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

func (m MigrationScenario) Existing() interface{} {
	return m.existing
}

func (m MigrationScenario) Expected() interface{} {
	return m.expected
}

func (m MigrationScenario) Handler() func(ctx sdk.Context, didKeeper didkeeper.Keeper) error {
	return m.handler
}

func (m MigrationScenario) Validate(actual interface{}) error {
	return m.validate(actual)
}

type CheqdMigrationScenario struct {
	name     string
	setup    func() didtestssetup.TestSetup
	existing didtestssetup.MinimalDidDocInfo
	expected didtypes.DidDoc
	handler  func(ctx sdk.Context, didKeeper didkeeper.Keeper) error
	validate func(actual didtypes.DidDoc) error
}

func NewCheqdMigrationScenario(name string, setup func() didtestssetup.TestSetup, existing didtestssetup.MinimalDidDocInfo, expected didtypes.DidDoc, handler func(ctx sdk.Context, didKeeper didkeeper.Keeper) error, validate func(actual didtypes.DidDoc) error) CheqdMigrationScenario {
	return CheqdMigrationScenario{
		name:     name,
		setup:    setup,
		existing: existing,
		expected: expected,
		handler:  handler,
		validate: validate,
	}
}

func (m CheqdMigrationScenario) Name() string {
	return m.name
}

func (m CheqdMigrationScenario) Setup() didtestssetup.TestSetup {
	return m.setup()
}

func (m CheqdMigrationScenario) Existing() didtestssetup.MinimalDidDocInfo {
	return m.existing
}

func (m CheqdMigrationScenario) Expected() didtypes.DidDoc {
	return m.expected
}

func (m CheqdMigrationScenario) Handler() func(ctx sdk.Context, didKeeper didkeeper.Keeper) error {
	return m.handler
}

func (m CheqdMigrationScenario) Validate(actual didtypes.DidDoc) error {
	return m.validate(actual)
}

type ResourceMigrationScenario struct {
	name     string
	setup    func() resourcetestssetup.TestSetup
	existing resourcetypes.MsgCreateResourcePayload
	didInfo  didtestssetup.MinimalDidDocInfoV1
	expected resourcetypes.Metadata
	handler  func(ctx sdk.Context, didKeeper didkeeper.Keeper, resourceKeeper resourcekeeper.Keeper) error
	validate func(actual resourcetypes.Metadata) error
}

func NewResourceMigrationScenario(name string, setup func() resourcetestssetup.TestSetup, existing resourcetypes.MsgCreateResourcePayload, didInfo didtestssetup.MinimalDidDocInfoV1, expected resourcetypes.Metadata, handler func(ctx sdk.Context, didKeeper didkeeper.Keeper, resourceKeeper resourcekeeper.Keeper) error, validate func(actual resourcetypes.Metadata) error) ResourceMigrationScenario {
	return ResourceMigrationScenario{
		name:     name,
		setup:    setup,
		existing: existing,
		didInfo:  didInfo,
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

func (m ResourceMigrationScenario) Existing() resourcetypes.MsgCreateResourcePayload {
	return m.existing
}

func (m ResourceMigrationScenario) DidInfo() didtestssetup.MinimalDidDocInfoV1 {
	return m.didInfo
}

func (m ResourceMigrationScenario) Expected() resourcetypes.Metadata {
	return m.expected
}

func (m ResourceMigrationScenario) Handler() func(ctx sdk.Context, didKeeper didkeeper.Keeper, resourceKeeper resourcekeeper.Keeper) error {
	return m.handler
}

func (m ResourceMigrationScenario) Validate(actual resourcetypes.Metadata) error {
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
		_, err := setup.CreateDid(migration.existing.Msg, []didtestssetup.SignInput{migration.existing.SignInput})
		if err != nil {
			return err
		}
		err = migration.Handler()(setup.SdkCtx, setup.Keeper)
		if err != nil {
			return err
		}
		data, err := setup.Keeper.GetDidDoc(&setup.SdkCtx, migration.existing.Msg.Id)
		if err != nil {
			return err
		}
		actual := data.DidDoc
		if err != nil {
			return err
		}
		err = migration.Validate(*actual)
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
		_, err := setup.CreateDidV1(migration.didInfo.Msg, []didtestssetup.SignInput{migration.didInfo.SignInput})
		if err != nil {
			return err
		}
		_, err = setup.CreateResource(&migration.existing, []didtestssetup.SignInput{migration.didInfo.SignInput})
		if err != nil {
			return err
		}
		err = migration.Handler()(setup.SdkCtx, setup.Keeper, setup.ResourceKeeper)
		if err != nil {
			return err
		}
		actual, err := setup.QueryResource(migration.existing.CollectionId, migration.existing.Id)
		if err != nil {
			return err
		}
		err = migration.Validate(*actual.Resource.Metadata)
		if err != nil {
			return err
		}
	}
	return nil
}
