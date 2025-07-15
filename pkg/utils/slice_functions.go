package utils

func RemoveElement[T comparable](slice []T, element T) []T {
	for i, v := range slice {
		if v == element {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func ReplaceElement[T comparable](slice []T, element T, newElement T) []T {
	for i, v := range slice {
		if v == element {
			slice[i] = newElement
			return slice
		}
	}
	return slice
}
