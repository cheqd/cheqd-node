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

func Subtract(minuend []string, subtrahend []string) []string {
	m := map[string]bool{}

	for _, v := range minuend {
		m[v] = true
	}

	for _, v := range subtrahend {
		delete(m, v)
	}

	result := make([]string, 0, len(m))

	for k, _ := range m {
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

	for k, _ := range m {
		result = append(result, k)
	}

	return result
}
