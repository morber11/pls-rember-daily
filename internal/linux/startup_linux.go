//go:build linux

package linux

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

func writeStartupEntry(configDirectory string, executable string) error {
	exec, err := quoteDesktopEntryArgument(executable)
	if err != nil {
		return err
	}

	autostartDirectory := filepath.Join(configDirectory, "autostart")
	if err := os.MkdirAll(autostartDirectory, 0700); err != nil {
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

	path := filepath.Join(autostartDirectory, startupFileName)
	if err := os.WriteFile(path, []byte(entry), 0600); err != nil {
		return fmt.Errorf("write Linux autostart entry: %w", err)
	}

	if err := os.Chmod(path, 0600); err != nil {
		return fmt.Errorf("set Linux autostart entry permissions: %w", err)
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
