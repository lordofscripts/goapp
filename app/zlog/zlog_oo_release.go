//go:build zlog && !debug

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
)

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/

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
// NOTE: It is conditioned to "-tag debug" being set, else NO-OP.
func (blg *FlexLogger) Trace(v ...any) {}

// Trace level with format string
// NOTE: It is conditioned to "-tag debug" being set, else NO-OP.
func (blg *FlexLogger) Tracef(format string, v ...any) {}

// Debug level with variadic parameters
// NOTE: It is conditioned to "-tag debug" being set, else NO-OP.
func (blg *FlexLogger) Debug(v ...any) {}

// Debug level with format string
// NOTE: It is conditioned to "-tag debug" being set, else NO-OP.
func (blg *FlexLogger) Debugf(format string, v ...any) {}

// Information level with variadic parameters
// NOTE: It is conditioned to "-tag debug" being set, else NO-OP.
func (blg *FlexLogger) Info(v ...any) {}

// Information level with format string
// NOTE: It is conditioned to "-tag debug" being set, else NO-OP.
func (blg *FlexLogger) Infof(format string, v ...any) {}

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
// NOTE: It is conditioned to "-tag debug" being set, else NO-OP.
func (blg *FlexLogger) TraceT(message string, v ...mtag.ILogKeyValuePair) {

}

// Tagged level-aware DEBUG message.
// NOTE: It is conditioned to "-tag debug" being set, else NO-OP.
func (blg *FlexLogger) DebugT(message string, v ...mtag.ILogKeyValuePair) {

}

// Tagged level-aware INFO message.
// NOTE: It is conditioned to "-tag debug" being set, else NO-OP.
func (blg *FlexLogger) InfoT(message string, v ...mtag.ILogKeyValuePair) {

}

// @note these output as log as mlog is compiled regardless of debug

// (Log Level: WARN) Warning level with message and variadic MLog tags.
// NOTE: It is always operational regardless of "-tag debug"
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
// NOTE: It is always operational regardless of "-tag debug"
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
// NOTE: It is always operational regardless of "-tag debug"
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
//			@note These are conditional to -tag debug

// NOTE: It is conditional to "-tag debug" being set, else NO-OP.
func (blg *FlexLogger) SetCatheterFile(filename string) bool {
	return false
}

// NOTE: It is conditional to "-tag debug" being set, else NO-OP.
func (blg *FlexLogger) PrintCatheter(message string, v ...mtag.ILogKeyValuePair) {

}

//	****	implements ILogEventDriven 	****

// Use a separate Call Tree log file for extended (ILogEventDriven) logging.
// The default is the same as the regular log output.
// NOTE: It is conditional to "-tag debug" being set, else NO-OP (returns nil).
func (blg *FlexLogger) WithCallTree(filename string) error {
	// Greet()
	return nil
}

// For use to log that an object constructor is executing.
// NOTE: Only operational with "-tag debug", else NO-OP.
func (blg *FlexLogger) Ctor() {}

// For use at entering a function/method of an EVENT callback.
// NOTE: Only operational with "-tag debug", else NO-OP.
func (blg *FlexLogger) EventEnter() {}

// For use at the end of a function/method of an EVENT callback.
// NOTE: Only operational with "-tag debug", else NO-OP.
func (blg *FlexLogger) EventLeave() {}

// For use at entering a function/method.
// NOTE: Only operational with "-tag debug", else NO-OP.
func (blg *FlexLogger) Enter() {}

// For use at the end of a function/method.
// NOTE: Only operational with "-tag debug", else NO-OP.
func (blg *FlexLogger) Leave() {}

// Just logging a visit without Enter/Leave.
// NOTE: Only operational with "-tag debug", else NO-OP.
func (blg *FlexLogger) Visit() {}

// A procedural step within a code block/function/method.
// NOTE: Only operational with "-tag debug", else NO-OP.
func (blg *FlexLogger) Step(string) {}

// a function/method result of some kind.
// NOTE: Only operational with "-tag debug", else NO-OP.
func (blg *FlexLogger) Result(format string, v ...any) {}

// a validation callback.
// NOTE: Only operational with "-tag debug", else NO-OP.
func (blg *FlexLogger) OnValidating() {}

// a value changed callback.
// NOTE: Only operational with "-tag debug", else NO-OP.
func (blg *FlexLogger) OnChanged(toValue ...any) {}

// a value change as a result of another trigger (?).
// NOTE: Only operational with "-tag debug", else NO-OP.
func (blg *FlexLogger) OnUpdate() {}

// ?
// NOTE: Only operational with "-tag debug", else NO-OP.
func (blg *FlexLogger) OnCascade(string, any) {}

// User Interface widget clicked event.
// NOTE: Only operational with "-tag debug", else NO-OP.
func (blg *FlexLogger) OnClick(...mtag.ILogKeyValuePair) {}

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
	return fmt.Sprintf("%s   ZLog (production)", UC_COG_GEAR_COLOR)
}
