package cli

import (
	"fmt"

	"github.com/fatih/color"
)

var Silent bool = false

// Format the given value the that it will be printed with value coloring.
func SprintValue(value any) string {
	return color.CyanString("%v", value)
}

// Format the given value so that it will be printed with alert coloring.
func SprintAlert(value any) string {
	return color.HiRedString("%v", value)
}

// Will print out the given values. Will omit printing when the Silent setting is true.
// Even if Silent is true, when one of the given values is an error, all of the values will be printed!
func Println(value ...any) {

	var isError bool

	for _, v := range value {
		if _, isError = v.(error); isError == true {
			break
		}
	}

	if Silent == false || isError {
		fmt.Println(value...)
	}
}
