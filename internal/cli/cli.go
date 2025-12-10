// Package cli provides helper methods to format output for the console
package cli

import (
	"fmt"

	"github.com/fatih/color"
)

// Silent controls if logging operates in silent mode
var Silent bool = false

// SprintValue will format the given value for the console so that it can be identified as a value.
func SprintValue(value any) string {
	return color.CyanString("%v", value)
}

// SprintAlert will format the given value for the console  so that it can be identified an alert information.
func SprintAlert(value any) string {
	return color.HiRedString("%v", value)
}

// Println will print out the given values to the console only when `Silent` is false, or when the given values contain an error.
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
