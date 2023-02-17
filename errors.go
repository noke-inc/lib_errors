package errors

import (
	"fmt"
	"io"
)

const (
	keyPrefix byte = '_'

	msgKey   = `_msg`
	stackKey = `_stack`
)

var errNilFlag = fmt.Errorf("nil flag")

type KVPairs map[string]interface{}

type DataKeeper interface {
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

// New 
func New(msg string) error {
	return doWrap(1, errNilFlag, nil, msg)
}

func Errorf(format string, args ...interface{}) error {
	return doWrap(1, errNilFlag, nil, format, args...)
}

func ErrorfD(data KVPairs, format string, args ...interface{}) error {
	return doWrap(1, errNilFlag, data, format, args...)
}

// Wrap adds a message and stack trace simultaneously
func Wrap(err error, format string, args ...interface{}) error {
	return doWrap(1, err, nil, format, args...)
}

// WrapD adds a message, debug data, and stack trace simultaneously
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
		data:  KVPairs{
			stackKey: callers(level),
			msgKey: formatMsg(format, args...),
		},
	}
	if data != nil {
		e.SetData(data)
	}

	return e
}

func WithMessage(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	} 
	
	if err == errNilFlag {
		err = nil
	}

	return &Basic{
		error: err,
		data:  KVPairs{msgKey: formatMsg(format, args...)},
	}
}

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

///// BASIC ANNOTATABLE ERROR STUFF ////////////////////////////////////////////////////

type Basic struct {
	error
	data KVPairs
}

// Error returns error message
func (e *Basic) Error() string {
	m := ``
	if msg, ok := e.data[msgKey]; ok {
		m += msg.(string)
		if e.error != nil {
			m += ": "
		}
	}
	if e.error != nil {
		m += e.error.Error()
	}

	return m
}

func (e *Basic) Unwrap() error { return e.error }

func (e *Basic) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			err := e.Unwrap();
			msg, mOk := e.data[msgKey];
			if err != nil {
				//fmt.Fprintf(s, "******** unwrap *******\n") // for debugging
				fmt.Fprintf(s, "%+v", e.Unwrap())
			}
			if  mOk {
				if err != nil {
					io.WriteString(s, "\n")
				}
				io.WriteString(s, msg.(string))
			}
			d := make(KVPairs)
			e.addMyData(d)
			if len(d) > 0 {
				if err != nil || mOk {
					io.WriteString(s, "\n")
				}
				fmt.Fprintf(s, "ERROR DATA: %+v", d)
			}
			if stk, ok := e.data[stackKey]; ok {
				stk.(*stack).Format(s, verb)
			}
			return
		}
		fallthrough
	case 's':
		if s.Flag('^') {
			if msg, ok := e.data[msgKey]; ok {
				io.WriteString(s, msg.(string))
			} else {
				fmt.Fprintf(s, "%^s", e.Unwrap())
			}
		}
		io.WriteString(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}

func (e *Basic) StackTrace() StackTrace {
	if stk, ok := e.data[stackKey]; ok {
		return stk.(*stack).StackTrace()
	} 
	return nil
}

////////// DATA STUFF //////////////////////////////

// SetKeyVal sets a single key/value pair.
// Error is returned if key is a reserved string.
func (e *Basic) SetKeyVal(k string, v interface{}) error {
	if e.isReservedKey(k) {
		return Errorf("cannot use a reserved key (string starting with '_')")
	}
	e.data[k] = v
	return nil
}

// SetData sets all the key/value pairs found in the given map.
// Reserved key strings are skipped.
func (e *Basic) SetData(d KVPairs) {
	for key, val := range d {
		if e.isReservedKey(key) {
			continue
		}
		e.data[key] = val
	}
}

// GetValue returns the value of the given key and true if the key is found.
// If not found, returns (nil, false).
func (e *Basic) GetValue(key string) (val interface{}, found bool) {
	if !e.isReservedKey(key) {
		val, found = e.data[key]
		if found {
			return
		}
	}
	if e.error == nil {
		return nil, false
	}

	var de DataKeeper
	switch x := e.error.(type) {
	case interface{ Unwrap() error }:
		if As(x.(error), &de) {
			if val, found = de.GetValue(key); found {
				return
			}
		}
	case interface{ Unwrap() []error }:
		for _, err := range x.Unwrap() {
			if As(err, &de) {
				if val, found = de.GetValue(key); found {
					return
				}
			}
		}
	}
	return nil, false
}

func (e *Basic) GetAllData() KVPairs {
	d := make(KVPairs)

	if e.error != nil {
		var de DataKeeper
		switch x := e.error.(type) {
		case interface{ Unwrap() error }:
			if As(x.(error), &de) {
				d = de.GetAllData()
			}
		case interface{ Unwrap() []error }:
			for _, err := range x.Unwrap() {
				if As(err, &de) {
					tmp := de.GetAllData()
					for k, v := range tmp {
						d[k] = v
					}
				}
			}
		}
	}

	e.addMyData(d)
	return d
}

func (e *Basic) addMyData(d KVPairs) {
	for key, val := range e.data {
		if e.isReservedKey(key) {
			continue
		}
		d[key] = val
	}
}

func (e *Basic) isReservedKey(key string) bool {
	return key[0] == keyPrefix
}
