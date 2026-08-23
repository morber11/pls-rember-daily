//go:build windows

package windows

import (
	"errors"
	"testing"

	"plsremberdaily/assets"
)

func TestPlayCallsOnStartedAfterPreparingVideo(t *testing.T) {
	wantErr := errors.New("started")
	called := false

	err := Play(assets.GetDefaultVideo(), func() error {
		called = true
		return wantErr
	})

	if !errors.Is(err, wantErr) {
		t.Fatalf("Play() error = %v, want %v", err, wantErr)
	}
	if !called {
		t.Fatal("Play() did not call onStarted")
	}
}
