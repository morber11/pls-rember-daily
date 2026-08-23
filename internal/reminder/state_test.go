package reminder

import (
	"path/filepath"
	"testing"
	"time"
)

func TestLoadStateMissing(t *testing.T) {
	saved, err := loadState(filepath.Join(t.TempDir(), "state.json"))
	if err != nil {
		t.Fatalf("loadState() error = %v", err)
	}

	if saved.LastShown != "" {
		t.Fatalf("loadState() lastShown = %q, want empty", saved.LastShown)
	}
}

func TestSaveAndLoadState(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "state.json")
	now := time.Date(2026, time.August, 11, 9, 30, 0, 0, time.Local)

	if err := saveState(path, now); err != nil {
		t.Fatalf("saveState() error = %v", err)
	}

	saved, err := loadState(path)
	if err != nil {
		t.Fatalf("loadState() error = %v", err)
	}
	if saved.LastShown != "2026-08-11" {
		t.Errorf("loadState() lastShown = %q, want %q", saved.LastShown, "2026-08-11")
	}
}
