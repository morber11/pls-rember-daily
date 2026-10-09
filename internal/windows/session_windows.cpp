// go:build windows

#include "session_windows.h"

#include <windows.h>
#include <wtsapi32.h>

constexpr DWORD unlock_poll_interval_ms = 1000;

// windows can sign the user in automatically and then show the lock screen
// so wait until the user has unlocked the session before playing
extern "C" long wait_for_unlocked_session(void)
{
    while (true)
    {
        LPWSTR buffer = nullptr;
        DWORD bytes = 0;

        if (!WTSQuerySessionInformationW(
                WTS_CURRENT_SERVER_HANDLE,
                WTS_CURRENT_SESSION,
                WTSSessionInfoEx,
                &buffer,
                &bytes))
        {
            return HRESULT_FROM_WIN32(GetLastError());
        }

        const auto *info = reinterpret_cast<const WTSINFOEXW *>(buffer);
        const LONG session_state = info->Data.WTSInfoExLevel1.SessionFlags;
        WTSFreeMemory(buffer);

        if (session_state == WTS_SESSIONSTATE_UNLOCK)
        {
            return S_OK;
        }

        Sleep(unlock_poll_interval_ms);
    }
}
