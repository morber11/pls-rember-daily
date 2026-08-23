//go:build windows

package player

/*
#cgo windows CXXFLAGS: -std=c++17
#cgo windows LDFLAGS: -lmfplat -lmfuuid -lole32 -loleaut32 -luuid -luser32 -lgdi32
#include "player_windows.h"
*/
import "C"

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"syscall"
	"time"
	"unsafe"

	"plsremberdaily/assets"
)

const (
	windowWidth  = 480
	windowHeight = 270
	videoPath    = "/video"
)

func Play(video assets.Video, onStarted func() error) error {
	videoURL, closeServer, err := serveVideo(video)
	if err != nil {
		return err
	}
	defer closeServer()

	if err := onStarted(); err != nil {
		return fmt.Errorf("mark reminder as shown: %w", err)
	}

	wideVideoURL, err := syscall.UTF16PtrFromString(videoURL)
	if err != nil {
		return fmt.Errorf("encode video URL: %w", err)
	}

	result := C.play_video(
		(*C.wchar_t)(unsafe.Pointer(wideVideoURL)),
		C.int(windowWidth),
		C.int(windowHeight),
	)

	if result < 0 {
		return fmt.Errorf("play video: HRESULT 0x%08X", uint32(result))
	}

	return nil
}

// set up a temp HTTP server to serve the video
// might be overkill but it's straightforward with how MF API works
func serveVideo(video assets.Video) (string, func(), error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", nil, fmt.Errorf("listen for video content: %w", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc(videoPath, func(response http.ResponseWriter, request *http.Request) {
		http.ServeContent(response, request, video.Name, time.Time{}, bytes.NewReader(video.Data))
	})

	server := &http.Server{Handler: mux}
	go server.Serve(listener)

	videoURL := (&url.URL{
		Scheme: "http",
		Host:   listener.Addr().String(),
		Path:   videoPath,
	}).String()

	return videoURL, func() {
		server.Close()
	}, nil
}
