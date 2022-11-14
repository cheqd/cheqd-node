package setup

import (
	"fmt"
	"strconv"

	didtypesv2 "github.com/cheqd/cheqd-node/x/did/types"
	didtypesv1 "github.com/cheqd/cheqd-node/x/did/types/v1"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/libs/log"
)

type KeeperV1 struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
}

func NewKeeperV1(cdc codec.BinaryCodec, storeKey storetypes.StoreKey) *KeeperV1 {
	return &KeeperV1{
		cdc:      cdc,
		storeKey: storeKey,
	}
}

func (k KeeperV1) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", didtypesv2.ModuleName))
}

// GetDidNamespace get did namespace
func (k KeeperV1) GetDidNamespace(ctx *sdk.Context) string {
	return k.GetFromState(ctx, didtypesv2.DidNamespaceKey)
}

// GetFromState - get State value
func (k KeeperV1) GetFromState(ctx *sdk.Context, stateKey string) string {
	store := ctx.KVStore(k.storeKey)
	byteKey := []byte(stateKey)
	bz := store.Get(byteKey)

	// Parse bytes
	namespace := string(bz)
	return namespace
}

// GetDidCount get the total number of did
func (k KeeperV1) GetDidCount(ctx *sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(didtypesv2.DidCountKey))
	byteKey := []byte(didtypesv2.DidCountKey)
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
func (k KeeperV1) SetDidCount(ctx *sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(didtypesv2.DidCountKey))
	byteKey := []byte(didtypesv2.DidCountKey)
	bz := []byte(strconv.FormatUint(count, 10))
	store.Set(byteKey, bz)
}

// GetDid returns a did from its id
func (k KeeperV1) GetDid(ctx *sdk.Context, id string) (didtypesv1.StateValue, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(didtypesv2.DidKey))

	if !k.HasDid(ctx, id) {
		return didtypesv1.StateValue{}, sdkerrors.ErrNotFound.Wrap(id)
	}

	var value didtypesv1.StateValue
	bytes := store.Get(GetDidIDBytes(id))
	if err := k.cdc.Unmarshal(bytes, &value); err != nil {
		return didtypesv1.StateValue{}, sdkerrors.Wrap(sdkerrors.ErrInvalidType, err.Error())
	}

	return value, nil
}

// SetDid set a specific did in the store. Updates DID counter if the DID is new.
func (k KeeperV1) SetDid(ctx *sdk.Context, stateValue *didtypesv1.StateValue) error {
	// Unpack
	did, err := stateValue.UnpackDataAsDid()
	if err != nil {
		return err
	}

	// Update counter
	if !k.HasDid(ctx, did.Id) {
		count := k.GetDidCount(ctx)
		k.SetDidCount(ctx, count+1)
	}

	// Create the did
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(didtypesv2.DidKey))
	b := k.cdc.MustMarshal(stateValue)
	store.Set(GetDidIDBytes(did.Id), b)
	return nil
}

// HasDid checks if the did exists in the store
func (k KeeperV1) HasDid(ctx *sdk.Context, id string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(didtypesv1.DidKey))
	return store.Has(GetDidIDBytes(id))
}

// GetDidIDBytes returns the byte representation of the ID
func GetDidIDBytes(id string) []byte {
	return []byte(id)
}
