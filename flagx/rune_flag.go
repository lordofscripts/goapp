/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * A custom Rune Flag for the GO flag package. We can now use rune
 * flags in the command line. For that we implement the flag.Value
 * interface.
 * This implementation works with both single and multi-byte runes.
 *-----------------------------------------------------------------*/
package flagx

import (
	"flag"
	"fmt"
	"unicode/utf8"
)

/* ----------------------------------------------------------------
 *                     I N T E R F A C E S
 *-----------------------------------------------------------------*/

var _ flag.Value = (*RuneFlag)(nil)

/* ----------------------------------------------------------------
 *                         T Y P E S
 *-----------------------------------------------------------------*/

type RuneFlag struct {
	Value rune
	IsSet bool
}

/* ----------------------------------------------------------------
 *                        M E T H O D S
 *-----------------------------------------------------------------*/

func (r *RuneFlag) String() string {
	if r.IsSet {
		return string(r.Value)
	}
	return ""
}

func (r *RuneFlag) Set(value string) error {
	if utf8.RuneCountInString(value) != 1 {
		return fmt.Errorf("invalid rune: %s", value)
	}

	r.Value = []rune(value)[0]
	r.IsSet = true
	return nil
}

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

func RegisterRuneVar(r *RuneFlag, name string, value rune, usage string) {
	r.Value = value
	r.IsSet = false
	flag.Var(r, name, usage)
}

/* ----------------------------------------------------------------
 *                  M A I N    |    D E M O
 *-----------------------------------------------------------------*/
/*
func DemoRuneFlag() {
	var myRune RuneFlag
	flag.Var(&myRune, "rune", "custom Rune value")
	flag.Parse()

	fmt.Printf("Rune value: %c\n", myRune.Value)
}
*/
