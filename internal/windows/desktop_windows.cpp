// go:build windows

#include "desktop_windows.h"

#include <windows.h>

// windows starts a user session even if you have not logged in yet
// this was causing the application to start playing before the user
// had logged in properly, instead wait for the user desktop before playing
extern "C" long wait_for_user_desktop(void)
{
    while (true)
    {
        HDESK desktop = GetThreadDesktop(GetCurrentThreadId());

        if (desktop == nullptr)
        {
            return HRESULT_FROM_WIN32(GetLastError());
        }

        BOOL user_desktop_is_active = FALSE;
        DWORD required_bytes = 0;

        if (!GetUserObjectInformationW(
                desktop,
                UOI_IO,
                &user_desktop_is_active,
                sizeof(user_desktop_is_active),
                &required_bytes))
        {
            return HRESULT_FROM_WIN32(GetLastError());
        }

        if (user_desktop_is_active)
        {
            return S_OK;
        }

        Sleep(500);
    }
}
