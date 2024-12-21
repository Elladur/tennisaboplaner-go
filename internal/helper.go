package internal


func isInSlice[T comparable](element T, slice []T) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}