package helpers

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable(
	"Positive/Negative entries for UUID style identifier compiling",
	func(id string, outputStr string) {
		Expect(MigrateUUIDId(id)).To(Equal(outputStr))
	},

	Entry("Valid: Lower case UUID should be kept as is", "cc490981-66a0-4d87-84b7-5aad0d699fd0", "cc490981-66a0-4d87-84b7-5aad0d699fd0"),
	Entry("Valid: Not uuid should be kept as is", "zGqsJraNJCojDzG4NXY2podMeaESVWvi", "zGqsJraNJCojDzG4NXY2podMeaESVWvi"),
	Entry("Valid: Upper case uuid should be hashed", "F62542C3-4F71-4C21-8A2B-AD8DA460A976", "587bf72c-1963-5ec4-a5ac-0ac5fd8521ce"),
	Entry("Valid: Mixed case uuid should be hashed", "F62542C3-4F71-4C21-8a2b-ad8da460a976", "3e2bc7a8-873c-5dd4-b713-6c9f96bf6a5f"),
)
