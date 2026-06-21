package Date

import (
	"time"

	Error "github.com/go-composites/error/src"
	Result "github.com/go-composites/result/src"
)

/*
Date is a composite over a calendar date — a year, month and day with no
time-of-day component (Ruby's Date as the reference).

It follows Composition-Oriented Programming: it is interface-first, its fallible
constructors (FromYMD, Parse, AddDays) return a Result.Interface so that
failures are values rather than panics, and it ships a Null-Object variant so
callers never have to test for a bare nil.

Date is deterministic by construction: there is no Today() in any covered path.
Instances are built only from explicit values (FromYMD, Parse), which keeps
behaviour reproducible across architectures and test runs.

Internally a Date is a Go time.Time pinned to midnight UTC; the time-of-day is
never exposed and is used only for calendar validation and day arithmetic.
*/
type Interface interface {
	Year() int
	Month() int
	Day() int
	Weekday() string
	ToGoString() string
	Before(Interface) bool
	After(Interface) bool
	Equal(Interface) bool
	AddDays(n int) Result.Interface
	DaysBetween(Interface) int
	IsNull() bool
}

// isoLayout is the ISO calendar-date layout (YYYY-MM-DD).
const isoLayout = "2006-01-02"

type data struct {
	value time.Time
}

/*
FromYMD is the Date constructor from explicit year, month and day values.

The calendar date is validated by a time.Date round-trip: a date that
normalises to different components (e.g. Feb 30, month 13) is rejected. On
success the Result carries the Date as its payload; otherwise it carries an
Error — the operation never panics and never returns nil.

	r := Date.FromYMD(2024, 2, 29) // valid: 2024 is a leap year
	if r.HasError() {
	    // r.Error().Message()
	}
*/
func FromYMD(year, month, day int) Result.Interface {
	value := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if value.Year() != year ||
		int(value.Month()) != month ||
		value.Day() != day {
		return Result.New(
			Result.WithError(
				Error.New(`invalid calendar date`),
			),
		)
	}
	return Result.New(
		Result.WithPayload(fromGo(value)),
	)
}

/*
fromGo wraps a Go time.Time, pinned to midnight UTC, into a Date composite.
*/
func fromGo(value time.Time) Interface {
	midnight := time.Date(
		value.Year(), value.Month(), value.Day(),
		0, 0, 0, 0, time.UTC,
	)
	return &data{value: midnight}
}

/*
Parse builds a Date from an ISO "YYYY-MM-DD" value.

On success the Result carries the Date as its payload; when the value does not
match the ISO layout the Result carries an Error instead — the operation never
panics and never returns nil.

	r := Date.Parse("2026-06-21")
	if r.HasError() {
	    // r.Error().Message()
	}
*/
func Parse(value string) Result.Interface {
	parsed, err := time.Parse(isoLayout, value)
	if err != nil {
		return Result.New(
			Result.WithError(
				Error.New(err.Error()),
			),
		)
	}
	return Result.New(
		Result.WithPayload(fromGo(parsed)),
	)
}

/*
Year returns the calendar year.
*/
func (d data) Year() int {
	return d.value.Year()
}

/*
Month returns the calendar month (1–12).
*/
func (d data) Month() int {
	return int(d.value.Month())
}

/*
Day returns the day of the month (1–31).
*/
func (d data) Day() int {
	return d.value.Day()
}

/*
Weekday returns the English name of the day of the week (e.g. "Monday").
*/
func (d data) Weekday() string {
	return d.value.Weekday().String()
}

/*
ToGoString returns the ISO "YYYY-MM-DD" representation of the date.
*/
func (d data) ToGoString() string {
	return d.value.Format(isoLayout)
}

/*
Before reports whether the receiver falls strictly before other.
*/
func (d data) Before(other Interface) bool {
	return d.value.Before(toMidnight(other))
}

/*
After reports whether the receiver falls strictly after other.
*/
func (d data) After(other Interface) bool {
	return d.value.After(toMidnight(other))
}

/*
Equal reports whether the receiver and other denote the same calendar date.
*/
func (d data) Equal(other Interface) bool {
	return d.value.Equal(toMidnight(other))
}

/*
AddDays returns a Result whose payload is a new Date n days later (n may be
negative). The operation never panics and never returns nil.
*/
func (d data) AddDays(n int) Result.Interface {
	shifted := d.value.AddDate(0, 0, n)
	return Result.New(
		Result.WithPayload(fromGo(shifted)),
	)
}

/*
DaysBetween returns the signed number of days from the receiver to other
(positive when other is later, negative when earlier).
*/
func (d data) DaysBetween(other Interface) int {
	delta := toMidnight(other).Sub(d.value)
	return int(delta.Hours() / 24)
}

/*
IsNull reports whether the Date is the Null-Object variant.

A concrete Date is never null.
*/
func (d data) IsNull() bool {
	return false
}

/*
toMidnight projects an Interface back onto a midnight-UTC time.Time so that
comparisons and arithmetic are purely calendar-based.
*/
func toMidnight(other Interface) time.Time {
	return time.Date(
		other.Year(), time.Month(other.Month()), other.Day(),
		0, 0, 0, 0, time.UTC,
	)
}

type null struct{}

/*
Null returns the Null-Object variant of Date.

It satisfies Interface so callers never have to test for a bare nil: its
components are zero, its string form is empty, its comparisons are false, its
AddDays yields a Result wrapping another null Date, DaysBetween is zero, and
IsNull() returns true.
*/
func Null() Interface {
	return &null{}
}

func (n null) Year() int {
	return 0
}

func (n null) Month() int {
	return 0
}

func (n null) Day() int {
	return 0
}

func (n null) Weekday() string {
	return ``
}

func (n null) ToGoString() string {
	return ``
}

func (n null) Before(Interface) bool {
	return false
}

func (n null) After(Interface) bool {
	return false
}

func (n null) Equal(other Interface) bool {
	return other.IsNull()
}

func (n null) AddDays(int) Result.Interface {
	return Result.New(
		Result.WithPayload(Null()),
	)
}

func (n null) DaysBetween(Interface) int {
	return 0
}

func (n null) IsNull() bool {
	return true
}
