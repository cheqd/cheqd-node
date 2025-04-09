package keeper

import (
	"context"
	"strconv"
	"strings"

	storetypes "cosmossdk.io/store/types"
	"github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cheqd/cheqd-node/x/did/utils"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// GetDidCount get the total number of did
func (k Keeper) GetDidDocCount(ctx *context.Context) (uint64, error) {
	store := k.storeService.OpenKVStore(*ctx)

	key := utils.StrBytes(types.DidDocCountKey)
	valueBytes, err := store.Get(key)
	if err != nil {
		return 0, err
	}

	// Count doesn't exist: no element
	if valueBytes == nil {
		return 0, nil
	}

	// Parse bytes
	count, err := strconv.ParseUint(string(valueBytes), 10, 64)
	if err != nil {
		// Panic because the count should be always formattable to iint64
		panic("cannot decode count")
	}

	return count, nil
}

// SetDidCount set the total number of did
func (k Keeper) SetDidDocCount(ctx *context.Context, count uint64) error {
	store := k.storeService.OpenKVStore(*ctx)

	key := utils.StrBytes(types.DidDocCountKey)
	valueBytes := []byte(strconv.FormatUint(count, 10))

	return store.Set(key, valueBytes)
}

func (k Keeper) AddNewDidDocVersion(ctx *context.Context, didDoc *types.DidDocWithMetadata) error {
	// Check if the diddoc version already exists
	hasDidDocVersion, err := k.HasDidDocVersion(ctx, didDoc.DidDoc.Id, didDoc.Metadata.VersionId)
	if err != nil {
		return err
	}
	if hasDidDocVersion {
		return types.ErrDidDocExists.Wrapf("diddoc version already exists for did %s, version %s", didDoc.DidDoc.Id, didDoc.Metadata.VersionId)
	}

	// Link to the previous version if it exists
	hasDidDoc, err := k.HasDidDoc(ctx, didDoc.DidDoc.Id)
	if hasDidDoc {
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
	err = k.SetLatestDidDocVersion(ctx, didDoc.DidDoc.Id, didDoc.Metadata.VersionId)
	if err != nil {
		return err
	}

	// Write new version (no override)
	return k.SetDidDocVersion(ctx, didDoc, false)
}

func (k Keeper) GetLatestDidDoc(ctx *context.Context, did string) (types.DidDocWithMetadata, error) {
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
func (k Keeper) SetDidDocVersion(ctx *context.Context, value *types.DidDocWithMetadata, override bool) error {
	hasdidVersion, err := k.HasDidDocVersion(ctx, value.DidDoc.Id, value.Metadata.VersionId)
	if err != nil {
		return err
	}
	if !override && hasdidVersion {
		return types.ErrDidDocExists.Wrap("diddoc version already exists")
	}

	// Create the diddoc version
	store := k.storeService.OpenKVStore(*ctx)

	key := types.GetDidDocVersionKey(value.DidDoc.Id, value.Metadata.VersionId)
	valueBytes := k.cdc.MustMarshal(value)
	store.Set(key, valueBytes)

	return nil
}

// GetDid returns a did from its id
func (k Keeper) GetDidDocVersion(ctx *context.Context, id, version string) (types.DidDocWithMetadata, error) {
	store := k.storeService.OpenKVStore(*ctx)
	hasdidVersion, err := k.HasDidDocVersion(ctx, id, version)
	if err != nil {
		return types.DidDocWithMetadata{}, err
	}
	if !hasdidVersion {
		return types.DidDocWithMetadata{}, sdkerrors.ErrNotFound.Wrap("diddoc version not found")
	}

	var value types.DidDocWithMetadata
	valueBytes, err := store.Get(types.GetDidDocVersionKey(id, version))
	if err != nil {
		return types.DidDocWithMetadata{}, err
	}
	k.cdc.MustUnmarshal(valueBytes, &value)

	return value, nil
}

func (k Keeper) GetAllDidDocVersions(ctx *context.Context, did string) ([]*types.Metadata, error) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(*ctx))

	result := make([]*types.Metadata, 0)

	versionIterator := storetypes.KVStorePrefixIterator(store, types.GetDidDocVersionsPrefix(did))
	defer closeIteratorOrPanic(versionIterator)

	for ; versionIterator.Valid(); versionIterator.Next() {
		// Get the diddoc
		var didDoc types.DidDocWithMetadata
		k.cdc.MustUnmarshal(versionIterator.Value(), &didDoc)

		result = append(result, didDoc.Metadata)
	}

	return result, nil
}

// SetLatestDidDocVersion sets the latest version id value for a diddoc
func (k Keeper) SetLatestDidDocVersion(ctx *context.Context, did, version string) error {
	// Update counter. We use latest version as existence indicator.
	hasVersion, err := k.HasLatestDidDocVersion(ctx, did)
	if err != nil {
		return err
	}
	if !hasVersion {
		count, err := k.GetDidDocCount(ctx)
		if err != nil {
			return err
		}
		k.SetDidDocCount(ctx, count+1)
	}

	store := k.storeService.OpenKVStore(*ctx)
	key := types.GetLatestDidDocVersionKey(did)
	valueBytes := utils.StrBytes(version)
	return store.Set(key, valueBytes)
}

// GetLatestDidDocVersion returns the latest version id value for a diddoc
func (k Keeper) GetLatestDidDocVersion(ctx *context.Context, id string) (string, error) {
	store := k.storeService.OpenKVStore(*ctx)
	hasVersion, err := k.HasLatestDidDocVersion(ctx, id)
	if !hasVersion {
		return "", sdkerrors.ErrNotFound.Wrap(id)
	}
	value, err := store.Get(types.GetLatestDidDocVersionKey(id))
	if err != nil {
		return "", err
	}
	return string(value), nil
}

func (k Keeper) HasDidDoc(ctx *context.Context, id string) (bool, error) {
	return k.HasLatestDidDocVersion(ctx, id)
}

func (k Keeper) HasLatestDidDocVersion(ctx *context.Context, id string) (bool, error) {
	store := k.storeService.OpenKVStore(*ctx)
	return store.Has(types.GetLatestDidDocVersionKey(id))
}

func (k Keeper) HasDidDocVersion(ctx *context.Context, id, version string) (bool, error) {
	store := k.storeService.OpenKVStore(*ctx)
	return store.Has(types.GetDidDocVersionKey(id, version))
}

func (k Keeper) IterateDids(ctx *context.Context, callback func(did string) (continue_ bool)) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(*ctx))
	latestVersionIterator := storetypes.KVStorePrefixIterator(store, types.GetLatestDidDocVersionPrefix())
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

func (k Keeper) IterateDidDocVersions(ctx *context.Context, did string, callback func(version types.DidDocWithMetadata) (continue_ bool)) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(*ctx))
	versionIterator := storetypes.KVStorePrefixIterator(store, types.GetDidDocVersionsPrefix(did))
	defer closeIteratorOrPanic(versionIterator)

	for ; versionIterator.Valid(); versionIterator.Next() {
		var didDoc types.DidDocWithMetadata
		k.cdc.MustUnmarshal(versionIterator.Value(), &didDoc)

		if !callback(didDoc) {
			break
		}
	}
}

func (k Keeper) IterateAllDidDocVersions(ctx *context.Context, callback func(version types.DidDocWithMetadata) (continue_ bool)) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(*ctx))
	allVersionsIterator := storetypes.KVStorePrefixIterator(store, []byte(types.DidDocVersionKey))
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
func (k Keeper) GetAllDidDocs(ctx *context.Context) ([]*types.DidDocVersionSet, error) {
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

func closeIteratorOrPanic(iterator storetypes.Iterator) {
	err := iterator.Close()
	if err != nil {
		panic(err.Error())
	}
}
