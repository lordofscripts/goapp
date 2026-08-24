//go:build zlog && develop

/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package zlog

import (
	"fmt"
	"os"
	"strings"

	"github.com/lordofscripts/goapp/app/mtag"
	"github.com/lordofscripts/goapp/osx"
)

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/

const (
	UC_CLICK    rune = '↯'
	UC_CHANGE   rune = '🢱'
	UC_VALIDATE rune = '🢰'
	UC_ENTER    rune = '❯'
	UC_LEAVE    rune = '❮'
)

/* ----------------------------------------------------------------
 *                     I N T E R F A C E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                         T Y P E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                   P U B L I C    T Y P E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                   P R I V A T E    T Y P E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                    C O N S T R U C T O R S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                        M E T H O D S
 *-----------------------------------------------------------------*/

//	****	implements ILogPrioritized 	****

// Trace level with variadic parameters
// NOTE: It is conditioned to "-tag develop" being set, else NO-OP.
func (blg *FlexLogger) Trace(v ...any) {
	if blg.minLogLevel <= LevelTrace {
		v1 := append([]any{tagTRACE}, v...)
		blg.ilogger.Print(v1...)
	}
}

// Trace level with format string
// NOTE: It is conditioned to "-tag develop" being set, else NO-OP.
func (blg *FlexLogger) Tracef(format string, v ...any) {
	if blg.minLogLevel <= LevelTrace {
		blg.ilogger.Printf(tagTRACE+format, v...)
	}
}

// Debug level with variadic parameters
// NOTE: It is conditioned to "-tag develop" being set, else NO-OP.
func (blg *FlexLogger) Debug(v ...any) {
	if blg.minLogLevel <= LevelDebug {
		v1 := append([]any{tagDEBUG}, v...)
		blg.ilogger.Print(v1...)
	}
}

// Debug level with format string
// NOTE: It is conditioned to "-tag develop" being set, else NO-OP.
func (blg *FlexLogger) Debugf(format string, v ...any) {
	if blg.minLogLevel <= LevelDebug {
		blg.ilogger.Printf(tagDEBUG+format, v...)
	}
}

// Information level with variadic parameters
// NOTE: It is conditioned to "-tag develop" being set, else NO-OP.
func (blg *FlexLogger) Info(v ...any) {
	if blg.minLogLevel <= LevelInfo {
		v1 := append([]any{tagINFO}, v...)
		blg.ilogger.Print(v1...)
	}
}

// Information level with format string
// NOTE: It is conditioned to "-tag develop" being set, else NO-OP.
func (blg *FlexLogger) Infof(format string, v ...any) {
	if blg.minLogLevel <= LevelInfo {
		blg.ilogger.Printf(tagINFO+format, v...)
	}
}

// (Log Level: WARN) Warning level with variadic parameters
func (blg *FlexLogger) Warn(v ...any) {
	if blg.minLogLevel <= LevelWarning {
		v1 := append([]any{tagWARN}, v...)
		blg.ilogger.Print(v1...)
	}
}

// (Log Level: WARN) Warning level with format string
func (blg *FlexLogger) Warnf(format string, v ...any) {
	if blg.minLogLevel <= LevelWarning {
		blg.ilogger.Printf(tagWARN+format, v...)
	}
}

// (Log Level: ERROR) Error level with variadic parameters
func (blg *FlexLogger) Error(v ...any) {
	if blg.minLogLevel <= LevelError {
		v1 := append([]any{tagERROR}, v...)
		blg.ilogger.Print(v1...)
	}
}

// (Log Level: ERROR) Error level with format string
func (blg *FlexLogger) Errorf(format string, v ...any) {
	if blg.minLogLevel <= LevelError {
		blg.ilogger.Printf(tagERROR+format, v...)
	}
}

// (Log Level: ERROR) Error level limited to the error itself
func (blg *FlexLogger) ErrorE(err error) {
	if blg.minLogLevel <= LevelError {
		blg.ilogger.Println(tagERROR, err.Error())
	}
}

// (Log Level: FATAL) Fatal level with variadic parameters and exitCode
// for terminating the application.
func (blg *FlexLogger) Fatal(exitCode int, v ...any) {
	if blg.minLogLevel <= LevelFatal {
		v1 := append([]any{tagFATAL}, v...)
		blg.ilogger.Print(v1...)
	}

	os.Exit(exitCode)
}

// (Log Level: FATAL) Fatal level with format string and exitCode for terminating
// the application.
func (blg *FlexLogger) Fatalf(exitCode int, format string, v ...any) {
	if blg.minLogLevel <= LevelFatal {
		blg.ilogger.Printf(tagFATAL+format, v...)
	}

	os.Exit(exitCode)
}

//	****	implements ILogTaggable 	****

// Tagged level-aware TRACE message.
// NOTE: It is conditioned to "-tag develop" being set, else NO-OP.
func (blg *FlexLogger) TraceT(message string, v ...mtag.ILogKeyValuePair) {
	if blg.minLogLevel <= LevelTrace {
		var sb strings.Builder
		sb.WriteString(tagTRACE)
		sb.WriteString(message)
		for _, t := range v {
			sb.WriteString(" ")
			sb.WriteString(t.String())
		}
		blg.ilogger.Print(sb.String())
	}
}

// Tagged level-aware DEBUG message.
// NOTE: It is conditioned to "-tag develop" being set, else NO-OP.
func (blg *FlexLogger) DebugT(message string, v ...mtag.ILogKeyValuePair) {
	if blg.minLogLevel <= LevelDebug {
		var sb strings.Builder
		sb.WriteString(tagDEBUG)
		sb.WriteString(message)
		for _, t := range v {
			sb.WriteString(" ")
			sb.WriteString(t.String())
		}
		blg.ilogger.Print(sb.String())
	}
}

// Tagged level-aware INFO message.
// NOTE: It is conditioned to "-tag develop" being set, else NO-OP.
func (blg *FlexLogger) InfoT(message string, v ...mtag.ILogKeyValuePair) {
	if blg.minLogLevel <= LevelInfo {
		var sb strings.Builder
		sb.WriteString(tagINFO)
		sb.WriteString(message)
		for _, t := range v {
			sb.WriteString(" ")
			sb.WriteString(t.String())
		}
		blg.ilogger.Print(sb.String())
	}
}

// @note these output as log as mlog is compiled regardless of develop

// (Log Level: WARN) Warning level with message and variadic MLog tags.
// NOTE: It is always operational regardless of "-tag develop"
func (blg *FlexLogger) WarnT(message string, v ...mtag.ILogKeyValuePair) {
	if blg.minLogLevel <= LevelWarning {
		var sb strings.Builder
		sb.WriteString(tagWARN)
		sb.WriteString(message)
		for _, t := range v {
			sb.WriteString(" ")
			sb.WriteString(t.String())
		}
		blg.ilogger.Print(sb.String())
	}
}

// (Log Level: ERROR) Error level with message and variadic MLog tags.
// NOTE: It is always operational regardless of "-tag develop"
func (blg *FlexLogger) ErrorT(message string, v ...mtag.ILogKeyValuePair) {
	if blg.minLogLevel <= LevelError {
		var sb strings.Builder
		sb.WriteString(tagERROR)
		sb.WriteString(message)
		for _, t := range v {
			sb.WriteString(" ")
			sb.WriteString(t.String())
		}
		blg.ilogger.Print(sb.String())
	}
}

// (Log Level: FATAL) Fatal level with message and variadic MLog tags.
// it terminates execution with exitCode.
// NOTE: It is always operational regardless of "-tag develop"
func (blg *FlexLogger) FatalT(exitCode int, message string, v ...mtag.ILogKeyValuePair) {
	if blg.minLogLevel <= LevelFatal {
		var sb strings.Builder
		sb.WriteString(tagFATAL)
		sb.WriteString(message)
		for _, t := range v {
			sb.WriteString(" ")
			sb.WriteString(t.String())
		}
		blg.ilogger.Print(sb.String())
	}

	os.Exit(exitCode)
}

//	****	implements ILogCatheterable	****
//			@note These are conditional to -tag develop

// NOTE: It is conditional to "-tag develop" being set, else NO-OP.
func (blg *FlexLogger) SetCatheterFile(filename string) bool {
	var err error = nil
	if blg.catFile != nil {
		return false
	}
	blg.catFile, err = blg.openLogFile(filename, false, true)

	return err == nil
}

// NOTE: It is conditional to "-tag develop" being set, else NO-OP.
func (blg *FlexLogger) PrintCatheter(message string, v ...mtag.ILogKeyValuePair) {
	if blg.catFile != nil {
		var sb strings.Builder
		sb.WriteString(tagCATHE)
		sb.WriteString(message)
		for _, t := range v {
			sb.WriteString(" ")
			sb.WriteString(t.String())
		}

		blg.catFile.WriteString(sb.String() + "\n")
	}
}

//	****	implements ILogEventDriven 	****

// Use a separate Call Tree log file for extended (ILogEventDriven) logging.
// The default is the same as the regular log output.
// NOTE: It is conditional to "-tag develop" being set, else NO-OP (returns nil).
// - The file is created or truncated.
func (blg *FlexLogger) WithCallTree(filename string) error {
	var err error
	if blg.treeFile, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644); err == nil {
		blg.treeFile.WriteString(_CALLTREE_FILE_LEADER)
		//Greet(blg.treeFile)
	} else {
		return err
	}

	return nil
}

// For use to log that an object constructor is executing.
// NOTE: Only operational with "-tag develop", else NO-OP.
func (blg *FlexLogger) Ctor() {
	if loc, _, allowed := getCallerFlexFiltered(2, blg.filterSvc); allowed {
		blg.ilogger.Print("⟫ (Ctor) ", loc)
	}
}

// For use at entering a function/method of an EVENT callback.
// NOTE: Only operational with "-tag develop", else NO-OP.
func (blg *FlexLogger) EventEnter() {
	if loc, _, allowed := getCallerFlexFiltered(2, blg.filterSvc); allowed {
		blg.ilogger.Print("❯ (Event) ", loc)

		nl, nm := GetNestingLevel(3)
		blg.writeCallTree(fmt.Sprintf("%3d ENT %c %s%s", nl, UC_ENTER, strings.Repeat(" ", nl), nm))
	}
}

// For use at the end of a function/method of an EVENT callback.
// NOTE: Only operational with "-tag develop", else NO-OP.
func (blg *FlexLogger) EventLeave() {
	if loc, _, allowed := getCallerFlexFiltered(2, blg.filterSvc); allowed {
		blg.ilogger.Print("❮ (Event) ", loc)

		nl, nm := GetNestingLevel(3)
		blg.writeCallTree(fmt.Sprintf("%3d LEA %c %s%s", nl, UC_LEAVE, strings.Repeat(" ", nl), nm))
	}
}

// For use at entering a function/method.
// NOTE: Only operational with "-tag develop", else NO-OP.
func (blg *FlexLogger) Enter() {
	if loc, _, allowed := getCallerFlexFiltered(2, blg.filterSvc); allowed {
		blg.ilogger.Print("❯ ", loc)

		nl, nm := GetNestingLevel(3)
		blg.writeCallTree(fmt.Sprintf("%3d ENT %c %s%s", nl, UC_ENTER, strings.Repeat(" ", nl), nm))
	}
}

// For use at the end of a function/method.
// NOTE: Only operational with "-tag develop", else NO-OP.
func (blg *FlexLogger) Leave() {
	if loc, _, allowed := getCallerFlexFiltered(2, blg.filterSvc); allowed {
		blg.ilogger.Print("❮ ", loc)

		nl, nm := GetNestingLevel(3)
		blg.writeCallTree(fmt.Sprintf("%3d LEA %c %s%s", nl, UC_LEAVE, strings.Repeat(" ", nl), nm))
	}
}

// Just logging a visit without Enter/Leave.
// NOTE: Only operational with "-tag develop", else NO-OP.
func (blg *FlexLogger) Visit() {
	if loc, _, allowed := getCallerFlexFiltered(2, blg.filterSvc); allowed {
		blg.ilogger.Print("❮❯ ", loc)
	}
}

// A procedural step within a code block/function/method.
// NOTE: Only operational with "-tag develop", else NO-OP.
func (blg *FlexLogger) Step(message string) {
	if _, _, allowed := getCallerFlexFiltered(2, blg.filterSvc); allowed {
		blg.ilogger.Printf("%c %s", UC_FOOTSTEPS, message)
	}
}

// a function/method result of some kind.
// NOTE: Only operational with "-tag develop", else NO-OP.
func (blg *FlexLogger) Result(format string, v ...any) {
	if loc, _, allowed := getCallerFlexFiltered(2, blg.filterSvc); allowed {
		blg.ilogger.Printf(string(UC_ARROWS3)+" "+loc+" "+format, v...)
	}
}

// a validation callback.
// NOTE: Only operational with "-tag develop", else NO-OP.
func (blg *FlexLogger) OnValidating() {
	nl, nm := GetNestingLevel(3)
	blg.ilogger.Printf("%c VAL (%d) %s", UC_VALIDATE, nl, nm)
	blg.writeCallTree(fmt.Sprintf("%3d VAL %c %s%s", nl, UC_VALIDATE, strings.Repeat(" ", nl), nm))
}

// a value changed callback.
// NOTE: Only operational with "-tag develop", else NO-OP.
func (blg *FlexLogger) OnChanged(toValue ...any) {
	nl, nm := GetNestingLevel(3)
	var valueStr string = ""
	if len(toValue) == 1 {
		valueStr = fmt.Sprintf(" %c %v", UC_ARROWS3, toValue)
	}
	blg.ilogger.Printf("%c CHG (%d) %s", UC_CHANGE, nl, nm)
	blg.writeCallTree(fmt.Sprintf("%3d CHG %c %s%s%s", nl, UC_CHANGE, strings.Repeat(" ", nl), nm, valueStr))
}

// a value change as a result of another trigger (?).
// NOTE: Only operational with "-tag develop", else NO-OP.
func (blg *FlexLogger) OnUpdate() {
	nl, nm := GetNestingLevel(3)
	blg.ilogger.Printf("%c UPD (%d) %s", UC_OBSERVER, nl, nm)
	blg.writeCallTree(fmt.Sprintf("%3d UPD %c %s%s", nl, UC_OBSERVER, strings.Repeat(" ", nl), nm))
}

// ?
// NOTE: Only operational with "-tag develop", else NO-OP.
func (blg *FlexLogger) OnCascade(to string, val any) {
	nl, nm := GetNestingLevel(3)
	blg.ilogger.Printf("%c SET (%d) %s TO %s = %v", UC_ARROWS3, nl, nm, to, val)
	blg.writeCallTree(fmt.Sprintf("%3d SET %c %s%s %s 🟰 %v", nl, UC_ARROWS3, strings.Repeat(" ", nl), nm, to, val))
}

// User Interface widget clicked event.
// NOTE: Only operational with "-tag develop", else NO-OP.
func (blg *FlexLogger) OnClick(widg ...mtag.ILogKeyValuePair) {
	nl, nm := GetNestingLevel(3)
	if len(widg) > 0 {
		blg.ilogger.Printf("%c CLK (%d) %s %s", UC_CLICK, nl, nm, widg[0].String())
		blg.writeCallTree(fmt.Sprintf("%3d CLK %c %s%s", nl, UC_CLICK, strings.Repeat(" ", nl), nm))
	} else {
		blg.ilogger.Printf("%c CLK (%d) %s %s", UC_CLICK, nl, nm, widg[0].String())
		blg.writeCallTree(fmt.Sprintf("%3d CLK %c %s%s", nl, UC_CLICK, strings.Repeat(" ", nl), nm))
	}
}

// Writes to the call tree (separate file or default output).
// The data string is expected WITHOUT line feed ("\n") as
// it will be appended here.
func (blg *FlexLogger) writeCallTree(data string) {
	data = data + osx.EOL
	if blg.treeFile == nil {
		blg.ilogger.Print(data)
	} else {
		blg.treeFile.WriteString(data)
	}
}

/* ----------------------------------------------------------------
 *                 P U B L I C    M E T H O D S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                 P R I V A T E    M E T H O D S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

// The current build of the ZLog library (development|production)
func GetMode() string {
	return fmt.Sprintf("%s   ZLog (debug)", UC_COG_GEAR_COLOR)
}
