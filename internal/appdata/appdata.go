package appdata

import (
	"fmt"
	"os"
	"path/filepath"
)

const Name = "PlsRemberDaily"

func Directory() (string, error) {
	configDirectory, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("find user configuration directory: %w", err)
	}

	return filepath.Join(configDirectory, Name), nil
}
