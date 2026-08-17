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
	"time"
)

/* ----------------------------------------------------------------
 *                   P R I V A T E    T Y P E S
 *-----------------------------------------------------------------*/

// a chronometer timer for calculating a duration with an optional context
type chronometerItem struct {
	Name     string
	Started  time.Time
	Duration time.Duration
}

/* ----------------------------------------------------------------
 *                    C O N S T R U C T O R S
 *-----------------------------------------------------------------*/

// (Ctor) a new instance of a chronometer timer for measuring duration,
// such as executing an operation.
func NewChronometerItem(name string, autoStart bool) *chronometerItem {
	var current time.Time
	if autoStart {
		current = time.Now()
	}
	return &chronometerItem{
		Name:     name,
		Started:  current,
		Duration: 0,
	}
}

/* ----------------------------------------------------------------
 *                        M E T H O D S
 *-----------------------------------------------------------------*/

// Start a timer (only if it has not started, else ignores it)
func (ci *chronometerItem) Start() {
	if ci.Started.IsZero() {
		ci.Started = time.Now()
	}
}

// Stops a running timer and returns the calculated duration
func (ci *chronometerItem) Stop() time.Duration {
	if !ci.Started.IsZero() {
		ci.Duration = time.Since(ci.Started)
		return ci.Duration
	}
	return time.Duration(0)
}

// time has been started but it may (or may not) have been stopped
func (ci *chronometerItem) IsActive() bool {
	return !ci.Started.IsZero()
}

// Whether the timer has never been started
func (ci *chronometerItem) IsNotStarted() bool {
	return ci.Started.IsZero()
}

// time has been started but not stopped
func (ci *chronometerItem) IsRunning() bool {
	return !ci.Started.IsZero() && ci.Duration == 0
}

// Whether the timer has been started and then stopped
func (ci *chronometerItem) IsStopped() bool {
	return !ci.Started.IsZero() && ci.Duration > 0
}

// implements fmt.Stringer by rendering an informative icon
// together with the timer name and duration.
func (ci *chronometerItem) String() string {
	const VARIATION_SELECTOR_16 rune = rune(0xfe0f)
	const ICON_TIMER_RUNNING rune = rune(0x23F3) // ⏳
	const ICON_TIMER_ACTIVE rune = rune(0x23F1)  // ⏱
	const ICON_TIMER_STOPPED rune = rune(0x23F9) // ⏹
	const ICON_TIMER_ACTIVE_COLOR string = "\u23F1\uFE0F"
	const ICON_TIMER_STOPPED_COLOR string = "\u23F9\uFE0F"
	const ICON_RED_CIRCLE rune = rune(0x1f534) // 🔴
	var visual string
	switch {
	// Never started
	case ci.IsNotStarted():
		return fmt.Sprintf("%c  %-20s not-started", ICON_RED_CIRCLE, ci.Name)

		// Started and running
	case ci.IsRunning():
		return fmt.Sprintf("%c  %-20s %s...", ICON_TIMER_RUNNING, ci.Name, ci.Duration)

		// Started and stopped
	case ci.IsStopped():
		return fmt.Sprintf("%s  %-20s %s", ICON_TIMER_ACTIVE_COLOR, ci.Name, ci.Duration)

	default:
		visual = ""
	}
	return visual
}

/* ----------------------------------------------------------------
 *                  M A I N    |    D E M O
 *-----------------------------------------------------------------*/

/*
func demoTimer() {
	tmrX := NewChronometerItem("Sample", true)
	// Run a specific timer
	tmrX.Start()
	//LongRunningOperation()
	tmrX.Stop()
	fmt.Println("Duration", tmrX)
}
*/
