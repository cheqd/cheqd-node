package migrations

import (
	"encoding/json"

	"github.com/cheqd/cheqd-node/app/migrations/helpers"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	didtypesv1 "github.com/cheqd/cheqd-node/x/did/types/v1"
	didutils "github.com/cheqd/cheqd-node/x/did/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MigrateDidProtobuf(sctx sdk.Context, mctx MigrationContext) error {
	println("Protobuf migration for dids. Start")
	codec := NewLegacyProtoCodec()
	store := sctx.KVStore(mctx.didStoreKey)

	// Erase old broken count key
	store.Delete([]byte(didtypesv1.DidCountKey + didtypesv1.DidCountKey))

	println("Protobuf migration for dids. Read all keys")
	didKeys := helpers.ReadAllKeys(store, didutils.StrBytes(didtypesv1.DidKey))

	for _, didKey := range didKeys {
		var stateValue didtypesv1.StateValue
		codec.MustUnmarshal(store.Get(didKey), &stateValue)

		newDidDocWithMetadata, err := MigrateStateValue(&stateValue)
		if err != nil {
			return err
		}

		// Remove old DID Doc
		store.Delete(didKey)

		// Set new DID Doc
		err = mctx.didKeeperNew.AddNewDidDocVersion(&sctx, &newDidDocWithMetadata)
		if err != nil {
			return err
		}
	}

	// Migrate DID namespace (at least make sure it's not changed)
	if didtypesv1.DidNamespaceKey != didtypes.DidNamespaceKey {
		panic("DID namespace key is changed")
	}
	println("Protobuf migration for dids. End")

	return nil
}

func NewLegacyProtoCodec() *codec.ProtoCodec {
	ir := codectypes.NewInterfaceRegistry()

	ir.RegisterInterface("StateValueData", (*didtypesv1.StateValueData)(nil))
	ir.RegisterImplementations((*didtypesv1.StateValueData)(nil), &didtypesv1.Did{})

	return codec.NewProtoCodec(ir)
}

func MigrateStateValue(stateValue *didtypesv1.StateValue) (didtypes.DidDocWithMetadata, error) {
	oldDidDoc, err := stateValue.UnpackDataAsDid()
	if err != nil {
		return didtypes.DidDocWithMetadata{}, err
	}

	oldMetadata := stateValue.Metadata

	newDidDoc := MigrateDidDoc(oldDidDoc)
	newMetadata := MigrateMetadata(oldMetadata)

	return didtypes.DidDocWithMetadata{
		DidDoc:   &newDidDoc,
		Metadata: &newMetadata,
	}, nil
}

func MigrateMetadata(metadata *didtypesv1.Metadata) didtypes.Metadata {
	return didtypes.Metadata{
		Created:           metadata.Created,
		Updated:           metadata.Updated,
		Deactivated:       metadata.Deactivated,
		VersionId:         metadata.VersionId,
		NextVersionId:     "",
		PreviousVersionId: "",
	}
}

func MigrateDidDoc(oldDid *didtypesv1.Did) didtypes.DidDoc {
	vms := []*didtypes.VerificationMethod{}
	for _, vm := range oldDid.VerificationMethod {
		vms = append(
			vms,
			&didtypes.VerificationMethod{
				Id:                   vm.Id,
				Type:                 MigrateType(vm.Type),
				Controller:           vm.Controller,
				VerificationMaterial: MigrateVerificationMaterial(vm),
			})
	}

	srvs := []*didtypes.Service{}
	for _, srv := range oldDid.Service {
		srvs = append(
			srvs,
			&didtypes.Service{
				Id:              srv.Id,
				Type:            srv.Type,
				ServiceEndpoint: []string{srv.ServiceEndpoint},
			})
	}

	return didtypes.DidDoc{
		Id:                   oldDid.Id,
		VerificationMethod:   vms,
		Service:              srvs,
		Context:              oldDid.Context,
		Controller:           oldDid.Controller,
		Authentication:       oldDid.Authentication,
		AssertionMethod:      oldDid.AssertionMethod,
		CapabilityDelegation: oldDid.CapabilityDelegation,
		CapabilityInvocation: oldDid.CapabilityInvocation,
		KeyAgreement:         oldDid.KeyAgreement,
		AlsoKnownAs:          oldDid.AlsoKnownAs,
	}
}

func MigrateType(t string) string {
	switch t {
	case didtypesv1.Ed25519VerificationKey2020:
		return didtypes.Ed25519VerificationKey2020{}.Type()
	case didtypesv1.JsonWebKey2020:
		return didtypes.JsonWebKey2020{}.Type()
	default:
		panic("Unknown type")
	}
}

func MigrateVerificationMaterial(vm *didtypesv1.VerificationMethod) string {
	switch vm.Type {
	case didtypesv1.JsonWebKey2020:
		jwk := make(map[string]string)
		for _, kv := range vm.PublicKeyJwk {
			jwk[kv.Key] = kv.Value
		}
		res, err := json.Marshal(jwk)
		if err != nil {
			panic(err)
		}

		jwk2020 := didtypes.JsonWebKey2020{
			PublicKeyJwk: res,
		}
		res, err = json.Marshal(jwk2020)
		if err != nil {
			panic(err)
		}

		return string(res)

	case didtypesv1.Ed25519VerificationKey2020:
		pk_multi := didtypes.Ed25519VerificationKey2020{
			PublicKeyMultibase: vm.PublicKeyMultibase,
		}

		res, err := json.Marshal(pk_multi)
		if err != nil {
			panic(err)
		}

		return string(res)

	default:
		panic("Unknown type")
	}
}
