package schedule

import (
	"testing"
	"time"
)

func TestVideo(t *testing.T) {
	tests := []struct {
		name     string
		now      time.Time
		wantName string
	}{
		{
			name:     "Sember day",
			now:      time.Date(2026, time.September, 21, 9, 30, 0, 0, time.Local),
			wantName: "sember.mp4",
		},
		{
			name:     "September 9",
			now:      time.Date(2026, time.September, 9, 9, 30, 0, 0, time.Local),
			wantName: "cirno-whatsapp 2.mp4",
		},
		{
			name:     "Friday the thirteenth",
			now:      time.Date(2026, time.February, 13, 9, 30, 0, 0, time.Local),
			wantName: "evil.webm",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Video(test.now); got.Name != test.wantName {
				t.Errorf("Video(%v) = %q, want %q", test.now, got.Name, test.wantName)
			}
		})
	}
}
