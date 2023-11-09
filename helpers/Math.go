package helpers

func Between(v int8, min int8, max int8) bool {
	return v >= min && v <= max
}
