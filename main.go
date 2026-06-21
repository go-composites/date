package main

import (
	"fmt"

	Date "github.com/go-composites/date/src"
	Result "github.com/go-composites/result/src"
)

func main() {
	// FromYMD is fallible: an impossible calendar date is a value, not a panic.
	feb29NonLeap := Date.FromYMD(2026, 2, 29)
	fmt.Println("2026-02-29 has error:", feb29NonLeap.HasError())
	fmt.Println(feb29NonLeap.Error().Message())

	// 2024 is a leap year, so Feb 29 is a valid calendar date.
	leap := Date.FromYMD(2024, 2, 29)
	fmt.Println("2024-02-29:", payload(leap).ToGoString())

	// Parse reads ISO "YYYY-MM-DD".
	today := payload(Date.Parse("2026-06-21"))
	fmt.Println("parsed:", today.ToGoString())
	fmt.Println("weekday:", today.Weekday())

	// AddDays returns a Result — shifting a date forward (or backward).
	later := payload(today.AddDays(10))
	fmt.Println("10 days later:", later.ToGoString())

	// DaysBetween is a signed day count.
	fmt.Println("days between:", today.DaysBetween(later))

	// The Null-Object variant is never a bare nil.
	fmt.Println("null date is null:", Date.Null().IsNull())
}

func payload(result Result.Interface) Date.Interface {
	return result.Payload().(Date.Interface)
}
