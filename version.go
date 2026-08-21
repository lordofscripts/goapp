/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                         goApp Module
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Module version descriptor. Useful if the client application needs
 * to send problem reports by querying each of my module's version
 * as in this case goapp.ModuleVersion.String()
 *-----------------------------------------------------------------*/
package goapp

import "github.com/lordofscripts/goapp/app"

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/

const (
	_NAME    string = "goApp"
	_DESC    string = "Library of useful GO packages"
	_VERSION string = "1.5.0"
)

var (
	ModuleVersion app.PackageVersion = app.NewReleaseVersion(_NAME, _DESC, _VERSION)
)

var (
	// All GO modules that import GoApp to use PackageVersion,
	// can (optionally) use the Register method to have itself
	// added here. The end-user application can enumerate this
	// in the Help/Support information.
	CustomImports []string = make([]string, 0)
)

/* ----------------------------------------------------------------
 *                    I N I T I A L I A Z E R
 *-----------------------------------------------------------------*/

func init() {
	RegisterModule(ModuleVersion)
}

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

// Register so that it can be enumerated in Help info
func RegisterModule(mod app.PackageVersion) {
	CustomImports = append(CustomImports, mod.Info())
}

/* DO NOT REMOVE - Used by Makefile - v2
//>>>BEGIN Versioner
package main

import (
    "os"
    "fmt"
    "strings"
    "github.com/lordofscripts/goapp"
)

func main() {
    if len(os.Args) == 2 && strings.EqualFold(os.Args[1], "short") {
        fmt.Println(goapp.ModuleVersion.Short())
	} else if len(os.Args) == 2 && strings.EqualFold(os.Args[1], "version") {
		fmt.Println(goapp.ModuleVersion.Version())
    } else {
        fmt.Println(goapp.ModuleVersion)
    }
}
//>>>END Versioner
*/
