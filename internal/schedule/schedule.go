package schedule

import (
	"math/rand/v2"
	"time"

	"plsremberdaily/assets"
)

const anyWeekday time.Weekday = -1

type scheduledVideo struct {
	month   time.Month
	day     int
	weekday time.Weekday
	video   assets.Video
}

var scheduledVideos = []scheduledVideo{
	{month: time.September, day: 21, weekday: anyWeekday, video: assets.SemberVideo()}, // sember
	{day: 13, weekday: time.Friday, video: assets.EvilVideo()},                         // fri 13th - EVIL
	{day: 9, month: time.September, video: assets.CirnoDayVideo()},                     // cirno day
}

func Video(now time.Time) assets.Video {
	for _, scheduled := range scheduledVideos {
		if now.Day() == scheduled.day &&
			(scheduled.month == 0 || now.Month() == scheduled.month) &&
			(scheduled.weekday == anyWeekday || now.Weekday() == scheduled.weekday) {
			return scheduled.video
		}
	}

	return selectFallbackVideo()
}

func selectFallbackVideo() assets.Video {
	if rand.IntN(100) == 0 {
		return assets.CoolAnimation()
	}

	return assets.GetDefaultVideo()
}
