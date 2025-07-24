package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	Errors []error
}

func (e *MultiError) Error() string {
	if e == nil || len(e.Errors) == 0 {
		return ""
	}

	var b strings.Builder
	fmt.Fprintf(&b, "%d errors occured:\n", len(e.Errors))

	for _, err := range e.Errors {
		fmt.Fprintf(&b, "\t* %s", err.Error())
	}

	b.WriteString("\n")
	return b.String()
}

func Append(err error, errs ...error) *MultiError {
	multiErr := &MultiError{}

	var me *MultiError

	if err != nil {
		if errors.As(err, &me) {
			multiErr.Errors = append(multiErr.Errors, me.Errors...)
		} else {
			multiErr.Errors = append(multiErr.Errors, err)
		}
	}

	for _, e := range errs {
		if e == nil {
			continue
		}

		if errors.As(e, &me) {
			multiErr.Errors = append(multiErr.Errors, me.Errors...)
		} else {
			multiErr.Errors = append(multiErr.Errors, e)
		}
	}

	return multiErr
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}
