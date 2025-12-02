package keeper

import (
	"context"
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
		oracleKeeper  types.OracleKeeper
		authority     string
		Schema        collections.Schema

		DidNamespace     collections.Item[string]
		DidCount         collections.Item[uint64]
		LatestDidVersion collections.Map[string, string]
		DidDocuments     collections.Map[collections.Pair[string, string], types.DidDocWithMetadata]
		Paramstore       collections.Item[types.FeeParams]
	}
)

func NewKeeper(cdc codec.BinaryCodec, storeService store.KVStoreService, paramSpace types.ParamSubspace, ak types.AccountKeeper, bk types.BankKeeper, sk types.StakingKeeper, ok types.OracleKeeper, authority string) *Keeper {
	sb := collections.NewSchemaBuilder(storeService)
	k := &Keeper{
		cdc:              cdc,
		storeService:     storeService,
		paramSpace:       paramSpace,
		accountKeeper:    ak,
		bankkeeper:       bk,
		stakingKeeper:    sk,
		oracleKeeper:     ok,
		authority:        authority,
		DidNamespace:     collections.NewItem(sb, types.DidNamespaceKeyPrefix, "did_namespace", collections.StringValue),
		DidCount:         collections.NewItem(sb, types.DidDocCountKeyPrefix, "did_count", collections.Uint64Value),
		LatestDidVersion: collections.NewMap(sb, types.LatestDidDocVersionKeyPrefix, "latest_did", collections.StringKey, collections.StringValue),
		DidDocuments:     collections.NewMap(sb, types.DidDocVersionKeyPrefix, "did_version", collections.PairKeyCodec(collections.StringKey, collections.StringKey), codec.CollValue[types.DidDocWithMetadata](cdc)),
		Paramstore:       collections.NewItem(sb, types.ParamStoreKeyFeeParams, "params", codec.CollValue[types.FeeParams](cdc)),
	}
	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema
	return k
}

func (k Keeper) Logger(ctx context.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetAuthority() string {
	return k.authority
}

// GetParams gets the auth module's parameters.
func (k Keeper) GetParams(ctx context.Context) (types.FeeParams, error) {
	return k.Paramstore.Get(ctx)
}

func (k Keeper) SetParams(ctx context.Context, params types.FeeParams) error {
	err := k.Paramstore.Set(ctx, params)
	if err != nil {
		return err
	}
	return nil
}
