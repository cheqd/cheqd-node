package types

import (
	"time"

	"github.com/cheqd/cheqd-node/x/did/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewMetadataFromContext(ctx sdk.Context) Metadata {
	created := ctx.BlockTime().Format(time.RFC3339)
	txHash := utils.GetTxHash(ctx.TxBytes())

	return Metadata{Created: created, Deactivated: false, VersionId: txHash}
}

func (m *Metadata) Update(ctx sdk.Context) {
	m.Updated = ctx.BlockTime().Format(time.RFC3339)
}
