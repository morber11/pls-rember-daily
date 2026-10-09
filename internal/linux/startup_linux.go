//go:build linux

package linux

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"plsremberdaily/internal/appdata"
	"plsremberdaily/internal/filemode"
)

const startupFileName = applicationID + ".desktop"

func Register() error {
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not find executable path: %w", err)
	}

	configDirectory, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("could not find user configuration directory: %w", err)
	}

	return writeStartupEntry(configDirectory, executable)
}

func Uninstall() error {
	configDirectory, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("could not find user configuration directory: %w", err)
	}

	appDirectory, err := appdata.Directory()
	if err != nil {
		return err
	}

	return removeStartupEntryAndAppData(configDirectory, appDirectory)
}

func startupEntryPath(configDirectory string) string {
	return filepath.Join(configDirectory, "autostart", startupFileName)
}

func writeStartupEntry(configDirectory string, executable string) error {
	exec, err := quoteDesktopEntryArgument(executable)
	if err != nil {
		return err
	}

	path := startupEntryPath(configDirectory)

	if err := os.MkdirAll(filepath.Dir(path), filemode.OwnerReadWriteExecute); err != nil {
		return fmt.Errorf("create Linux autostart directory: %w", err)
	}

	entry := strings.Join([]string{
		"[Desktop Entry]",
		"Type=Application",
		"Name=Pls Rember Daily",
		"Exec=" + exec,
		"Terminal=false",
		"",
	}, "\n")

	if err := os.WriteFile(path, []byte(entry), filemode.OwnerReadWrite); err != nil {
		return fmt.Errorf("write Linux autostart entry: %w", err)
	}

	if err := os.Chmod(path, filemode.OwnerReadWrite); err != nil {
		return fmt.Errorf("set Linux autostart entry permissions: %w", err)
	}

	return nil
}

func removeStartupEntryAndAppData(configDirectory string, appDirectory string) error {
	if err := os.Remove(startupEntryPath(configDirectory)); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("remove Linux autostart entry: %w", err)
	}

	if err := os.RemoveAll(appDirectory); err != nil {
		return fmt.Errorf("remove app directory: %w", err)
	}

	return nil
}

func quoteDesktopEntryArgument(argument string) (string, error) {
	if strings.Contains(argument, "=") {
		return "", fmt.Errorf("Linux autostart executable path contains an equals sign")
	}

	// desktop entry values are unescaped before the Exec command is parsed
	quotedArgument := strings.NewReplacer(
		`\`, `\\`,
		`"`, `\"`,
		"`", "\\`",
		`$`, `\$`,
		`%`, `%%`,
	).Replace(argument)

	desktopEntryValue := strings.NewReplacer(
		`\`, `\\`,
		"\n", `\n`,
		"\r", `\r`,
		"\t", `\t`,
	).Replace(quotedArgument)

	return `"` + desktopEntryValue + `"`, nil
}
