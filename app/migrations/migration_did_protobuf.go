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
	sctx.Logger().Debug("MigrateDidProtobuf: Starting migration")

	codec := NewLegacyProtoCodec()
	store := sctx.KVStore(mctx.didStoreKey)

	sctx.Logger().Debug("MigrateDidProtobuf: Erasing old count key")
	// Erase old broken count key
	store.Delete([]byte(didtypesv1.DidCountKey + didtypesv1.DidCountKey))

	sctx.Logger().Debug("MigrateResourceProtobuf: Reading all keys")
	didKeys := helpers.ReadAllKeys(store, didutils.StrBytes(didtypesv1.DidKey))

	for _, didKey := range didKeys {
		sctx.Logger().Debug("MigrateDidProtobuf: Starting migration for didKey: " + string(didKey))

		var stateValue didtypesv1.StateValue
		sctx.Logger().Debug("MigrateDidProtobuf: Reading StateValue of DID from store")
		codec.MustUnmarshal(store.Get(didKey), &stateValue)

		sctx.Logger().Debug("MigrateDidProtobuf: Migrating StateValue for DidDocWithMetadata")
		newDidDocWithMetadata, err := MigrateStateValue(sctx, mctx, &stateValue)
		if err != nil {
			return err
		}

		// Remove old DID Doc
		store.Delete(didKey)

		sctx.Logger().Debug("MigrateDidProtobuf: Setting DidDocWithMetadata to store")
		// Set new DID Doc
		err = mctx.didKeeperNew.AddNewDidDocVersion(&sctx, &newDidDocWithMetadata)
		if err != nil {
			return err
		}
		sctx.Logger().Debug("MigrateDidProtobuf: Migration finished for didKey: " + string(didKey))
	}

	// Migrate DID namespace (at least make sure it's not changed)
	if didtypesv1.DidNamespaceKey != didtypes.DidNamespaceKey {
		panic("DID namespace key is changed")
	}
	sctx.Logger().Debug("MigrateDidProtobuf: Migration finished")

	return nil
}

func NewLegacyProtoCodec() *codec.ProtoCodec {
	ir := codectypes.NewInterfaceRegistry()

	ir.RegisterInterface("StateValueData", (*didtypesv1.StateValueData)(nil))
	ir.RegisterImplementations((*didtypesv1.StateValueData)(nil), &didtypesv1.Did{})

	return codec.NewProtoCodec(ir)
}

func MigrateStateValue(sctx sdk.Context, mctx MigrationContext, stateValue *didtypesv1.StateValue) (didtypes.DidDocWithMetadata, error) {
	oldDidDoc, err := stateValue.UnpackDataAsDid()
	if err != nil {
		return didtypes.DidDocWithMetadata{}, err
	}

	oldMetadata := stateValue.Metadata

	sctx.Logger().Debug("MigrateDidProtobuf: OldMetadata: " + string(mctx.codec.MustMarshalJSON(oldMetadata)))
	newDidDoc := MigrateDidDoc(oldDidDoc)
	newMetadata := MigrateMetadata(oldMetadata)
	sctx.Logger().Debug("MigrateDidProtobuf: NewMetadata: " + string(mctx.codec.MustMarshalJSON(&newMetadata)))

	return didtypes.DidDocWithMetadata{
		DidDoc:   &newDidDoc,
		Metadata: &newMetadata,
	}, nil
}

func MigrateMetadata(metadata *didtypesv1.Metadata) didtypes.Metadata {
	updated := helpers.MustParseFromStringTimeToGoTime(metadata.Updated)
	return didtypes.Metadata{
		Created:           helpers.MustParseFromStringTimeToGoTime(metadata.Created),
		Updated:           &updated,
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
				Id:                     vm.Id,
				VerificationMethodType: MigrateType(vm.Type),
				Controller:             vm.Controller,
				VerificationMaterial:   MigrateVerificationMaterial(vm),
			})
	}

	srvs := []*didtypes.Service{}
	for _, srv := range oldDid.Service {
		srvs = append(
			srvs,
			&didtypes.Service{
				Id:              srv.Id,
				ServiceType:     srv.Type,
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
		return didtypes.Ed25519VerificationKey2020Type
	case didtypesv1.JSONWebKey2020:
		return didtypes.JSONWebKey2020Type
	default:
		panic("Unknown type")
	}
}

func MigrateVerificationMaterial(vm *didtypesv1.VerificationMethod) string {
	switch vm.Type {
	case didtypesv1.JSONWebKey2020:
		jwk := make(map[string]string)
		for _, kv := range vm.PublicKeyJwk {
			jwk[kv.Key] = kv.Value
		}
		res, err := json.Marshal(jwk)
		if err != nil {
			panic(err)
		}
		return string(res)

	case didtypesv1.Ed25519VerificationKey2020:
		pkMulti, err := helpers.GenerateEd25519VerificationKey2020VerificationMaterial(vm.PublicKeyMultibase)
		if err != nil {
			panic(err)
		}

		return pkMulti

	default:
		panic("Unknown type")
	}
}
