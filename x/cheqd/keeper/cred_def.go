package keeper

import (
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strconv"
)

// GetCredDefCount get the total number of credDef
func (k Keeper) GetCredDefCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CredDefCountKey))
	byteKey := types.KeyPrefix(types.CredDefCountKey)
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

// SetCredDefCount set the total number of credDef
func (k Keeper) SetCredDefCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CredDefCountKey))
	byteKey := types.KeyPrefix(types.CredDefCountKey)
	bz := []byte(strconv.FormatUint(count, 10))
	store.Set(byteKey, bz)
}

// AppendCredDef appends a CredDef in the store with a new id and update the count
func (k Keeper) AppendCredDef(
	ctx sdk.Context,
	CredDef types.CredDef,
	metadata *types.Metadata,
) (*string, error) {
	// Create the CredDef
	count := k.GetCredDefCount(ctx)

	err := k.SetCredDef(ctx, CredDef, metadata)
	if err != nil {
		return nil, err
	}

	// Update CredDef count
	k.SetCredDefCount(ctx, count+1)

	return &CredDef.Id, nil
}

// SetCredDef set a specific cred def in the store
func (k Keeper) SetCredDef(ctx sdk.Context, CredDef types.CredDef, metadata *types.Metadata) error {
	stateValue, err := types.NewStateValue(&CredDef, metadata)
	if err != nil {
		return types.ErrSetToState.Wrap(err.Error())
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CredDefKey))
	b := k.cdc.MustMarshal(stateValue)
	store.Set(GetCredDefIDBytes(CredDef.Id), b)
	return nil
}

// GetCredDef returns a credDef from its id
func (k Keeper) GetCredDef(ctx sdk.Context, id string) (*types.StateValue, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CredDefKey))

	if !k.HasCredDef(ctx, id) {
		return nil, sdkerrors.ErrNotFound
	}

	var value types.StateValue
	var bytes = store.Get(GetCredDefIDBytes(id))
	if err := k.cdc.Unmarshal(bytes, &value); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidType, err.Error())
	}

	return &value, nil
}

// HasCredDef checks if the credDef exists in the store
func (k Keeper) HasCredDef(ctx sdk.Context, id string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CredDefKey))
	return store.Has(GetCredDefIDBytes(id))
}

// GetCredDefIDBytes returns the byte representation of the ID
func GetCredDefIDBytes(id string) []byte {
	return []byte(id)
}
