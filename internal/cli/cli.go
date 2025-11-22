package cli

import (
	"fmt"

	"github.com/fatih/color"
)

var Silent bool = false

func SprintValue(value any) string {
	return color.CyanString("%v", value)
}

func SprintAlert(value any) string {
	return color.HiRedString("%v", value)
}

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
