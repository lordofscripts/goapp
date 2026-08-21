/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                       goApp::zlog
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * A versatile -yet simple- logging package that expands on the standard
 * `log` adding features like a separate call tree log useful for
 * event-driven applications, log levels as well as type-tagged log levels.
 * It will replace both `mlog` and `logx` by merging their functionalities.
 *-----------------------------------------------------------------*/
package zlog

import (
	"io"

	"github.com/lordofscripts/goapp/app/mtag"
)

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                     I N T E R F A C E S
 *-----------------------------------------------------------------*/

// These are implemented by *os.File and easy to implement by any
// ui-driven application for logging purposes.
type ILoggable interface {
	io.Writer
	io.StringWriter
	io.Closer
}

// Current capabilities of this Log Engine
type IFlexLogger interface {
	ILogConfigurable
	ILogPrioritized
	ILogTaggable
	ILogEventDriven
	ILogFilterable // fine-grained filtering on ILogEventDriven
	ILogCatheterable
}

// A log object that has some configurable elements. Loosely follows
// the standard "log" package
type ILogConfigurable interface {
	io.Writer
	io.Closer
	SetPrefix(prefix string)
	SetOutput(w io.Writer)
	// the same flags used in the standard 'log': Ldate, Ltime, Lmicroseconds,
	// Llongfile, Lshortfile, LUTC, Lmsgprefix & LstdFlags (Ldate|Ltime). In
	// its default configuration we use log.Ldate|log.Ltime|log.Lmsgprefix
	SetFlags(log_flags int)
	Flags() int

	WithWriter(w ILoggable) *FlexLogger   // for standard log
	WithUiWriter(w ILoggable) *FlexLogger // for ILogEventDriven output
	GetMode() string
	Greet()

	AttentionErr(short string, err error)
}

// Interface for a log-level aware log object that implements
// the Trace, Debug, Info, Warning, Error & Fatal levels.
// NOTES:
// - The lowest (most verbose) level is Trace.
// - The highest (least verbose) are Error & Fatal.
// - Warning, Error & Fatal levels are always operational.
// - Trace, Debug & Info only operational with "-tags debug"
type ILogPrioritized interface {
	SetLevel(newLevel LogLevel) LogLevel

	Trace(v ...any)
	Tracef(format string, v ...any)
	Debug(v ...any)
	Debugf(format string, v ...any)
	Info(v ...any)
	Infof(format string, v ...any)

	Warn(v ...any)
	Warnf(format string, v ...any)
	Error(v ...any)
	Errorf(format string, v ...any)
	ErrorE(err error)
	Fatal(exitCode int, v ...any)
	Fatalf(exitCode int, format string, v ...any)
}

// Interface for a log-level-aware log object that supports
// logging using simple messages aggregated with key-value
// pairs as an alternative to message/string formatting.
// NOTES:
// - The lowest (most verbose) level is Trace.
// - The highest (least verbose) are Error & Fatal.
// - Warning, Error & Fatal levels are always operational.
// - Trace, Debug & Info only operational with "-tags debug"
type ILogTaggable interface {
	TraceT(message string, v ...mtag.ILogKeyValuePair)
	DebugT(message string, v ...mtag.ILogKeyValuePair)
	InfoT(message string, v ...mtag.ILogKeyValuePair)

	// Tagged Warning logs a warning and continues.
	// NOTE: It is always operational regardless of "-tag debug"
	WarnT(message string, v ...mtag.ILogKeyValuePair)

	// Tagged Error logs an error and continues.
	// NOTE: It is always operational regardless of "-tag debug"
	ErrorT(message string, v ...mtag.ILogKeyValuePair)

	// Tagged Fatal logs a Fatal error and exits the application.
	// NOTE: It is always operational regardless of "-tag debug"
	FatalT(exitCode int, message string, v ...mtag.ILogKeyValuePair)
}

// Interface for a log object that has extended functions suitable for
// event-driven applications such as Graphical User Interfaces. These
// functions allow for an alternate log in which a call tree with
// cause-effect can be easily seen.
// NOTE: Only operational IF "-tags debug" is set.
type ILogEventDriven interface {
	WithCallTree(filename string) error // The file is created OR truncated.

	Ctor()                          // For use to log that an object constructor is executing
	EventEnter()                    // For use at entering a function/method of an EVENT callback
	EventLeave()                    // For use at the end of a function/method of an EVENT callback
	Enter()                         // For use at entering a function/method
	Leave()                         // For use at the end of a function/method
	Visit()                         // Just logging a visit without Enter/Leave
	Step(string)                    // A procedural step within a code block/function/method
	Result(format string, v ...any) // a function/method result of some kind
	OnValidating()                  // a validation callback
	OnChanged(toValue ...any)       // a value changed callback
	OnUpdate()                      // a value change as a result of another trigger (?)
	OnCascade(string, any)          // ?
	OnClick()                       // User Interface widget clicked event
}

// Interface for a log object that implements filtering
// log output by Module/Package/Object name.
// NOTE: Only operational IF "-tags debug" is set.
type ILogFilterable interface {
	// By default a null filter service does no filtering. Else use
	// NewLogFilterService() to pass here, or set to nil to revert.
	// It automatically calls the filter's LoadFilter() method.
	UseFilter(ILogFilterService) error
}

// Interface for a log object that has a Catheter file. That file
// is an alternate output file for data other than logging or
// end-user output.
// NOTE: Only operational IF "-tags debug" is set.
type ILogCatheterable interface {
	SetCatheterFile(filename string) bool
	PrintCatheter(message string, v ...mtag.ILogKeyValuePair)
}

/* ----------------------------------------------------------------
 *                  M A I N    |    D E M O
 *-----------------------------------------------------------------*/

/*
import (
	_ "github.com/lordofscripts/goapp/app/mtag"
	"github.com/lordofscripts/goapp/app/zlog"
)

// declare once for use throughout your package
var ß zlog.IFlexLogEngine

// Skeleton for logging to a Fyne text widget
var _ zlog.ILoggable = (*WidgetLogger)(nil)

type WidgetLogger struct {
}

// implement io.Writer
func (wl *WidgetLogger) Write(p []byte) (n int, e error) {
	fyne.Do(func() {
		str = string(p)
		// append text to your log widget (multi-line Entry)
	})
}

// implement io.StringWriter
func (wl *WidgetLogger) WriteString(str string) (n int, e error) {
	return wl.Write([]byte(str))
}

// implement io.Closer
func (wl *WidgetLogger) Close() error {
	return nil
}

// Early package-wide initialization
func init() {
	//ß = NewMLogger() // logs to StdErr console as usual
	//ß.SetOutput(os.Stdout) // if you dislike stderr
	ß = NewMLoggerWithFile("/tmp/logdemo.log")
	ß.WithUiWriter(&WidgetLogger{}) // for a Fyne or other GUI application
	ß.SetLevel(LevelDebug)
	ß.SetPrefix("DEMO")
}

func oneFunction() {
	ß.Enter()
	defer ß.Leave()

	ß.InfoT("currently at", At())
	a := 5 * 7
	ß.Result("A is ", a)
}

func DemoFlexLogger() {
	defer ß.Close()

	oneFunction()
	var byVal byte = 5
	var iVal int = 800
	var boVal bool = false
	var stVal string = "Thiß Über Alles"
	var rVal rune = '@'
	var err error = fmt.Errorf("oops, bad thing happened.")

	ß.Trace("this", 500, "is")
	ß.Debugf("This value is %t", true)
	ß.InfoT("tagged logging", Byte("ByteVal", byVal), Rune("Letter", rVal))
	ß.WarnT("a warning", Int("IntVal", iVal), String("A string", stVal))
	ß.Error("a free format error", 583, 65.3, "variadic")
	ß.ErrorE(err)
	ß.ErrorT("invalid value", Bool("Condition", boVal))
	ß.Errorf("or the usual %d format", 700)
	ß.ErrorT("tagged as well", Rune("What", rVal), Err(err))
	ß.Fatal(1, "have to go") // also Fatalf and FatalT available
}
*/
