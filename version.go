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
	_DESC    string = "Library of useful GO gadgets"
	_VERSION string = "1.4.4"
)

var (
	ModuleVersion app.PackageVersion = app.NewReleaseVersion(_NAME, _DESC, _VERSION)
)

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

/* DO NOT REMOVE - Used by Makefile
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
    } else {
        fmt.Println(goapp.ModuleVersion)
    }
}
//>>>END Versioner
*/
