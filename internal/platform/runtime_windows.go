//go:build windows

package platform

import (
	"plsremberdaily/internal/application"
	"plsremberdaily/internal/windows"
)

func NewRuntime() application.Platform {
	return windows.Runtime{}
}

func Uninstall() error {
	err := windows.Uninstall()
	windows.ShowUninstallResult(err)

	return err
}
