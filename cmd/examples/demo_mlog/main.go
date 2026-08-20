/*
 * Build instructions:
 *	(a) For debug version: -tags mlog,debug
 *	(b) For release version: -tags mlog
 */
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/lordofscripts/goapp"
	"github.com/lordofscripts/goapp/app/mlog"
	. "github.com/lordofscripts/goapp/app/mtag"
)

/* ----------------------------------------------------------------
 *				M o d u l e   I n i t i a l i z a t i o n
 *-----------------------------------------------------------------*/

// It is advisable to setup Logging in the initializer
func init() {
	mlog.SetOutput(os.Stderr)
	mlog.SetLevel(mlog.LevelTrace)
	mlog.SetPrefix("DEMO")
	mlog.WithCallTree("/tmp/call_tree_demo.log")
	mlog.Greet()
}

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type DemoObject struct {
	text string
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (do *DemoObject) SampleMethod(str string, nr int) {
	mlog.Enter()
	defer mlog.Leave()

	if str != do.text {

	}
}

func (do *DemoObject) SetText(newStr string) {
	mlog.Enter()
	defer mlog.Leave()

	mlog.OnUpdate()
	mlog.Debugf("Setting text to '%s'", newStr)

	if newStr != do.text {
		do.changedCB(do.text, newStr)
	}
}

func (do *DemoObject) changedCB(oldStr, newStr string) {
	mlog.EventEnter()
	defer mlog.EventLeave()

	mlog.InfoT("text changed", String("Old", oldStr), String("New", newStr))
	mlog.AttentionAlways("are you reading?", io.EOF)
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

func oneFunction() int {
	mlog.Enter()
	defer mlog.Leave()

	x := 5
	y := 20
	result := x * y
	mlog.Result("Returning with", result)
	return result
}

func Help() {
	flag.PrintDefaults()
}

/* ----------------------------------------------------------------
 *                  M A I N    |    D E M O
 *-----------------------------------------------------------------*/

/*
 * Optional environment variables:
 *	LOG_LEVEL_CX : trace|debug|info|warn|error|fatal
 *	LOG_FILE_CX	 : output log filename (else stderr)
 * Those can be overriden by what is initialized in main.
 */
func main() {
	const CHR_TRIDENT rune = rune(0x1f531) // 🔱
	// - will close all opened log file handles
	defer mlog.Close()

	goapp.ModuleVersion.Copyright("lordofscripts", CHR_TRIDENT)

	// A. Command-line flags
	var noLog, noTimestamp, help bool
	var level mlog.LogLevel = mlog.LevelTrace
	var levelStr string
	flag.BoolVar(&help, "help", false, "Help me")
	flag.BoolVar(&noLog, "nolog", false, "disable logging")
	flag.BoolVar(&noTimestamp, "no-timestamp", false, "log entries without timestamp")
	flag.StringVar(&levelStr, "level", "", "Set logging level")
	flag.Parse()

	// B. command-line override(s)
	var err error
	if noLog {
		// - to disable logging
		mlog.SetOutput(io.Discard) // by default it does os.Stderr
	}
	if len(levelStr) != 0 {
		// if the level name is invalid it defaults to the current value
		if level, err = level.Parse(levelStr); err == nil {
			mlog.SetLevel(level)
		}
	}
	if noTimestamp {
		// Remove Ldate|Ltime
		newFlags := mlog.Flags() &^ log.Ldate
		newFlags = newFlags &^ log.Ltime
		mlog.SetFlags(newFlags)
	}

	fmt.Printf("MLog Level: %s (%d)\n", level, level)

	// C. Stuff
	var byVal byte = 5
	var iVal int = 800
	var boVal bool = false
	var stVal string = "Thiß Über Alles"
	var rVal rune = 'ß'
	err = fmt.Errorf("oops, bad thing happened.")

	// D. Logging for Event-driven applications (some to try)
	myObj := &DemoObject{text: "initial"}
	myObj.SampleMethod(stVal, iVal)

	// E. Level aware logging (tagged and normal)
	mlog.Trace("test trace", rVal, stVal)
	mlog.Tracef("formatted %d trace %t", iVal, boVal)
	mlog.TraceT("tagged trace", Byte("byte", byVal), Bool("True?", boVal))

	oneFunction()

	mlog.Debug("test debug", rVal, stVal)
	mlog.Debugf("formatted %d debug %t", iVal, boVal)
	mlog.DebugT("tagged debug", Int("integer", iVal), String("text", stVal))

	mlog.Info("test info", rVal, stVal)
	mlog.Infof("formatted %d info %t", iVal, boVal)
	mlog.InfoT("tagged info", Rune("letter", rVal), Bool("True?", boVal))

	mlog.Warn("test warning", rVal, stVal)
	mlog.Warnf("formatted %d warning %t", iVal, boVal)
	mlog.WarnT("tagged warning", String("text", stVal), Err(err))

	mlog.Error("test error", rVal, stVal)
	mlog.Errorf("formatted %d error %t", iVal, boVal)
	mlog.ErrorT("tagged error", Int("integer", iVal), Rune("letter", rVal))
	mlog.ErrorE(err)

	// The first of these exits with that application exit code
	const EXIT_CODE_1 int = 1
	const EXIT_CODE_2 int = 2
	const EXIT_CODE_3 int = 3
	mlog.Fatal(EXIT_CODE_1, "test fatal", rVal, stVal)
	mlog.Fatalf(EXIT_CODE_2, "formatted %d fatal %t", iVal, boVal)
	mlog.FatalT(EXIT_CODE_3, "tagged fatal", Byte("byte", byVal), Bool("True?", boVal))
}
