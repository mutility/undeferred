# Undeferred

`undeferred` is a
[go/analysis](https://pkg.go.dev/golang.org/x/tools/go/analysis)-based tool that
flags suspicious early bindings of named results in a defer call.

It will flag references to a named return value that is passed to the invocation
of a defer statement (line 8 below), and references in the body of a deferred
function to variables that shadowed a named return outside that defer (line 15 below).

[![CI](https://github.com/mutility/undeferred/actions/workflows/build.yaml/badge.svg)](https://github.com/mutility/undeferred/actions/workflows/build.yaml)

## Example messages

Given the following source code `example.go`:

```go
     1  package example
     2
     3  import (
     4      "errors"
     5  )
     6
     7  func foo() (err error) {
     8      defer errors.Join(err)
     9      return errors.New("foo")
    10  }
    11
    12  func bar() (err error) {
    13      if err := errors.New("bar"); err != nil {
    14          defer func() {
    15              err = nil
    16          }()
    17      }
    18      return
    19  }
```

undeferred will report the following:

```console
$ undeferred ./...
.../example.go:8:20: defer captures current value of named result 'err'
.../example.go:15:4: defer references shadow of named result 'err'
.../example.go:13:5: shadows named result 'err' referenced in later defer
```

## Usage

Run from source with `go run github.com/mutility/undeferred@latest` or
install with `go install github.com/mutility/undeferred@latest` and run
undeferred from GOPATH/bin.

You can configure behvior at the command line by passing the flags below, or in
library use by setting fields on `undefer.Analyzer()`.

Flag | Field | Meaning
-|-|-
`-shadow` | ReferencedShadows | Report references to shadowed named results (default: true)

## Bug reports and feature contributions

`undeferred` is developed in spare time, so while bug reports and feature
contributions are welcomed, it may take a while for them to be reviewed. If
possible, try to find a minimal reproduction before reporting a bug. Bugs that
are difficult or impossible to reproduce will likely be closed.

All bug fixes will include tests to help ensure no regression; correspondingly
all contributions should include such tests.

## Mutility Analyzers

`undeferred` is part of [mutility-analyzers](https://github.com/mutility/analyzers).
