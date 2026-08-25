/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package util

import (
	"fmt"
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

// The runtime frame components parsed out of a runtime.Frame
type FrameComponents struct {
	Package           string
	Struct            string // or empty
	Function          string // function OR method
	Closure           string
	IsPointerReceiver bool
}

/* ----------------------------------------------------------------
 *                    C O N S T R U C T O R S
 *-----------------------------------------------------------------*/

// (Ctor) an instance of a runtime.Frame parser. On every parse the
// values are returned and updated into this same instance.
func NewFrameComponentParser() *FrameComponents {
	return &FrameComponents{}
}

/* ----------------------------------------------------------------
 *                        M E T H O D S
 *-----------------------------------------------------------------*/

// Parse the runtime frame and return its details. These same details
// are updated into the current instance as well.
func (fc *FrameComponents) Parse(frame runtime.Frame) FrameComponents {
	compies := ParseFrame(frame)
	fc.Package = compies.Package
	fc.Struct = compies.Struct
	fc.Function = compies.Function
	fc.Closure = compies.Closure
	fc.IsPointerReceiver = compies.IsPointerReceiver
	return compies
}

// implements fmt.Stringer and renders P/S/M/C or P/F/C
// taking into consideration if it is a function or method,
// and whether it is a pointer receiver or not.
func (fc *FrameComponents) String() string {
	if fc.Struct != "" {
		ptr := ' '
		if fc.IsPointerReceiver {
			ptr = '*'
		}
		return fmt.Sprintf("P: %s S:%c%s M: %s C: %s",
			fc.Package,
			ptr,
			fc.Struct,
			fc.Function,
			fc.Closure)
	} else {
		return fmt.Sprintf("P: %s F: %s C: %s",
			fc.Package,
			fc.Function,
			fc.Closure)
	}
}

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

// Parses a stack frame and returns the details
//   - main.MyFunction
//   - main.MyFunction.func1
//   - main.(*MyStruct).MyMethod.func15.1
//   - github.com/foot/bar.pkg.func2.3.4
func ParseFrame(frame runtime.Frame) FrameComponents {
	fullName := frame.Function
	if fullName == "" {
		return FrameComponents{}
	}

	var closure string
	closureParserRegex := regexp.MustCompile(`\.(func\d+(?:\.\d+)*)$`)
	// 1. detect and extract closure string
	closureMatches := closureParserRegex.FindStringSubmatch(fullName)
	if len(closureMatches) > 1 {
		closure = closureMatches[1]
		// strip the entire closure match with leading dot from fullName
		fullName = fullName[:strings.LastIndex(fullName, "."+closure)]
	}

	// 2. isolate package path
	lastSlash := strings.LastIndex(fullName, "/")
	var remaining string
	var pkgPath string

	if lastSlash == -1 {
		remaining = fullName
	} else {
		pkgPath = fullName[:lastSlash]
		remaining = fullName[lastSlash+1:]
	}

	//3. separate the base package name from the rest
	firstDot := strings.Index(remaining, ".")
	if firstDot != -1 {
		if pkgPath != "" {
			pkgPath += "/" + remaining[:firstDot]
		} else {
			pkgPath = remaining[:firstDot]
		}
		remaining = remaining[firstDot+1:]
	}

	// 4. extract struct and function names
	var structName, funcName string
	isPointer := false
	if strings.Contains(remaining, ".") {
		parts := strings.SplitN(remaining, ".", 2)
		structPart := parts[0]
		funcName = parts[1]

		structPart = strings.TrimPrefix(structPart, "(")
		structPart = strings.TrimSuffix(structPart, ")")
		if strings.ContainsRune(structPart, '*') {
			isPointer = true
		}
		structPart = strings.TrimPrefix(structPart, "*")
		structName = structPart
	} else {
		funcName = remaining
	}

	return FrameComponents{
		Package:           pkgPath,
		Struct:            structName,
		Function:          funcName,
		Closure:           closure,
		IsPointerReceiver: isPointer,
	}
}
