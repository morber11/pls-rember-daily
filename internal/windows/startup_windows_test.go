//go:build windows

package windows

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"plsremberdaily/internal/filemode"
)

func TestInstallExecutableCopiesIntoConfigDirectory(t *testing.T) {
	executable := filepath.Join(t.TempDir(), "downloaded.exe")
	content := []byte("executable content")
	if err := os.WriteFile(executable, content, filemode.OwnerReadWriteExecute); err != nil {
		t.Fatalf("write executable: %v", err)
	}

	configDirectory := t.TempDir()

	installed, err := installExecutable(configDirectory, executable)
	if err != nil {
		t.Fatalf("installExecutable() error = %v", err)
	}

	want := filepath.Join(configDirectory, "PlsRemberDaily", "PlsRemberDaily.exe")
	if installed != want {
		t.Errorf("installExecutable() = %q, want %q", installed, want)
	}

	data, err := os.ReadFile(installed)
	if err != nil {
		t.Fatalf("read installed executable: %v", err)
	}

	if !bytes.Equal(data, content) {
		t.Errorf("installed executable content = %q, want %q", data, content)
	}
}
