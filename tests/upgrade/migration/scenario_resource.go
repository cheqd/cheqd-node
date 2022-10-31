package migration

import (
	"bytes"
	"fmt"

	appmigrations "github.com/cheqd/cheqd-node/app/migrations"
	resourcetestssetup "github.com/cheqd/cheqd-node/x/resource/tests/setup"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var checksumResource = resourcetypes.Resource{
	Header: &resourcetypes.ResourceHeader{
		CollectionId: "collectionId",
		Id:           "resourceId",
		Name:         "name",
		ResourceType: "resourceType",
		MediaType:    "mediaType",
		Created:      "created",
		Checksum:     []byte("checksum"),
	},
	Data: []byte("data"),
}

var ResourceChecksumScenario = NewResourceMigrationScenario(
	"ResourceChecksum",
	resourcetestssetup.Setup,
	checksumResource,
	func(ctx sdk.Context) error {
		setup := resourcetestssetup.Setup()
		return appmigrations.MigrateResourceV1(ctx, setup.Keeper, setup.ResourceKeeper)
	},
	func(actual resourcetypes.Resource) error {
		if !bytes.Equal(actual.Header.Checksum, checksumResource.Header.Checksum) {
			return fmt.Errorf("expected checksum %v, got %v", checksumResource.Header.Checksum, actual.Header.Checksum)
		}
		return nil
	},
)
