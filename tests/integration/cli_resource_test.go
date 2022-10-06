//go:build integration

package integration

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("cheqd cli", func() {
	It("works", func() {
		Expect(true).To(BeTrue())
	})
})
