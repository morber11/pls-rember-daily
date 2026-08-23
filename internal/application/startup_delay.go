//go:build !skipstartupdelay

package application

import "time"

const startupDelay = 5 * time.Second

func waitForStartup(clock Clock) {
	clock.Sleep(startupDelay)
}
