[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/go-corelibs/rxp)
[![codecov](https://codecov.io/gh/go-corelibs/rxp/graph/badge.svg?token=)](https://codecov.io/gh/go-corelibs/rxp)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-corelibs/rxp)](https://goreportcard.com/report/github.com/go-corelibs/rxp)

# rxp

rxp is an experiment in doing regexp-like things, without actually using regexp
to do any of the work.

For most use cases, the regexp package is likely the correct choice as it is
fairly optimized and uses the familiar regular expression byte/string patterns
to compile and use to match and replace text.

rxp by contrast doesn't really have a compilation phase, rather it is simply
the declaration of a Pattern, which is simply a slice of Matcher functions,
and to do the neat things one needs to do with regular expressions, simply use
the methods on the Pattern list.

# Notice

This is the v0.4.x series with a little scaffolding still present (stuff devs
write just to get things working before refining and optimizing work).

This version is being published for the purpose of documenting the evolution of
rxp over time.

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

// using rxp:
output := rxp.Pipeline{}.
	Transform(strings.ToLower).
	ReplaceText(rxp.S("+"), " ").
	ReplaceText(rxp.Text("'"), "_").
	ReplaceText(rxp.Not(rxp.Or(rxp.W(), rxp.S()), "c"), "").
	Process(`Isn't  this  neat?`)
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
