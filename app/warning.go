/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * A warning object that can also be in the guise of an error instance
 * that can be checked/asserted via IsWarning().
 *-----------------------------------------------------------------*/
package app

import (
	"fmt"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

// Ensure the Warning type implements the error interface
var _ error = (*Warning)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

// The warning code number. In the user's application it can be used
// as an iota enumeration (custom warning codes).
type WarningCode = int

// Warning struct represents a custom warning type
type Warning struct {
	Message string
	Code    WarningCode
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// A warning in the guise of an error, use the IsWarning() helper
// function to check and type-assert.
func NewWarningAsErr(wcode WarningCode, format string, args ...any) error {
	return &Warning{
		Message: fmt.Sprintf(format, args...),
		Code:    wcode,
	}
}

// A pure warning but still implements errors.error
func NewWarning(wcode WarningCode, format string, args ...any) *Warning {
	return &Warning{
		Message: fmt.Sprintf(format, args...),
		Code:    wcode,
	}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// Implement the error interface
func (w *Warning) Error() string {
	return fmt.Sprintf("WARN-%05d: %s", w.Code, w.Message)
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// tells whether a given error instance represents a Warning and
// if so it casts it using type assertion. Else returns nil,false.
func IsWarning(err error) (*Warning, bool) {
	if err != nil {
		if warning, ok := err.(*Warning); ok {
			return warning, true
		}
	}

	return nil, false
}
