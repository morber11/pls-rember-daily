package reminder

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

const (
	appDataDir = "PlsRemberDaily"
	stateFile  = "state.json"
)

type state struct {
	LastShown string `json:"lastShown"`
}

func stateFilePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("find user config directory: %w", err)
	}

	return filepath.Join(configDir, appDataDir, stateFile), nil
}

func loadState(path string) (state, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return state{}, nil
	}

	if err != nil {
		return state{}, fmt.Errorf("read reminder state: %w", err)
	}

	var saved state
	if err := json.Unmarshal(data, &saved); err != nil {
		return state{}, fmt.Errorf("decode reminder state: %w", err)
	}

	return saved, nil
}

func saveState(path string, now time.Time) error {
	data, err := json.Marshal(state{LastShown: now.Local().Format(dateLayout)})
	if err != nil {
		return fmt.Errorf("encode reminder state: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create reminder state directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write reminder state: %w", err)
	}

	return nil
}
