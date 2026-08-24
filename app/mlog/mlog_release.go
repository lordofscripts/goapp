//go:build mlog && !develop

/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Lord of Scripts
 *							   goApp
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * MLog Release Build. A stripped-off version of MLog which only
 * logs errors, fatals, and warnings with everything else disabled.
 *-----------------------------------------------------------------*/
package mlog

import (
	"fmt"

	"github.com/lordofscripts/goapp/app/mtag"
)

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// The current build of the MLog library (development|production)
func GetMode() string {
	return fmt.Sprintf("%s   MLog (production)", UC_COG_GEAR_COLOR)
}

// The "catheter" feature is not enabled.
func SetCatheterFile(filename string) bool {
	return false
}

// The "catheter" feature is not enabled.
func PrintCatheter(message string, v ...mtag.ILogKeyValuePair) {}

/* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *				P r i v i l e g e d   L e v e l s
 *- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -*/

// Trace level with variadic parameters
func Trace(v ...any) {}

// Trace level with format string
func Tracef(format string, v ...any) {}

// Trace level with message and variadic MLog tags.
func TraceT(message string, v ...mtag.ILogKeyValuePair) {}

// Debug level with variadic parameters
func Debug(v ...any) {}

// Debug level with format string
func Debugf(format string, v ...any) {}

// Debug level with message and variadic MLog tags.
func DebugT(message string, v ...mtag.ILogKeyValuePair) {}

// Information level with variadic parameters
func Info(v ...any) {}

// Information level with format string
func Infof(format string, v ...any) {}

// Information level with message and variadic MLog tags.
func InfoT(message string, v ...mtag.ILogKeyValuePair) {}

/* ----------------------------------------------------------------
 *				E x t e n d e d 	F u n c t i o n s
 *					(For Event-driven apps)
 *-----------------------------------------------------------------*/

/**
 * For use to log that an object constructor is executing
 */
func Ctor() {}

/**
 * For use at entering a function/method of an EVENT callback
 */
func EventEnter() {}

/**
 * For use at the end of a function/method of an EVENT callback
 */
func EventLeave() {}

/**
 * For use at entering a function/method
 */
func Enter() {}

/**
 * For use at the end of a function/method
 */
func Leave() {}

/**
 * For use when entering a function/method where you don't want Enter+Leave
 * but just logging your visit.
 */
func Visit() {}

func Step(string) {}

// Log an error irrespective of the log level or filter
func AttentionAlways(short string, err error) {
	ilogger.Printf("%c ERR %s %s", UC_EYES, short, err)
}

// In Release it is the same as AttentionAlways.
// Log an error irrespective of the log level or filter
func Attention(short string, err error) {
	ilogger.Printf("%c ERR %s %s", UC_EYES, short, err)
}

func Result(format string, v ...any) {}

func OnValidating() {}

func OnChanged(toValue ...any) {}

func OnUpdate() {}

func OnCascade(string, any) {}

func OnClick() {}
