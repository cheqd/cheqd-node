package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	Keeper struct {
		cdc           codec.BinaryCodec
		storeService  store.KVStoreService
		paramSpace    types.ParamSubspace
		accountKeeper types.AccountKeeper
		bankkeeper    types.BankKeeper
		stakingKeeper types.StakingKeeper
		authority     string
		Schema        collections.Schema

		DidNamespace     collections.Item[string]
		DidCount         collections.Item[uint64]
		LatestDidVersion collections.Map[string, string]
		DidDocuments     collections.Map[collections.Pair[string, string], types.DidDocWithMetadata]
		Params           collections.Item[types.FeeParams]
	}
)

func NewKeeper(cdc codec.BinaryCodec, storeService store.KVStoreService, paramSpace types.ParamSubspace, ak types.AccountKeeper, bk types.BankKeeper, sk types.StakingKeeper, authority string) *Keeper {
	sb := collections.NewSchemaBuilder(storeService)
	k := &Keeper{
		cdc:              cdc,
		storeService:     storeService,
		paramSpace:       paramSpace,
		accountKeeper:    ak,
		bankkeeper:       bk,
		stakingKeeper:    sk,
		authority:        authority,
		DidNamespace:     collections.NewItem(sb, types.DidNamespaceKeyPrefix, "did-namespace:", collections.StringValue),
		DidCount:         collections.NewItem(sb, types.DidDocCountKeyPrefix, "did-count:", collections.Uint64Value),
		LatestDidVersion: collections.NewMap(sb, types.LatestDidDocVersionKeyPrefix, "latest-did", collections.StringKey, collections.StringValue),
		DidDocuments:     collections.NewMap(sb, types.DidDocVersionKeyPrefix, "did-version", collections.PairKeyCodec(collections.StringKey, collections.StringKey), codec.CollValue[types.DidDocWithMetadata](cdc)),
		Params:           collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.FeeParams](cdc)),
	}
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
