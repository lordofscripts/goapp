/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Lord of Scripts
 *							goApp::mtag
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Key-Value pair tags used in MLog and ZLog but also useful in your
 * everyday application.
 * Implements: Money tag (uses commas)
 *-----------------------------------------------------------------*/
package mtag

import (
	"fmt"

	"github.com/lordofscripts/goapp/util"
)

/* ----------------------------------------------------------------
 *						L o c a l s
 *-----------------------------------------------------------------*/

const (
	regionNone regionSpec = iota // render in compact form without thousands
	regionAmerican
	regionEuropean
	regionSwiss
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ ILogKeyValuePair = (*kvMoney)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type regionSpec byte

// A float-type key-value pair tag
type kvMoney struct {
	label          string
	amount         float64
	currencySymbol rune
	useSeparators  bool
	region         regionSpec
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// implements fmt.Stringer for mtag.Float() rendering the value with
// the number of whole digits and decimal digits specified in the
// Float() function call.
func (k *kvMoney) String() string {
	var str string
	if k.useSeparators {
		switch k.region {
		case regionAmerican:
			str = util.FormatMoney(k.amount, k.currencySymbol)

		case regionEuropean:
			str = util.FormatMoneyEUR(k.amount, k.currencySymbol)

		case regionSwiss:
			str = util.FormatMoneySwiss(k.amount, k.currencySymbol)
		}
	} else {
		switch k.region {
		case regionSwiss:
			fallthrough // same decimal separator, no thousands separator
		case regionAmerican:
			str = string(k.currencySymbol) + util.FormatFloat(k.amount, 2)

		case regionEuropean:
			str = string(k.currencySymbol) + util.FormatFloatEUR(k.amount, int(k.currencySymbol))
		}
	}
	return fmt.Sprintf("%s=%s", k.label, str)
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// log the currency amount with the nr of whole digits and 2 decimals
// in compact fashion (no thousands separator). American format.
// Example: MoneyF("Balance", '$', 4531.58) yields $4531.58
func Money(key string, currencySymbol rune, value float64) ILogKeyValuePair {
	return &kvMoney{label: key, amount: value, currencySymbol: currencySymbol, useSeparators: false, region: regionAmerican}
}

// log the currency amount with the European format: €4.850,23
// Example: MoneyEUR("Balance", '€', 4531.58) yields €4.531,58
func MoneyEUR(key string, currencySymbol rune, value float64) ILogKeyValuePair {
	return &kvMoney{label: key, amount: value, currencySymbol: currencySymbol, useSeparators: true, region: regionEuropean}
}

// log the currency amount with the Swiss & Liechtenstein format: €4'850,23
// Example: MoneyEUR("Balance", '€', 4531.58) yields €4.531,58
func MoneyCHE(key string, currencySymbol rune, value float64) ILogKeyValuePair {
	return &kvMoney{label: key, amount: value, currencySymbol: currencySymbol, useSeparators: true, region: regionSwiss}
}
