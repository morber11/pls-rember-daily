package reminder

import "time"

type Service struct{}

func (Service) AllowMultipleRuns() bool {
	return AllowMultipleRuns()
}

func (Service) ShouldShowToday(now time.Time) (bool, error) {
	return ShouldShowToday(now)
}

func (Service) MarkShownToday(now time.Time) error {
	return MarkShownToday(now)
}
