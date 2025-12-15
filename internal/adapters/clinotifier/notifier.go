package clinotifier

import (
	"fmt"

	"github.com/fatih/color"
)

var isSilent bool
var section bool

type UserNotifier struct {
}

func NewUserNotifier() UserNotifier {
	return UserNotifier{}
}

func (c *UserNotifier) formaValue(value any) string {
	return color.CyanString("%v", value)
}

func (c *UserNotifier) formatError(value any) string {
	return color.HiRedString("%v", value)
}

func (c *UserNotifier) format(values ...any) []any {
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

func (c *UserNotifier) NotifyLoud(values ...any) {
	section = false
	fmt.Println(values...)
}

func (c *UserNotifier) NotifyfLoud(format string, values ...any) {
	section = false
	fmt.Printf(format+"\n", c.format(values...)...)
}

func (c *UserNotifier) Notifyf(format string, values ...any) {
	if isSilent == false {
		section = false
		fmt.Printf(format+"\n", c.format(values...)...)
	}
}

func (c *UserNotifier) Notify(values ...any) {

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

func (c *UserNotifier) Section() {
	if section == false && isSilent == false {
		section = true
		fmt.Println()
	}
}

func (c *UserNotifier) SectionLoud() {
	if section == false {
		section = true
		fmt.Println()
	}
}

func (c *UserNotifier) SetSilent(silent bool) {
	isSilent = silent
}

func (c *UserNotifier) Silent() bool {
	return isSilent
}
