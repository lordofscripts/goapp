//go:build !logx
// +build !logx

/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Lord of Scripts
 *							   goApp
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Production/Release version. Log module equivalents are mimicked
 * as-is, extra functionality is empty.
 *-----------------------------------------------------------------*/
package logx

import (
	"io"
	"log"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

func Prefix() string {
	return log.Prefix()
}

func SetPrefix(prefix string) {
	log.SetPrefix(prefix)
	log.SetFlags(log.LstdFlags | log.Lmsgprefix)
}

func Flags() int {
	return log.Flags()
}

func Print(v ...any) {
	log.Print(v...)
}

func Printf(format string, v ...any) {
	log.Printf(format, v...)
}

func Println(v ...any) {
	log.Println(v...)
}

func Fatal(v ...any) {
	log.Fatal(v...)
}

func Fatalf(format string, v ...any) {
	log.Fatalf(format, v...)
}

func Fatalln(v ...any) {
	log.Fatalln(v...)
}

//func Output(calldepth int, s string) error {}

func Panic(v ...any) {
	log.Panic(v...)
}

func Panicf(format string, v ...any) {
	log.Panicf(format, v...)
}

func Panicln(v ...any) {
	log.Panicln(v...)
}

func SetFlags(flag int) {
	log.SetFlags(flag)
}

func SetOutput(w io.Writer) {
	log.SetOutput(w)
}

/* ----------------------------------------------------------------
 *				E x t e n d e d 	F u n c t i o n s
 *-----------------------------------------------------------------*/

/**
 * For use to log that an object constructor is executing
 */
func Ctor() {}

/**
 * For use at entering a function/method of an EVENT callback
 */
func EventEnter() {}

/**
 * For use at the end of a function/method of an EVENT callback
 */
func EventLeave() {}

/**
 * For use at entering a function/method
 */
func Enter() {}

/**
 * For use at the end of a function/method
 */
func Leave() {}

/**
 * For use when entering a function/method where you don't want Enter+Leave
 * but just logging your visit.
 */
func Visit() {}

func Step(string) {}

// Log an error irrespective of the log level
func AttentionAlways(short string, err error) {
	log.Printf("%c ERR %s %s", UC_EYES, short, err)
}

// In Release it is the same as AttentionAlways
func Attention(short string, err error) {
	AttentionAlways(short, err)
}

func Result(format string, v ...any) {}

func OnValidating() {}

func OnChanged(toValue ...any) {}

func OnUpdate() {}

func OnCascade(string, any) {}

func OnClick() {}

/* ----------------------------------------------------------------
 *				P r i v a t e	F u n c t i o n s
 *-----------------------------------------------------------------*/
