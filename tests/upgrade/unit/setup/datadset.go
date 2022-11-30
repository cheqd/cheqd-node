package setup

import (
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	didtypesv1 "github.com/cheqd/cheqd-node/x/did/types/v1"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	resourcetypesv1 "github.com/cheqd/cheqd-node/x/resource/types/v1"

	. "github.com/onsi/gomega"
)

// Existing

type ExistingDataset struct {
	setup  TestSetup
	loader Loader

	DidDocsV1   []didtypesv1.StateValue
	ResourcesV1 []resourcetypesv1.Resource

	DidDocsV2   []didtypes.DidDocWithMetadata
	ResourcesV2 []resourcetypes.ResourceWithMetadata
}

func NewExistingDataset(setup TestSetup) *ExistingDataset {
	return &ExistingDataset{
		setup:  setup,
		loader: Loader{},
	}
}

func (d *ExistingDataset) AddDidDocV1(pathToDir, prefix string) error {
	files, err := d.loader.GetLsitOfFiles(pathToDir, prefix)
	if err != nil {
		return err
	}

	for _, path_to_file := range files {
		var existingDidDoc didtypesv1.StateValue
		err = d.loader.LoadFile(path_to_file, &existingDidDoc, d.setup)
		if err != nil {
			return err
		}

		d.DidDocsV1 = append(d.DidDocsV1, existingDidDoc)
	}

	return nil
}

func (d *ExistingDataset) MustAddDidDocV1(pathToDir, prefix string) {
	err := d.AddDidDocV1(pathToDir, prefix)
	Expect(err).To(BeNil())
}

func (d *ExistingDataset) AddDidDocV2(pathToDir, prefix string) error {
	files, err := d.loader.GetLsitOfFiles(pathToDir, prefix)
	if err != nil {
		return err
	}

	for _, path_to_file := range files {
		var existingDidDoc didtypes.DidDocWithMetadata
		err = d.loader.LoadFile(path_to_file, &existingDidDoc, d.setup)
		if err != nil {
			return err
		}

		d.DidDocsV2 = append(d.DidDocsV2, existingDidDoc)
	}

	return nil
}

func (d *ExistingDataset) MustAddDidDocV2(pathToDir, prefix string) {
	err := d.AddDidDocV2(pathToDir, prefix)
	Expect(err).To(BeNil())
}

func (d *ExistingDataset) AddResourceV1(pathToDir, prefix string) error {
	files, err := d.loader.GetLsitOfFiles(pathToDir, prefix)
	if err != nil {
		return err
	}

	for _, path_to_file := range files {
		var existingResource resourcetypesv1.Resource
		err = d.loader.LoadFile(path_to_file, &existingResource, d.setup)
		if err != nil {
			return err
		}

		d.ResourcesV1 = append(d.ResourcesV1, existingResource)
	}
	return nil
}

func (d *ExistingDataset) MustAddResourceV1(pathToDir, prefix string) {
	err := d.AddResourceV1(pathToDir, prefix)
	Expect(err).To(BeNil())
}

func (d *ExistingDataset) AddResourceV2(pathToDir, prefix string) error {
	files, err := d.loader.GetLsitOfFiles(pathToDir, prefix)
	if err != nil {
		return err
	}

	for _, path_to_file := range files {
		var existingResource resourcetypes.ResourceWithMetadata
		err = d.loader.LoadFile(path_to_file, &existingResource, d.setup)
		if err != nil {
			return err
		}

		d.ResourcesV2 = append(d.ResourcesV2, existingResource)
	}

	return nil
}

func (d *ExistingDataset) MustAddResourceV2(pathToDir, prefix string) {
	err := d.AddResourceV2(pathToDir, prefix)
	Expect(err).To(BeNil())
}

func (d *ExistingDataset) FillStore() error {
	for _, didDoc := range d.DidDocsV1 {
		err := d.setup.DidKeeperV1.SetDid(&d.setup.SdkCtx, &didDoc)
		if err != nil {
			return err
		}
	}
	for _, resource := range d.ResourcesV1 {
		err := d.setup.ResourceKeeperV1.SetResource(&d.setup.SdkCtx, &resource)
		if err != nil {
			return err
		}
	}

	for _, did_doc := range d.DidDocsV2 {
		err := d.setup.DidKeeper.AddNewDidDocVersion(&d.setup.SdkCtx, &did_doc)
		if err != nil {
			return err
		}
	}
	for _, resource := range d.ResourcesV2 {
		err := d.setup.ResourceKeeper.SetResource(&d.setup.SdkCtx, &resource)
		if err != nil {
			return err
		}
	}
	return nil
}

// Expected

type ExpectedDataset struct {
	setup  TestSetup
	loader Loader

	DidDocs   []didtypes.DidDocWithMetadata
	Resources []resourcetypes.ResourceWithMetadata
}

func NewExpectedDataset(setup TestSetup) *ExpectedDataset {
	return &ExpectedDataset{
		setup:  setup,
		loader: Loader{},
	}
}

func (d *ExpectedDataset) AddDidDocV2(pathToDir, prefix string) error {
	files, err := d.loader.GetLsitOfFiles(pathToDir, prefix)
	if err != nil {
		return err
	}

	for _, path_to_file := range files {
		var expectedDidDoc didtypes.DidDocWithMetadata
		err := d.loader.LoadFile(path_to_file, &expectedDidDoc, d.setup)
		if err != nil {
			return err
		}

		d.DidDocs = append(d.DidDocs, expectedDidDoc)
	}

	return nil
}

func (d *ExpectedDataset) MustAddDidDocV2(pathToDir, prefix string) {
	err := d.AddDidDocV2(pathToDir, prefix)
	Expect(err).To(BeNil())
}

func (d *ExpectedDataset) AddResourceV2(pathToDir, prefix string) error {
	files, err := d.loader.GetLsitOfFiles(pathToDir, prefix)
	if err != nil {
		return err
	}

	for _, path_to_file := range files {
		var expectedResource resourcetypes.ResourceWithMetadata
		err = d.loader.LoadFile(path_to_file, &expectedResource, d.setup)
		if err != nil {
			return err
		}

		d.Resources = append(d.Resources, expectedResource)
	}
	return nil
}

func (d *ExpectedDataset) MustAddResourceV2(pathToDir, prefix string) {
	err := d.AddResourceV2(pathToDir, prefix)
	Expect(err).To(BeNil())
}

func (d *ExpectedDataset) CheckStore() error {
	for _, expectedDidDoc := range d.DidDocs {
		didDoc, err := d.setup.DidKeeper.GetLatestDidDoc(&d.setup.SdkCtx, expectedDidDoc.DidDoc.Id)
		if err != nil {
			return err
		}

		if didDoc.DidDoc.Context == nil {
			didDoc.DidDoc.Context = []string{}
		}
		if didDoc.DidDoc.Authentication == nil {
			didDoc.DidDoc.Authentication = []string{}
		}
		if didDoc.DidDoc.AssertionMethod == nil {
			didDoc.DidDoc.AssertionMethod = []string{}
		}
		if didDoc.DidDoc.CapabilityInvocation == nil {
			didDoc.DidDoc.CapabilityInvocation = []string{}
		}
		if didDoc.DidDoc.CapabilityDelegation == nil {
			didDoc.DidDoc.CapabilityDelegation = []string{}
		}
		if didDoc.DidDoc.KeyAgreement == nil {
			didDoc.DidDoc.KeyAgreement = []string{}
		}
		if didDoc.DidDoc.Service == nil {
			didDoc.DidDoc.Service = []*didtypes.Service{}
		}
		if didDoc.DidDoc.AlsoKnownAs == nil {
			didDoc.DidDoc.AlsoKnownAs = []string{}
		}

		Expect(didDoc.DidDoc).To(Equal(expectedDidDoc.DidDoc))
		Expect(didDoc.Metadata.VersionId).To(Equal(expectedDidDoc.Metadata.VersionId))
		Expect(didDoc.Metadata.Deactivated).To(Equal(expectedDidDoc.Metadata.Deactivated))
	}

	for _, expectedResource := range d.Resources {
		resource, err := d.setup.ResourceKeeper.GetResource(&d.setup.SdkCtx,
			expectedResource.Metadata.CollectionId,
			expectedResource.Metadata.Id)
		if err != nil {
			return err
		}
		Expect(resource).To(Equal(expectedResource))
	}

	return nil
}
