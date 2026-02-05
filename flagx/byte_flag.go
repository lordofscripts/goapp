/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * A custom Byte flag for the flag package that can be displayed as rune.
 *-----------------------------------------------------------------*/
package flagx

import (
	"flag"
	"fmt"
	"strconv"
)

/* ----------------------------------------------------------------
 *                     I N T E R F A C E S
 *-----------------------------------------------------------------*/

var _ flag.Value = (*ByteFlag)(nil)

/* ----------------------------------------------------------------
 *                         T Y P E S
 *-----------------------------------------------------------------*/

type ByteFlag struct {
	Value byte
	IsSet bool
}

/* ----------------------------------------------------------------
 *                        M E T H O D S
 *-----------------------------------------------------------------*/

func (r *ByteFlag) String() string {
	if r.IsSet {
		return string(r.Value)
	}
	return ""
}

func (r *ByteFlag) Set(value string) error {
	val, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("not an integer value: %s", value)
	}
	if val < 0 || val >= 256 {
		return fmt.Errorf("not a byte value: %d", val)
	}

	r.Value = byte(val)
	r.IsSet = true
	return nil
}

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

func RegisterByteVarSet(r *ByteFlag, name string, value byte, usage string) {
	r.Value = value
	r.IsSet = false
	flag.Var(r, name, usage)
}

func RegisterByteVar(r *ByteFlag, name string, usage string) {
	r.Value = 0
	r.IsSet = false
	flag.Var(r, name, usage)
}

/* ----------------------------------------------------------------
 *                  M A I N    |    D E M O
 *-----------------------------------------------------------------*/
/*
func DemoByteFlag() {
	var myByte ByteFlag
	flag.Var(&myByte, "byte", "custom Byte value")
	flag.Parse()

	fmt.Printf("Byte value: %c\n", myByte.Value)
}
*/
