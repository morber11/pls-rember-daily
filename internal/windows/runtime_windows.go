//go:build windows

package windows

import (
	"plsremberdaily/assets"
)

type Runtime struct{}

func (Runtime) Acquire() (func(), bool, error) {
	handle, acquired, err := Acquire()
	if err != nil || !acquired {
		return nil, acquired, err
	}

	return func() {
		Release(handle)
	}, true, nil
}

func (Runtime) RegisterStartup() error {
	return Register()
}

func (Runtime) Play(video assets.Video, onStarted func() error) error {
	return Play(video, onStarted)
}
