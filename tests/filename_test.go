/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2024 Dídimo Grimaldo T.
 *                           GoApp
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *						U n i t   T e s t
 *-----------------------------------------------------------------*/
package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/lordofscripts/goapp/util"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	TBASE string = "Filename{} "
	EXT   string = ".ini"

	BAD_PATH1     string = "/usr/share/bin/matusalen.cfg"
	BAD_DIR1      string = "/usr/share/bin"
	BAD_FILENAME1 string = "matusalen.cfg"
	BAD_BASE1     string = "matusalen"
	BAD_EXT1      string = ".cfg"

	FILE_PRESENT  string = "/etc/passwd"
	FILE_PRESENT1 string = "testdata/fifty.bin"
	FILE_PRESENT2 string = "testdata/fifty_rnd.bin"
	FILE_ABSENT   string = "/tmp/dummydummy.123"
	DIR_PRESENT   string = "/etc"
)

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

func TestFilename_Ctor(t *testing.T) {
	fno := util.NewFilename(BAD_PATH1)
	if fno.GetFullPath() != BAD_PATH1 {
		t.Errorf("Ctor expected %q", BAD_PATH1)
	}
}

// the object's stringer should return the name given on the constructor
// with any manipulations it may have suffered outside a non-commited query.
func TestFilename_String(t *testing.T) {
	fno := util.NewFilename(BAD_PATH1)
	fno.WithFile(FILE_ABSENT) // query operates on this rather than PATH1
	if fno.String() != BAD_PATH1 {
		t.Fatal("fmt.Stringer should have returned original path")
	}
}

func TestFilename_Decompose(t *testing.T) {
	fno := util.NewFilename(BAD_PATH1)
	dir, base, ext := fno.Decompose()

	Show("Input:", BAD_PATH1)
	Show("Dir  :", BAD_DIR1)
	Show("Base :", BAD_BASE1)
	Show("Ext  :", BAD_EXT1)

	if dir != BAD_DIR1 {
		t.Errorf("Dir exp %q got %q", BAD_DIR1, dir)
	}
	if base != BAD_BASE1 {
		t.Errorf("Base exp %q got %q", BAD_BASE1, base)
	}
	if ext != BAD_EXT1 {
		t.Errorf("Ext exp %q got %q", BAD_EXT1, ext)
	}
}

func TestFilename_GetName(t *testing.T) {
	Title(TBASE + "Name")
	fno := util.NewFilename(BAD_PATH1)
	exp := BAD_FILENAME1
	got := fno.GetName()

	Show("Input:", BAD_PATH1)
	Show("Name :", got)

	if exp != got {
		t.Errorf("Name() exp %q got %q", exp, got)
	}
}

func TestFilename_ReplaceExt(t *testing.T) {
	Title(TBASE + "ReplaceExt")
	fno := util.NewFilename(BAD_PATH1)
	name := fno.ReplaceExt(EXT)
	Show("Input : ", BAD_PATH1)
	Show("Output: ", name)

	if name != "/usr/share/bin/matusalen.ini" {
		t.Errorf("Got %s", name)
	}
}

func TestFilename_RemoveExt(t *testing.T) {
	Title(TBASE + "RemoveExt")
	fno := util.NewFilename(BAD_PATH1)
	name := fno.RemoveExt()

	Show("Input", BAD_PATH1)
	Show("Output", name)

	if name != "/usr/share/bin/matusalen" {
		t.Errorf("Got %s", name)
	}
}

func TestFilename_ReplaceDirectory(t *testing.T) {
	fno := util.NewFilename("/tmp/var/dummy.log")
	got := fno.ReplaceDirectory("/usr/share")
	if got != "/usr/share/dummy.log" {
		t.Fatal("ReplaceDirectory failed")
	}
}

func TestFilename_IsNotThere(t *testing.T) {
	Title(TBASE + "NotThere")
	fno := util.NewFilename(BAD_PATH1)
	if notThere, err := fno.IsNotThere(); notThere != true {
		t.Errorf("expected true on IsNotThere: %v", err)
	}
}

func TestFilename_IsThere(t *testing.T) {
	Title(TBASE + "IsThere")

	fno := util.NewFilename(SureFile())
	if isThere, err := fno.IsThere(); isThere != true {
		t.Fatalf("expected true on IsThere(%s): %v", fno.String(), err)
	}

	fno = util.NewFilename(FILE_ABSENT)
	if isThere, err := fno.IsThere(); isThere == true {
		t.Fatalf("%s isn't supposed to exist: %v", FILE_ABSENT, err)
	}
}

func TestFilename_IsInPath(t *testing.T) {
	Title(TBASE + "IsInPath")

	fno := util.NewFilename(BAD_PATH1)
	if fno.IsInPath(BAD_DIR1) != true {
		t.Error("expected true on IsInPath")
	}
}

func TestFilename_IsDirectory(t *testing.T) {
	Title(TBASE + "IsDirectory")
	fno := util.NewFilename(SureDirectory())
	if ok, err := fno.IsDirectory(); err != nil {
		t.Error(err)
	} else if !ok {
		t.Errorf("Says %q is not directory", DIR_PRESENT)
	}
}

// @todo IsSizeGreater() IsSizeSmaller()

/* ----------------------------------------------------------------
 *				U n i t  T e s t   F u n c t i o n s
 *				  Operate on IFluentFilename query
 *-----------------------------------------------------------------*/

func TestFilename_Query_With(t *testing.T) {
	fno := util.NewFilename(BAD_PATH1)
	q := fno.With()

	// an attempt to call With*() without terminating previous returns NIL
	if fno.With() != nil {
		t.Fatal("expected nil with ongoing query")
	}

	if q.QueryStrResult(false) != BAD_PATH1 {
		t.Fatal("failed With() with valid path")
	}
}

func TestFilename_Query_WithFile(t *testing.T) {
	fno := util.NewFilename(BAD_PATH1)
	q := fno.WithFile(FILE_ABSENT) // query operates on this rather than PATH1

	// (a) when a query is active, any With*() call returns nil
	if fno.WithFile(FILE_ABSENT) != nil {
		t.Fatal("expected nil with ongoing query")
	}

	// (b) obtain result of query
	sV, bV := q.QueryResult(false)
	if sV != FILE_ABSENT {
		t.Fatal("after query init, the result should be the original query")
	}
	if bV != true {
		t.Fatal("after query init, expected true for current bool query result")
	}
}

/* ----------------------------------------------------------------
 *				U n i t  T e s t   F u n c t i o n s
 *				  Operate on IFluentFile query
 *-----------------------------------------------------------------*/

func TestFilename_Query_HasExt(t *testing.T) {
	Title(TBASE + "Q.HasExt")
	Extensions := []string{".txt", ".123"}
	fno := util.NewFilename(FILE_ABSENT)
	if fno.HasExt(Extensions) != nil {
		t.Fatal("call to HasExt without ongoing query should be nil")
	}

	if fno.With().HasExt(Extensions).QueryBoolResult() != true {
		t.Fatal("call to HasExt without ongoing query should be nil")
	}
}

func TestFilename_Query_NotThere(t *testing.T) {
	Title(TBASE + "Q.NotThere")
	fno := util.NewFilename(FILE_PRESENT1).WithFile(FILE_ABSENT)
	if fno.NotThere().QueryBoolResult() != true {
		t.Fatal("expected true on Q.NotThere")
	}

	fno = util.NewFilename(FILE_ABSENT).WithFile(FILE_PRESENT1)
	if fno.NotThere().QueryBoolResult() != false {
		t.Fatal("expected false on Q.NotThere")
	}
}

func TestFilename_Query_Exists(t *testing.T) {
	Title(TBASE + "Q.Exists")
	fno := util.NewFilename(FILE_PRESENT1).WithFile(FILE_ABSENT)
	if fno.Exists().QueryBoolResult() != false {
		t.Fatal("expected false on Q.Exists")
	}

	fno = util.NewFilename(FILE_ABSENT).WithFile(FILE_PRESENT1)
	if fno.Exists().QueryBoolResult() != true {
		t.Fatal("expected true on Q.Exists")
	}
}

func TestFilename_Query_ChangeDirectory(t *testing.T) {
	q := util.NewFilename("/tmp/var/dummy.log").
		With().
		ChangeExt("txt").
		ChangeDirectory("/usr/share")
	got := q.QueryStrResult(false)
	if got != "/usr/share/dummy.txt" {
		t.Fatal("ChangeDirectory failed")
	}
}

func TestFilename_Query_InPath(t *testing.T) {
	Title(TBASE + "Q.InPath")

	fno := util.NewFilename(BAD_PATH1).With()
	if fno.InPath(BAD_DIR1).QueryBoolResult() != true {
		t.Error("expected true on Q.InPath")
	}
}

func TestFilename_Query_Directory(t *testing.T) {
	Title(TBASE + "Q.Directory")
	directory := SureDirectory()
	fno := util.NewFilename(directory).With()
	if fno.Directory().QueryBoolResult() != true {
		t.Errorf("Says %q is not directory", directory)
	}
}

func TestFilename_Query_Greater(t *testing.T) {
	Title(TBASE + "Q.Greater")
	const FSIZE = 50 // actual file size
	fno := util.NewFilename(FILE_PRESENT1).With()
	if fno.Greater(FSIZE).QueryBoolResult() != false {
		t.Fatalf("got true on Q.Greater(50) for %s", FILE_PRESENT1)
	}

	// after QueryBoolResult we need to initiate a new one
	fno = util.NewFilename(FILE_PRESENT1).With()
	if fno.Greater(FSIZE-10).QueryBoolResult() != true {
		t.Fatalf("got false on Q.Greater(50) for %s", FILE_PRESENT1)
	}
}

func TestFilename_Query_Smaller(t *testing.T) {
	Title(TBASE + "Q.Greater")
	const FSIZE = 50 // actual file size
	fno := util.NewFilename(FILE_PRESENT1).With()
	if fno.Smaller(FSIZE).QueryBoolResult() != false {
		t.Fatalf("got true on Q.Smaller(50) for %s", FILE_PRESENT1)
	}

	// after QueryBoolResult we need to initiate a new one
	fno = util.NewFilename(FILE_PRESENT1).With()
	if fno.Smaller(FSIZE+10).QueryBoolResult() != true {
		t.Fatalf("got false on Q.Smaller(50) for %s", FILE_PRESENT1)
	}
}

/* ----------------------------------------------------------------
 *					H e l p e r   F u n c t i o n s
 *-----------------------------------------------------------------*/

// Test setup/teardown
func TestMain(m *testing.M) {
	// setup
	fmt.Println("*** Filename{} tests ***")

	// run the tests
	result := m.Run()

	// teardown
	//

	os.Exit(result)
}
