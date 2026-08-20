/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Lord of Scripts
 *							goApp::mtag
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Key-Value pair tags used in MLog and ZLog but also useful in your
 * everyday application.
 * Implements: Rune tag
 *-----------------------------------------------------------------*/
package mtag

import (
	"fmt"
	"unicode"
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ ILogKeyValuePair = (*kvRune)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

// A Rune-type key-value pair tag. It Stringer implementation renders
// taking into consideration whether the rune is printable or not.
type kvRune struct {
	k string
	v rune
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// implements fmt.Stringer for mtag.Rune()
func (k *kvRune) String() string {
	var out string
	if unicode.IsPrint(k.v) {
		out = fmt.Sprintf("%s='%c' (0x%X)", k.k, k.v, k.v)
	} else {
		out = fmt.Sprintf("%s=*** (0x%X)", k.k, k.v)
	}
	return out
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// log the Rune=value
func Rune(key string, value rune) ILogKeyValuePair {
	return &kvRune{key, value}
}
