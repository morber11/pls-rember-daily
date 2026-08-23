//go:build windows

package windows

import (
	"errors"
	"fmt"

	xwindows "golang.org/x/sys/windows"
)

const mutexName = `Local\PlsRemberDaily`

func Acquire() (xwindows.Handle, bool, error) {
	name, err := xwindows.UTF16PtrFromString(mutexName)
	if err != nil {
		return 0, false, fmt.Errorf("create instance with mutex name: %w", err)
	}

	handle, err := xwindows.CreateMutex(nil, true, name)
	if errors.Is(err, xwindows.ERROR_ALREADY_EXISTS) {
		xwindows.CloseHandle(handle)
		return 0, false, nil
	}

	if err != nil {
		return 0, false, fmt.Errorf("create instance with mutex: %w", err)
	}

	return handle, true, nil
}

func Release(handle xwindows.Handle) {
	xwindows.CloseHandle(handle)
}
