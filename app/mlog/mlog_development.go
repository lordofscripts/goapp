//go:build mlog && develop

/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Lord of Scripts
 *							   goApp
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * MLog Development Build. A full-featured version of MLog.
 *-----------------------------------------------------------------*/
package mlog

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"

	"github.com/lordofscripts/goapp/app/mtag"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

const (
	UC_CLICK    rune = '↯'
	UC_CHANGE   rune = '🢱'
	UC_VALIDATE rune = '🢰'
)

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// The current build of the MLog library (development|production)
func GetMode() string {
	return fmt.Sprintf("%s   MLog (development)", UC_COG_GEAR_COLOR)
}

// Create a "catheter" log file. It is a supplementary lifeline for
// exceptional logging and contains no format. Use PrintCathether()
// for writing output.
func SetCatheterFile(filename string) bool {
	var err error = nil
	if catFile != nil {
		return false
	}
	catFile, err = openLogFile(filename, false, true)

	return err == nil
}

// Print to the catheter file.
func PrintCatheter(message string, v ...mtag.ILogKeyValuePair) {
	if catFile != nil {
		var sb strings.Builder
		sb.WriteString(tagCATHE)
		sb.WriteString(message)
		for _, t := range v {
			sb.WriteString(" ")
			sb.WriteString(t.String())
		}

		catFile.WriteString(sb.String() + "\n")
	}
}

/* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *				P r i v i l e g e d   L e v e l s
 *- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -*/

// (Log Level: TRACE) Trace level with variadic parameters
func Trace(v ...any) {
	if minLogLevel <= LevelTrace {
		v1 := append([]any{tagTRACE}, v...)
		ilogger.Print(v1...)
	}
}

// (Log Level: TRACE) Trace level with format string
func Tracef(format string, v ...any) {
	if minLogLevel <= LevelTrace {
		ilogger.Printf(tagTRACE+format, v...)
	}
}

// (Log Level: TRACE) Trace level with message and variadic MLog tags.
func TraceT(message string, v ...mtag.ILogKeyValuePair) {
	if minLogLevel <= LevelTrace {
		var sb strings.Builder
		sb.WriteString(tagTRACE)
		sb.WriteString(message)
		for _, t := range v {
			sb.WriteString(" ")
			sb.WriteString(t.String())
		}
		ilogger.Print(sb.String())
	}
}

// (Log Level: DEBUG) Debug level with variadic parameters
func Debug(v ...any) {
	if minLogLevel <= LevelDebug {
		v1 := append([]any{tagDEBUG}, v...)
		ilogger.Print(v1...)
	}
}

// (Log Level: DEBUG) Debug level with format string
func Debugf(format string, v ...any) {
	if minLogLevel <= LevelDebug {
		ilogger.Printf(tagDEBUG+format, v...)
	}
}

// (Log Level: DEBUG) Debug level with message and variadic MLog tags.
func DebugT(message string, v ...mtag.ILogKeyValuePair) {
	if minLogLevel <= LevelDebug {
		var sb strings.Builder
		sb.WriteString(tagDEBUG)
		sb.WriteString(message)
		for _, t := range v {
			sb.WriteString(" ")
			sb.WriteString(t.String())
		}
		ilogger.Print(sb.String())
	}
}

// (Log Level: INFO) Information level with variadic parameters
func Info(v ...any) {
	if minLogLevel <= LevelInfo {
		v1 := append([]any{tagINFO}, v...)
		ilogger.Print(v1...)
	}
}

// (Log Level: INFO) Information level with format string
func Infof(format string, v ...any) {
	if minLogLevel <= LevelInfo {
		ilogger.Printf(tagINFO+format, v...)
	}
}

// (Log Level: INFO) Information level with message and variadic MLog tags.
func InfoT(message string, v ...mtag.ILogKeyValuePair) {
	if minLogLevel <= LevelInfo {
		var sb strings.Builder
		sb.WriteString(tagINFO)
		sb.WriteString(message)
		for _, t := range v {
			sb.WriteString(" ")
			sb.WriteString(t.String())
		}
		ilogger.Print(sb.String())
	}
}

/* ----------------------------------------------------------------
 *				E x t e n d e d 	F u n c t i o n s
 *					(For Event-driven apps)
 *-----------------------------------------------------------------*/

// (Event-driven App Logging) Log when an object constructor executes.
func Ctor() {
	if loc, _, allowed := getCallerFlexFiltered(2); allowed {
		ilogger.Print("⟫ (Ctor) ", loc)
	}
}

// (Event-driven App Logging) For use at entering a function/method of
// an EVENT callback
func EventEnter() {
	if loc, _, allowed := getCallerFlexFiltered(2); allowed {
		ilogger.Print("❯ (Event) ", loc)
	}
}

// (Event-driven App Logging) For use at the end of a function/method
// of an EVENT callback
func EventLeave() {
	if loc, _, allowed := getCallerFlexFiltered(2); allowed {
		ilogger.Print("❮ (Event) ", loc)
	}
}

// (Event-driven App Logging) For use at entering a function/method
func Enter() {
	if loc, _, allowed := getCallerFlexFiltered(2); allowed {
		ilogger.Print("❯ ", loc)
	}
}

// (Event-driven App Logging) For use at the end of a function/method
func Leave() {
	if loc, _, allowed := getCallerFlexFiltered(2); allowed {
		ilogger.Print("❮ ", loc)
	}
}

// (Event-driven App Logging) For use when entering a function/method
// where you don't want Enter+Leave but just logging your visit.
func Visit() {
	if loc, _, allowed := getCallerFlexFiltered(2); allowed {
		ilogger.Print("❮❯ ", loc)
	}
}

// (Event-driven App Logging) A procedural step within a block
func Step(message string) {
	if _, _, allowed := getCallerFlexFiltered(2); allowed {
		ilogger.Printf("%c %s", UC_FOOTSTEPS, message)
	}
}

// Log an error irrespective of the log level
func AttentionAlways(short string, err error) {
	_, pkgShort, _ := getCallerFlexFiltered(2)
	ilogger.Printf("%c %s ERR %s", UC_EYES, pkgShort, err)
}

// Log an error but conditionally filtered. See AttentionAlways()
func Attention(short string, err error) {
	if _, pkgShort, allowed := getCallerFlexFiltered(2); allowed {
		ilogger.Printf("%c %s ERR %s", UC_EYES, pkgShort, err)
	}
}

// (Event-driven App Logging) To display the result or return value
// of a function/method.
func Result(format string, v ...any) {
	if loc, _, allowed := getCallerFlexFiltered(2); allowed {
		ilogger.Printf(string(UC_ARROWS3)+" "+loc+" "+format, v...)
	}
}

// (Event-driven App Logging) Validator event logging
func OnValidating() {
	nl, nm := GetNestingLevel(3)
	ilogger.Printf("%c VAL (%d) %s", UC_VALIDATE, nl, nm)
	writeCallTree(fmt.Sprintf("%3d VAL %c %s%s", nl, UC_VALIDATE, strings.Repeat(" ", nl), nm))
}

// (Event-driven App Logging) value Changed event logging
func OnChanged(toValue ...any) {
	nl, nm := GetNestingLevel(3)
	var valueStr string = ""
	if len(toValue) == 1 {
		valueStr = fmt.Sprintf(" %c %v", UC_ARROWS3, toValue)
	}
	ilogger.Printf("%c CHG (%d) %s", UC_CHANGE, nl, nm)
	writeCallTree(fmt.Sprintf("%3d CHG %c %s%s%s", nl, UC_CHANGE, strings.Repeat(" ", nl), nm, valueStr))
}

// (Event-driven App Logging) value updated logging
func OnUpdate() {
	nl, nm := GetNestingLevel(3)
	ilogger.Printf("%c UPD (%d) %s", UC_OBSERVER, nl, nm)
	writeCallTree(fmt.Sprintf("%3d UPD %c %s%s", nl, UC_OBSERVER, strings.Repeat(" ", nl), nm))
}

// (Event-driven App Logging) Cascading action logging
func OnCascade(to string, val any) {
	nl, nm := GetNestingLevel(3)
	ilogger.Printf("%c SET (%d) %s TO %s = %v", UC_ARROWS3, nl, nm, to, val)
	writeCallTree(fmt.Sprintf("%3d SET %c %s%s %s 🟰 %v", nl, UC_ARROWS3, strings.Repeat(" ", nl), nm, to, val))
}

// (Event-driven App Logging) UI Click event logging
func OnClick() {
	nl, nm := GetNestingLevel(3)
	ilogger.Printf("%c CLK (%d) %s", UC_CLICK, nl, nm)
	writeCallTree(fmt.Sprintf("%3d CLK %c %s%s", nl, UC_CLICK, strings.Repeat(" ", nl), nm))
}

/* ----------------------------------------------------------------
 *				P r i v a t e	F u n c t i o n s
 *-----------------------------------------------------------------*/

func getCallerFlexFiltered(stackIdx int) (string, string, bool) {
	pc, _, _, ok := runtime.Caller(stackIdx) // PC,file,line,ok
	details := runtime.FuncForPC(pc)

	if ok && details != nil {
		info := details.Name()
		lastSlash := strings.LastIndexByte(info, '/')
		if lastSlash < 0 {
			lastSlash = 0
		}
		lastDot := strings.LastIndexByte(info[lastSlash:], '.') + lastSlash

		//fmt.Printf("INFO %s\n", info)
		//fileS := filepath.Base(filename)
		packageS := info[:lastDot] // in tests it returns 'command-line-*'
		if strings.HasPrefix(packageS, "command-line") {
			packageS = "main"
		}

		var objectS string = ""
		if idx := strings.Index(packageS, ".("); idx != -1 {
			objectS = packageS[idx+1:]
			packageS = packageS[:idx]
			regex := regexp.MustCompile(`\(\*`)
			objectS = regex.ReplaceAllString(objectS, "")
			regex = regexp.MustCompile(`\.\(`)
			objectS = regex.ReplaceAllString(objectS, "")
			regex = regexp.MustCompile(`\)$`)
			objectS = regex.ReplaceAllString(objectS, "") //+ "."
		}

		// Filtered out? don't bother with the rest.
		// @todo implement logx filters in mlog.
		if IsFilteredObject(packageS, objectS) {
			return "", "", false
		}

		// Shortened package name for Log prefix
		var shortPkg = packageS
		if strings.Contains(packageS, "/") {
			shortPkg = packageS[strings.LastIndex(packageS, "/")+1:]
		}

		var funcS string
		if strings.Contains(info, ".init.") {
			funcS = info[:lastDot]
		} else {
			funcS = info[lastDot+1:]
			if objectS != "" {
				funcS = "." + funcS
			}
		}
		//return fmt.Sprintf("@%s.%s()#%d\n", packageS, funcS, line)
		return fmt.Sprintf("%s:%s%s()", shortPkg, objectS, funcS), shortPkg, true
		//fmt.Printf("called from %s\n", details.Name())
		//return filename + "+" + string(line) + " " + packageS + ":" + funcS
	}
	return "", "", false
}
