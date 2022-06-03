package utils

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndexOf(t *testing.T) {
	cases := []struct {
		array          []string
		searchElement  string
		fromIndex      int
		expectedResult int
	}{
		{[]string{}, "", 0, -1},
		{nil, "", 0, -1},
		{[]string{"1", "2"}, "1", 0, 0},
		{[]string{"1", "2", "3"}, "3", 0, 2},
		{[]string{"1", "2", "3"}, "4", 0, -1},
		{[]string{"4", "1", "6", "2", "3", "4"}, "4", 0, 0},
		{[]string{"4", "1", "6", "2", "3", "4"}, "4", 1, 5},
		{[]string{"4", "1", "6", "2", "3", "4"}, "4", 3, 5},
	}

	for _, tc := range cases {
		actual := IndexOf(tc.array, tc.searchElement, tc.fromIndex)
		require.Equal(t, tc.expectedResult, actual)
	}
}

func TestContains(t *testing.T) {
	cases := []struct {
		array          []string
		searchElement  string
		expectedResult bool
	}{
		{[]string{}, "", false},
		{nil, "", false},
		{[]string{"1", "2"}, "1", true},
		{[]string{"1", "2", "3"}, "2", true},
		{[]string{"1", "2", "3"}, "3", true},
		{[]string{"1", "2", "3"}, "123", false},
	}

	for _, tc := range cases {
		actual := Contains(tc.array, tc.searchElement)
		require.Equal(t, tc.expectedResult, actual)
	}
}

func TestSubtract(t *testing.T) {
	cases := []struct {
		first    []string
		second   []string
		expected []string
	}{
		{[]string{}, []string{}, []string{}},
		{nil, []string{}, []string{}},
		{nil, nil, []string{}},
		{[]string{"1", "2"}, []string{"1", "2"}, []string{}},
		{[]string{"1", "2", "3"}, []string{}, []string{"1", "2", "3"}},
		{[]string{"1", "2", "3"}, nil, []string{"1", "2", "3"}},
		{[]string{"1", "2", "3"}, []string{"4", "5", "6"}, []string{"1", "2", "3"}},
		{[]string{"1", "2", "3"}, []string{"1", "5", "2"}, []string{"3"}},
		{[]string{"4", "1", "6", "2", "3"}, []string{"1", "5", "2"}, []string{"3", "4", "6"}},
	}

	for _, tc := range cases {
		actual := Subtract(tc.first, tc.second)
		// We can't compare arrays directly cause result of `subtract` is not deterministic
		sort.Strings(actual)
		require.Equal(t, tc.expected, actual)
	}
}

func TestUnique(t *testing.T) {
	cases := []struct {
		array    []string
		expected []string
	}{
		{[]string{}, []string{}},
		{nil, []string{}},
		{[]string{"1", "2"}, []string{"1", "2"}},
		{[]string{"1", "3", "2"}, []string{"1", "2", "3"}},
		{[]string{"4", "1", "6", "2", "3", "1", "3", "1"}, []string{"1", "2", "3", "4", "6"}},
	}

	for _, tc := range cases {
		actual := Unique(tc.array)
		// We can't compare arrays directly cause result of `unique` is not deterministic
		sort.Strings(actual)
		require.Equal(t, tc.expected, actual)
	}
}

func TestReplaceInList(t *testing.T) {
	list := []string{"1", "2", "3", "2"}
	ReplaceInSlice(list, "2", "3")

	require.Equal(t, []string{"1", "3", "3", "3"}, list)
}

func TestUniqueSorted(t *testing.T) {
	cases := []struct {
		name   string
		input  []string
		output []string
	}{
		{"General alphabet list", []string{"aa", "bb"}, []string{"aa", "bb"}},
		{"General alphabet reverse list", []string{"bb", "aa"}, []string{"aa", "bb"}},
		{"General number list", []string{"22", "11"}, []string{"11", "22"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := UniqueSorted(tc.input)
			require.Equal(t, res, tc.output)
		})
	}
}
