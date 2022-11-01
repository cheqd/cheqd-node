package migration

import (
	"fmt"
	"encoding/base64"
	"bytes"

	appmigrations "github.com/cheqd/cheqd-node/app/migrations"
	cheqdtypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	cheqdtestssetup "github.com/cheqd/cheqd-node/x/cheqd/tests/setup"
	resourcetestssetup "github.com/cheqd/cheqd-node/x/resource/tests/setup"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var err error
var didDoc cheqdtypes.MsgCreateDidPayload
var didInfo cheqdtestssetup.MinimalDidInfo
var existingChecksumResource resourcetypes.MsgCreateResourcePayload
var expectedChecksumResourceHeader resourcetypes.ResourceHeader
var ResourceChecksumScenario ResourceMigrationScenario


func InitResourceChecksumScenario() error {
	err = Loader("generated/payload/diddoc_multibase_16.json", &didDoc)
	if err != nil {
		return err
	}
	var signInput SignInput
	err = Loader("generated/keys/signinput_multibase_16.json", &signInput)
	if err != nil {
		return err
	}
	err = Loader("generated/payload/existing/resource_checksum.json", &existingChecksumResource)
	if err != nil {
		return err
	}
	err = Loader("generated/payload/expected/resource_checksum.json", &expectedChecksumResourceHeader)
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
			Key: privateKey,
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
