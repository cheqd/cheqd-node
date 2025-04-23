package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	store "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
)

type Keeper struct {
	cdc codec.BinaryCodec
	// storeKey     storetypes.StoreKey
	storeService store.KVStoreService
	paramSpace   types.ParamSubspace
	portKeeper   types.PortKeeper
	scopedKeeper exported.ScopedKeeper
	Schema       collections.Schema

	Port collections.Item[string]
	// ResourceCount stores the total number of resources
	ResourceCount collections.Item[uint64]

	// Resources stores resource metadata by collection ID and resource ID
	ResourceMetadata collections.Map[collections.Pair[string, string], types.Metadata]

	// ResourceData stores resource data by collection ID and resource ID
	ResourceData collections.Map[collections.Pair[string, string], []byte]

	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/resource module account.
	authority string

	Params collections.Item[types.FeeParams]
}

func NewKeeper(cdc codec.BinaryCodec, storeService store.KVStoreService, paramSpace types.ParamSubspace, portKeeper types.PortKeeper, scopedKeeper exported.ScopedKeeper, authority string) *Keeper {
	sb := collections.NewSchemaBuilder(storeService)
	// Define the port collection

	k := &Keeper{
		cdc: cdc,
		// storeKey:     storeKey,
		storeService: storeService,
		paramSpace:   paramSpace,
		portKeeper:   portKeeper,
		scopedKeeper: scopedKeeper,
		authority:    authority,
		Port: collections.NewItem(
			sb,
			collections.NewPrefix(types.ResourcePortIDKey), // prefix for the port
			"port",                  // key for the port
			collections.StringValue, // type of the value
		),
		ResourceCount: collections.NewItem(
			sb,
			collections.NewPrefix(types.ResourceCountKey),
			"resource_count",
			collections.Uint64Value,
		),

		ResourceMetadata: collections.NewMap(
			sb,
			collections.NewPrefix(types.ResourceMetadataKey),
			"resource_metadata",
			collections.PairKeyCodec(collections.StringKey, collections.StringKey),
			codec.CollValue[types.Metadata](cdc),
		),

		ResourceData: collections.NewMap(
			sb,
			collections.NewPrefix(types.ResourceDataKey),
			"resource_data",
			collections.PairKeyCodec(collections.StringKey, collections.StringKey),
			collections.BytesValue,
		),
		Params: collections.NewItem(
			sb,
			collections.NewPrefix(types.ParamStoreKeyFeeParams),
			"params",
			codec.CollValue[types.FeeParams](cdc),
		),
	}

	// Build the schema
	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema
	return k
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetAuthority() string {
	return k.authority
}

func (k Keeper) SetParams(ctx context.Context, params types.FeeParams) error {
	err := k.Params.Set(ctx, params)
	if err != nil {
		return err
	}
	return nil
}

// GetParams gets the x/resource module parameters.
func (k Keeper) GetParams(ctx context.Context) (params types.FeeParams, err error) {
	return k.Params.Get(ctx)
}
