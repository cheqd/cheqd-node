package utils

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

func Subtract(from []string, value []string) []string {
	return Filter(from, func(s string) bool {
		return !Contains(value, s)
	})
}
