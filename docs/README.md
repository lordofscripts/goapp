# Go App - Spice up your GO applications

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/lordofscripts/goapp)
[![Go Report Card](https://goreportcard.com/badge/github.com/lordofscripts/goapp?style=flat-square)](https://goreportcard.com/report/github.com/lordofscripts/goapp)
![Build](https://github.com/lordofscripts/goapp/actions/workflows/go.yml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/lordofscripts/goapp.svg)](https://pkg.go.dev/github.com/lordofscripts/goapp)
[![GitHub release (with filter)](https://img.shields.io/github/v/release/lordofscripts/goapp)](https://github.com/lordofscripts/goapp/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)


When I was developing on the Microsoft Windows OS, I was a hard-core 
Microsoft .NET developer. I lived and breathed *.NET* and *Visual Studio*.

However, a few years ago I sort of abandoned Microsoft Windows as my main
development platform. I did what was sensible, I returned to my roots: Linux.
I had been a Linux-based OS developer since its early versions, prior to
kernel 1.0 when you basically had to construct the car before you could
dream of driving it. Back then I used C, C++ and Java, plus a lot of
scripting languages, Perl was my favorite.

When I returned to my Linux roots, I found out .NET (Mono) was neither
stable nor viable serious alternative on Linux. I learned Python, but
I must say I dislike interpreted languages.

But then I discovered Go (somewhat late I must say!) and I have been
doing a lot of Go development ever since! I write a lot of command-line
applications, some GUI (using Fyne) too. But one thing was obvious, I
found myself needing the same things over and over again, so why not
bundle them in a single module and rewire my applications to use this
small contribution instead? I finally decided to it when developing
*Go CaesarDisk* , so here it is.

*This is stable but still work in progress as I move reusable code
from some of my applications into this Go module.*

## Install

> go get github.com/lordofscripts/goapp@latest

## Features

* Your choice of enhanced, non-structured logging packages: `app/logx`
  or `app/mlog`. These are enhanced versions of the standard `log`
  and `log/slog` packages. Or also level-agnostic colored logging.
* Multiplatform utility functions such as `DirExists` and `FileExists`,
  others will be added as I see fit.
* Assertion functions: `Assert*()`
* Peaceful Death functions when your app needs to rest in peace: `Die*()`
* Check whether your input is from a pipe: `IsPipeInput()`
* Obtain a platform-agnostic configuration directory: `GetConfigDir()`
  and ensure it is present or created for your application: `EnsureConfigDir()`

### Logging

You have two choices of enhanced logging ready to use:

* [MLog](./MLOG.md) for tagged, logging like the standard `log/slog` but
  with extra tags and support for logging levels, or
* [Colored MLog](./MLOG.md#colored-logging) is a variant of MLog that
  always logs to a colored console, it supports log levels by prefixing
  them but the output happens *regardless* of the actual log level.
* [Logx](LOGX.md) preformatted logging with call-tree support, it was
  instrumental to get one of my Fyne GUI applications to work.
