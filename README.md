# Parameter Swap

`undeferred` is a
[go/analysis](https://pkg.go.dev/golang.org/x/tools/go/analysis)-based tool that
flags suspicious early bindings of parameters in a defer call.

[![CI](https://github.com/mutility/undeferred/actions/workflows/build.yaml/badge.svg)](https://github.com/mutility/undeferred/actions/workflows/build.yaml)

## Example messages

Given the following source code `example.go`:

```go
     1  package example
     2
     3  func foo() (err error) {
     4      defer cleanup(err)
     5      return io.EOF
     6  }
     7
     8  func main() {
     9      foo()
    10  }
```

undeferred will report the following:

```console
$ undeferred ./...
.../example.go:4:17: defer captures current value of named result 'err'
exit status 3
```

## Usage

Run from source with `go run github.com/mutility/undeferred@latest` or
install with `go install github.com/mutility/undeferred@latest` and run
undeferred from GOPATH/bin.

You can configure behvior at the command line by passing the flags below, or in
library use by setting fields on `undefer.Analyzer()`.

Flag | Field | Meaning
-|-|-
`-exact` | ExactTypeOnly | Suppress reports of mismatched parameters of mismatching types
`-gen` | IncludeGeneratedFiles | Include reports from generated files

## Bug reports and feature contributions

`undeferred` is developed in spare time, so while bug reports and feature
contributions are welcomed, it may take a while for them to be reviewed. If
possible, try to find a minimal reproduction before reporting a bug. Bugs that
are difficult or impossible to reproduce will likely be closed.

All bug fixes will include tests to help ensure no regression; correspondingly
all contributions should include such tests.

## Mutility Analyzers

`undeferred` is part of [mutility-analyzers](https://github.com/mutility/analyzers).
