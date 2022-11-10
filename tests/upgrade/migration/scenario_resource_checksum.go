package migration

import (
	"bytes"
	"encoding/base64"
	"fmt"

	appmigrations "github.com/cheqd/cheqd-node/app/migrations"
	cheqdtestssetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	didtypes "github.com/cheqd/cheqd-node/x/did/types/v1"
	resourcetestssetup "github.com/cheqd/cheqd-node/x/resource/tests/setup"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	err                            error
	didDoc                         didtypes.MsgCreateDidPayload
	didInfo                        cheqdtestssetup.MinimalDidInfo
	existingChecksumResource       resourcetypes.MsgCreateResourcePayload
	expectedChecksumResourceHeader resourcetypes.ResourceHeader
	ResourceChecksumScenario       ResourceMigrationScenario
)

func InitResourceChecksumScenario() error {
	err = Loader(GENERATED_JSON_DIR+"/payload/existing/diddoc_multibase_16.json", &didDoc)
	if err != nil {
		return err
	}
	var signInput SignInput
	err = Loader(GENERATED_JSON_DIR+"/keys/signinput_multibase_16.json", &signInput)
	if err != nil {
		return err
	}
	err = Loader(GENERATED_JSON_DIR+"/payload/existing/resource_checksum.json", &existingChecksumResource)
	if err != nil {
		return err
	}
	err = Loader(GENERATED_JSON_DIR+"/payload/expected/resource_checksum.json", &expectedChecksumResourceHeader)
	if err != nil {
		return err
	}
	privateKey, err := base64.StdEncoding.DecodeString(signInput.PrivateKey)
	if err != nil {
		return err
	}
	didInfo = cheqdtestssetup.MinimalDidInfo{
		Msg: &didDoc,
		SignInput: cheqdtestssetup.SignInput{
			VerificationMethodId: signInput.VerificationMethodId,
			Key:                  privateKey,
		},
	}

	ResourceChecksumScenario = NewResourceMigrationScenario(
		"ResourceChecksum",
		resourcetestssetup.Setup,
		existingChecksumResource,
		didInfo,
		expectedChecksumResourceHeader,
		func(ctx sdk.Context) error {
			setup := resourcetestssetup.Setup()
			return appmigrations.MigrateResourceV1(ctx, setup.Keeper, setup.ResourceKeeper)
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
