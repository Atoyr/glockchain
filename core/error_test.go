package core_test

import (
	"errors"
	"testing"

	"github.com/atoyr/glockchain/core"
)

func Test_GlockchainError(t *testing.T) {
	err := errors.New("hoge")
	gError := core.NewGlockchainError()

	if core.IsGlockchainCause(err) {
		t.Fatal(err)
	}

	if !core.IsGlockchainCause(gError) {
		t.Fatal(gError)
	}
}
