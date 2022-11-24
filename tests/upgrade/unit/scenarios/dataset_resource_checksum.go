package scenarios

import (
	"bytes"
	"fmt"
	"path/filepath"

	migrationsetup "github.com/cheqd/cheqd-node/tests/upgrade/unit/setup"

	didtypes "github.com/cheqd/cheqd-node/x/did/types"

	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
)

type ChecksumBuilder struct {
	setup  migrationsetup.TestSetup
	loader Loader
	cds    ChecksumDataSet
}

func NewChecksumBuilder(setup migrationsetup.TestSetup) ChecksumBuilder {
	return ChecksumBuilder{
		setup:  setup,
		loader: Loader{},
		cds:    NewChecksumDataSet(setup),
	}
}

type ChecksumDataSet struct {
	setup             migrationsetup.TestSetup
	existingDIDDocs   []didtypes.DidDocWithMetadata
	existingResources []resourcetypes.ResourceWithMetadata
	expectedDidDocs   []didtypes.DidDocWithMetadata
	expectedResources []resourcetypes.ResourceWithMetadata
}

func NewChecksumDataSet(setup migrationsetup.TestSetup) ChecksumDataSet {
	return ChecksumDataSet{
		setup:             setup,
		existingDIDDocs:   []didtypes.DidDocWithMetadata{},
		existingResources: []resourcetypes.ResourceWithMetadata{},
		expectedDidDocs:   []didtypes.DidDocWithMetadata{},
		expectedResources: []resourcetypes.ResourceWithMetadata{},
	}
}

func (cb *ChecksumBuilder) BuildDataSet(setup migrationsetup.TestSetup) (ChecksumDataSet, error) {
	err := cb.buildExistingDids()

	if err != nil {
		return ChecksumDataSet{}, err
	}
	err = cb.buildExistingResources()
	if err != nil {
		return ChecksumDataSet{}, err
	}
	err = cb.buildExpectedDids()
	if err != nil {
		return ChecksumDataSet{}, err
	}
	err = cb.buildExpectedResources()
	if err != nil {
		return ChecksumDataSet{}, err
	}

	return cb.cds, err
}

func (cb *ChecksumBuilder) buildExistingDids() error {
	return nil
}

func (cb *ChecksumBuilder) buildExistingResources() error {
	var existingResource resourcetypes.ResourceWithMetadata
	files, err := cb.loader.GetLsitOfFiles(
		filepath.Join(GENERATED_JSON_DIR, "payload", "existing", "v2", "checksum"),
		"resource")
	if err != nil {
		return err
	}
	for _, path_to_file := range files {
		err := cb.loader.LoadFile(
			path_to_file,
			&existingResource,
			cb.cds.setup,
		)
		if err != nil {
			return err
		}
		cb.cds.existingResources = append(cb.cds.existingResources, existingResource)
	}
	return nil
}

func (cb *ChecksumBuilder) buildExpectedDids() error {
	return nil
}

func (cb *ChecksumBuilder) buildExpectedResources() error {
	var expectedResource resourcetypes.ResourceWithMetadata
	files, err := cb.loader.GetLsitOfFiles(
		filepath.Join(GENERATED_JSON_DIR, "payload", "expected", "v2", "checksum"),
		"resource")
	if err != nil {
		return err
	}
	for _, path_to_file := range files {
		err := cb.loader.LoadFile(
			path_to_file,
			&expectedResource,
			cb.cds.setup,
		)
		if err != nil {
			return err
		}
		cb.cds.expectedResources = append(cb.cds.expectedResources, expectedResource)
	}
	return nil
}

func (cds *ChecksumDataSet) Prepare() error {
	for _, resource := range cds.existingResources {
		err := cds.setup.ResourceKeeper.SetResource(&cds.setup.SdkCtx, &resource)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cds *ChecksumDataSet) Validate() error {
	var (
		resourceWithMetadata resourcetypes.ResourceWithMetadata
		err                  error
	)

	for _, expectedResource := range cds.expectedResources {
		resourceWithMetadata, err = cds.setup.ResourceKeeper.GetResource(
			&cds.setup.SdkCtx,
			expectedResource.Metadata.CollectionId,
			expectedResource.Metadata.Id)
		if err != nil {
			return err
		}

		if !bytes.Equal(resourceWithMetadata.Metadata.Checksum, expectedResource.Metadata.Checksum) {
			return fmt.Errorf("checksum is not migrated correctly")
		}
	}
	return nil
}
