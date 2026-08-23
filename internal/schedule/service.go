package schedule

import (
	"time"

	"plsremberdaily/assets"
)

type Service struct{}

func (Service) Video(now time.Time) assets.Video {
	return Video(now)
}
