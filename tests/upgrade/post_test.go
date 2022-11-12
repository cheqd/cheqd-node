//go:build upgrade

package upgrade

import (
	"path/filepath"

	cli "github.com/cheqd/cheqd-node/tests/upgrade/cli"
	didtypesv1 "github.com/cheqd/cheqd-node/x/did/types/v1"
	resourcetypesv1 "github.com/cheqd/cheqd-node/x/resource/types/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Upgrade - Post", func() {
	Context("After a software upgrade execution has concluded", func() {
		var DidDocCreateRecord didtypesv1.Did
		var DidDocUpdateRecord didtypesv1.Did
		var ResourceCreateRecord resourcetypesv1.ResourceHeader
		var err error

		It("should wait for node catching up", func() {
			By("pinging the node status until catching up is flagged as false")
			err := cli.WaitForCaughtUp(cli.VALIDATOR0, cli.CLI_BINARY_NAME, cli.VOTING_PERIOD*6)
			Expect(err).To(BeNil())
		})

		It("should load and run expected diddoc payloads - case: create", func() {
			By("matching the glob pattern for existing diddoc payloads")
			ExpectedDidDocCreateRecords, err = Glob(filepath.Join(GENERATED_JSON_DIR, "existing", "diddoc", "create", "*.json"))
			Expect(err).To(BeNil())

			for _, payload := range ExpectedDidDocCreateRecords {
				testCase, _ := GetCase(payload)
				By("Running: " + testCase)
				err = Loader(payload, &DidDocCreateRecord)
				Expect(err).To(BeNil())

				// TODO: Switch to QueryDid, after migration handlers have been implemented
				res, err := cli.QueryDidLegacy(DidDocCreateRecord.Id, cli.VALIDATOR0)
				Expect(err).To(BeNil())
				Expect(res.Did.Id).To(Equal(DidDocCreateRecord.Id))

				// TODO: Add v1 -> v2 deep comparison cases, after defining the migration handlers.
				// e.g.: Migration to Indy format, uuid lowercasing, etc.
			}
		})

		It("should load and run expected diddoc payloads - case: update", func() {
			By("matching the glob pattern for existing diddoc payloads")
			ExpectedDidDocUpdateRecords, err = Glob(filepath.Join(GENERATED_JSON_DIR, "existing", "diddoc", "update", "*.json"))
			Expect(err).To(BeNil())

			for _, payload := range ExpectedDidDocUpdateRecords {
				testCase, _ := GetCase(payload)
				By("Running: " + testCase)
				err = Loader(payload, &DidDocUpdateRecord)
				Expect(err).To(BeNil())

				// TODO: Switch to QueryDid, after migration handlers have been implemented
				res, err := cli.QueryDidLegacy(DidDocUpdateRecord.Id, cli.VALIDATOR0)
				Expect(err).To(BeNil())
				Expect(res.Did.Id).To(Equal(DidDocUpdateRecord.Id))

				// TODO: Add v1 -> v2 deep comparison cases, after defining the migration handlers.
				// e.g.: Migration to Indy format, uuid lowercasing, etc.
			}
		})

		It("should load and run expected resource payloads - case: create", func() {
			By("matching the glob pattern for existing resource payloads")
			ExpectedResourceCreateRecords, err = Glob(filepath.Join(GENERATED_JSON_DIR, "existing", "resource", "create", "*.json"))
			Expect(err).To(BeNil())

			for _, payload := range ExpectedResourceCreateRecords {
				testCase, _ := GetCase(payload)
				By("Running: " + testCase)
				err = Loader(payload, &ResourceCreateRecord)
				Expect(err).To(BeNil())

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
