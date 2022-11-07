package keeper

import (
	"strconv"
	"strings"

	"github.com/cheqd/cheqd-node/x/did/types"
	. "github.com/cheqd/cheqd-node/x/did/utils"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// GetDidCount get the total number of did
func (k Keeper) GetDidDocCount(ctx *sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), StrBytes(types.DidDocCountKey))

	key := StrBytes(types.DidDocCountKey)
	valueBytes := store.Get(key)

	// Count doesn't exist: no element
	if valueBytes == nil {
		return 0
	}

	// Parse bytes
	count, err := strconv.ParseUint(string(valueBytes), 10, 64)
	if err != nil {
		// Panic because the count should be always formattable to iint64
		panic("cannot decode count")
	}

	return count
}

// SetDidCount set the total number of did
func (k Keeper) SetDidDocCount(ctx *sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), StrBytes(types.DidDocCountKey))

	key := StrBytes(types.DidDocCountKey)
	valueBytes := []byte(strconv.FormatUint(count, 10))

	store.Set(key, valueBytes)
}

// SetDid set a specific did in the store. Updates DID counter if the DID is new.
func (k Keeper) SetDidDocVersion(ctx *sdk.Context, value *types.DidDocWithMetadata) error {
	if k.HasDidDocVersion(ctx, value.DidDoc.Id, value.Metadata.VersionId) {
		return types.ErrDidDocExists.Wrap("diddoc version already exists")
	}

	// Create the diddoc version
	store := ctx.KVStore(k.storeKey)

	key := types.GetDidDocVersionKey(value.DidDoc.Id, value.Metadata.VersionId)
	valueBytes := k.cdc.MustMarshal(value)
	store.Set(key, valueBytes)

	return nil
}

// GetDid returns a did from its id
func (k Keeper) GetDidDocVersion(ctx *sdk.Context, id, version string) (types.DidDocWithMetadata, error) {
	store := ctx.KVStore(k.storeKey)

	if !k.HasDidDocVersion(ctx, id, version) {
		return types.DidDocWithMetadata{}, sdkerrors.ErrNotFound.Wrap("diddoc version not found")
	}

	var value types.DidDocWithMetadata
	key := store.Get(StrBytes(id))
	k.cdc.MustUnmarshal(key, &value)

	return value, nil
}

// SetDidDocLatestVersion sets the latest version id value for a diddoc
func (k Keeper) SetLatestDidDocVersion(ctx *sdk.Context, did, version string) error {
	// Update counter. We use latest version as existence indicator.
	if !k.HasLatestDidDocVersion(ctx, did) {
		count := k.GetDidDocCount(ctx)
		k.SetDidDocCount(ctx, count+1)
	}

	store := ctx.KVStore(k.storeKey)

	key := types.GetLatestDidDocVersionKey(did)
	valueBytes := StrBytes(version)
	store.Set(key, valueBytes)

	return nil
}

// GetDidDocLatestVersion returns the latest version id value for a diddoc
func (k Keeper) GetLatestDidDocVersion(ctx *sdk.Context, id string) (string, error) {
	store := ctx.KVStore(k.storeKey)

	if !k.HasLatestDidDocVersion(ctx, id) {
		return "", sdkerrors.ErrNotFound.Wrap(id)
	}

	return string(store.Get(types.GetLatestDidDocVersionKey(id))), nil
}

// HasDid checks if the did exists in the store
func (k Keeper) HasLatestDidDocVersion(ctx *sdk.Context, id string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetLatestDidDocVersionKey(id))
}

// HasDid checks if the did exists in the store
func (k Keeper) HasDidDocVersion(ctx *sdk.Context, id, version string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetDidDocVersionKey(id, version))
}

// GetAllDidDocs returns all did
// Loads all DIDs in memory. Use only for genesis export.
func (k Keeper) GetAllDidDocs(ctx *sdk.Context) []types.DidDocVersionSet {
	store := ctx.KVStore(k.storeKey)
	latestVresionIterator := sdk.KVStorePrefixIterator(store, types.GetLatestDidDocVersionPrefix())
	defer closeIteratorOrPanic(latestVresionIterator)

	var didDocs []types.DidDocVersionSet

	for ; latestVresionIterator.Valid(); latestVresionIterator.Next() {
		// Get did from key
		key := string(latestVresionIterator.Key())
		did := strings.Split(key, ":")[1]

		// Get the latest version
		latestVersion := string(latestVresionIterator.Value())

		didDocVersionSet := types.DidDocVersionSet{
			LatestVersion: latestVersion,
		}

		// Get all versions
		versionIterator := sdk.KVStorePrefixIterator(store, types.GetDidDocVersionsPrefix(did))
		defer closeIteratorOrPanic(versionIterator)

		for ; versionIterator.Valid(); versionIterator.Next() {
			// Get the diddoc
			var didDoc types.DidDocWithMetadata
			k.cdc.MustUnmarshal(versionIterator.Value(), &didDoc)

			didDocVersionSet.DidDocs = append(didDocVersionSet.DidDocs, &didDoc)
		}
	}

	return didDocs
}

func closeIteratorOrPanic(iterator sdk.Iterator) {
	err := iterator.Close()
	if err != nil {
		panic(err.Error())
	}
}
