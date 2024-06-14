[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/go-corelibs/rxp)
[![codecov](https://codecov.io/gh/go-corelibs/rxp/graph/badge.svg?token=H2SU6zZMqN)](https://codecov.io/gh/go-corelibs/rxp)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-corelibs/rxp)](https://goreportcard.com/report/github.com/go-corelibs/rxp)

# rxp

rxp is an experiment in doing regexp-like things, without actually using regexp
to do any of the work.

For most use cases, the regexp package is likely the correct choice as it is
fairly optimized and uses the familiar regular expression byte/string patterns
to compile and use to match and replace text.

rxp by contrast doesn't really have a compilation phase, rather it is simply
the declaration of a Pattern, which is really just a slice of Matcher functions,
and to do the neat things one needs to do with regular expressions, simply use
the methods on the Pattern list.

# Notice

This is the v0.10.x series, it works but likely not exactly as one would expect.
For example, the greedy-ness of things is incorrect, however, there are always
ways to write the patterns differently such that the greedy-ness issue is
irrelevant.

Please do not blindly use this project without at least writing specific unit
tests for all Patterns and methods required.

There are no safeguards against footguns and other such pitfalls.

# Installation

``` shell
> go get github.com/go-corelibs/rxp@latest
```

# Examples

## Find all words at the start of any line of input

``` go
// regexp version:
m := regexp.
    MustCompile(`(?:m)^\s*(\w+)\b`).
    FindAllStringSubmatch(input, -1)

// equivalent rxp version
m := rxp.Pattern{}.
    Caret("m").S("*").W("+", "c").B().
    FindAllStringSubmatch(input, -1)
```

## Perform a series of text transformations

For whatever reason, some text needs to be transformed and these transformations
must satisfy four requirements: lowercase everything, consecutive spaces become
one space, single quotes must be turned into underscores and all
non-alphanumeric-underscore-or-spaces be removed.

These requirements can be explored with the traditional Perl substitution
syntax, as in the following table:

 | # | Perl Expression      | Description                  |
 |---|----------------------|------------------------------|
 | 1 | s/[A-Z]+/\L${1}\e/mg | lowercase all letters        |
 | 2 | s/\s+/ /mg           | collapse all spaces          |
 | 3 | s/[']/_/mg           | single quotes to underscores |
 | 4 | s/[^\w\s]+//mg       | delete non-word-or-spaces    |

The result of the above should take: `Isn't this neat?` and transform it into:
`isn_t this neat`.

``` go
// using regexp:
output := strings.ToLower(`Isn't  this  neat?`)
output = regexp.MustCompile(`\s+`).ReplaceAllString(output, " ")
output = regexp.MustCompile(`[']`).ReplaceAllString(output, "_")
output = regexp.MustCompile(`[^\w ]`).ReplaceAllString(output, "")
// output == "isn_t this neat"

// using rxp:
output := rxp.Pipeline{}.
	Transform(strings.ToLower).
	Literal(rxp.S("+"), " ").
	Literal(rxp.Text("'"), "_").
	Literal(rxp.Not(rxp.W(), rxp.S(), "c"), "").
	Process(`Isn't  this  neat?`)
// output == "isn_t this neat"
```

# Benchmarks

These benchmarks can be regenerated using `make benchmark`.

## Historical (make benchstats-historical)

Given that performance is basically the entire point of the rxp package, here's
some benchmark statistics showing the evolution of the rxp package itself from
v0.1.0 to the current v0.10.0. Each of these releases are present in separate
pre-release branches so that curious developers can easily study the progression
of this initial development cycle.

```
goos: linux
goarch: arm64
pkg: github.com/go-corelibs/rxp
                     │     v0.1.0      │                 v0.2.0                  │                 v0.4.0                  │                 v0.8.0                  │                 v0.10.0                 │
                     │     sec/op      │     sec/op       vs base                │     sec/op       vs base                │     sec/op       vs base                │     sec/op       vs base                │
_FindAllString_Rxp      0.004292n ± 0%    0.003496n ± 1%  -18.56% (n=50)            0.002005n ± 0%  -53.30% (n=50)            0.001868n ± 1%  -56.48% (n=50)            0.002162n ± 1%  -49.62% (n=50)
_Pipeline_Combo_Rxp    0.0002862n ± 1%   0.0002866n ± 1%        ~ (p=0.920 n=50)   0.0002945n ± 3%   +2.86% (p=0.037 n=50)   0.0002910n ± 2%   +1.66% (p=0.010 n=50)   0.0002858n ± 1%        ~ (p=0.348 n=50)
_Pipeline_Readme_Rxp    0.047985n ± 0%    0.053185n ± 1%  +10.84% (p=0.000 n=50)    0.010055n ± 0%  -79.05% (n=50)            0.006152n ± 0%  -87.18% (n=50)            0.007117n ± 1%  -85.17% (n=50)
_Replace_ToUpper_Rxp     0.07639n ± 1%     0.08496n ± 1%  +11.22% (p=0.000 n=50)     0.03293n ± 0%  -56.90% (n=50)             0.02835n ± 0%  -62.89% (n=50)             0.02339n ± 1%  -69.38% (n=50)
geomean                 0.008192n         0.008203n        +0.13%                   0.003739n       -54.36%                   0.003120n       -61.91%                   0.003185n       -61.12%
```

## Versus Regexp (make benchstats-regexp)

These benchmarks are loosely comparing regexp with rxp in "as similar as can be
done" cases. While rxp seems to outperform regexp, note that a poorly crafted
rxp Pattern can easily tank performance, as in the pipeline readme case below.

```
goos: linux
goarch: arm64
pkg: github.com/go-corelibs/rxp
                 │     regexp      │                   rxp                   │
                 │     sec/op      │     sec/op       vs base                │
_FindAllString      0.005376n ± 0%    0.002162n ± 1%  -59.78% (n=50)
_Pipeline_Combo    0.0028000n ± 1%   0.0002858n ± 1%  -89.79% (n=50)
_Pipeline_Readme    0.005760n ± 1%    0.007117n ± 1%  +23.56% (p=0.000 n=50)
_Replace_ToUpper     0.02497n ± 0%     0.02339n ± 1%   -6.29% (p=0.000 n=50)
geomean             0.006821n         0.003185n       -53.31%
```

# Go-CoreLibs

[Go-CoreLibs] is a repository of shared code between the [Go-Curses] and
[Go-Enjin] projects.

# License

```
Copyright 2024 The Go-CoreLibs Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use file except in compliance with the License.
You may obtain a copy of the license at

 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

[Go-CoreLibs]: https://github.com/go-corelibs
[Go-Curses]: https://github.com/go-curses
[Go-Enjin]: https://github.com/go-enjin
