//go:build upgrade

package upgrade

import (
	"path/filepath"

	integrationtestdata "github.com/cheqd/cheqd-node/tests/integration/testdata"
	cli "github.com/cheqd/cheqd-node/tests/upgrade/cli"
	didtypesv1 "github.com/cheqd/cheqd-node/x/did/types/v1"
	resourcetypesv1 "github.com/cheqd/cheqd-node/x/resource/types/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Upgrade - Pre", func() {
	Context("Before a softare upgrade execution is initiated", func() {
		var DidDocCreatePayload didtypesv1.MsgCreateDidPayload
		var DidDocCreateSignInput []cli.SignInput
		var DidDocUpdatePayload didtypesv1.MsgUpdateDidPayload
		var DidDocUpdateSignInput []cli.SignInput
		var DidDocDeactivatePayload didtypesv1.MsgDeactivateDidPayload
		var DidDocDeactivateSignInput []cli.SignInput
		var ResourceCreatePayload resourcetypesv1.MsgCreateResourcePayload
		var ResourceCreateSignInput []cli.SignInput
		var err error
		It("should load and run existing diddoc payloads - case: create", func() {
			By("matching the glob pattern for existing diddoc payloads")
			ExistingDidDocCreatePayloads, err = Glob(filepath.Join(GENERATED_JSON_DIR, "existing", "diddoc", "create", "*.json"))
			Expect(err).To(BeNil())

			By("matching the glob pattern for existing diddoc sign input")
			ExistingSignInputCreatePayloads, err = Glob(filepath.Join(GENERATED_JSON_DIR, "existing", "diddoc", "create", "signinput", "*.json"))
			Expect(err).To(BeNil())

			for i, payload := range ExistingDidDocCreatePayloads {
				testCase, _ := GetCase(payload)
				By("Running: " + testCase)
				err = Loader(payload, &DidDocCreatePayload)
				Expect(err).To(BeNil())

				err = Loader(ExistingSignInputCreatePayloads[i], &DidDocCreateSignInput)
				Expect(err).To(BeNil())

				res, err := cli.CreateDidLegacy(DidDocCreatePayload, DidDocCreateSignInput, cli.VALIDATOR0)
				Expect(err).To(BeNil())
				Expect(res.Code).To(BeEquivalentTo(0))
			}
		})

		It("should load and run existing resource payloads - case: create", func() {
			By("matching the glob pattern for existing resource payloads")
			ExistingResourceCreatePayloads, err = Glob(filepath.Join(GENERATED_JSON_DIR, "existing", "resource", "create", "*.json"))
			Expect(err).To(BeNil())

			By("matching the glob pattern for existing resource sign input")
			ExistingSignInputCreatePayloads, err = Glob(filepath.Join(GENERATED_JSON_DIR, "existing", "resource", "create", "signinput", "*.json"))
			Expect(err).To(BeNil())

			By("copying the existing resource file to the container")
			ResourceFile, err := integrationtestdata.CreateTestJson(GinkgoT().TempDir())
			Expect(err).To(BeNil())
			_, err = cli.LocalnetExecCopyAbsoluteWithPermissions(ResourceFile, cli.DOCKER_HOME, cli.VALIDATOR0)
			Expect(err).To(BeNil())

			for i, payload := range ExistingResourceCreatePayloads {
				testCase, _ := GetCase(payload)
				By("Running: " + testCase)
				err = Loader(payload, &ResourceCreatePayload)
				Expect(err).To(BeNil())

				err = Loader(ExistingSignInputCreatePayloads[i], &ResourceCreateSignInput)
				Expect(err).To(BeNil())

				// TODO: Add resource file. Right now, it is not possible to create a resource without a file. So we need to copy a file to the container home directory.
				res, err := cli.CreateResource(ResourceCreatePayload.CollectionId, ResourceCreatePayload.Id, ResourceCreatePayload.Name, ResourcePayload.ResourceType, ResourceFile, ResourceCreateSignInput, cli.VALIDATOR0)
				Expect(err).To(BeNil())
				Expect(res.Code).To(BeEquivalentTo(0))
			}
		})

		It("should load and run existing diddoc payloads - case: update", func() {
			By("matching the glob pattern for existing diddoc payloads")
			ExistingDidDocUpdatePayloads, err = Glob(filepath.Join(GENERATED_JSON_DIR, "existing", "diddoc", "update", "*.json"))
			Expect(err).To(BeNil())

			By("matching the glob pattern for existing diddoc sign input")
			ExistingSignInputUpdatePayloads, err = Glob(filepath.Join(GENERATED_JSON_DIR, "existing", "diddoc", "update", "signinput", "*.json"))
			Expect(err).To(BeNil())

			for i, payload := range ExistingDidDocUpdatePayloads {
				testCase, _ := GetCase(payload)
				By("Running: " + testCase)
				err = Loader(payload, &DidDocUpdatePayload)
				Expect(err).To(BeNil())

				err = Loader(ExistingSignInputUpdatePayloads[i], &DidDocUpdateSignInput)
				Expect(err).To(BeNil())

				q, err := cli.QueryDidLegacy(DidDocUpdatePayload.Id, cli.VALIDATOR0)
				Expect(err).To(BeNil())
				Expect(q.Did.Id).To(BeEquivalentTo(DidDocUpdatePayload.Id))

				DidDocUpdatePayload.VersionId = q.Metadata.VersionId

				res, err := cli.UpdateDidLegacy(DidDocUpdatePayload, DidDocUpdateSignInput, cli.VALIDATOR0)
				Expect(err).To(BeNil())
				Expect(res.Code).To(BeEquivalentTo(0))
			}
		})

		It("should load and run existing diddoc payloads - case: deactivate", func() {
			By("matching the glob pattern for existing diddoc payloads")
			ExistingDidDocDeactivatePayloads, err = Glob(filepath.Join(GENERATED_JSON_DIR, "existing", "diddoc", "deactivate", "*.json"))
			Expect(err).To(BeNil())

			By("matching the glob pattern for existing diddoc sign input")
			ExistingSignInputDeactivatePayloads, err = Glob(filepath.Join(GENERATED_JSON_DIR, "existing", "diddoc", "deactivate", "signinput", "*.json"))
			Expect(err).To(BeNil())

			for i, payload := range ExistingDidDocDeactivatePayloads {
				testCase, _ := GetCase(payload)
				By("Running: " + testCase)
				err = Loader(payload, &DidDocDeactivatePayload)
				Expect(err).To(BeNil())

				err = Loader(ExistingSignInputDeactivatePayloads[i], &DidDocDeactivateSignInput)
				Expect(err).To(BeNil())

				res, err := cli.DeactivateDidLegacy(DidDocDeactivatePayload, DidDocDeactivateSignInput, cli.VALIDATOR0)
				Expect(err).To(BeNil())
				Expect(res.Code).To(BeEquivalentTo(0))
			}
		})
	})
})
