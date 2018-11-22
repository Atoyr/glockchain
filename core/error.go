package core

import "github.com/pkg/errors"

type glockchainError interface {
	IsError() bool
}

func isGlockchainError(err error) bool {
	gerror, ok := errors.Cause(err).(glockchainError)
	return ok && gerror.IsError()
}
