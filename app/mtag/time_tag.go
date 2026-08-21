/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Lord of Scripts
 *							goApp::mtag
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Key-Value pair tags used in MLog and ZLog but also useful in your
 * everyday application.
 * Implements: TimeOnly tag
 *-----------------------------------------------------------------*/
package mtag

import (
	"fmt"
	"time"
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ ILogKeyValuePair = (*kvTime)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

// A Time-only key-value pair tag.
type kvTime struct {
	k   string
	v   time.Time
	utc bool
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// implements fmt.Stringer for mtag.Time() rendering Local or UTC
func (k *kvTime) String() string {
	const CUSTOM_TIME_FORMAT = "15:04:05"
	if !k.utc {
		return fmt.Sprintf("%s=%s", k.k, k.v.Format(CUSTOM_TIME_FORMAT))
	} else {
		return fmt.Sprintf("%s=%s", k.k, k.v.UTC().Format(CUSTOM_TIME_FORMAT))
	}
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// log the time key=value pair
func Time(key string, value time.Time) ILogKeyValuePair {
	return &kvTime{key, value, false}
}

// log the time key=value pair in UTC/Zulu time
func TimeUTC(key string, value time.Time) ILogKeyValuePair {
	return &kvTime{key, value, true}
}
