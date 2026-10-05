//go:build windows

package windows

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"plsremberdaily/internal/filemode"
)

func TestInstallExecutableCopiesIntoAppDirectory(t *testing.T) {
	executable := filepath.Join(t.TempDir(), "downloaded.exe")
	content := []byte("executable content")
	if err := os.WriteFile(executable, content, filemode.OwnerReadWriteExecute); err != nil {
		t.Fatalf("write executable: %v", err)
	}

	appDirectory := filepath.Join(t.TempDir(), "PlsRemberDaily")

	installed, err := installExecutable(appDirectory, executable)
	if err != nil {
		t.Fatalf("installExecutable() error = %v", err)
	}

	want := filepath.Join(appDirectory, "PlsRemberDaily.exe")
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
