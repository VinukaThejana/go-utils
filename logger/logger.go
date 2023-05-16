// Package logger .
// Used to validate structs, log errors and error messages
// in a more log freindly manner
package logger

import (
	"log"

	"github.com/fatih/color"
	"github.com/go-playground/validator/v10"
)

// Logger is a struct that is used fot logging output given
// the message status
type Logger struct{}

// Error is a method on logger to log the error messages to the stdout
func (Logger) Error(err error, msg *string) {
	color.Set(color.FgRed)
	if msg != nil {
		log.Println(*msg)
	}
	log.Println(err)
	color.Unset()
}

// Errorf is a method on logger to log the error messages to the stdout
// while panicing the programming upon the process
func (Logger) Errorf(err error, msg *string) {
	color.Set(color.FgRed)
	if msg != nil {
		log.Println(*msg)
	}
	log.Fatalln(err)
	color.Unset()
}

// Success is a method on logger to log the success messages to the stdout
func (Logger) Success(msg string) {
	color.Set(color.FgGreen)
	log.Println(msg)
	color.Unset()
}

// Validate is used to validate the given struct and notify wehter the validation was
// failed or wether the validation was successful
func (Logger) Validate(s interface{}) bool {
	v := validator.New()
	err := v.Struct(s)
	if err != nil {
		color.Set(color.FgRed)
		log.Println(err)
		color.Unset()
		return false
	}

	return true
}

// Validatef is used to validate the given structs and panic if it fails
// the validation
func (Logger) Validatef(s interface{}) {
	v := validator.New()
	err := v.Struct(s)
	if err != nil {
		color.Set(color.FgRed)
		log.Fatalln(err)
		color.Unset()
	}
	return
}
