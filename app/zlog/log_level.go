/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                        goApp:zlog
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * LogLevel enumeration and methods. It provides a custom JSON
 * marshaller so that it is stored/retrieved by its string value
 * rather than the numeric value.
 *-----------------------------------------------------------------*/
package zlog

import (
	"fmt"
	"strings"
)

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/

const (
	// Logging level enumeration (0..6)
	LevelNone LogLevel = iota
	LevelTrace
	LevelDebug
	LevelInfo
	LevelWarning
	LevelError
	LevelFatal
)

/* ----------------------------------------------------------------
 *                     I N T E R F A C E S
 *-----------------------------------------------------------------*/

var _ fmt.Stringer = (*LogLevel)(nil)

/* ----------------------------------------------------------------
 *                         T Y P E S
 *-----------------------------------------------------------------*/

// Logging Level: Trace, Debug, Info, Warning, Error, Fatal
type LogLevel int

/* ----------------------------------------------------------------
 *                        M E T H O D S
 *-----------------------------------------------------------------*/

// implements fmt.Stringer for LogLevel returning its name
// without the "Level" prefix.
func (l LogLevel) String() string {
	vmap := map[LogLevel]string{
		LevelNone:    "None",
		LevelTrace:   "Trace",
		LevelDebug:   "Debug",
		LevelInfo:    "Info",
		LevelWarning: "Warn",
		LevelError:   "Error",
		LevelFatal:   "Fatal",
	}

	if v, ok := vmap[l]; ok {
		return v
	}

	return "++InvalidLevel++"
}

// Parse a string as a LogLevel and return it.
// if it fails it returns the current value.
// Example: x := LevelInfo; value, err := x.Parse("Trace")
func (l *LogLevel) Parse(str string) (LogLevel, error) {
	str = strings.ToLower(strings.TrimSpace(str))
	tmap := map[string]LogLevel{
		"disabled": LevelNone,
		"none":     LevelNone,
		"trace":    LevelTrace,
		"debug":    LevelDebug,
		"info":     LevelInfo,
		"warn":     LevelWarning,
		"warning":  LevelWarning,
		"error":    LevelError,
		"fatal":    LevelFatal,
	}

	if t, ok := tmap[str]; ok {
		return t, nil
	} else {
		return *l, fmt.Errorf("could not parse '%s' as LogLevel, returning current value %s", str, *l)
	}
}

// implements encoding.TextMarshaler for both JSON & YAML
func (l LogLevel) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

// implements encoding.TextUnmarshaler for both JSON & YAML
func (l *LogLevel) UnmarshalText(text []byte) error {
	var err error
	*l, err = l.Parse(string(text))
	return err
}

/*
// JSON marshaller so that it is written for its Text value
// instead of numeric value.
func (l LogLevel) MarshallJSON() ([]byte, error) {
	return  json.Marshal(l.String())
}

// JSON unmarshaller that assumes it was marshalled as text
// rather than numeric.
func (l *LogLevel) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if lev, err := l.Parse(s); err != nil {
		return err
	} else {
		*l = lev
	}
	return nil
}

// Custom YAML marshaller to string rather than number.
func (l *LogLevel) UnmarshalYAML(value *yaml.Node) error {
	// it should be a string
	var name string
	if err := value.Decode(&name); err != nil {
		return err
	}

	// it was a string, continue
	if lev, err := l.Parse(name); err != nil {
		return err
	} else {
		*l = lev
	}
	return nil
}

// Custom YAML marshalling of enumeration, otherwise it appears as integer.
func (l LogLevel) MarshalYAML() (interface{}, error) {
	return l.String(), nil
}
*/

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

// parse a string to convert it to a LogLevel value
func parseLevel(s string) LogLevel { // @audit deprecate
	var lvl LogLevel
	s = strings.Trim(s, " \t")

	switch {
	case strings.EqualFold(s, "none") || strings.EqualFold(s, "disabled"):
		lvl = LevelNone

	case strings.EqualFold(s, "trace"):
		lvl = LevelTrace

	case strings.EqualFold(s, "debug"):
		lvl = LevelDebug

	case strings.EqualFold(s, "info"):
		lvl = LevelInfo

	case strings.EqualFold(s, "warning"):
		fallthrough
	case strings.EqualFold(s, "warn"):
		lvl = LevelWarning

	case strings.EqualFold(s, "error"):
		lvl = LevelError

	case strings.EqualFold(s, "fatal"):
		lvl = LevelFatal

	default:
		lvl = LevelFatal
	}

	return lvl
}
