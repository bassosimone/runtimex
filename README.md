# Golang Runtime Extensions

[![GoDoc](https://pkg.go.dev/badge/github.com/bassosimone/runtimex)](https://pkg.go.dev/github.com/bassosimone/runtimex) [![Build Status](https://github.com/bassosimone/runtimex/actions/workflows/go.yml/badge.svg)](https://github.com/bassosimone/runtimex/actions) [![codecov](https://codecov.io/gh/bassosimone/runtimex/branch/main/graph/badge.svg)](https://codecov.io/gh/bassosimone/runtimex)

The `runtimex` Go package contains helpers for code paths:

1. that are not expected to fail;

2. where failure indicates a programmer error or an unrecoverable condition.

For example:

```Go
import "github.com/bassosimone/runtimex"

// 1. Runtime assertions invoking `panic` when invariants are not met.
runtimex.Assert(txp != nil)

// 2. Quick error unwrapping for functions that can't fail.
data := runtimex.PanicOnError1(json.Marshal("always marshals"))

// Avoiding `if err != nil { panic(err) }` in packages only used for testing.
req := runtimex.PanicOnError1(http.NewRequest("GET", URL, nil))

// Avoiding `if err != nil { log.Fatal(err) }` in `main` code.
resp := runtimex.LogFatalOnError1(txp.RoundTrip(req))
```

## Installation

To add this package as a dependency to your module:

```sh
go get github.com/bassosimone/runtimex
```

## Development

To run the tests:
```sh
go test -v .
```

To measure test coverage:
```sh
go test -v -cover .
```

## License

```
SPDX-License-Identifier: GPL-3.0-or-later
```

## History

Adapted from [ooni/probe-cli/internal/runtimex](https://github.com/ooni/probe-cli/blob/647b03f4270eb758106523fe6273e3ebdbcd599c/internal/runtimex/runtimex.go) and inspired by [m-lab/go/rtx](https://github.com/m-lab/go/blob/bff8a3a44a1a38ebec6db2521826fdb6c50fe975/rtx/rtx.go).
