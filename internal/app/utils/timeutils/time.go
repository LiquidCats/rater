package timeutils

import "time"

const (
	FiveMinuteBucket int = 5
)

// RoundToNearest rounds a time to the nearest 5-minute interval.
func RoundToNearest(t time.Time, b int) time.Time {
	// Get the minute component
	minute := t.Minute()

	// Calculate the remainder when divided by 5
	remainder := minute % b

	// Determine whether to round up or down
	var roundedMinute int
	if remainder < 3 { //nolint:mnd
		// Round down
		roundedMinute = minute - remainder
	} else {
		// Round up
		roundedMinute = minute + (5 - remainder) //nolint:mnd
	}

	// Create a new time with the rounded minutes
	rounded := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), roundedMinute, 0, 0, t.Location())

	// Handle the case where rounding up goes to the next hour
	if roundedMinute >= 60 { //nolint:mnd
		rounded = rounded.Add(-time.Duration(roundedMinute-60) * time.Minute) //nolint:mnd
	}

	return rounded
}
