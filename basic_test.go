package errors

import (
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
		key string
		val interface{}
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
			"int": 100,
		},
	}, {
		"struct",
		struct{ 
			a bool 
			b int
		}{
			a:true,
			b:-5,
		},
		map[string]interface{}{
			"string": "I'm a string",
			"int": 100,
			"struct": struct{ 
				a bool 
				b int
			}{
				a:true,
				b:-5,
			},
		},
	}, {
		"int",
		200,
		map[string]interface{}{
			"string": "I'm a string",
			"int": 200,
			"struct": struct{ 
				a bool 
				b int
			}{
				a:true,
				b:-5,
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
			"greet": "Welcome to testing",
		},
	}, {
		KVPairs{
			"number": 144,
			"pi": 3.14,
			"popcorn": time.Date(1970, time.January, 1, 0, 1, 45, 0, loc),
			"_special": "reserved",
		},
		KVPairs{
			"number": 144,
			"greet": "Welcome to testing",
			"pi": 3.14,
			"popcorn": time.Date(1970, time.January, 1, 0, 1, 45, 0, loc),
		},
	}}

	for i, tt := range tests {
		e.SetData(tt.data)
		assert.Exactly(t, tt.want, e.data, "test %d failed\n", i+1)
	}
}


