package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"github.com/cheqd/cheqd-node/x/did/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// GetDidCount get the total number of did
func (k Keeper) GetDidDocCount(ctx context.Context) (uint64, error) {
	has, err := k.DidCount.Has(ctx)
	if err != nil {
		return 0, err
	}

	if !has {
		if setErr := k.DidCount.Set(ctx, 0); setErr != nil {
			return 0, setErr
		}
	}
	count, err := k.DidCount.Get(ctx)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SetDidCount set the total number of did
func (k Keeper) SetDidDocCount(ctx context.Context, count uint64) error {
	return k.DidCount.Set(ctx, count)
}

func (k Keeper) AddNewDidDocVersion(ctx context.Context, didDoc *types.DidDocWithMetadata) error {
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
	if err != nil {
		return err
	}
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

func (k Keeper) GetLatestDidDoc(ctx context.Context, did string) (types.DidDocWithMetadata, error) {
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
func (k Keeper) SetDidDocVersion(ctx context.Context, value *types.DidDocWithMetadata, override bool) error {
	if !override {
		hasdidVersion, err := k.HasDidDocVersion(ctx, value.DidDoc.Id, value.Metadata.VersionId)
		if err != nil {
			return err
		}

		if hasdidVersion {
			return types.ErrDidDocExists.Wrap("diddoc version already exists")
		}
	}

	// Create the diddoc version
	return k.DidDocuments.Set(ctx, collections.Join(value.DidDoc.Id, value.Metadata.VersionId), *value)
}

// GetDid returns a did from its id
func (k Keeper) GetDidDocVersion(ctx context.Context, id, version string) (types.DidDocWithMetadata, error) {
	hasdidVersion, err := k.HasDidDocVersion(ctx, id, version)
	if err != nil {
		return types.DidDocWithMetadata{}, err
	}
	if !hasdidVersion {
		return types.DidDocWithMetadata{}, sdkerrors.ErrNotFound.Wrap("diddoc version not found")
	}

	return k.DidDocuments.Get(ctx, collections.Join(id, version))
}

func (k Keeper) GetAllDidDocVersions(ctx context.Context, did string) ([]*types.Metadata, error) {
	var metadataList []*types.Metadata
	rng := collections.NewPrefixedPairRange[string, string](did)

	iter, err := k.DidDocuments.Iterate(ctx, rng)
	if err != nil {
		return nil, err
	}

	kvs, err := iter.KeyValues()
	if err != nil {
		return nil, err
	}
	for _, kv := range kvs {
		metadataList = append(metadataList, kv.Value.Metadata)
	}
	return metadataList, nil
}

// SetLatestDidDocVersion sets the latest version id value for a diddoc
func (k Keeper) SetLatestDidDocVersion(ctx context.Context, did, version string) error {
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
		err = k.SetDidDocCount(ctx, count+1)
		if err != nil {
			return err
		}
	}

	return k.LatestDidVersion.Set(ctx, did, version)
}

// GetLatestDidDocVersion returns the latest version id value for a diddoc
func (k Keeper) GetLatestDidDocVersion(ctx context.Context, id string) (string, error) {
	hasVersion, err := k.HasLatestDidDocVersion(ctx, id)
	if err != nil {
		return "", err
	}
	if !hasVersion {
		return "", sdkerrors.ErrNotFound.Wrap(id)
	}
	value, err := k.LatestDidVersion.Get(ctx, id)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (k Keeper) HasDidDoc(ctx context.Context, id string) (bool, error) {
	return k.HasLatestDidDocVersion(ctx, id)
}

func (k Keeper) HasLatestDidDocVersion(ctx context.Context, id string) (bool, error) {
	return k.LatestDidVersion.Has(ctx, id)
}

func (k Keeper) HasDidDocVersion(ctx context.Context, id, version string) (bool, error) {
	return k.DidDocuments.Has(ctx, collections.Join(id, version))
}

func (k Keeper) IterateDids(ctx context.Context, callback func(did string) (continue_ bool)) {
	iter, err := k.LatestDidVersion.Iterate(ctx, nil)
	if err != nil {
		panic(err)
	}

	kvs, err := iter.KeyValues()
	if err != nil {
		panic(err)
	}

	for _, kv := range kvs {
		if !callback(kv.Key) {
			break
		}
	}
}

func (k Keeper) IterateDidDocVersions(ctx context.Context, did string, callback func(version types.DidDocWithMetadata) (continue_ bool)) {
	rng := collections.NewPrefixedPairRange[string, string](did)

	iter, err := k.DidDocuments.Iterate(ctx, rng)
	if err != nil {
		panic(err)
	}

	kvs, err := iter.KeyValues()
	if err != nil {
		panic(err)
	}

	for _, kv := range kvs {
		if !callback(kv.Value) {
			break
		}
	}
}

func (k Keeper) IterateAllDidDocVersions(ctx context.Context, callback func(version types.DidDocWithMetadata) (continue_ bool)) {
	iter, err := k.DidDocuments.Iterate(ctx, nil)
	if err != nil {
		panic(err)
	}

	kvs, err := iter.KeyValues()
	if err != nil {
		panic(err)
	}

	for _, kv := range kvs {
		if !callback(kv.Value) {
			break
		}
	}
}

// GetAllDidDocs returns all did
// Loads all DIDs in memory. Use only for genesis export.
func (k Keeper) GetAllDidDocs(ctx context.Context) ([]*types.DidDocVersionSet, error) {
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
