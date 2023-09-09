# prsteps

## Installation

```shell
$ go install github.com/knightso/prsteps@latest
$ prsteps --version
1.0.0
```

## Usage

```shell
$ prsteps knightso/base 9
file,additions,deletions
.examples/unusedroot/compare.go,22,0
.examples/unusedroot/go.mod,13,0
.examples/unusedroot/go.sum,13,0
.examples/unusedroot/type_assertion.go,22,0
errors/cmd/unusedroot/README.md,8,0
errors/cmd/unusedroot/main.go,10,0
errors/internal/unusedroot/analyzer.go,76,0
errors/internal/unusedroot/root_used_for_compare.go,42,0
errors/internal/unusedroot/root_used_for_type_assertion.go,33,0
go.mod,20,0
go.sum,42,0
```

## Accessing to private repositories

prsteps tries to get GitHub Personal Access Tokens(PAT) in the following order.

1. from the output of `gh auth token`
2. from the value of `--pat={pat}` flag
