package application

import (
	"fmt"
	"time"

	"plsremberdaily/assets"
)

const startupDelay = 5 * time.Second

type Platform interface {
	Acquire() (release func(), acquired bool, err error)
	RegisterStartup() error
	Play(assets.Video, func() error) error
}

type Reminder interface {
	AllowMultipleRuns() bool
	ShouldShowToday(time.Time) (bool, error)
	MarkShownToday(time.Time) error
}

type Schedule interface {
	Video(time.Time) assets.Video
}

type Clock interface {
	Now() time.Time
	Sleep(time.Duration)
}

type Logger interface {
	Printf(string, ...any)
}

type Runner struct {
	platform Platform
	reminder Reminder
	schedule Schedule
	clock    Clock
	logger   Logger
}

func NewRunner(
	platform Platform,
	reminder Reminder,
	schedule Schedule,
	clock Clock,
	logger Logger,
) Runner {
	return Runner{
		platform: platform,
		reminder: reminder,
		schedule: schedule,
		clock:    clock,
		logger:   logger,
	}
}

func (runner Runner) Run() error {
	release, acquired, err := runner.platform.Acquire()
	if err != nil {
		return err
	}
	if !acquired {
		return nil
	}
	defer release()

	if err := runner.platform.RegisterStartup(); err != nil {
		runner.logger.Printf("register startup: %v", err)
	}

	if !runner.reminder.AllowMultipleRuns() {
		show, err := runner.reminder.ShouldShowToday(runner.clock.Now())
		if err != nil {
			return err
		}
		if !show {
			return nil
		}
	}

	runner.clock.Sleep(startupDelay)

	video := runner.schedule.Video(runner.clock.Now())
	return runner.platform.Play(video, func() error {
		if err := runner.reminder.MarkShownToday(runner.clock.Now()); err != nil {
			return fmt.Errorf("mark reminder as shown: %w", err)
		}

		return nil
	})
}

type SystemClock struct{}

func (SystemClock) Now() time.Time {
	return time.Now()
}

func (SystemClock) Sleep(duration time.Duration) {
	time.Sleep(duration)
}
