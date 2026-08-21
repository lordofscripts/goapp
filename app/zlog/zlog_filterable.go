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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/

const _FILTER_SCHEMA_VERSION uint = 1

/* ----------------------------------------------------------------
 *                     I N T E R F A C E S
 *-----------------------------------------------------------------*/

// The filter service is an extra layer over the Log Level aware
// logging to filter out (skip logging) for selected GO packages or
// to a certain object within a package.
// NOTES:
//   - IsFiltered*() returns true to skip the log entry (blacklisted)
//   - Any Log call is filtered out IF there are no filters loaded
//   - Any Log call is filtered out IF the package OR package/object
//     is NOT listed in the Filter file
//   - If the Package or Package/Object IS listed in the filter, then
//     it evaluates the filter level.
type ILogFilterService interface {
	// This only serves to create the FIRST filter configuration as
	// implemented by the service, i.e. to setup an empty filter config.
	// If the application already has an existing filter file (or whatever
	// storage) do NOT call this or you will lose your tweaked filters.
	Init() error
	// Save the current Log Filter configuration for this application. If
	// they have been modified by the app /(unlikely).
	SaveFilters() error
	// Load this application's Log Filters.
	LoadFilters() error
	// Check if this GO package name is filtered (blacklisted) or
	// filtered to None. If blacklisted the log output is skipped.
	IsFiltered(packageName string) bool
	// Check if this GO object in that package is filtered (blacklisted)
	// or filtered to None. If blacklisted the log output is skipped.
	IsFilteredObject(packageName, objectName string) bool
	// Load/Save configuration file in YAML rather than JSON.
	UseYaml()

	// It should return the location of the filter specification.
	// For a real file the fully-qualified filename, for a database-driven
	// config "db::SERVER:PORT/DATABASE_NAME"
	fmt.Stringer
}

var _ ILogFilterService = (*LogFilterService)(nil)
var _ ILogFilterService = (*NullFilterService)(nil)

/* ----------------------------------------------------------------
 *                   P U B L I C    T Y P E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                   P R I V A T E    T Y P E S
 *-----------------------------------------------------------------*/

// It provides the Log Filter configuration service which includes
// storage, retrieval, validation and (un)marshalling of the filter's
// JSON configuration file.
type LogFilterService struct {
	Schema       uint                 `json:"schema_version"`
	AppName      string               `json:"appname"`
	CurrentLevel LogLevel             `json:"level"` // allowed log level: None, Trace, Debug, Info, Warning, Error, Fatal
	Filters      map[string]LogFilter `json:"filters"`
	fd           *os.File
	fdCallTree   *os.File
	configSub    string
	useYaml      bool
}

// Part of the JSON log filter configuration file structure.
type LogFilter struct {
	LogLevel     LogLevel `json:"log_level"`
	Specifically string   `json:"specifically,omitempty"`
}

type NullFilterService struct{}

/* ----------------------------------------------------------------
 *                    C O N S T R U C T O R S
 *-----------------------------------------------------------------*/

// The configSubGroup is a subdirectory of the application configuration
// directory. The appname is detected from startup arguments.
func NewLogFilterService(appName, configSubGroup string) *LogFilterService {
	me := newLogFilterService(appName, configSubGroup)
	if len(appName) == 0 {
		me.AppName = me.getAppName()
	}
	return me
}

// don't use this for debugging in VSCode because the AppName will
// default to __debug******.bin
func NewLogFilterServiceAuto(configSubGroup string) *LogFilterService {
	me := newLogFilterService("", configSubGroup)
	me.AppName = me.getAppName()
	return me
}

func newLogFilterService(appName, configSubGroup string) *LogFilterService {
	return &LogFilterService{
		Schema:    _FILTER_SCHEMA_VERSION,
		AppName:   appName,
		configSub: configSubGroup,
		Filters:   make(map[string]LogFilter),
		useYaml:   false,
	}
}

// A null filter does not blacklist any package or object and pretends
// it saves and loads configuration.
func NewNullLogFilterService() *NullFilterService {
	return &NullFilterService{}
}

/* ----------------------------------------------------------------
 *                 P U B L I C    M E T H O D S
 *-----------------------------------------------------------------*/

func (svc *LogFilterService) Init() error {
	clear(svc.Filters)
	svc.Filters["main"] = LogFilter{LogLevel: LevelDebug, Specifically: ""}
	svc.Filters["module/package1"] = LogFilter{LogLevel: LevelInfo, Specifically: "StructA,StructB"}
	svc.Filters["module/package2"] = LogFilter{LogLevel: LevelNone, Specifically: "StructC"}

	return svc.SaveFilters()
}

// Saves a sample Log Filter configuration file in the user's configuration
// directory. The filename is APPNAME.logfilter.
func (svc *LogFilterService) SaveFilters() error {
	buffer := new(bytes.Buffer)
	if svc.useYaml {
		encoder := yaml.NewEncoder(buffer)
		if err := encoder.Encode(svc); err != nil {
			println("YAML encoder failure.", err)
			return err
		}
	} else {
		encoder := json.NewEncoder(buffer)
		encoder.SetEscapeHTML(false)
		encoder.SetIndent("", "  ")

		if err := encoder.Encode(svc); err != nil {
			println("JSON encoder failure.", err)
			return err
		}
	}

	if fd, err := os.Create(svc.getLogFilterFile()); err == nil {
		defer fd.Close()
		fd.WriteString(buffer.String())
	} else {
		println("couldn't save filter file")
		return err
	}
	return nil
}

// Attempts to load APPNAME.logfilters which contains a list of application/module
// packages white-listed for log output. It may contain a fine-grained configuration
// of which objects (structs) of that package are allowed. If empty all are allowed,
// else whatever struct name that is not listed becomes black-listed for logging.
func (svc *LogFilterService) LoadFilters() error {
	var all LogFilterService
	binDataBytes, err := os.ReadFile(svc.getLogFilterFile())
	if err != nil {
		println("Error reading filter file:", err)
		return err
	}

	if svc.useYaml {
		if err = yaml.Unmarshal(binDataBytes, &all); err != nil {
			println("error marshalling YAML filter:", err)
			return err
		}
	} else {
		if err = json.Unmarshal(binDataBytes, &all); err != nil {
			println("error marshalling JSON filter:", err)
			return err
		}
	}

	svc.Filters = all.Filters
	return nil
}

// Checks if the packageName (GO package name format) is
// black-listed. Black-listed packages do not produce log
// output under LOGX.
func (svc *LogFilterService) IsFiltered(packageName string) bool {
	if filter, exists := svc.Filters[packageName]; !exists ||
		exists && filter.LogLevel == LevelNone {
		return true
	}

	return false
}

// Checks if the packageName (GO package name format) is
// black-listed. Black-listed packages do not produce log
// output under LOGX.
func (svc *LogFilterService) IsFilteredObject(packageName, objectName string) bool {
	var blackListed bool = true
	filter, exists := svc.Filters[packageName]
	if exists {
		if filter.Specifically == "*" {
			blackListed = false
		} else if filter.LogLevel != LevelNone && strings.Contains(filter.Specifically, objectName) {
			blackListed = false
		}
	}
	/*
		if filter, exists := l.Filters[packageName]; !exists ||
			exists && filter.LogLevel == "" ||
			exists && filter.LogLevel != "" && !strings.Contains(filter.Specifically, objectName) {
			return true
		}

		return false
	*/
	return blackListed
}

func (svc *LogFilterService) UseYaml() {
	svc.useYaml = true
}

// It should return the location of the filter specification.
// For a real file the fully-qualified filename, for a database-driven
// config "db::SERVER:PORT/DATABASE_NAME"
func (svc *LogFilterService) String() string {
	return svc.getLogFilterFile()
}

// No-Op
func (null *NullFilterService) Init() error {
	return nil
}

// No-Op
func (null *NullFilterService) SaveFilters() error {
	return nil
}

// No-Op
func (null *NullFilterService) LoadFilters() error {
	return nil
}

// Returns false, the NULL filter does not blacklist
func (null *NullFilterService) IsFiltered(packageName string) bool {
	return false
}

// Returns false, the NULL filter does not blacklist
func (null *NullFilterService) IsFilteredObject(packageName, objectName string) bool {
	return false
}

// No-Op
func (null *NullFilterService) UseYaml() {}

func (null *NullFilterService) String() string {
	return "(dummy)"
}

/* ----------------------------------------------------------------
 *                 P R I V A T E    M E T H O D S
 *-----------------------------------------------------------------*/

// Derive the application name from the executable. Strip ".exe"
// if present.
func (svc *LogFilterService) getAppName() string {
	if svc.AppName != "" {
		return svc.AppName
	}
	appName := filepath.Base(os.Args[0])
	appName = strings.TrimSuffix(appName, ".exe")

	return appName
}

// Full path to JSON (appname_logfilter.json) OR
// YAML (appname_logfilter.yml) file with Log Filter configuration.
func (svc *LogFilterService) getLogFilterFile() string {
	var cfgPath string
	if usrConfig, err := os.UserConfigDir(); err != nil {
		return ""
	} else {
		var basename string = svc.getAppName() + "_logfilter"
		if svc.useYaml {
			basename += ".yaml"
		} else {
			basename += ".json"
		}
		cfgPath = filepath.Join(usrConfig, svc.configSub, basename)
	}

	return cfgPath
}

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/
