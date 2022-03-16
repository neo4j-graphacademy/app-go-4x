package ioutils_test

import (
	"errors"
	"fmt"
	"github.com/neo4j-graphacademy/neoflix/pkg/ioutils"
)

func ExampleDeferredClose() {
	anErr := &anError{msg: "something wrong occurred"}
	err := func() (err error) {
		closer := &aCloser{err: fmt.Errorf("failed closing")}
		defer func() {
			err = ioutils.DeferredClose(closer, err)
		}()
		return anErr
	}()

	fmt.Println(err)
	fmt.Println(errors.Is(err, anErr))
	// Output:
	// close error: failed closing after something wrong occurred
	// true
}

func ExampleDeferredClose_withoutPreviousError() {
	err := func() (err error) {
		closer := &aCloser{err: fmt.Errorf("failed closing")}
		defer func() {
			err = ioutils.DeferredClose(closer, err)
		}()
		return nil
	}()

	fmt.Println(err)
	// Output:
	// failed closing
}

func ExampleDeferredClose_withoutErrors() {
	err := func() (err error) {
		closer := &aCloser{}
		defer func() {
			err = ioutils.DeferredClose(closer, err)
		}()
		return nil
	}()

	fmt.Println(err == nil)
	// Output:
	// true
}

type aCloser struct {
	err error
}

func (d *aCloser) Close() error {
	return d.err
}

type anError struct {
	msg string
}

func (s *anError) Error() string {
	return s.msg
}
