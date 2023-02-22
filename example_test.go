package errors_test

import (
	"fmt"

	"github.com/pantsmann/errors"
)

func ExampleNew() {
	err := errors.New("whoops")
	fmt.Println(err)

	// Output: whoops
}

func ExampleNew_printf() {
	err := errors.New("whoops")
	fmt.Printf("%+v", err)

	// Example output:
	// whoops
	// github.com/pantsmann/errors_test.ExampleNew_printf
	//         /home/dfc/src/github.com/pantsmann/errors/example_test.go:17
	// testing.runExample
	//         /home/dfc/go/src/testing/example.go:114
	// testing.RunExamples
	//         /home/dfc/go/src/testing/example.go:38
	// testing.(*M).Run
	//         /home/dfc/go/src/testing/testing.go:744
	// main.main
	//         /github.com/pantsmann/errors/_test/_testmain.go:106
	// runtime.main
	//         /home/dfc/go/src/runtime/proc.go:183
	// runtime.goexit
	//         /home/dfc/go/src/runtime/asm_amd64.s:2059
}

func ExampleWithMessage() {
	cause := errors.New("whoops")
	err := errors.WithMessage(cause, "oh noes")
	fmt.Println(err)

	// Output: oh noes: whoops
}

func ExampleWithStack() {
	cause := errors.New("whoops")
	err := errors.WithStack(cause)
	fmt.Println(err)

	// Output: whoops
}

func ExampleWithStack_printf() {
	cause := errors.New("whoops")
	err := errors.WithStack(cause)
	fmt.Printf("%+v", err)

	// Example Output:
	// whoops
	// github.com/pantsmann/errors_test.ExampleWithStack_printf
	//         /home/fabstu/go/src/github.com/pantsmann/errors/example_test.go:55
	// testing.runExample
	//         /usr/lib/go/src/testing/example.go:114
	// testing.RunExamples
	//         /usr/lib/go/src/testing/example.go:38
	// testing.(*M).Run
	//         /usr/lib/go/src/testing/testing.go:744
	// main.main
	//         github.com/pantsmann/errors/_test/_testmain.go:106
	// runtime.main
	//         /usr/lib/go/src/runtime/proc.go:183
	// runtime.goexit
	//         /usr/lib/go/src/runtime/asm_amd64.s:2086
	// github.com/pantsmann/errors_test.ExampleWithStack_printf
	//         /home/fabstu/go/src/github.com/pantsmann/errors/example_test.go:56
	// testing.runExample
	//         /usr/lib/go/src/testing/example.go:114
	// testing.RunExamples
	//         /usr/lib/go/src/testing/example.go:38
	// testing.(*M).Run
	//         /usr/lib/go/src/testing/testing.go:744
	// main.main
	//         github.com/pantsmann/errors/_test/_testmain.go:106
	// runtime.main
	//         /usr/lib/go/src/runtime/proc.go:183
	// runtime.goexit
	//         /usr/lib/go/src/runtime/asm_amd64.s:2086
}

func ExampleWrap() {
	cause := errors.New("whoops")
	err := errors.Wrap(cause, "oh noes")
	fmt.Println(err)

	// Output: oh noes: whoops
}

func fn() error {
	e1 := errors.New("error")
	e2 := errors.Wrap(e1, "inner")
	e3 := errors.Wrap(e2, "middle")
	return errors.Wrap(e3, "outer")
}

func ExampleCause() {
	err := fn()
	fmt.Println(err)
	fmt.Println(errors.Cause(err))

	// Output: outer: middle: inner: error
	// error
}

func ExampleWrap_extended() {
	err := fn()
	fmt.Printf("%+v\n", err)

	// Example output:
	// error
	// github.com/pantsmann/errors_test.fn
	//         /home/dfc/src/github.com/pantsmann/errors/example_test.go:47
	// github.com/pantsmann/errors_test.ExampleCause_printf
	//         /home/dfc/src/github.com/pantsmann/errors/example_test.go:63
	// testing.runExample
	//         /home/dfc/go/src/testing/example.go:114
	// testing.RunExamples
	//         /home/dfc/go/src/testing/example.go:38
	// testing.(*M).Run
	//         /home/dfc/go/src/testing/testing.go:744
	// main.main
	//         /github.com/pantsmann/errors/_test/_testmain.go:104
	// runtime.main
	//         /home/dfc/go/src/runtime/proc.go:183
	// runtime.goexit
	//         /home/dfc/go/src/runtime/asm_amd64.s:2059
	// github.com/pantsmann/errors_test.fn
	// 	  /home/dfc/src/github.com/pantsmann/errors/example_test.go:48: inner
	// github.com/pantsmann/errors_test.fn
	//        /home/dfc/src/github.com/pantsmann/errors/example_test.go:49: middle
	// github.com/pantsmann/errors_test.fn
	//      /home/dfc/src/github.com/pantsmann/errors/example_test.go:50: outer
}

func ExampleWrap_formatted() {
	cause := errors.New("whoops")
	err := errors.Wrap(cause, "oh noes #%d", 2)
	fmt.Println(err)

	// Output: oh noes #2: whoops
}

func ExampleErrorf_extended() {
	err := errors.Errorf("whoops: %s", "foo")
	fmt.Printf("%+v", err)

	// Example output:
	// whoops: foo
	// github.com/pantsmann/errors_test.ExampleErrorf
	//         /home/dfc/src/github.com/pantsmann/errors/example_test.go:101
	// testing.runExample
	//         /home/dfc/go/src/testing/example.go:114
	// testing.RunExamples
	//         /home/dfc/go/src/testing/example.go:38
	// testing.(*M).Run
	//         /home/dfc/go/src/testing/testing.go:744
	// main.main
	//         /github.com/pantsmann/errors/_test/_testmain.go:102
	// runtime.main
	//         /home/dfc/go/src/runtime/proc.go:183
	// runtime.goexit
	//         /home/dfc/go/src/runtime/asm_amd64.s:2059
}

func Example_stackTrace() {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	err, ok := errors.Cause(fn()).(stackTracer)
	if !ok {
		panic("oops, err does not implement stackTracer")
	}

	if st := err.StackTrace(); st != nil {
		fmt.Printf("%+v", st[0:2]) // top two frames
	}

	// Example output:
	// github.com/pantsmann/errors_test.fn
	//	/home/dfc/src/github.com/pantsmann/errors/example_test.go:47
	// github.com/pantsmann/errors_test.Example_stackTrace
	//	/home/dfc/src/github.com/pantsmann/errors/example_test.go:127
}

func ExampleCause_printf() {
	err := errors.Wrap(func() error {
		return func() error {
			return errors.Errorf("hello %s", fmt.Sprintf("world"))
		}()
	}(), "failed")

	fmt.Printf("%v", err)

	// Output: failed: hello world
}

func NewFancy() fancy {
	return fancy{&errors.Basic{}}
}

type fancy struct {
	*errors.Basic
}

func (f fancy) Fancy() bool {
	return true
}

func Example_typechecking() {
	// Assuming an error called fancy is defined as follows:
	//
	//	func NewFancy() fancy {
	//		return fancy{&errors.Basic{}}
	//	}
	//
	//	type fancy struct {
	//		*errors.Basic
	//	}
	//
	//	func (f fancy) Fancy() bool {
	//		return true
	//	}
	//
	// You can wrap an error in a special type like so:

	plain := errors.New("inside")
	fancy := NewFancy()
	fancy.Wrap(plain)
	layer3 := errors.WithMessage(fancy, "outside")

	//		The code above, including the code in the first comment block,
	//		would be defined in the function or file where the error gets
	//		created.
	//		The code below would be implemented in a different function that is a 
	//		consumer of the error.
	//		For example, the code above could be part of the model while the code 
	//		below is in an endpoint handler.
	//		Note the the Fancier interface is not explicitly defined, but it is
	//		implicitly part of the package fancy is defined in.

	// To test for Fancier errors further down the call stack you could do the following:
	isFancier := func(err error) bool {
		var fancier interface{ Fancy() bool }
		return errors.As(err, &fancier) && fancier.Fancy()
	}

	fmt.Printf("plain is a Fancier: %v\n", isFancier(plain))
	fmt.Printf("fancy is a Fancier: %v\n", isFancier(fancy))
	fmt.Printf("layer3 is a Fancier: %v\n", isFancier(layer3))

	// Output: plain is a Fancier: false
	// fancy is a Fancier: true
	// layer3 is a Fancier: true
}
