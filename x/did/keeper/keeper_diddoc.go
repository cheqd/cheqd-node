package keeper

import (
	"strconv"
	"strings"

	"github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cheqd/cheqd-node/x/did/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// GetDidCount get the total number of did
func (k Keeper) GetDidDocCount(ctx *sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)

	key := utils.StrBytes(types.DidDocCountKey)
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
	store := ctx.KVStore(k.storeKey)

	key := utils.StrBytes(types.DidDocCountKey)
	valueBytes := []byte(strconv.FormatUint(count, 10))

	store.Set(key, valueBytes)
}

func (k Keeper) AddNewDidDocVersion(ctx *sdk.Context, didDoc *types.DidDocWithMetadata) error {
	// Check if the diddoc version already exists
	if k.HasDidDocVersion(ctx, didDoc.DidDoc.Id, didDoc.Metadata.VersionId) {
		return types.ErrDidDocExists.Wrapf("diddoc version already exists for did %s, version %s", didDoc.DidDoc.Id, didDoc.Metadata.VersionId)
	}

	// Link to the previous version if it exists
	if k.HasDidDoc(ctx, didDoc.DidDoc.Id) {
		latestVersionID, err := k.GetLatestDidDocVersion(ctx, didDoc.DidDoc.Id)
		if err != nil {
			return err
		}

		latestVersion, err := k.GetDidDocVersion(ctx, didDoc.DidDoc.Id, latestVersionID)
		if err != nil {
			return err
		}

		// Update version links
		latestVersion.Metadata.NextVersionId = didDoc.Metadata.VersionId
		didDoc.Metadata.PreviousVersionId = latestVersion.Metadata.VersionId

		// Update previous version with override
		err = k.SetDidDocVersion(ctx, &latestVersion, true)
		if err != nil {
			return err
		}
	}

	// Update latest version
	err := k.SetLatestDidDocVersion(ctx, didDoc.DidDoc.Id, didDoc.Metadata.VersionId)
	if err != nil {
		return err
	}

	// Write new version (no override)
	return k.SetDidDocVersion(ctx, didDoc, false)
}

func (k Keeper) GetLatestDidDoc(ctx *sdk.Context, did string) (types.DidDocWithMetadata, error) {
	latestVersionID, err := k.GetLatestDidDocVersion(ctx, did)
	if err != nil {
		return types.DidDocWithMetadata{}, err
	}

	latestVersion, err := k.GetDidDocVersion(ctx, did, latestVersionID)
	if err != nil {
		return types.DidDocWithMetadata{}, err
	}

	return latestVersion, nil
}

// SetDid set a specific did in the store. Updates DID counter if the DID is new.
func (k Keeper) SetDidDocVersion(ctx *sdk.Context, value *types.DidDocWithMetadata, override bool) error {
	if !override && k.HasDidDocVersion(ctx, value.DidDoc.Id, value.Metadata.VersionId) {
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
	valueBytes := store.Get(types.GetDidDocVersionKey(id, version))
	k.cdc.MustUnmarshal(valueBytes, &value)

	return value, nil
}

func (k Keeper) GetAllDidDocVersions(ctx *sdk.Context, did string) ([]*types.Metadata, error) {
	store := ctx.KVStore(k.storeKey)

	result := make([]*types.Metadata, 0)

	versionIterator := sdk.KVStorePrefixIterator(store, types.GetDidDocVersionsPrefix(did))
	defer closeIteratorOrPanic(versionIterator)

	for ; versionIterator.Valid(); versionIterator.Next() {
		// Get the diddoc
		var didDoc types.DidDocWithMetadata
		k.cdc.MustUnmarshal(versionIterator.Value(), &didDoc)

		result = append(result, didDoc.Metadata)
	}

	return result, nil
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
	valueBytes := utils.StrBytes(version)
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

func (k Keeper) HasDidDoc(ctx *sdk.Context, id string) bool {
	return k.HasLatestDidDocVersion(ctx, id)
}

func (k Keeper) HasLatestDidDocVersion(ctx *sdk.Context, id string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetLatestDidDocVersionKey(id))
}

func (k Keeper) HasDidDocVersion(ctx *sdk.Context, id, version string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetDidDocVersionKey(id, version))
}

func (k Keeper) IterateDids(ctx *sdk.Context, callback func(did string) (continue_ bool)) {
	store := ctx.KVStore(k.storeKey)
	latestVersionIterator := sdk.KVStorePrefixIterator(store, types.GetLatestDidDocVersionPrefix())
	defer closeIteratorOrPanic(latestVersionIterator)

	for ; latestVersionIterator.Valid(); latestVersionIterator.Next() {
		// Get did from key
		key := string(latestVersionIterator.Key())
		did := strings.Join(strings.Split(key, ":")[1:], ":")

		if !callback(did) {
			break
		}
	}
}

func (k Keeper) IterateDidDocVersions(ctx *sdk.Context, did string, callback func(version types.DidDocWithMetadata) (continue_ bool)) {
	store := ctx.KVStore(k.storeKey)
	versionIterator := sdk.KVStorePrefixIterator(store, types.GetDidDocVersionsPrefix(did))
	defer closeIteratorOrPanic(versionIterator)

	for ; versionIterator.Valid(); versionIterator.Next() {
		var didDoc types.DidDocWithMetadata
		k.cdc.MustUnmarshal(versionIterator.Value(), &didDoc)

		if !callback(didDoc) {
			break
		}
	}
}

func (k Keeper) IterateAllDidDocVersions(ctx *sdk.Context, callback func(version types.DidDocWithMetadata) (continue_ bool)) {
	store := ctx.KVStore(k.storeKey)
	allVersionsIterator := sdk.KVStorePrefixIterator(store, []byte(types.DidDocVersionKey))
	defer closeIteratorOrPanic(allVersionsIterator)

	for ; allVersionsIterator.Valid(); allVersionsIterator.Next() {
		var didDoc types.DidDocWithMetadata
		k.cdc.MustUnmarshal(allVersionsIterator.Value(), &didDoc)

		if !callback(didDoc) {
			break
		}
	}
}

// GetAllDidDocs returns all did
// Loads all DIDs in memory. Use only for genesis export.
func (k Keeper) GetAllDidDocs(ctx *sdk.Context) ([]*types.DidDocVersionSet, error) {
	var didDocs []*types.DidDocVersionSet
	var err error

	k.IterateDids(ctx, func(did string) bool {
		var latestVersion string
		latestVersion, err = k.GetLatestDidDocVersion(ctx, did)
		if err != nil {
			return false
		}

		didDocVersionSet := types.DidDocVersionSet{
			LatestVersion: latestVersion,
		}

		k.IterateDidDocVersions(ctx, did, func(version types.DidDocWithMetadata) bool {
			didDocVersionSet.DidDocs = append(didDocVersionSet.DidDocs, &version)

			return true
		})

		didDocs = append(didDocs, &didDocVersionSet)

		return true
	})

	if err != nil {
		return nil, err
	}

	return didDocs, nil
}

func closeIteratorOrPanic(iterator sdk.Iterator) {
	err := iterator.Close()
	if err != nil {
		panic(err.Error())
	}
}
