//go:build zlog

/*
 * Build instructions (With Makefile or TaskFile):
 *	(a) For debug version: -tags zlog,develop
 *	(b) For release version: -tags zlog
 * Usage:
 *	demo_zlog -help
 *	demo_zlog
 *	demo_zlog [-nolog][--no-timestamp][-level LEVEL]
 * Where:
 *	LEVEL ::= trace|debug|info|warn|warning|error|fatal
 *
 * Application-specific log filter ($HOME/.config/coralys/demo_zlog_logfilter.yaml):
 * 		filters:
 *			main:
 *				loglevel: Warn
 *				specifically: "DemoObject"
 */
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/lordofscripts/goapp"
	"github.com/lordofscripts/goapp/app/mtag"
	. "github.com/lordofscripts/goapp/app/mtag"
	"github.com/lordofscripts/goapp/app/zlog"
)

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/

const EXIT_CODE_1 int = 1
const EXIT_CODE_2 int = 2
const EXIT_CODE_3 int = 3
const USER_CONFIG_SUBDIR = "coralys" // ~/home/you/.config/coralys/zlog_demo

// declare once for use throughout your package
var (
	ß     zlog.IFlexLogger
	level zlog.LogLevel = zlog.LevelTrace // overrides default and ENV
)

var (
	byVal byte   = 5
	iVal  int    = 800
	boVal bool   = false
	stVal string = "Thiß Über Alles"
	rVal  rune   = 'ß'
	err   error  = fmt.Errorf("oops, bad thing happened.")
)

/* ----------------------------------------------------------------
 *				M o d u l e   I n i t i a l i z a t i o n
 *-----------------------------------------------------------------*/

// Early package-wide initialization. It is advisable to setup
// application logging in the main initializer.
func init() {
	// A. Instantiate the ZLogger
	//ß = zlog.NewMLoggerWithFile("/tmp/logdemo.log")
	ß = zlog.NewZLogger() // logs to StdErr console as usual

	// A.1 Extra setup for tweaking defaults
	ß.SetLevel(zlog.LevelDebug)
	ß.SetPrefix("ZLOG ")
	// (optional) for event-driven application logging to a file
	ß.WithCallTree("/tmp/call_tree_demo.log") // to a file
	//ß.WithUiWriter(&WidgetLogger{}) // to a Fyne widget

	// A.2 Decide where you want the output (default stderr)
	//ß.SetOutput(os.Stdout) // or any open file descriptor

	// A.3 Greet and meet (optional) This frames your log output
	// with Local time, UTC time, Timezone, UTC offset, Log build mode.
	ß.Greet()
}

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

// To demonstrate *some* ILogEventDriven capabilities
type DemoObject struct {
	text      string
	Validator func(string) bool
}

/* ----------------------------------------------------------------
 *                    C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func newDemoObject(initial string) *DemoObject {
	ß.Ctor()

	return &DemoObject{text: initial}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (do *DemoObject) SampleMethod(str string, nr int) {
	ß.Enter()
	defer ß.Leave()

	if str != do.text {
		ß.Debug("Sample method", str, do.text)
	}
}

func (do *DemoObject) SetText(newStr string) {
	ß.Enter()
	defer ß.Leave()

	if do.Validator != nil {
		valid := do.Validator(newStr)
		if valid {
			fmt.Printf("We like '%s'\n", newStr)
		}
	}

	ß.OnUpdate()
	ß.Debugf("Setting text to '%s'", newStr)

	if newStr != do.text {
		do.changedCB(do.text, newStr)
	}
}

func (do *DemoObject) changedCB(oldStr, newStr string) {
	ß.EventEnter()
	defer ß.EventLeave()

	ß.InfoT("text changed", String("Old", oldStr), String("New", newStr))
	ß.AttentionErr("are you reading?", io.EOF)
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

const EXPECTED_ALT_STR string = "A LEG"

func nonInlinedValidator(s string) bool {
	ß.OnValidating() // writes to call tree (file)

	return s == EXPECTED_ALT_STR
}

// catether functionality (extra, free-form log file)
func doCatether(Γ zlog.ILogCatheterable) {
	Γ.SetCatheterFile("/tmp/zlog_footprints.log")
	Γ.PrintCatheter("This is complimentary information", mtag.MoneyEUR("Balance", '€', 4850.283))
}

// event-driven logging
func oneFunction(Ω zlog.ILogEventDriven) int {
	Ω.Enter()
	defer Ω.Leave()

	Ω.Visit()

	Ω.Step("1. Definition")
	x := 5
	y := 20

	Ω.Step("2. Operation")
	result := x * y

	Ω.Result("Returning with %d", result)

	return result
}

// Level-aware tagged logging
func twoFunction(Ω zlog.ILogTaggable) {
	Ω.TraceT("tagged trace", Byte("byte", byVal), Bool("True?", boVal), At())
	Ω.DebugT("tagged debug", Int("integer", iVal), String("text", stVal))
	Ω.InfoT("tagged info", Rune("letter", rVal), Bool("True?", boVal))
	Ω.WarnT("tagged warning", String("text", stVal), Err(err))
	Ω.ErrorT("tagged error", Int("integer", iVal), Rune("letter", rVal))
	Ω.FatalT(EXIT_CODE_3, "tagged fatal", Byte("byte", byVal), Bool("True?", boVal))
}

// Level-aware logging
func threeFunction(Ω zlog.ILogPrioritized) {
	Ω.Trace("test trace", rVal, stVal)
	Ω.Tracef("formatted %d trace %t", iVal, boVal)

	Ω.Debug("test debug", rVal, stVal)
	Ω.Debugf("formatted %d debug %t", iVal, boVal)

	Ω.Info("test info", rVal, stVal)
	Ω.Infof("formatted %d info %t", iVal, boVal)

	Ω.Warn("test warning", rVal, stVal)
	Ω.Warnf("formatted %d warning %t", iVal, boVal)

	Ω.Error("test error", rVal, stVal)
	Ω.Errorf("formatted %d error %t", iVal, boVal)
	Ω.ErrorE(err)

	// Leaving them out because it would exit the demo
	//Ω.Fatal(EXIT_CODE_1, "test fatal", rVal, stVal) // our Fatal exits with OUR exit code
	// unreachable because the previous calls os.Exit(EXIT_CODE_1)
	//Ω.Fatalf(EXIT_CODE_2, "formatted %d fatal %t", iVal, boVal)
}

// try all the tags but standalone
func tagsFunction() {
	fmt.Println("At", mtag.At())
	fmt.Println("Boolean", mtag.Bool("True?", boVal))
	fmt.Println("Byte", mtag.Byte("Binary", byVal))
	fmt.Println("Date-only", mtag.Date("Today", time.Now()))
	fmt.Println("Date-only UTC", mtag.DateUTC("Today", time.Now()))
	fmt.Println("Sample Error", mtag.Err(err))
	fmt.Println("Float", mtag.Float("Value", 4850.385, 4, 2))
	fmt.Println("Whole number", mtag.Int("Offset", 2800))
	fmt.Println("Letter/Rune", mtag.Rune("Letter", rVal))
	fmt.Println("String value", mtag.String("Name", stVal))
	fmt.Println("Time only", mtag.Time("Now", time.Now()))
	fmt.Println("Time only", mtag.TimeUTC("Now UTC", time.Now()))
	fmt.Println("Yes/No", mtag.YesNo("Answer", boVal))
	fmt.Println("Money (American)", mtag.Money("Balance", '$', 8700.568))
	fmt.Println("Money (European)", mtag.MoneyEUR("Balance", '€', 8700.568))
	fmt.Println("Money (Swiss)", mtag.MoneyCHE("Balance", '$', 8700.568))
}

// ZLog comes with built-in defaults. During startup, the ZLog package checks if
// there are Environment variables to override Logging level & Filename output.
// Furthermore, this demo shows how those can then be overriden in the main
// application during startup.
func OverrideLogConfig(omitLogging bool, useLevelStr string, omitTimestamp bool) {
	var err error
	if omitLogging {
		// - to disable logging
		ß.SetOutput(io.Discard) // by default it does os.Stderr
	}
	if len(useLevelStr) != 0 {
		// if the level name is invalid it defaults to the current value
		if level, err = level.Parse(useLevelStr); err == nil {
			ß.SetLevel(level)
		}
	}
	if omitTimestamp {
		// Remove Ldate|Ltime
		newFlags := ß.Flags() &^ log.Ldate
		newFlags = newFlags &^ log.Ltime
		ß.SetFlags(newFlags)
	}
}

func OverrideFilterConfig(Ω zlog.ILogFilterable, create bool) {
	// BTW you can implement your own in a database as long as
	// your service implements ILogFilterService. Being explicit
	// down here for instructional purposes...
	var filterService zlog.ILogFilterService
	filterService = zlog.NewLogFilterService("demo_zlog", USER_CONFIG_SUBDIR)
	filterService.UseYaml()

	if create {
		if err := filterService.Init(); err != nil {
			println("Could not initialize Log Filter: ", err)
			os.Exit(100)
		} else {
			fmt.Println("Log Filter File (Please tweak)")
			fmt.Printf("\tStatus: created\n")
			fmt.Printf("\tLocation: %s\n", filterService.String())
			os.Exit(0)
		}
		// Since the initial configuration is not tailored to this
		// or any app, we are not going to load it. Thus we remain
		// with the default NullLogFilterService.
	} else {
		// the developer already has a tailored Filter config file.
		Ω.UseFilter(filterService)
	}
}

// Because people usually don't read...
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
	defer ß.Close()

	goapp.ModuleVersion.Copyright("lordofscripts", CHR_TRIDENT)

	// A. Command-line flags
	var noLog, noTimestamp, firstTime, help bool
	var levelStr string
	flag.BoolVar(&help, "help", false, "Help me")
	flag.BoolVar(&noLog, "nolog", false, "disable logging")
	flag.BoolVar(&noTimestamp, "no-timestamp", false, "log entries without timestamp")
	flag.BoolVar(&firstTime, "init", false, "Initialize filter config file")
	flag.StringVar(&levelStr, "level", "", "Set logging level")
	flag.Parse()

	// B. command-line override(s)
	//	- NOTE: init() has already executed so ß is already initialized.
	OverrideLogConfig(noLog, levelStr, noTimestamp)
	OverrideFilterConfig(ß, firstTime)

	fmt.Println("ZLog is a Flexible Logger from GoApp")
	fmt.Printf("ZLog Level: %s (%d)\n", level, level)

	// C. Logging for Event-driven applications (some to try)
	const EXPECTED_STR string = "JAMON"
	myObj := newDemoObject("initial")
	myObj.Validator = func(s string) bool {
		ß.OnValidating() // writes to call tree (file)

		return s == EXPECTED_STR
	}
	myObj.SampleMethod(stVal, iVal)
	myObj.SetText("Turkey")     // will not validate
	myObj.SetText(EXPECTED_STR) // validates
	myObj.Validator = nonInlinedValidator
	myObj.SetText(EXPECTED_ALT_STR) // validates

	// D. Try by restricting to one interface
	oneFunction(ß)   // event-driven demo
	twoFunction(ß)   // level-aware taggable
	threeFunction(ß) // level-aware (non-tagged)
	tagsFunction()   // try standalone mtag.*

	doCatether(ß) // extra free-form log

	// The first of these exits with that application exit code

	ß.FatalT(EXIT_CODE_3, "tagged fatal", Byte("byte", byVal), Bool("True?", boVal))

	fmt.Print("This will never be printed")
}

/*
  I. The Call Tree will look like this:

		>>>>>  ZLog Call Tree <<<<<<
		⏳  Local time: 2026-08-21 18:13:08
		⏳  Zulu time : 2026-08-21 23:13:08 UTC
		⏳  Timezone  : EST UTC-05:00
		⚙️   ZLog (debug)
		2 VAL 🢰   main▪mai⯈func1()
		1 UPD ⏿  main▪DemoObject🡪SetText()
		2 VAL 🢰   main▪mai⯈func1()
		1 UPD ⏿  main▪DemoObject🡪SetText()
		2 VAL 🢰   main▪nonInlinedValidator()
		1 UPD ⏿  main▪DemoObject🡪SetText()
		1 ENT ❯  main▪oneFunction()
		1 LEA ❮  main▪oneFunction()

  II. The console log may look like this:

		⏳  Local time: 2026-08-21 18:18:06
		⏳  Zulu time : 2026-08-21 23:18:06 UTC
		⏳  Timezone  : EST UTC-05:00
		⚙️   ZLog (debug)
			🔱 goApp v1.4.6 lordofscripts 🔱
		ZLog is a Flexible Logger from GoApp
		ZLog Level: Trace (1)
		2026/08/21 18:18:06 ZLOG ⟫ (Ctor) main:newDemoObject()
		2026/08/21 18:18:06 ZLOG [DBG] Sample methodThiß Über Allesinitial
		2026/08/21 18:18:06 ZLOG 🢰 VAL (2) main▪mai⯈func1()
		2026/08/21 18:18:06 ZLOG ⏿ UPD (1) main▪DemoObject🡪SetText()
		2026/08/21 18:18:06 ZLOG [DBG] Setting text to 'Turkey'
		2026/08/21 18:18:06 ZLOG [INF] text changed Old='initial' New='Turkey'
		2026/08/21 18:18:06 ZLOG are you reading?EOF
		2026/08/21 18:18:06 ZLOG 🢰 VAL (2) main▪mai⯈func1()
		We like 'JAMON'
		2026/08/21 18:18:06 ZLOG ⏿ UPD (1) main▪DemoObject🡪SetText()
		2026/08/21 18:18:06 ZLOG [DBG] Setting text to 'JAMON'
		2026/08/21 18:18:06 ZLOG [INF] text changed Old='initial' New='JAMON'
		2026/08/21 18:18:06 ZLOG are you reading?EOF
		2026/08/21 18:18:06 ZLOG 🢰 VAL (2) main▪nonInlinedValidator()
		We like 'A LEG'
		2026/08/21 18:18:06 ZLOG ⏿ UPD (1) main▪DemoObject🡪SetText()
		2026/08/21 18:18:06 ZLOG [DBG] Setting text to 'A LEG'
		2026/08/21 18:18:06 ZLOG [INF] text changed Old='initial' New='A LEG'
		2026/08/21 18:18:06 ZLOG are you reading?EOF
		2026/08/21 18:18:06 ZLOG ❯ main:oneFunction()
		2026/08/21 18:18:06 ZLOG ❮❯ main:oneFunction()
		2026/08/21 18:18:06 ZLOG 👣 1. Definition
		2026/08/21 18:18:06 ZLOG 👣 2. Operation
		2026/08/21 18:18:06 ZLOG ⇶ main:oneFunction() Returning with 100
		2026/08/21 18:18:06 ZLOG ❮ main:oneFunction()
		2026/08/21 18:18:06 ZLOG [DBG] tagged debug integer=800 text='Thiß Über Alles'
		2026/08/21 18:18:06 ZLOG [INF] tagged info letter='ß' (0xDF) True?=false
		2026/08/21 18:18:06 ZLOG [WRN] tagged warning text='Thiß Über Alles' Error=*errors.errorString=>oops, bad thing happened.
		2026/08/21 18:18:06 ZLOG [ERR] tagged error integer=800 letter='ß' (0xDF)
		2026/08/21 18:18:06 ZLOG [DIE] tagged fatal byte=0x05 True?=false
*/
