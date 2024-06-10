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

func TestPattern_Split(t *testing.T) {

	c.Convey("batch", t, func() {

		for idx, test := range []struct {
			input   string
			count   int
			pattern Pattern
			output  []string
		}{

			{input: "", count: 0, pattern: nil, output: []string(nil)},
			{input: "", count: -1, pattern: nil, output: []string(nil)},
			{input: "contents", count: -1, pattern: nil, output: []string(nil)},
			{input: "one two many more", count: 1, pattern: Pattern{}.S("+"), output: []string{
				"one two many more",
			}},
			{input: "one two many more", count: 2, pattern: Pattern{}.S("+"), output: []string{
				"one", "two many more",
			}},

			{
				input:   "",
				count:   -1,
				pattern: Pattern{}.Text("a", "*"),
				output:  []string{""},
			},

			{
				input:   "bb",
				count:   -1,
				pattern: Pattern{}.Text("a", "*"),
				output: []string{
					"b", "b",
				},
			},

			{
				input:   "ababa",
				count:   -1,
				pattern: Pattern{}.Text("a", "*"),
				output: []string{
					"", "b", "b", "",
				},
			},

			{
				input:   "abaaba",
				count:   -1,
				pattern: Pattern{}.Text("a", "*"),
				output: []string{
					"", "b", "b", "",
				},
			},

			{
				input:   "abaabacc",
				count:   -1,
				pattern: Pattern{}.Text("a", "*"),
				output: []string{
					"", "b", "b", "c", "c",
				},
			},

			{
				input:   "abaabacca",
				count:   -1,
				pattern: Pattern{}.Text("a", "*"),
				output: []string{
					"", "b", "b", "c", "c", "",
				},
			},

			{
				input:   "abaabacca",
				count:   -1,
				pattern: Pattern{}.Text("a", "*"),
				output: []string{
					"", "b", "b", "c", "c", "",
				},
			},

			{
				//       v v  v vv vvvvv
				input:   "abaabaccadaaae",
				count:   5,
				pattern: Pattern{}.Text("a", "*"),
				output: []string{
					"", "b", "b", "c", "cadaaae",
				},
			},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q (count=%d)", idx, test.input, test.count),
				test.pattern.Split(test.input, test.count),
				c.ShouldEqual,
				test.output,
			)
		}

	})

}

func TestPattern_MatchString(t *testing.T) {

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
				test.pattern.MatchString(test.input),
				c.ShouldEqual,
				test.ok,
			)
		}

	})

}

func TestPattern_FindString(t *testing.T) {

	c.Convey("batch", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  string
		}{

			{input: "", pattern: nil, output: ""},
			{input: "\naa", pattern: Pattern{}.Dot("{1}", "c"), output: "a"},
			{input: "\naa", pattern: Pattern{}.Dot("{1}", "c"), output: "a"},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.FindString(test.input),
				c.ShouldEqual,
				test.output,
			)
		}

	})

}

func TestPattern_FindIndex(t *testing.T) {

	c.Convey("batch", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  []int
		}{

			{input: "", pattern: nil, output: []int(nil)},
			{input: "aa", pattern: Pattern{}.Dot("{1}", "c"), output: []int{0, 1}},
			{input: "aa", pattern: Pattern{}.Dot("{1}", "c"), output: []int{0, 1}},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.FindIndex(test.input),
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
			output  [][]int
		}{

			{input: "", pattern: nil, count: -1, output: [][]int{{0, 0}}},
			{input: "aa", pattern: Pattern{}.Dot("{1}", "c"), count: 1, output: [][]int{{0, 1}}},
			{input: "aa", pattern: Pattern{}.Dot("{1}", "c"), count: -1, output: [][]int{{0, 1}, {1, 2}}},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.FindAllStringIndex(test.input, test.count),
				c.ShouldEqual,
				test.output,
			)
		}

	})

}

func TestPattern_FindStringSubmatch(t *testing.T) {

	c.Convey("batch", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  []string
		}{

			{input: "", pattern: nil, output: []string(nil)},
			{input: "\naa", pattern: Pattern{}.Dot("{1}", "c"), output: []string{"a", "a"}},
			{input: "a\naa", pattern: Pattern{}.Dot("{1}", "c"), output: []string{"a", "a"}},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.FindStringSubmatch(test.input),
				c.ShouldEqual,
				test.output,
			)
		}

	})

}

func TestPattern_FindAllString(t *testing.T) {

	c.Convey("batch", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			count   int
			output  []string
		}{

			{input: "", pattern: nil, count: -1, output: []string(nil)},
			{input: "aa", pattern: Pattern{}.Dot("{1}", "c"), count: 1, output: []string{"a"}},
			{input: "aa", pattern: Pattern{}.Dot("{1}", "c"), count: -1, output: []string{"a", "a"}},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.FindAllString(test.input, test.count),
				c.ShouldEqual,
				test.output,
			)
		}

	})

}

func TestPattern_FindAllStringSubmatchIndex(t *testing.T) {

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
				test.pattern.FindAllStringSubmatchIndex(test.input, test.count),
				c.ShouldEqual,
				test.output,
			)
		}

	})

}

func TestPattern_ReplaceAllString(t *testing.T) {

	c.Convey("batch", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			replace Replace
			output  string
		}{

			{
				input: "module@1.0.0/thing.txt",
				pattern: Pattern{
					Text("@"),
					Text("/", "^", "+?"),
					Text("/"),
				},
				replace: Replace{}.WithText("/"),
				output:  "module@1.0.0/thing.txt",
			},

			{
				input: "module@1.0.0/thing.txt",
				pattern: Pattern{
					Text("@"),
					Text("/", "^", "+"),
					Text("/"),
				},
				replace: Replace{}.WithText("/"),
				output:  "module/thing.txt",
			},

			{
				input:   "one",
				pattern: Pattern{}.Dot("+", "c"),
				replace: Replace{}.ToUpper(),
				output:  "ONE",
			},

			{
				input:   "one",
				pattern: Pattern{}.Dot("{1}", "c"),
				replace: Replace{}.ToUpper(),
				output:  "ONE",
			},

			{input: "", pattern: nil, replace: nil, output: ""},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.ReplaceAllString(test.input, test.replace),
				c.ShouldEqual,
				test.output,
			)
		}

	})

}

func TestPattern_ReplaceAllStringFunc(t *testing.T) {

	c.Convey("batch", t, func() {

		for idx, test := range []struct {
			input     string
			pattern   Pattern
			transform Transform
			output    string
		}{

			{input: "one", pattern: Pattern{}.Dot("+", "c"), transform: func(input string) (output string) {
				return input
			}, output: "one"},
			{input: "one", pattern: Pattern{}.Dot("{1}", "c"), transform: func(input string) (output string) {
				return strings.ToUpper(input)
			}, output: "ONE"},
			{input: "one\ntwo", pattern: Pattern{}.Dot("+", "s", "c"), transform: func(input string) (output string) {
				return strings.ToUpper(input)
			}, output: "ONE\nTWO"},
			{input: "", pattern: nil, output: ""},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.ReplaceAllStringFunc(test.input, test.transform),
				c.ShouldEqual,
				test.output,
			)
		}

	})

}

func TestPattern_ScanStrings(t *testing.T) {

	c.Convey("batch", t, func() {

		c.Convey("Strings", func() {

			for idx, test := range []struct {
				input   string
				pattern Pattern
				output  []string
			}{

				{input: "", pattern: nil, output: []string{""}},
				{input: "one", pattern: Pattern{}.Dot("+", "c"), output: []string{"one"}},
				{input: "one", pattern: Pattern{}.Dot("{1}", "c"), output: []string{"o", "n", "e"}},
				{input: "one", pattern: Pattern{}.Text("nope", "c"), output: []string{"one"}},
			} {
				c.SoMsg(
					fmt.Sprintf("test #%d - %q", idx, test.input),
					test.pattern.ScanStrings(test.input).Strings(),
					c.ShouldEqual,
					test.output,
				)
			}

		})

		c.Convey("Indexes", func() {

			for idx, test := range []struct {
				input   string
				pattern Pattern
				output  [][]int
			}{

				{input: "", pattern: nil, output: [][]int{{0, 0}}},
				{input: "one", pattern: Pattern{}.Dot("+", "c"), output: [][]int{{0, 3}}},
				{input: "one", pattern: Pattern{}.Dot("{1}", "c"), output: [][]int{{0, 1}, {1, 2}, {2, 3}}},
				{input: "one", pattern: Pattern{}.Text("nope", "c"), output: [][]int{{0, 3}}},
			} {
				c.SoMsg(
					fmt.Sprintf("test #%d - %q", idx, test.input),
					test.pattern.ScanStrings(test.input).Indexes(),
					c.ShouldEqual,
					test.output,
				)
			}

		})

	})

}
