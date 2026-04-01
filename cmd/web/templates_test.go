package main

import (
	"testing"
	"time"

	"github.com/Nobsmoke123/snippetbox/internal/assert"
)

func TestHumanDate(t *testing.T) {
	// Create a slice of annonymous stuct
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2024, 3, 17, 10, 15, 0, 0, time.UTC),
			want: "17 Mar 2024 at 10:15am",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2024, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Mar 2024 at 09:15am",
		},
	}

	for _, tt := range tests {
		// Use the t.Run() function to run a sub-test for each test case. The
		// First parameter to this is the name of the test (which is used to identify
		// the sub-test in any log output) and the second parameter is an annonymous
		// function containing the actual test for each case
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)

			assert.Equal(t, hd, tt.want)

			// if hd != tt.want {
			// 	t.Errorf("got %q; want %q", hd, tt.want)
			// }
		})
	}

	// Initialize a new time.Time object and pass it to the humandate function.
	// tm := time.Date(2026, 4, 1, 1, 39, 0, 0, time.UTC)
	// hd := humanDate(tm)

	// Check that the output from humanDate function is in the format
	// we expect. If it isn't what we expect, use the t.Errorf() function to
	// indicate that the test has failed and log the expected and actual values.
	// if hd != "01 Apr 2026 at 01:39am" {
	// 	t.Errorf("got %q; want %q", hd, "01 Apr 2026 at 01:39am")
	// }
}
