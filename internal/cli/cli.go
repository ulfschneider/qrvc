package cli

import "github.com/fatih/color"

func SprintValue(value any) string {
	return color.CyanString("%v", value)
}

func SprintAlert(value any) string {
	return color.HiRedString("%v", value)
}
