/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025-2026 Dídimo Grimaldo T.
 *							   go-app
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Application-related functions.
 *-----------------------------------------------------------------*/
package app

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

const (
	UC_RED_EXCLAMATION = rune(0x2757) // ❗Dingbats
	UC_RED_CROSS       = rune(0x274C) // ❌
)

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

// Returns true if the application input is not from a character device (tty)
// but instead from a piped input like "cat textfile.txt | yourapp -encrypt".
// When true you can use a bufio.Scanner to read text lines one by one and
// process them accordingly.
func IsPipedInput() bool {
	fi, _ := os.Stdin.Stat()
	// This is a workaround for VSCode debug sessions which would otherwise
	// make this function return true. For VSCode debugger set the env
	// in launch.json to DD_NOT_PIPED=1
	if os.Getenv("DD_NOT_PIPED") == "1" {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) == 0
}

// platform-agnostic function to obtain the user's configuration directory.
// In Linux "~/.config/orgName, appName", Windows "APPDATA\orgName\appName" and
// MacOS "~/Library/Application Support/orgName/appName"
func GetOSConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch runtime.GOOS {
	case "windows":
		return os.Getenv("APPDATA"), nil
	case "darwin": // macOS
		return filepath.Join(homeDir, "Library", "Application Support"), nil
	default: // Other platforms (Linux, etc.)
		return filepath.Join(homeDir, ".config"), nil
	}
}

// platform-agnostic function to obtain the user's configuration directory.
// In Linux "~/.config/orgName, appName", Windows "APPDATA\orgName\appName" and
// MacOS "~/Library/Application Support/orgName/appName"
func GetConfigDir(orgName, appName string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("APPDATA"), orgName, appName), nil
	case "darwin": // macOS
		return filepath.Join(homeDir, "Library", "Application Support", orgName, appName), nil
	default: // Other platforms (Linux, etc.)
		return filepath.Join(homeDir, ".config", orgName, appName), nil
	}
}

// Ensures a directory and all its parents exist and create them if necessary.
// Default permissions is 0750.
func EnsureConfigDir(path string) error {
	// Create the config directory if it doesn't exist
	err := os.MkdirAll(path, 0750) // 0755 permissions: rwxr-x---
	if err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	return nil
}

// Ensures the platform-aware configuration file and directory exist
// for orgName and appName. The application configuration file has
// the same name as the application (appName) plus the extension.
// If fileExtension is missing the leading period, it is automatically
// added. It returns the fully-qualified configuration filename and
// whether an error occurred or not.
func EnsureConfig(orgName, appName, fileExtension string) (string, error) {
	var err error = nil
	var cfgPath string
	// build name of platform-aware configuration directory
	if cfgPath, err = GetConfigDir(orgName, appName); err == nil {
		if err = EnsureConfigDir(cfgPath); err == nil {
			// ensure correct file extension with leadig period
			fileExtension = strings.Trim(fileExtension, " \t")
			if !strings.HasPrefix(fileExtension, ".") {
				fileExtension = "." + fileExtension
			}
			// the config file has the name of the app plus extension
			filePath := filepath.Join(cfgPath, appName+fileExtension)
			if err = CheckFileExistsAndReadable(filePath); err == nil {
				return filePath, nil
			} else {
				// at least returned the expected filename
				return filePath, err
			}
		}
	}

	return cfgPath, err // no home directory
}

// Same as EnsureConfig except it adds the suffix to the basename. Thus,
// with orgName "ACME", appName "goapp" and nameSuffix "_template" and fileExtension "json",
// the resulting name would be ~/PathToConfig/ACME/goapp_template.json
func EnsureConfigWithSuffix(orgName, appName, nameSuffix, fileExtension string) (string, error) {
	var err error = nil
	var cfgPath string
	// build name of platform-aware configuration directory
	if cfgPath, err = GetConfigDir(orgName, appName); err == nil {
		if err = EnsureConfigDir(cfgPath); err == nil {
			// ensure correct file extension with leadig period
			fileExtension = strings.Trim(fileExtension, " \t")
			if !strings.HasPrefix(fileExtension, ".") {
				fileExtension = "." + fileExtension
			}
			// the config file has the name of the app plus extension
			filePath := filepath.Join(cfgPath, appName+nameSuffix+fileExtension)
			if err = CheckFileExistsAndReadable(filePath); err == nil {
				return filePath, nil
			} else {
				// at least returned the expected filename
				return filePath, err
			}
		}
	}

	return cfgPath, err // no home directory
}

// Checks whether the file exists and is readable.
func CheckFileExistsAndReadable(filePath string) error {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", filePath)
		}
		return fmt.Errorf("error checking file: %w", err)
	}
	return nil
}
