//go:build go1.20
// +build go1.20

package errors

import (
	stderrors "errors"
)

func Join(errs ...error) error {
	return stderrors.Join(errs...)
}