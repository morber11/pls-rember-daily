//go:build !windows && !linux

package platform

import (
	"errors"

	"plsremberdaily/assets"
	"plsremberdaily/internal/application"
)

var errUnsupported = errors.New("operating system is not supported")

type unsupportedRuntime struct{}

func NewRuntime() application.Platform {
	return unsupportedRuntime{}
}

func Uninstall() error {
	return errUnsupported
}

func (unsupportedRuntime) Acquire() (func(), bool, error) {
	return nil, false, errUnsupported
}

func (unsupportedRuntime) RegisterStartup() error {
	return errUnsupported
}

func (unsupportedRuntime) Play(assets.Video, func() error) error {
	return errUnsupported
}
