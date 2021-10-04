package keeper

import (
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
)

// GetDidCount get the total number of did
func (k Keeper) GetDidCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DidCountKey))
	byteKey := types.KeyPrefix(types.DidCountKey)
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

// SetDidCount set the total number of did
func (k Keeper) SetDidCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DidCountKey))
	byteKey := types.KeyPrefix(types.DidCountKey)
	bz := []byte(strconv.FormatUint(count, 10))
	store.Set(byteKey, bz)
}

// AppendDid appends a did in the store with a new id and update the count
func (k Keeper) AppendDid(
	ctx sdk.Context,
	id string,
	controller []string,
	verificationMethod []*types.VerificationMethod,
	authentication []string,
	assertionMethod []string,
	capabilityInvocation []string,
	capabilityDelegation []string,
	keyAgreement []string,
	alsoKnownAs []string,
	service []*types.DidService,
) string {
	// Create the did
	count := k.GetDidCount(ctx)
	var did = types.Did{
		Id:                   id,
		Controller:           controller,
		VerificationMethod:   verificationMethod,
		Authentication:       authentication,
		AssertionMethod:      assertionMethod,
		CapabilityInvocation: capabilityInvocation,
		CapabilityDelegation: capabilityDelegation,
		KeyAgreement:         keyAgreement,
		AlsoKnownAs:          alsoKnownAs,
		Service:              service,
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DidKey))
	value := k.cdc.MustMarshal(&did)
	store.Set(GetDidIDBytes(did.Id), value)

	// Update did count
	k.SetDidCount(ctx, count+1)

	return id
}

// SetDid set a specific did in the store
func (k Keeper) SetDid(ctx sdk.Context, did types.Did) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DidKey))
	b := k.cdc.MustMarshal(&did)
	store.Set(GetDidIDBytes(did.Id), b)
}

// GetDid returns a did from its id
func (k Keeper) GetDid(ctx sdk.Context, id string) types.Did {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DidKey))
	var did types.Did
	k.cdc.MustUnmarshal(store.Get(GetDidIDBytes(id)), &did)
	return did
}

// areOwners returns a bool are received authors can control this DID
func (k Keeper) areDidOwners(ctx sdk.Context, id string, authors []string) bool {
	return utils.CompareOwners(authors, k.GetDid(ctx, id).Controller)
}

// HasDid checks if the did exists in the store
func (k Keeper) HasDid(ctx sdk.Context, id string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DidKey))
	return store.Has(GetDidIDBytes(id))
}

// RemoveDid removes a did from the store
func (k Keeper) RemoveDid(ctx sdk.Context, id string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DidKey))
	store.Delete(GetDidIDBytes(id))
}

// GetAllDid returns all did
func (k Keeper) GetAllDid(ctx sdk.Context) (list []types.Did) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DidKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Did
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetDidIDBytes returns the byte representation of the ID
func GetDidIDBytes(id string) []byte {
	return []byte(id)
}

// GetDidIDFromBytes returns ID in uint64 format from a byte array
func GetDidIDFromBytes(bz []byte) string {
	return string(bz)
}
