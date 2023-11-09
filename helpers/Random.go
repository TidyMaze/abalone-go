package helpers

import "math/rand"

func RandIn[T any](slice []T) T {
	return slice[rand.Intn(len(slice))]
}

func RandWeight() float64 {
	sign := rand.Intn(2) - 1
	value := rand.Float64() * 100
	return float64(sign) * value
}
