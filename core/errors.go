package core

import (
	"fmt"

	"github.com/pkg/errors"
)

type glockchainCause interface {
	GlockchainCause() bool
}

// GlockchainError is error type for GlockchainError
type GlockchainError struct {
	code    int
	message string
}

// NewGlockchainError is GlockchainError constructor
func NewGlockchainError(code int) error {
	glockchainError := &GlockchainError{code, getErrorMessage(code)}
	return glockchainError
}

func (e *GlockchainError) Error() string {
	return fmt.Sprintf("Glockchain error %d : %s", e.code, e.message)
}

// GlockchainCause is decision GlockchainCause
func (e *GlockchainError) GlockchainCause() bool { return true }

// IsGlockchainCause is Error handling interface
func IsGlockchainCause(err error) bool {
	error, ok := errors.Cause(err).(glockchainCause)
	return ok && error.GlockchainCause()
}

func getErrorMessage(code int) string {
	var message string
	switch code {
	case 91001:
		message = "DB file not exist"
	case 91002:
		message = "DB file exist"
	case 91003:
		message = "DB update error"
	case 91004:
		message = "DB view error"
	case 93001:
		message = "Not have data into TXPool"
	case 93002:
		message = "Create coinbase tramsaction error"
	case 93003:
		message = "Create tramsaction error"
	case 93005:
		message = "Transaction signature error"
	case 94001:
		message = "Wallet not found"
	default:
		message = ""
	}
	return message
}
