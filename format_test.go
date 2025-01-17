package errors

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

func TestFormatNew(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		New("error"),
		"%s",
		"error",
	}, {
		New("error"),
		"%v",
		"error",
	}, {
		New("error"),
		"%+v",
		"error\n" +
			"github.com/noke-inc/lib_errors.TestFormatNew\n" +
			fmt.Sprintf("\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-4)),
	}, {
		New("error"),
		"%q",
		`"error"`,
	}}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.error, tt.format, tt.want)
	}
}

func TestFormatErrorf(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		Errorf("%s", "error"),
		"%s",
		"error",
	}, {
		Errorf("%s", "error"),
		"%v",
		"error",
	}, {
		Errorf("%s", "error"),
		"%+v",
		"error\n" +
			"github.com/noke-inc/lib_errors.TestFormatErrorf\n" +
			fmt.Sprintf("\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-4)),
	}}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.error, tt.format, tt.want)
	}
}

func TestFormatWrap(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		Wrap(New("error"), "error2"),
		"%s",
		"error2: error",
	}, {
		Wrap(New("error"), "error2"),
		"%v",
		"error2: error",
	}, {
		Wrap(New("error"), "error2"),
		"%+v",
		"error\n" +
			"github.com/noke-inc/lib_errors.TestFormatWrap\n" +
			fmt.Sprintf("\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-4)),
	}, {
		Wrap(io.EOF, "error"),
		"%s",
		"error: EOF",
	}, {
		Wrap(io.EOF, "error"),
		"%v",
		"error: EOF",
	}, {
		Wrap(io.EOF, "error"),
		"%+v",
		"EOF\n" +
			"error\n" +
			"github.com/noke-inc/lib_errors.TestFormatWrap\n" +
			fmt.Sprintf("\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-5)),
	}, {
		Wrap(Wrap(io.EOF, "error1"), "error2"),
		"%+v",
		"EOF\n" +
			"error1\n" +
			"github.com/noke-inc/lib_errors.TestFormatWrap\n" +
			fmt.Sprintf("\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-5)),
	}, {
		Wrap(New("error with space"), "context"),
		"%q",
		`"context: error with space"`,
	}}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.error, tt.format, tt.want)
	}
}

func TestFormatWrapf(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		Wrapf(io.EOF, "error%d", 2),
		"%s",
		"error2: EOF",
	}, {
		Wrapf(io.EOF, "error%d", 2),
		"%v",
		"error2: EOF",
	}, {
		Wrapf(io.EOF, "error%d", 2),
		"%+v",
		"EOF\n" +
			"error2\n" +
			"github.com/noke-inc/lib_errors.TestFormatWrapf\n" +
			fmt.Sprintf("\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-5)),
	}, {
		Wrapf(New("error"), "error%d", 2),
		"%s",
		"error2: error",
	}, {
		Wrapf(New("error"), "error%d", 2),
		"%v",
		"error2: error",
	}, {
		Wrapf(New("error"), "error%d", 2),
		"%+v",
		"error\n" +
			"github.com/noke-inc/lib_errors.TestFormatWrapf\n" +
			fmt.Sprintf("\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-4)),
	}}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.error, tt.format, tt.want)
	}
}

func TestFormatWithStack(t *testing.T) {
	tests := []struct {
		error
		format string
		want   []string
	}{{
		WithStack(io.EOF),
		"%s",
		[]string{"EOF"},
	}, {
		WithStack(io.EOF),
		"%v",
		[]string{"EOF"},
	}, {
		WithStack(io.EOF),
		"%+v",
		[]string{"EOF",
			"github.com/noke-inc/lib_errors.TestFormatWithStack\n" +
				fmt.Sprintf("\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-4))},
	}, {
		WithStack(New("error")),
		"%s",
		[]string{"error"},
	}, {
		WithStack(New("error")),
		"%v",
		[]string{"error"},
	}, {
		WithStack(New("error")),
		"%+v",
		[]string{"error",
			"github.com/noke-inc/lib_errors.TestFormatWithStack\n" +
				fmt.Sprintf("\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-4)),
			"github.com/noke-inc/lib_errors.TestFormatWithStack\n" +
				fmt.Sprintf("\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-6))},
	}, {
		WithStack(WithStack(io.EOF)),
		"%+v",
		[]string{"EOF",
			"github.com/noke-inc/lib_errors.TestFormatWithStack\n" +
				fmt.Sprintf("\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-4)),
			"github.com/noke-inc/lib_errors.TestFormatWithStack\n" +
				fmt.Sprintf("\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-6))},
	}, {
		WithStack(WithStack(Wrapf(io.EOF, "message"))),
		"%+v",
		[]string{"EOF",
			"message",
			"github.com/noke-inc/lib_errors.TestFormatWithStack\n" +
				fmt.Sprintf("\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-5)),
			"github.com/noke-inc/lib_errors.TestFormatWithStack\n" +
				fmt.Sprintf("\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-7)),
			"github.com/noke-inc/lib_errors.TestFormatWithStack\n" +
				fmt.Sprintf("\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-9))},
	}, {
		WithStack(Errorf("error%d", 1)),
		"%+v",
		[]string{"error1",
			"github.com/noke-inc/lib_errors.TestFormatWithStack\n" +
				fmt.Sprintf("\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-4)),
			"github.com/noke-inc/lib_errors.TestFormatWithStack\n" +
				fmt.Sprintf("\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-6))},
	}}

	for i, tt := range tests {
		testFormatCompleteCompare(t, i, tt.error, tt.format, tt.want, true)
	}
}

func TestFormatWithMessage(t *testing.T) {
	tests := []struct {
		error
		format string
		want   []string
	}{{
		WithMessage(New("error"), "error2"),
		"%s",
		[]string{"error2: error"},
	}, {
		WithMessage(New("error"), "error2"),
		"%v",
		[]string{"error2: error"},
	}, {
		WithMessage(New("error"), "error2"),
		"%+v",
		[]string{
			"error",
			"github.com/noke-inc/lib_errors.TestFormatWithMessage\n" +
				fmt.Sprintf("\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-5)),
			"error2"},
	}, {
		WithMessage(io.EOF, "addition1"),
		"%s",
		[]string{"addition1: EOF"},
	}, {
		WithMessage(io.EOF, "addition1"),
		"%v",
		[]string{"addition1: EOF"},
	}, {
		WithMessage(io.EOF, "addition1"),
		"%+v",
		[]string{"EOF", "addition1"},
	}, {
		WithMessage(WithMessage(io.EOF, "addition1"), "addition2"),
		"%v",
		[]string{"addition2: addition1: EOF"},
	}, {
		WithMessage(WithMessage(io.EOF, "addition1"), "addition2"),
		"%+v",
		[]string{"EOF", "addition1", "addition2"},
	}, {
		Wrap(WithMessage(io.EOF, "error1"), "error2"),
		"%+v",
		[]string{"EOF", "error1", "error2",
			"github.com/noke-inc/lib_errors.TestFormatWithMessage\n" +
				fmt.Sprintf("\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-4))},
	}, {
		WithMessage(Errorf("error%d", 1), "error2"),
		"%+v",
		[]string{"error1",
			"github.com/noke-inc/lib_errors.TestFormatWithMessage\n" +
				fmt.Sprintf("\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-4)),
			"error2"},
	}, {
		WithMessage(WithStack(io.EOF), "error"),
		"%+v",
		[]string{
			"EOF",
			"github.com/noke-inc/lib_errors.TestFormatWithMessage\n" +
				fmt.Sprintf("\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-5)),
			"error"},
	}, {
		WithMessage(Wrap(WithStack(io.EOF), "inside-error"), "outside-error"),
		"%+v",
		[]string{
			"EOF",
			"github.com/noke-inc/lib_errors.TestFormatWithMessage\n" +
				fmt.Sprintf("\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-5)),
			"inside-error",
			"github.com/noke-inc/lib_errors.TestFormatWithMessage\n" +
				fmt.Sprintf("\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-8)),
			"outside-error"},
	}}

	for i, tt := range tests {
		testFormatCompleteCompare(t, i, tt.error, tt.format, tt.want, true)
	}
}

func TestFormatWithData(t *testing.T) {
	tests := []struct {
		error
		format string
		want   []string
	}{{
		WithData(io.EOF, "key", "val"),
		"%s",
		[]string{"EOF"},
	}, {
		WithData(io.EOF, "key", "val"),
		"%v",
		[]string{"EOF"},
	}, {
		WithData(io.EOF, "key", "val"),
		"%+v",
		[]string{"EOF",
			"error data: map[key:val]"},
	}, {
		WithData(io.EOF, 1, 2, "nodata"), // 1 is not a valid key, "nodata" is a valid key but there is no val
		"%+v",
		[]string{"EOF"},
	}, {
		WithData(New("error"), "key", "val"),
		"%s",
		[]string{"error"},
	}, {
		WithData(New("error"), "key", "val"),
		"%v",
		[]string{"error"},
	}, {
		WithData(New("error"), "key", "val"),
		"%+v",
		[]string{"error",
			fmt.Sprintf("github.com/noke-inc/lib_errors.TestFormatWithData\n"+
				"\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-4)),
			"error data: map[key:val]"},
	}, {
		WithData(WithStack(io.EOF), "key", "val"),
		"%+v",
		[]string{"EOF",
			fmt.Sprintf("github.com/noke-inc/lib_errors.TestFormatWithData\n"+
				"\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-4)),
			"error data: map[key:val]"},
	}, {
		WithData(WithStack(Wrapf(io.EOF, "message")), "key", "val"),
		"%+v",
		[]string{"EOF",
			"message",
			fmt.Sprintf("github.com/noke-inc/lib_errors.TestFormatWithData\n"+
				"\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-5)),
			fmt.Sprintf("github.com/noke-inc/lib_errors.TestFormatWithData\n"+
				"\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-7)),
			"error data: map[key:val]"},
	}, {
		WithData(Errorf("error%d", 1), "key", "val"),
		"%+v",
		[]string{"error1",
			fmt.Sprintf("github.com/noke-inc/lib_errors.TestFormatWithData\n"+
				"\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-4)),
			"error data: map[key:val]"},
	}, {
		WithData(
			WithStack(
				WithMessage(
					WithData(
						Wrapf(io.EOF, "message"),
						"key", 1),
					"message2"),
			),
			"key", "val"),
		"%+v",
		[]string{"EOF",
			"message",
			fmt.Sprintf("github.com/noke-inc/lib_errors.TestFormatWithData\n"+
				"\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-9)),
			"error data: map[key:1]",
			"message2",
			fmt.Sprintf("github.com/noke-inc/lib_errors.TestFormatWithData\n"+
				"\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-16)),
			"error data: map[key:val]"},
	}}

	for i, tt := range tests {
		testFormatCompleteCompare(t, i, tt.error, tt.format, tt.want, true)
	}
}

func TestFormatWrapWithData(t *testing.T) {
	tests := []struct {
		error
		format string
		want   []string
	}{{ // 1
		WrapWithData(io.EOF, "message", "key", "val"),
		"%s",
		[]string{"message: EOF"},
	}, { // 2
		WrapWithData(io.EOF, "message", "key", "val"),
		"%v",
		[]string{"message: EOF"},
	}, { // 3
		WrapWithData(WrapWithData(io.EOF, "message", "key", "val"), "message2", "key2", "val2"),
		"%v",
		[]string{"message2: message: EOF"},
	}, { // 4
		WrapWithData(WrapWithData(io.EOF, "message", "key", "val"), "message2", "key2", "val2"),
		"%+v",
		[]string{"EOF",
			"message",
			"error data: map[key:val]",
			"github.com/noke-inc/lib_errors.TestFormatWrapWithData\n" +
				fmt.Sprintf("\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-6)),
			"message2",
			"error data: map[key2:val2]",
			"github.com/noke-inc/lib_errors.TestFormatWrapWithData\n" +
				fmt.Sprintf("\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-10))},
	}, { // 5
		WrapWithData(io.EOF, "message", 1, 2, "nodata"), // 1 is not a valid key, "nodata" is a valid key but there is no val
		"%+v",
		[]string{"EOF",
			"message",
			fmt.Sprintf("github.com/noke-inc/lib_errors.TestFormatWrapWithData\n"+
				"\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-5))},
	}, { // 6
		WrapWithData(New("error"), "message", "key", "val"),
		"%s",
		[]string{"message: error"},
	}, { // 7
		WrapWithData(New("error"), "message", "key", "val"),
		"%v",
		[]string{"message: error"},
	}, { // 8
		WrapWithData(New("error"), "message", "key", "val"),
		"%+v",
		[]string{"error",
			fmt.Sprintf("github.com/noke-inc/lib_errors.TestFormatWrapWithData\n"+
				"\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-4)),
			"message",
			"error data: map[key:val]",
			fmt.Sprintf("github.com/noke-inc/lib_errors.TestFormatWrapWithData\n"+
				"\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-8))},
	}, { // 9
		WrapWithData(WithStack(io.EOF), "message", "key", "val"),
		"%+v",
		[]string{"EOF",
			fmt.Sprintf("github.com/noke-inc/lib_errors.TestFormatWrapWithData\n"+
				"\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-4)),
			"message",
			"error data: map[key:val]",
			fmt.Sprintf("github.com/noke-inc/lib_errors.TestFormatWrapWithData\n"+
				"\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-8))},
	}, { // 10
		WrapWithData(WithStack(Wrapf(io.EOF, "message")), "message2", "key", "val"),
		"%+v",
		[]string{"EOF",
			"message",
			fmt.Sprintf("github.com/noke-inc/lib_errors.TestFormatWrapWithData\n"+
				"\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-5)),
			fmt.Sprintf("github.com/noke-inc/lib_errors.TestFormatWrapWithData\n"+
				"\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-7)),
			"message2",
			"error data: map[key:val]",
			fmt.Sprintf("github.com/noke-inc/lib_errors.TestFormatWrapWithData\n"+
				"\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-11))},
	}, { // 11
		WrapWithData(Errorf("error%d", 1), "message", "key", "val"),
		"%+v",
		[]string{"error1",
			fmt.Sprintf("github.com/noke-inc/lib_errors.TestFormatWrapWithData\n"+
				"\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-4)),
			"message",
			"error data: map[key:val]",
			fmt.Sprintf("github.com/noke-inc/lib_errors.TestFormatWrapWithData\n"+
				"\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-8))},
	}, { // 12
		WrapWithData(
			WithStack(
				WithMessage(
					WrapWithData(
						Wrapf(io.EOF, "message"),
						"messageA", "key", 1),
					"message2"),
			),
			"messageB", "key", "val"),
		"%+v",
		[]string{"EOF",
			"message",
			fmt.Sprintf("github.com/noke-inc/lib_errors.TestFormatWrapWithData\n"+
				"\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-9)),
			"messageA",
			"error data: map[key:1]",
			fmt.Sprintf("github.com/noke-inc/lib_errors.TestFormatWrapWithData\n"+
				"\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-14)),
			"message2",
			fmt.Sprintf("github.com/noke-inc/lib_errors.TestFormatWrapWithData\n"+
				"\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-19)),
			"messageB",
			"error data: map[key:val]",
			fmt.Sprintf("github.com/noke-inc/lib_errors.TestFormatWrapWithData\n"+
				"\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-24))},
	}}

	for i, tt := range tests {
		testFormatCompleteCompare(t, i, tt.error, tt.format, tt.want, true)
	}
}

func TestFormatGeneric(t *testing.T) {
	starts := []struct {
		err  error
		want []string
	}{
		{New("new-error"), []string{
			"new-error",
			"github.com/noke-inc/lib_errors.TestFormatGeneric\n" +
				fmt.Sprintf("\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-3))},
		}, {Errorf("errorf-error"), []string{
			"errorf-error",
			"github.com/noke-inc/lib_errors.TestFormatGeneric\n" +
				fmt.Sprintf("\t.+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-3))},
		}, {errors.New("errors-new-error"), []string{
			"errors-new-error"},
		},
	}

	wrappers := []wrapper{
		{
			func(err error) error { return WithMessage(err, "with-message") },
			[]string{"with-message"},
		}, {
			func(err error) error { return WithStack(err) },
			[]string{
				"github.com/noke-inc/lib_errors.(func·002|TestFormatGeneric.func2)\n\t" +
					fmt.Sprintf(".+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-3)),
			},
		}, {
			func(err error) error { return Wrap(err, "wrap-error") },
			[]string{
				"wrap-error",
				"github.com/noke-inc/lib_errors.(func·003|TestFormatGeneric.func3)\n\t" +
					fmt.Sprintf(".+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-4)),
			},
		}, {
			func(err error) error { return Wrapf(err, "wrapf-error%d", 1) },
			[]string{
				"wrapf-error1",
				"github.com/noke-inc/lib_errors.(func·004|TestFormatGeneric.func4)\n\t" +
					fmt.Sprintf(".+/github.com/noke-inc/lib_errors/format_test.go:%d", lineNum(-4)),
			},
		}, {
			func(err error) error { return WithData(err, "key", "val") },
			[]string{
				"error data: map[key:val]",
			},
		},
	}

	for s := range starts {
		err := starts[s].err
		want := starts[s].want
		testFormatCompleteCompare(t, s, err, "%+v", want, false)
		testGenericRecursive(t, err, want, wrappers, 3)
	}
}

func wrappedNew(message string) error { // This function will be mid-stack inlined in go 1.12+
	return New(message)
}

func TestFormatWrappedNew(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		wrappedNew("error"),
		"%+v",
		"error\n" +
			"github.com/noke-inc/lib_errors.wrappedNew\n" +
			fmt.Sprintf("\t.+/github.com/noke-inc/lib_errors/format_test.go:%d\n", lineNum(-13)) +
			"github.com/noke-inc/lib_errors.TestFormatWrappedNew\n" +
			fmt.Sprintf("\t.+/github.com/noke-inc/lib_errors/format_test.go:%d\n", lineNum(-6)),
	}}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.error, tt.format, tt.want)
	}
}

func lineNum(shift int) int {
	_, _, ln, _ := runtime.Caller(1)
	return ln + shift
}

func testFormatRegexp(t *testing.T, n int, arg interface{}, format, want string) {
	t.Helper()
	got := fmt.Sprintf(format, arg)
	gotLines := strings.SplitN(got, "\n", -1)
	wantLines := strings.SplitN(want, "\n", -1)

	if len(wantLines) > len(gotLines) {
		t.Errorf("test %d: wantLines(%d) > gotLines(%d):\n got: %q\nwant: %q", n+1, len(wantLines), len(gotLines), got, want)
		return
	}

	for i, w := range wantLines {
		match, err := regexp.MatchString(w, gotLines[i])
		if err != nil {
			t.Fatal(err)
		}
		if !match {
			t.Errorf("test %d: line %d: fmt.Sprintf(%q, err):\n got: %q\nwant: %q", n+1, i+1, format, got, want)
		}
	}
}

var stackLineR = regexp.MustCompile(`\.`)

// parseBlocks parses input into a slice, where:
//  - incase entry contains a newline, its a stacktrace
//  - incase entry contains no newline, its a solo line.
//
// Detecting stack boundaries only works incase the WithStack-calls are
// to be found on the same line, thats why it is optionally here.
//
// Example use:
//
// for _, e := range blocks {
//   if strings.ContainsAny(e, "\n") {
//     // Match as stack
//   } else {
//     // Match as line
//   }
// }
//
func parseBlocks(input string, detectStackboundaries bool) ([]string, error) {
	var blocks []string

	stack := ""
	wasStack := false
	lines := map[string]bool{} // already found lines

	for _, l := range strings.Split(input, "\n") {
		isStackLine := stackLineR.MatchString(l)

		switch {
		case !isStackLine && wasStack:
			blocks = append(blocks, stack, l)
			stack = ""
			lines = map[string]bool{}
		case isStackLine:
			if wasStack {
				// Detecting two stacks after another, possible cause lines match in
				// our tests due to WithStack(WithStack(io.EOF)) on same line.
				if detectStackboundaries {
					if lines[l] {
						if len(stack) == 0 {
							return nil, errors.New("len of block must not be zero here")
						}

						blocks = append(blocks, stack)
						stack = l
						lines = map[string]bool{l: true}
						continue
					}
				}

				stack = stack + "\n" + l
			} else {
				stack = l
			}
			lines[l] = true
		case !isStackLine && !wasStack:
			blocks = append(blocks, l)
		default:
			return nil, errors.New("must not happen")
		}

		wasStack = isStackLine
	}

	// Use up stack
	if stack != "" {
		blocks = append(blocks, stack)
	}
	return blocks, nil
}

func testFormatCompleteCompare(t *testing.T, n int, arg interface{}, format string, want []string, detectStackBoundaries bool) {
	gotStr := fmt.Sprintf(format, arg)

	got, err := parseBlocks(gotStr, detectStackBoundaries)
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != len(want) {
		t.Fatalf("test %d: fmt.Sprintf(%s, err) -> wrong number of blocks: got(%d) want(%d)\n got: %s\nwant: %s\ngotStr: %q",
			n+1, format, len(got), len(want), prettyBlocks(got), prettyBlocks(want), gotStr)
	}

	for i := range got {
		if strings.ContainsAny(want[i], "\n") {
			// Match as stack
			match, err := regexp.MatchString(want[i], got[i])
			if err != nil {
				t.Fatal(err)
			}
			if !match {
				t.Fatalf("test %d: block %d: fmt.Sprintf(%q, err):\ngot:\n%q\nwant:\n%q\nall-got:\n%s\nall-want:\n%s\n",
					n+1, i+1, format, got[i], want[i], prettyBlocks(got), prettyBlocks(want))
			}
		} else {
			// Match as message
			if got[i] != want[i] {
				t.Fatalf("test %d: fmt.Sprintf(%s, err) at block %d got != want:\n got: %q\nwant: %q", n+1, format, i+1, got[i], want[i])
			}
		}
	}
}

type wrapper struct {
	wrap func(err error) error
	want []string
}

func prettyBlocks(blocks []string) string {
	var out []string

	for _, b := range blocks {
		out = append(out, fmt.Sprintf("%v", b))
	}

	return "   " + strings.Join(out, "\n   ")
}

func testGenericRecursive(t *testing.T, beforeErr error, beforeWant []string, list []wrapper, maxDepth int) {
	if len(beforeWant) == 0 {
		panic("beforeWant must not be empty")
	}
	for _, w := range list {
		if len(w.want) == 0 {
			panic("want must not be empty")
		}

		err := w.wrap(beforeErr)

		// Copy required cause append(beforeWant, ..) modified beforeWant subtly.
		beforeCopy := make([]string, len(beforeWant))
		copy(beforeCopy, beforeWant)

		beforeWant := beforeCopy
		last := len(beforeWant) - 1
		var want []string

		// Merge two stacks behind each other.
		if strings.ContainsAny(beforeWant[last], "\n") && strings.ContainsAny(w.want[0], "\n") {
			want = append(beforeWant[:last], append([]string{beforeWant[last] + "((?s).*)" + w.want[0]}, w.want[1:]...)...)
		} else {
			want = append(beforeWant, w.want...)
		}

		testFormatCompleteCompare(t, maxDepth, err, "%+v", want, false)
		if maxDepth > 0 {
			testGenericRecursive(t, err, want, list, maxDepth-1)
		}
	}
}
