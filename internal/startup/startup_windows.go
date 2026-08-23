//go:build windows

package startup

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows/registry"
)

const (
	runKey  = `Software\Microsoft\Windows\CurrentVersion\Run`
	appName = "PlsRemberDaily"
)

func Register() error {
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not find executable path: %w", err)
	}

	key, _, err := registry.CreateKey(registry.CURRENT_USER, runKey, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("opening Windows startup key: %w", err)
	}
	defer key.Close()

	if err := key.SetStringValue(appName, fmt.Sprintf(`"%s"`, executable)); err != nil {
		return fmt.Errorf("register Windows startup: %w", err)
	}

	return nil
}
