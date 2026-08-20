/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Lord of Scripts
 *							goApp::mtag
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Key-Value pair tags used in MLog and ZLog but also useful in your
 * everyday application.
 * Implements: At tag (uses reflection to get code location info)
 *-----------------------------------------------------------------*/
package mtag

import "fmt"

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ ILogKeyValuePair = (*kvAt)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

// A Code location tag that uses location to obtain information about
// the location (package, object, method/function, line) where it
// was called.
type kvAt struct {
	v *CallerInfo
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// implements fmt.Stringer for mlog.At()
func (k *kvAt) String() string {
	return fmt.Sprintf("At=%s", k.v)
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// log the current package/method/function/line location
func At() ILogKeyValuePair {
	return &kvAt{RetrieveCallerInfo(FRAMENR_CALLER)}
}
