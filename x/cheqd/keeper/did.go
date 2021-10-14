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
	did := types.Did{
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

	created := ctx.BlockTime().String()
	txHash := base64.StdEncoding.EncodeToString(tmhash.Sum(ctx.TxBytes()))

	metadata := types.Metadata{
		Created:     created,
		Updated:     created,
		Deactivated: false,
		VersionId:   txHash,
	}

	stateValue := types.StateValue{
		Metadata: &metadata,
		Data: &types.StateValue_Did{
			Did: &did,
		},
		Timestamp: created,
		TxHash:    txHash,
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DidKey))
	value := k.cdc.MustMarshal(&stateValue)
	store.Set(GetDidIDBytes(did.Id), value)

	// Update did count
	k.SetDidCount(ctx, count+1)

	return id
}

// SetDid set a specific did in the store
func (k Keeper) SetDid(ctx sdk.Context, did types.Did, metadata *types.Metadata) {
	updated := ctx.BlockTime().String()
	txHash := base64.StdEncoding.EncodeToString(tmhash.Sum(ctx.TxBytes()))

	metadata = &types.Metadata{
		Created:     metadata.Created,
		Updated:     updated,
		Deactivated: metadata.Deactivated,
		VersionId:   txHash,
	}

	stateValue := types.StateValue{
		Metadata: metadata,
		Data: &types.StateValue_Did{
			Did: &did,
		},
		Timestamp: updated,
		TxHash:    txHash,
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DidKey))
	b := k.cdc.MustMarshal(&stateValue)
	store.Set(GetDidIDBytes(did.Id), b)
}

// GetDid returns a did from its id
func (k Keeper) GetDid(ctx *sdk.Context, id string) (*types.Did, *types.Metadata, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DidKey))

	var value types.StateValue
	err := k.cdc.Unmarshal(store.Get(GetDidIDBytes(id)), &value)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(sdkerrors.ErrInvalidType, err.Error())
	}

	switch data := value.Data.(type) {
	case *types.StateValue_Did:
		return data.Did, value.Metadata, nil
	default:
		return nil, nil, sdkerrors.Wrap(sdkerrors.ErrInvalidType, fmt.Sprintf("State has unexpected type %T", data))
	}
}

// areOwners returns a bool are received authors can control this DID
func (k Keeper) areDidOwners(ctx sdk.Context, id string, authors []string) bool {
	//return utils.CompareOwners(authors, k.GetDid(ctx, id).Controller)
	return true
}

// HasDid checks if the did exists in the store
func (k Keeper) HasDid(ctx sdk.Context, id string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DidKey))
	return store.Has(GetDidIDBytes(id))
}

// GetDidIDBytes returns the byte representation of the ID
func GetDidIDBytes(id string) []byte {
	return []byte(id)
}
