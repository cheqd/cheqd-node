package migration

import (
	"bytes"
	"fmt"
	"path/filepath"

	appmigrations "github.com/cheqd/cheqd-node/app/migrations"
	migrationsetup "github.com/cheqd/cheqd-node/tests/upgrade/migration/setup"
	didtestssetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	didtypesv1 "github.com/cheqd/cheqd-node/x/did/types/v1"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	resourcetypesv1 "github.com/cheqd/cheqd-node/x/resource/types/v1"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	didDoc                         didtypesv1.MsgCreateDidPayload
	didInfo                        migrationsetup.MinimalDidDocInfoV1
	existingChecksumResource       resourcetypesv1.MsgCreateResourcePayload
	expectedChecksumResourceHeader resourcetypes.Metadata
	ResourceChecksumScenario       ResourceMigrationScenario
)

func InitResourceChecksumScenario() error {
	err := Loader(filepath.Join("payload", "existing", "diddoc_uuid.json"), &didDoc)
	if err != nil {
		fmt.Println("Error loading didDoc")
		return err
	}
	var signInput SignInput
	err = Loader(filepath.Join("keys", "signinput_uuid.json"), &signInput)
	if err != nil {
		fmt.Println("Error loading signInput")
		return err
	}
	err = Loader(filepath.Join("payload", "existing", "resource_checksum.json"), &existingChecksumResource)
	if err != nil {
		fmt.Println("Error loading existingChecksumResource")
		return err
	}
	err = Loader(filepath.Join("payload", "expected", "resource_checksum.json"), &expectedChecksumResourceHeader)
	if err != nil {
		fmt.Println("Error loading expectedChecksumResourceHeader")
		return err
	}

	didInfo = migrationsetup.MinimalDidDocInfoV1{
		Msg: &didDoc,
		SignInput: didtestssetup.SignInput{
			VerificationMethodId: signInput.VerificationMethodId,
			Key:                  signInput.PrivateKey,
		},
	}

	ResourceChecksumScenario = NewResourceMigrationScenario(
		"ResourceChecksum",
		migrationsetup.NewExtendedSetup,
		existingChecksumResource,
		didInfo,
		expectedChecksumResourceHeader,
		func(ctx sdk.Context, migrationCtx appmigrations.MigrationContext) error {
			return appmigrations.MigrateResourceV1(ctx, migrationCtx)
		},
		func(actual resourcetypes.Metadata) error {
			if !bytes.Equal(actual.Checksum, expectedChecksumResourceHeader.Checksum) {
				return fmt.Errorf("expected checksum %v, got %v", expectedChecksumResourceHeader.Checksum, actual.Checksum)
			}
			return nil
		},
	)
	return nil
}
