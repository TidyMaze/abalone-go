package helpers

func Between(v int, min int, max int) bool {
	return v >= min && v <= max
}
