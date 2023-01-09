package utils

import "sort"

func IndexOf(array []string, searchElement string, fromIndex int) int {
	for i, v := range array[fromIndex:] {
		if v == searchElement {
			return fromIndex + i
		}
	}

	return -1
}

func Contains(vs []string, t string) bool {
	return IndexOf(vs, t, 0) >= 0
}

func Filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func Subtract(minuend []string, subtrahend []string) []string {
	m := map[string]bool{}

	for _, v := range minuend {
		m[v] = true
	}

	for _, v := range subtrahend {
		delete(m, v)
	}

	result := make([]string, 0, len(m))

	for k := range m {
		result = append(result, k)
	}

	return result
}

// Unique returns a copy of the passed array with duplicates removed
func Unique(array []string) []string {
	m := map[string]bool{}

	for _, v := range array {
		m[v] = true
	}

	result := make([]string, 0, len(m))

	for k := range m {
		result = append(result, k)
	}

	return result
}

func IsUnique(list []string) bool {
	set := map[string]bool{}

	for _, did := range list {
		set[did] = true
	}

	return len(list) == len(set)
}

func ToInterfaces(list []string) []interface{} {
	res := make([]interface{}, len(list))

	for i := range list {
		res[i] = list[i]
	}

	return res
}

func ReplaceInSlice(list []string, old, new string) {
	for i := range list {
		if list[i] == old {
			list[i] = new
		}
	}
}

func UniqueSorted(ls []string) []string {
	tmp := Unique(ls)
	sort.Strings(tmp)
	return tmp
}

func StrBytes(p string) []byte {
	return []byte(p)
}
