package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewMetadataFromContext(ctx sdk.Context, version string) Metadata {
	created := ctx.BlockTime().Format(time.RFC3339)

	return Metadata{Created: created, Updated: created, Deactivated: false, VersionId: version}
}

func (m *Metadata) Update(ctx sdk.Context, version string) {
	m.Updated = ctx.BlockTime().Format(time.RFC3339)
	m.VersionId = version
}
