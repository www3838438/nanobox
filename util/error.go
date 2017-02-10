package util

import (
	"flag"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/bugsnag/bugsnag-go"
	"github.com/jcelliott/lumber"
)

type (
	Err struct {
		Message string
		Stack   []string
	}
)

// satisfy the error interface
func (eh Err) Error() string {
	if len(eh.Stack) == 0 {
		return eh.Message
	}
	return fmt.Sprintf("%s: %s", strings.Join(eh.Stack, ": "), eh.Message)
}

// report an issue to bugsnag this has to be done when the error is first created by us
// so it can have a valid stacktrace
func (eh Err) report() {
	// dont report if we are testing
	if flag.Lookup("test.v") != nil {
		return
	}

	bugsnagErr := bugsnag.Notify(eh, bugsnag.User{Id: UniqueID()}, bugsnag.SeverityInfo)
	if bugsnagErr != nil {
		lumber.Error("Bugsnag error: %s", bugsnagErr)
	}

}

// log the error we ran into into our log file
func (eh Err) log() {
	// dont log if we are testing
	if flag.Lookup("test.v") != nil {
		return
	}

	lumber.Error(eh.Error())
	lumber.Error("%s\n", debug.Stack())
}

// Write an error message simular to Printf but logs the error to
// the log file
func ErrorfQuiet(fmtStr string, args ...interface{}) error {
	err := Err{
		Message: fmt.Sprintf(fmtStr, args...),
		Stack:   []string{},
	}
	err.log()
	return err
}

// creates an error the same fmt does but also reports errors to bugsnag
func Errorf(fmt string, args ...interface{}) error {
	eh := ErrorfQuiet(fmt, args...).(Err)
	eh.report()
	return eh
}

// create an error but do not report to bugsnag
func ErrorQuiet(err error) error {
	if err == nil {
		return err
	}

	if er, ok := err.(Err); ok {
		return er
	}

	er := Err{
		Message: err.Error(),
		Stack:   []string{},
	}
	er.log()
	return er
}

// createson of our errors from a external error
func Error(err error) error {
	if err == nil {
		return err
	}

	eh := ErrorQuiet(err).(Err)
	eh.report()
	return eh
}

// prepend the new message to the stack on our error messages
// this is usefull because delimiting stack elements by :
// is not sufficient
func ErrorAppend(err error, fmtStr string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	msg := fmt.Sprintf(fmtStr, args...)

	// if it is one of our errors
	if er, ok := err.(Err); ok {
		er.Stack = append([]string{msg}, er.Stack...)
		return er
	}

	// make sure when we get any new error that isnt ours
	// we log and report it
	return ErrorAppend(Error(err), fmtStr, args...)
}
