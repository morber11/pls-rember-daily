//go:build linux

package linux

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"plsremberdaily/assets"
)

func TestCreateVideoMemoryFileContainsVideoData(t *testing.T) {
	videoData := []byte("video data")

	file, err := createVideoMemoryFile(videoData)
	if err != nil {
		t.Fatalf("createVideoMemoryFile() error = %v", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("read in-memory video file: %v", err)
	}

	if !bytes.Equal(data, videoData) {
		t.Fatalf("in-memory video data = %q, want %q", data, videoData)
	}
}

func TestPlayReturnsOnStartedErrorWithoutStartingPlayback(t *testing.T) {
	wantErr := errors.New("started")
	runtime := Runtime{}

	err := runtime.Play(assets.GetDefaultVideo(), func() error {
		return wantErr
	})

	if !errors.Is(err, wantErr) {
		t.Fatalf("Play() error = %v, want %v", err, wantErr)
	}
}
