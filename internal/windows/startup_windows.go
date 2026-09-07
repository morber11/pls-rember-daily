//go:build windows

package windows

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/sys/windows/registry"
)

const (
	runKey      = `Software\Microsoft\Windows\CurrentVersion\Run`
	appName     = "PlsRemberDaily"
	taskName    = `\PlsRemberDaily`
	schtasksExe = "schtasks.exe"
)

func Register() error {
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not find executable path: %w", err)
	}

	taskCommand := fmt.Sprintf(`"%s"`, executable)
	command := exec.Command(
		schtasksExe,
		"/Create",
		"/TN", taskName,
		"/TR", taskCommand,
		"/SC", "ONLOGON",
		"/IT",
		"/RL", "LIMITED",
		"/F",
	)
	output, err := command.CombinedOutput()
	if err != nil {
		return fmt.Errorf("register Windows logon task: %w: %s", err, strings.TrimSpace(string(output)))
	}

	key, _, err := registry.CreateKey(registry.CURRENT_USER, runKey, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("opening Windows startup key for cleanup: %w", err)
	}
	defer key.Close()

	if err := key.DeleteValue(appName); err != nil && !errors.Is(err, registry.ErrNotExist) {
		return fmt.Errorf("remove legacy Windows startup registration: %w", err)
	}

	return nil
}
