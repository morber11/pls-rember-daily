package reminder

import "time"

// go only accepts certain formats so don't touch this completely random dateLayout
// https://pkg.go.dev/time#Time.Format
const dateLayout = "2006-01-02"

func shouldShow(lastShown string, now time.Time) bool {
	return lastShown != now.Local().Format(dateLayout)
}

func ShouldShowToday(now time.Time) (bool, error) {
	path, err := stateFilePath()
	if err != nil {
		return false, err
	}

	saved, err := loadState(path)
	if err != nil {
		return false, err
	}

	return shouldShow(saved.LastShown, now), nil
}

func MarkShownToday(now time.Time) error {
	path, err := stateFilePath()
	if err != nil {
		return err
	}

	return saveState(path, now)
}
