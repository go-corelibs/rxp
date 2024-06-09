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

type Flags uint16

const (
	DefaultFlags Flags = 0
	NegatedFlag  Flags = 1 << iota
	MultilineFlag
	DotNewlineFlag
	AnyCaseFlag
	CaptureFlag
	ZeroOrMoreFlag
	ZeroOrOneFlag
	OneOrMoreFlag
	LessFlag
)

// ParseOptions accepts Pattern, Matcher and string options and recasts them
// into their specific types
//
// ParseOptions will panic with any type other than Pattern, Matcher or string
func ParseOptions(options ...interface{}) (pattern Pattern, flags []string, argv []interface{}) {

	for idx, value := range options {

		switch t := value.(type) {

		case nil:
			// allow nops
			continue

		case string:
			flags = append(flags, t)

		case Matcher:
			pattern = append(pattern, t)

		case Pattern:
			pattern = append(pattern, t...)

		default:
			panic(fmt.Errorf("invalid argument #%d: %#+v", idx, t))
		}

	}

	return
}

// ParseFlags parses a regexp-like option string into a Flags instance and two
// integers, the low and high range of repetitions
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
func ParseFlags(flags ...string) (Reps, Flags) {
	var f Flags
	var reps Reps

	for _, flag := range flags {

		flag = strings.ToLower(strings.TrimSpace(flag))
		lower := []rune(flag)
		if size := len(lower); size == 1 {
			if flg, lh, ok := f.parseFlag(lower[0]); ok {
				if lh != nil {
					reps = lh
				}
				f = flg
				continue
			}
		} else if flg, lh, ok := f.parseFlags(lower); ok {
			if lh != nil {
				reps = lh
			}
			f = flg
			continue
		}

		panic(fmt.Errorf("invalid flag: %q", flags))
	}

	return reps, f
}

func (f Flags) Set(flag Flags) Flags {
	return f | flag
}

func (f Flags) Unset(flag Flags) Flags {
	return f &^ flag
}

func (f Flags) Has(flag Flags) bool {
	return f&flag == flag
}

func (f Flags) Negated() bool {
	return f&NegatedFlag == NegatedFlag
}

func (f Flags) Multiline() bool {
	return f&MultilineFlag == MultilineFlag
}

func (f Flags) DotNL() bool {
	return f&DotNewlineFlag == DotNewlineFlag
}

func (f Flags) AnyCase() bool {
	return f&AnyCaseFlag == AnyCaseFlag
}

func (f Flags) Capture() bool {
	return f&CaptureFlag == CaptureFlag
}

func (f Flags) Less() bool {
	return f&LessFlag == LessFlag
}

func (f Flags) ZeroOrMore() bool {
	return f&ZeroOrMoreFlag == ZeroOrMoreFlag
}

func (f Flags) ZeroOrOne() bool {
	return f&ZeroOrOneFlag == ZeroOrOneFlag
}

func (f Flags) OneOrMore() bool {
	return f&OneOrMoreFlag == OneOrMoreFlag
}

func (f Flags) SetNegated() Flags {
	return f | NegatedFlag
}

func (f Flags) SetCapture() Flags {
	return f | CaptureFlag
}

func (f Flags) Merge(other Flags) Flags {
	return f | other
}

func (f Flags) String() string {
	var buf strings.Builder
	switch {
	case f.ZeroOrMore():
		buf.WriteRune('*')
	case f.OneOrMore():
		buf.WriteRune('+')
	case f.ZeroOrOne():
		buf.WriteRune('?')
	}
	if f.Less() {
		buf.WriteRune('?')
	}
	if f.Negated() {
		buf.WriteRune('^')
	}
	if f.Multiline() {
		buf.WriteRune('m')
	}
	if f.DotNL() {
		buf.WriteRune('s')
	}
	if f.AnyCase() {
		buf.WriteRune('i')
	}
	if f.Capture() {
		buf.WriteRune('c')
	}
	return buf.String()
}

func (f Flags) private(_ Flags) {}

func (f Flags) parseFlag(lower rune) (flags Flags, reps Reps, ok bool) {
	flags = f
	switch lower {
	case '^':
		flags = flags.Set(NegatedFlag)

	case 'm':
		flags = flags.Set(MultilineFlag)

	case 's':
		flags = flags.Set(DotNewlineFlag)

	case 'i':
		flags = flags.Set(AnyCaseFlag)

	case 'c':
		flags = flags.Set(CaptureFlag)

	case '*':
		reps = Reps{-1, -1}
		flags = flags.Unset(LessFlag).Set(ZeroOrMoreFlag)

	case '+':
		reps = Reps{1, -1}
		flags = flags.Unset(LessFlag).Set(OneOrMoreFlag)

	case '?':
		reps = Reps{0, 1}
		flags = flags.Unset(LessFlag).Set(ZeroOrOneFlag)

	case 0, ' ':
		// nop is ok

	default:
		// error
		return f, nil, false

	}
	return flags, reps, true
}

func (f Flags) parseFlags(input []rune) (flags Flags, reps Reps, ok bool) {
	flags = f

	for idx := 0; idx < len(input); idx += 1 {
		var next rune
		this := input[idx]
		if idx+1 < len(input) {
			next = input[idx+1]
		}

		switch this {

		case '*':
			reps = Reps{-1, -1}
			flags = flags.Set(ZeroOrMoreFlag)
			if next == '?' {
				idx += 1
				flags = flags.Set(LessFlag)
			} else {
				flags = flags.Unset(LessFlag)
			}

		case '+':
			reps = Reps{1, -1}
			flags = flags.Set(OneOrMoreFlag)
			if next == '?' {
				idx += 1
				flags = flags.Set(LessFlag)
			} else {
				flags = flags.Unset(LessFlag)
			}

		case '?':
			reps = Reps{0, 1}
			flags = flags.Set(ZeroOrOneFlag)
			if next == '?' {
				idx += 1
				flags = flags.Set(LessFlag)
			} else {
				flags = flags.Unset(LessFlag)
			}

		case '{':

			if v, flg, lh, ok := flags.parseRangeFlag(idx, input); ok {
				if lh != nil {
					reps = lh
				}
				idx = v
				flags = flg
				continue
			}

			return f, nil, false

		case ' ':
		// nop is allowed

		case '^', 'm', 's', 'i', 'c':
			flags, _, ok = flags.parseFlag(this)
			continue

		default:
			// error
			return f, nil, false
		}

	}

	return flags, reps, true
}

func (f Flags) parseRangeFlag(index int, input []rune) (idx int, flags Flags, reps Reps, ok bool) {
	idx = index
	flags = f

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
			flags = flags.Set(LessFlag)
		}
	}

	if jdx != idx+1 {
		idx = jdx // skip ahead
	}

	switch len(pair) {
	case 1:
		if exact, err := strconv.Atoi(pair[0]); err == nil {
			reps = Reps{exact, exact}
			return idx, flags, reps, true // keep going
		}
	case 2:
		if minimum, err := strconv.Atoi(pair[0]); err == nil {
			if pair[1] == "" {
				reps = Reps{minimum, -1}
				return idx, flags, reps, true // keep going
			} else if maximum, ee := strconv.Atoi(pair[1]); ee == nil {
				if minimum <= maximum {
					reps = Reps{minimum, maximum}
					return idx, flags, reps, true // keep going
				}
			}
		}

	}

	return 0, f, nil, false
}
