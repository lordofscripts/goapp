/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Lord of Scripts
 *							   goApp
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * MLog is an enhanced version of the log/slog package with logging
 * levels, extra log tags such as String, Bool, YesNo, Int, Byte,
 * Rune and At. But MLog is not for structured logging (JSON, etc.)
 * because I like traditional logging messages without clutter.
 *   Log level can be set in an environment variable and an optional
 * log filename too. It also has a supplementary log file called
 * "catheter" which is not formatted.
 *   If a log file is used, the main log is appended whereas the
 * catheter is truncated.
 *   As of v1.4.6 MLog includes the extended functions for Event-driven
 * (GUI) applications as supported by Logx.
 *-----------------------------------------------------------------*/
package mlog

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/lordofscripts/goapp/app/mtag"
	"github.com/lordofscripts/goapp/osx"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	// Logging level enumeration
	LevelTrace LogLevel = iota
	LevelDebug
	LevelInfo
	LevelWarning
	LevelError
	LevelFatal
)

const (
	defaultPrefix string   = ""
	defaultLevel  LogLevel = LevelError

	tagCATHE string = "[CAT] "
	tagTRACE string = "[TRC] "
	tagDEBUG string = "[DBG] "
	tagINFO  string = "[INF] "
	tagWARN  string = "[WRN] "
	tagERROR string = "[ERR] "
	tagFATAL string = "[DIE] "

	// environment variable that overrides the default (Error) Log Level for CaesarX
	LOG_LEVEL_ENV string = "LOG_LEVEL_CX"
	// environment variable that indicates the log output filename for CaesarX (default stderr)
	LOG_FILE_ENV string = "LOG_FILE_CX"

	_LOG_LEADER           string = "[BEG]\t> > > >   T h e   B e g i n n i n g   < < < <\n"
	_LOG_TRAILER          string = "[END]\t> > > >   T h e   E n d   < < < <\n"
	_CALLTREE_FILE_LEADER string = ">>>>> " + " MLog Call Tree <<<<<<\n"
)

const (
	UC_FOOTSTEPS      rune   = rune(0x1f463) // 👣
	UC_EXCLAMATION    rune   = rune(0x2757)  // ❗
	UC_CROSSMARK      rune   = rune(0x274c)  // ❌
	UC_CROSS          rune   = rune(0x2716)  // ✖
	UC_CHECK          rune   = rune(0x1f5f8) // 🗸
	UC_ARROWS3        rune   = rune(0x21f6)  // ⇶
	UC_OBSERVER       rune   = rune(0x23ff)  // ⏿
	UC_EYES           rune   = rune(0x1f440) // 👀
	UC_EYE            rune   = rune(0x1f441) // 👁
	UC_TIMER_RUNNING  rune   = rune(0x23F3)
	UC_COG_GEAR       rune   = rune(0x2699)   // ⚙ (monochrome)
	UC_COG_GEAR_COLOR string = "\u2699\uFE0F" // ⚙️ (emoji)
)

var (
	logMutex    sync.Mutex
	minLogLevel LogLevel    = LevelDebug
	ilogger     *log.Logger = nil
	logFile     *os.File    = nil // the standard log output file
	catFile     *os.File    = nil // the "catheter output file"
	treeFile    *os.File    = nil // the "call tree" file for extended UI functions
	// UTF8 BOM (Byte Order Mark)
	UTF8_BOM []byte = []byte{0xEF, 0xBB, 0xBF}
)

/* ----------------------------------------------------------------
 *				M o d u l e   I n i t i a l i z a t i o n
 *-----------------------------------------------------------------*/

func init() {
	levelString := os.Getenv(LOG_LEVEL_ENV)
	if levelString != "" {
		minLogLevel = parseLevel(levelString)
	} else {
		minLogLevel = defaultLevel
	}

	const CUSTOM_TIME_FORMAT = "2006-01-02 15:04:05"
	cw := newCustomLogWriter(os.Stderr, CUSTOM_TIME_FORMAT)

	//ilogger = log.New(os.Stderr, defaultPrefix, log.Ldate|log.Ltime|log.Lshortfile)
	ilogger = log.New(os.Stderr, defaultPrefix, log.Ldate|log.Ltime|log.Lmsgprefix)
	outputLogFilename := os.Getenv(LOG_FILE_ENV)
	if len(outputLogFilename) != 0 {
		if fd, err := openLogFile(outputLogFilename, true, false); err != nil {
			ilogger.SetOutput(cw) // fallback to stderr
		} else {
			ilogger.SetOutput(fd)
		}
	} else {
		ilogger.SetOutput(cw)
	}
}

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type LogLevel int

type LevelLogger struct { // @audit deprecate not used or move vars here
	*log.Logger
	MinLevel LogLevel
}

type customLogWriter struct {
	writer io.Writer
	format string
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func newCustomLogWriter(w io.Writer, timeStampFormat string) *customLogWriter {
	return &customLogWriter{writer: w, format: timeStampFormat}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (clw *customLogWriter) Write(p []byte) (n int, err error) {
	timestamp := time.Now().Format(clw.format)
	formattedMessage := fmt.Sprintf("%s %s", timestamp, string(p))
	return clw.writer.Write([]byte(formattedMessage))
}

// implements fmt.Stringer for LogLevel returning its name
// without the "Level" prefix.
func (l LogLevel) String() string {
	vmap := map[LogLevel]string{
		LevelTrace:   "Trace",
		LevelDebug:   "Debug",
		LevelInfo:    "Info",
		LevelWarning: "Warn",
		LevelError:   "Error",
		LevelFatal:   "Fatal",
	}

	if v, ok := vmap[l]; ok {
		return v
	}

	return "++InvalidLevel++"
}

// Parse a string as a LogLevel and return it.
// if it fails it returns the current value.
func (l *LogLevel) Parse(str string) (LogLevel, error) {
	str = strings.ToLower(strings.TrimSpace(str))
	tmap := map[string]LogLevel{
		"trace":   LevelTrace,
		"debug":   LevelDebug,
		"info":    LevelInfo,
		"warn":    LevelWarning,
		"warning": LevelWarning,
		"error":   LevelError,
		"fatal":   LevelFatal,
	}

	if t, ok := tmap[str]; ok {
		return t, nil
	} else {
		return *l, fmt.Errorf("could not parse '%s' as LogLevel, returning current value %s", str, *l)
	}
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

func Greet() {
	if ilogger != nil {
		greetLog(ilogger.Writer())
	}
	if logFile != nil {
		greetLog(logFile)
	}
	if catFile != nil {
		greetLog(catFile)
	}
	if treeFile != nil {
		greetLog(treeFile)
	}
}

// the same flags used in the standard 'log': Ldate, Ltime, Lmicroseconds,
// Llongfile, Lshortfile, LUTC, Lmsgprefix & LstdFlags (Ldate|Ltime). In
// its default configuration we use log.Ldate|log.Ltime|log.Lmsgprefix
func SetFlags(log_flags int) {
	ilogger.SetFlags(log_flags)
}

// The log flags that are currently used
func Flags() int {
	return ilogger.Flags()
}

// SetLevel sets the current logging level. Unlike log and slog
// the mlog package supports logging levels.
func SetLevel(newLevel LogLevel) LogLevel {
	logMutex.Lock()
	defer logMutex.Unlock()

	oldLevel := minLogLevel
	minLogLevel = newLevel
	return oldLevel
}

// SetPrefix sets the prefix to appear on all log entries
func SetPrefix(prefix string) {
	logMutex.Lock()
	defer logMutex.Unlock()

	ilogger.SetPrefix(prefix)
}

// SetOutput sets the logging output writer instance. By
// default mlog uses stderr.
func SetOutput(w io.Writer) {
	logMutex.Lock()
	defer logMutex.Unlock()

	ilogger.SetOutput(w)
}

// Same as CloseLogFiles() but will replace it at some point.
// use defer Close() after setting up mlog.
func Close() {
	CloseLogFiles()
}

// CloseLogFiles to close the log file. Call this in a defer statement in your
// main() IF you specified a log filename in the LOG_FILENAME environment var.
// It does nothing if you used SetOutput() with your own file writer.
func CloseLogFiles() {
	if logFile != nil {
		ilogger.Print(_LOG_TRAILER)
		err := logFile.Close()
		if err != nil {
			ilogger.Printf("Error closing log file: %v", err)
		}
	}

	if catFile != nil {
		catFile.WriteString(_LOG_TRAILER)
		err := catFile.Close()
		if err != nil {
			ilogger.Printf("Error closing catheter file: %v", err)
		}
	}

	if treeFile != nil {
		treeFile.WriteString(_LOG_TRAILER)
		err := treeFile.Close()
		if err != nil {
			ilogger.Printf("Error closing call tree file: %v", err)
		}
	}
}

// Use a separate Call Tree log file for extended log functions.
// The default is the same as the regular log output.
func WithCallTree(filename string) error {
	var err error
	if treeFile, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644); err == nil {
		treeFile.WriteString(_CALLTREE_FILE_LEADER)
		greetLog(treeFile)
	}
	return err
}

// Writes to the call tree (separate file or default output).
// The data string is expected WITHOUT line feed ("\n") as
// it will be appended here.
func writeCallTree(data string) {
	data = data + "\n"
	if treeFile == nil {
		ilogger.Print(data)
	} else {
		treeFile.WriteString(data)
	}
}

// opens the log file and outputs the first message to delimit
// multiple application runs.
func openLogFile(filePath string, isMainLog, truncate bool) (*os.File, error) {
	fileFlags := os.O_CREATE | os.O_WRONLY
	if truncate {
		fileFlags |= os.O_TRUNC
	} else {
		fileFlags |= os.O_APPEND
	}

	logFileX, err := os.OpenFile(filePath, fileFlags, 0666)
	if err != nil {
		return nil, err
	}

	if isMainLog {
		ilogger.Print(string(UTF8_BOM), _LOG_LEADER)
		greetLog(ilogger.Writer())
	} else {
		logFileX.WriteString(string(UTF8_BOM) + _LOG_LEADER)
		greetLog(logFileX)
	}

	return logFileX, nil
}

// Outputs a greeter (to selected log's file descriptor) showing Local & UTC
// date/time for comparative reference, and the time zone with UTC offset.
func greetLog(w io.Writer) {
	timeRef := getTimeReference()
	timeOfs := getOffsetUTC()
	if wStr, ok := w.(io.StringWriter); ok {
		wStr.WriteString(timeRef[0] + osx.EOL)
		wStr.WriteString(timeRef[1] + osx.EOL)
		wStr.WriteString(timeOfs + osx.EOL)
		wStr.WriteString(GetMode() + osx.EOL)
	}
}

// Returns a time-reference string that includes both Local & UTC times
func getTimeReference() []string {
	timeNow := time.Now()
	timeRefLocal := fmt.Sprintf("%c  Local time: %s",
		UC_TIMER_RUNNING,
		timeNow.Format(time.DateTime))
	timeRefUTC := fmt.Sprintf("%c  Zulu time : %s UTC",
		UC_TIMER_RUNNING,
		timeNow.UTC().Format(time.DateTime))
	return []string{timeRefLocal, timeRefUTC}
}

// Gets the Timezone name followed by the UTC offset in HH:MM
func getOffsetUTC() string {
	zoneName, offsetSecs := time.Now().Zone()
	hours := offsetSecs / 3600
	mins := (offsetSecs % 3600) / 60

	if mins < 0 {
		mins = -mins
	}
	return fmt.Sprintf("%c  Timezone  : %s %03d:%02d", UC_TIMER_RUNNING, zoneName, hours, mins)
}

// parse a string to convert it to a LogLevel value
func parseLevel(s string) LogLevel {
	var lvl LogLevel
	s = strings.Trim(s, " \t")

	switch {
	case strings.EqualFold(s, "trace"):
		lvl = LevelTrace

	case strings.EqualFold(s, "debug"):
		lvl = LevelDebug

	case strings.EqualFold(s, "info"):
		lvl = LevelInfo

	case strings.EqualFold(s, "warning"):
		fallthrough
	case strings.EqualFold(s, "warn"):
		lvl = LevelWarning

	case strings.EqualFold(s, "error"):
		lvl = LevelError

	case strings.EqualFold(s, "fatal"):
		lvl = LevelFatal

	default:
		lvl = LevelFatal
	}

	return lvl
}

/* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *			N o n - P r i v i l e g e d   L e v e l s
 *- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -*/

// (Log Level: WARN) Warning level with variadic parameters
func Warn(v ...any) {
	if minLogLevel <= LevelWarning {
		v1 := append([]any{tagWARN}, v...)
		ilogger.Print(v1...)
	}
}

// (Log Level: WARN) Warning level with format string
func Warnf(format string, v ...any) {
	if minLogLevel <= LevelWarning {
		ilogger.Printf(tagWARN+format, v...)
	}
}

// (Log Level: WARN) Warning level with message and variadic MLog tags.
func WarnT(message string, v ...mtag.ILogKeyValuePair) {
	if minLogLevel <= LevelWarning {
		var sb strings.Builder
		sb.WriteString(tagWARN)
		sb.WriteString(message)
		for _, t := range v {
			sb.WriteString(" ")
			sb.WriteString(t.String())
		}
		ilogger.Print(sb.String())
	}
}

// (Log Level: ERROR) Error level with variadic parameters
func Error(v ...any) {
	if minLogLevel <= LevelError {
		v1 := append([]any{tagERROR}, v...)
		ilogger.Print(v1...)
	}
}

// (Log Level: ERROR) Error level with format string
func Errorf(format string, v ...any) {
	if minLogLevel <= LevelError {
		ilogger.Printf(tagERROR+format, v...)
	}
}

// (Log Level: ERROR) Error level with message and variadic MLog tags.
func ErrorT(message string, v ...mtag.ILogKeyValuePair) {
	if minLogLevel <= LevelError {
		var sb strings.Builder
		sb.WriteString(tagERROR)
		sb.WriteString(message)
		for _, t := range v {
			sb.WriteString(" ")
			sb.WriteString(t.String())
		}
		ilogger.Print(sb.String())
	}
}

// (Log Level: ERROR) Error level limited to the error itself
func ErrorE(err error) {
	if minLogLevel <= LevelError {
		ilogger.Println(tagERROR, err.Error())
	}
}

// (Log Level: FATAL) Fatal level with variadic parameters and exitCode
// for terminating the application.
func Fatal(exitCode int, v ...any) {
	if minLogLevel <= LevelFatal {
		v1 := append([]any{tagFATAL}, v...)
		ilogger.Print(v1...)
	}

	os.Exit(exitCode)
}

// (Log Level: FATAL) Fatal level with format string and exitCode for terminating
// the application.
func Fatalf(exitCode int, format string, v ...any) {
	if minLogLevel <= LevelFatal {
		ilogger.Printf(tagFATAL+format, v...)
	}

	os.Exit(exitCode)
}

// (Log Level: FATAL) Fatal level with message and variadic MLog tags.
// it terminates execution with exitCode.
func FatalT(exitCode int, message string, v ...mtag.ILogKeyValuePair) {
	if minLogLevel <= LevelFatal {
		var sb strings.Builder
		sb.WriteString(tagFATAL)
		sb.WriteString(message)
		for _, t := range v {
			sb.WriteString(" ")
			sb.WriteString(t.String())
		}
		ilogger.Print(sb.String())
	}

	os.Exit(exitCode)
}

/* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *			T o   B e   I m p l e m e n t e d
 *- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -*/

// will implement ILogFiltered
func SaveFilters()                                         {}
func LoadFilters() error                                   { return fmt.Errorf("Not Implemented (Pending migration)") }
func IsFiltered(packageName string) bool                   { return false }
func IsFilteredObject(packageName, objectName string) bool { return false }
