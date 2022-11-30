package setup

import (
	resourcetestssetup "github.com/cheqd/cheqd-node/x/resource/tests/setup"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
)

type DidKeeperV1 struct {
	// nolint: structcheck, unused
	cdc codec.BinaryCodec
	// nolint: structcheck, unused
	storeKey storetypes.KVStoreKey
}

type DidMsgServerV1 struct {
	keeper DidKeeperV1
}

type DidQueryServerV1 struct {
	keeper DidKeeperV1
}

type ResourceKeeperV1 struct {
	// nolint: structcheck, unused
	cdc codec.BinaryCodec
	// nolint: structcheck, unused
	storeKey storetypes.KVStoreKey
}

type ResourceMsgServerV1 struct {
	keeper ResourceKeeperV1
}

type ResourceQueryServerV1 struct {
	keeper ResourceKeeperV1
}

type ExtendedTestSetup struct {
	resourcetestssetup.TestSetup
	didKeeperV1           DidKeeperV1           // TODO: replace with actual type implementation
	didMsgServerV1        DidMsgServerV1        // TODO: replace with actual type implementation
	didQueryServerV1      DidQueryServerV1      // TODO: replace with actual type implementation
	resourceKeeperV1      ResourceKeeperV1      // TODO: replace with actual type implementation
	resourceMsgServerV1   ResourceMsgServerV1   // TODO: replace with actual type implementation
	resourceQueryServerV1 ResourceQueryServerV1 // TODO: replace with actual type implementation
}

func NewExtendedSetup() ExtendedTestSetup {
	setup := resourcetestssetup.Setup()
	didKeeperV1 := DidKeeperV1{}
	resourceKeeperV1 := ResourceKeeperV1{}
	return ExtendedTestSetup{
		TestSetup:             setup,
		didKeeperV1:           didKeeperV1,
		didMsgServerV1:        DidMsgServerV1{keeper: didKeeperV1},
		didQueryServerV1:      DidQueryServerV1{keeper: didKeeperV1},
		resourceKeeperV1:      resourceKeeperV1,
		resourceMsgServerV1:   ResourceMsgServerV1{keeper: resourceKeeperV1},
		resourceQueryServerV1: ResourceQueryServerV1{keeper: resourceKeeperV1},
	}
}
