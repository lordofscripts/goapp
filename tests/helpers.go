/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2024 Dídimo Grimaldo T.
 *                           GoApp
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *				U n i t   T e s t H e l p e r s
 *-----------------------------------------------------------------*/
package tests

import (
	"fmt"
	"runtime"
	"strings"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/
const (
	OK = "[OK]"
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *					H e l p e r   F u n c t i o n s
 *-----------------------------------------------------------------*/

// unit test suite title
func SuiteTitle(t string) {
	fmt.Println(strings.Repeat("·", 10), t, strings.Repeat("·", 10))
}

// test case title
func Title(t string) {
	fmt.Println("\tTest Case (", t, ")")
}

func Show(k string, v any) {
	fmt.Printf("\t~%s: %v\n", k, v)
}

// a file that always exists on a native OS
func SureFile() string {
	filename := ""
	switch runtime.GOOS {
	case "windows":
		filename = `C:\Windows\cmd.exe`
	case "darwin":
		filename = `/System/Library/CoreServices/Finder.app`
	default:
		filename = `/usr/bin/sh`
	}
	return filename
}

// a directory that always exists on the native OS
func SureDirectory() string {
	filename := ""
	switch runtime.GOOS {
	case "windows":
		filename = `C:\Windows`
	case "darwin":
		filename = "/System/Library"
	default:
		filename = "/usr"
	}
	return filename
}
