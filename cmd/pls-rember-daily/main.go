package main

import (
	"log"
	"time"

	"plsremberdaily/internal/instance"
	"plsremberdaily/internal/player"
	"plsremberdaily/internal/reminder"
	"plsremberdaily/internal/schedule"
	"plsremberdaily/internal/startup"
)

const startupDelay = 5 * time.Second

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	mutex, acquired, err := instance.Acquire()

	if err != nil {
		return err
	}

	if !acquired {
		return nil
	}

	defer instance.Release(mutex)

	if err := startup.Register(); err != nil {
		log.Printf("register Windows startup: %v", err)
	}

	if !reminder.AllowMultipleRuns() {
		show, err := reminder.ShouldShowToday(time.Now())
		if err != nil {
			return err
		}
		if !show {
			return nil
		}
	}

	time.Sleep(startupDelay)

	video := schedule.Video(time.Now())
	if err := player.Play(video, func() error {
		return reminder.MarkShownToday(time.Now())
	}); err != nil {
		return err
	}

	return nil
}
