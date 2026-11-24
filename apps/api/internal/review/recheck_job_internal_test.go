package review

import (
	"testing"
	"time"
)

func TestQuarterly(t *testing.T) {
	for in, want := range map[string]string{
		"2026-01-01T00:00:00Z": "2026-04-01", "2026-02-15T10:00:00Z": "2026-04-01",
		"2026-09-30T23:59:59Z": "2026-10-01", "2026-11-24T12:00:00Z": "2027-01-01",
	} {
		at, _ := time.Parse(time.RFC3339, in)
		if got := (quarterly{}).Next(at).Format(time.DateOnly); got != want {
			t.Errorf("Next(%s) = %s, want %s", in, got, want)
		}
	}
}
