//go:build upgrade_integration

package integration

import (
	"path/filepath"

	cli "github.com/cheqd/cheqd-node/tests/upgrade/integration/cli"
	didtypesv1 "github.com/cheqd/cheqd-node/x/did/types/v1"
	resourcetypesv1 "github.com/cheqd/cheqd-node/x/resource/types/v1"
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
			err := Loader(filepath.Join(GENERATED_JSON_DIR, "expected", "module_version_map", "v1.json"), &expected)

			By("matching the expected module version map")
			actual, err := cli.QueryModuleVersionMap(cli.VALIDATOR0)
			Expect(err).To(BeNil())

			Expect(actual.ModuleVersions).To(Equal(expected.ModuleVersions), "module version map mismatch")
		})

		It("should load and run expected diddoc payloads - case: create", func() {
			By("matching the glob pattern for existing diddoc payloads")
			ExpectedDidDocCreateRecords, err := RelGlob(GENERATED_JSON_DIR, "expected", "diddoc", "create", "*.json")
			Expect(err).To(BeNil())

			for _, payload := range ExpectedDidDocCreateRecords {
				var DidDocCreateRecord didtypesv1.Did

				testCase := GetCaseName(payload)
				By("Running: query " + testCase)
				err = Loader(payload, &DidDocCreateRecord)
				Expect(err).To(BeNil())

				// TODO: Implement v1 -> v2 protobuf migration handlers.
				// Right now, this will fail.
				res, err := cli.QueryDid(DidDocCreateRecord.Id, cli.VALIDATOR0)
				Expect(err).To(BeNil())
				Expect(res.Value.DidDoc.Id).To(Equal(DidDocCreateRecord.Id))

				// TODO: Add v1 -> v2 deep comparison cases, after defining the migration handlers.
				// e.g.: Migration to Indy format, uuid lowercasing, etc.
			}
		})

		It("should load and run expected diddoc payloads - case: update", func() {
			By("matching the glob pattern for existing diddoc payloads")
			ExpectedDidDocUpdateRecords, err := RelGlob(GENERATED_JSON_DIR, "expected", "diddoc", "update", "*.json")
			Expect(err).To(BeNil())

			for _, payload := range ExpectedDidDocUpdateRecords {
				var DidDocUpdateRecord didtypesv1.Did

				testCase := GetCaseName(payload)
				By("Running: query " + testCase)
				err = Loader(payload, &DidDocUpdateRecord)
				Expect(err).To(BeNil())

				// TODO: Implement v1 -> v2 protobuf migration handlers.
				// Right now, this will fail.
				res, err := cli.QueryDid(DidDocUpdateRecord.Id, cli.VALIDATOR0)
				Expect(err).To(BeNil())
				Expect(res.Value.DidDoc.Id).To(Equal(DidDocUpdateRecord.Id))

				// TODO: Add v1 -> v2 deep comparison cases, after defining the migration handlers.
				// e.g.: Migration to Indy format, uuid lowercasing, etc.
			}
		})

		It("should load and run expected resource payloads - case: create", func() {
			By("matching the glob pattern for existing resource payloads")
			ExpectedResourceCreateRecords, err := RelGlob(GENERATED_JSON_DIR, "expected", "resource", "create", "*.json")
			Expect(err).To(BeNil())

			for _, payload := range ExpectedResourceCreateRecords {
				var ResourceCreateRecord resourcetypesv1.ResourceHeader

				testCase := GetCaseName(payload)
				By("Running: query " + testCase)
				err = Loader(payload, &ResourceCreateRecord)
				Expect(err).To(BeNil())

				// TODO: Implement v1 -> v2 protobuf migration handlers.
				// Right now, this will fail.
				// Specifically, the resource is written successfully, but the collectionId will report the resource as not found.
				res, err := cli.QueryResourceLegacy(ResourceCreateRecord.CollectionId, ResourceCreateRecord.Id, cli.VALIDATOR0)
				Expect(err).To(BeNil())
				Expect(res.Resource.Header.Id).To(Equal(ResourceCreateRecord.Id))

				// TODO: Add v1 -> v2 deep comparison cases, after defining the migration handlers.
				// e.g.: Migration to Indy format, uuid lowercasing, etc.
				// Checksum migration is already defined as an e2e example.
			}
		})
	})
})
