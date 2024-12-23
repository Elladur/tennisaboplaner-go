package internal

import "math/rand"

func isInSlice[T comparable](element T, slice []T) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}

func pop[T any](slice []T) ([]T, T) {
	el := slice[len(slice)-1]
	slice = slice[:len(slice)-1]
	return slice, el
}

func shuffle[T any](slice []T) []T {
	for i := range slice {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

func convertToSlice[K comparable, V any](mapping map[K]V) []V {
	values := make([]V, 0, len(mapping))
	for _, value := range mapping {
		values = append(values, value)
	}
	return values
}
