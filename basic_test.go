package errors

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStackTrace_nil(t *testing.T) {
	e := &Basic{}
	assert.Nil(t, e.StackTrace())
}

func TestSetKeyVal(t *testing.T) {
	e := &Basic{}

	// test reserved words are not allowed
	assert.ErrorContains(t, e.SetKeyVal("_reserved", "this should fail"), "cannot use a reserved key (string starting with '_')")

	// test setting and resetting various key/value pairs
	tests := []struct {
		key  string
		val  interface{}
		want KVPairs
	}{{
		"string",
		"I'm a string",
		map[string]interface{}{"string": "I'm a string"},
	}, {
		"int",
		100,
		map[string]interface{}{
			"string": "I'm a string",
			"int":    100,
		},
	}, {
		"struct",
		struct {
			a bool
			b int
		}{
			a: true,
			b: -5,
		},
		map[string]interface{}{
			"string": "I'm a string",
			"int":    100,
			"struct": struct {
				a bool
				b int
			}{
				a: true,
				b: -5,
			},
		},
	}, {
		"int",
		200,
		map[string]interface{}{
			"string": "I'm a string",
			"int":    200,
			"struct": struct {
				a bool
				b int
			}{
				a: true,
				b: -5,
			},
		},
	}}

	for i, tt := range tests {
		e.SetKeyVal(tt.key, tt.val)
		assert.Exactly(t, tt.want, e.data, "test %d failed\n", i+1)
	}

}

func TestSetData(t *testing.T) {
	e := &Basic{}
	loc := time.Now().Location()

	// test setting and resetting various key/value pairs
	tests := []struct {
		data KVPairs
		want KVPairs
	}{{
		nil,
		KVPairs{},
	}, {
		KVPairs{},
		KVPairs{},
	}, {
		KVPairs{
			"number": 72,
		},
		KVPairs{
			"number": 72,
		},
	}, {
		KVPairs{
			"greet": "Welcome to testing",
		},
		KVPairs{
			"number": 72,
			"greet":  "Welcome to testing",
		},
	}, {
		KVPairs{
			"number":   144,
			"pi":       3.14,
			"popcorn":  time.Date(1970, time.January, 1, 0, 1, 45, 0, loc),
			"_special": "reserved",
		},
		KVPairs{
			"number":  144,
			"greet":   "Welcome to testing",
			"pi":      3.14,
			"popcorn": time.Date(1970, time.January, 1, 0, 1, 45, 0, loc),
		},
	}}

	for i, tt := range tests {
		e.SetData(tt.data)
		assert.Exactly(t, tt.want, e.data, "test %d failed\n", i+1)
	}
}

func TestGetValue(t *testing.T) {
	e, e2, e4 := getTestErrors()

	tests := []struct {
		err   error
		key   string
		val   interface{}
		found bool
	}{
		{ // 1) no key, no map
			&Basic{},
			"none",
			nil,
			false,
		}, { // 2) no key
			e,
			"none",
			nil,
			false,
		}, { // 3) key exists, but is reserved key
			e2,
			msgKey,
			nil,
			false,
		}, { // 4) key exists
			e,
			"magic",
			144,
			true,
		}, { // 5) key exists
			e,
			"lonely",
			1,
			true,
		}, { // 6) key overridden
			e2,
			"lonely",
			0,
			true,
		}, { // 7) key found in outer
			e2,
			"moon",
			"io",
			true,
		}, { // 8) key found in inner
			e2,
			"dessert",
			"quark belchen",
			true,
		}, { // 9) key found only in inner of first
			e4,
			"magic",
			144,
			true,
		}, { // 10) key found in inner of first and in last
			e4,
			"dessert",
			"quark belchen",
			true,
		}, { // 11) key found in last only
			e4,
			"joke",
			"A man visits Chuck Norris and sees an aquarium with a mantis shrimp. He says to Chuck, \"I didn't know you have a pet!\"\n To which the shrimp replies, \"What's it to ya? I can have any pet I want!\"",
			true,
		},
	}

	for i, tt := range tests {
		err := tt.err.(*Basic)
		v, found := err.GetValue(tt.key)
		assert.Exactly(t, tt.val, v, "test %d.1 failed\n", i+1)
		assert.Exactly(t, tt.found, found, "test %d.2 failed\n", i+1)
	}
}

func TestGetAllData(t *testing.T) {
	e, _, e4 := getTestErrors()

	tests := []struct {
		err  error
		data KVPairs
	}{
		{ // 1) no map
			&Basic{},
			KVPairs{},
		}, { // 2) no data
			&Basic{data: KVPairs{}},
			KVPairs{},
		}, { // 3)
			New("the sky is falling"),
			KVPairs{},
		}, { // 4) no sub-data
			e,
			KVPairs{
				"magic":   144,
				"charm":   "quark",
				"dessert": "quark belchen",
				"lonely":  1,
			},
		}, { // 5) multiple, depth and width
			e4,
			KVPairs{
				"joke":    "A man visits Chuck Norris and sees an aquarium with a mantis shrimp. He says to Chuck, \"I didn't know you have a pet!\"\n To which the shrimp replies, \"What's it to ya? I can have any pet I want!\"",
				"moon":    "io",
				"magic":   144,
				"charm":   "bracelet",
				"dessert": "quark belchen",
				"lonely":  0,
			},
		},
	}

	for i, tt := range tests {
		err := tt.err.(*Basic)
		assert.Exactly(t, tt.data, err.GetAllData(), "test %d failed\n", i+1)
	}
}

func TestWithMessageD(t *testing.T) {
	err := WithMessageD(errNilFlag, KVPairs{"dummy": "crash test"}, "data error without stack trace")
	assert.Equal(t, 
		"data error without stack trace\nERROR DATA: map[dummy:crash test]", 
		fmt.Sprintf("%+v", err))
}

func getTestErrors() (error, error, error) {
	e := &Basic{}

	e.SetData(
		KVPairs{
			"magic":   144,
			"charm":   "quark",
			"dessert": "quark belchen",
			"lonely":  1,
		},
	)

	e2 := WrapD(
		e,
		KVPairs{
			"charm":  "bracelet",
			"lonely": 0,
			"moon":   "io",
		},
		"wrap blank error %s",
		"'cause it's embarassing",
	)

	e3 := ErrorfD(
		KVPairs{
			"dessert": "blueberry3.14",
			"joke":    "A man visits Chuck Norris and sees an aquarium with a mantis shrimp. He says to Chuck, \"I didn't know you have a pet!\"\n To which the shrimp replies, \"What's it to ya? I can have any pet I want!\"",
		},
		"I'm the younger brother of %#s",
		e2,
	)

	e4 := Wrap(Join(e2, e3), "multiple errors")
	return e, e2, e4
}
