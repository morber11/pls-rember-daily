//go:build skipstartupdelay

package application

import "time"

const startupDelay time.Duration = 0

func waitForStartup(Clock) {}
