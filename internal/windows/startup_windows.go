//go:build windows

package windows

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"

	"plsremberdaily/internal/filemode"
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

	configDirectory, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("could not find user configuration directory: %w", err)
	}

	installed, err := installExecutable(configDirectory, executable)
	if err != nil {
		return err
	}

	key, _, err := registry.CreateKey(registry.CURRENT_USER, runKey, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("opening Windows startup key: %w", err)
	}
	defer key.Close()

	if err := key.SetStringValue(appName, fmt.Sprintf(`"%s"`, installed)); err != nil {
		return fmt.Errorf("register Windows startup: %w", err)
	}

	return nil
}

func installExecutable(configDirectory string, executable string) (string, error) {
	installed := filepath.Join(configDirectory, appName, appName+".exe")

	if strings.EqualFold(executable, installed) {
		return installed, nil
	}

	data, err := os.ReadFile(executable)
	if err != nil {
		return "", fmt.Errorf("read running executable: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(installed), filemode.OwnerReadWriteExecute); err != nil {
		return "", fmt.Errorf("create install directory: %w", err)
	}

	if err := os.WriteFile(installed, data, filemode.OwnerReadWriteExecute); err != nil {
		return "", fmt.Errorf("copy executable to install directory: %w", err)
	}

	return installed, nil
}
