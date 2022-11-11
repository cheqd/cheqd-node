//go:build upgrade

package upgrade

import (
	"path/filepath"

	cli "github.com/cheqd/cheqd-node/tests/upgrade/cli"
	didtypesv1 "github.com/cheqd/cheqd-node/x/did/types/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Upgrade - Pre", func() {
	Context("Before a softare upgrade execution is initiated", func() {
		var DidDocCreatePayload didtypesv1.MsgCreateDidPayload
		var DidDocCreateSignInput []cli.SignInput
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

		// TODO: Add more cases here, namely update diddoc, deactivate diddoc, create resource
	})
})
