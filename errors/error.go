package errors

import (
	"fmt"
	"io"
	"runtime"
	"strings"
)

var (
//NotFound = &Error{code: "400"}
)

type Error interface {
	error
	WithCode(string) Error
}

type trait struct {
	id    uint64
	label string
}

type innerError struct {
	err     error
	code    string
	trait   string
	arg     []interface{}
	markers []string
	frames  []eframe
}

func (e *innerError) Error() string {
	if e == nil || e.err == nil {
		return ""
	}
	var b strings.Builder
	if len(strings.TrimSpace(e.code)) > 0 {
		b.WriteString(fmt.Sprintf("[%s] - ", e.code))
	}
	for i := len(e.markers) - 1; i >= 0; i-- {
		b.WriteString(e.markers[i])
		b.WriteString(" : ")
	}
	b.WriteString(e.err.Error())
	return b.String()
}

func (e *innerError) hashTrait(trait trait) bool {
	//todo export error later.
	return false
}

/*
func (e *Error) Is(err error) bool {
	//todo
	return false
}

func (e *Error) As(target interface{}) bool {
	//todo
	return false
}
*/

func (e *innerError) Unwrap() error {
	return e.err
}

func hashTrait(err error, trait trait) bool {
	//todo export func after.
	return false
}

//
func (e *innerError) WithCode(code string) Error {
	e.code = code
	return e
}

func (e *innerError) Errorf(format string, args ...interface{}) *innerError {
	frames := trace(2)
	e.arg = args
	e.err = fmt.Errorf(format, args...)
	e.frames = frames
	return e
}

func output(err error) string {
	if err == nil {
		return ""
	}
	e, ok := err.(*innerError)
	if !ok {
		return err.Error()
	}
	//before, after, withSource := calcRows(nums)
	frames := e.frames
	expectedRows := len(frames) + 1
	//if withSource {
	//	expectedRows = (before+after+3)*len(frames) + 2
	//}
	rows := make([]string, 0, expectedRows)
	rows = append(rows, e.Error())
	//if withSource {
	//	rows = append(rows, "")
	//}
	for _, frame := range frames {
		message := frame.String()
		//if colorized {
		//	message = aurora.Bold(message).String()
		//}
		rows = append(rows, message)
		//if withSource {
		//	rows = sourceRows(rows, frame, before, after, colorized)
		//}
	}
	return strings.Join(rows, "\n")
}

func (e *innerError) Format(s fmt.State, verb rune) {
	//todo
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, output(e))
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.err.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.err.Error())
	}
}

//
func New(msg string) Error {
	return &innerError{
		frames: trace(2),
		err:    fmt.Errorf(msg),
	}
}

//
func Errorf(format string, args ...interface{}) Error {
	return &innerError{
		frames: trace(2),
		err:    fmt.Errorf(format, args...),
		arg:    args,
	}
}

func Wrap(err error) Error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*innerError); ok {
		return e
	}
	return &innerError{
		frames: trace(2),
		err:    err, //replace new err
		//arg:args,
	}
}

func WrapMark(err error, marker string) Error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*innerError); ok {
		e.markers = append(e.markers, marker)
		return e
	}
	return &innerError{
		frames: trace(2),
		err:    err, //replace new err
		//arg:args,
		markers: append(make([]string, 0), marker),
	}
}

func Unwrap(err error) error {
	u, ok := err.(interface {
		Unwrap() error
	})
	if !ok {
		return nil
	}
	return u.Unwrap()
}

func GetDefaultLocalizedMsg(code string, arg ...interface{}) string {
	return ""
}

func GetLocalizedMsgf(local string, code string, arg ...interface{}) string {
	return ""
}

// Frame is a single step in stack trace.
type eframe struct {
	// Func contains a function name.
	xfunc string
	// Line contains a line number.
	line int
	// Path contains a file path.
	path string
}

// String formats Frame to string.
func (f eframe) String() string {
	return fmt.Sprintf("%s:%d %s()", f.path, f.line, f.xfunc)
}

func trace(skip int) []eframe {
	frames := make([]eframe, 0, 20)
	for {
		pc, path, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		frame := eframe{
			xfunc: fn.Name(),
			line:  line,
			path:  path,
		}
		frames = append(frames, frame)
		skip++
	}
	return frames
}

func Printf(err error) {
	fmt.Printf("%+v\n", err)
}
