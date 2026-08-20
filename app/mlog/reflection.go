/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                          GoApp
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package mlog

import (
	"regexp"
	"runtime"
	"strings"
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

/* ----------------------------------------------------------------
 *                 P U B L I C    M E T H O D S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                 P R I V A T E    M E T H O D S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

/**
 * Gets the current nesting level with main.main() as reference and
 * a string with the location (package/struct/function).
 */
func GetNestingLevel(frCnt ...int) (int, string) {
	// moved here from app/logx package
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
	// moved here from app/logx package
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
