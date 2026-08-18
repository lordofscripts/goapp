/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package app

import (
	"fmt"
)

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/

const (
	// don't change
	DevStatusAlpha    DevStatus = "Alpha"
	DevStatusBeta     DevStatus = "Beta"
	DevStatusRC       DevStatus = "RC" // Release Candidate
	DevStatusReleased DevStatus = ""
)

/* ----------------------------------------------------------------
 *                     I N T E R F A C E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                         T Y P E S
 *-----------------------------------------------------------------*/

// Development code status: Alpha,Beta,RC,Released
type DevStatus = string

// Package/Module/Application version descriptor
type PackageVersion struct {
	n      string    // name
	v      string    // version tag
	s      DevStatus // status
	sv     int       // Alpha/Beta/RC-### sequence
	dsc    string    // description
	author string    // module's author (optional)
}

/* ----------------------------------------------------------------
 *                   P U B L I C    T Y P E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                   P R I V A T E    T Y P E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                    C O N S T R U C T O R S
 *-----------------------------------------------------------------*/

func NewPackageVersion(name, description string, verStr string, status DevStatus) PackageVersion {
	return PackageVersion{
		n:   name,
		v:   verStr,
		s:   status,
		sv:  1,
		dsc: description,
	}
}

// ctor for Release versions of a package
func NewReleaseVersion(name, description string, verStr string) PackageVersion {
	return PackageVersion{
		n:   name,
		v:   verStr,
		s:   DevStatusReleased,
		sv:  1,
		dsc: description,
	}
}

// ctor for Alpha versions of a package
func NewAlphaVersion(name, description string, verStr string, alphaNum int) PackageVersion {
	return PackageVersion{
		n:   name,
		v:   verStr,
		s:   DevStatusAlpha,
		sv:  alphaNum,
		dsc: description,
	}
}

// ctor for Beta versions of a package
func NewBetaVersion(name, description string, verStr string, betaNum int) PackageVersion {
	return PackageVersion{
		n:   name,
		v:   verStr,
		s:   DevStatusBeta,
		sv:  betaNum,
		dsc: description,
	}
}

// ctor for Release Candidate versions of a package, i.e. 1.0.0-RC.1
func NewReleaseCandidateVersion(name, description string, verStr string, rcNum int) PackageVersion {
	return PackageVersion{
		n:   name,
		v:   verStr,
		s:   DevStatusRC,
		sv:  rcNum,
		dsc: description,
	}
}

/* ----------------------------------------------------------------
 *                        M E T H O D S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                 P U B L I C    M E T H O D S
 *-----------------------------------------------------------------*/

// Specify a revision number in case of DevStatusRC, DevStatusAlpha
// or DevStatusBeta. Ignored in the output for DevStatusRelease.
func (pv *PackageVersion) WithRevision(rev int) PackageVersion {
	pv.sv = rev
	return *pv
}

// (Optional) register the module's author for Support information
func (pv *PackageVersion) WithAuthor(author string) PackageVersion {
	pv.author = author
	return *pv
}

func (pv PackageVersion) GetAuthor() string {
	return pv.author
}

// Like String() but without the name, just the version number
// without the "v" prefix, i.e. "1.5.4".
func (pv PackageVersion) Version() string {
	var ver string

	switch pv.s {
	case DevStatusAlpha:
		fallthrough
	case DevStatusBeta:
		fallthrough
	case DevStatusRC:
		ver = fmt.Sprintf("%s-%s.%d", pv.v, pv.s, pv.sv)
	default:
		ver = fmt.Sprintf("%s", pv.v)
	}

	return ver
}

// Like String() but without the name, just the version number
// prefixed with "v". Or like Version() with "v" prefix.
func (pv PackageVersion) Short() string {
	return "v" + pv.Version()
}

// implements fmt.Stringer and returns the name of the package/module/app
// followed by its version number.
func (pv PackageVersion) String() string {
	return fmt.Sprintf("%s %s", pv.n, pv.Short())
}

// Same as String() but with the description. Suitable for an application's
// Support information
func (pv PackageVersion) Info() string {
	return fmt.Sprintf("%s : %s", pv.String(), pv.dsc)
}

// Example: Version.Copyright
func (pv PackageVersion) Copyright(owner string, adornment rune) {
	fmt.Printf("\t%c %s %s %c\n", adornment, pv.String(), owner, adornment)
}

// Hey! My time costs money too!
func (pv PackageVersion) BuyMeCoffee(coffee4 string) {
	const COFFEE_CUP rune = '\u2615' // ☕

	fmt.Printf("\t%c Buy me a Coffee? https://www.buymeacoffee/%s\n", COFFEE_CUP, coffee4)
}

/* ----------------------------------------------------------------
 *                 P R I V A T E    M E T H O D S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/
