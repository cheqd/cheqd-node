package migration

import (
	"fmt"
	"path/filepath"

	. "github.com/onsi/gomega"

	migrationsetup "github.com/cheqd/cheqd-node/tests/upgrade/migration/setup"

	// didtestssetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	didtypesv1 "github.com/cheqd/cheqd-node/x/did/types/v1"

	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	resourcetypesv1 "github.com/cheqd/cheqd-node/x/resource/types/v1"
)

type ProtobufDataSet struct {
	setup             migrationsetup.TestSetup
	existingDIDDocs   []didtypesv1.StateValue
	existingResources []resourcetypesv1.Resource
	expectedDidDocs   []didtypes.DidDocWithMetadata
	expectedResources []resourcetypes.ResourceWithMetadata
}

func NewProtobufDataSet(setup migrationsetup.TestSetup) ProtobufDataSet {
	return ProtobufDataSet{
		setup:             setup,
		existingDIDDocs:   []didtypesv1.StateValue{},
		existingResources: []resourcetypesv1.Resource{},
		expectedDidDocs:   []didtypes.DidDocWithMetadata{},
		expectedResources: []resourcetypes.ResourceWithMetadata{},
	}
}

func (pds *ProtobufDataSet) Load() error {
	var (
		existingProtobufDidDoc   didtypesv1.StateValue
		existingProtobufResource resourcetypesv1.Resource

		expectedProtobufDidDoc   didtypes.DidDocWithMetadata
		expectedProtobufResource resourcetypes.ResourceWithMetadata
	)
	// Get DIDDoc v1
	err := Loader(
		filepath.Join("payload", "existing", "v1", "protobuf", "diddoc.json"),
		&existingProtobufDidDoc,
		pds.setup)
	if err != nil {
		fmt.Println("Error loading didDoc")
		return err
	}
	// Get resource v1
	err = Loader(
		filepath.Join("payload", "existing", "v1", "protobuf", "resource.json"),
		&existingProtobufResource,
		pds.setup)
	if err != nil {
		fmt.Println("Error loading resource v1")
		return err
	}
	// Get DIDDoc v2
	err = Loader(
		filepath.Join("payload", "expected", "v2", "protobuf", "diddoc.json"),
		&expectedProtobufDidDoc,
		pds.setup)
	if err != nil {
		fmt.Println("Error loading didDoc v2")
		return err
	}

	// Get resource v2
	err = Loader(
		filepath.Join("payload", "expected", "v2", "protobuf", "resource.json"),
		&expectedProtobufResource,
		pds.setup)
	if err != nil {
		fmt.Println("Error loading expectedChecksumResource")
		return err
	}

	pds.existingDIDDocs = append(pds.existingDIDDocs, existingProtobufDidDoc)
	pds.existingResources = append(pds.existingResources, existingProtobufResource)

	pds.expectedDidDocs = append(pds.expectedDidDocs, expectedProtobufDidDoc)
	pds.expectedResources = append(pds.expectedResources, expectedProtobufResource)
	return nil
}

func (pds *ProtobufDataSet) Prepare() error {
	for _, resource := range pds.existingResources {
		err := pds.setup.ResourceKeeperV1.SetResource(&pds.setup.SdkCtx, &resource)
		if err != nil {
			return err
		}
	}
	for _, didDoc := range pds.existingDIDDocs {
		err := pds.setup.DidKeeperV1.SetDid(&pds.setup.SdkCtx, &didDoc)
		if err != nil {
			return err
		}
	}
	return nil
}

func (pds *ProtobufDataSet) Validate() error {
	var (
		expectedDidDoc  didtypes.DidDocWithMetadata
		expectedResource resourcetypes.ResourceWithMetadata
	)
	for _, expectedDidDoc = range pds.expectedDidDocs {
		didDoc, err := pds.setup.DidKeeper.GetLatestDidDoc(&pds.setup.SdkCtx, expectedDidDoc.DidDoc.Id)
		if err != nil {
			return err
		}
		Expect(didDoc.DidDoc.Id).To(Equal(expectedDidDoc.DidDoc.Id))
		Expect(didDoc.DidDoc.Context).To(Equal(expectedDidDoc.DidDoc.Context))
		Expect(didDoc.DidDoc.Service).To(Equal(expectedDidDoc.DidDoc.Service))
		Expect(didDoc.DidDoc.VerificationMethod).To(Equal(expectedDidDoc.DidDoc.VerificationMethod))
		Expect(didDoc.DidDoc.Authentication).To(Equal(expectedDidDoc.DidDoc.Authentication))
		Expect(didDoc.DidDoc.AssertionMethod).To(Equal(expectedDidDoc.DidDoc.AssertionMethod))
		Expect(didDoc.DidDoc.CapabilityInvocation).To(Equal(expectedDidDoc.DidDoc.CapabilityInvocation))
		Expect(didDoc.DidDoc.CapabilityDelegation).To(Equal(expectedDidDoc.DidDoc.CapabilityDelegation))
		Expect(didDoc.DidDoc.KeyAgreement).To(Equal(expectedDidDoc.DidDoc.KeyAgreement))
		Expect(didDoc.DidDoc.Service).To(Equal(expectedDidDoc.DidDoc.Service))
		Expect(didDoc.DidDoc.AlsoKnownAs).To(Equal(expectedDidDoc.DidDoc.AlsoKnownAs))
		Expect(didDoc.Metadata).To(Equal(expectedDidDoc.Metadata))
	}

	for _, expectedResource = range pds.expectedResources {
		resource, err := pds.setup.ResourceKeeper.GetResource(&pds.setup.SdkCtx, 
			expectedResource.Metadata.CollectionId,
			expectedResource.Metadata.Id)
		if err != nil {
			return err
		}
		Expect(resource.Metadata.Id).To(Equal(expectedResource.Metadata.Id))
		Expect(resource.Metadata.CollectionId).To(Equal(expectedResource.Metadata.CollectionId))
		Expect(resource.Metadata.Name).To(Equal(expectedResource.Metadata.Name))
		Expect(resource.Metadata.Version).To(Equal(expectedResource.Metadata.Version))
		Expect(resource.Metadata.ResourceType).To(Equal(expectedResource.Metadata.ResourceType))
		Expect(resource.Metadata.AlsoKnownAs).To(Equal(expectedResource.Metadata.AlsoKnownAs))
		Expect(resource.Metadata.MediaType).To(Equal(expectedResource.Metadata.MediaType))
		Expect(resource.Metadata.Created).To(Equal(expectedResource.Metadata.Created))
		Expect(resource.Metadata.Checksum).To(Equal(expectedResource.Metadata.Checksum))
	}

	return nil
}
