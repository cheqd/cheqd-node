package ante_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAnte(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ante Suite")
}
