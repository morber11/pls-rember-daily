//go:build linux

package platform

import (
	"plsremberdaily/internal/application"
	"plsremberdaily/internal/linux"
)

func NewRuntime() application.Platform {
	return &linux.Runtime{}
}

func Uninstall() error {
	return linux.Uninstall()
}
