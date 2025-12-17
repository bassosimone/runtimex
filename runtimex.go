//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Adapted from: https://github.com/ooni/probe-cli/blob/647b03f4270eb758106523fe6273e3ebdbcd599c/internal/runtimex/runtimex.go
//

// Package runtimex contains helpers for code paths that are not expected to fail and
// where failure indicates a programmer error or an unrecoverable condition.
//
// # When to use what
//
// AssertTrue/AssertNotError: For enforcing invariants in library code. Document the
// invariant and its justification in a comment above the assertion. Use these for
// conditions that should be impossible if the program is correct.
//
// Try0/Try1/Try2/Try3: For unwrapping (value, error) returns where the error cannot
// occur in correct usage. These are syntactic sugar over AssertNotError but improve
// readability when chaining operations.
//
// ExitOnError: In main() functions when you want to exit silently on error. Use when
// the error has already been logged or displayed elsewhere.
//
// LogFatalOnError: In main() functions when you want to log and exit. The error should
// already contain sufficient context - use the optional message parameters only for
// simple qualification like "loading config", not for complex formatting.
//
// This package was originally inspired by [github.com/m-lab/go/rtx].
package runtimex

import (
	"errors"
	"log"
	"os"
)

// AssertTrue panics if the given value is false. The value passed
// to `panic()` is an error constructed using [errors.New].
//
// You typically use this function to assert runtime invariants in your codebase
// to make it more robust. Document the invariant and its justification in a
// comment above the assertion. For example:
//
//	// Invariant: msg is never nil after successful initialization. This is
//	// guaranteed by the constructor, which validates input.
//	runtimex.AssertTrue(cfg.msg != nil)
//
// The correct approach is to assert for conditions that should be
// impossible if the program is correct.
func AssertTrue(value bool) {
	if !value {
		panic(errors.New("expected true, got false"))
	}
}

// AssertNotError panics if the given err is not nil. The value passed
// to `panic()` is the given err value.
//
// You typically use this function to assert runtime invariants
// in your codebase to make it more robust. For example:
//
//	runtimex.AssertNotError(operationThatCannotFail())
//
// The correct approach is to assert for errors that cannot
// possibly happen (e.g., [json.Marshal] applied to a struct
// that can always be marshalled to a JSON string).
//
// This function is aliased as [Try0] for consistency with the
// Try family of functions.
func AssertNotError(err error) {
	if err != nil {
		panic(err)
	}
}

// Try0 panics if the given err is not nil. The value passed
// to `panic()` is the given err value.
//
// This code is equivalent to and is an alias for [AssertNotError]
// providing a semantics consistent with [Try1], [Try2], etc.
var Try0 = AssertNotError

// Try1 panics if the given err is not nil. The value passed
// to `panic()` is the given err value.
//
// This code is equivalent to:
//
//	v1, err := fx()
//	runtimex.AssertNotError(err)
//
// but is more compact and improves readability when chaining operations.
func Try1[T1 any](v1 T1, err error) T1 {
	AssertNotError(err)
	return v1
}

// Try2 panics if the given err is not nil. The value passed
// to `panic()` is the given err value.
//
// This code is equivalent to:
//
//	v1, v2, err := fx()
//	runtimex.AssertNotError(err)
//
// but is more compact and improves readability when chaining operations.
func Try2[T1, T2 any](v1 T1, v2 T2, err error) (T1, T2) {
	AssertNotError(err)
	return v1, v2
}

// Try3 panics if the given err is not nil. The value passed
// to `panic()` is the given err value.
//
// This code is equivalent to:
//
//	v1, v2, v3, err := fx()
//	runtimex.AssertNotError(err)
//
// but is more compact and improves readability when chaining operations.
func Try3[T1, T2, T3 any](v1 T1, v2 T2, v3 T3, err error) (T1, T2, T3) {
	AssertNotError(err)
	return v1, v2, v3
}

// osExit allows testing [ExitOnError].
var osExit = os.Exit

// ExitOnError invokes [os.Exit] if the given err is not nil.
//
// This function DOES NOT print any error message. Use this in main()
// functions when the error has already been logged or displayed elsewhere.
// Use [LogFatalOnError] if you want to log the error before exiting.
//
// For example:
//
//	data, err := os.ReadFile("/etc/mytool/config.json")
//	runtimex.ExitOnError(err)
//
// The exit code is `1`, which indicates generic failure.
func ExitOnError(err error) {
	if err != nil {
		osExit(1)
	}
}

// logFatal allows testing [LogFatalOnError].
var logFatal = log.Fatal

// LogFatalOnError invokes [log.Fatal] if the given err is not nil.
//
// This function logs the error message before exiting. Use this in main()
// functions when you want to display the error. Use [ExitOnError] if you
// want to exit silently.
//
// The error should already contain sufficient context from error wrapping
// upstream. The optional message parameters are for simple qualification only,
// not for complex formatting. If not specified, the function will just log
// the error. Otherwise, the output contains all the provided strings separated
// by space followed by ": " and by the string representation of the error.
//
// For example:
//
//	data, err := os.ReadFile("/etc/mytool/config.json")
//	runtimex.LogFatalOnError(err, "loading config")
//
// On failure, this code would print something like:
//
//	loading config: open /etc/mytool/config.json: no such file or directory
//
// The exit code is `1`, which indicates generic failure.
func LogFatalOnError(err error, msgs ...string) {
	if err != nil {
		arguments := make([]any, 0, 1+len(msgs))
		for idx, msg := range msgs {
			if idx == len(msgs)-1 {
				msg += ":"
			}
			arguments = append(arguments, msg)
		}
		arguments = append(arguments, err.Error())
		logFatal(arguments...)
	}
}
