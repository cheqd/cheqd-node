package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewMetadataFromContext(ctx sdk.Context, version string) Metadata {
	created := ctx.BlockTime()

	return Metadata{Created: created, Deactivated: false, VersionId: version}
}

func (m *Metadata) Update(ctx sdk.Context, version string) {
	updated := ctx.BlockTime()
	m.Updated = &updated
	m.VersionId = version
}
