package reminder

import (
	"testing"
	"time"
)

func TestShouldShow(t *testing.T) {
	now := time.Date(2026, time.August, 11, 9, 30, 0, 0, time.Local)

	tests := []struct {
		name      string
		lastShown string
		want      bool
	}{
		{name: "today", lastShown: "2026-08-11", want: false},
		{name: "previous day", lastShown: "2026-08-10", want: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := shouldShow(test.lastShown, now); got != test.want {
				t.Errorf("shouldShow(%q, %v) = %t, want %t", test.lastShown, now, got, test.want)
			}
		})
	}
}
