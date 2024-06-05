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

	"github.com/go-corelibs/values"
)

var _ Flags = (*cFlags)(nil)

type cFlags uint64

type Flags interface {
	Negated() bool
	Multiline() bool
	DotNL() bool
	AnyCase() bool
	Capture() bool
	Less() bool

	SetNegated()
	SetCapture()

	Clone() Flags
	Merge(other Flags)
	Equal(other Flags) bool
	String() string

	private(f *cFlags)
}

const (
	fDefault cFlags = 0
	fNegated cFlags = 1 << iota
	fMultiline
	fDotNL
	fAnyCase
	fCapture
	fZeroOrMore
	fZeroOrOne
	fOneOrMore
	fLess
)

func mkDefaultFlags() *cFlags {
	return values.Ref(fDefault)
}

var (
	gDefaultFlags = values.Ref(cFlags(0))
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
	f := mkDefaultFlags()
	var reps Reps

	for _, flag := range flags {

		flag = strings.ToLower(strings.TrimSpace(flag))
		lower := []rune(flag)
		if size := len(lower); size == 1 {
			if lh, ok := f.parseFlag(lower[0]); ok {
				if lh != nil {
					reps = lh
				}
				continue
			}
		} else if lh, ok := f.parseFlags(lower); ok {
			if lh != nil {
				reps = lh
			}
			continue
		}

		panic(fmt.Errorf("invalid flag: %q", flags))
	}

	return reps, f
}

func (f *cFlags) set(flag cFlags) {
	*f = (*f) | flag
}

func (f *cFlags) unset(flag cFlags) {
	*f = (*f) &^ flag
}

func (f *cFlags) has(flag cFlags) bool {
	return (*f)&flag == flag
}

func (f *cFlags) Negated() bool {
	return f.has(fNegated)
}

func (f *cFlags) Multiline() bool {
	return f.has(fMultiline)
}

func (f *cFlags) DotNL() bool {
	return f.has(fDotNL)
}

func (f *cFlags) AnyCase() bool {
	return f.has(fAnyCase)
}

func (f *cFlags) Capture() bool {
	return f.has(fCapture)
}

func (f *cFlags) Less() bool {
	return f.has(fLess)
}

func (f *cFlags) ZeroOrMore() bool {
	return f.has(fZeroOrMore)
}

func (f *cFlags) ZeroOrOne() bool {
	return f.has(fZeroOrOne)
}

func (f *cFlags) OneOrMore() bool {
	return f.has(fOneOrMore)
}

func (f *cFlags) SetNegated() {
	f.set(fNegated)
}

func (f *cFlags) SetCapture() {
	f.set(fCapture)
}

func (f *cFlags) Clone() Flags {
	v := *f
	return &v
}

func (f *cFlags) clone() *cFlags {
	v := *f
	return &v
}

func (f *cFlags) Merge(other Flags) {
	if o, ok := other.(*cFlags); ok {
		*f |= *o
	}
}

func (f *cFlags) Equal(other Flags) bool {
	if o, ok := other.(*cFlags); ok {
		return *f == *o
	}
	return false
}

func (f *cFlags) String() string {
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

func (f *cFlags) private(_ *cFlags) {}

func (f *cFlags) parseFlag(lower rune) (reps Reps, ok bool) {
	switch lower {
	case '^':
		f.set(fNegated)

	case 'm':
		f.set(fMultiline)

	case 's':
		f.set(fDotNL)

	case 'i':
		f.set(fAnyCase)

	case 'c':
		f.set(fCapture)

	case '*':
		reps = Reps{-1, -1}
		f.unset(fLess)
		f.set(fZeroOrMore)

	case '+':
		reps = Reps{1, -1}
		f.unset(fLess)
		f.set(fOneOrMore)

	case '?':
		reps = Reps{0, 1}
		f.unset(fLess)
		f.set(fZeroOrOne)

	case 0, ' ':
		// nop is ok

	default:
		// error
		return nil, false

	}
	return reps, true
}

func (f *cFlags) parseFlags(input []rune) (reps Reps, ok bool) {

	for idx := 0; idx < len(input); idx += 1 {
		var next rune
		this := input[idx]
		if idx+1 < len(input) {
			next = input[idx+1]
		}

		switch this {

		case '*':
			reps = Reps{-1, -1}
			f.set(fZeroOrMore)
			if next == '?' {
				idx += 1
				f.set(fLess)
			} else {
				f.unset(fLess)
			}

		case '+':
			reps = Reps{1, -1}
			f.set(fOneOrMore)
			if next == '?' {
				idx += 1
				f.set(fLess)
			} else {
				f.unset(fLess)
			}

		case '?':
			reps = Reps{0, 1}
			f.set(fZeroOrOne)
			if next == '?' {
				idx += 1
				f.set(fLess)
			} else {
				f.unset(fLess)
			}

		case '{':

			if v, lh, ok := f.parseRangeFlag(idx, input); ok {
				if lh != nil {
					reps = lh
				}
				idx = v
				continue
			}

			return nil, false

		case ' ':
		// nop is allowed

		case '^', 'm', 's', 'i', 'c':
			if _, ok := f.parseFlag(this); ok {
				continue
			}

		default:
			// error
			return nil, false
		}

	}

	return reps, true
}

func (f *cFlags) parseRangeFlag(index int, input []rune) (idx int, reps Reps, ok bool) {
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
			f.set(fLess)
		}
	}

	if jdx != idx+1 {
		idx = jdx // skip ahead
	}

	switch len(pair) {
	case 1:
		if exact, err := strconv.Atoi(pair[0]); err == nil {
			reps = Reps{exact, exact}
			return idx, reps, true // keep going
		}
	case 2:
		if minimum, err := strconv.Atoi(pair[0]); err == nil {
			if pair[1] == "" {
				reps = Reps{minimum, -1}
				return idx, reps, true // keep going
			} else if maximum, ee := strconv.Atoi(pair[1]); ee == nil {
				if minimum <= maximum {
					reps = Reps{minimum, maximum}
					return idx, reps, true // keep going
				}
			}
		}

	}

	return 0, nil, false
}
