package utils

func CompareOwners(authors []string, controllers []string) bool {
	type void struct{}
	var member void

	controllerSet := make(map[string]void)
	for _, author := range authors {
		controllerSet[author] = member
	}
	result := true
	for _, controller := range controllers {
		_, exists := controllerSet[controller]
		result = result && exists
	}

	return result
}
