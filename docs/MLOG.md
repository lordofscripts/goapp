# MLOG - Mary Logs

What can I say, sometimes I have different logging needs in my applications.
What is certain is that I am no fan of structured logs, i.e. with XML/JSON.

*Mary Logs* is another variation of the standard GO `log/slog` package, it builds
upon that rather than rewriting it. Simply because I liked the ability of
`slog` to use key-value pairs; however, it turned out it missed some things
I needed so I created `app/mlog`.

## Features

* Improved version of the standard `log/slog`
* Supports *logging levels* (`slog` does not)
* Added extra logging key-value tags: String, Bool, YesNo, Int, Byte, Rune and At.
* Log output appears on `stderr`
* Log lines prefixed with timestamp: format `2006-01-02 15:04:05`
* Colored logging to the console (not to a file)

## Installing

First get the package added to your application:

> go get github.com/lordofscripts/goapp@latest

Then in your application source:

> import "github.com/lordofscripts/goapp/app/mlog"

### Usage

To enable logging on your executable:

> go build -tags mlog PATH

When the tag is omitted, only strictly-needed log output is preserved.
This helps differentiating between development and release builds.

Set environment variables to configure `app/mlog`, I prefer to use `.env`:

> LOG_LEVEL_CX=debug
> LOG_FILE_CX=/tmp/myapp.log

If you don't want to log to a file, and let the log output to go to
`stderr` then leave `LOG_FILE_CX` undefined or empty. Normally I rig up my 
VSCode `launch.json` with:

> "envFile": "${workspaceRoot}/.env"

Then instrument your code accordingly to output log messages by using a
combination of the following:

```go
	LevelTrace LogLevel = iota  // trace
	LevelDebug      // debug
	LevelInfo       // info
	LevelWarning    // warning | warn
	LevelError      // error
	LevelFatal      // fatal
```

The default level is `error` meaning that by default Error and Fatal
log messags get outputed.

Each log level has a set of logging functions to suit different purposes:

> func Warn(v ...any)

Uses variadic parameters in free form, i.e. no format string.

> func Warnf(format string, v ...any)

Uses variadic parameters but with a format string just like `fmt.Printf`.
As a behavioral bonus, the variadic parameters, given that are of 
type `any` means it also accepts `ILogKeyValuePair` tags like the next
function.

> func WarnT(message string, v ...ILogKeyValuePair)

This one also uses variadic parameters, but without a format string. It
does have a short message before the parameters. In this case however,
it follows the `log/slog` form of log key-value tags. 

#### Key-Value Tags

These are key-value tags/pairs that can be used in the logging function
parameters:

To output a string key-value:

> func String(key, value string) ILogKeyValuePair

To output a `rune` key-value that is printed as a character rather
than its underlying integer:

> func Rune(key string, value rune) ILogKeyValuePair

To output an integer key-value pair:

> func Int(key string, value int) ILogKeyValuePair

To output a boolean key-value pair:

> func Bool(key string, value bool) ILogKeyValuePair

To output a boolean but as a Yes/No value:

> func YesNo(key string, value bool) ILogKeyValuePair

To output a byte (`uint8`):

> func Byte(key string, value byte) ILogKeyValuePair

To output the location from which the log function was called
(package/struct/method/function):

> func At() ILogKeyValuePair

To log an error:

> func Err(err error) ILogKeyValuePair

#### Colored Logging

If you feel like logging messages to the text console with a flair
for color, then use the global `app/mlog.Console` object.

```go
	mlog.Console.Trace("Trace %d\n", 1)
	mlog.Console.Debug("Debug %d\n", 2)
	mlog.Console.Info("Info %d\n", 3)
	mlog.Console.Warn("Warn %d\n", 4)
	mlog.Console.Error("Error %d\n", 5)
	mlog.Console.Fatal(110, "Fatal %d %s\n", 6, mlog.At())
```

Better yet, you can also use the MLog key-value tags introduced
earlier in this document. Keep in mind that this will output
the various levels by prefixing the log entry with the level's
shortform: TRC, DBG, INF, WRN, ERR, DIE, CAT. The output will
be printed *regardless* of the actual logging level.