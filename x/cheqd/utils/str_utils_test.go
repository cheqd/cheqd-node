package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIndexOf(t *testing.T) {
	cases := []struct {
		array          []string
		searchElement  string
		fromIndex	   int
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


func TestComplement(t *testing.T) {
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
		{[]string{"4", "1", "6", "2", "3"}, []string{"1", "5", "2"}, []string{"4", "6", "3"}},
	}

	for _, tc := range cases {
		actual := Subtract(tc.first, tc.second)
		require.Equal(t, tc.expected, actual)
	}
}
