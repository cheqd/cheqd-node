package migration

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Migration - Unit", func() {
	It("should unit test migration scenario handlers", func() {
		err := AssertHandlers()
		Expect(err).To(BeNil())
	})
})
