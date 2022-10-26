//go:build upgrade

package upgrade

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Pre() is a function that runs before the upgrade test suite.
// Idiomatically, it is called from the upgrade_suite_test.go file, in the BeforeSuite() function.
// We will keep both BeforeSuite() and Pre() callback here for easiness of conceptual understanding.
var _ = BeforeSuite(func() {
	Pre()
})

func Pre() {
	// TODO: Implement this function
	Expect(true).To(BeTrue())
}
