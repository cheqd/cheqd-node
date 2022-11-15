package migrations_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTestsMigrations(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "App Migrations unit tests")
}