//go:build !mlog

/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Lord of Scripts
 *							   goApp
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * MLog No-Op Build. Fallback when mlog tag is not used.
 *-----------------------------------------------------------------*/
package mlog

import "github.com/lordofscripts/goapp/app/mtag"

// The current build of the MLog library (development|production)
func GetMode() string {
	return "[disabled]   MLog (disabled)"
}

// The "catheter" feature is not available.
func SetCatheterFile(filename string) bool {
	return false
}

// The "catheter" feature is not available.
func PrintCatheter(message string, v ...mtag.ILogKeyValuePair) {}

/* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *			P r i v i l e g e d   L e v e l s
 *- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -*/

func Trace(v ...any)                                    {}
func Tracef(format string, v ...any)                    {}
func TraceT(message string, v ...mtag.ILogKeyValuePair) {}

func Debug(v ...any)                                    {}
func Debugf(format string, v ...any)                    {}
func DebugT(message string, v ...mtag.ILogKeyValuePair) {}

func Info(v ...any)                                    {}
func Infof(format string, v ...any)                    {}
func InfoT(message string, v ...mtag.ILogKeyValuePair) {}

/* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *			E x t e n d e d 	F u n c t i o n s
 *				(For Event-driven apps)
 *- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -*/

func Ctor()       {}
func EventEnter() {}
func EventLeave() {}
func Enter()      {}
func Leave()      {}
func Visit()      {}
func Step(string) {}

func AttentionAlways(short string, err error) {}
func Attention(short string, err error)       {}

func Result(format string, v ...any) {}
func OnValidating()                  {}
func OnChanged(toValue ...any)       {}
func OnUpdate()                      {}
func OnCascade(string, any)          {}
func OnClick()                       {}
