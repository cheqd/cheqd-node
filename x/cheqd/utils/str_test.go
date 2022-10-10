package utils_test

import (
	"sort"

	. "github.com/cheqd/cheqd-node/x/cheqd/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Str utils functionality", func() {
	DescribeTable("IndexOf function",

		func(array []string, searchElement string, fromIndex int, expectedResult int) {
			actual := IndexOf(array, searchElement, fromIndex)
			Expect(actual).To(Equal(expectedResult))
		},

		Entry("Emty array, Empty element, Expected: -1", []string{}, "", 0, -1),
		Entry("Nil as array, Empty element, Expected: -1", nil, "", 0, -1),
		Entry("Desired element is the first. Expected: 0", []string{"1", "2"}, "1", 0, 0),
		Entry("Desired element is the latest one. Expected: 2", []string{"1", "2", "3"}, "3", 0, 2),
		Entry("Desired element is absent. Expected: -1", []string{"1", "2", "3"}, "4", 0, -1),
		Entry("There are more then 1 such elements but search should be started from the beginning. Expected: 0", []string{"4", "1", "6", "2", "3", "4"}, "4", 0, 0),
		Entry("There are more then 1 such elements but search should be started from the index 1. Expected: 5", []string{"4", "1", "6", "2", "3", "4"}, "4", 1, 5),
		Entry("There are more then 1 such elements but search should be started from the index 3. Expected: 5", []string{"4", "1", "6", "2", "3", "4"}, "4", 3, 5),
	)

	DescribeTable("Contains function",

		func(array []string, searchElement string, doesContain bool) {
			Expect(Contains(array, searchElement)).To(Equal(doesContain))
		},
		Entry("Emty array, Empty element, Expected: false", []string{}, "", false),
		Entry("Nil as array, Empty element, Expected: false", nil, "", false),
		Entry("Desired element exists at the position 1. Expected: true", []string{"1", "2"}, "1", true),
		Entry("Desired element exists at the position 2. Expected: true", []string{"1", "2", "3"}, "2", true),
		Entry("Desired element exists at the position 3. Expected: true", []string{"1", "2", "3"}, "3", true),
		Entry("Desired element is absent. Expected: false", []string{"1", "2", "3"}, "123", false),
	)

	DescribeTable("Substract function",

		func(first []string, second []string, expected []string) {
			actual := Subtract(first, second)
			// We can't compare arrays directly cause result of `subtract` is not deterministic
			sort.Strings(actual)
			Expect(expected).To(Equal(actual))
		},
		Entry("Empty first and empty second. Expected empty list", []string{}, []string{}, []string{}),
		Entry("nil as the first and the empty second. Expected empty list", nil, []string{}, []string{}),
		Entry("nil as the first and nil as the second. Expected empty list", nil, nil, []string{}),
		Entry("The same lists. Expected empty list", []string{"1", "2"}, []string{"1", "2"}, []string{}),
		Entry("Substract with empty list. Exected: first array", []string{"1", "2", "3"}, []string{}, []string{"1", "2", "3"}),
		Entry("Substract with nil as the second. Exected: first array", []string{"1", "2", "3"}, nil, []string{"1", "2", "3"}),
		Entry("Substract with totally different list. Expected: first array", []string{"1", "2", "3"}, []string{"4", "5", "6"}, []string{"1", "2", "3"}),
		Entry("Substract. General case. Expected: [1, 2]", []string{"1", "2", "3"}, []string{"1", "5", "2"}, []string{"3"}),
		Entry("Substract. General case. Expected: [3, 4, 6]", []string{"4", "1", "6", "2", "3"}, []string{"1", "5", "2"}, []string{"3", "4", "6"}),
	)

	DescribeTable("Unique function",

		func(array []string, expected []string) {
			actual := Unique(array)
			// We can't compare arrays directly cause result of `unique` is not deterministic
			sort.Strings(actual)
			Expect(expected).To(Equal(actual))
		},
		Entry("Empty array. Expected empty list", []string{}, []string{}),
		Entry("nil as array. Expected empty list", nil, []string{}),
		Entry("Unique array. Expected the same array", []string{"1", "2"}, []string{"1", "2"}),
		Entry("Unique array with length 3. Expected the same array", []string{"1", "3", "2"}, []string{"1", "2", "3"}),
		Entry("General case. Expeceted array with unique elements", []string{"4", "1", "6", "2", "3", "1", "3", "1"}, []string{"1", "2", "3", "4", "6"}),
	)

	DescribeTable("ReplaceInList function",

		func(list []string, oldVal string, newVal string, expected []string) {
			ReplaceInSlice(list, "2", "3")
			Expect(expected).To(Equal(list))
		},
		Entry("Replace 2 with 3", []string{"1", "2", "3", "2"}, "2", "3", []string{"1", "3", "3", "3"}),
	)

	DescribeTable("UniqueSorted function",

		func(input []string, expected []string) {
			Expect(UniqueSorted(input)).To(Equal(expected))
		},
		Entry("General alphabet list", []string{"aa", "bb"}, []string{"aa", "bb"}),
		Entry("General alphabet reverse list", []string{"bb", "aa"}, []string{"aa", "bb"}),
		Entry("General number list", []string{"22", "11"}, []string{"11", "22"}),
	)
})
