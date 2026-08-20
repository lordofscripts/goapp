/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           go-App
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Utility functions for Error-related handling in a pretty fashion.
 *-----------------------------------------------------------------*/
package app

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/

// Beware that non-monospace fonts will render a small gap between
// the box characters in between lines!
const (
	CHR_HORIZ       rune = rune(0x2500)
	CHR_VERT        rune = rune(0x2502)
	CHR_INTERSEX    rune = rune(0x253C)
	CHR_INTERSEX_L  rune = rune(0x251C)
	CHR_INTERSEX_R  rune = rune(0x2524)
	CHR_UPPER_LEFT       = rune(0x250C)
	CHR_UPPER_RIGHT      = rune(0x2510)
	CHR_LOWER_LEFT       = rune(0x2514)
	CHR_LOWER_RIGHT      = rune(0x2518)
)

const (
	_LOG_WITH_DEFAULT preferredLogService = iota
	_LOG_WITH_MLOG
	_LOG_WITH_LOGX
)

// False by default, but when true, every Die*() call also logs
// the message via the `app.mlog` package.
var LogOnDeath bool = false

// preferred renderer for Messages
var messageRender messageRenderer = errorMsgHeadingBoxed

// preferred renderer for Errors
var errorRender errorRenderer = errorHeadingBoxed

/* ----------------------------------------------------------------
 *                         T Y P E S
 *-----------------------------------------------------------------*/

type preferredLogService = byte

type messageRenderer func(message string, exitCode int, isTerminal bool)
type errorRenderer func(err error, exitCode int, isTerminal bool)

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

// by default they render as a box, else a simpler non-boxed version
func ErrorMessageRenderAsBox(asBox bool) {
	if asBox {
		errorRender = errorHeadingBoxed
		messageRender = errorMsgHeadingBoxed
	} else {
		errorRender = errorHeading
		messageRender = errorMsgHeading
	}
}

// Death of an application by outputting a good-bye and setting
// the OS exit code. It is logged as fatal.
func Die(message string, exitCode int) {
	messageRender(message, exitCode, true)
	if LogOnDeath {
		altMsg := fmt.Sprintf("Died: YES ExitCode: %d Message: %s", exitCode, message)
		log.Print(altMsg) // because log*.Fatal uses log.Fatal which uses exitCode=1
		os.Exit(exitCode)
	} else {
		os.Exit(exitCode)
	}
}

func DieWith(exitCode int, format string, args ...any) {
	messageRender(fmt.Sprintf(format, args...), exitCode, true)
	os.Exit(exitCode)
}

// display the error and die with an exit code, logging it as Fatal.
func DieWithError(err error, exitCode int) {
	errorRender(err, exitCode, true)
	if LogOnDeath {
		altMsg := fmt.Sprintf("Died: YES ExitCode: %d Error: %v", exitCode, err)
		log.Print(altMsg) // because logx.Fatal uses log.Fatal which uses exitCode=1
		os.Exit(exitCode)
	} else {
		os.Exit(exitCode)
	}
}

// When the condition is met the warning message is printed
func Assert(condition bool, warnMessage string) {
	if condition {
		fmt.Fprintf(os.Stderr, "\n\t%c Assert Failed:\n\t%s\n", UC_RED_EXCLAMATION, warnMessage)
	}
}

// If the condition is met, the death message is printed and the
// application terminates with the exit code.
func AssertOrDie(condition bool, deathMessage string, exitCode int) {
	if condition {
		fmt.Fprintf(os.Stderr, "\n\t%c Assert Failed:", UC_RED_EXCLAMATION)
		messageRender(deathMessage, exitCode, true)
		os.Exit(exitCode)
	}
}

// prints the error message with the exit code but does NOT exit.
func AnnounceErrorMessage(message string, exitCode int) {
	messageRender(message, exitCode, false)
}

// prints the error and exit code but does NOT exit the application.
func AnnounceError(err error, exitCode int) {
	errorRender(err, exitCode, false)
}

/* ----------------------------------------------------------------
 *                 L O C A L   F U N C T I O N S
 *-----------------------------------------------------------------*/

// Pretty formatting of a (wrapped) error. The skeleton line is only
// outputted if isTerminal is true, meaning certain death of the application.
// The exitCode appears after "ERR-" and it automatically detects the
// package name and line number from which it was called.
// 💥 ERROR (ERR-005) 💥
// 🎯 From: main #57
// ·  a sad error message string, one per wrapped error
// 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀
// ------------------------------------------------------------
func errorHeading(err error, exitCode int, isTerminal bool) {
	fmt.Fprintf(os.Stderr, "💥 ERROR (ERR-%03d) 💥\n", exitCode)
	fmt.Fprintf(os.Stderr, "🎯 From: %s\n", callerPackage())
	for errC := err; errC != nil; errC = errors.Unwrap(errC) {
		fmt.Println("· ", errC)
	}
	if isTerminal {
		fmt.Fprintln(os.Stderr, strings.Repeat("💀 ", 20)) //🚫
	}
	fmt.Fprintln(os.Stderr, strings.Repeat("-", 60))
}

// Pretty formatting of a (wrapped) error. The skeleton line is only
// outputted if isTerminal is true, meaning certain death of the application.
// The exitCode appears after "ERR-" and it automatically detects the
// package name and line number from which it was called.
// 💥 ERROR (ERR-005) 💥
// 🎯 From: main #57
// ·  a sad error message string, one per wrapped error
// 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀 💀
// ------------------------------------------------------------
func errorMsgHeading(message string, exitCode int, isTerminal bool) {
	fmt.Fprintf(os.Stderr, "💥 ERROR (ERR-%03d) 💥\n", exitCode)
	fmt.Fprintf(os.Stderr, "🎯 From: %s\n", callerPackage())
	fmt.Println("· ", message)
	if isTerminal {
		fmt.Fprintln(os.Stderr, strings.Repeat("💀 ", 20)) //🚫
	}
	fmt.Fprintln(os.Stderr, strings.Repeat("-", 60))
}

func errorHeadingBoxed(err error, exitCode int, isTerminal bool) {
	topLINE(75)
	headingLINE(75, fmt.Sprintf("💥 ERROR (ERR-%03d) 💥", exitCode))
	midLINE(75)
	contentLINE(74, fmt.Sprintf("🎯 From: %s", callerPackage()))
	for errC := err; errC != nil; errC = errors.Unwrap(errC) {
		contentLINE(75, BulletWrap(75, errC.Error(), ""))
	}
	if isTerminal {
		//headingLINE(75, "  💀 💀 💀  ") //🚫
		contentLINE(75, Center("X x X", 75))
	}
	bottomLINE(75)
}

func errorMsgHeadingBoxed(message string, exitCode int, isTerminal bool) {
	topLINE(75)
	headingLINE(75, fmt.Sprintf("💥 ERROR (ERR-%03d) 💥", exitCode))
	midLINE(75)
	contentLINE(74, fmt.Sprintf("🎯 From: %s", callerPackage()))
	contentLINE(75, BulletWrap(75, message, ""))
	if isTerminal {
		//headingLINE(75, "  💀 💀 💀  ") //🚫
		contentLINE(75, Center("X x X", 75))
	}
	bottomLINE(75)
}

// Returns a string with the caller's package and line number.
func callerPackage() string {
	// PC, File, Line, OK
	const FROM_HERE = 2
	const VIA_DIE = 3
	pc, _, line, ok := runtime.Caller(VIA_DIE) // 0=this func, 1=Caller wrapper, 2=your caller
	if !ok {
		return ""
	}

	fn := runtime.FuncForPC(pc) // e.g. "github.com/me/mod/pkg.(*T).Method"
	if fn == nil {
		return ""
	}
	name := fn.Name()

	// Package path is everything up to the last '.'
	// (works for "path/pkg.Func" and "path/pkg.Type.Method")
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '.' {
			return fmt.Sprintf("%s #%d", name[:i], line)
		}
	}
	return ""
}

func Center(s string, width int) string {
	if len(s) >= width {
		return s
	}
	pad := width - len(s)
	left := pad / 2
	right := pad - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}

func BulletWrap(width int, s string, bullet string) string {
	if width <= 0 {
		return s
	}
	if bullet == "" {
		bullet = "•"
	}

	lines := strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
	if len(lines) == 0 {
		return ""
	}

	indent := strings.Repeat(" ", len([]rune(bullet))+1) // bullet + one space
	rendered := make([]string, 0, len(lines))

	// first line: bullet + first line text
	first := strings.TrimRight(lines[0], " \t")
	firstRunes := []rune(first)

	availFirst := max(width-(len([]rune(bullet))+1), 1)

	if len(firstRunes) <= availFirst {
		rendered = append(rendered, bullet+" "+first)
	} else {
		rendered = append(rendered, bullet+" "+string(firstRunes[:availFirst]))
		rest := firstRunes[availFirst:]
		// continue wrapping rest with indent, but only for this first line
		rendered = append(rendered, wrapWithIndent(indent, rest, width)...)
	}

	// remaining lines: if longer than width, indent past the bullet
	for i := 1; i < len(lines); i++ {
		txt := strings.TrimRight(lines[i], " \t")
		r := []rune(txt)
		if len(r) == 0 {
			rendered = append(rendered, "")
			continue
		}
		if len(r) <= width {
			rendered = append(rendered, string(r))
		} else {
			rendered = append(rendered, wrapWithIndent(indent, r, width)...)
		}
	}

	return strings.Join(rendered, "\n")
}

func wrapWithIndent(indent string, runes []rune, width int) []string {
	// wrap "already-specified" text: fill width- len(indent) on subsequent lines
	avail := max(width-len([]rune(indent)), 1)

	out := []string{}
	for len(runes) > 0 {
		n := avail
		if len(runes) < n {
			n = len(runes)
		}
		out = append(out, indent+string(runes[:n]))
		runes = runes[n:]
	}
	return out
}

// error/message box top line
func topLINE(width int) {
	fmt.Fprintf(os.Stderr, "    %c%s%c\n", CHR_UPPER_LEFT, strings.Repeat(string(CHR_HORIZ), width), CHR_UPPER_RIGHT)
}

// error/message box mid line between box heading and box content
func midLINE(width int) {
	fmt.Fprintf(os.Stderr, "    %c%s%c\n", CHR_INTERSEX_L, strings.Repeat(string(CHR_HORIZ), width), CHR_INTERSEX_R)
}

// error/message box bottom line
func bottomLINE(width int) {
	fmt.Fprintf(os.Stderr, "    %c%s%c\n", CHR_LOWER_LEFT, strings.Repeat(string(CHR_HORIZ), width), CHR_LOWER_RIGHT)
}

// error/message box heading line
func headingLINE(width int, content string) {
	fmt.Fprintf(os.Stderr, "    %c%s    %c\n", CHR_VERT, Center(content, width), CHR_VERT)
}

// error/message box content line
func contentLINE(width int, content string) {
	fmt.Fprintf(os.Stderr, "    %c%-*s%c\n", CHR_VERT, width, content, CHR_VERT)
}
