package utils

func SliceContains[T comparable](slice []T, item T) bool {
	for _, x := range slice {
		if x == item {
			return true
		}
	}
	return false
}
