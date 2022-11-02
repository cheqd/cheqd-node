package migration

import (
	"bytes"
	"fmt"
	"path/filepath"

	appmigrations "github.com/cheqd/cheqd-node/app/migrations"
	cheqdkeeper "github.com/cheqd/cheqd-node/x/cheqd/keeper"
	cheqdtestssetup "github.com/cheqd/cheqd-node/x/cheqd/tests/setup"
	cheqdtypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	resourcekeeper "github.com/cheqd/cheqd-node/x/resource/keeper"
	resourcetestssetup "github.com/cheqd/cheqd-node/x/resource/tests/setup"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	didDoc                         cheqdtypes.MsgCreateDidPayload
	didInfo                        cheqdtestssetup.MinimalDidInfo
	existingChecksumResource       resourcetypes.MsgCreateResourcePayload
	expectedChecksumResourceHeader resourcetypes.ResourceHeader
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

	didInfo = cheqdtestssetup.MinimalDidInfo{
		Msg: &didDoc,
		SignInput: cheqdtestssetup.SignInput{
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
		func(ctx sdk.Context, cheqdKeeper cheqdkeeper.Keeper, resourceKeeper resourcekeeper.Keeper) error {
			return appmigrations.MigrateResourceV1(ctx, cheqdKeeper, resourceKeeper)
		},
		func(actual resourcetypes.ResourceHeader) error {
			if !bytes.Equal(actual.Checksum, expectedChecksumResourceHeader.Checksum) {
				return fmt.Errorf("expected checksum %v, got %v", expectedChecksumResourceHeader.Checksum, actual.Checksum)
			}
			return nil
		},
	)
	return nil
}
