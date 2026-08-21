/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2024 Dídimo Grimaldo T.
 *                           GoApp
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package util

import (
	"fmt"
	"strings"
)

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

/*
 * NOTES ON FORMATTING DECIMALS & THOUSANDS
 *	Applies to:
 *		- FormatMoney, FormatMoneyEUR,
 *		- FormatThousands, FormatThousandsEUR
 *	Normal convention ($18,584.56):
 *		- "," for thousands separator
 *		- "." for decimal separator
 *		- Used in most of the world.
 *	European convention (€18.584,56):
 *		- "." for thousands separator
 *		- "," for decimal separator
 *		- Europe: Germany, France, Italy, Spain, Netherlands, Belgium,
 *				Denmark, Norway, Sweden, Portugal, Austria, Romania,
 *				Turkey and Russia.
 *		- South America: Brazil, Argentina, Chile, Colombia(?), Peru, Uruguay & Ecuador.
 *		- Asia: Indonesia & Vietnam
 * NOTE:
 *	- In some regions the currency symbol is place after the number
 *	- In some countries the currency symbol has multiple characters (not a rune)
 */

// Formats a float value with the specified precision rendering it
// in American format (thousands with "," decimals with ".")
func FormatMoney(val float64, currencySymbol rune) string {
	return string(currencySymbol) + formatWithThousands(val, 2)
}

// Formats a float value with the specified precision rendering it
// in European format (thousands with "." decimals with ",")
func FormatMoneyEUR(val float64, currencySymbol rune) string {
	return string(currencySymbol) + formatWithThousands(val, 2, '.')
}

// Formats a float value with the specified precision rendering it
// in Swiss & Liechtenstein format (thousands with "'" decimals with ".")
func FormatMoneySwiss(val float64, currencySymbol rune) string {
	return string(currencySymbol) + formatWithThousands(val, 2, '\'')
}

// Formats a float value with the specified precision rendering it
// in American format (thousands with "," decimals with ".")
func FormatThousands(val float64, precision int) string {
	return formatWithThousands(val, precision)
}

// Formats a float value with the specified precision rendering it
// in European format (thousands with "." decimals with ",")
func FormatThousandsEUR(val float64, precision int) string {
	return formatWithThousands(val, precision, '.')
}

// Formats a float value with the specified precision rendering it
// in European format (thousands with "." decimals with ",")
func FormatThousandsSwiss(val float64, precision int) string {
	return formatWithThousands(val, precision, '\'')
}

// Format float in American format without thousands separator
func FormatFloat[T float32 | float64](val T, precision int) string {
	str := fmt.Sprintf("%.*f", precision, val)
	return str
}

// Format float in European format without thousands separator
func FormatFloatEUR[T float32 | float64](val T, precision int) string {
	str := fmt.Sprintf("%.*f", precision, val)
	str = strings.Replace(str, ".", ",", 1)
	return str
}

// Formats a currency (float64) with the desired precision. By default, if the last
// parameter thousandsSep is omitted, it uses the American format "18,250.23" but if
// the thousandsSep is '.' it uses European format "18.250,23". If the thousandSep
// is anything else than ',' or '.' it will use the Swiss & Liechtenstein exception
// which uses "'" as thousands separator and "." decimal separator, i.e. 18'250.23
// Usage:
//   - formatWithThousands(18250.23, 6) yields 18,250.23
//   - formatWithThousands(18250.23, 6, ',') yields 18,250.23
//   - formatWithThousands(18250.23, 6, '.') yields 18.250,23
func formatWithThousands(val float64, precision int, thousandsSep ...rune) string {
	// 1. Determine thousands and decimals separators
	var sepThousands string = ","
	var sepDecimals string = "."
	if thousandsSep != nil || len(thousandsSep) > 0 {
		switch thousandsSep[0] {
		case '.':
			// use weird European format
			sepThousands = "."
			sepDecimals = ","
		case '\'':
			// use even weirder Swiss/Liechtenstein format
			sepThousands = "'"
			sepDecimals = "."
		}
	}
	// 2. Format to a standard string with desired precision
	str := fmt.Sprintf("%.*f", precision, val)
	// 3. Separate whole numbers from decimals
	const NATIVE_DECIMAL_SEP string = "."
	const NEGATIVE_SIGN string = "-"
	parts := strings.Split(str, NATIVE_DECIMAL_SEP)
	whole := parts[0]
	// 4. Handle negative signs
	isNegative := strings.HasPrefix(whole, NEGATIVE_SIGN)
	if isNegative {
		whole = whole[1:]
	}
	// 5. Loop backward and inject separators
	var result []string
	length := len(whole)
	for i := length; i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		result = append([]string{whole[start:i]}, result...)
	}

	// 6. Reconstruct the final string: 1000's
	finalWhole := strings.Join(result, sepThousands)
	if isNegative {
		finalWhole = NEGATIVE_SIGN + finalWhole
	}

	if len(parts) > 1 {
		return finalWhole + sepDecimals + parts[1]
	}

	return finalWhole
}

/* ----------------------------------------------------------------
 *                          T E S T S
 *-----------------------------------------------------------------*/
