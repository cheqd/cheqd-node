package migration

import (
	appmigrations "github.com/cheqd/cheqd-node/app/migrations"
	migrationsetup "github.com/cheqd/cheqd-node/tests/upgrade/migration/setup"
	didtestssetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	resourcetypesv1 "github.com/cheqd/cheqd-node/x/resource/types/v1"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MigrationScenario struct {
	name     string
	setup    func()
	existing interface{}
	expected interface{}
	handler  func(ctx sdk.Context, migrationCtx appmigrations.MigrationContext) error
	validate func(actual interface{}) error
}

func NewMigrationScenario(name string, setup func(), existing interface{}, expected interface{}, handler func(ctx sdk.Context, migrationCtx appmigrations.MigrationContext) error, validate func(actual interface{}) error) MigrationScenario {
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

func (m MigrationScenario) Handler() func(ctx sdk.Context, migrationCtx appmigrations.MigrationContext) error {
	return m.handler
}

func (m MigrationScenario) Validate(actual interface{}) error {
	return m.validate(actual)
}

type DidMigrationScenario struct {
	name     string
	setup    func() migrationsetup.ExtendedTestSetup
	existing migrationsetup.MinimalDidDocInfoV1
	expected didtypes.DidDoc
	handler  func(ctx sdk.Context, migrationCtx appmigrations.MigrationContext) error
	validate func(actual didtypes.DidDoc) error
}

func NewDidMigrationScenario(name string, setup func() migrationsetup.ExtendedTestSetup, existing migrationsetup.MinimalDidDocInfoV1, expected didtypes.DidDoc, handler func(ctx sdk.Context, migrationCtx appmigrations.MigrationContext) error, validate func(actual didtypes.DidDoc) error) DidMigrationScenario {
	return DidMigrationScenario{
		name:     name,
		setup:    setup,
		existing: existing,
		expected: expected,
		handler:  handler,
		validate: validate,
	}
}

func (m DidMigrationScenario) Name() string {
	return m.name
}

func (m DidMigrationScenario) Setup() migrationsetup.ExtendedTestSetup {
	return m.setup()
}

func (m DidMigrationScenario) Existing() migrationsetup.MinimalDidDocInfoV1 {
	return m.existing
}

func (m DidMigrationScenario) Expected() didtypes.DidDoc {
	return m.expected
}

func (m DidMigrationScenario) Handler() func(ctx sdk.Context, migrationCtx appmigrations.MigrationContext) error {
	return m.handler
}

func (m DidMigrationScenario) Validate(actual didtypes.DidDoc) error {
	return m.validate(actual)
}

type ResourceMigrationScenario struct {
	name     string
	setup    func() migrationsetup.ExtendedTestSetup
	existing resourcetypesv1.MsgCreateResourcePayload
	didInfo  migrationsetup.MinimalDidDocInfoV1
	expected resourcetypes.Metadata
	handler  func(ctx sdk.Context, migrationCtx appmigrations.MigrationContext) error
	validate func(actual resourcetypes.Metadata) error
}

func NewResourceMigrationScenario(name string, setup func() migrationsetup.ExtendedTestSetup, existing resourcetypesv1.MsgCreateResourcePayload, didInfo migrationsetup.MinimalDidDocInfoV1, expected resourcetypes.Metadata, handler func(ctx sdk.Context, migrationCtx appmigrations.MigrationContext) error, validate func(actual resourcetypes.Metadata) error) ResourceMigrationScenario {
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

func (m ResourceMigrationScenario) Setup() migrationsetup.ExtendedTestSetup {
	return m.setup()
}

func (m ResourceMigrationScenario) Existing() resourcetypesv1.MsgCreateResourcePayload {
	return m.existing
}

func (m ResourceMigrationScenario) DidInfo() migrationsetup.MinimalDidDocInfoV1 {
	return m.didInfo
}

func (m ResourceMigrationScenario) Expected() resourcetypes.Metadata {
	return m.expected
}

func (m ResourceMigrationScenario) Handler() func(ctx sdk.Context, migrationCtx appmigrations.MigrationContext) error {
	return m.handler
}

func (m ResourceMigrationScenario) Validate(actual resourcetypes.Metadata) error {
	return m.validate(actual)
}

type Migrator interface {
	Migrate(ctx sdk.Context) error
}

type DidMigrator struct {
	migrations []DidMigrationScenario
}

func NewDidMigrator(migrations []DidMigrationScenario) DidMigrator {
	return DidMigrator{
		migrations: migrations,
	}
}

func (m DidMigrator) Migrate(ctx sdk.Context) error {
	for _, migration := range m.migrations {
		setup := migration.Setup()
		_, err := setup.CreateDidV1(migration.existing.Msg, []didtestssetup.SignInput{migration.existing.SignInput})
		if err != nil {
			return err
		}
		migrationCtx := appmigrations.NewMigrationContext(setup.Cdc, setup.Keeper, setup.ResourceKeeper)
		err = migration.Handler()(setup.SdkCtx, migrationCtx)
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
		_, err = setup.CreateResourceV1(&migration.existing, []didtestssetup.SignInput{migration.didInfo.SignInput})
		if err != nil {
			return err
		}
		migrationCtx := appmigrations.NewMigrationContext(setup.Cdc, setup.Keeper, setup.ResourceKeeper)
		err = migration.Handler()(setup.SdkCtx, migrationCtx)
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
