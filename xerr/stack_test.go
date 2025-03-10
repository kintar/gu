package xerr

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

//go:noinline
func innerStackTraceFunc() StackTrace {
	return NewStackTrace()
}

//go:noinline
func outerStackTraceFunc() StackTrace {
	return innerStackTraceFunc()
}

func TestStackTrace_Format_Verb_V(t *testing.T) {
	st := outerStackTraceFunc()
	result := fmt.Sprintf("%v\n", st)
	lines := strings.Split(result, "\n")
	assert.Equal(t, 6, len(lines))
	lines = lines[:3]
	expectedLines := []string{
		"stack_test.go:13",
		"stack_test.go:18",
		"stack_test.go:22",
	}
	assert.Equal(t, expectedLines, lines)
}

func TestStackTrace_Format_Verb_N(t *testing.T) {
	st := outerStackTraceFunc()
	result := fmt.Sprintf("%n\n", st)
	lines := strings.Split(result, "\n")
	assert.Equal(t, 6, len(lines))
	lines = lines[:3]
	expectedLines := []string{
		"github.com/kintar/gu/xerr.innerStackTraceFunc",
		"github.com/kintar/gu/xerr.outerStackTraceFunc",
		"github.com/kintar/gu/xerr.TestStackTrace_Format_Verb_N",
	}
	assert.Equal(t, expectedLines, lines)
}

func TestStackTraceErr(t *testing.T) {
	err := StackTraceErr("err")
	var enh *StackTraceError
	assert.True(t, errors.As(err, &enh))
	assert.NotNil(t, enh.stackTrace)
	lines := strings.Split(fmt.Sprintf("%v", enh.stackTrace), "\n")
	assert.Equal(t, 3, len(lines))
	assert.Equal(t, "stack_test.go:50", lines[0])
}

func TestAddStackTrace(t *testing.T) {
	err := AddStackTrace(errors.New("some error"))
	var enh *StackTraceError
	assert.True(t, errors.As(err, &enh))
	assert.NotNil(t, enh.stackTrace)
	lines := strings.Split(fmt.Sprintf("%v", enh.stackTrace), "\n")
	assert.Equal(t, 3, len(lines))
	assert.Equal(t, "stack_test.go:60", lines[0])
}

func TestGetStackTrace(t *testing.T) {
	err := StackTraceErr("err")
	st := GetStackTrace(err)
	lines := strings.Split(fmt.Sprintf("%v", st), "\n")
	assert.Equal(t, 3, len(lines))
	assert.Equal(t, "stack_test.go:70", lines[0])
}
