package migration

import (
	"bytes"
	"fmt"
	"path/filepath"

	appmigrations "github.com/cheqd/cheqd-node/app/migrations"
	didkeeper "github.com/cheqd/cheqd-node/x/did/keeper"
	didtestssetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	didtypesv1 "github.com/cheqd/cheqd-node/x/did/types"
	resourcekeeper "github.com/cheqd/cheqd-node/x/resource/keeper"
	resourcetestssetup "github.com/cheqd/cheqd-node/x/resource/tests/setup"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	err                            error
	didDoc                         didtypes.MsgCreateDidPayload
	didInfo                        cheqdtestssetup.MinimalDidInfo
	existingChecksumResource       resourcetypes.MsgCreateResourcePayload
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

	didInfo = didtestssetup.MinimalDidDocInfoV1{
		Msg: &didDoc,
		SignInput: didtestssetup.SignInput{
			VerificationMethodId: signInput.VerificationMethodId,
			Key:                  signInput.PrivateKey,
		},
	}

	ResourceChecksumScenario = NewResourceMigrationScenario(
		"ResourceChecksum",
		resourcetestssetup.Setup,
		existingChecksumResource,
		didInfo,
		expectedChecksumResourceHeader,
		func(ctx sdk.Context, didKeeper didkeeper.Keeper, resourceKeeper resourcekeeper.Keeper) error {
			return appmigrations.MigrateResourceV1(ctx, didKeeper, resourceKeeper)
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
