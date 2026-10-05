//go:build windows

package windows

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"

	"plsremberdaily/internal/appdata"
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

	appDirectory, err := appdata.Directory()
	if err != nil {
		return err
	}

	installed, err := installExecutable(appDirectory, executable)
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

func installExecutable(appDirectory string, executable string) (string, error) {
	installed := filepath.Join(appDirectory, appName+".exe")

	if strings.EqualFold(executable, installed) {
		return installed, nil
	}

	data, err := os.ReadFile(executable)
	if err != nil {
		return "", fmt.Errorf("read running executable: %w", err)
	}

	if err := os.MkdirAll(appDirectory, filemode.OwnerReadWriteExecute); err != nil {
		return "", fmt.Errorf("create install directory: %w", err)
	}

	if err := os.WriteFile(installed, data, filemode.OwnerReadWriteExecute); err != nil {
		return "", fmt.Errorf("copy executable to install directory: %w", err)
	}

	return installed, nil
}
