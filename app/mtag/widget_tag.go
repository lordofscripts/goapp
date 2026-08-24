/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Lord of Scripts
 *							goApp::mtag
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Key-Value pair tags used in MLog and ZLog but also useful in your
 * everyday application.
 * Implements: Widget (Menu/Button/Widget) tag
 *-----------------------------------------------------------------*/
package mtag

import "fmt"

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ ILogKeyValuePair = (*kvWidget)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

// A string-type key-value pair tag
type kvWidget struct {
	k string
	v string
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// implements fmt.Stringer for mtag.String()
func (k *kvWidget) String() string {
	return fmt.Sprintf("%s='%s'", k.k, k.v)
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// log the Button item key=value pair.
// Menu("Open") renders as "BUTN Open"
func Button(btnStr string) ILogKeyValuePair {
	return &kvWidget{"BUTN", btnStr}
}

// log the Menu item key=value pair.
// Menu("File > Open") renders as "MENU File > Open"
func Menu(menuStr string) ILogKeyValuePair {
	return &kvWidget{"MENU", menuStr}
}
