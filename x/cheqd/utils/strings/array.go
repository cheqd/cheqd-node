package strings

func Index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}

	return -1
}

func Include(vs []string, t string) bool {
	return Index(vs, t) >= 0
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

func Complement(vs []string, ts []string) []string {
	return Filter(vs, func(s string) bool {
		return !Include(ts, s)
	})
}
