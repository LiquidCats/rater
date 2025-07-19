package timeutils_test

import (
	"testing"
	"time"

	"github.com/LiquidCats/rater/internal/app/utils/timeutils"
)

func TestRoundToNearest(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Test rounding down (remainder < 3)
		{
			name:     "Round down from 0 minutes",
			input:    "2024-01-15 10:00:30",
			expected: "2024-01-15 10:00:00",
		},
		{
			name:     "Round down from 1 minute",
			input:    "2024-01-15 10:01:45",
			expected: "2024-01-15 10:00:00",
		},
		{
			name:     "Round down from 2 minutes",
			input:    "2024-01-15 10:02:59",
			expected: "2024-01-15 10:00:00",
		},
		{
			name:     "Round down from 6 minutes",
			input:    "2024-01-15 10:06:00",
			expected: "2024-01-15 10:05:00",
		},
		{
			name:     "Round down from 7 minutes",
			input:    "2024-01-15 10:07:30",
			expected: "2024-01-15 10:05:00",
		},

		// Test rounding up (remainder >= 3)
		{
			name:     "Round up from 3 minutes",
			input:    "2024-01-15 10:03:00",
			expected: "2024-01-15 10:05:00",
		},
		{
			name:     "Round up from 4 minutes",
			input:    "2024-01-15 10:04:59",
			expected: "2024-01-15 10:05:00",
		},
		{
			name:     "Round up from 8 minutes",
			input:    "2024-01-15 10:08:00",
			expected: "2024-01-15 10:10:00",
		},
		{
			name:     "Round up from 9 minutes",
			input:    "2024-01-15 10:09:45",
			expected: "2024-01-15 10:10:00",
		},

		// Test exact 5-minute intervals
		{
			name:     "Exact 5 minutes",
			input:    "2024-01-15 10:05:00",
			expected: "2024-01-15 10:05:00",
		},
		{
			name:     "Exact 10 minutes",
			input:    "2024-01-15 10:10:00",
			expected: "2024-01-15 10:10:00",
		},
		{
			name:     "Exact 15 minutes",
			input:    "2024-01-15 10:15:00",
			expected: "2024-01-15 10:15:00",
		},

		// Test hour boundaries
		{
			name:     "Round up crossing hour boundary",
			input:    "2024-01-15 10:58:00",
			expected: "2024-01-15 11:00:00",
		},
		{
			name:     "Round up crossing hour boundary at 59 minutes",
			input:    "2024-01-15 10:59:30",
			expected: "2024-01-15 11:00:00",
		},
		{
			name:     "Round down near hour boundary",
			input:    "2024-01-15 10:57:00",
			expected: "2024-01-15 10:55:00",
		},

		// Test midnight boundary
		{
			name:     "Round up crossing midnight",
			input:    "2024-01-15 23:58:00",
			expected: "2024-01-16 00:00:00",
		},
		{
			name:     "Round up crossing midnight at 59 minutes",
			input:    "2024-01-15 23:59:00",
			expected: "2024-01-16 00:00:00",
		},
		{
			name:     "Round down before midnight",
			input:    "2024-01-15 23:57:00",
			expected: "2024-01-15 23:55:00",
		},

		// Test start of day
		{
			name:     "Start of day",
			input:    "2024-01-15 00:00:00",
			expected: "2024-01-15 00:00:00",
		},
		{
			name:     "Round up from 3 minutes after midnight",
			input:    "2024-01-15 00:03:30",
			expected: "2024-01-15 00:05:00",
		},

		// Test all possible minute remainders when divided by 5
		{
			name:     "Minute ending in 0",
			input:    "2024-01-15 14:20:45",
			expected: "2024-01-15 14:20:00",
		},
		{
			name:     "Minute ending in 1",
			input:    "2024-01-15 14:21:30",
			expected: "2024-01-15 14:20:00",
		},
		{
			name:     "Minute ending in 2",
			input:    "2024-01-15 14:22:15",
			expected: "2024-01-15 14:20:00",
		},
		{
			name:     "Minute ending in 3",
			input:    "2024-01-15 14:23:00",
			expected: "2024-01-15 14:25:00",
		},
		{
			name:     "Minute ending in 4",
			input:    "2024-01-15 14:24:45",
			expected: "2024-01-15 14:25:00",
		},

		// Test different hours throughout the day
		{
			name:     "Early morning",
			input:    "2024-01-15 03:13:00",
			expected: "2024-01-15 03:15:00",
		},
		{
			name:     "Noon",
			input:    "2024-01-15 12:02:00",
			expected: "2024-01-15 12:00:00",
		},
		{
			name:     "Evening",
			input:    "2024-01-15 18:48:00",
			expected: "2024-01-15 18:50:00",
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse input time
			layout := "2006-01-02 15:04:05"
			inputTime, err := time.Parse(layout, tt.input)
			if err != nil {
				t.Fatalf("Failed to parse input time: %v", err)
			}

			// Parse expected time
			expectedTime, err := time.Parse(layout, tt.expected)
			if err != nil {
				t.Fatalf("Failed to parse expected time: %v", err)
			}

			// Call the function
			result := timeutils.RoundToNearest(inputTime, timeutils.FiveMinuteBucket)

			// Compare results
			if !result.Equal(expectedTime) {
				t.Errorf("RoundToNearest(%s) = %s; want %s",
					tt.input,
					result.Format(layout),
					tt.expected)
			}
		})
	}
}

func TestRoundToNearestWithTimezones(t *testing.T) {
	// Test with different timezones
	timezones := []string{
		"America/New_York",
		"Europe/London",
		"Asia/Tokyo",
		"Australia/Sydney",
		"UTC",
	}

	for _, tzName := range timezones {
		t.Run("Timezone_"+tzName, func(t *testing.T) {
			loc, err := time.LoadLocation(tzName)
			if err != nil {
				t.Skipf("Skipping timezone %s: %v", tzName, err)
			}

			// Create a time in the specific timezone
			testTime := time.Date(2024, 1, 15, 10, 13, 30, 0, loc)
			expectedTime := time.Date(2024, 1, 15, 10, 15, 0, 0, loc)

			result := timeutils.RoundToNearest(testTime, timeutils.FiveMinuteBucket)

			// Verify the timezone is preserved
			if result.Location() != loc {
				t.Errorf("Timezone not preserved. Got %v, want %v",
					result.Location(), loc)
			}

			// Verify the rounding is correct
			if !result.Equal(expectedTime) {
				t.Errorf("RoundToNearest(%s) = %s; want %s",
					testTime.Format("15:04:05"),
					result.Format("15:04:05"),
					expectedTime.Format("15:04:05"))
			}
		})
	}
}

func TestRoundToNearestEdgeCases(t *testing.T) {
	t.Run("Leap year boundary", func(t *testing.T) {
		// Test rounding across Feb 29 in a leap year
		input := time.Date(2024, 2, 29, 23, 58, 0, 0, time.UTC)
		expected := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
		result := timeutils.RoundToNearest(input, timeutils.FiveMinuteBucket)

		if !result.Equal(expected) {
			t.Errorf("Failed leap year test: got %s, want %s",
				result.Format("2006-01-02 15:04:05"),
				expected.Format("2006-01-02 15:04:05"))
		}
	})

	t.Run("Year boundary", func(t *testing.T) {
		// Test rounding across year boundary
		input := time.Date(2023, 12, 31, 23, 59, 0, 0, time.UTC)
		expected := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		result := timeutils.RoundToNearest(input, timeutils.FiveMinuteBucket)

		if !result.Equal(expected) {
			t.Errorf("Failed year boundary test: got %s, want %s",
				result.Format("2006-01-02 15:04:05"),
				expected.Format("2006-01-02 15:04:05"))
		}
	})
}

func BenchmarkRoundToNearest(b *testing.B) {
	testTime := time.Now()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = timeutils.RoundToNearest(testTime, timeutils.FiveMinuteBucket)
	}
}
