//go:build logx
// +build logx

/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Lord of Scripts
 *							   goApp
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Development/Debug version. Log module equivalents are mimicked
 * as-is, extra functionality and the equivalent Print*() functions
 * are filter-enabled.
 *-----------------------------------------------------------------*/
package logx

import (
	"fmt"
	"io"
	"log"
	"regexp"
	"runtime"
	"strings"
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
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

func Prefix() string {
	return log.Prefix()
}

func SetPrefix(prefix string) {
	log.SetPrefix(prefix)
	log.SetFlags(log.LstdFlags | log.Lmsgprefix)
}

func Flags() int {
	return log.Flags()
}

func Print(v ...any) {
	if _, _, allowed := getCallerFlexFiltered(2); allowed {
		v = append([]any{string(UC_FOOTSTEPS), " "}, v...)
		log.Print(v...)
	}
}

func Printf(format string, v ...any) {
	if _, _, allowed := getCallerFlexFiltered(2); allowed {
		log.Printf(string(UC_FOOTSTEPS)+" "+format, v...)
	}
}

func Println(v ...any) {
	if _, _, allowed := getCallerFlexFiltered(2); allowed {
		v = append([]any{string(UC_FOOTSTEPS), " "}, v...)
		log.Println(v...)
	}
}

func Fatal(v ...any) {
	log.Fatal(v...)
}

func Fatalf(format string, v ...any) {
	log.Fatalf(format, v...)
}

func Fatalln(v ...any) {
	log.Fatalln(v...)
}

//func Output(calldepth int, s string) error {}

func Panic(v ...any) {
	log.Panic(v...)
}

func Panicf(format string, v ...any) {
	log.Panicf(format, v...)
}

func Panicln(v ...any) {
	log.Panicln(v...)
}

func SetFlags(flag int) {
	log.SetFlags(flag)
}

func SetOutput(w io.Writer) {
	log.SetOutput(w)
}

/* ----------------------------------------------------------------
 *				E x t e n d e d 	F u n c t i o n s
 *-----------------------------------------------------------------*/

/**
 * For use to log that an object constructor is executing
 */
func Ctor() {
	if loc, _, allowed := getCallerFlexFiltered(2); allowed {
		log.Print("⟫ (Ctor) ", loc)
	}
}

/**
 * For use at entering a function/method of an EVENT callback
 */
func EventEnter() {
	if loc, _, allowed := getCallerFlexFiltered(2); allowed {
		log.Print("❯ (Event) ", loc)
	}
}

/**
 * For use at the end of a function/method of an EVENT callback
 */
func EventLeave() {
	if loc, _, allowed := getCallerFlexFiltered(2); allowed {
		log.Print("❮ (Event) ", loc)
	}
}

/**
 * For use at entering a function/method
 */
func Enter() {
	if loc, _, allowed := getCallerFlexFiltered(2); allowed {
		log.Print("❯ ", loc)
	}
}

/**
 * For use at the end of a function/method
 */
func Leave() {
	if loc, _, allowed := getCallerFlexFiltered(2); allowed {
		log.Print("❮ ", loc)
	}
}

/**
 * For use when entering a function/method where you don't want Enter+Leave
 * but just logging your visit.
 */
func Visit() {
	if loc, _, allowed := getCallerFlexFiltered(2); allowed {
		log.Print("❮❯ ", loc)
	}
}

func Step(message string) {
	if _, _, allowed := getCallerFlexFiltered(2); allowed {
		log.Printf("%c %s", UC_FOOTSTEPS, message)
	}
}

// Log an error irrespective of the log level
func AttentionAlways(short string, err error) {
	_, pkgShort, _ := getCallerFlexFiltered(2)
	log.Printf("%c %s ERR %s", UC_EYES, pkgShort, err)
}

// Log an error but conditionally filtered. See AttentionAlways()
func Attention(short string, err error) {
	if _, pkgShort, allowed := getCallerFlexFiltered(2); allowed {
		log.Printf("%c %s ERR %s", UC_EYES, pkgShort, err)
	}
}

/**
 * To display the result or return value of a function/method.
 */
func Result(format string, v ...any) {
	if loc, _, allowed := getCallerFlexFiltered(2); allowed {
		log.Printf(string(UC_ARROWS3)+" "+loc+" "+format, v...)
	}
}

func OnValidating() {
	nl, nm := GetNestingLevel(3)
	log.Printf("%c VAL (%d) %s", UC_VALIDATE, nl, nm)
	SingLogGate.WriteCallTree(fmt.Sprintf("%3d VAL %c %s%s", nl, UC_VALIDATE, strings.Repeat(" ", nl), nm))
}

func OnChanged(toValue ...any) {
	nl, nm := GetNestingLevel(3)
	var valueStr string = ""
	if len(toValue) == 1 {
		valueStr = fmt.Sprintf(" %c %v", UC_ARROWS3, toValue)
	}
	log.Printf("%c CHG (%d) %s", UC_CHANGE, nl, nm)
	SingLogGate.WriteCallTree(fmt.Sprintf("%3d CHG %c %s%s%s", nl, UC_CHANGE, strings.Repeat(" ", nl), nm, valueStr))
}

func OnUpdate() {
	nl, nm := GetNestingLevel(3)
	log.Printf("%c UPD (%d) %s", UC_OBSERVER, nl, nm)
	SingLogGate.WriteCallTree(fmt.Sprintf("%3d UPD %c %s%s", nl, UC_OBSERVER, strings.Repeat(" ", nl), nm))
}

func OnCascade(to string, val any) {
	nl, nm := GetNestingLevel(3)
	log.Printf("%c SET (%d) %s TO %s = %v", UC_ARROWS3, nl, nm, to, val)
	SingLogGate.WriteCallTree(fmt.Sprintf("%3d SET %c %s%s %s 🟰 %v", nl, UC_ARROWS3, strings.Repeat(" ", nl), nm, to, val))
}

func OnClick() {
	nl, nm := GetNestingLevel(3)
	log.Printf("%c CLK (%d) %s", UC_CLICK, nl, nm)
	SingLogGate.WriteCallTree(fmt.Sprintf("%3d CLK %c %s%s", nl, UC_CLICK, strings.Repeat(" ", nl), nm))
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

		// Filtered out? don't bother with the rest
		if SingLogGate.IsFilteredObject(packageS, objectS) {
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

/*
// @todo fix this output format
// @lordofscripts/bulkrename/flow.(*Workflow).recurseDirectory()#139
func getCallerFlex(stackIdx int) string {
	// 3 works when test.Announce() is used
	// 2 works when ShowCaseOK/Failed called directly from Test*()
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
			objectS = regex.ReplaceAllString(objectS, "") + "."
		}

		var funcS string
		if strings.Contains(info, ".init.") {
			funcS = info[:lastDot]
		} else {
			funcS = info[lastDot+1:]
		}
		//return fmt.Sprintf("@%s.%s()#%d\n", packageS, funcS, line)
		return fmt.Sprintf("%s%s()\n", objectS, funcS)
		//fmt.Printf("called from %s\n", details.Name())
		//return filename + "+" + string(line) + " " + packageS + ":" + funcS
	}
	return ""
}
*/
