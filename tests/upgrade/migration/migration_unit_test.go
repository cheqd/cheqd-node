package migration

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/cheqd/cheqd-node/tests/upgrade/migration/scenarios"
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
})
