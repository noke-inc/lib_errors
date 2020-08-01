// Package errors provides simple error handling primitives including stack, annotations, and data.
//
// Original package created by Dave Cheney
// Copyright (c) 2015, Dave Cheney <dave@cheney.net>
//
// Modifications of the original package by Friends of Go
// Copyright (c) 2019, Friends of Go <contact@friendsofgo.tech>
//
// Further modifications of the Friends of Go version by Adam Manwaring
// Copyright (c) 2019, Adam Manwaring <pantsmann@byu.net>
//
// The traditional error handling idiom in Go is roughly akin to
//
//     if err != nil {
//             return err
//     }
//
// which when applied recursively up the call stack results in error reports
// without context or debugging information. The errors package allows
// programmers to add context to the failure path in their code in a way
// that does not destroy the original value of the error.
//
// Adding context to an error
//
// The errors.Wrap function returns a new error that adds context to the
// original error by recording a stack trace at the point Wrap is called,
// together with the supplied message. For example
//
//     _, err := ioutil.ReadAll(r)
//     if err != nil {
//             return errors.Wrap(err, "read failed")
//     }
//
// If additional control is required, the errors.WithStack and
// errors.WithMessage functions destructure errors.Wrap into its component
// operations: annotating an error with a stack trace and with a message,
// respectively.
//
// Retrieving the cause of an error
//
// Using errors.Wrap constructs a stack of errors, adding context to the
// preceding error. Depending on the nature of the error it may be necessary
// to reverse the operation of errors.Wrap to retrieve the original error
// for inspection. Any error value which implements this interface
//
//     type causer interface {
//             Cause() error
//     }
//
// can be inspected by errors.Cause. errors.Cause will recursively retrieve
// the topmost error that does not implement causer, which is assumed to be
// the original cause. For example:
//
//     switch err := errors.Cause(err).(type) {
//     case *MyError:
//             // handle specifically
//     default:
//             // unknown error
//     }
//
// Although the causer interface is not exported by this package, it is
// considered a part of its stable public interface.
//
// With the new standard package error we have two new ways to figure what is the cause of
// our error:
//
//     var target *MyError
//     if errors.As(err, &target) {
//            // handle specifically
//     } else {
//            // unknown error
//     }
//
// or even with sentinel errors:
//
//     var ErrMyError = errors.New("my sentinel error")
//     if errors.Is(err, ErrMyError) {
//           // handle specifically
//     } else {
//          // unknown error
//     }
//
// Formatted printing of errors
//
// All error values returned from this package implement fmt.Formatter and can
// be formatted by the fmt package. The following verbs are supported:
//
//     %s    print the error. If the error is a wrapper it will be
//           printed recursively.
//     %v    see %s
//     %+v   extended format. Each Frame of the error's StackTrace will
//           be printed in detail. Any key/value pairs recorded WithData()
//           will also be printed.
//
// Retrieving the stack trace of an error or wrapper
//
// New, Errorf, Wrap, and Wrapf record a stack trace at the point they are
// invoked. This information can be retrieved with the following interface:
//
//     type stackTracer interface {
//             StackTrace() errors.StackTrace
//     }
//
// The returned errors.StackTrace type is defined as
//
//     type StackTrace []Frame
//
// The Frame type represents a call site in the stack trace. Frame supports
// the fmt.Formatter interface that can be used for printing information about
// the stack trace of this error. For example:
//
//     if err, ok := err.(stackTracer); ok {
//             for _, f := range err.StackTrace() {
//                     fmt.Printf("%+s:%d\n", f, f)
//             }
//     }
//
// Although the stackTracer interface is not exported by this package, it is
// considered a part of its stable public interface.
//
// See the documentation for Frame.Format for more details.
//
// Retrieving data from an error chain
//
// Key/value pairs recorded (using WithData or WrapWithData) in errors of the error chain can
// be retrieved with the following interface:
//
//     type dataCacher interface {
//          DataCache() map[string]interface{}
//     }
//
// The returned map will contain all the key/value pairs from any dataCachers of
// equal or greater depth in the error chain. For example
// (assuming the dataError definition above):
//
//     var d dataError
//     if errors.As(err, d) {
//          kv := d.DataCache()
//     }
//
// Note that if more than one error in the chain contains the same key, DataCache()
// returns only the shallowest value for that key. That is:
//     err = errors.WithData(err, "deep", 1, "dup", 1)
//     kv := WithData(err, "dup" 2, "shallow", 2).DataCache()
//     fmt.Printf("key/value pairs: %v", kv)
//
//     // Example Output:
//     // key/value pairs: map[deep:1 dup:2 shallow:2]
//
// The when using %+v to format the error all keys and values are output at the level
// they were set.
package errors

import (
	"fmt"
	"io"
)

// Base is an error with basic annotatable functionality.
// It can be used to simply add Annotatable to special error behaviors by extending
// an embedded error to fulfill Annotatable. For example given the following definitions:
//    type mySpecial struct { errors.Annotatable }
//    func (m mySpecial) Special() bool { return true }
// the following would add Special() as a behavior to someGenericError and this would be
// findable using errors.As and one would still be able to output stack traces inside of
// someGenericError (if they exist).
//    specialErr := mySpecial{errors.Base{someGenericError}}
type Base struct{ Err error }

// Error returns the error as a string
func (b Base) Error() string { return b.Err.Error() }

// Unwrap returns the internal error
func (b Base) Unwrap() error { return b.Err }

// Format outputs the error based on the given format state and verbs
func (b Base) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if f.Flag('+') {
			fmt.Fprintf(f, "%+v", b.Err)
			return
		}
		fallthrough
	case 's':
		io.WriteString(f, b.Err.Error())
	case 'q':
		fmt.Fprintf(f, "%q", b.Err)
	}
}

// New returns an error with the supplied message.
// New also records the stack trace at the point it was called.
func New(message string) error {
	return &fundamental{
		msg:   message,
		stack: callers(),
	}
}

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
// Errorf also records the stack trace at the point it was called.
func Errorf(format string, args ...interface{}) error {
	return &fundamental{
		msg:   fmt.Sprintf(format, args...),
		stack: callers(),
	}
}

// fundamental is an error that has a message and a stack, but no caller.
type fundamental struct {
	msg string
	*stack
}

func (f *fundamental) Error() string { return f.msg }

func (f *fundamental) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, f.msg)
			f.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, f.msg)
	case 'q':
		fmt.Fprintf(s, "%q", f.msg)
	}
}

// WithStack annotates err with a stack trace at the point WithStack was called.
// If err is nil, WithStack returns nil.
func WithStack(err error) error {
	if err == nil {
		return nil
	}
	return &withStack{
		err,
		callers(),
	}
}

type withStack struct {
	error
	*stack
}

func (w *withStack) Unwrap() error { return w.error }

func (w *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", w.Unwrap())
			w.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
	}
}

// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	err = &withMessage{
		error: err,
		msg:   message,
	}
	return &withStack{
		err,
		callers(),
	}
}

// Wrapf returns an error annotating err with a stack trace
// at the point Wrapf is called, and the format specifier.
// If err is nil, Wrapf returns nil.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	err = &withMessage{
		error: err,
		msg:   fmt.Sprintf(format, args...),
	}
	return &withStack{
		err,
		callers(),
	}
}

// WithMessage annotates err with a new message.
// If err is nil, WithMessage returns nil.
func WithMessage(err error, message string) error {
	if err == nil {
		return nil
	}
	return &withMessage{
		error: err,
		msg:   message,
	}
}

// WithMessagef annotates err with the format specifier.
// If err is nil, WithMessagef returns nil.
func WithMessagef(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &withMessage{
		error: err,
		msg:   fmt.Sprintf(format, args...),
	}
}

type withMessage struct {
	error
	msg string
}

func (w *withMessage) Error() string { return w.msg + ": " + w.error.Error() }

func (w *withMessage) Unwrap() error { return w.error }

func (w *withMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", w.Unwrap())
			io.WriteString(s, w.msg)
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, w.Error())
	}
}

// WithData annotates err with a map of key/value pairs.
// keyVals should be passed in as pairs; the first of each pair being a string (the key).
// If an odd number of keyVals are passed in, the last one is ignored.
// If a key value is not a string, it and the next keyVal element (the value) are skipped.
// If err is nil, WithData returns nil.
func WithData(err error, keyVals ...interface{}) error {
	if err == nil {
		return nil
	}
	e := &withData{
		err,
		make(map[string]interface{}),
	}
	for i := 0; (i + 1) < len(keyVals); i += 2 {
		if key, ok := keyVals[i].(string); !ok {
			continue
		} else {
			e.data[key] = keyVals[i+1]
		}
	}
	return e
}

// WrapWithData returns an error annotating err with a stack trace
// at the point WrapWithData is called, the supplied message, and
// any supplied key/value pairs.
// If err is nil, WrapWithData returns nil.
func WrapWithData(err error, message string, keyVals ...interface{}) error {
	if err == nil {
		return nil
	}
	err = &withMessage{
		error: err,
		msg:   message,
	}
	err = WithData(err, keyVals...)
	return &withStack{
		err,
		callers(),
	}
}

type withData struct {
	error
	data map[string]interface{}
}

// Unwrap provides compatibility for Go 1.13 error chains.
func (w *withData) Unwrap() error { return w.error }

// DataCache returns all key/value pairs in the error (including from wrapped errors)
func (w *withData) DataCache() map[string]interface{} {
	type dataCacher interface {
		DataCache() map[string]interface{}
	}

	var data dataCacher
	var kv map[string]interface{}

	if As(w.error, &data) {
		kv = data.DataCache()
	} else {
		kv = make(map[string]interface{})
	}
	for k, v := range w.data {
		kv[k] = v
	}
	return kv
}

func (w *withData) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			if len(w.data) > 0 {
				fmt.Fprintf(s, "%+v\nERROR DATA: %v", w.Unwrap(), w.data)
			} else {
				fmt.Fprintf(s, "%+v", w.Unwrap())
			}
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
	}
}

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the standard
// errors.Wrapper interface:
//
//     type Wrapper interface {
//            Unwrap() error
//     }
//
// If the error does not implement Wrapper, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause(err error) error {
	for err != nil {
		var c Wrapper
		if !As(err, &c) {
			break
		}
		err = c.Unwrap()
	}
	return err
}
