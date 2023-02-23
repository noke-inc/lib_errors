package errors

import (
	"fmt"
	"io"
)

///// BASIC ANNOTATE-ABLE ERROR STUFF ////////////////////////////////////////////////////

type Basic struct {
	error
	data KVPairs
}

func (e *Basic) Wrap(err error) {
	e.error = err
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

// Unwrap returns the inner error if there is one, otherwise nil.
func (e *Basic) Unwrap() error { return e.error }

// Format prints the Basic error.
//
// "%s" - outputs all annotations from the outside in. Same as err.Error().
// "%q" - same as %s, but quoted.
// "%v" - same as %s.
//
// "%#s" - outputs the outermost message only.
// "%+v" - outputs each annotation with any associated message, debug data, and/or stack trace.
func (e *Basic) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			err := e.Unwrap()
			msg, mOk := e.data[msgKey]
			if err != nil {
				//fmt.Fprintf(s, "******** unwrap *******\n") // for debugging
				fmt.Fprintf(s, "%+v", e.Unwrap())
			}
			if err != nil {
				io.WriteString(s, "\n")
			}
			if mOk {
				io.WriteString(s, msg.(string))
			}
			io.WriteString(s, " >>>")
			d := make(KVPairs)
			e.addMyData(d)
			if len(d) > 0 {
				if err != nil || mOk {
					io.WriteString(s, "\n\t")
				}
				fmt.Fprintf(s, "ERROR DATA: %+v", d)
			}
			if stk, ok := e.data[stackKey]; ok {
				io.WriteString(s, "\n\tSTACK TRACE:")
				stk.(*stack).Format(s, verb)
			} else if stk, ok := e.data[abbrStackKey]; ok {
				io.WriteString(s, "\n\tSTACK TRACE (abbr.):")
				stk.(*stack).Format(s, verb)
			}
			return
		}
		fallthrough
	case 's':
		if s.Flag('#') {
			if msg, ok := e.data[msgKey]; ok {
				io.WriteString(s, msg.(string))
			} else {
				fmt.Fprintf(s, "%#s", e.Unwrap())
			}
		} else {
			io.WriteString(s, e.Error())
		}
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}

// StackTrace returns the first full stack trace found in the error chain.
// Search stops if e.Unwrap() ever returns []error.
// Returns nil if no full stack trace found.
func (e *Basic) StackTrace() StackTrace {
	if stk, ok := e.data[stackKey]; ok {
		return stk.(*stack).StackTrace()
	}

	if e.error == nil {
		return nil
	}

	var be *Basic
	switch x := e.error.(type) {
	case interface{ Unwrap() error }:
		if As(x.(error), &be) {
			return be.StackTrace()
		}
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
	if e.data == nil {
		e.data = make(KVPairs)
	}
	e.data[k] = v
	return nil
}

// SetData sets all the key/value pairs found in the given map.
// Reserved key strings are skipped.
func (e *Basic) SetData(d KVPairs) {
	if e.data == nil {
		e.data = make(KVPairs)
	}
	for key, val := range d {
		if e.isReservedKey(key) {
			continue
		}
		e.data[key] = val
	}
}

// GetValue returns the value for the first instance of key found in the entire error graph.
// The return value 'found' is true if key is found, false otherwise. (pre-order, depth-first)
func (e *Basic) GetValue(key string) (val interface{}, found bool) {
	if e.isReservedKey(key) {
		return nil, false
	}
	if e.data == nil {
		e.data = make(KVPairs)
	}

	val, found = e.data[key]
	if found {
		return
	}

	if e.error == nil {
		return nil, false
	}

	var de DataError
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

// GetAllData return a map of all key/value pairs for the entire error graph.
// First instance of each key found is used. (pre-order, depth-first)
func (e *Basic) GetAllData() KVPairs {
	d := make(KVPairs)

	if e.error != nil {
		var de DataError
		switch x := e.error.(type) {
		case interface{ Unwrap() error }:
			if As(x.(error), &de) {
				d = de.GetAllData()
			}
		case interface{ Unwrap() []error }:
			multi := x.Unwrap()
			last := len(multi)-1
			for i := range multi {
				if As(multi[last-i], &de) {
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
