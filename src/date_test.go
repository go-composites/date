package Date_test

import (
	Date "github.com/go-composites/date/src"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Date", func() {

	ginkgo.Describe("constructors", func() {
		ginkgo.Describe("FromYMD", func() {
			ginkgo.It("builds a valid calendar date", func() {
				r := Date.FromYMD(2024, 2, 29)
				gomega.Expect(r.HasError()).To(gomega.BeFalse())
				d := r.Payload().(Date.Interface)
				gomega.Expect(d.Year()).To(gomega.Equal(2024))
				gomega.Expect(d.Month()).To(gomega.Equal(2))
				gomega.Expect(d.Day()).To(gomega.Equal(29))
				gomega.Expect(d.IsNull()).To(gomega.BeFalse())
			})
			ginkgo.It("rejects a non-existent leap day", func() {
				r := Date.FromYMD(2026, 2, 29)
				gomega.Expect(r.HasError()).To(gomega.BeTrue())
				gomega.Expect(r.Error().Message()).NotTo(gomega.BeEmpty())
			})
			ginkgo.It("rejects an out-of-range month", func() {
				r := Date.FromYMD(2026, 13, 1)
				gomega.Expect(r.HasError()).To(gomega.BeTrue())
			})
			ginkgo.It("rejects an out-of-range day", func() {
				r := Date.FromYMD(2026, 4, 31)
				gomega.Expect(r.HasError()).To(gomega.BeTrue())
			})
		})

		ginkgo.Describe("Parse", func() {
			ginkgo.It("parses a well-formed ISO value", func() {
				r := Date.Parse("2026-06-21")
				gomega.Expect(r.HasError()).To(gomega.BeFalse())
				d := r.Payload().(Date.Interface)
				gomega.Expect(d.ToGoString()).To(gomega.Equal("2026-06-21"))
			})
			ginkgo.It("returns a Result carrying an error on a malformed value", func() {
				r := Date.Parse("not-a-date")
				gomega.Expect(r.HasError()).To(gomega.BeTrue())
				gomega.Expect(r.Error().Message()).NotTo(gomega.BeEmpty())
			})
		})
	})

	ginkgo.Describe("accessors", func() {
		var d = Date.FromYMD(2026, 6, 21).Payload().(Date.Interface)

		ginkgo.It("reports the year, month and day", func() {
			gomega.Expect(d.Year()).To(gomega.Equal(2026))
			gomega.Expect(d.Month()).To(gomega.Equal(6))
			gomega.Expect(d.Day()).To(gomega.Equal(21))
		})
		ginkgo.It("names the weekday", func() {
			gomega.Expect(d.Weekday()).To(gomega.Equal("Sunday"))
		})
		ginkgo.It("renders ISO from ToGoString", func() {
			gomega.Expect(d.ToGoString()).To(gomega.Equal("2026-06-21"))
		})
	})

	ginkgo.Describe("comparisons", func() {
		var early = Date.FromYMD(2026, 1, 1).Payload().(Date.Interface)
		var late = Date.FromYMD(2026, 12, 31).Payload().(Date.Interface)

		ginkgo.It("reports before", func() {
			gomega.Expect(early.Before(late)).To(gomega.BeTrue())
			gomega.Expect(late.Before(early)).To(gomega.BeFalse())
		})
		ginkgo.It("reports after", func() {
			gomega.Expect(late.After(early)).To(gomega.BeTrue())
			gomega.Expect(early.After(late)).To(gomega.BeFalse())
		})
		ginkgo.It("reports equality", func() {
			gomega.Expect(early.Equal(Date.FromYMD(2026, 1, 1).Payload().(Date.Interface))).To(gomega.BeTrue())
			gomega.Expect(early.Equal(late)).To(gomega.BeFalse())
		})
	})

	ginkgo.Describe("arithmetic", func() {
		var d = Date.FromYMD(2026, 6, 21).Payload().(Date.Interface)

		ginkgo.It("AddDays shifts a Date forward", func() {
			r := d.AddDays(10)
			gomega.Expect(r.HasError()).To(gomega.BeFalse())
			gomega.Expect(r.Payload().(Date.Interface).ToGoString()).To(gomega.Equal("2026-07-01"))
		})
		ginkgo.It("AddDays shifts a Date backward with a negative count", func() {
			r := d.AddDays(-21)
			gomega.Expect(r.HasError()).To(gomega.BeFalse())
			gomega.Expect(r.Payload().(Date.Interface).ToGoString()).To(gomega.Equal("2026-05-31"))
		})
		ginkgo.It("DaysBetween reports a positive gap to a later date", func() {
			later := Date.FromYMD(2026, 7, 1).Payload().(Date.Interface)
			gomega.Expect(d.DaysBetween(later)).To(gomega.Equal(10))
		})
		ginkgo.It("DaysBetween reports a negative gap to an earlier date", func() {
			earlier := Date.FromYMD(2026, 6, 11).Payload().(Date.Interface)
			gomega.Expect(d.DaysBetween(earlier)).To(gomega.Equal(-10))
		})
	})

	ginkgo.Describe("the Null-Object variant", func() {
		var n = Date.Null()

		ginkgo.It("satisfies the Date interface", func() {
			var _ Date.Interface = n
		})
		ginkgo.It("reports IsNull() true", func() {
			gomega.Expect(n.IsNull()).To(gomega.BeTrue())
		})
		ginkgo.It("reports zero components", func() {
			gomega.Expect(n.Year()).To(gomega.Equal(0))
			gomega.Expect(n.Month()).To(gomega.Equal(0))
			gomega.Expect(n.Day()).To(gomega.Equal(0))
		})
		ginkgo.It("Weekday is the empty string", func() {
			gomega.Expect(n.Weekday()).To(gomega.Equal(``))
		})
		ginkgo.It("ToGoString is the empty string", func() {
			gomega.Expect(n.ToGoString()).To(gomega.Equal(``))
		})
		ginkgo.It("Before is always false", func() {
			gomega.Expect(n.Before(Date.FromYMD(2026, 1, 1).Payload().(Date.Interface))).To(gomega.BeFalse())
		})
		ginkgo.It("After is always false", func() {
			gomega.Expect(n.After(Date.FromYMD(2026, 1, 1).Payload().(Date.Interface))).To(gomega.BeFalse())
		})
		ginkgo.It("Equal is true only against another null", func() {
			gomega.Expect(n.Equal(Date.Null())).To(gomega.BeTrue())
			gomega.Expect(n.Equal(Date.FromYMD(2026, 1, 1).Payload().(Date.Interface))).To(gomega.BeFalse())
		})
		ginkgo.It("AddDays returns a Result whose payload is a null Date", func() {
			r := n.AddDays(5)
			gomega.Expect(r.HasError()).To(gomega.BeFalse())
			gomega.Expect(r.Payload().(Date.Interface).IsNull()).To(gomega.BeTrue())
		})
		ginkgo.It("DaysBetween is zero", func() {
			gomega.Expect(n.DaysBetween(Date.FromYMD(2026, 1, 1).Payload().(Date.Interface))).To(gomega.Equal(0))
		})
	})
})
