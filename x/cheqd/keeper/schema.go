package keeper

import (
	"encoding/base64"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"strconv"
)

// GetSchemaCount get the total number of schema
func (k Keeper) GetSchemaCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SchemaCountKey))
	byteKey := types.KeyPrefix(types.SchemaCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	count, err := strconv.ParseUint(string(bz), 10, 64)
	if err != nil {
		// Panic because the count should be always formattable to iint64
		panic("cannot decode count")
	}

	return count
}

// SetSchemaCount set the total number of schema
func (k Keeper) SetSchemaCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SchemaCountKey))
	byteKey := types.KeyPrefix(types.SchemaCountKey)
	bz := []byte(strconv.FormatUint(count, 10))
	store.Set(byteKey, bz)
}

// AppendSchema appends a schema in the store with a new id and update the count
func (k Keeper) AppendSchema(
	ctx sdk.Context,
	id string,
	schemaType string,
	name string,
	version string,
	attrNames []string,
	controller []string,
) string {
	// Create the schema
	count := k.GetSchemaCount(ctx)
	var schema = types.Schema{
		Id:         id,
		Name:       name,
		Type:       schemaType,
		Version:    version,
		AttrNames:  attrNames,
		Controller: controller,
	}

	created := ctx.BlockTime().String()
	txHash := base64.StdEncoding.EncodeToString(tmhash.Sum(ctx.TxBytes()))

	stateValue := types.StateValue{
		Data: &types.StateValue_Schema{
			Schema: &schema,
		},
		Timestamp: created,
		TxHash:    txHash,
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SchemaKey))
	value := k.cdc.MustMarshal(&stateValue)
	store.Set(GetSchemaIDBytes(schema.Id), value)

	// Update schema count
	k.SetSchemaCount(ctx, count+1)

	return id
}

// GetSchema returns a schema from its id
func (k Keeper) GetSchema(ctx sdk.Context, id string) (*types.StateValue, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SchemaKey))

	if !k.HasSchema(ctx, id) {
		return nil, sdkerrors.ErrNotFound
	}

	var value types.StateValue
	var bytes = store.Get(GetSchemaIDBytes(id))
	if err := k.cdc.Unmarshal(bytes, &value); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidType, err.Error())
	}

	return &value, nil
}

// HasSchema checks if the schema exists in the store
func (k Keeper) HasSchema(ctx sdk.Context, id string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SchemaKey))
	return store.Has(GetSchemaIDBytes(id))
}

// GetSchemaIDBytes returns the byte representation of the ID
func GetSchemaIDBytes(id string) []byte {
	return []byte(id)
}
