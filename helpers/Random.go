package helpers

import "math/rand"

func RandIn[T any](slice []T) T {
	return slice[rand.Intn(len(slice))]
}
