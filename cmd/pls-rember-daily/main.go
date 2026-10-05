package main

import (
	"flag"
	"log"

	"plsremberdaily/internal/application"
	"plsremberdaily/internal/platform"
	"plsremberdaily/internal/reminder"
	"plsremberdaily/internal/schedule"
)

func main() {
	uninstall := flag.Bool("uninstall", false, "remove startup registration and app data, then exit")
	flag.Parse()

	if *uninstall {
		if err := platform.Uninstall(); err != nil {
			log.Fatal(err)
		}

		return
	}

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
