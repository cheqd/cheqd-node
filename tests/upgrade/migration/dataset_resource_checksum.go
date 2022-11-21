package migration

import (
	"bytes"
	"fmt"
	"path/filepath"

	migrationsetup "github.com/cheqd/cheqd-node/tests/upgrade/migration/setup"

	didtypes "github.com/cheqd/cheqd-node/x/did/types"

	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
)

type ChecksumDataSet struct {
	setup migrationsetup.TestSetup
	existingDIDDocs   []didtypes.DidDocWithMetadata
	existingResources []resourcetypes.ResourceWithMetadata
	expectedDidDocs   []didtypes.DidDocWithMetadata
	expectedResources []resourcetypes.ResourceWithMetadata
}

func NewChecksumDataSet(setup migrationsetup.TestSetup) ChecksumDataSet {
	return ChecksumDataSet{
		setup: setup,
		existingDIDDocs: []didtypes.DidDocWithMetadata{},
		existingResources: []resourcetypes.ResourceWithMetadata{},
		expectedDidDocs: []didtypes.DidDocWithMetadata{},
		expectedResources: []resourcetypes.ResourceWithMetadata{},
	}
}

func (cds *ChecksumDataSet) Load() error {
	var (
		existingChecksumResource resourcetypes.ResourceWithMetadata
		expectedChecksumResource resourcetypes.ResourceWithMetadata
	)
	err := Loader(
		filepath.Join("payload", "existing", "v2", "checksum", "resource.json"),
		&existingChecksumResource,
		cds.setup)
	if err != nil {
		fmt.Println("Error loading existingChecksumResource")
		return err
	}
	err = Loader(
		filepath.Join("payload", "expected", "v2", "checksum", "resource.json"),
		&expectedChecksumResource,
		cds.setup)
	if err != nil {
		fmt.Println("Error loading expectedChecksumResource")
		return err
	}
	cds.existingResources = append(cds.existingResources, existingChecksumResource)
	cds.expectedResources = append(cds.expectedResources, expectedChecksumResource)
	return nil
}

func (cds *ChecksumDataSet) Prepare() error {
	for _, resource := range cds.existingResources {
		err := cds.setup.ResourceKeeper.SetResource(&cds.setup.SdkCtx, &resource)
		if err != nil {
			return err
		}
	}
	for _, didDoc := range cds.existingDIDDocs {
		err := cds.setup.DidKeeper.SetDidDocVersion(&cds.setup.SdkCtx, &didDoc, false)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cds *ChecksumDataSet) Validate() error {
	var (
		resourceWithMetadata resourcetypes.ResourceWithMetadata
		err error
	)

	for _, expectedResource := range cds.expectedResources {
		resourceWithMetadata, err = cds.setup.ResourceKeeper.GetResource(
			&cds.setup.SdkCtx,
			expectedResource.Metadata.CollectionId,
			expectedResource.Metadata.Id)
		if err != nil {
			return err
		}

		if bytes.Compare(
			resourceWithMetadata.Metadata.Checksum,
			expectedResource.Metadata.Checksum) != 0 {
			return fmt.Errorf("Checksum is not migrated correctly")
		}
	}
	return nil
}
