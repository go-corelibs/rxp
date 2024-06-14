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
	"strings"
	"testing"

	c "github.com/smartystreets/goconvey/convey"
)

func TestPattern_MatchRunes(t *testing.T) {

	c.Convey("batch", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			ok      bool
		}{

			{input: "", pattern: nil, ok: false},
			{input: "\naa", pattern: Pattern{}.Dot("{1}", "c"), ok: true},
			{input: "\naa", pattern: Pattern{}.Dot("{1}", "c"), ok: true},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.MatchRunes([]rune(test.input)),
				c.ShouldEqual,
				test.ok,
			)
		}

	})

}

func TestPattern_FindRunes(t *testing.T) {

	c.Convey("batch", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  []rune
		}{

			{input: "", pattern: nil, output: []rune(nil)},
			{input: "\naa", pattern: Pattern{}.Dot("{1}", "c"), output: []rune("a")},
			{input: "\naa", pattern: Pattern{}.Dot("{1}", "c"), output: []rune("a")},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.FindRunes([]rune(test.input)),
				c.ShouldEqual,
				test.output,
			)
		}

	})

}

func TestPattern_FindRunesIndex(t *testing.T) {

	c.Convey("batch", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  [2]int
		}{

			{input: "", pattern: nil, output: [2]int{}},
			{input: "aa", pattern: Pattern{}.Dot("{1}", "c"), output: [2]int{0, 1}},
			{input: "aa", pattern: Pattern{}.Dot("{1}", "c"), output: [2]int{0, 1}},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.FindRunesIndex([]rune(test.input)),
				c.ShouldEqual,
				test.output,
			)
		}

	})

}

func TestPattern_FindAllRunesIndex(t *testing.T) {

	c.Convey("batch", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			count   int
			output  [][2]int
		}{

			{input: "", pattern: nil, count: -1, output: [][2]int{{0, 0}}},
			{input: "aa", pattern: Pattern{}.Dot("{1}", "c"), count: 1, output: [][2]int{{0, 1}}},
			{input: "aa", pattern: Pattern{}.Dot("{1}", "c"), count: -1, output: [][2]int{{0, 1}, {1, 2}}},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.FindAllRunesIndex([]rune(test.input), test.count),
				c.ShouldEqual,
				test.output,
			)
		}

	})

}

func TestPattern_FindRunesSubmatch(t *testing.T) {

	c.Convey("batch", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  [][]rune
		}{

			{input: "", pattern: nil, output: [][]rune(nil)},
			{input: "\naa", pattern: Pattern{}.Dot("{1}", "c"), output: [][]rune{[]rune("a"), []rune("a")}},
			{input: "a\naa", pattern: Pattern{}.Dot("{1}", "c"), output: [][]rune{[]rune("a"), []rune("a")}},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.FindRunesSubmatch([]rune(test.input)),
				c.ShouldEqual,
				test.output,
			)
		}

	})

}

func TestPattern_FindAllRunes(t *testing.T) {

	c.Convey("batch", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			count   int
			output  [][]rune
		}{

			{input: "", pattern: nil, count: -1, output: [][]rune(nil)},
			{input: "aa", pattern: Pattern{}.Dot("{1}", "c"), count: 1, output: [][]rune{[]rune("a")}},
			{input: "aa", pattern: Pattern{}.Dot("{1}", "c"), count: -1, output: [][]rune{[]rune("a"), []rune("a")}},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.FindAllRunes([]rune(test.input), test.count),
				c.ShouldEqual,
				test.output,
			)
		}

	})

}

func TestPattern_FindAllRunesSubmatch(t *testing.T) {

	c.Convey("batch", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			count   int
			output  [][][]rune
		}{

			{input: "", pattern: nil, count: -1, output: [][][]rune(nil)},
			{input: "aa", pattern: Pattern{}.Dot("{1}", "c"), count: 1, output: [][][]rune{{{97}, {97}}}},
			{input: "aa", pattern: Pattern{}.Dot("{1}", "c"), count: -1, output: [][][]rune{{{97}, {97}}, {{97}, {97}}}},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.FindAllRunesSubmatch([]rune(test.input), test.count),
				c.ShouldEqual,
				test.output,
			)
		}

	})

}

func TestPattern_FindAllRunesSubmatchIndex(t *testing.T) {

	c.Convey("batch", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			count   int
			output  [][][2]int
		}{

			{input: "", pattern: nil, count: -1, output: [][][2]int(nil)},
			{input: "aa", pattern: Pattern{}.Dot("{1}", "c"), count: 1, output: [][][2]int{{{0, 1}, {0, 1}}}},
			{input: "aa", pattern: Pattern{}.Dot("{1}", "c"), count: -1, output: [][][2]int{
				{{0, 1}, {0, 1}},
				{{1, 2}, {1, 2}},
			}},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.FindAllRunesSubmatchIndex([]rune(test.input), test.count),
				c.ShouldEqual,
				test.output,
			)
		}

	})

}

func TestPattern_ReplaceAllRunes(t *testing.T) {

	c.Convey("batch", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			replace Replace[[]rune]
			output  string
		}{

			{ // testing that +? does not match input
				input: "module@1.0.0/thing.txt",
				pattern: Pattern{
					Group("c",
						Text("@"),
						Text("/", "^", "+?"),
						Text("/"),
					),
				},
				replace: Replace[[]rune]{}.WithLiteral([]rune("/")),
				output:  "module@1.0.0/thing.txt",
			},

			{ // testing that + does match input
				input: "module@1.0.0/thing.txt",
				pattern: Pattern{
					Group("c",
						Text("@"),
						Text("/", "^", "+"),
						Text("/"),
					),
				},
				replace: Replace[[]rune]{}.WithLiteral([]rune("/")),
				output:  "module/thing.txt",
			},

			{
				input:   "one",
				pattern: Pattern{}.Dot("+", "c"),
				replace: Replace[[]rune]{}.ToUpper(),
				output:  "ONE",
			},

			{
				input:   "one",
				pattern: Pattern{}.Dot("{1}", "c"),
				replace: Replace[[]rune]{}.ToUpper(),
				output:  "ONE",
			},

			{
				input:   "ONE",
				pattern: Pattern{}.Dot("+", "c"),
				replace: Replace[[]rune]{}.ToLower(),
				output:  "one",
			},

			{input: "", pattern: nil, replace: nil, output: ""},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.ReplaceAllRunes([]rune(test.input), test.replace),
				c.ShouldEqual,
				[]rune(test.output),
			)
		}

	})

}

func TestPattern_ReplaceAllLiteralRunes(t *testing.T) {

	c.Convey("batch", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			replace []rune
			output  string
		}{

			{ // testing that +? does not match input
				input: "module@1.0.0/thing.txt",
				pattern: Pattern{
					Group("c",
						Text("@"),
						Text("/", "^", "+?"),
						Text("/"),
					),
				},
				replace: []rune("/"),
				output:  "module@1.0.0/thing.txt",
			},

			{ // testing that + does match input
				input: "module@1.0.0/thing.txt",
				pattern: Pattern{
					Group("c",
						Text("@"),
						Text("/", "^", "+"),
						Text("/"),
					),
				},
				replace: []rune("/"),
				output:  "module/thing.txt",
			},

			{input: "", pattern: nil, replace: nil, output: ""},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.ReplaceAllLiteralRunes([]rune(test.input), test.replace),
				c.ShouldEqual,
				[]rune(test.output),
			)
		}

	})

}

func TestPattern_ReplaceAllRunesFunc(t *testing.T) {

	c.Convey("batch", t, func() {

		for idx, test := range []struct {
			input     string
			pattern   Pattern
			transform Transform[[]rune]
			output    string
		}{

			{input: "", pattern: nil, output: ""},
			{input: "one two.", pattern: Pattern{}.W("+", "c"), transform: func(input []rune) (output []rune) {
				return []rune(string(input) + "!")
			}, output: "one! two!."},
			{input: "one", pattern: Pattern{}.Dot("+", "c"), transform: func(input []rune) (output []rune) {
				return input
			}, output: "one"},
			{input: "one", pattern: Pattern{}.Dot("{1}", "c"), transform: func(input []rune) (output []rune) {
				return []rune(strings.ToUpper(string(input)))
			}, output: "ONE"},
			{input: "one\ntwo", pattern: Pattern{}.Dot("+", "s", "c"), transform: func(input []rune) (output []rune) {
				return []rune(strings.ToUpper(string(input)))
			}, output: "ONE\nTWO"},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.ReplaceAllRunesFunc([]rune(test.input), test.transform),
				c.ShouldEqual,
				[]rune(test.output),
			)
		}

	})

}

func TestPattern_SplitRunes(t *testing.T) {

	c.Convey("batch", t, func() {

		for idx, test := range []struct {
			input   string
			count   int
			pattern Pattern
			output  [][]rune
		}{

			{input: "", count: 0, pattern: nil, output: [][]rune(nil)},
			{input: "", count: -1, pattern: nil, output: [][]rune(nil)},
			{input: "contents", count: -1, pattern: nil, output: [][]rune(nil)},
			{input: "one two many more", count: 1, pattern: Pattern{}.S("+"), output: [][]rune{
				[]rune("one two many more"),
			}},
			{input: "one two many more", count: 2, pattern: Pattern{}.S("+"), output: [][]rune{
				[]rune("one"),
				[]rune("two many more"),
			}},

			{
				input:   "",
				count:   -1,
				pattern: Pattern{}.Text("a", "*"),
				output:  [][]rune(nil),
			},

			{
				input:   "bb",
				count:   -1,
				pattern: Pattern{}.Text("a", "*"),
				output: [][]rune{
					[]rune("b"),
					[]rune("b"),
				},
			},

			{
				input:   "ababa",
				count:   -1,
				pattern: Pattern{}.Text("a", "*"),
				output: [][]rune{
					[]rune(""),
					[]rune("b"),
					[]rune("b"),
					[]rune(""),
				},
			},

			{
				input:   "abaaba",
				count:   -1,
				pattern: Pattern{}.Text("a", "*"),
				output: [][]rune{
					[]rune(""),
					[]rune("b"),
					[]rune("b"),
					[]rune(""),
				},
			},

			{
				input:   "abaabacc",
				count:   -1,
				pattern: Pattern{}.Text("a", "*"),
				output: [][]rune{
					[]rune(""),
					[]rune("b"),
					[]rune("b"),
					[]rune("c"),
					[]rune("c"),
				},
			},

			{
				input:   "abaabacca",
				count:   -1,
				pattern: Pattern{}.Text("a", "*"),
				output: [][]rune{
					[]rune(""),
					[]rune("b"),
					[]rune("b"),
					[]rune("c"),
					[]rune("c"),
					[]rune(""),
				},
			},

			{
				input:   "abaabacca",
				count:   -1,
				pattern: Pattern{}.Text("a", "+"),
				output: [][]rune{
					[]rune(""),
					[]rune("b"),
					[]rune("b"),
					[]rune("cc"),
					[]rune(""),
				},
			},

			{
				//       v v  v vv vvvvv
				input:   "abaabaccadaaae",
				count:   5,
				pattern: Pattern{}.Text("a", "*"),
				output: [][]rune{
					[]rune(""),
					[]rune("b"),
					[]rune("b"),
					[]rune("c"),
					[]rune("cadaaae"),
				},
			},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q (count=%d)", idx, test.input, test.count),
				test.pattern.SplitRunes([]rune(test.input), test.count),
				c.ShouldEqual,
				test.output,
			)
		}

	})

}
