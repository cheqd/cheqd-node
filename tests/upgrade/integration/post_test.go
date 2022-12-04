package integration

import (
	"path/filepath"

	cli "github.com/cheqd/cheqd-node/tests/upgrade/integration/cli"
	didtypesv2 "github.com/cheqd/cheqd-node/x/did/types"
	resourcetypesv2 "github.com/cheqd/cheqd-node/x/resource/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Upgrade - Post", func() {
	Context("After a software upgrade execution has concluded", func() {
		It("should wait for node catching up", func() {
			By("pinging the node status until catching up is flagged as false")
			err := cli.WaitForCaughtUp(cli.VALIDATOR0, cli.CLI_BINARY_NAME, cli.VOTING_PERIOD*6)
			Expect(err).To(BeNil())
		})

		It("should match the expected module version map", func() {
			By("loading the expected module version map")
			var expected upgradetypes.QueryModuleVersionsResponse
			_, err := Loader(filepath.Join(GENERATED_JSON_DIR, "post", "responses", "module_version_map", "v1.json"), &expected)
			Expect(err).To(BeNil())

			By("matching the expected module version map")
			actual, err := cli.QueryModuleVersionMap(cli.VALIDATOR0)
			Expect(err).To(BeNil())

			Expect(actual.ModuleVersions).To(Equal(expected.ModuleVersions), "module version map mismatch")
		})

		It("should load and run expected diddoc payloads", func() {
			By("matching the glob pattern for existing diddoc payloads")
			ExpectedDidDocUpdateRecords, err := RelGlob(GENERATED_JSON_DIR, "post", "responses", "payloads", "diddoc", "*.json")
			Expect(err).To(BeNil())

			for _, payload := range ExpectedDidDocUpdateRecords {
				var DidDocUpdateRecord didtypesv2.DidDoc

				testCase := GetCaseName(payload)
				By("Running: query " + testCase)
				_, err = Loader(payload, &DidDocUpdateRecord)
				Expect(err).To(BeNil())

				res, err := cli.QueryDid(DidDocUpdateRecord.Id, cli.VALIDATOR0)
				Expect(err).To(BeNil())

				if DidDocUpdateRecord.Context == nil {
					DidDocUpdateRecord.Context = []string{}
				}
				if DidDocUpdateRecord.Authentication == nil {
					DidDocUpdateRecord.Authentication = []string{}
				}
				if DidDocUpdateRecord.AssertionMethod == nil {
					DidDocUpdateRecord.AssertionMethod = []string{}
				}
				if DidDocUpdateRecord.CapabilityInvocation == nil {
					DidDocUpdateRecord.CapabilityInvocation = []string{}
				}
				if DidDocUpdateRecord.CapabilityDelegation == nil {
					DidDocUpdateRecord.CapabilityDelegation = []string{}
				}
				if DidDocUpdateRecord.KeyAgreement == nil {
					DidDocUpdateRecord.KeyAgreement = []string{}
				}
				if DidDocUpdateRecord.Service == nil {
					DidDocUpdateRecord.Service = []*didtypesv2.Service{}
				}
				if DidDocUpdateRecord.AlsoKnownAs == nil {
					DidDocUpdateRecord.AlsoKnownAs = []string{}
				}

				Expect(*res.Value.DidDoc).To(Equal(DidDocUpdateRecord))

			}
		})

		It("should load and run expected resource payloads", func() {
			By("matching the glob pattern for existing resource payloads")
			ExpectedResourceCreateRecords, err := RelGlob(GENERATED_JSON_DIR, "post", "responses", "payloads", "resource", "*.json")
			Expect(err).To(BeNil())

			for _, payload := range ExpectedResourceCreateRecords {
				var ResourceCreateRecord resourcetypesv2.ResourceWithMetadata

				testCase := GetCaseName(payload)
				By("Running: query " + testCase)
				_, err = Loader(payload, &ResourceCreateRecord)
				Expect(err).To(BeNil())

				// TODO: Implement v1 -> v2 protobuf migration handlers.
				// Right now, this will fail.
				// Specifically, the resource is written successfully, but the collectionId will report the resource as not found.
				res, err := cli.QueryResource(ResourceCreateRecord.Metadata.CollectionId, ResourceCreateRecord.Metadata.Id, cli.VALIDATOR0)

				Expect(err).To(BeNil())
				Expect(res.Resource.Metadata.Id).To(Equal(ResourceCreateRecord.Metadata.Id))
				Expect(res.Resource.Metadata.CollectionId).To(Equal(ResourceCreateRecord.Metadata.CollectionId))
				Expect(res.Resource.Metadata.Name).To(Equal(ResourceCreateRecord.Metadata.Name))
				Expect(res.Resource.Metadata.Version).To(Equal(ResourceCreateRecord.Metadata.Version))
				Expect(res.Resource.Metadata.ResourceType).To(Equal(ResourceCreateRecord.Metadata.ResourceType))
				Expect(res.Resource.Metadata.AlsoKnownAs).To(Equal(ResourceCreateRecord.Metadata.AlsoKnownAs))
				Expect(res.Resource.Metadata.MediaType).To(Equal(ResourceCreateRecord.Metadata.MediaType))
				// Created fills while creating. We just ignoring it while checking.
				// Expect(res.Resource.Metadata.Created).To(Equal(ResourceCreateRecord.Metadata.Created))
				Expect(res.Resource.Metadata.Checksum).To(Equal(ResourceCreateRecord.Metadata.Checksum))
				Expect(res.Resource.Metadata.PreviousVersionId).To(Equal(ResourceCreateRecord.Metadata.PreviousVersionId))
				Expect(res.Resource.Metadata.NextVersionId).To(Equal(ResourceCreateRecord.Metadata.NextVersionId))

				// TODO: Add v1 -> v2 deep comparison cases, after defining the migration handlers.
				// e.g.: Migration to Indy format, uuid lowercasing, etc.
				// Checksum migration is already defined as an e2e example.
			}
		})
	})
})
