/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2024 Dídimo Grimaldo T.
 *                           GoApp
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Filename manipulations (v2) with and without a fluent API query.
 * NOTE: This object was promoted and improved from my Secretum module.
 *-----------------------------------------------------------------*/
package util

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/

const (
	_INITIAL_BOOL_QUERY_RESULT bool = true
)

/* ----------------------------------------------------------------
 *                     I N T E R F A C E S
 *-----------------------------------------------------------------*/

// Methods to initiate, chain and obtain results of (multiple)  queries
// on a filename object.
type IFilenameQuery interface {
	// (A) Fluent API filename query initiators
	// ------------------------------------------
	// initiate the query with a different filename than the one
	// used in the constructor.
	WithFile(filename string) IFilenameQuery
	// initiate the query with the currently known filename
	With() IFilenameQuery

	// (B) Fluent API filename query terminators
	// ------------------------------------------
	// Current result of the query with option to commit its value
	// to the object.
	QueryStrResult(commit bool) string
	// Current result of bool query
	QueryBoolResult() bool
	// Returns the outcome of both QueryStrResult and QueryBoolResult
	// It is useful because a call to any of those two results in the
	// entire query (bool & string) to be reset.
	QueryResult(commit bool) (string, bool)

	// (C) Fluent API file queries whose intermediate result is a
	// YES/NO (bool) to obtain with QueryBoolResult()
	HasExt(extensions []string) IFilenameQuery
	NotThere() IFilenameQuery
	Exists() IFilenameQuery
	InPath(path string) IFilenameQuery
	Directory() IFilenameQuery
	ChangeDirectory(newDir string) IFilenameQuery
	ChangeExt(newExt string) IFilenameQuery
	Greater(size int64) IFilenameQuery
	Smaller(size int64) IFilenameQuery
}

var _ IFilenameQuery = (*Filename)(nil)

/* ----------------------------------------------------------------
 *                   P U B L I C    T Y P E S
 *-----------------------------------------------------------------*/

// An object that allows several chained manipulations on a filename
// or certain attributes of a file/directory.
type Filename struct {
	fqn             string     // general non-fluent filename operations
	queryResultStr  string     // working value for fluent API
	queryResultBool bool       // fluent API boolean queries
	onQuery         bool       // a query is in progress. Cancelled with QueryResult*()
	qmux            sync.Mutex // query Mutex for IFilenameQuery
	nmux            sync.Mutex // normal Mutex
}

/* ----------------------------------------------------------------
 *                    C O N S T R U C T O R S
 *-----------------------------------------------------------------*/

// (Ctor) Filename operations
func NewFilename(fullyQualified string) *Filename {
	return &Filename{
		fqn:             fullyQualified,
		queryResultStr:  "",
		queryResultBool: false,
		onQuery:         false,
	}
}

/* ----------------------------------------------------------------
 *                 P U B L I C    M E T H O D S
 *-----------------------------------------------------------------*/

// implements fmt.Stringer returning the current fully-qualified
// filename.
func (f *Filename) String() string {
	return f.fqn
}

// The fully-qualified filename with directory (if available).
// NOTE: does not require a query.
func (f *Filename) GetFullPath() string {
	return f.fqn
}

// The filename without its directory but including file extension.
// NOTE: does not require a query.
func (f *Filename) GetName() string {
	_, base, ext := f.Decompose()
	return base + ext
}

// Decompose a filename/path into its parts:
// Filename: /usr/bin/secretum.exe
// Dir: /usr/bin
// Base: secretum
// Ext: .exe
// NOTE: does not require a query.
func (f *Filename) Decompose() (dir, base, ext string) {
	f.nmux.Lock()
	defer f.nmux.Unlock()

	return f.decompose(f.fqn)
}

// Replaces the extension by 'ext'
// NOTE: does not require a query.
func (f *Filename) ReplaceExt(ext string) string {
	f.nmux.Lock()
	defer f.nmux.Unlock()

	return f.replaceExt(f.fqn, ext)
}

// Replaces the directory by 'newDir'
// NOTE: does not require a query.
func (f *Filename) ReplaceDirectory(newDir string) string {
	f.nmux.Lock()
	defer f.nmux.Unlock()

	return f.replaceDirectory(f.fqn, newDir)
}

// The fully-qualified filename (path & basename) without the extension
// NOTE: does not require a query.
func (f *Filename) RemoveExt() string {
	f.nmux.Lock()
	defer f.nmux.Unlock()

	dir, base, _ := f.decompose(f.fqn)
	return filepath.Join(dir, base)
}

// The file is NOT there, but if there is an error of some kind
// that is not ErrNotExist, then the false result is inconclusive!
// NOTE: does not require a query.
func (f *Filename) IsNotThere() (bool, error) {
	f.nmux.Lock()
	defer f.nmux.Unlock()

	return f.isNotThere(f.fqn)
}

// Whether the file exists. If it returns true the file exists and
// is accessible. If it returns false without error then it does not
// exist, but false with error means it may or may not exist due
// to possible file permission errors.
// NOTE: does not require a query.
func (f *Filename) IsThere() (bool, error) {
	f.nmux.Lock()
	defer f.nmux.Unlock()

	return f.isThere(f.fqn)
}

// The file lives within that directory
// NOTE: does not require a query.
func (f *Filename) IsInPath(path string) bool {
	f.nmux.Lock()
	defer f.nmux.Unlock()

	return strings.HasPrefix(f.fqn, path)
}

// It is a directory
// NOTE: does not require a query.
func (f *Filename) IsDirectory() (bool, error) {
	f.nmux.Lock()
	defer f.nmux.Unlock()

	return f.isDirectory(f.fqn)
}

// the file's size is greater than
// NOTE: does not require a query.
func (f *Filename) IsSizeGreater(size int64) (bool, error) {
	f.nmux.Lock()
	defer f.nmux.Unlock()

	return f.isSizeGreater(f.fqn, size)
}

// the file's size is smaller than
// NOTE: does not require a query.
func (f *Filename) IsSizeSmaller(size int64) (bool, error) {
	f.nmux.Lock()
	defer f.nmux.Unlock()

	return f.isSizeSmaller(f.fqn, size)
}

/* ----------------------------------------------------------------
 *             P U B L I C   A P I   M E T H O D S
 *-----------------------------------------------------------------*/

// Begins a new (string/bool-based) Fluent API filename manipulation query.
// It must be ended with a call to QueryResult(commit).
func (f *Filename) WithFile(filename string) IFilenameQuery {
	f.qmux.Lock()
	defer f.qmux.Unlock()

	if f.queryInProgress() {
		fmt.Fprintf(os.Stderr, "a query is already in progress (forgot to call QueryResult?), ignoring request.")
		return nil
	} else {
		f.queryResultStr = filename
		f.queryResultBool = _INITIAL_BOOL_QUERY_RESULT // all Bool queries are AND and thus begin with true to work.
		f.onQuery = true
		return f
	}
}

// Begins a new (string-based) Fluent API filename manipulation query
// using the currently known filename stored in the object by a commited query result.
// It must be ended with a call to QueryResult(commit).
func (f *Filename) With() IFilenameQuery {
	return f.WithFile(f.fqn)
}

// Ends a Fluent API filename query and returns the result. If commitValue
// is true, it replaces the object's current value overriding the one given
// in the constructor or in any other commited query.
func (f *Filename) QueryStrResult(commitValue bool) string {
	f.qmux.Lock()
	defer f.qmux.Unlock()

	var result = ""
	if !f.queryInProgress() {
		fmt.Fprintf(os.Stderr, "there is no IFilename query in progress")
	} else {
		result = f.queryResultStr
	}

	if commitValue {
		f.fqn = f.queryResultStr
	}

	f.queryResultStr = ""
	f.queryResultBool = false
	f.onQuery = false
	return result
}

// Result of an IFluentFile query
func (f *Filename) QueryBoolResult() bool {
	f.qmux.Lock()
	defer f.qmux.Unlock()

	if f.queryBoolInProgress() {
		result := f.queryResultBool
		f.queryResultBool = false
		f.onQuery = false
		return result
	} else {
		fmt.Fprintf(os.Stderr, "invalid filename bool query result, none in progress!")
		os.Exit(127)
	}

	return false
}

// Returns the outcome of both QueryStrResult and QueryBoolResult
// It is useful because a call to any of those two results in the
// entire query (bool & string) to be reset.
func (f *Filename) QueryResult(commit bool) (s string, b bool) {
	f.qmux.Lock()
	defer f.qmux.Unlock()

	s = f.queryResultStr
	b = f.queryResultBool
	f.onQuery = false
	return
}

// Checks that the filename in the current query value contains
// any of the file extensions in the list.
// NOTE: Requires a started query (See With() and WithFile()).
func (f *Filename) HasExt(extensions []string) IFilenameQuery {
	f.qmux.Lock()
	defer f.qmux.Unlock()

	if !f.queryBoolInProgress() {
		return nil
	}

	hit := f.hasExt(f.fqn, extensions)
	f.queryResultBool = f.queryResultBool && hit
	return f
}

// The file or dir is not there.
// NOTE: Requires a started query (See With() and WithFile()).
func (f *Filename) NotThere() IFilenameQuery {
	f.qmux.Lock()
	defer f.qmux.Unlock()

	if !f.queryBoolInProgress() {
		return nil
	}

	if hit, err := f.isNotThere(f.queryResultStr); err != nil {
		return nil
	} else {
		f.queryResultBool = f.queryResultBool && hit
	}

	return f
}

// The file or dir exists.
// NOTE: Requires a started query (See With() and WithFile()).
func (f *Filename) Exists() IFilenameQuery {
	f.qmux.Lock()
	defer f.qmux.Unlock()

	if !f.queryBoolInProgress() {
		return nil
	}

	if hit, err := f.isThere(f.queryResultStr); err != nil {
		return nil
	} else {
		f.queryResultBool = f.queryResultBool && hit
	}
	return f
}

// The file is on that path.
// NOTE: Requires a started query (See With() and WithFile()).
func (f *Filename) InPath(path string) IFilenameQuery {
	f.qmux.Lock()
	defer f.qmux.Unlock()

	if !f.queryBoolInProgress() {
		return nil
	}

	f.queryResultBool = f.queryResultBool && strings.HasPrefix(f.queryResultStr, path)
	return f
}

// Replaces the directory by 'newDir'
// NOTE: does REQUIRE an active query.
func (f *Filename) ChangeDirectory(newDir string) IFilenameQuery {
	f.qmux.Lock()
	defer f.qmux.Unlock()

	if !f.queryInProgress() {
		return nil
	}
	f.queryResultStr = f.replaceDirectory(f.queryResultStr, newDir)
	return f
}

func (f *Filename) ChangeExt(newExt string) IFilenameQuery {
	f.qmux.Lock()
	defer f.qmux.Unlock()

	if !f.queryInProgress() {
		return nil
	}

	f.queryResultStr = f.replaceExt(f.queryResultStr, newExt)
	return f
}

// it is a directory.
// NOTE: Requires a started query (See With() and WithFile()).
func (f *Filename) Directory() IFilenameQuery {
	f.qmux.Lock()
	defer f.qmux.Unlock()

	if !f.queryBoolInProgress() {
		return nil
	}

	if ok, err := f.isDirectory(f.queryResultStr); err != nil {
		return nil
	} else {
		f.queryResultBool = f.queryResultBool && ok
	}

	return f
}

// file is greater than.
// NOTE: Requires a started query (See With() and WithFile()).
func (f *Filename) Greater(size int64) IFilenameQuery {
	f.qmux.Lock()
	defer f.qmux.Unlock()

	if !f.queryBoolInProgress() {
		return nil
	}

	if ok, err := f.isSizeGreater(f.queryResultStr, size); err != nil {
		return nil
	} else {
		f.queryResultBool = f.queryResultBool && ok
	}

	return f
}

// file is smaller than.
// NOTE: Requires a started query (See With() and WithFile()).
func (f *Filename) Smaller(size int64) IFilenameQuery {
	f.qmux.Lock()
	defer f.qmux.Unlock()

	if !f.queryBoolInProgress() {
		return nil
	}

	if ok, err := f.isSizeSmaller(f.queryResultStr, size); err != nil {
		return nil
	} else {
		f.queryResultBool = f.queryResultBool && ok
	}

	return f
}

/* ----------------------------------------------------------------
 *                 P R I V A T E    M E T H O D S
 *-----------------------------------------------------------------*/

// that a string-based filename query is in progress
func (f *Filename) queryInProgress() bool {
	return f.onQuery
}

// that a boolean query is in progress, to avoid polluting the end
// result with an unknown previous value.
func (f *Filename) queryBoolInProgress() bool {
	return f.onQuery
}

// Decompose a filename/path into its parts:
// Filename: /usr/bin/secretum.exe
// Dir: /usr/bin
// Base: secretum
// Ext: .exe
func (f *Filename) decompose(name string) (dir, base, ext string) {
	dir = filepath.Dir(name)
	base = filepath.Base(name)
	ext = filepath.Ext(name)

	if len(ext) > 0 {
		base = base[:len(base)-len(ext)]
	}

	return dir, base, ext
}

// The file is NOT there, but if there is an error of some kind
// that is not ErrNotExist, then the false result is inconclusive!
func (f *Filename) isNotThere(name string) (bool, error) {
	if _, err := os.Stat(name); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return true, nil // categorically NOT there
		}

		return false, err // we may never know
	}

	return false, nil
}

// Whether the file exists. If it returns true the file exists and
// is accessible. If it returns false without error then it does not
// exist, but false with error means it may or may not exist due
// to possible file permission errors.
func (f *Filename) isThere(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err // with error, may or may not exist
}

// It is a directory
func (f *Filename) isDirectory(name string) (bool, error) {
	fi, err := os.Stat(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Filename error: %v", err)
		return false, err
	}

	return fi.IsDir(), nil
}

// Replaces the directory by 'newDir'
// NOTE: does not require a query.
func (f *Filename) replaceDirectory(name string, newDir string) string {
	_, base, ext := f.decompose(name)

	return filepath.Join(newDir, base+ext)
}

// Replaces the extension by 'ext'
// NOTE: does not require a query.
func (f *Filename) replaceExt(name, ext string) string {
	ext = strings.Trim(ext, " \t")
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	dir, base, _ := f.decompose(name)

	return filepath.Join(dir, base+ext)
}

// the file's size is greater than
func (f *Filename) isSizeGreater(name string, size int64) (bool, error) {
	fi, err := os.Stat(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "file access error: %v", err)
		return false, err
	}

	return fi.Size() > size, nil
}

// the file's size is smaller than
func (f *Filename) isSizeSmaller(name string, size int64) (bool, error) {
	fi, err := os.Stat(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "file access error: %v", err)
		return false, err
	}

	return fi.Size() < size, nil
}

// Checks that the filename in the current query value contains
// any of the file extensions in the list.
func (f *Filename) hasExt(name string, extensions []string) bool {
	hit := false
	_, _, cExt := f.decompose(name)
	cExt = strings.ToLower(cExt)
	for _, ext := range extensions {
		if cExt == strings.ToLower(ext) {
			hit = true
			break
		}
	}

	return hit
}

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                  M A I N    |    D E M O
 *-----------------------------------------------------------------*/

/*
func FilenameOpsDemo() {
	// Single operations mode
	fnS := NewFilename("/usr/share/apache/config/apache_config.cfg")
	fnS.RemoveExt() // /usr/share/apache/config/apache_config
	fnS.ReplaceExt(".txt") // /usr/share/apache/config/apache_config.txt
	ok, err := fnS.IsDirectory() // false
	fnS.IsInPath("/tmp") // false

	// Query mode
	fnQ := NewFilename("/tmp/apache_config.cfg").
		With().
		ChangeExt("txt").	// /tmp/apache_config.txt
		ChangeDirectory("/usr/share") // /usr/share/apache_config.txt
	result := fnQ.QueryStrResult(true)
}
*/
