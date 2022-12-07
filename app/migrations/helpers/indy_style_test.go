package helpers

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable(
	"Positive/Negative entries for checking indy style identifier compiling",
	func(id string, outputStr string) {
		Expect(MigrateIndyStyleId(id)).To(Equal(outputStr))
	},

	Entry("Valid: Real case: 16-symbol id", "zGqsJraNJCojDzG4", "QQHVWEaGae5Jts1quynR6M"),
	Entry("Valid: Real case: 32-symbol id", "zGqsJraNJCojDzG4NXY2podMeaESVWvi", "AamcX5kPatrjccMNmuJxSo"),
	Entry("Valid: UUID should not be changed", "F62542C3-4F71-4C21-8A2B-AD8DA460A976", "F62542C3-4F71-4C21-8A2B-AD8DA460A976"),
)
