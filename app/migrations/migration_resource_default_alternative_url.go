package migrations

import (
	"fmt"

	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MigrateResourceDefaultAlternativeUrl(sctx sdk.Context, mctx MigrationContext) error {
	namespace := mctx.didKeeperNew.GetDidNamespace(&sctx)

	return MigrateResourceSimple(sctx, mctx, func(resource *resourcetypes.ResourceWithMetadata) {
		resource.Metadata.AlsoKnownAs = append(resource.Metadata.AlsoKnownAs, &resourcetypes.AlternativeUri{
			Uri:         fmt.Sprintf("did:cheqd:%s:%s/resources/%s", namespace, resource.Metadata.CollectionId, resource.Metadata.Id),
			Description: "did-url",
		})
	})
}
