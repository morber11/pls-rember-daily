//go:build windows

package windows

/*
#cgo windows LDFLAGS: -lwtsapi32
#include "session_windows.h"
*/
import "C"

import "fmt"

func WaitForUnlock() error {
	if result := C.wait_for_unlocked_session(); result < 0 {
		return fmt.Errorf("wait for unlocked session: HRESULT 0x%08X", uint32(result))
	}

	return nil
}
