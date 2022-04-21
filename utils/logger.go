package utils

import (
	"os"

	"github.com/pterm/pterm"
)

type logger struct {
	VerboseEnabled bool
	DebugEnabled bool
}

var Log = newLogger()

func newLogger() logger {
	pterm.PrintDebugMessages = true
	
	return logger{
		VerboseEnabled: false,
		DebugEnabled: false,
	}
}

func (logger *logger) Success(format string, a ...interface{}) {
	pterm.Success.Printfln(format, a...)
}

func (logger *logger) Warning(format string, a ...interface{}) {
	pterm.Warning.Printfln(format, a...)
}

func (logger *logger) Error(format string, a ...interface{}) {
	pterm.Error.Printfln(format, a...)
}

func (logger *logger) Fatal(format string, a ...interface{}) {
	f := pterm.Fatal
	f.Fatal = logger.DebugEnabled

	f.Printfln(format, a...)
	os.Exit(1)
}

func (logger *logger) FatalNoExit(format string, a ...interface{}) {
	f := pterm.Fatal
	f.Fatal = false

	f.Printfln(format, a...)
}

func (logger *logger) Info(format string, a ...interface{}) {
	// if logger.VerboseEnabled {
		pterm.Success.Printfln(format, a...)
	// }
}

func (logger *logger) Debug(format string, a ...interface{}) {
	if logger.DebugEnabled {
		pterm.Debug.Printfln(format, a...)
	}
}