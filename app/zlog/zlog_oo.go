//go:build zlog

/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                        goApp::zlog
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * ZLog is the Flexible Logger, the next-generation GoApp logger that
 * has the features of LogX & MLog combined and will eventually
 * deprecate those two packages.
 *-----------------------------------------------------------------*/
package zlog

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/lordofscripts/goapp/osx"
)

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/

const (
	// Logging level log-line prefixes
	tagCATHE string = "[CAT] "
	tagTRACE string = "[TRC] "
	tagDEBUG string = "[DBG] "
	tagINFO  string = "[INF] "
	tagWARN  string = "[WRN] "
	tagERROR string = "[ERR] "
	tagFATAL string = "[DIE] "
)

// Names of environment variables used by this package
const (
	// environment variable that overrides the default (Error) Log Level for CaesarX
	LOG_LEVEL_ENV string = "LOG_LEVEL_CX"
	// environment variable that indicates the log output filename for CaesarX (default stderr)
	LOG_FILE_ENV string = "LOG_FILE_CX"
)

// Unicode code points & emojis
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

const (
	defaultPrefix string   = ""
	defaultFlags  int      = log.Ldate | log.Ltime | log.Lmsgprefix
	defaultLevel  LogLevel = LevelError

	_LOG_LEADER           string = "[BEG]\t> > > >   T h e   B e g i n n i n g   < < < <\n"
	_LOG_TRAILER          string = "[END]\t> > > >   T h e   E n d   < < < <\n"
	_CALLTREE_FILE_LEADER string = ">>>>> " + " ZLog Call Tree <<<<<<\n"
)

// UTF8 BOM (Byte Order Mark)
var UTF8_BOM []byte = []byte{0xEF, 0xBB, 0xBF}

/* ----------------------------------------------------------------
 *                     I N T E R F A C E S
 *-----------------------------------------------------------------*/

var _ ILogConfigurable = (*FlexLogger)(nil)
var _ ILogEventDriven = (*FlexLogger)(nil) // NO-OP if !develop
var _ ILogPrioritized = (*FlexLogger)(nil)
var _ ILogTaggable = (*FlexLogger)(nil)
var _ ILogFilterable = (*FlexLogger)(nil)
var _ ILogCatheterable = (*FlexLogger)(nil) // NO-OP if !develop

/* ----------------------------------------------------------------
 *                   P U B L I C    T Y P E S
 *-----------------------------------------------------------------*/

// MLog logger object implementing various logging interfaces for
// maximum flexibility. The following interfaces are available:
type FlexLogger struct {
	logMutex     sync.Mutex
	minLogLevel  LogLevel
	ilogger      *log.Logger
	logFile      ILoggable // the standard log output file
	catFile      ILoggable // the "catheter output file"
	treeFile     ILoggable // the "call tree" file for extended UI functions
	customWriter *customLogWriter
	filterSvc    ILogFilterService
}

/* ----------------------------------------------------------------
 *                   P R I V A T E    T Y P E S
 *-----------------------------------------------------------------*/

type customLogWriter struct {
	writer io.Writer
	format string
}

/* ----------------------------------------------------------------
 *                     I N I T I A L I Z E R
 *-----------------------------------------------------------------*/
func init() {}

/* ----------------------------------------------------------------
 *                    C O N S T R U C T O R S
 *-----------------------------------------------------------------*/

// (Ctor) new instance of Maximum Flexibility Logger. It defaults to
// Error level (Error+Fatal) unless overriden by an environment
// variable. Defaults to StdErr output unless overridden with an
// environment variable specifying the output log file.
// See Also: WithWriter()
func NewZLogger() *FlexLogger {
	const CUSTOM_TIME_FORMAT = "2006-01-02 15:04:05"

	myLogObj := &FlexLogger{
		minLogLevel:  LevelError,
		ilogger:      nil,
		logFile:      nil,
		catFile:      nil,
		treeFile:     nil,
		customWriter: newCustomLogWriter(os.Stderr, CUSTOM_TIME_FORMAT),
		filterSvc:    NewNullLogFilterService(), // no blacklisting
	}

	myLogObj.getLevelFromEnvironment(). // overrides minLogLevel
						getDefaultLogger(). // sets & configures ilogger
						getFileFromEnvironment()

	return myLogObj
}

// (Ctor) new instance of Maximum Flexibility Logger. It defaults to
// Error level (Error+Fatal) unless overriden by an environment
// variable. It attempts to use filename as the output file, if that
// fails it reverts to StdErr.
// See Also: WithWriter()
func NewZLoggerWithFile(filename string) *FlexLogger {
	myLogObj := &FlexLogger{
		minLogLevel: LevelError,
		ilogger:     nil,
		logFile:     nil,
		catFile:     nil,
		treeFile:    nil,
		filterSvc:   NewNullLogFilterService(), // no blacklisting
	}

	myLogObj.getLevelFromEnvironment(). // overrides minLogLevel
						getDefaultLogger(). // sets & configures ilogger
						setFileOutput(filename)

	return myLogObj
}

func newCustomLogWriter(w io.Writer, timeStampFormat string) *customLogWriter {
	return &customLogWriter{writer: w, format: timeStampFormat}
}

/* ----------------------------------------------------------------
 *                        M E T H O D S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                 P U B L I C    M E T H O D S
 *-----------------------------------------------------------------*/

// It tries to set the *standard* log output to the given w writer.
// It closes any previously open descriptor. This method is useful
// if you already have an open file, or to use a custom file writer,
// or even a custom writer that outputs logging to a Fyne text widget.
func (blg *FlexLogger) WithWriter(w ILoggable) *FlexLogger {
	if w != nil {
		// close currently active (if any)
		blg.closeFile(blg.logFile, "STDLOG")
		// redirect to the new writer
		blg.logFile = w
		blg.ilogger.SetOutput(w)
	}

	return blg
}

// It tries to set the *event-driven-app* log output to the given w writer.
// It closes any previously open descriptor. This method is useful
// if you already have an open file, or to use a custom file writer,
// or even a custom writer that outputs logging to a Fyne text widget.
// In a Fyne app you can therefore have ILogPrioritized & ILogTaggable
// to one widget and ILogEventDriven to another widget.
func (blg *FlexLogger) WithUiWriter(w ILoggable) *FlexLogger {
	if w != nil {
		// close currently active (if any)
		blg.closeFile(blg.treeFile, "UILOG")
		// redirect to the new writer
		blg.treeFile = w
		blg.ilogger.SetOutput(w)
	}

	return blg
}

// The current build of the MLog library (development|production)
func (blg *FlexLogger) GetMode() string {
	return GetMode()
}

// Outputs a greeter to the log outputs. The greeter has the local & UTC
// date & times and the Timezone + UTC offset information. This is usually
// called as a convenience after ZLog has been configured and the application
// is ready to start.
func (blg *FlexLogger) Greet() {
	if blg.ilogger != nil {
		Greet(blg.ilogger.Writer())
	}
	if blg.logFile != nil {
		Greet(blg.logFile)
	}
	if blg.catFile != nil {
		Greet(blg.catFile)
	}
	if blg.treeFile != nil {
		Greet(blg.treeFile)
	}
}

// Log an error irrespective of the log level
func (blg *FlexLogger) AttentionErr(short string, err error) {
	blg.ilogger.Print(short, err)
}

//	****	implements ILogConfigurable 	****

// implements io.Writer
func (blg *FlexLogger) Write(p []byte) (n int, err error) {
	timestamp := time.Now().Format(blg.customWriter.format)
	formattedMessage := fmt.Sprintf("%s %s", timestamp, string(p))
	return blg.customWriter.writer.Write([]byte(formattedMessage))
}

// implements io.Closer to close any active file descriptors.
// use defer Close()
func (blg *FlexLogger) Close() error {
	err1 := blg.closeFile(blg.logFile, "STDLOG")
	err2 := blg.closeFile(blg.catFile, "CATLOG")
	err3 := blg.closeFile(blg.treeFile, "UILOG")

	return errors.Join(err1, err2, err3)
}

// SetPrefix sets the prefix to appear on all log entries
func (blg *FlexLogger) SetPrefix(prefix string) {
	blg.logMutex.Lock()
	defer blg.logMutex.Unlock()

	blg.ilogger.SetPrefix(prefix)
}

// SetOutput sets the logging output writer instance. By
// default mlog uses stderr.
func (blg *FlexLogger) SetOutput(w io.Writer) {
	blg.logMutex.Lock()
	defer blg.logMutex.Unlock()

	blg.ilogger.SetOutput(w)
}

// the same flags used in the standard 'log': Ldate, Ltime, Lmicroseconds,
// Llongfile, Lshortfile, LUTC, Lmsgprefix & LstdFlags (Ldate|Ltime). In
// its default configuration we use log.Ldate|log.Ltime|log.Lmsgprefix
func (blg *FlexLogger) SetFlags(log_flags int) {
	blg.ilogger.SetFlags(log_flags)
}

// The log flags that are currently used
func (blg *FlexLogger) Flags() int {
	return blg.ilogger.Flags()
}

//	****	implements ILogPrioritized 	****

// SetLevel sets the current logging level. Unlike log and slog
// the mlog package supports logging levels.
func (blg *FlexLogger) SetLevel(newLevel LogLevel) LogLevel {
	blg.logMutex.Lock()
	defer blg.logMutex.Unlock()

	oldLevel := blg.minLogLevel
	blg.minLogLevel = newLevel
	return oldLevel
}

//	****	implements ILogFilterable	****

// Enable fine-grained package & object log filtering by providing
// an instance such as NewLogFilterService() for a default implementation.
// Or implement your own service. For no filtering pass nil to this
// method, that is the default setup. IF not nil, the filter's LoadFilter()
// is automatically called.
// NOTE: This filtering ONLY applies to the ILogEventDriven API of
// this logger.
func (blg *FlexLogger) UseFilter(svc ILogFilterService) error {
	blg.logMutex.Lock()
	defer blg.logMutex.Unlock()

	if svc != nil {
		blg.filterSvc = svc
		return svc.LoadFilters()
	} else {
		blg.filterSvc = NewNullLogFilterService()
	}
	return nil
}

/* ----------------------------------------------------------------
 *                 P R I V A T E    M E T H O D S
 *-----------------------------------------------------------------*/

// Sets the Logging Level by fetching it from an environment variable.
// If the value is unset or invalid, it uses the default level (Error level).
func (blg *FlexLogger) getLevelFromEnvironment() *FlexLogger {
	levelStr := os.Getenv(LOG_LEVEL_ENV)
	if levelStr != "" {
		blg.minLogLevel = parseLevel(levelStr)
	} else {
		blg.minLogLevel = defaultLevel
	}

	return blg
}

// Configures the default logger to StdErr using an empty prefix and
// the Date/Time/Prefix flags using an ISO timestamp in 24 hour format.
func (blg *FlexLogger) getDefaultLogger() *FlexLogger {
	myLogger := log.New(os.Stderr, defaultPrefix, defaultFlags)
	myLogger.SetOutput(blg.customWriter.writer)
	blg.ilogger = myLogger

	return blg
}

// If the environment variable specifies an output file for the standard log,
// and it is able to open it, then it redirects log output to that file.
// If that is not set or it fails, nothing happens.
// NOTE: ensure getDefaultLogger has been called earlier.
func (blg *FlexLogger) getFileFromEnvironment() *FlexLogger {
	envLogFilename := os.Getenv(LOG_FILE_ENV)

	if len(envLogFilename) != 0 {
		return blg.setFileOutput(envLogFilename)
	}

	return blg
}

// attempts to set the log output to the given filename, if it is unable
// to comply the previous configuration remains.
func (blg *FlexLogger) setFileOutput(filename string) *FlexLogger {
	if fd, err := blg.openLogFile(filename, true, false); err == nil {
		blg.logFile = fd
		blg.ilogger.SetOutput(fd)
	}
	return blg
}

// opens the log file and outputs the first message to delimit
// multiple application runs.
func (blg *FlexLogger) openLogFile(filePath string, isMainLog, truncate bool) (*os.File, error) {
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
		blg.ilogger.Print(string(UTF8_BOM), _LOG_LEADER)
		Greet(blg.ilogger.Writer())
	} else {
		logFileX.WriteString(string(UTF8_BOM) + _LOG_LEADER)
		Greet(logFileX)
	}

	return logFileX, nil
}

// closes the file by the file descriptor if it is not nil. If
// successful it sets the descriptor to nil.
func (blg *FlexLogger) closeFile(fd ILoggable, alias string) error {
	if fd != nil {
		fd.WriteString(_LOG_TRAILER)
		err := fd.Close()
		if err != nil {
			blg.ilogger.Printf("couldn't close '%s': %v", alias, err)
			return err
		}
		fd = nil
	}
	return nil
}

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

// Outputs a greeter (to selected log's file descriptor) showing Local & UTC
// date/time for comparative reference, and the time zone with UTC offset.
func Greet(w io.Writer) {
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
	return fmt.Sprintf("%c  Timezone  : %s UTC%+03d:%02d", UC_TIMER_RUNNING, zoneName, hours, mins)
}
