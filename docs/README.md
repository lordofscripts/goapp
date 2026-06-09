# Go App v1.3 - Spice up your GO applications

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

Requirements:

* GO `v1.21`

## Features

* There are several constructors in the `app` package to declare your
  software version. It follows *Semantic Versioning* allowing for
  Alpha, Beta, RC and (final) releases.
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
* custom flags for the `flag` standard package: `flagx.RuneFlag`, 
  `flagx.ByteFlag`, `flagx.StringsFlag` and `flagx.DateFlag`  
* A convenience function to detect **piped** input for Linux/Unix apps.
  If you are debugging piped tests in VSCode, set the environment
  variable `DD_PIPED_INPUT=1`.

### Logging

You have two choices of enhanced logging ready to use:

* [MLog](./MLOG.md) for tagged, logging like the standard `log/slog` but
  with extra tags and support for logging levels, or
* [Colored MLog](./MLOG.md#colored-logging) is a variant of MLog that
  always logs to a colored console, it supports log levels by prefixing
  them but the output happens *regardless* of the actual log level.
* [Logx](LOGX.md) preformatted logging with call-tree support, it was
  instrumental to get one of my Fyne GUI applications to work.

### Command-line applications

For declaring your application/package versions:

```go
 const BASE_VERSION string = "1.3.2"
 const REVISION int = 8
 For Alpha version v1.3.2-Alpha.8
 pver := app.NewAlphaVersion(NAME, DESC, BASE_VERSION, REVISION)
 For Beta version v1.3.2-Beta.8
 pver := app.NewBetaVersion(NAME, DESC, BASE_VERSION, REVISION)
 For Release Candidate version v1.3.2-RC.8
 pver := app.NewReleaseCandidateVersion(NAME, DESC, BASE_VERSION, REVISION)
```

> For final release version v1.3.2
> pver := app.NewReleaseVersion(NAME, DESC, BASE_VERSION)

Getting the version number as displayed above:

> pver.Short()

The app name plus its version number:

> fmt.Println(pver.String())

Outputting your application copyright on the CLI:

> app.Copyright("Jan Knowsalot", '\u2615')

Outputing your [Buy Me a Coffee](https://buymeacoffee.com/lostinwriting)
profile URL:

> app.BuyMeCoffee("lordofscripts")

Then, should your application need to rest in peace:

```go
  const EXIT_CODE int = 5
  app.Die("bad thing happened", EXIT_CODE)
  app.DieWith(EXIT_CODE, "%d bad things happened", 3)
  app.DieWithError(err, EXIT_CODE)
```  

And there are a few more utility functions. Please check the package
documentation.

### Application Configuration

The package has built-in utility functions to let your application
use an application configuration file in a platform-aware fashion.
Thus, it works on MacOS, Windows and Linux but the directory where
it resides varies depending on the OS. There are several functions
but the one that makes use of all of them to accomplish all tasks
in one is:

```go
 const ORG_NAME string = "acme"
 const APP_NAME string = "coyote"
 const CONFIG_FILE_EXT string = ".json" // or .yaml or .ini
 configFile, err := app.EnsureConfig(ORG_NAME, APP_NAME, CONFIG_FILE_EXT)
```

Even when an error happens, the `configFile` return value will always
have the expected fully-qualified name of the configuration file.

The configuration *directory* varies per OS:

* **Linux & Unix**: `~/.config/ORG_NAME/APP_NAME`
* **MacOS**: `~/Library/Application Support/ORG_NAME/APP_NAME`
* **Microsoft Windows**: `C:\\Users\USERNAME\APPDATA\ORG_NAME\APP_NAME`
