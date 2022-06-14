package keeper

import (
	cheqd_types "github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServer returns an implementation of the MsgServer interface for the provided Keeper.
func NewMsgServer(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func FindResource(k *Keeper, ctx *sdk.Context, inMemoryResources map[string]types.Resource, collectionId string, id string) (res types.Resource, found bool, err error) {
	// Look in inMemory dict
	value, found := inMemoryResources[collectionId+id]
	if found {
		return value, true, nil
	}

	// Look in state
	if k.HasResource(ctx, collectionId, id) {
		value, err := k.GetResource(ctx, collectionId, id)
		if err != nil {
			return types.Resource{}, false, err
		}

		return value, true, nil
	}

	return types.Resource{}, false, nil
}

func VerifySignature(k *Keeper, ctx *sdk.Context, inMemoryResources map[string]types.Resource, message []byte, signature cheqd_types.SignInfo) error {
	//TODO: implement
	return nil
}
