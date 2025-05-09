package types

import (
	context "context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewMetadataFromContext(ctx context.Context, version string) Metadata {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	created := sdkCtx.BlockTime()

	return Metadata{Created: created, Deactivated: false, VersionId: version}
}

func (m *Metadata) Update(ctx context.Context, version string) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	updated := sdkCtx.BlockTime()
	m.Updated = &updated
	m.VersionId = version
}
