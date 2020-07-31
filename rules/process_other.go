// +build !darwin,!linux

package rules

import (
	C "github.com/paradiseduo/clashr/constant"
)

func NewProcess(process string, adapter string) (C.Rule, error) {
	return nil, ErrPlatformNotSupport
}
