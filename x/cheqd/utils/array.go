package utils

func StringArrayContains(array []string, item string) bool {
	for _, i := range array {
		if item == i {
			return true
		}
	}

	return false
}
