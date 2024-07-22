# strftime

> fork from https://github.com/lestrrat-go/strftime.git

Fast strftime for Go

[![Build Status](https://travis-ci.org/lestrrat-go/strftime.png?branch=master)](https://travis-ci.org/lestrrat-go/strftime)

[![GoDoc](https://godoc.org/github.com/lestrrat-go/strftime?status.svg)](https://godoc.org/github.com/lestrrat-go/strftime)

# SYNOPSIS

```go
f, err := strftime.New(`.... pattern ...`)
if err := f.Format(buf, time.Now()); err != nil {
    log.Println(err.Error())
}
```

# DESCRIPTION

The goals for this library are

* Optimized for the same pattern being called repeatedly
* Be flexible about destination to write the results out
* Be as complete as possible in terms of conversion specifications

# API

## Format(string, time.Time) (string, error)

Takes the pattern and the time, and formats it. This function is a utility function that recompiles the pattern every time the function is called. If you know beforehand that you will be formatting the same pattern multiple times, consider using `New` to create a `Strftime` object and reuse it.

## New(string) (\*Strftime, error)

Takes the pattern and creates a new `Strftime` object.

## obj.Pattern() string

Returns the pattern string used to create this `Strftime` object

## obj.Format(io.Writer, time.Time) error

Formats the time according to the pre-compiled pattern, and writes the result to the specified `io.Writer`

## obj.FormatString(time.Time) string

Formats the time according to the pre-compiled pattern, and returns the result string.

# SUPPORTED CONVERSION SPECIFICATIONS

| pattern | description |
|:--------|:------------|
| %A      | national representation of the full weekday name |
| %a      | national representation of the abbreviated weekday |
| %B      | national representation of the full month name |
| %b      | national representation of the abbreviated month name |
| %C      | (year / 100) as decimal number; single digits are preceded by a zero |
| %c      | national representation of time and date |
| %D      | equivalent to %m/%d/%y |
| %d      | day of the month as a decimal number (01-31) |
| %e      | the day of the month as a decimal number (1-31); single digits are preceded by a blank |
| %F      | equivalent to %Y-%m-%d |
| %H      | the hour (24-hour clock) as a decimal number (00-23) |
| %h      | same as %b |
| %I      | the hour (12-hour clock) as a decimal number (01-12) |
| %j      | the day of the year as a decimal number (001-366) |
| %k      | the hour (24-hour clock) as a decimal number (0-23); single digits are preceded by a blank |
| %l      | the hour (12-hour clock) as a decimal number (1-12); single digits are preceded by a blank |
| %M      | the minute as a decimal number (00-59) |
| %m      | the month as a decimal number (01-12) |
| %n      | a newline |
| %p      | national representation of either "ante meridiem" (a.m.)  or "post meridiem" (p.m.)  as appropriate. |
| %R      | equivalent to %H:%M |
| %r      | equivalent to %I:%M:%S %p |
| %S      | the second as a decimal number (00-60) |
| %T      | equivalent to %H:%M:%S |
| %t      | a tab |
| %U      | the week number of the year (Sunday as the first day of the week) as a decimal number (00-53) |
| %u      | the weekday (Monday as the first day of the week) as a decimal number (1-7) |
| %V      | the week number of the year (Monday as the first day of the week) as a decimal number (01-53) |
| %v      | equivalent to %e-%b-%Y |
| %W      | the week number of the year (Monday as the first day of the week) as a decimal number (00-53) |
| %w      | the weekday (Sunday as the first day of the week) as a decimal number (0-6) |
| %X      | national representation of the time |
| %x      | national representation of the date |
| %Y      | the year with century as a decimal number |
| %y      | the year without century as a decimal number (00-99) |
| %Z      | the time zone name |
| %z      | the time zone offset from UTC |
| %%      | a '%' |

# EXTENSIONS / CUSTOM SPECIFICATIONS

This library in general tries to be POSIX compliant, but sometimes you just need that
extra specification or two that is relatively widely used but is not included in the
POSIX specification.

For example, POSIX does not specify how to print out milliseconds,
but popular implementations allow `%f` or `%L` to achieve this.

For those instances, `strftime.Strftime` can be configured to use a custom set of
specifications:

```
ss := strftime.NewSpecificationSet()
ss.Set('L', ...) // provide implementation for `%L`

// pass this new specification set to the strftime instance
p, err := strftime.New(`%L`, strftime.WithSpecificationSet(ss))
p.Format(..., time.Now())
```

The implementation must implement the `Appender` interface, which is

```
type Appender interface {
  Append([]byte, time.Time) []byte
}
```

For commonly used extensions such as the millisecond example and Unix timestamp, we provide a default
implementation so the user can do one of the following:

```
// (1) Pass a specification byte and the Appender
//     This allows you to pass arbitrary Appenders
p, err := strftime.New(
  `%L`,
  strftime.WithSpecification('L', strftime.Milliseconds),
)

// (2) Pass an option that knows to use strftime.Milliseconds
p, err := strftime.New(
  `%L`,
  strftime.WithMilliseconds('L'),
)
```

Similarly for Unix Timestamp:
```
// (1) Pass a specification byte and the Appender
//     This allows you to pass arbitrary Appenders
p, err := strftime.New(
  `%s`,
  strftime.WithSpecification('s', strftime.UnixSeconds),
)

// (2) Pass an option that knows to use strftime.UnixSeconds
p, err := strftime.New(
  `%s`,
  strftime.WithUnixSeconds('s'),
)
```

If a common specification is missing, please feel free to submit a PR
(but please be sure to be able to defend how "common" it is)

## List of available extensions

- [`Milliseconds`](https://pkg.go.dev/github.com/lestrrat-go/strftime?tab
