// Examples of using packages in GoApp module.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lordofscripts/goapp/app"
	"github.com/lordofscripts/goapp/flagx"
)

const (
	ORG_NAME       string      = "coralys"
	APP_NAME       string      = "goapptest"
	META_FILE_MODE os.FileMode = 0644
	MANUAL_VERSION             = "1.0.0"

	CHR_TRIDENT = rune(0x1f531) // 🔱
)

const (
	// you can create your own WarningCode enumeration
	WarnNotFound app.WarningCode = iota
	WarnInvalid
)

var (
	Version app.PackageVersion = app.NewReleaseVersion(APP_NAME, "Just a test", MANUAL_VERSION)
)

func main() {
	// ############ flagx Package ##############
	var mySelect flagx.StringsFlag
	var myByte flagx.ByteFlag
	var myRune flagx.RuneFlag
	var myDate flagx.DateFlag = *flagx.NewDateVar("2006-Jan-02")

	mySelect.Strict(true)
	mySelect.SetChoices([]string{"custom", "english", "spanish"})

	flag.Var(&mySelect, "alpha", "any of custom|english|spanish")
	flag.Var(&myByte, "byte", "custom Byte value")
	flag.Var(&myRune, "rune", "custom Rune value")
	flag.Var(&myDate, "date", "custom Date value")
	flag.Parse()

	// ########### app Package ##############
	Version.Copyright("lordofscripts", CHR_TRIDENT)

	// ############ flagx Package ##############
	fmt.Printf("Selected value: %s\n", mySelect.Value)
	fmt.Printf("Byte value: %c\n", myByte.Value)
	fmt.Printf("Rune value: %c\n", myRune.Value)
	fmt.Printf("Date value: %s\n", myDate.Value)

	// ########### app Warning ##############
	warn1 := app.NewWarning(WarnNotFound, "could not find %s", os.Args[0])
	warn2 := app.NewWarningAsErr(WarnInvalid, "something invalid %s", os.Args[0])
	fmt.Println("Plain warning: ", warn1)
	if warnT, isWarn := app.IsWarning(warn2); isWarn {
		fmt.Println("Disguised warning: ", warnT)
	} else {
		fmt.Printf("Warning as error: %v\n", warn2)
	}

	// ########### app Package ##############
	Version.BuyMeCoffee("lostinwriting")
}
