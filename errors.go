package errors

import (
	"fmt"
)

const (
	keyPrefix byte = '_'

	msgKey   = `_msg`
	stackKey = `_stack`
	abbrStackKey = `_abbrstk`
)

// since we use one struct to represent errors and wrapper errors this variable
// come in handy internally to tell doWrap() whether to accept a nil error or not.
var errNilFlag = fmt.Errorf("nil flag")

// KVPairs is an abbreviation for a map of key/value pairs
type KVPairs map[string]interface{}

// DataError is an error with debug data included.
type DataError interface {
	error
	SetKeyVal(k string, v interface{}) error
	SetData(d KVPairs)
	GetValue(key string) (val interface{}, found bool)
	GetAllData() KVPairs
}

func formatMsg(format string, args ...interface{}) string {
	m := ``
	if len(args) <= 0 {
		m = format
	} else {
		m = fmt.Sprintf(format, args...)
	}
	return m
}

/////////////////////// ERROR CREATORS //////////////////////////

// New returns an error that formats as the given text.
// Each call to New returns a distinct error value even if the text is identical.
func New(msg string) error {
	return doWrap(1, errNilFlag, nil, msg)
}

// Errorf formats according to a format specifier and returns the string as a
// value that satisfies error.
//
// If the format specifier includes a %w verb with an error operand, the returned
// error will implement an Unwrap method returning the operand.
// If there is more than one %w verb, the returned error will implement an
// Unwrap method returning a []error containing all the %w operands in the
// order they appear in the arguments.
// It is invalid to supply the %w verb with an operand that does not implement the
// error interface.
// The %w verb is otherwise a synonym for %v.
func Errorf(format string, args ...interface{}) error {
	return doWrap(1, errNilFlag, nil, format, args...)
}

// ErrorfD functions the same as Errorf except that it also sets data as the
// internal debug data.
func ErrorfD(data KVPairs, format string, args ...interface{}) error {
	return doWrap(1, errNilFlag, data, format, args...)
}

///////////// ERROR WRAPPERS //////////////////////////////

// Wrap augments error with a message and stack trace simultaneously.
func Wrap(err error, format string, args ...interface{}) error {
	return doWrap(1, err, nil, format, args...)
}

// WrapD augments error with a message, debug data, and stack trace simultaneously.
func WrapD(err error, data KVPairs, format string, args ...interface{}) error {
	return doWrap(1, err, data, format, args...)
}

func doWrap(level int, err error, data KVPairs, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	if err == errNilFlag {
		err = nil
	}

	e := &Basic{
		error: err,
		data: KVPairs{
			msgKey:   formatMsg(format, args...),
		},
	}

	if abbr, stk := useAbbreviatedStack(err, callers(level)); abbr{
		e.data[abbrStackKey] = stk
	} else {
		e.data[stackKey] = stk
	}

	if data != nil {
		e.SetData(data)
	}

	return e
}

func useAbbreviatedStack(err error, s *stack) (bool, *stack) {
	var st interface{StackTrace() StackTrace}
	if As(err, &st) {
		inner := st.StackTrace()
		outer := s.StackTrace()
		lastIn := len(inner)-1
		lastOut := len(outer)-1
		var i int
		foundDiff := false
		for i = range outer {
			if outer[lastOut-i] != inner[lastIn-i] {
				foundDiff = true
				break
			}
		}
		switch {
		case i == 0: //do nothing
		case i == 1: 
			fallthrough
		case !foundDiff:
			stk := (*s)[:len(outer)-i]
			return true, &stk
		default:
			stk := (*s)[:len(outer)-i+1]
			return true, &stk
		}
	}
	return false, s
}

// WithMessage augments error with a message only.
func WithMessage(err error, format string, args ...interface{}) error {
	return WithMessageD(err, nil, format, args...)
}

// WithMessageD augments error with a message and debug data simultaneously.
func WithMessageD(err error, data KVPairs, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	if err == errNilFlag {
		err = nil
	}

	e := &Basic{
		error: err,
		data:  KVPairs{msgKey: formatMsg(format, args...)},
	}
	if data != nil {
		e.SetData(data)
	}
	return e
}

// WithStack augments error with a stack trace only.
func WithStack(err error) error {
	if err == nil {
		return nil
	}

	s := callers(0)

	if err == errNilFlag {
		err = nil
	}

	return &Basic{
		error: err,
		data:  KVPairs{stackKey: s},
	}
}

// WithData augments error with debug data only.
func WithData(err error, d map[string]interface{}) error {
	if err == nil {
		return nil
	}

	if err == errNilFlag {
		err = nil
	}

	e := &Basic{
		error: err,
		data:  KVPairs{},
	}
	e.SetData(d)
	return e
}

//////////////// OTHER STUFF /////////////////////////////////

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the standard
// errors.Wrapper interface:
//
//	type Wrapper interface {
//	       Unwrap() error
//	}
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
		e := c.Unwrap()
		if e == nil {
			break
		}
		err = e
	}
	return err
}
