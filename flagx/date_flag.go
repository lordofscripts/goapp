/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * A custom Date Flag for the GO flag package. You can specify dates
 * in various formats as defined by your application.
 *-----------------------------------------------------------------*/
package flagx

import (
	"flag"
	"fmt"
	"strings"
	"time"
)

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/
const ()

var ()

/* ----------------------------------------------------------------
 *                     I N T E R F A C E S
 *-----------------------------------------------------------------*/

var _ flag.Value = (*DateFlag)(nil)

/* ----------------------------------------------------------------
 *                         T Y P E S
 *-----------------------------------------------------------------*/

type DateFlag struct {
	Value                     time.Time
	IsSet                     bool
	formats                   []string
	hasYear, hasMonth, hasDay bool
}

/* ----------------------------------------------------------------
 *                    C O N S T R U C T O R S
 *-----------------------------------------------------------------*/

func NewDateVar(formats ...string) *DateFlag {
	var fmts []string
	if len(formats) == 0 {
		fmts = []string{"2006-01-02"} // yyyy-MM-dd
	} else {
		fmts = formats
	}

	return &DateFlag{
		Value:    time.Now().UTC(),
		IsSet:    false,
		formats:  fmts,
		hasYear:  false,
		hasMonth: false,
		hasDay:   false,
	}
}

/* ----------------------------------------------------------------
 *                        M E T H O D S
 *-----------------------------------------------------------------*/

// implements fmt.Stringer rendering the set date in the first format
// of the constructor list, else defaults to YYYY-MM-DD
func (r *DateFlag) String() string {
	if r.IsSet {
		return r.Value.Format(r.formats[0])
	}

	return ""
}

/* ----------------------------------------------------------------
 *                 P U B L I C    M E T H O D S
 *-----------------------------------------------------------------*/

// invoked by flag.Parse() to set the value. We attempt to convert the
// date with "now" or "ahora" to a full date/time, "today" and "hoy"
// to today's date without time, or to the first interpretation of the
// formats given at the constructor. Else an error is produced. Time
// is in UTC.
func (r *DateFlag) Set(value string) error {
	switch value {
	case "now", "ahora":
		r.hasYear, r.hasMonth, r.hasDay = true, true, true
		r.Value = time.Now().UTC()
		r.IsSet = true

	case "today", "hoy":
		r.hasYear, r.hasMonth, r.hasDay = true, true, true
		now := time.Now().UTC()
		r.Value = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		r.IsSet = true

	default:
		isValid := false
		for i, layout := range r.formats {
			if fecha, err := time.Parse(layout, value); err == nil {
				r.Value = fecha
				isValid = true
				r.IsSet = true
				r.formats[i] = r.formats[0]
				r.formats[0] = layout // becomes the output format
				r.analyse()
				break
			}
		}

		if !isValid {
			return fmt.Errorf("none of the %d layouts was good for parsing the date", len(r.formats))
		}
	}

	return nil
}

// returns an indication of what the default format has (Year,Month,Day)
func (r *DateFlag) Has() (year, month, day bool) {
	year = r.hasYear
	month = r.hasMonth
	day = r.hasDay
	return
}

/* ----------------------------------------------------------------
 *                 P R I V A T E    M E T H O D S
 *-----------------------------------------------------------------*/

// analyze the 1st format in the list
func (r *DateFlag) analyse() {
	r.hasYear = strings.Contains(r.formats[0], "2006")
	r.hasMonth = strings.Contains(r.formats[0], "01")
	r.hasDay = strings.Contains(r.formats[0], "02")
}

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

// registers an instance of DateFlag with the flag package for parsing.
func RegisterDateVar(r *DateFlag, name string, usage string, formats ...string) {
	r.Value = time.Now().UTC()
	r.IsSet = false
	if len(formats) == 0 {
		r.formats = []string{"2006-01-02"}
	} else {
		r.formats = formats
	}
	flag.Var(r, name, usage)
}

// A (local) date that corresponds to January 1st of the current year, no time.
func ThisYear() time.Time {
	return time.Date(time.Now().Year(), time.January, 1, 0, 0, 0, 0, time.Local)
}

// A (local) date that corresponds to the 1st of the current month, no time.
func ThisMonth() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
}

// A (local) date that corresponds to today. No time.
func Today() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
}

/* ----------------------------------------------------------------
 *                  M A I N    |    D E M O
 *-----------------------------------------------------------------*/
/*
func DemoDateFlag() {
	var myDate1 DateFlag = NewDateVar("2006-Jan-02")
	flag.Var(&myDate, "date", "custom Date value")
	flag.Parse()

	fmt.Printf("Date value: %s\n", myDate.Value)
}
*/
