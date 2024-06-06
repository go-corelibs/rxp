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
	"unicode"
)

func (p Pattern) Add(matcher Matcher) Pattern {
	return append(p, matcher)
}

func (p Pattern) Or(options ...interface{}) Pattern {
	return append(p, Or(options...))
}

func (p Pattern) Not(options ...interface{}) Pattern {
	return append(p, Not(options...))
}

func (p Pattern) Group(options ...interface{}) Pattern {
	return append(p, Group(options...))
}

func (p Pattern) Dot(flags ...string) Pattern {
	return append(p, Dot(flags...))
}

func (p Pattern) Text(text string, flags ...string) Pattern {
	return append(p, Text(text, flags...))
}

func (p Pattern) Caret(flags ...string) Pattern {
	return append(p, Caret(flags...))
}

func (p Pattern) Dollar(flags ...string) Pattern {
	return append(p, Dollar(flags...))
}

func (p Pattern) A() Pattern {
	return append(p, A())
}

func (p Pattern) B() Pattern {
	return append(p, B())
}

func (p Pattern) Z() Pattern {
	return append(p, Z())
}

func (p Pattern) D(flags ...string) Pattern {
	return append(p, D(flags...))
}

func (p Pattern) S(flags ...string) Pattern {
	return append(p, S(flags...))
}

func (p Pattern) W(flags ...string) Pattern {
	return append(p, W(flags...))
}

func (p Pattern) Alnum(flags ...string) Pattern {
	return append(p, Alnum(flags...))
}

func (p Pattern) Alpha(flags ...string) Pattern {
	return append(p, Alpha(flags...))
}

func (p Pattern) Ascii(flags ...string) Pattern {
	return append(p, Ascii(flags...))
}

func (p Pattern) Blank(flags ...string) Pattern {
	return append(p, Blank(flags...))
}

func (p Pattern) Cntrl(flags ...string) Pattern {
	return append(p, Cntrl(flags...))
}

func (p Pattern) Digit(flags ...string) Pattern {
	return append(p, Digit(flags...))
}

func (p Pattern) Graph(flags ...string) Pattern {
	return append(p, Graph(flags...))
}

func (p Pattern) Lower(flags ...string) Pattern {
	return append(p, Lower(flags...))
}

func (p Pattern) Print(flags ...string) Pattern {
	return append(p, Print(flags...))
}

func (p Pattern) Punct(flags ...string) Pattern {
	return append(p, Punct(flags...))
}

func (p Pattern) Space(flags ...string) Pattern {
	return append(p, Space(flags...))
}

func (p Pattern) Upper(flags ...string) Pattern {
	return append(p, Upper(flags...))
}

func (p Pattern) Word(flags ...string) Pattern {
	return append(p, Word(flags...))
}

func (p Pattern) Xdigit(flags ...string) Pattern {
	return append(p, Xdigit(flags...))
}

func (p Pattern) NamedClass(name AsciiNames, flags ...string) Pattern {
	return append(p, NamedClass(name, flags...))
}

func (p Pattern) RangeTable(table *unicode.RangeTable, flags ...string) Pattern {
	return append(p, IsUnicodeRange(table, flags...))
}

func (p Pattern) R(characters string, flags ...string) Pattern {
	return append(p, R(characters, flags...))
}
