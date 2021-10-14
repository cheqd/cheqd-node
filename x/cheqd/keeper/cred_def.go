package keeper

import (
	"encoding/base64"
	"fmt"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto/tmhash"
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

// AppendCredDef TODO add Value
// AppendCredDef appends a credDef in the store with a new id and update the count
func (k Keeper) AppendCredDef(
	ctx sdk.Context,
	id string,
	schemaId string,
	tag string,
	signatureType string,
	clValue *types.CredDef_ClType,
) string {
	// Create the credDef
	count := k.GetCredDefCount(ctx)

	var credDef = types.CredDef{
		Id:            id,
		SchemaId:      schemaId,
		Tag:           tag,
		SignatureType: signatureType,
		Value:         clValue,
	}

	created := ctx.BlockTime().String()
	txHash := base64.StdEncoding.EncodeToString(tmhash.Sum(ctx.TxBytes()))

	stateValue := types.StateValue{
		Data: &types.StateValue_CredDef{
			CredDef: &credDef,
		},
		Timestamp: created,
		TxHash:    txHash,
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CredDefKey))
	value := k.cdc.MustMarshal(&stateValue)
	store.Set(GetCredDefIDBytes(credDef.Id), value)

	// Update credDef count
	k.SetCredDefCount(ctx, count+1)

	return id
}

// SetCredDef set a specific credDef in the store
func (k Keeper) SetCredDef(ctx sdk.Context, credDef types.CredDef) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CredDefKey))
	b := k.cdc.MustMarshal(&credDef)
	store.Set(GetCredDefIDBytes(credDef.Id), b)
}

// GetCredDef returns a credDef from its id
func (k Keeper) GetCredDef(ctx sdk.Context, id string) (*types.CredDef, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CredDefKey))

	var value types.StateValue
	k.cdc.MustUnmarshal(store.Get(GetCredDefIDBytes(id)), &value)

	switch data := value.Data.(type) {
	case *types.StateValue_CredDef:
		return data.CredDef, nil
	default:
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidType, fmt.Sprintf("State has unexpected type %T", data))
	}
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
