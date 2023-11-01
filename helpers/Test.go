package helpers

import "fmt"

func AssertEqual(expected interface{}, actual interface{}) {
	if expected != actual {
		panic(fmt.Sprintf("Expected %v, got %v", expected, actual))
	}
}
