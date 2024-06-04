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
	"testing"

	c "github.com/smartystreets/goconvey/convey"
)

func TestMatchers(t *testing.T) {

	c.Convey("Panics", t, func() {

		c.So(func() {
			_ = newMatch(nil, nil)
		}, c.ShouldPanic)

	})

	c.Convey("Custom", t, func() {

		c.So(
			Pattern{}.
				Add(func(s MatchState) (stop, ok bool) {
					s.Capture()
					if prev, ok := s.Prev(); ok {
						if prev == 'o' {
							if this, ok := s.This(); ok {
								if this == 'n' {
									if next, ok := s.Next(); ok {
										if next == 'e' {
											s.Consume(1)
											return true, true
										}
									}
								}
							}
						}
					}
					return
				}).
				FindAllStringSubmatch("one", -1),
			c.ShouldEqual,
			[][]string{{"n", "n"}},
		)

		p := Pattern{}.
			Caret().
			Text("/w/").
			//Not(Text("/"), "+", "c").
			Text("/", "^", "+", "c").
			Text("/", "??").
			Dollar()
		for _, test := range []struct {
			input  string
			output [][]string
		}{
			{"/w/core", [][]string{{"/w/core", "core"}}},
			{"/w/core/", [][]string{{"/w/core/", "core"}}},
			{"/w//", [][]string(nil)},
			{"/nope/core/", [][]string(nil)},
		} {
			c.So(p.FindAllStringSubmatch(test.input, -1), c.ShouldEqual, test.output)
		}

		// ^/b/([^/]+)/??$
		p = Pattern{}.
			Caret().
			Text("/b/").
			Not(Text("/"), "+", "c").
			Text("/", "??").
			Dollar()
		for _, test := range []struct {
			input  string
			output [][]string
		}{
			{"/b/core", [][]string{{"/b/core", "core"}}},
			{"/b/core/", [][]string{{"/b/core/", "core"}}},
			{"/b//", [][]string(nil)},
			{"/nope/core/", [][]string(nil)},
		} {
			c.So(p.FindAllStringSubmatch(test.input, -1), c.ShouldEqual, test.output)
		}

		// ^/build/([a-zA-Z0-9])?/??
		p = Pattern{}.
			Caret().
			Text("/build/").
			Alnum("?", "c").
			Text("/", "??").
			Dollar()
		for _, test := range []struct {
			input  string
			output [][]string
		}{
			{"/build/c", [][]string{{"/build/c", "c"}}},
			{"/build/c/", [][]string{{"/build/c/", "c"}}},
			{"/build//", [][]string{{"/build//"}}},
			{"/nope/core/", [][]string(nil)},
		} {
			c.So(p.FindAllStringSubmatch(test.input, -1), c.ShouldEqual, test.output)
		}

	})
}

func TestDot(t *testing.T) {
	c.Convey("Dot", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  [][]string
		}{
			//{ // the .+? consumes the / and thus the final Matcher fails the match the pattern
			//  // the problem is that + and * do not look to the next Matcher for advice?
			//	input:   "/@func/",
			//	pattern: Pattern{}.Text("@").Dot("+?").Text("/"),
			//	output:  [][]string{{"@func/"}},
			//},

			// rxInvalidFuncName = regexp.MustCompile(`^\s*(\d+|func\d+)\s*$`)
			{
				input: " func1  ",
				pattern: Pattern{}.
					Caret().
					S("*").
					Or("c",
						D("+"),
						Group(Text("func"),
							D("+"),
						),
					).
					S("*").
					Dollar(),
				output: [][]string{{" func1  ", "func1"}},
			},

			{
				input:   "/@1.0.0/",
				pattern: Pattern{}.Text("@").Text("/", "+", "^", "c").Text("/"),
				output:  [][]string{{"@1.0.0/", "1.0.0"}},
			},

			{
				input: "stuff @func/more/stuff",
				pattern: Pattern{}.Add(func(m MatchState) (next, keep bool) {
					this, okt := m.This()
					if okt && this == '@' {
						consume := 1

						input := m.Input()
						total := len(input)
						start := m.Index()

						var found bool
						for idx := start + consume; idx < total; idx += 1 {

							consume += 1
							if found = input[idx] == '/'; found {
								break
							}

						}

						if found {
							//m.Capture()
							m.Consume(consume)
							return true, true
						}
					}
					return
				}),
				output: [][]string{{"@func/"}},
			},

			{input: "", pattern: nil, output: [][]string(nil)},

			{
				input:   "a\nb",
				pattern: Pattern{}.Add(Dot("c")),
				output:  [][]string{{"a", "a"}, {"b", "b"}},
			},

			{
				input:   "a\nb",
				pattern: Pattern{}.Dot("+", "s", "c"),
				output:  [][]string{{"a\nb", "a\nb"}},
			},

			{
				input:   "abc",
				pattern: Pattern{}.Dot("{2,3}?", "c"),
				output:  [][]string{{"ab", "ab"}},
			},

			{
				input:   "abc",
				pattern: Pattern{}.Dot("{2,3}", "c"),
				output:  [][]string{{"abc", "abc"}},
			},

			{
				input:   "abc\nn\no\np\ne",
				pattern: Pattern{}.Dot("{2,3}", "c"),
				output:  [][]string{{"abc", "abc"}},
			},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d", idx),
				test.pattern.FindAllStringSubmatch(test.input, -1),
				c.ShouldEqual,
				test.output,
			)
		}

	})
}

func TestText(t *testing.T) {
	c.Convey("Text", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  [][]string
		}{

			{
				input:   "aBabbb",
				pattern: Pattern{}.Text("b", "{1,2}", "ic"),
				output:  [][]string{{"B", "B"}, {"bb", "bb"}, {"b", "b"}},
			},

			{
				input:   "a\nb",
				pattern: Pattern{}.Text("b", "c"),
				output:  [][]string{{"b", "b"}},
			},

			{
				input:   "a\nb",
				pattern: Pattern{}.Text("\n", "c"),
				output:  [][]string{{"\n", "\n"}},
			},

			{
				input:   "abaaa",
				pattern: Pattern{}.Text("a", "{2,3}", "c"),
				output:  [][]string{{"aaa", "aaa"}},
			},

			{
				input:   "abaaa",
				pattern: Pattern{}.Text("a", "{2,3}?", "c"),
				output:  [][]string{{"aa", "aa"}},
			},

			{
				input:   "abaaa",
				pattern: Pattern{}.Text("b", "c"),
				output:  [][]string{{"b", "b"}},
			},

			{
				input:   "abaaa",
				pattern: Pattern{}.Text("", "c"),
				output:  [][]string(nil),
			},

			{
				input:   "abaaa",
				pattern: Pattern{}.Text("aa", "+", "c"),
				output:  [][]string{{"aa", "aa"}},
			},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d", idx),
				test.pattern.FindAllStringSubmatch(test.input, -1),
				c.ShouldEqual,
				test.output)
		}

	})
}
