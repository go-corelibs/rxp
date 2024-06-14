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

func TestPattern_Match(t *testing.T) {

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
				test.pattern.MatchBytes([]byte(test.input)),
				c.ShouldEqual,
				test.ok,
			)
		}

	})

}

func TestPattern_Find(t *testing.T) {

	c.Convey("batch", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  []byte
		}{

			{input: "", pattern: nil, output: []byte(nil)},
			{input: "\naa", pattern: Pattern{}.Dot("{1}", "c"), output: []byte("a")},
			{input: "\naa", pattern: Pattern{}.Dot("{1}", "c"), output: []byte("a")},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.FindBytes([]byte(test.input)),
				c.ShouldEqual,
				[]byte(test.output),
			)
		}

	})

}

func TestPattern_FindIndex(t *testing.T) {

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
				test.pattern.FindBytesIndex([]byte(test.input)),
				c.ShouldEqual,
				test.output,
			)
		}

	})

}

func TestPattern_FindAllIndex(t *testing.T) {

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
				test.pattern.FindAllBytesIndex([]byte(test.input), test.count),
				c.ShouldEqual,
				test.output,
			)
		}

	})

}

func TestPattern_FindSubmatch(t *testing.T) {

	c.Convey("batch", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  [][]byte
		}{

			{input: "", pattern: nil, output: [][]byte(nil)},
			{input: "\naa", pattern: Pattern{}.Dot("{1}", "c"), output: [][]byte{[]byte("a"), []byte("a")}},
			{input: "a\naa", pattern: Pattern{}.Dot("{1}", "c"), output: [][]byte{[]byte("a"), []byte("a")}},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.FindBytesSubmatch([]byte(test.input)),
				c.ShouldEqual,
				test.output,
			)
		}

	})

}

func TestPattern_FindAll(t *testing.T) {

	c.Convey("batch", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			count   int
			output  [][]byte
		}{

			{input: "", pattern: nil, count: -1, output: [][]byte(nil)},
			{input: "aa", pattern: Pattern{}.Dot("{1}", "c"), count: 1, output: [][]byte{[]byte("a")}},
			{input: "aa", pattern: Pattern{}.Dot("{1}", "c"), count: -1, output: [][]byte{[]byte("a"), []byte("a")}},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.FindAllBytes([]byte(test.input), test.count),
				c.ShouldEqual,
				test.output,
			)
		}

	})

}

func TestPattern_FindAllSubmatch(t *testing.T) {

	c.Convey("batch", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			count   int
			output  [][][]byte
		}{

			{input: "", pattern: nil, count: -1, output: [][][]byte(nil)},
			{input: "aa", pattern: Pattern{}.Dot("{1}", "c"), count: 1, output: [][][]byte{{{97}, {97}}}},
			{input: "aa", pattern: Pattern{}.Dot("{1}", "c"), count: -1, output: [][][]byte{{{97}, {97}}, {{97}, {97}}}},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.FindAllBytesSubmatch([]byte(test.input), test.count),
				c.ShouldEqual,
				test.output,
			)
		}

	})

}

func TestPattern_FindAllSubmatchIndex(t *testing.T) {

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
				test.pattern.FindAllBytesSubmatchIndex([]byte(test.input), test.count),
				c.ShouldEqual,
				test.output,
			)
		}

	})

}

func TestPattern_ReplaceAll(t *testing.T) {

	c.Convey("batch", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			replace Replace[[]byte]
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
				replace: Replace[[]byte]{}.WithLiteral([]byte("/")),
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
				replace: Replace[[]byte]{}.WithLiteral([]byte("/")),
				output:  "module/thing.txt",
			},

			{
				input:   "one",
				pattern: Pattern{}.Dot("+", "c"),
				replace: Replace[[]byte]{}.ToUpper(),
				output:  "ONE",
			},

			{
				input:   "one",
				pattern: Pattern{}.Dot("{1}", "c"),
				replace: Replace[[]byte]{}.ToUpper(),
				output:  "ONE",
			},

			{
				input:   "ONE",
				pattern: Pattern{}.Dot("+", "c"),
				replace: Replace[[]byte]{}.ToLower(),
				output:  "one",
			},

			{input: "", pattern: nil, replace: nil, output: ""},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.ReplaceAllBytes([]byte(test.input), test.replace),
				c.ShouldEqual,
				[]byte(test.output),
			)
		}

	})

}

func TestPattern_ReplaceAllLiteral(t *testing.T) {

	c.Convey("batch", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			replace []byte
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
				replace: []byte("/"),
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
				replace: []byte("/"),
				output:  "module/thing.txt",
			},

			{input: "", pattern: nil, replace: nil, output: ""},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.ReplaceAllLiteralBytes([]byte(test.input), test.replace),
				c.ShouldEqual,
				[]byte(test.output),
			)
		}

	})

}

func TestPattern_ReplaceAllFunc(t *testing.T) {

	c.Convey("batch", t, func() {

		for idx, test := range []struct {
			input     string
			pattern   Pattern
			transform Transform[[]byte]
			output    string
		}{

			{input: "", pattern: nil, output: ""},
			{input: "one two.", pattern: Pattern{}.W("+", "c"), transform: func(input []byte) (output []byte) {
				return []byte(string(input) + "!")
			}, output: "one! two!."},
			{input: "one", pattern: Pattern{}.Dot("+", "c"), transform: func(input []byte) (output []byte) {
				return input
			}, output: "one"},
			{input: "one", pattern: Pattern{}.Dot("{1}", "c"), transform: func(input []byte) (output []byte) {
				return []byte(strings.ToUpper(string(input)))
			}, output: "ONE"},
			{input: "one\ntwo", pattern: Pattern{}.Dot("+", "s", "c"), transform: func(input []byte) (output []byte) {
				return []byte(strings.ToUpper(string(input)))
			}, output: "ONE\nTWO"},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.ReplaceAllBytesFunc([]byte(test.input), test.transform),
				c.ShouldEqual,
				[]byte(test.output),
			)
		}

	})

}

func TestPattern_SplitBytes(t *testing.T) {

	c.Convey("batch", t, func() {

		for idx, test := range []struct {
			input   string
			count   int
			pattern Pattern
			output  [][]byte
		}{

			{input: "", count: 0, pattern: nil, output: [][]byte(nil)},
			{input: "", count: -1, pattern: nil, output: [][]byte(nil)},
			{input: "contents", count: -1, pattern: nil, output: [][]byte(nil)},
			{input: "one two many more", count: 1, pattern: Pattern{}.S("+"), output: [][]byte{
				[]byte("one two many more"),
			}},
			{input: "one two many more", count: 2, pattern: Pattern{}.S("+"), output: [][]byte{
				[]byte("one"),
				[]byte("two many more"),
			}},

			{
				input:   "",
				count:   -1,
				pattern: Pattern{}.Text("a", "*"),
				output:  [][]byte(nil),
			},

			{
				input:   "bb",
				count:   -1,
				pattern: Pattern{}.Text("a", "*"),
				output: [][]byte{
					[]byte("b"),
					[]byte("b"),
				},
			},

			{
				input:   "ababa",
				count:   -1,
				pattern: Pattern{}.Text("a", "*"),
				output: [][]byte{
					[]byte(""),
					[]byte("b"),
					[]byte("b"),
					[]byte(""),
				},
			},

			{
				input:   "abaaba",
				count:   -1,
				pattern: Pattern{}.Text("a", "*"),
				output: [][]byte{
					[]byte(""),
					[]byte("b"),
					[]byte("b"),
					[]byte(""),
				},
			},

			{
				input:   "abaabacc",
				count:   -1,
				pattern: Pattern{}.Text("a", "*"),
				output: [][]byte{
					[]byte(""),
					[]byte("b"),
					[]byte("b"),
					[]byte("c"),
					[]byte("c"),
				},
			},

			{
				input:   "abaabacca",
				count:   -1,
				pattern: Pattern{}.Text("a", "*"),
				output: [][]byte{
					[]byte(""),
					[]byte("b"),
					[]byte("b"),
					[]byte("c"),
					[]byte("c"),
					[]byte(""),
				},
			},

			{
				input:   "abaabacca",
				count:   -1,
				pattern: Pattern{}.Text("a", "+"),
				output: [][]byte{
					[]byte(""),
					[]byte("b"),
					[]byte("b"),
					[]byte("cc"),
					[]byte(""),
				},
			},

			{
				//       v v  v vv vvvvv
				input:   "abaabaccadaaae",
				count:   5,
				pattern: Pattern{}.Text("a", "*"),
				output: [][]byte{
					[]byte(""),
					[]byte("b"),
					[]byte("b"),
					[]byte("c"),
					[]byte("cadaaae"),
				},
			},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q (count=%d)", idx, test.input, test.count),
				test.pattern.SplitBytes([]byte(test.input), test.count),
				c.ShouldEqual,
				test.output,
			)
		}

	})

}
