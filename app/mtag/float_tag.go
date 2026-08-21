/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Lord of Scripts
 *							goApp::mtag
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Key-Value pair tags used in MLog and ZLog but also useful in your
 * everyday application.
 * Implements: Float tag
 *-----------------------------------------------------------------*/
package mtag

import "fmt"

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ ILogKeyValuePair = (*kvFloat[float64])(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

// A float-type key-value pair tag
type kvFloat[T float32 | float64] struct {
	k           string
	v           T
	floatFormat string
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// implements fmt.Stringer for mtag.Float() rendering the value with
// the number of whole digits and decimal digits specified in the
// Float() function call.
func (k *kvFloat[T]) String() string {
	formatted := fmt.Sprintf(k.floatFormat, k.v)
	return fmt.Sprintf("%s=%s", k.k, formatted)
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// log the float32/64 key=value pair. Specify both the nr. of whole
// digits (left side of decimal point), and nr. of decimals.
func Float[T float32 | float64](key string, value T, whole, decimals byte) ILogKeyValuePair {
	format := fmt.Sprintf("%%%d.%df", whole+decimals+1, decimals) // %11.6f
	return &kvFloat[T]{key, value, format}
}
