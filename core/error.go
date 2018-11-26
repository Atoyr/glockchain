package core

import (
	"fmt"

	"github.com/pkg/errors"
)

type glockchainCause interface {
	GlockchainCause() bool
}

type GlockchainError struct {
	code    int
	message string
}

func NewGlockchainError() error {
	glockchainError := &GlockchainError{-1, "hoge"}
	return glockchainError
}

func (e *GlockchainError) Error() string {
	return fmt.Sprintf("Glockchain error %d : %s", e.code, e.message)
}

func (e *GlockchainError) GlockchainCause() bool { return true }

func IsGlockchainCause(err error) bool {
	error, ok := errors.Cause(err).(glockchainCause)
	return ok && error.GlockchainCause()
}
