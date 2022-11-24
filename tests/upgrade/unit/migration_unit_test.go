package migration

import (
	. "github.com/cheqd/cheqd-node/tests/upgrade/unit/scenarios"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Migration - Unit", func() {
	It("checks that Checksum migration handler works", func() {
		By("Ensuring the Checksum migration scenario is successful")

		// Run checksum migration scenario
		err := RunChecksumScenario()
		Expect(err).To(BeNil())
	})

	It("checks that Protobuf migration handler works", func() {
		By("Ensuring the Protobuf migration handler is working as expected")

		// Run Protobuf migration
		err := RunProtobufScenario()
		Expect(err).To(BeNil())
	})

	It("checks IndyStyle Migration", func() {
		By("Ensuring the IndyStyle migration handler is working as expected")

		// Run IndyStyle migration
		err := RunIndyStyleScenario()
		Expect(err).To(BeNil())
	})
})
