package ioutils

import (
	"fmt"
	"io"
)

// DeferredClose handles errors that happen with deferred calls
// to Closer on the provided closer.
// Note: since deferred function arguments are evaluated immediately, this
// function should always be called within an anonymous function.
func DeferredClose(closer io.Closer, err error) error {
	closeErr := closer.Close()
	if closeErr == nil {
		return err
	}
	if err == nil {
		return closeErr
	}
	return fmt.Errorf("close error: %v after %w", closeErr, err)
}

// PanicOnError panics on provided errors if they are non-nil
// Note: use sparingly, designed for use within the main package
func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
