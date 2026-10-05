//go:build windows

package windows

import xwindows "golang.org/x/sys/windows"

func ShowUninstallResult(uninstallErr error) {
	text := "PlsRemberDaily was uninstalled."
	var style uint32 = xwindows.MB_OK | xwindows.MB_ICONINFORMATION

	if uninstallErr != nil {
		text = "PlsRemberDaily could not be uninstalled.\n\n" + uninstallErr.Error()
		style = xwindows.MB_OK | xwindows.MB_ICONERROR
	}

	xwindows.MessageBox(0, xwindows.StringToUTF16Ptr(text), xwindows.StringToUTF16Ptr(appName), style)
}
