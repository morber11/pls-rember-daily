//go:build linux

package linux

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"plsremberdaily/internal/filemode"
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

func TestRemoveStartupEntryAndAppData(t *testing.T) {
	configDirectory := t.TempDir()
	if err := writeStartupEntry(configDirectory, "/opt/pls-rember-daily"); err != nil {
		t.Fatalf("writeStartupEntry() error = %v", err)
	}

	appDirectory := filepath.Join(configDirectory, "PlsRemberDaily")
	if err := os.MkdirAll(appDirectory, filemode.OwnerReadWriteExecute); err != nil {
		t.Fatalf("create app directory: %v", err)
	}

	if err := os.WriteFile(filepath.Join(appDirectory, "state.json"), nil, filemode.OwnerReadWrite); err != nil {
		t.Fatalf("write state: %v", err)
	}

	if err := removeStartupEntryAndAppData(configDirectory, appDirectory); err != nil {
		t.Fatalf("removeStartupEntryAndAppData() error = %v", err)
	}

	for _, path := range []string{startupEntryPath(configDirectory), appDirectory} {
		if _, err := os.Stat(path); !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("%s still exists, stat error = %v", path, err)
		}
	}
}
