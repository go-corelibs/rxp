// Copyright (c) 2024  The Go-CoreLibs Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rxp

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseOptions accepts Pattern, Matcher and string options and recasts them
// into their specific types
//
// ParseOptions will panic with any type other than Pattern, Matcher or string
func ParseOptions(options ...interface{}) (pattern Pattern, flags []string, argv []interface{}) {

	for idx, value := range options {

		switch t := value.(type) {

		case string:
			flags = append(flags, t)

		case Matcher:
			pattern = append(pattern, t)

		case Pattern:
			pattern = append(pattern, t...)

		case nil:
			continue

		default:
			panic(fmt.Errorf("invalid argument #%d: %#+v", idx, t))
		}

	}

	return
}

// ParseFlags takes a regexp-like option string and builds a corresponding
// Config instance
//
//	|  Flags  | Description                                                                             |
//	|---------|-----------------------------------------------------------------------------------------|
//	|    ^    | Invert the meaning of this match group                                                  |
//	|    m    | Multiline mode Caret and Dollar match begin/end of line in addition to begin/end text   |
//	|    s    | DotNL allows Dot to match newlines (\n)                                                 |
//	|    i    | AnyCase is case-insensitive matching of unicode text                                    |
//	|    c    | Capture allows this Matcher to be included in Pattern substring results                 |
//	|    *    | zero or more repetitions, prefer more                                                   |
//	|    +    | one or more repetitions, prefer more                                                    |
//	|    ?    | zero or one repetition, prefer one                                                      |
//	|  {l,h}  | range of repetitions, l minimum and up to h maximum, prefer more                        |
//	|  {l,}   | range of repetitions, l minimum, prefer more                                            |
//	|  {l}    | range of repetitions, l minimum, prefer more                                            |
//	|   *?    | zero or more repetitions, prefer less                                                   |
//	|   +?    | one or more repetitions, prefer less                                                    |
//	|   ??    | zero or one repetition, prefer zero                                                     |
//	|  {l,h}? | range of repetitions, l minimum and up to h maximum, prefer less                        |
//	|  {l,}?  | range of repetitions, l minimum, prefer less                                            |
//	|  {l}?   | range of repetitions, l minimum, prefer less                                            |
//
// The flags presented above can be combined into a single string argument, or
// can be individually given to ParseFlags
//
// Any parsing errors will result in a runtime panic
func ParseFlags(flags ...string) (cfg *Config) {
	cfg = DefaultConfig()

	for _, flag := range flags {

		if cfg.Valid() {
			flag = strings.ToLower(strings.TrimSpace(flag))
			lower := []rune(flag)
			if size := len(lower); size == 1 {
				if cfg.parseFlag(lower[0]) {
					continue
				}
			} else if cfg.parseFlags(lower) {
				continue
			}
		}

		panic(fmt.Errorf("invalid config: %q", flags))
	}

	return
}

func (cfg *Config) parseFlag(lower rune) (ok bool) {
	switch lower {
	case '^':
		cfg.Negated = true

	case 'm':
		cfg.Multiline = true

	case 's':
		cfg.DotNL = true

	case 'i':
		cfg.AnyCase = true

	case 'c':
		cfg.Capture = true

	case '*':
		cfg.Minimum = -1
		cfg.Maximum = -1
		cfg.Less = false

	case '+':
		cfg.Minimum = 1
		cfg.Maximum = -1

	case '?':
		cfg.Minimum = 0
		cfg.Maximum = 1

	case 0, ' ':
		// nop is ok

	default:
		// error
		return false

	}
	return true
}

func (cfg *Config) parseFlags(input []rune) (ok bool) {

	for idx := 0; idx < len(input); idx += 1 {
		var next rune
		this := input[idx]
		if idx+1 < len(input) {
			next = input[idx+1]
		}

		switch this {

		case '*':
			cfg.Minimum = -1
			cfg.Maximum = -1
			// default config is already zero-or-more, prefer more
			if next == '?' {
				idx += 1
				cfg.Less = true
				// now prefers less
			}

		case '+':
			cfg.Minimum = 1
			cfg.Maximum = -1
			if next == '?' {
				idx += 1
				cfg.Less = true
			}

		case '?':
			cfg.Minimum = 0
			cfg.Maximum = 1
			if next == '?' {
				cfg.Less = true
			}

		case '{':

			if v, ok := cfg.parseRangeFlag(idx, input); ok {
				idx = v
				continue
			}

			// error, failed to keep going
			return false

		case ' ':
			// nop is allowed

		default:
			// error
			if cfg.parseFlag(this) {
				continue
			}
			return false
		}

	}

	return true
}

func (cfg *Config) parseRangeFlag(index int, input []rune) (idx int, ok bool) {
	idx = index

	pair := []string{""}
	var jdx int
	for jdx = idx + 1; jdx < len(input)-1 && input[jdx] != '}'; jdx += 1 {
		if input[jdx] == ',' {
			pair = append(pair, "")
		} else {
			pair[len(pair)-1] += string(input[jdx])
		}
	}
	if jdx+1 < len(input) {
		if input[jdx+1] == '?' {
			jdx += 1
			cfg.Less = true
		}
	}
	if jdx != idx+1 {
		idx = jdx // skip ahead
	}
	switch len(pair) {
	case 1:
		if exact, err := strconv.Atoi(pair[0]); err == nil {
			cfg.Minimum = exact
			cfg.Maximum = exact
			return idx, true // keep going
		}
	case 2:
		if minimum, err := strconv.Atoi(pair[0]); err == nil {
			if pair[1] == "" {
				cfg.Minimum = minimum
				cfg.Maximum = -1
				return idx, true // keep going
			} else if maximum, ee := strconv.Atoi(pair[1]); ee == nil {
				if minimum <= maximum {
					cfg.Minimum = minimum
					cfg.Maximum = maximum
					return idx, true // keep going
				}
			}
		}

	}

	return 0, false
}
