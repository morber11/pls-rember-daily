//go:build linux

package linux

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteStartupEntry(t *testing.T) {
	configDirectory := t.TempDir()
	executable := `/opt/Pls Rember Daily/pls-rember-daily`

	if err := writeStartupEntry(configDirectory, executable); err != nil {
		t.Fatalf("writeStartupEntry() error = %v", err)
	}

	path := filepath.Join(configDirectory, "autostart", startupFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read startup entry: %v", err)
	}

	entry := string(data)
	for _, field := range []string{
		"[Desktop Entry]\n",
		"Type=Application\n",
		"Name=Pls Rember Daily\n",
		`Exec="/opt/Pls Rember Daily/pls-rember-daily"` + "\n",
		"Terminal=false\n",
	} {
		if !strings.Contains(entry, field) {
			t.Errorf("startup entry does not contain %q", field)
		}
	}
}

func TestWriteStartupEntryEscapesExecutablePath(t *testing.T) {
	configDirectory := t.TempDir()
	executable := `/opt/a\b"c` + "`" + "$d%\te\nf\rg/app"

	if err := writeStartupEntry(configDirectory, executable); err != nil {
		t.Fatalf("writeStartupEntry() error = %v", err)
	}

	path := filepath.Join(configDirectory, "autostart", startupFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read startup entry: %v", err)
	}

	want := `Exec="/opt/a\\\\b\\"c\\` + "`" + `\\$d%%\te\nf\rg/app"` + "\n"
	if !strings.Contains(string(data), want) {
		t.Errorf("startup entry does not contain %q", want)
	}
}
