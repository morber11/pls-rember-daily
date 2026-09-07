//go:build linux

package linux

/*
#cgo pkg-config: gtk4 gtk4-x11
#include <stdlib.h>
#include "gtk_linux.h"
*/
import "C"

import (
	"fmt"
	"io"
	"os"
	goruntime "runtime"
	"unsafe"

	"golang.org/x/sys/unix"

	"plsremberdaily/assets"
)

const videoMemoryFileName = "pls-rember-daily-video"

type Runtime struct {
	application unsafe.Pointer
}

func init() {
	goruntime.LockOSThread()
}

func (runtime *Runtime) Acquire() (func(), bool, error) {
	nativeApplicationID := C.CString(applicationID)
	defer C.free(unsafe.Pointer(nativeApplicationID))

	application := C.linux_application_new(nativeApplicationID)

	var errorMessage *C.char
	if C.linux_application_register(application, &errorMessage) == 0 {
		C.linux_application_release(application)
		return nil, false, nativeError("register GTK application", errorMessage)
	}

	if C.linux_application_is_remote(application) != 0 {
		C.linux_application_release(application)
		return nil, false, nil
	}

	runtime.application = application

	return func() {
		C.linux_application_release(runtime.application)
		runtime.application = nil
	}, true, nil
}

func (runtime *Runtime) RegisterStartup() error {
	return Register()
}

func (runtime *Runtime) Play(video assets.Video, onStarted func() error) error {
	videoFile, err := createVideoMemoryFile(video.Data)
	if err != nil {
		return err
	}
	defer videoFile.Close()

	if err := onStarted(); err != nil {
		return err
	}

	videoPath := fmt.Sprintf("/proc/self/fd/%d", videoFile.Fd())
	nativeVideoPath := C.CString(videoPath)
	defer C.free(unsafe.Pointer(nativeVideoPath))

	errorMessage := C.linux_play_video(runtime.application, nativeVideoPath)
	if errorMessage != nil {
		return nativeError("play video", errorMessage)
	}

	return nil
}

func createVideoMemoryFile(videoData []byte) (*os.File, error) {
	descriptor, err := unix.MemfdCreate(videoMemoryFileName, unix.MFD_CLOEXEC)
	if err != nil {
		return nil, fmt.Errorf("create in-memory video file: %w", err)
	}

	file := os.NewFile(uintptr(descriptor), videoMemoryFileName)
	if _, err := file.Write(videoData); err != nil {
		file.Close()
		return nil, fmt.Errorf("write in-memory video file: %w", err)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		file.Close()
		return nil, fmt.Errorf("rewind in-memory video file: %w", err)
	}

	return file, nil
}

func nativeError(operation string, message *C.char) error {
	detail := C.GoString(message)
	C.linux_error_free(message)
	return fmt.Errorf("%s: %s", operation, detail)
}
