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
	_DESC    string = "a minimal GO CLI mini-framework"
	_VERSION string = "1.4.0"
)

var (
	ModuleVersion app.PackageVersion = app.NewReleaseVersion(_NAME, _DESC, _VERSION)
)

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/
