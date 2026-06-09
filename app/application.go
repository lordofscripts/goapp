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

	"github.com/lordofscripts/goapp/app/mlog"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

const (
	UC_RED_EXCLAMATION = rune(0x2757) // Dingbats
)

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

// Death of an application by outputting a good-bye and setting
// the OS exit code. It is logged as fatal.
func Die(message string, exitCode int) {
	fmt.Fprintln(os.Stderr, "\n", "\t💀 x 💀 x 💀\n\t", message, "\n\tExit code: ", exitCode)
	mlog.FatalT(exitCode, message, mlog.YesNo("Died", true), mlog.Int("Code", exitCode))
}

func DieWith(exitCode int, format string, args ...any) {
	err := fmt.Errorf(format, args...)
	DieWithError(err, exitCode)
}

// display the error and die with an exit code, logging it as Fatal.
func DieWithError(err error, exitCode int) {
	fmt.Fprintln(os.Stderr, "\n", "\t💀 x 💀 x 💀\n\t", err.Error(), "\n\tExit code: ", exitCode)
	mlog.FatalT(exitCode, err.Error(), mlog.YesNo("Died", true), mlog.Int("Code", exitCode))
}

// When the condition is met the warning message is printed
func Assert(condition bool, warnMessage string) {
	if condition {
		fmt.Fprintf(os.Stderr, "\n\t%c Assertion Failed:\n\t%s\n", UC_RED_EXCLAMATION, warnMessage)
	}
}

// If the condition is met, the death message is printed and the
// application terminates with the exit code.
func AssertOrDie(condition bool, deathMessage string, exitCode int) {
	if condition {
		fmt.Fprintf(os.Stderr, "\n\t%c Assertion Failed:", UC_RED_EXCLAMATION)
		Die(deathMessage, exitCode)
	}
}

// prints the error message with the exit code but does NOT exit.
func AnnounceErrorMessage(message string, exitCode int) {
	fmt.Fprintln(os.Stderr, "\n", "\t💀 x 💀 x 💀\n\t", message, "\n\tExit code: ", exitCode)
}

// prints the error and exit code but does NOT exit the application.
func AnnounceError(err error, exitCode int) {
	fmt.Fprintln(os.Stderr, "\n", "\t💀 x 💀 x 💀\n\t", err.Error(), "\n\tExit code: ", exitCode)
}

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
			}
		}
	}

	return cfgPath, err
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
