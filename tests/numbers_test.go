/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           GoApp
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *						U n i t   T e s t
 *-----------------------------------------------------------------*/
package tests

import (
	"testing"

	"github.com/lordofscripts/goapp/util"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				U n i t  T e s t   F u n c t i o n s
 *						Operate on Current
 *-----------------------------------------------------------------*/

// Target: util.FormatThousands
// Formatting floats in American format (thousands:"," decimals:".")
func Test_FormatThousands_American(t *testing.T) {
	testCases := []struct {
		money     float64
		precision int
		expected  string
	}{
		{150.385, 3, "150.385"},
		{28563.485, 2, "28,563.49"},
		{4828563.485, 2, "4,828,563.49"},
		{4828563.484, 2, "4,828,563.48"},
	}

	for i, tCase := range testCases {
		got := util.FormatThousands(tCase.money, tCase.precision)
		if got != tCase.expected {
			t.Fatalf("#%d American failed Exp: %s Got: %s", i+1, tCase.expected, got)
		}
	}
}

// Target: util.FormatThousandsEUR
// Formatting floats in European format (thousands:"." decimals:",")
func Test_FormatThousands_European(t *testing.T) {
	testCases := []struct {
		money     float64
		precision int
		expected  string
	}{
		{150.385, 3, "150,385"},
		{28563.485, 2, "28.563,49"},
		{4828563.485, 2, "4.828.563,49"},
		{4828563.484, 2, "4.828.563,48"},
	}

	for i, tCase := range testCases {
		got := util.FormatThousandsEUR(tCase.money, tCase.precision)
		if got != tCase.expected {
			t.Fatalf("#%d European failed Exp: %s Got: %s", i+1, tCase.expected, got)
		}
	}
}

// Target: util.FormatThousandsSwiss
// Formatting floats in Swiss/Liechtenstein format (thousands:"'" decimals:".")
func Test_FormatThousands_Swiss(t *testing.T) {
	testCases := []struct {
		money     float64
		precision int
		expected  string
	}{
		{150.385, 3, "150.385"},
		{28563.485, 2, "28'563.49"},
		{4828563.485, 2, "4'828'563.49"},
		{4828563.484, 2, "4'828'563.48"},
	}

	for i, tCase := range testCases {
		got := util.FormatThousandsSwiss(tCase.money, tCase.precision)
		if got != tCase.expected {
			t.Fatalf("#%d European failed Exp: %s Got: %s", i+1, tCase.expected, got)
		}
	}
}

/* ----------------------------------------------------------------
 *					H e l p e r   F u n c t i o n s
 *-----------------------------------------------------------------*/
