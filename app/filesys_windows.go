//go:build !unix

/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Lord of Scripts
 *							   goApp
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Windoze-specific code for FileSys.
 *-----------------------------------------------------------------*/
package app

import (
	"os"
	"syscall"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

var (
	kernel32         = syscall.MustLoadDLL("kernel32.dll")
	procSetStdHandle = kernel32.MustFindProc("SetStdHandle")
)

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

func GetUserTempDir() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		return `C:\`
	}
	return dir
}

// Dup2 provides an equivalent of the Unix dup2 function for Windows.
func Dup2(oldfd, newfd any) error {
	var r0 uintptr
	var e1 syscall.Errno

	switch oldfd.(type) {
	case int:
		var oldFD, newFD int
		oldFD, _ = oldfd.(int)
		newFD, _ = newfd.(int)
		r0, _, e1 = syscall.SyscallN(procSetStdHandle.Addr(), 2, uintptr(oldFD), uintptr(newFD), 0)
	case uintptr:
		var oldFD, newFD uintptr
		oldFD, _ = oldfd.(uintptr)
		newFD, _ = newfd.(uintptr)
		r0, _, e1 = syscall.SyscallN(procSetStdHandle.Addr(), 2, oldFD, newFD, 0)
	}

	if r0 == 0 {
		if e1 != 0 {
			return error(e1)
		}
		return syscall.EINVAL
	}
	return nil
}
