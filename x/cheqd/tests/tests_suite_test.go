package tests_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTestsGeneral(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "General Tests Suite")
}

