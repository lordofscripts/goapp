/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Lord of Scripts
 *							goApp::mtag
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Key-Value pair tags used in MLog and ZLog but also useful in your
 * everyday application.
 * Implements: Yes/No tag
 *-----------------------------------------------------------------*/
package mtag

import "fmt"

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ ILogKeyValuePair = (*kvYesNo)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

// A variation of a Bool tag that renders as a Yes (true) & No (false)
type kvYesNo struct {
	k string
	v bool
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// implements fmt.Stringer for mtag.YesNo()
func (k *kvYesNo) String() string {
	s := "No"
	if k.v {
		s = "Yes"
	}
	return fmt.Sprintf("%s=%s", k.k, s)
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// log the boolean key=value pair as a Yes/No value
func YesNo(key string, value bool) ILogKeyValuePair {
	return &kvYesNo{key, value}
}
