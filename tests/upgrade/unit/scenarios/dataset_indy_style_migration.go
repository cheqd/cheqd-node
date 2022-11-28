package scenarios

import (
	"path/filepath"

	. "github.com/onsi/gomega"

	migrationsetup "github.com/cheqd/cheqd-node/tests/upgrade/unit/setup"

	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
)

type IndyStyleBuilder struct {
	setup      migrationsetup.TestSetup
	loader     Loader
	is_dataset IndyStyleDataSet
}

func NewIndyStyleBuilder(setup migrationsetup.TestSetup) IndyStyleBuilder {
	return IndyStyleBuilder{
		setup:      setup,
		loader:     Loader{},
		is_dataset: NewIndyStyleDataSet(setup),
	}
}

type IndyStyleDataSet struct {
	setup             migrationsetup.TestSetup
	existingDIDDocs   []didtypes.DidDocWithMetadata
	existingResources []resourcetypes.ResourceWithMetadata
	expectedDidDocs   []didtypes.DidDocWithMetadata
	expectedResources []resourcetypes.ResourceWithMetadata
}

func NewIndyStyleDataSet(setup migrationsetup.TestSetup) IndyStyleDataSet {
	return IndyStyleDataSet{
		setup:             setup,
		existingDIDDocs:   []didtypes.DidDocWithMetadata{},
		existingResources: []resourcetypes.ResourceWithMetadata{},
		expectedDidDocs:   []didtypes.DidDocWithMetadata{},
		expectedResources: []resourcetypes.ResourceWithMetadata{},
	}
}

func (is_builder *IndyStyleBuilder) BuildDataSet(setup migrationsetup.TestSetup) (IndyStyleDataSet, error) {
	err := is_builder.buildExistingDids()
	if err != nil {
		return IndyStyleDataSet{}, err
	}
	err = is_builder.buildExistingResources()
	if err != nil {
		return IndyStyleDataSet{}, err
	}
	err = is_builder.buildExpectedDids()
	if err != nil {
		return IndyStyleDataSet{}, err
	}
	err = is_builder.buildExpectedResources()
	if err != nil {
		return IndyStyleDataSet{}, err
	}

	return is_builder.is_dataset, err
}

func (is_builder *IndyStyleBuilder) buildExistingDids() error {
	var existingDidDoc didtypes.DidDocWithMetadata
	files, err := is_builder.loader.GetLsitOfFiles(
		filepath.Join(GENERATED_JSON_DIR, "payload", "existing", "v2", "indy_style"),
		"diddoc")
	if err != nil {
		return err
	}
	for _, path_to_file := range files {
		err = is_builder.loader.LoadFile(
			path_to_file,
			&existingDidDoc,
			is_builder.setup,
		)
		if err != nil {
			return err
		}
		is_builder.is_dataset.existingDIDDocs = append(is_builder.is_dataset.existingDIDDocs, existingDidDoc)
	}
	return nil
}

func (is_builder *IndyStyleBuilder) buildExistingResources() error {
	files, err := is_builder.loader.GetLsitOfFiles(
		filepath.Join(GENERATED_JSON_DIR, "payload", "existing", "v2", "indy_style"),
		"resource")
	if err != nil {
		return err
	}
	for _, path_to_file := range files {
		var existingResource resourcetypes.ResourceWithMetadata
		err = is_builder.loader.LoadFile(
			path_to_file,
			&existingResource,
			is_builder.setup,
		)
		if err != nil {
			return err
		}
		is_builder.is_dataset.existingResources = append(is_builder.is_dataset.existingResources, existingResource)
	}
	return nil
}

func (is_builder *IndyStyleBuilder) buildExpectedDids() error {
	var expectedDidDoc didtypes.DidDocWithMetadata
	files, err := is_builder.loader.GetLsitOfFiles(
		filepath.Join(GENERATED_JSON_DIR, "payload", "expected", "v2", "indy_style"),
		"diddoc")
	if err != nil {
		return err
	}
	for _, path_to_file := range files {
		err := is_builder.loader.LoadFile(
			path_to_file,
			&expectedDidDoc,
			is_builder.setup,
		)
		if err != nil {
			return err
		}
		is_builder.is_dataset.expectedDidDocs = append(is_builder.is_dataset.expectedDidDocs, expectedDidDoc)
	}
	return nil
}

func (is_builder *IndyStyleBuilder) buildExpectedResources() error {
	var expectedResource resourcetypes.ResourceWithMetadata
	files, err := is_builder.loader.GetLsitOfFiles(
		filepath.Join(GENERATED_JSON_DIR, "payload", "expected", "v2", "indy_style"),
		"resource")
	if err != nil {
		return err
	}
	for _, path_to_file := range files {
		err = is_builder.loader.LoadFile(
			path_to_file,
			&expectedResource,
			is_builder.setup,
		)
		if err != nil {
			return err
		}
		is_builder.is_dataset.expectedResources = append(is_builder.is_dataset.expectedResources, expectedResource)
	}
	return nil
}

func (is_dataset *IndyStyleDataSet) Prepare() error {
	for _, did_doc := range is_dataset.existingDIDDocs {
		err := is_dataset.setup.DidKeeper.AddNewDidDocVersion(&is_dataset.setup.SdkCtx, &did_doc)
		if err != nil {
			return err
		}
	}
	for _, resource := range is_dataset.existingResources {
		err := is_dataset.setup.ResourceKeeper.SetResource(&is_dataset.setup.SdkCtx, &resource)
		if err != nil {
			return err
		}
	}
	return nil
}

func (is_dataset *IndyStyleDataSet) Validate() error {
	var (
		expectedDidDoc   didtypes.DidDocWithMetadata
		expectedResource resourcetypes.ResourceWithMetadata
	)
	for _, expectedDidDoc = range is_dataset.expectedDidDocs {
		didDoc, err := is_dataset.setup.DidKeeper.GetLatestDidDoc(&is_dataset.setup.SdkCtx, expectedDidDoc.DidDoc.Id)
		if err != nil {
			return err
		}
		Expect(didDoc.DidDoc).To(Equal(expectedDidDoc.DidDoc))
	}

	for _, expectedResource = range is_dataset.expectedResources {
		resource, err := is_dataset.setup.ResourceKeeper.GetResource(&is_dataset.setup.SdkCtx,
			expectedResource.Metadata.CollectionId,
			expectedResource.Metadata.Id)
		if err != nil {
			return err
		}
		Expect(resource).To(Equal(expectedResource))
	}

	return nil
}
