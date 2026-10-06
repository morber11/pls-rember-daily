//go:build windows

package windows

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"

	"plsremberdaily/internal/appdata"
	"plsremberdaily/internal/filemode"
)

const runKey = `Software\Microsoft\Windows\CurrentVersion\Run`

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

	if err := key.SetStringValue(appdata.Name, fmt.Sprintf(`"%s"`, installed)); err != nil {
		return fmt.Errorf("register Windows startup: %w", err)
	}

	return nil
}

func Uninstall() error {
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not find executable path: %w", err)
	}

	appDirectory, err := appdata.Directory()
	if err != nil {
		return err
	}

	key, err := registry.OpenKey(registry.CURRENT_USER, runKey, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("opening Windows startup key: %w", err)
	}
	defer key.Close()

	if err := key.DeleteValue(appdata.Name); err != nil && !errors.Is(err, registry.ErrNotExist) {
		return fmt.Errorf("remove Windows startup: %w", err)
	}

	return removeAppDirectory(appDirectory, executable)
}

func installedExecutablePath(appDirectory string) string {
	return filepath.Join(appDirectory, appdata.Name+".exe")
}

func installExecutable(appDirectory string, executable string) (string, error) {
	installed := installedExecutablePath(appDirectory)

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

func removeAppDirectory(appDirectory string, executable string) error {
	installed := installedExecutablePath(appDirectory)

	// windows cannot delete a running executable but it can move/rename it for some reason
	if strings.EqualFold(executable, installed) {
		if err := os.Rename(installed, filepath.Join(os.TempDir(), appdata.Name+".exe")); err != nil {
			return fmt.Errorf("move running executable out of app directory: %w", err)
		}
	}

	if err := os.RemoveAll(appDirectory); err != nil {
		return fmt.Errorf("remove app directory: %w", err)
	}

	return nil
}
