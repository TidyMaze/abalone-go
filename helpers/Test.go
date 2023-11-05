package helpers

import (
	"fmt"
	"reflect"
)

func AssertEqual(expected interface{}, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		panic(fmt.Sprintf("Expected %v, got %v", expected, actual))
	}
}
