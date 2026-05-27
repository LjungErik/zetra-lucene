package utils

func Grow[T any](s []T, minLength int) []T {
	if minLength <= len(s) {
		return s
	}

	newLength := minLength + (minLength >> 3)
	if newLength < minLength+3 {
		newLength = minLength + 3
	}

	grown := make([]T, newLength)
	copy(grown, s)
	return grown
}
