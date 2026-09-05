package handlers

import (
	"math"
	"testing"
	"time"
)

func TestExecutivePeriodBoundsMonth(t *testing.T) {
	location := time.FixedZone("ICT", 7*60*60)
	now := time.Date(2026, time.January, 15, 10, 30, 0, 0, location)

	currentStart, currentEnd, previousStart, currentLabel, previousLabel := executivePeriodBounds(now, "month")

	assertExecutiveTime(t, currentStart, time.Date(2026, time.January, 1, 0, 0, 0, 0, location))
	assertExecutiveTime(t, currentEnd, time.Date(2026, time.February, 1, 0, 0, 0, 0, location))
	assertExecutiveTime(t, previousStart, time.Date(2025, time.December, 1, 0, 0, 0, 0, location))
	if currentLabel != "Tháng 01/2026" {
		t.Fatalf("unexpected current period label: %q", currentLabel)
	}
	if previousLabel != "Tháng 12/2025" {
		t.Fatalf("unexpected previous period label: %q", previousLabel)
	}
}

func TestExecutiveDeltaPct(t *testing.T) {
	tests := []struct {
		name     string
		current  float64
		previous float64
		want     float64
	}{
		{name: "increase", current: 120, previous: 100, want: 20},
		{name: "decrease", current: 75, previous: 100, want: -25},
		{name: "both zero", current: 0, previous: 0, want: 0},
		{name: "new activity", current: 10, previous: 0, want: 100},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := executiveDeltaPct(test.current, test.previous)
			if math.Abs(got-test.want) > 0.0001 {
				t.Fatalf("executiveDeltaPct(%v, %v) = %v, want %v", test.current, test.previous, got, test.want)
			}
		})
	}
}

func assertExecutiveTime(t *testing.T, got, want time.Time) {
	t.Helper()
	if !got.Equal(want) {
		t.Fatalf("time = %s, want %s", got, want)
	}
}
