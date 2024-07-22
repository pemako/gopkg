package strftime

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// These are all of the standard, POSIX compliant specifications.
// Extensions should be in extensions.go
var (
	fullWeekDayName             = StdlibFormat("Monday")                   // 完整的星期名称
	abbrvWeekDayName            = StdlibFormat("Mon")                      // 缩写的星期名称
	fullMonthName               = StdlibFormat("January")                  // 完整的月份名称
	abbrvMonthName              = StdlibFormat("Jan")                      // 缩写的月份名称
	centuryDecimal              = AppendFunc(appendCentury)                // 世纪十进制
	timeAndDate                 = StdlibFormat("Mon Jan _2 15:04:05 2006") // 时间与日期
	mdy                         = StdlibFormat("01/02/06")                 // 月日年
	dayOfMonthZeroPad           = StdlibFormat("02")                       // 月份中的某天0填充
	dayOfMonthSpacePad          = StdlibFormat("_2")                       // 月份中的某天空格填充
	ymd                         = StdlibFormat("2006-01-02")               // 年月日
	twentyFourHourClockZeroPad  = &hourPadded{twelveHour: false, pad: '0'} // 24小时进制
	twelveHourClockZeroPad      = &hourPadded{twelveHour: true, pad: '0'}  // 12小时进制
	dayOfYear                   = AppendFunc(appendDayOfYear)              // 一年中的一天
	twentyFourHourClockSpacePad = &hourPadded{twelveHour: false, pad: ' '} // 24小时进制空格填充
	twelveHourClockSpacePad     = &hourPadded{twelveHour: true, pad: ' '}  // 12小时进制空格填充
	minutesZeroPad              = StdlibFormat("04")                       // 某一分
	monthNumberZeroPad          = StdlibFormat("01")                       // 某一月
	newline                     = Verbatim("\n")                           // 换行
	ampm                        = StdlibFormat("PM")                       // ap mp
	hm                          = StdlibFormat("15:04")                    // 时分
	imsp                        = hmsWAMPM{}
	secondsNumberZeroPad        = StdlibFormat("05")       // 某一秒
	hms                         = StdlibFormat("15:04:05") // 时分秒
	tab                         = Verbatim("\t")           // 换行
	weekNumberSundayOrigin      = weeknumberOffset(true)   // week number of the year, Sunday first
	weekdayMondayOrigin         = weekday(1)
	// monday as the first day, and 01 as the first value
	weekNumberMondayOriginOneOrigin = AppendFunc(appendWeekNumber)
	eby                             = StdlibFormat("_2-Jan-2006")
	// monday as the first day, and 00 as the first value
	weekNumberMondayOrigin = weeknumberOffset(false) // week number of the year, Monday first
	weekdaySundayOrigin    = weekday(0)
	natReprTime            = StdlibFormat("15:04:05") // national representation of the time XXX is this correct?
	natReprDate            = StdlibFormat("01/02/06") // national representation of the date XXX is this correct?
	year                   = StdlibFormat("2006")     // year with century
	yearNoCentury          = StdlibFormat("06")       // year w/o century
	timezone               = StdlibFormat("MST")      // time zone name
	timezoneOffset         = StdlibFormat("-0700")    // time zone ofset from UTC
	percent                = Verbatim("%")
)

// Appender is the interface that must be fulfilled by components that
// implement the translation of specifications to actual time value.
//
// The Append method takes the accumulated byte buffer, and the time to
// use to generate the textual representation. The resulting byte
// sequence must be returned by this method, normally by using the
// append() builtin function.
type Appender interface {
	Append([]byte, time.Time) []byte
}

// AppendFunc is an utility type to allow users to create a
// function-only version of an Appender
type AppendFunc func([]byte, time.Time) []byte

func (af AppendFunc) Append(b []byte, t time.Time) []byte {
	return af(b, t)
}

type appenderList []Appender

type dumper interface {
	dump(io.Writer)
}

func (l appenderList) dump(out io.Writer) {
	var buf bytes.Buffer
	ll := len(l)
	for i, a := range l {
		if dumper, ok := a.(dumper); ok {
			dumper.dump(&buf)
		} else {
			fmt.Fprintf(&buf, "%#v", a)
		}

		if i < ll-1 {
			fmt.Fprintf(&buf, ",\n")
		}
	}
	if _, err := buf.WriteTo(out); err != nil {
		panic(err)
	}
}

// does the time.Format thing
type stdlibFormat struct {
	s string
}

// StdlibFormat returns an Appender that simply goes through `time.Format()`
// For example, if you know you want to display the abbreviated month name for %b,
// you can create a StdlibFormat with the pattern `Jan` and register that
// for specification `b`:
//
// a  := StdlibFormat(`Jan`)
// ss := NewSpecificationSet()
// ss.Set('b', a) // does %b -> abbreviated month name
func StdlibFormat(s string) Appender {
	return &stdlibFormat{s: s}
}

func (v stdlibFormat) Append(b []byte, t time.Time) []byte {
	return t.AppendFormat(b, v.s)
}

func (v stdlibFormat) str() string {
	return v.s
}

func (v stdlibFormat) canCombine() bool {
	return true
}

func (v stdlibFormat) combine(w combiner) Appender {
	return StdlibFormat(v.s + w.str())
}

func (v stdlibFormat) dump(out io.Writer) {
	fmt.Fprintf(out, "stdlib: %s", v.s)
}

type verbatimw struct {
	s string
}

// Verbatim returns an Appender suitable for generating static text.
// For static text, this method is slightly favorable than creating
// your own appender, as adjacent verbatim blocks will be combined
// at compile time to produce more efficient Appenders
func Verbatim(s string) Appender {
	return &verbatimw{s: s}
}

func (v verbatimw) Append(b []byte, _ time.Time) []byte {
	return append(b, v.s...)
}

func (v verbatimw) canCombine() bool {
	return canCombine(v.s)
}

func (v verbatimw) combine(w combiner) Appender {
	if _, ok := w.(*stdlibFormat); ok {
		return StdlibFormat(v.s + w.str())
	}
	return Verbatim(v.s + w.str())
}

func (v verbatimw) str() string {
	return v.s
}

func (v verbatimw) dump(out io.Writer) {
	fmt.Fprintf(out, "verbatim: %s", v.s)
}

// These words below, as well as any decimal character
var combineExclusion = []string{
	"Mon",
	"Monday",
	"Jan",
	"January",
	"MST",
	"PM",
	"pm",
}

func canCombine(s string) bool {
	if strings.ContainsAny(s, "0123456789") {
		return false
	}
	for _, word := range combineExclusion {
		if strings.Contains(s, word) {
			return false
		}
	}
	return true
}

type combiner interface {
	canCombine() bool
	combine(combiner) Appender
	str() string
}

// this is container for the compiler to keep track of appenders,
// and combine them as we parse and compile the pattern
type combiningAppend struct {
	list           appenderList
	prev           Appender
	prevCanCombine bool
}

func (ca *combiningAppend) Append(w Appender) {
	if ca.prevCanCombine {
		if wc, ok := w.(combiner); ok && wc.canCombine() {
			ca.prev = ca.prev.(combiner).combine(wc)
			ca.list[len(ca.list)-1] = ca.prev
			return
		}
	}

	ca.list = append(ca.list, w)
	ca.prev = w
	ca.prevCanCombine = false
	if comb, ok := w.(combiner); ok {
		if comb.canCombine() {
			ca.prevCanCombine = true
		}
	}
}

func appendCentury(b []byte, t time.Time) []byte {
	n := t.Year() / 100
	if n < 10 {
		b = append(b, '0')
	}
	return append(b, strconv.Itoa(n)...)
}

type weekday int

func (v weekday) Append(b []byte, t time.Time) []byte {
	n := int(t.Weekday())
	if n < int(v) {
		n += 7
	}
	return append(b, byte(n+48))
}

type weeknumberOffset bool

func (v weeknumberOffset) Append(b []byte, t time.Time) []byte {
	offset := int(t.Weekday())
	if v {
		offset = 6 - offset
	} else if offset != 0 {
		offset = 7 - offset
	}
	n := (t.YearDay() + offset) / 7
	if n < 10 {
		b = append(b, '0')
	}
	return append(b, strconv.Itoa(n)...)
}

func appendWeekNumber(b []byte, t time.Time) []byte {
	_, n := t.ISOWeek()
	if n < 10 {
		b = append(b, '0')
	}
	return append(b, strconv.Itoa(n)...)
}

func appendDayOfYear(b []byte, t time.Time) []byte {
	n := t.YearDay()
	if n < 10 {
		b = append(b, '0', '0')
	} else if n < 100 {
		b = append(b, '0')
	}
	return append(b, strconv.Itoa(n)...)
}

type hourPadded struct {
	pad        byte
	twelveHour bool
}

func (v hourPadded) Append(b []byte, t time.Time) []byte {
	h := t.Hour()
	if v.twelveHour && h > 12 {
		h = h - 12
	}
	if v.twelveHour && h == 0 {
		h = 12
	}

	if h < 10 {
		b = append(b, v.pad)
		b = append(b, byte(h+48))
	} else {
		b = unrollTwoDigits(b, h)
	}
	return b
}

func unrollTwoDigits(b []byte, v int) []byte {
	b = append(b, byte((v/10)+48))
	b = append(b, byte((v%10)+48))
	return b
}

type hmsWAMPM struct{}

func (v hmsWAMPM) Append(b []byte, t time.Time) []byte {
	h := t.Hour()
	var am bool

	if h == 0 {
		b = append(b, '1')
		b = append(b, '2')
		am = true
	} else {
		switch {
		case h == 12:
			// no op
		case h > 12:
			h = h - 12
		default:
			am = true
		}
		b = unrollTwoDigits(b, h)
	}
	b = append(b, ':')
	b = unrollTwoDigits(b, t.Minute())
	b = append(b, ':')
	b = unrollTwoDigits(b, t.Second())

	b = append(b, ' ')
	if am {
		b = append(b, 'A')
	} else {
		b = append(b, 'P')
	}
	b = append(b, 'M')

	return b
}
