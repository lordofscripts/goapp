/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package zlog

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/lordofscripts/goapp/osx"
)

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                     I N T E R F A C E S
 *-----------------------------------------------------------------*/

// Exactly what the standard library "log" implements
type PlainLogger interface {
	Fatal(v ...any)
	Fatalf(format string, v ...any)
	Fatalln(v ...any)
	Flags() int
	Output(calldepth int, s string) error
	Panic(v ...any)
	Panicf(format string, v ...any)
	Panicln(v ...any)
	Prefix() string
	Print(v ...any)
	Printf(format string, v ...any)
	Println(v ...any)
	SetFlags(flag int)
	SetOutput(w io.Writer)
	SetPrefix(prefix string)
	Writer() io.Writer
}

var _ PlainLogger = (*BufferedLogger)(nil)

/* ----------------------------------------------------------------
 *                         T Y P E S
 *-----------------------------------------------------------------*/

// A temporary buffered logger to delay the initial log output until
// after the application's greeter has been outputed.
// NOTE: Not used yet.
type BufferedLogger struct {
	mu       sync.Mutex
	buffer   bytes.Buffer
	output   *log.Logger
	released bool
}

/* ----------------------------------------------------------------
 *                    C O N S T R U C T O R S
 *-----------------------------------------------------------------*/

// (Ctor) instance of a buffered logger that can be flushed.
func NewBufferedLogger() *BufferedLogger {
	return &BufferedLogger{
		output: log.New(os.Stderr, "", log.LstdFlags),
	}
}

/* ----------------------------------------------------------------
 *                        M E T H O D S
 *-----------------------------------------------------------------*/

// Ready() releases the in-memory buffered log entries into the
// mainstream log output. It is usually called just after the
// application welcome stanza is completed.
// Pre-condition:
//   - All log calls are stored in-memory
//
// Post-condition:
//   - Current in-memory log entries flushed to actual log stream.
//   - All log calls are directed to the real log stream.
func (l *BufferedLogger) Ready() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.released = true
	_, _ = l.output.Writer().Write(l.buffer.Bytes())
	l.buffer.Reset()
}

// Flush() causes any in-memory (buffered) log entries to be
// flushed into the toilet, I mean, the real log stream.
func (l *BufferedLogger) Flush() {
	l.mu.Lock()
	defer l.mu.Unlock()

	_, _ = l.output.Writer().Write(l.buffer.Bytes())
	l.buffer.Reset()
}

// If not Ready() it writes to memory, else to the actual log stream.
// NOTE: if format is empty, it behaves like Print().
func (l *BufferedLogger) Printf(format string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	unformatted := len(format) == 0
	if l.released {
		if unformatted {
			l.output.Print(args...)
		} else {
			l.output.Printf(format, args...)
		}
		return
	}

	fmt.Fprintf(&l.buffer, format+"\n", args...)
}

// If not Ready() it writes to memory, else to the actual log stream.
// Exits with code 1.
func (l *BufferedLogger) Print(args ...any) {
	l.print(nil, args...)
}

// If not Ready() it writes to memory, else to the actual log stream.
// Exits with code 1.
func (l *BufferedLogger) Println(args ...any) {
	args = append(args, osx.EOL)
	l.print(nil, args...)
}

func (l *BufferedLogger) print(format *string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	unformatted := format == nil
	if l.released {
		if unformatted {
			l.output.Print(args...)
		} else {
			l.output.Printf(*format, args...)
		}
		return
	}

	// Still buffered (not Ready)
	if unformatted {
		fmt.Fprint(&l.buffer, args...)
		fmt.Fprint(&l.buffer, osx.EOL)
	} else {
		fmt.Fprintf(&l.buffer, *format+osx.EOL, args...)
	}
}

// If not Ready() it writes to memory, else to the actual log stream.
// Exits with code 1.
func (l *BufferedLogger) Fatal(args ...any) {
	l.Flush() // Required: os.Exit skips defers.
	l.output.Print(args...)
	os.Exit(1)
}

// If not Ready() it writes to memory, else to the actual log stream.
// Exits with code 1.
func (l *BufferedLogger) Fatalln(args ...any) {
	args = append(args, osx.EOL)
	l.Fatal(args...)
}

// If not Ready() it writes to memory, else to the actual log stream.
// Exits with code 1.
func (l *BufferedLogger) Fatalf(format string, args ...any) {
	l.Flush() // Required: os.Exit skips defers.
	l.output.Printf(format, args...)
	os.Exit(1)
}

// If not Ready() it writes to memory, else to the actual log stream.
// Panics after data is written.
func (l *BufferedLogger) Panicf(format string, args ...any) {
	l.Flush() // Optional if a caller has no recovery defer.
	panic(fmt.Sprintf(format, args...))
}

func (l *BufferedLogger) Panic(args ...any) {
	l.Flush() // Optional if a caller has no recovery defer.
	panic(fmt.Sprint(args...))
}

func (l *BufferedLogger) Panicln(args ...any) {
	l.Flush() // Optional if a caller has no recovery defer.
	args = append(args, osx.EOL)
	panic(fmt.Sprint(args...))
}

func (l *BufferedLogger) SetFlags(flag int) {
	l.output.SetFlags(flag)
}

func (l *BufferedLogger) Flags() int {
	return l.output.Flags()
}

func (l *BufferedLogger) Output(calldepth int, s string) error {
	return l.output.Output(calldepth, s)
}

func (l *BufferedLogger) Prefix() string {
	return l.output.Prefix()
}

func (l *BufferedLogger) SetPrefix(prefix string) {
	l.output.SetPrefix(prefix)
}

func (l *BufferedLogger) SetOutput(w io.Writer) {
	l.output.SetOutput(w)
}

func (l *BufferedLogger) Writer() io.Writer {
	return l.output.Writer()
}

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/
