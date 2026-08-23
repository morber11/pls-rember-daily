//go:build allowmultipleruns

package reminder

// by default - runs once per day. while debugging, we might want to run it multiple times
func AllowMultipleRuns() bool {
	return true
}
