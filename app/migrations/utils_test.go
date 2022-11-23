package migrations_test

import (
	. "github.com/cheqd/cheqd-node/app/migrations"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("IndyStyleIdentifier", func() {
	DescribeTable(
		"Positive/Negative entries for checking indy style identifier compiling",
		func(id string, outputStr string) {
			Expect(IndyStyleId(id)).To(Equal(outputStr))
		},

		Entry("Valid: Real case: 16-symbol id", "zGqsJraNJCojDzG4", "QQHVWEaGae5Jts1quynR6M"),
		Entry("Valid: Real case: 32-symbol id", "zGqsJraNJCojDzG4NXY2podMeaESVWvi", "AamcX5kPatrjccMNmuJxSo"),
		Entry("Valid: UUID should not be changed", "F62542C3-4F71-4C21-8A2B-AD8DA460A976", "F62542C3-4F71-4C21-8A2B-AD8DA460A976"),
	)

	DescribeTable("Replace HeaderKey to DataKey", func(headerKey, expectedDataKey []byte) {
		Expect(ResourceV1HeaderkeyToDataKey(headerKey)).To(Equal(expectedDataKey))
	},
		Entry(
			"Valid: Expected behaviour", 
			[]byte("resource-header:zGqsJraNJCojDzG4:ba62c728-cb15-498b-8e9e-9259cc242186"), 
			[]byte("resource-data:zGqsJraNJCojDzG4:ba62c728-cb15-498b-8e9e-9259cc242186")),
		)
})
