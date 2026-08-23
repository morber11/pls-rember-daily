//go:build !allowmultipleruns

package reminder

// by default - runs once per day
func AllowMultipleRuns() bool {
	return false
}
