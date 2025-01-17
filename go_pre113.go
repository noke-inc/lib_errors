//go:build !go1.13
// +build !go1.13

package errors

import (
	"fmt"

	"golang.org/x/xerrors"
)

// Annotatable represents errors that are unwrappable and formatable.
// Errors that implement special behaviors should implement this interface
// either directly or through embedding in order to avoid breaking Wrap() semantics.
// As a simple example, given the following definitions:
//    type mySpecial struct { errors.Annotatable }
//    func (m mySpecial) Special() bool { return true }
// the following code would add the Special() behavior to a wrapped error
//    specialErr := mySpecial{errors.Wrap(someGenericError, "You are now a special error").(errors.Annotatable)}
type Annotatable interface {
	error
	fmt.Formatter
	Wrapper
}

// Is reports whether any error in err's chain matches target.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
func Is(err error, target error) bool {
	return xerrors.Is(err, target)
}

// As finds the first error in err's chain that matches the type to which target
// points, and if so, sets the target to its value and returns true. An error
// matches a type if it is assignable to the target type, or if it has a method
// As(interface{}) bool such that As(target) returns true. As will panic if target
// is not a non-nil pointer to a type which implements error or is of interface type.
//
// The As method should set the target to its value and return true if err
// matches the type to which target points.
func As(err error, target interface{}) bool {
	return xerrors.As(err, target)
}

// Unwrap returns the result of calling the Unwrap method on err, if err implements
// Unwrap. Otherwise, Unwrap returns nil.
func Unwrap(err error) error {
	return xerrors.Unwrap(err)
}

// A Wrapper provides context around another error.
type Wrapper xerrors.Wrapper
