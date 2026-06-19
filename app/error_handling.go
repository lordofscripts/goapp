/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           go-App
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Utility functions for Error-related handling in a pretty fashion.
 *-----------------------------------------------------------------*/
package app

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/lordofscripts/goapp/app/mlog"
)

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/

// False by default, but when true, every Die*() call also logs
// the message via the `app.mlog` package.
var LogOnDeath bool = false

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

// Death of an application by outputting a good-bye and setting
// the OS exit code. It is logged as fatal.
func Die(message string, exitCode int) {
	errorMsgHeading(message, exitCode, true)
	if LogOnDeath {
		mlog.FatalT(exitCode, message, mlog.YesNo("Died", true), mlog.Int("Code", exitCode))
	} else {
		os.Exit(exitCode)
	}
}

func DieWith(exitCode int, format string, args ...any) {
	errorMsgHeading(fmt.Sprintf(format, args...), exitCode, true)
	os.Exit(exitCode)
}

// display the error and die with an exit code, logging it as Fatal.
func DieWithError(err error, exitCode int) {
	errorHeading(err, exitCode, true)
	if LogOnDeath {
		mlog.FatalT(exitCode, err.Error(), mlog.YesNo("Died", true), mlog.Int("Code", exitCode))
	} else {
		os.Exit(exitCode)
	}
}

// When the condition is met the warning message is printed
func Assert(condition bool, warnMessage string) {
	if condition {
		fmt.Fprintf(os.Stderr, "\n\t%c Assert Failed:\n\t%s\n", UC_RED_EXCLAMATION, warnMessage)
	}
}

// If the condition is met, the death message is printed and the
// application terminates with the exit code.
func AssertOrDie(condition bool, deathMessage string, exitCode int) {
	if condition {
		fmt.Fprintf(os.Stderr, "\n\t%c Assert Failed:", UC_RED_EXCLAMATION)
		errorMsgHeading(deathMessage, exitCode, true)
		os.Exit(exitCode)
	}
}

// prints the error message with the exit code but does NOT exit.
func AnnounceErrorMessage(message string, exitCode int) {
	errorMsgHeading(message, exitCode, false)
}

// prints the error and exit code but does NOT exit the application.
func AnnounceError(err error, exitCode int) {
	errorHeading(err, exitCode, false)
}

/* ----------------------------------------------------------------
 *                 L O C A L   F U N C T I O N S
 *-----------------------------------------------------------------*/

// Pretty formatting of a (wrapped) error. The skeleton line is only
// outputted if isTerminal is true, meaning certain death of the application.
// The exitCode appears after "ERR-" and it automatically detects the
// package name and line number from which it was called.
// 💥 ERROR (ERR-005) 💥
// 🎯 From: main #57
// ·  a sad error message string, one per wrapped error
// 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀
// ------------------------------------------------------------
func errorHeading(err error, exitCode int, isTerminal bool) {
	fmt.Fprintf(os.Stderr, "💥 ERROR (ERR-%03d) 💥\n", exitCode)
	fmt.Fprintf(os.Stderr, "🎯 From: %s\n", callerPackage())
	for errC := err; errC != nil; errC = errors.Unwrap(errC) {
		fmt.Println("· ", errC)
	}
	if isTerminal {
		fmt.Fprintln(os.Stderr, strings.Repeat("💀 ", 20)) //🚫
	}
	fmt.Fprintln(os.Stderr, strings.Repeat("-", 60))
}

// Pretty formatting of a (wrapped) error. The skeleton line is only
// outputted if isTerminal is true, meaning certain death of the application.
// The exitCode appears after "ERR-" and it automatically detects the
// package name and line number from which it was called.
// 💥 ERROR (ERR-005) 💥
// 🎯 From: main #57
// ·  a sad error message string, one per wrapped error
// 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀
// ------------------------------------------------------------
func errorMsgHeading(message string, exitCode int, isTerminal bool) {
	fmt.Fprintf(os.Stderr, "💥 ERROR (ERR-%03d) 💥\n", exitCode)
	fmt.Fprintf(os.Stderr, "🎯 From: %s\n", callerPackage())
	fmt.Println("· ", message)
	if isTerminal {
		fmt.Fprintln(os.Stderr, strings.Repeat("💀 ", 20)) //🚫
	}
	fmt.Fprintln(os.Stderr, strings.Repeat("-", 60))
}

// Returns a string with the caller's package and line number.
func callerPackage() string {
	// PC, File, Line, OK
	const FROM_HERE = 2
	const VIA_DIE = 3
	pc, _, line, ok := runtime.Caller(VIA_DIE) // 0=this func, 1=Caller wrapper, 2=your caller
	if !ok {
		return ""
	}

	fn := runtime.FuncForPC(pc) // e.g. "github.com/me/mod/pkg.(*T).Method"
	if fn == nil {
		return ""
	}
	name := fn.Name()

	// Package path is everything up to the last '.'
	// (works for "path/pkg.Func" and "path/pkg.Type.Method")
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '.' {
			return fmt.Sprintf("%s #%d", name[:i], line)
		}
	}
	return ""
}
