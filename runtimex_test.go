// SPDX-License-Identifier: GPL-3.0-or-later

package runtimex

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssertTrue(t *testing.T) {
	t.Run("with true value does not panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			AssertTrue(true)
		})
	})

	t.Run("with false value panics", func(t *testing.T) {
		assert.PanicsWithError(t, "expected true, got false", func() {
			AssertTrue(false)
		})
	})
}

func TestAssertNotError(t *testing.T) {
	t.Run("with nil error does not panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			AssertNotError(nil)
		})
	})

	t.Run("with non-nil error panics", func(t *testing.T) {
		expectedErr := errors.New("test error")
		assert.PanicsWithValue(t, expectedErr, func() {
			AssertNotError(expectedErr)
		})
	})
}

func TestTry0(t *testing.T) {
	t.Run("with nil error does not panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			Try0(nil)
		})
	})

	t.Run("with non-nil error panics", func(t *testing.T) {
		expectedErr := errors.New("test error")
		assert.PanicsWithValue(t, expectedErr, func() {
			Try0(expectedErr)
		})
	})
}

func TestTry1(t *testing.T) {
	t.Run("with nil error returns value", func(t *testing.T) {
		expectedValue := "test value"
		var result string
		assert.NotPanics(t, func() {
			result = Try1(expectedValue, nil)
		})
		assert.Equal(t, expectedValue, result)
	})

	t.Run("with non-nil error panics", func(t *testing.T) {
		expectedErr := errors.New("test error")
		assert.PanicsWithValue(t, expectedErr, func() {
			Try1("value", expectedErr)
		})
	})

	t.Run("works with different types", func(t *testing.T) {
		intResult := Try1(42, nil)
		assert.Equal(t, 42, intResult)

		stringResult := Try1("hello", nil)
		assert.Equal(t, "hello", stringResult)

		type customStruct struct{ field int }
		structResult := Try1(customStruct{field: 123}, nil)
		assert.Equal(t, customStruct{field: 123}, structResult)
	})
}

func TestTry2(t *testing.T) {
	t.Run("with nil error returns values", func(t *testing.T) {
		expectedV1 := "first"
		expectedV2 := 42
		var v1 string
		var v2 int
		assert.NotPanics(t, func() {
			v1, v2 = Try2(expectedV1, expectedV2, nil)
		})
		assert.Equal(t, expectedV1, v1)
		assert.Equal(t, expectedV2, v2)
	})

	t.Run("with non-nil error panics", func(t *testing.T) {
		expectedErr := errors.New("test error")
		assert.PanicsWithValue(t, expectedErr, func() {
			Try2("value1", "value2", expectedErr)
		})
	})

	t.Run("works with different type combinations", func(t *testing.T) {
		s, i := Try2("hello", 123, nil)
		assert.Equal(t, "hello", s)
		assert.Equal(t, 123, i)

		b, f := Try2(true, 3.14, nil)
		assert.Equal(t, true, b)
		assert.Equal(t, 3.14, f)
	})
}

func TestTry3(t *testing.T) {
	t.Run("with nil error returns values", func(t *testing.T) {
		expectedV1 := "first"
		expectedV2 := 42
		expectedV3 := true
		var v1 string
		var v2 int
		var v3 bool
		assert.NotPanics(t, func() {
			v1, v2, v3 = Try3(expectedV1, expectedV2, expectedV3, nil)
		})
		assert.Equal(t, expectedV1, v1)
		assert.Equal(t, expectedV2, v2)
		assert.Equal(t, expectedV3, v3)
	})

	t.Run("with non-nil error panics", func(t *testing.T) {
		expectedErr := errors.New("test error")
		assert.PanicsWithValue(t, expectedErr, func() {
			Try3("value1", "value2", "value3", expectedErr)
		})
	})

	t.Run("works with different type combinations", func(t *testing.T) {
		s, i, b := Try3("hello", 123, true, nil)
		assert.Equal(t, "hello", s)
		assert.Equal(t, 123, i)
		assert.Equal(t, true, b)
	})
}

func TestExitOnError(t *testing.T) {
	// Save original osExit and restore after test
	originalOsExit := osExit
	defer func() { osExit = originalOsExit }()

	t.Run("with nil error does not exit", func(t *testing.T) {
		exitCalled := false
		osExit = func(code int) {
			exitCalled = true
		}

		ExitOnError(nil)
		assert.False(t, exitCalled, "osExit should not have been called")
	})

	t.Run("with non-nil error exits with code 1", func(t *testing.T) {
		var exitCode int
		exitCalled := false
		osExit = func(code int) {
			exitCode = code
			exitCalled = true
		}

		ExitOnError(errors.New("test error"))
		assert.True(t, exitCalled, "osExit should have been called")
		assert.Equal(t, 1, exitCode, "exit code should be 1")
	})
}

func TestLogFatalOnError(t *testing.T) {
	// Save original logFatal and restore after test
	originalLogFatal := logFatal
	defer func() { logFatal = originalLogFatal }()

	t.Run("with nil error does not log", func(t *testing.T) {
		logFatalCalled := false
		logFatal = func(v ...any) {
			logFatalCalled = true
		}

		LogFatalOnError(nil)
		assert.False(t, logFatalCalled, "logFatal should not have been called")
	})

	t.Run("with non-nil error and no messages logs only error", func(t *testing.T) {
		var loggedArgs []any
		logFatal = func(v ...any) {
			loggedArgs = v
		}

		expectedErr := errors.New("test error")
		LogFatalOnError(expectedErr)
		assert.Equal(t, []any{"test error"}, loggedArgs)
	})

	t.Run("with non-nil error and single message logs message and error", func(t *testing.T) {
		var loggedArgs []any
		logFatal = func(v ...any) {
			loggedArgs = v
		}

		expectedErr := errors.New("test error")
		LogFatalOnError(expectedErr, "loading config")
		assert.Equal(t, []any{"loading config:", "test error"}, loggedArgs)
	})

	t.Run("with non-nil error and multiple messages logs all messages and error", func(t *testing.T) {
		var loggedArgs []any
		logFatal = func(v ...any) {
			loggedArgs = v
		}

		expectedErr := errors.New("test error")
		LogFatalOnError(expectedErr, "fatal:", "cannot open", "config file")
		assert.Equal(t, []any{"fatal:", "cannot open", "config file:", "test error"}, loggedArgs)
	})

	t.Run("formats messages with colon after last message", func(t *testing.T) {
		var loggedArgs []any
		logFatal = func(v ...any) {
			loggedArgs = v
		}

		expectedErr := errors.New("file not found")
		LogFatalOnError(expectedErr, "error", "reading", "data")

		// Verify the last message has a colon appended
		assert.Len(t, loggedArgs, 4)
		assert.Equal(t, "error", loggedArgs[0])
		assert.Equal(t, "reading", loggedArgs[1])
		assert.Equal(t, "data:", loggedArgs[2])
		assert.Equal(t, "file not found", loggedArgs[3])
	})
}
