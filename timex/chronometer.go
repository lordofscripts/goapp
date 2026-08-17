/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                       GoApp Library
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * NOTE: Promoted from Secretum NG
 *-----------------------------------------------------------------*/
package timex

import (
	"fmt"
	"strings"
)

/* ----------------------------------------------------------------
 *                   P U B L I C    T Y P E S
 *-----------------------------------------------------------------*/

// The callback function when Running a chronometer timer
type ChronoRunCallback func() (err error, exitCode int)

// The Error callback when a chronometer's Run task ends in error
type ChronoErrorCallback func(title string, errx error, exitCode int)

// Keeps timing several items by their duration and name.
type Chronometer struct {
	Items   map[string]*chronometerItem
	OnError ChronoErrorCallback
}

/* ----------------------------------------------------------------
 *                    C O N S T R U C T O R S
 *-----------------------------------------------------------------*/

// Prepare a multi-chronometer
func NewChronometer() *Chronometer {
	return &Chronometer{
		Items:   make(map[string]*chronometerItem, 0),
		OnError: nil}
}

/* ----------------------------------------------------------------
 *                 P U B L I C    M E T H O D S
 *-----------------------------------------------------------------*/

// Defines and auto-starts a new timer. It does NOT restart an existing timer.
func (ch *Chronometer) Start(name string) *chronometerItem {
	var item *chronometerItem = nil
	norm := strings.ToUpper(name)
	if _, ok := ch.Items[norm]; !ok {
		item = NewChronometerItem(name, true)
		ch.Items[norm] = item
	}
	return item
}

// Run a synchronous task that returns an error with an OS exit code
// or nil and zero on success.
func (ch *Chronometer) Run(name string, callback ChronoRunCallback) {
	const CHECKMARK_GREEN string = "\u2714\uFE0F" // ✔️
	// Start timer...
	if timer := ch.Start(name); timer != nil {
		// Execute action
		errx, exitCode := callback()
		// Stop and calculate duration
		timer.Stop()
		// Success or not
		if errx == nil {
			fmt.Printf("\t%s %s\n", timer.String(), CHECKMARK_GREEN)
		} else {
			// Error reporting
			fmt.Printf("\t%s\n", timer.String())
			if ch.OnError != nil {
				ch.OnError(name, errx, exitCode)
			}
		}
	}
}

// Implements fmt.Stringer on the chronometer, returning a detail
// of all timers and their state/duration as a string.
func (ch *Chronometer) String() string {
	var sb strings.Builder
	for v := range ch.Items {
		fmt.Fprintf(&sb, "%s\n", ch.Items[v])
	}
	return sb.String()
}

/* ----------------------------------------------------------------
 *                 P R I V A T E    M E T H O D S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                  M A I N    |    D E M O
 *-----------------------------------------------------------------*/

/*
func demoTime() {
	// Define the Chronometer manager
	chrono := util.NewChronometer()
	// Specify an error handling callback function
	chrono.OnError = func(title string, errx error, exitCode int) {
		fmt.Println("❌ ERROR", title, errx)
		os.Exit(exitCode)
	}
	// Run a task
	chrono.Run("Create-Shares", func() (errT error, exitCode int) {
		errT = CreateSharesStream(secretReader, masterSecret, TOTAL_SHARES, THRESHOLD_SHARES, fileShares)
		if errT != nil {
			exitCode = 1 // the exitCode to be passed to OnError
		}
		return errT, 0
	})

	// Run a specific timer
	tmrX := chrono.Start("Sample")
	LongRunningOperation()
	tmrX.Stop()
	fmt.Println("Duration", tmrX)
}
*/
