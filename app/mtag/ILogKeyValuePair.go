/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Lord of Scripts
 *							   goApp
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * MLog variadic tags: String, Rune, Int, Bool, YesNo, Byte & At.
 * Tags like the log/slog package but enhanced.
 *-----------------------------------------------------------------*/
package mtag

import "fmt"

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

// ILogKeyValuePair defines the interface for mlog tags that can
// appear as variadic parameters to the mlog logging functions.
type ILogKeyValuePair interface {
	fmt.Stringer
}

/* ----------------------------------------------------------------
 *						M A I N | E X A M P L E
 *-----------------------------------------------------------------*/
/*

// For this specific situation, you might as well import it in
// the current package as it contains no other functionality
//import _ "github.com/lordofscripts/goapp/app/mtag"
import "github.com/lordofscripts/goapp/app/mtag"

func DemoMLog() {
	mlog.SetLevel(mlog.LevelDebug)
	mlog.Info("Useful information")
	mlog.Infof("For your info %c", '⥖')
	mlog.Error("Error happened")
	mlog.Errorf("%s with %d", message, value)
	err := fmt.Errorf("random error")
	mlog.ErrorE(err)
	mlog.Fatal(-5, "Terrible thing happened")
	mlog.DebugT("lazy programmer", mlog.String("Key","value"),
					mtag.Int("Key", 5),
					mtag.Rune("Rune", 'x'),
					mtag.Bool("Key", true),
					mtag.YesNo("Key", false),
					mtag.Err(err),
					mtag.At())

mlog.SetLevel(mlog.LevelDebug)
	mlog.Info("DidimusCommand R/T")
	mlog.Infof("This is %c", '⥖')
	mlog.InfoT("Tagged error",
		mtag.Rune("Rune", 'E'),
		mtag.YesNo("Bad", true),
		mtag.Int("Value", 5),
		mtag.String("String", "text here"))
}
*/
