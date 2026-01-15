/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Lord of Scripts
 *							   goApp
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Convenience Logging functions present when "logx" tag is used.
 *-----------------------------------------------------------------*/
package logx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
)

/* ----------------------------------------------------------------
 *					P u b l i c		G l o b a l s
 *-----------------------------------------------------------------*/
const (
	UC_FOOTSTEPS   rune = rune(0x1f463) // 👣
	UC_EXCLAMATION rune = rune(0x2757)  // ❗
	UC_CROSSMARK   rune = rune(0x274c)  // ❌
	UC_CROSS       rune = rune(0x2716)  // ✖
	UC_CHECK       rune = rune(0x1f5f8) // 🗸
	UC_ARROWS3     rune = rune(0x21f6)  // ⇶
	UC_OBSERVER    rune = rune(0x23ff)  // ⏿
	UC_EYES        rune = rune(0x1f440) // 👀
	UC_EYE         rune = rune(0x1f441) // 👁
)

var (
	SingLogGate *LogGate // Singleton for App CLI parsing and Filtering
)

/* ----------------------------------------------------------------
 *					P r i v a t e	G l o b a l s
 *-----------------------------------------------------------------*/

var (
	initializedLogGate atomic.Uint32 // singleton management
	instanceLogGate    *LogGate      // singleton management
	onceLogGate        sync.Mutex    // singleton management
)

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/

type ICallTree interface {
}

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

type LogGate struct {
	AppName      string               `json:"appname"`
	CurrentLevel string               `json:"level"` // allowed log level: NONE,FATAL,ERROR,WARN,INFO,DEBUG
	Filters      map[string]LogFilter `json:"filters"`
	fd           *os.File
	fdCallTree   *os.File
	configSub    string
}

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

type LogFilter struct {
	LogLevel     string `json:"log_level"`
	Specifically string `json:"specifically,omitempty"`
}

/* ----------------------------------------------------------------
 *				I n i t i a l i z e r
 *-----------------------------------------------------------------*/

/*
func init() {
	SingLogGate = GetLogGateInstance()
}
*/

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

/**
 * Manages the Logging configuration for the application lifecycle
 * Allowing to easily set it up via a CLI parameter to disable
 * logging, or enable to a terminal or a named file.
 *
 * @returns LogGate{} singleton
 */
func GetLogGateInstance(names ...string) *LogGate {
	if initializedLogGate.Load() == 1 {
		return instanceLogGate
	}

	onceLogGate.Lock()
	defer onceLogGate.Unlock()

	if initializedLogGate.Load() == 0 && instanceLogGate == nil {
		var applicationName = path.Base(os.Args[0])
		var groupName = ""
		if names != nil {
			if names[0] != "" {
				applicationName = names[0]
			}
			if names[1] != "" {
				groupName = names[1]
			}
		}
		welcome := fmt.Sprintf(">>>>>> Welcome to %s <<<<<<", applicationName)
		instanceLogGate = newLogGate()
		instanceLogGate.AppName = applicationName
		instanceLogGate.configSub = groupName
		initializedLogGate.Store(1)
		instanceLogGate.LoadFilters()

		log.Print(welcome)
	}

	return instanceLogGate
}

func newLogGate() *LogGate {
	return &LogGate{
		CurrentLevel: "", // @todo To be implemented
		Filters:      make(map[string]LogFilter, 0),
		fd:           nil,
		fdCallTree:   nil,
		configSub:    ""}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

func (l *LogGate) Help() string {
	return "Log modes: (empty)|none|term|FOLDER_NAME"
}

func (l *LogGate) Setup(mode string) {
	switch strings.ToLower(mode) {
	case "":
		fallthrough
	case "none":
		l.DisableLogging()

	case "term":
		l.EnableLogging()

	default:
		l.EnableLoggingToDir(mode)
	}
}

func (l *LogGate) SetAppName(name string) *LogGate {
	l.AppName = name
	return l
}

func (l *LogGate) WithConfigSubdirectory(subdir string) *LogGate {
	l.configSub = subdir
	return l
}

func (l *LogGate) WithCallTree(filename string) {
	var err error
	if l.fdCallTree, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644); err == nil {
		l.fdCallTree.WriteString(">>>>> " + l.AppName + " LogX Call Tree <<<<<<\n")
	}
}

// EnableLogging enables logging for tests
func (l *LogGate) EnableLogging() {
	//log.SetOutput(os.Stdout) // Enable logging output
	log.SetOutput(os.Stderr) // Enable logging but still to STDERR
}

/**
 * Enable logging to file.
 * @returns (*os.File) so that the caller can defer Close()
 */
func (l *LogGate) EnableLoggingToDir(dirname string) error {
	filename := path.Join(dirname, os.Args[0]+".log")
	return l.EnableLoggingToFile(filename)
}

func (l *LogGate) EnableLoggingToFile(filename string) error {
	var err error
	if l.fd, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644); err == nil {
		syscall.Dup2(int(l.fd.Fd()), 1) // stdout
		syscall.Dup2(int(l.fd.Fd()), 2) // stderr
	}

	return err
}

// DisableLogging disables logging in the application
func (l *LogGate) DisableLogging() {
	log.SetOutput(io.Discard) // Disable logging output
}

/**
 * This call is only necessary if the application is logging to a file.
 * In that case it is better to put a call to this on a defer statement
 * on the main function.
 */
func (l *LogGate) Close() {
	if l.fdCallTree != nil {
		l.fdCallTree.WriteString("#### Terminated ####\n")
		if err := l.fdCallTree.Close(); err != nil {
			fmt.Printf("error closing call tree log: %v\n", err)
		}
	}

	if l.fd != nil {
		if err := l.fd.Close(); err != nil {
			fmt.Printf("error closing log: %v\n", err)
		}
	}
}

func (l *LogGate) WriteCallTree(data string) {
	l.fdCallTree.WriteString(data + "\n")
}

/**
 * Saves a sample Log Filter configuration file in the user's configuration
 * directory. The filename is APPNAME.logfilter.
 */
func (l *LogGate) SaveFilters() {
	l.Filters["main"] = LogFilter{LogLevel: "debug", Specifically: ""}
	l.Filters["lordofscripts/demo"] = LogFilter{LogLevel: "info", Specifically: "StructA,StructB"}

	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(l); err != nil {
		log.Print("marshall filter failure: ", err)
	} else {
		if fd, err := os.Create(l.getLogFilterFile()); err == nil {
			defer fd.Close()
			fd.WriteString(buffer.String())
		} else {
			log.Print("couldn't save filter file")
		}
	}
}

/**
 * Attempts to load APPNAME.logfilters which contains a list of application/module
 * packages white-listed for log output. It may contain a fine-grained configuration
 * of which objects (structs) of that package are allowed. If empty all are allowed,
 * else whatever struct name that is not listed becomes black-listed for logging.
 */
func (l *LogGate) LoadFilters() error {
	var all LogGate
	jsonDataBytes, err := os.ReadFile(l.getLogFilterFile())
	if err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}

	err = json.Unmarshal(jsonDataBytes, &all)
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return err
	}

	l.Filters = all.Filters
	return nil
}

/**
 * Checks if the packageName (GO package name format) is
 * black-listed. Black-listed packages do not produce log
 * output under LOGX.
 */
func (l *LogGate) IsFiltered(packageName string) bool {
	if filter, exists := l.Filters[packageName]; !exists ||
		exists && filter.LogLevel == "" {
		return true
	}

	return false
}

/**
 * Checks if the packageName (GO package name format) is
 * black-listed. Black-listed packages do not produce log
 * output under LOGX.
 */
func (l *LogGate) IsFilteredObject(packageName, objectName string) bool {
	var blackListed bool = true
	filter, exists := l.Filters[packageName]
	if exists {
		if filter.Specifically == "*" {
			blackListed = false
		} else if filter.LogLevel != "" && strings.Contains(filter.Specifically, objectName) {
			blackListed = false
		}
	}
	/*
		if filter, exists := l.Filters[packageName]; !exists ||
			exists && filter.LogLevel == "" ||
			exists && filter.LogLevel != "" && !strings.Contains(filter.Specifically, objectName) {
			return true
		}

		return false
	*/
	return blackListed
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

/**
 * Derive the application name from the executable. Strip ".exe"
 * if present.
 */
func (l *LogGate) getAppName() string {
	if l.AppName != "" {
		return l.AppName
	}
	appName := os.Args[0]
	appName = strings.TrimSuffix(appName, ".exe")
	return appName
}

/**
 * Full path to JSON file with Log Filter
 */
func (l *LogGate) getLogFilterFile() string {
	var cfgPath string
	if usrConfig, err := os.UserConfigDir(); err != nil {
		panic(err)
	} else {
		basename := l.getAppName() + ".logfilter"
		cfgPath = filepath.Join(usrConfig, l.configSub, basename)
	}

	return cfgPath
}

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

/**
 * Gets the current nesting level with main.main() as reference and
 * a string with the location (package/struct/function).
 */
func GetNestingLevel(frCnt ...int) (int, string) {
	spec := "%p%O%F"
	// Capture the call stack
	popCnt := 2
	if frCnt != nil {
		popCnt = frCnt[0]
	}
	pc := make([]uintptr, 10)        // Adjust size as needed
	n := runtime.Callers(popCnt, pc) // Skip the first two frames (getNestingLevel and the calling function)

	// Count how many frames until we reach main
	nestingLevel := 0
	var pk, s, f string
	var frame runtime.Frame
	for _, p := range pc[:n] {
		frame, _ = runtime.CallersFrames([]uintptr{p}).Next()
		if frame.Function == "main.main" {
			break
		}

		if nestingLevel == 0 {
			pk, s, f = getNames(frame.Function)
		}
		nestingLevel++
	}

	if nestingLevel == 0 && pk == "" && f == "" {
		pk = "main"
		s = ""
		f = "main"
	}

	const SEP_FUNC = "⯈"
	const SEP_METH_PTR = "🡪"
	const SEP_PKG = "▪" // ⮩🡂🡆⤷⮆🢂⭄🠞⏵➠⯈◾
	pretty := spec

	// transform package
	//	· %P	fully-qualified package name
	//	· %p	package base-name (last part)
	if strings.Contains(spec, "%p") {
		if idx := strings.LastIndex(pk, "/"); idx != -1 {
			pk = pk[idx+1:]
		}
	}
	pk = pk + SEP_PKG
	pretty = strings.Replace(pretty, "%P", pk, 1)
	pretty = strings.Replace(pretty, "%p", pk, 1)

	// transform struct/object (if any)
	isPointer := strings.HasPrefix(s, ".(*")
	if s != "" {
		if isPointer {
			s = s[3 : len(s)-1]
			s = s + SEP_METH_PTR
		} else {
			s = s[1:len(s)-1] + SEP_FUNC
		}
	}
	pretty = strings.Replace(pretty, "%O", s, 1)
	pretty = strings.Replace(pretty, "%o", s, 1)

	// transform function
	if strings.Contains(pretty, "%F") {
		f = f + "()"
	}
	pretty = strings.Replace(pretty, "%F", f, 1)
	pretty = strings.Replace(pretty, "%f", f, 1)

	//return nestingLevel, fmt.Sprintf("%03d %s", nestingLevel, pretty)
	return nestingLevel, pretty
}

func getNames(fq string) (string, string, string) {
	dotCnt := strings.Count(fq, ".")
	index := strings.LastIndex(fq, ".")
	var namePkg, nameStruct, nameFunc string
	switch dotCnt {
	case 1:
		// A function
		nameFunc = fq[index+1:]
		nameStruct = ""
		namePkg = fq[:index]
	case 2:
		// A method
		nameFunc = fq[index+1:]
		otherPart := fq[:index]
		index = strings.LastIndex(otherPart, ".")
		nameStruct = otherPart[index:]
		namePkg = otherPart[:index]
	default:
		// github.com/lordofscripts/caesardisk/cmd/gui-app/gui.(*CipherModeGadget).Define.func1
		re := regexp.MustCompile(`^[-\w]+\.[A-Za-z]+/`)
		cleaned := re.ReplaceAllString(fq, "")
		parts := strings.Split(cleaned, ".")
		namePkg = (parts[0])[strings.LastIndex(parts[0], "/"):]
		nameStruct = parts[1]
		if len(parts) == 3 {
			nameFunc = parts[2]
		} else if len(parts) == 4 {
			nameFunc = parts[2] + "." + parts[3]
		} else {
			println(fq)
			panic("WTF")
		}
	}

	//fmt.Printf("\tP:%s S:%s F:%s\n", namePkg, nameStruct, nameFunc)
	return namePkg, nameStruct, nameFunc
}
