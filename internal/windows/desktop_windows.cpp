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
        HDESK desktop = OpenInputDesktop(0, FALSE, DESKTOP_READOBJECTS);

        if (desktop != nullptr)
        {
            CloseDesktop(desktop);
            return S_OK;
        }

        const DWORD error = GetLastError();

        if (error != ERROR_ACCESS_DENIED)
        {
            return HRESULT_FROM_WIN32(error);
        }

        Sleep(500);
    }
}
