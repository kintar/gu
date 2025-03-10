package xerr

import (
	"fmt"
	"io"
	"path"
	"runtime"
	"strconv"
)

type StackFrame runtime.Frame

// Format implements the fmt.Formatter interface for a StackFrame. It honors the following verbs and flags:
//
//	| Verb | Description                                         |
//	| -----|---------------------------------------------------- |
//	| %s   | source file                                         |
//	| %d   | source line                                         |
//	| %n   | function name                                       |
//	| %v   | equivalent to %s:%d                                 |
//	| %+s  | function name and compile-time path of source file, |
//	|      | if known                                            |
//	| %+v  | equivalent to %+s:%d                                |
//
// Note that the '+' flag will result in multiple lines for the frame, with the file path indented on the second line
// after the function name.
func (f StackFrame) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		switch {
		case s.Flag('+'):
			_, _ = io.WriteString(s, f.Function)
			_, _ = io.WriteString(s, "\n\r")
			_, _ = io.WriteString(s, f.File)
		default:
			_, _ = io.WriteString(s, path.Base(f.File))
		}
	case 'd':
		_, _ = io.WriteString(s, strconv.Itoa(f.Line))
	case 'n':
		_, _ = io.WriteString(s, f.Function)
	case 'v':
		// note that by design these calls are recursive
		f.Format(s, 's')
		_, _ = io.WriteString(s, ":")
		f.Format(s, 'd')
	}
}

// StackTrace is a stack of Frames from newest (innermost) to oldest (outermost)
type StackTrace []StackFrame

// NewStackTrace creates a new StackTrace struct with the current execution point's frames. If invoked with no arguments,
// the result has a maximum depth of 32 frames, starting at the function which called NewStackTrace.
//
// If called with arguments, the first argument is the number of callers to skip before collection begins (not counting
// NewStackTrace itself). If a second argument is present, this overrides the default max depth of 32. Further arguments
// are ignored.
func NewStackTrace(args ...int) StackTrace {
	var skip = 0
	var depth = 32

	// if the caller asked to skip
	if len(args) > 0 {
		skip = args[0] + skip
	}

	// if the caller asked to change the depth
	if len(args) > 1 {
		depth = args[1]
	}

	var pcs = make([]uintptr, depth)
	n := runtime.Callers(2+skip, pcs)
	if n == 0 {
		// No PCs available. This can happen if skip value is large.
		return nil
	}

	stack := make(StackTrace, 0, depth)
	frames := runtime.CallersFrames(pcs)
	for i := 0; i < depth; i++ {
		f, more := frames.Next()
		stack = append(stack, StackFrame(f))
		if !more {
			break
		}
	}

	return stack
}

// Format implements the fmt.Formatter interface for a StackTrace. It honors the same verbs as StackFrame.Format, and
// produces a newline after each frame
func (st StackTrace) Format(s fmt.State, verb rune) {
	if len(st) == 0 {
		return
	}

	// this loop only executes if there is more than one frame in the stack
	for f := 0; f < len(st)-1; f++ {
		st[f].Format(s, verb)
		_, _ = io.WriteString(s, "\n")
	}

	st[len(st)-1].Format(s, verb)
}
