package main

import (
	"log"

	"plsremberdaily/internal/application"
	"plsremberdaily/internal/platform"
	"plsremberdaily/internal/reminder"
	"plsremberdaily/internal/schedule"
)

func main() {
	runner := application.NewRunner(
		platform.NewRuntime(),
		reminder.Service{},
		schedule.Service{},
		application.SystemClock{},
		log.Default(),
	)

	if err := runner.Run(); err != nil {
		log.Fatal(err)
	}
}
