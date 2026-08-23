//go:build windows

package instance

import (
	"errors"
	"fmt"

	"golang.org/x/sys/windows"
)

const mutexName = `Local\PlsRemberDaily`

func Acquire() (windows.Handle, bool, error) {
	name, err := windows.UTF16PtrFromString(mutexName)
	if err != nil {
		return 0, false, fmt.Errorf("create instance with mutex name: %w", err)
	}

	handle, err := windows.CreateMutex(nil, true, name)
	if errors.Is(err, windows.ERROR_ALREADY_EXISTS) {
		windows.CloseHandle(handle)
		return 0, false, nil
	}

	if err != nil {
		return 0, false, fmt.Errorf("create instance with mutex: %w", err)
	}

	return handle, true, nil
}

func Release(handle windows.Handle) {
	windows.CloseHandle(handle)
}
