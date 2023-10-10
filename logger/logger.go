// Package logger is a package that is used to print logs, err messages
// to the stdout and validate structs with the help of go playground validator
package logger

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/VinukaThejana/go-utils/text"
	"github.com/charmbracelet/lipgloss"
	"github.com/go-playground/validator/v10"
)

// Error is a function that is used to log the error messages
func Error(err error) {
	now := time.Now()

	text.Text{}.ErrorWithPadding(
		text.P{
			Bottom: 0,
		},
		fmt.Sprintf(
			"[%d:%s:%d : [%d:%d:%d]] : %s\n",
			now.Year(),
			now.Month().String(),
			now.Day(),
			now.Hour(),
			now.Minute(),
			now.Second(),
			err.Error(),
		),
	)
}

// ErrorWithMsg is a function tha is used to log the error messages
// with custom messages
func ErrorWithMsg(err error, msg string) {
	now := time.Now()

	text.Text{}.ErrorWithPadding(
		text.P{
			Bottom: 0,
		},
		fmt.Sprintf(
			"[%d:%s:%d : [%d:%d:%d]] : %s\n%s\n",
			now.Year(),
			now.Month().String(),
			now.Day(),
			now.Hour(),
			now.Minute(),
			now.Second(),
			err.Error(),
			msg,
		),
	)
}

// Errorf is a function that gives the ability to log error messages to the stdout
// with exiting the program with error code 1
func Errorf(err error) {
	now := time.Now()

	text.Text{}.ErrorWithPadding(
		text.P{
			Top:    0,
			Bottom: 0,
			Right:  0,
			Left:   0,
		},
		fmt.Sprintf(
			"[%d:%s:%d : [%d:%d:%d]] : %s",
			now.Year(),
			now.Month().String(),
			now.Day(),
			now.Hour(),
			now.Minute(),
			now.Second(),
			err.Error(),
		),
	)

	text.Text{}.ErrorWithPadding(
		text.P{
			Top:    0,
			Bottom: 0,
			Right:  0,
			Left:   0,
		},
		fmt.Sprintf("\n STACKTRACE \n\n%s\n", string(debug.Stack())),
	)

	os.Exit(1)
}

// ErrorfWithMsg is a function to log the error message with the custom message
// to the stdout while exiting the program with error code 1
func ErrorfWithMsg(err error, msg string) {
	now := time.Now()

	text.Text{}.ErrorWithPadding(
		text.P{},
		fmt.Sprintf(
			"[%d:%s:%d : [%d:%d:%d]] : %s\n%s\n",
			now.Year(),
			now.Month().String(),
			now.Day(),
			now.Hour(),
			now.Minute(),
			now.Second(),
			err.Error(),
			msg,
		),
	)

	text.Text{}.ErrorWithPadding(
		text.P{
			Top:    0,
			Bottom: 0,
			Right:  0,
			Left:   0,
		},
		fmt.Sprintf("\n STACKTRACE \n%s\n", string(debug.Stack())),
	)
}

// Log is a function that is used to log a message to the stdout
func Log(msg string) {
	now := time.Now()

	fmt.Println(text.Text{}.P(text.Style{
		Color:   lipgloss.Color("#ffffff"),
		Padding: text.P{},
		Align:   lipgloss.Left,
		Bold:    false,
	},
		fmt.Sprintf(
			"[%d:%s:%d : [%d:%d:%d]] : %s\n",
			now.Year(),
			now.Month().String(),
			now.Day(),
			now.Hour(),
			now.Minute(),
			now.Second(),
			msg,
		),
	),
	)
}

// LogWithStyling is a function that is used to log the messages
// to the stdout with styling
func LogWithStyling(msg string, style text.Style) {
	now := time.Now()
	switch {
	case style.Align == 0:
		style.Align = lipgloss.Left
		fallthrough
	case style.Color == nil:
		style.Color = lipgloss.Color("#ffffff")
	}

	fmt.Println(text.Text{}.P(text.Style{
		Color:   style.Color,
		Padding: style.Padding,
		Align:   style.Align,
		Bold:    style.Bold,
	},
		fmt.Sprintf(
			"[%d:%s:%d : [%d:%d:%d]] : %s\n",
			now.Year(),
			now.Month().String(),
			now.Day(),
			now.Hour(),
			now.Minute(),
			now.Second(),
			msg,
		),
	),
	)
}

// Validate is a function that is used to validate wether a given struct satisfies
// a given condition
func Validate(s interface{}) (isValid bool, err error) {
	v := validator.New()
	err = v.Struct(s)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Validatef is a function that is used to validate wether a given struct satisfies
// a given condition, if the condition is not met this will cause the program to panic
// the validation
func Validatef(s interface{}) {
	v := validator.New()
	err := v.Struct(s)
	if err == nil {
		return
	}

	Errorf(err)
}
