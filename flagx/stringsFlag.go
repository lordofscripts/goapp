/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * A Flag for the GO flag package. Useful when you want a flag/CLI option
 * to be any of a pre-determined list of strings.
 *-----------------------------------------------------------------*/
package flagx

import (
	"flag"
	"fmt"
	"slices"
	"strings"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

// custom flags must implement this interface
var _ flag.Value = (*StringsFlag)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type StringsFlag struct {
	Value      string
	IsSet      bool
	anyOf      []string
	defaultVal string
	strict     bool // when strict Set() only accepts anyOf[] values
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// implements flag.Value and fmt.Stringer
func (r *StringsFlag) String() string {
	if r.IsSet {
		return string(r.Value)
	}
	return ""
}

// implements flag.Value
func (r *StringsFlag) Set(value string) error {
	if !r.strict {
		r.Value = value
		r.IsSet = true
	} else {
		var ok bool
		if ok, value = r.isValidValue(value); ok {
			r.Value = value
			r.IsSet = true
		} else {
			return fmt.Errorf("flag accepts only: %v", r.anyOf)
		}
	}
	return nil
}

// when strict is set, the Set() function only accepts the values
// defined in SetChoices(), else it accepts any; thus behaving
// like a simple flag.StringVar
func (r *StringsFlag) Strict(strict bool) *StringsFlag {
	r.strict = strict
	return r
}

// set the string choices that will be considered valid. Any other
// values will be rejected by Set(). Values will be lowercased.
func (r *StringsFlag) SetChoices(choices []string) *StringsFlag {
	for idx, str := range choices {
		choices[idx] = strings.ToLower(strings.Trim(str, " \t"))
	}
	r.anyOf = choices
	return r
}

// return the list of allowed values (choices) and default value
func (r *StringsFlag) Help() string {
	all := strings.Join(r.anyOf, ",")
	return fmt.Sprintf("%s (default:'%s')", all, r.defaultVal)
}

func (r *StringsFlag) IsValid() bool {
	ok, _ := r.isValidValue(r.Value)
	return ok
}

func (r *StringsFlag) isValidValue(value string) (bool, string) {
	value = strings.ToLower(strings.Trim(value, " \t"))
	ok := slices.ContainsFunc(r.anyOf, func(f string) bool {
		return strings.EqualFold(f, value)
	})
	return ok, value
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

func RegisterSelectVar(r *StringsFlag, name string, value string, usage string) *StringsFlag {
	r.Value = value
	r.IsSet = false
	r.strict = false
	r.defaultVal = value
	flag.Var(r, name, usage)
	return r
}

/* ----------------------------------------------------------------
 *						M A I N | E X A M P L E
 *-----------------------------------------------------------------*/

/*
func DemoStringsFlag() {
	var mySelect StringsFlag
	mySelect.Strict(true)
	flag.Var(&mySelect, "alpha", "any of custom|english|spanish")
	flag.Parse()

	fmt.Printf("Selected value: %c\n", mySelect.Value)
}
*/
