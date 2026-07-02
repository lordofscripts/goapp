//go:build unix

/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Lord of Scripts
 *							   goApp
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Unix-specific code for FileSys.
 *-----------------------------------------------------------------*/
package app

import (
	"path"
	"syscall"
)

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

func GetUserTempDir() string {
	return "/tmp"
}

// @todo Use golang.org/x/sys package for better portability?
func Dup2(oldfd, newfd any) error {
	var oldFD, newFD int
	oldFD, _ = oldfd.(int)
	newFD, _ = newfd.(int)
	return syscall.Dup2(oldFD, newFD)
}

// (Linux|Unix|Darwin) whether the file is hidden, i.e. its
// name starts with a leading period.
func IsHiddenFile(filename string) (bool, error) {
	return path.Base(filename)[0] == '.', nil
}
