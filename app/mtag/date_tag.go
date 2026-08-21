/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Lord of Scripts
 *							goApp::mtag
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Key-Value pair tags used in MLog and ZLog but also useful in your
 * everyday application.
 * Implements: Date tag
 *-----------------------------------------------------------------*/
package mtag

import (
	"fmt"
	"time"
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ ILogKeyValuePair = (*kvDate)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

// A Date-only key-value pair tag.
type kvDate struct {
	k   string
	v   time.Time
	utc bool
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// implements fmt.Stringer for mtag.Date() rendering Local or UTC
func (k *kvDate) String() string {
	const CUSTOM_DATE_FORMAT = "2006-01-02 15:04:05"
	if !k.utc {
		return fmt.Sprintf("%s=%s", k.k, k.v.Format(CUSTOM_DATE_FORMAT))
	} else {
		return fmt.Sprintf("%s=%s", k.k, k.v.UTC().Format(CUSTOM_DATE_FORMAT))
	}
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// log the date key=value pair
func Date(key string, value time.Time) ILogKeyValuePair {
	return &kvDate{key, value, false}
}

// log the date key=value pair in UTC/Zulu date
func DateUTC(key string, value time.Time) ILogKeyValuePair {
	return &kvDate{key, value, true}
}
