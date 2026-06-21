<p align="center"><img src="https://raw.githubusercontent.com/go-composites/brand/main/social/go-composites.png" alt="go-composites/date" width="720"></p>

# date

[![ci](https://github.com/go-composites/date/actions/workflows/ci.yml/badge.svg)](https://github.com/go-composites/date/actions/workflows/ci.yml)

A `Date` composite — a calendar date (year, month, day; no time-of-day, Ruby's
`Date` as the reference) for **Composition-Oriented Programming**. `Date` is
backed by a Go `time.Time` pinned to midnight UTC (used only for calendar
validation and day arithmetic) and exposes its **fallible operations as
`Result` values** — so an impossible date (the canonical example being
`FromYMD(2026, 2, 29)`) is a *value*, never a panic and never `nil`.

`Date` is deterministic by construction: there is **no `Today()`**. Instances
are built only from explicit values (`FromYMD`, `Parse`), which keeps behaviour
reproducible across architectures and test runs.

```golang
parsed := Date.Parse(value)
if parsed.HasError() {
    fmt.Println(parsed.Error().Message())
} else {
    fmt.Println(parsed.Payload().(Date.Interface).ToGoString())
}
```

`Date` follows the org's Null-Object / never-nil invariant (enforced by the
`nonnil` CI analyzer): `Date.Null()` satisfies the same `Interface` and reports
`IsNull() == true`.

## Install

```bash
export GOPRIVATE=github.com/go-composites GOPROXY=direct GOSUMDB=off
go get github.com/go-composites/date@main
```

## Usage

> [!NOTE] main.go

```golang
package main

import (
    "fmt"

    Date "github.com/go-composites/date/src"
)

func main() {
    // FromYMD is fallible: an impossible calendar date is a value, not a panic.
    fmt.Println(Date.FromYMD(2026, 2, 29).HasError()) // true (2026 is not a leap year)

    // 2024 is a leap year, so Feb 29 is valid.
    leap := Date.FromYMD(2024, 2, 29).Payload().(Date.Interface)
    fmt.Println(leap.ToGoString()) // 2024-02-29

    // Parse reads ISO "YYYY-MM-DD".
    today := Date.Parse("2026-06-21").Payload().(Date.Interface)
    fmt.Println(today.Weekday()) // Sunday

    // AddDays returns a Result; DaysBetween is a signed day count.
    later := today.AddDays(10).Payload().(Date.Interface)
    fmt.Println(today.DaysBetween(later)) // 10
}
```

```bash
$ task build
```

## API

### Date (`github.com/go-composites/date/src`, package `Date`)

Constructors

- `FromYMD(year, month, day int) Result.Interface` — a `Result` whose payload is
  a `Date`, validated by a `time.Date` round-trip, or an `Error.New(...)` when
  the components do not form a real calendar date (Feb 30, month 13, …).
- `Parse(value string) Result.Interface` — a `Result` whose payload is a `Date`
  parsed from ISO `"YYYY-MM-DD"`, or an `Error.New(...)` when `value` is
  malformed.
- `Null() Interface` — the Null-Object `Date` (`IsNull() == true`).

Accessors

- `Year() int`, `Month() int` (1–12), `Day() int` (1–31).
- `Weekday() string` — the English day name (e.g. `"Monday"`).
- `ToGoString() string` — ISO `"YYYY-MM-DD"`.

Comparisons (each returns `bool`)

- `Before(other)` / `After(other)` / `Equal(other)`.

Arithmetic

- `AddDays(n int) Result.Interface` — a `Result` whose payload is a new `Date`
  `n` days later (`n` may be negative).
- `DaysBetween(other Interface) int` — the signed number of days from the
  receiver to `other`.

Null-Object

- `IsNull() bool`.

## License

BSD-3-Clause — see [LICENSE](./LICENSE).
