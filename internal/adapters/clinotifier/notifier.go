package clinotifier

import (
	"fmt"

	"github.com/fatih/color"
)

var isSilent bool
var section bool

type CLINotifier struct {
}

func NewCLINotifier() CLINotifier {
	return CLINotifier{}
}

func (c *CLINotifier) formaValue(value any) string {
	return color.CyanString("%v", value)
}

func (c *CLINotifier) formatError(value any) string {
	return color.HiRedString("%v", value)
}

func (c *CLINotifier) format(values ...any) []any {
	formattedValues := []any{}
	for _, v := range values {
		if _, isError := v.(error); isError == true {
			formattedValues = append(formattedValues, c.formatError(v))
		} else {
			formattedValues = append(formattedValues, c.formaValue(v))
		}
	}

	return formattedValues
}

func (c *CLINotifier) NotifyLoud(values ...any) {
	section = false
	fmt.Println(values...)
}

func (c *CLINotifier) NotifyfLoud(format string, values ...any) {
	section = false
	fmt.Printf(format+"\n", c.format(values...)...)
}

func (c *CLINotifier) Notifyf(format string, values ...any) {
	if isSilent == false {
		section = false
		fmt.Printf(format+"\n", c.format(values...)...)
	}
}

func (c *CLINotifier) Notify(values ...any) {

	var isError bool
	section = false

	for _, v := range values {
		if _, isError = v.(error); isError == true {
			//the list of values contains at least one error
			break
		}
	}

	if isSilent == false || isError {
		fmt.Println(values...)
	}
}

func (c *CLINotifier) Section() {
	if section == false && isSilent == false {
		section = true
		fmt.Println()
	}
}

func (c *CLINotifier) SectionLoud() {
	if section == false {
		section = true
		fmt.Println()
	}
}

func (c *CLINotifier) SetSilent(silent bool) {
	isSilent = silent
}

func (c *CLINotifier) Silent() bool {
	return isSilent
}
