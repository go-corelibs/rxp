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
	"unicode"

	c "github.com/smartystreets/goconvey/convey"
)

func TestMatchersClassPerl(t *testing.T) {

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
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.FindAllStringSubmatch(test.input, -1),
				c.ShouldEqual,
				test.output)
		}

	})

	c.Convey("Dot", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  [][]string
		}{

			// rxInvalidFuncName = regexp.MustCompile(`^\s*(\d+|func\d+)\s*$`)
			{
				input: " func1  ",
				pattern: Pattern{}.
					Caret().
					S("*").
					Or("c",
						D("+"),
						Group(
							Text("func"),
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
				pattern: Pattern{}.Add(func(scope Flags, reps Reps, input *RuneBuffer, index int, sm SubMatches) (consumed int, captured bool, negated bool, proceed bool) {
					// not capturing on purpose
					this, okt := input.Get(index)
					if proceed = okt && this == '@'; proceed {
						consumed = 1

						total := input.Len()
						start := index

						for idx := start + consumed; idx < total; idx += 1 {
							consumed += 1
							if r, _ := input.Get(idx); r == '/' {
								break
							}
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

			{
				input:   "abc\nn\no\np\ne",
				pattern: Pattern{}.Dot("m", "+", "c"),
				output:  [][]string{{"abc", "abc"}, {"n", "n"}, {"o", "o"}, {"p", "p"}, {"e", "e"}},
			},

			{
				input:   "1",
				pattern: Pattern{}.D().Dot(),
				output:  [][]string(nil),
			},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.FindAllStringSubmatch(test.input, -1),
				c.ShouldEqual,
				test.output,
			)
		}

	})

	c.Convey("D", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  [][]string
		}{

			{
				input:   "a11",
				pattern: Pattern{}.D("+", "c").B(),
				output:  [][]string{{"11", "11"}},
			},

			{
				input:   "",
				pattern: Pattern{}.D(),
				output:  [][]string(nil),
			},

			{
				input:   "a1a",
				pattern: Pattern{}.D("c"),
				output:  [][]string{{"1", "1"}},
			},

			{
				input:   "aa",
				pattern: Pattern{}.D("c"),
				output:  [][]string(nil),
			},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.FindAllStringSubmatch(test.input, -1),
				c.ShouldEqual,
				test.output)
		}
	})

	c.Convey("S", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  [][]string
		}{

			{
				input:   "a11",
				pattern: Pattern{}.S("+", "c").B(),
				output:  [][]string(nil),
			},

			{
				input:   "",
				pattern: Pattern{}.S(),
				output:  [][]string(nil),
			},

			{
				input:   "a a",
				pattern: Pattern{}.S("c"),
				output:  [][]string{{" ", " "}},
			},

			{
				input:   "a \t a",
				pattern: Pattern{}.S("+", "c"),
				output:  [][]string{{" \t ", " \t "}},
			},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.FindAllStringSubmatch(test.input, -1),
				c.ShouldEqual,
				test.output)
		}
	})

	c.Convey("W", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  [][]string
		}{

			{
				input:   "",
				pattern: Pattern{},
				output:  [][]string(nil),
			},

			{
				input:   "#!/usr/bin/env perl",
				pattern: Pattern{}.W("+", "c"),
				output: [][]string{
					{"usr", "usr"},
					{"bin", "bin"},
					{"env", "env"},
					{"perl", "perl"},
				},
			},

			{
				input:   "one",
				pattern: Pattern{}.W("", "c"),
				output: [][]string{
					{"o", "o"},
					{"n", "n"},
					{"e", "e"},
				},
			},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.FindAllStringSubmatch(test.input, -1),
				c.ShouldEqual,
				test.output)
		}
	})

	c.Convey("NamedClass", t, func() {

		c.So(func() {
			_ = NamedClass("nope")
		}, c.ShouldPanic)

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  [][]string
		}{

			{
				input:   "",
				pattern: Pattern{}.NamedClass(ASCII),
				output:  [][]string(nil),
			},

			{
				input:   "a",
				pattern: Pattern{}.NamedClass(ASCII, "c"),
				output:  [][]string{{"a", "a"}},
			},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.FindAllStringSubmatch(test.input, -1),
				c.ShouldEqual,
				test.output)
		}
	})

	c.Convey("IsUnicodeRange", t, func() {

		c.So(func() {
			_ = IsUnicodeRange(nil)
		}, c.ShouldPanic)

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  [][]string
		}{

			{
				input:   "",
				pattern: Pattern{}.RangeTable(unicode.Common),
				output:  [][]string(nil),
			},

			{
				input:   "a",
				pattern: Pattern{}.RangeTable(unicode.Latin, "c"),
				output:  [][]string{{"a", "a"}},
			},

			{
				input:   "a",
				pattern: Pattern{}.RangeTable(unicode.Latin, "^", "c"),
				output:  [][]string(nil),
			},

			{
				input:   "a",
				pattern: Pattern{}.Text("a").RangeTable(unicode.Latin, "c"),
				output:  [][]string(nil),
			},

			{
				input:   "日本語",
				pattern: Pattern{}.RangeTable(unicode.Latin, "c"),
				output:  [][]string(nil),
			},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.FindAllStringSubmatch(test.input, -1),
				c.ShouldEqual,
				test.output)
		}
	})

	c.Convey("R", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  []string
		}{

			{
				input:   `Test x, y and z`,
				pattern: Pattern{}.R("[xyza-f]", "c"),
				output:  []string{"e", "x", "y", "a", "d", "z"},
			},

			{
				input:   `[x-z]`,
				pattern: Pattern{}.R("-[xyz]", "c"),
				output:  []string{"[", "x", "-", "z", "]"},
			},

			{ // this isn't regexp, the square brackets are just runes but the
				// dash is not the first key and so that's the range of [ to x!
				input:   `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXY1234567890-_=+[]{}()*&^%$#@!`,
				pattern: Pattern{}.R("[-x!]", "c"),
				output: []string{
					"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x",
					"_", "[", "]", "^", "!",
				}, // there is no dash in this output and because the value of [ is less than the value of x...
			},

			{ // place the literal dash at the start to match dashes
				input:   `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXY1234567890-_=+[]{}()*&^%$#@!`,
				pattern: Pattern{}.R("-[x!]", "c"),
				output:  []string{"x", "-", "[", "]", "!"},
			},

			{ // place the literal dash at the start to match dashes
				input:   `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXY1234567890-_=+[]{}()*&^%$#@!`,
				pattern: Pattern{}.R("c-a", "c"),
				output:  []string{"a", "b", "c"},
			},
		} {

			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.FindAllString(test.input, -1),
				c.ShouldEqual,
				test.output,
			)

		}

	})

	c.Convey("Alpha", t, func() {

		c.So(
			Pattern{}.Alpha().FindAllString(`01a`, -1),
			c.ShouldEqual,
			[]string{"a"},
		)

	})

	c.Convey("Ascii", t, func() {

		c.So(
			Pattern{}.Ascii().FindAllString("日本語eh", -1),
			c.ShouldEqual,
			[]string{"e", "h"},
		)

	})

	c.Convey("Blank", t, func() {

		c.So(
			Pattern{}.Blank().FindAllString("space tab\tnewline\n", -1),
			c.ShouldEqual,
			[]string{" ", "\t"},
		)

	})

	c.Convey("Cntrl", t, func() {

		c.So(
			Pattern{}.Cntrl().FindAllString("alarm\a", -1),
			c.ShouldEqual,
			[]string{"\a"},
		)

	})

	c.Convey("Digit", t, func() {

		c.So(
			Pattern{}.Digit().FindAllString("one 1৩", -1),
			c.ShouldEqual,
			[]string{"1"},
		)

	})

	c.Convey("Graph", t, func() {

		c.So(
			Pattern{}.Graph().FindAllString("o৩", -1),
			c.ShouldEqual,
			[]string{"o"},
		)

	})

	c.Convey("Lower", t, func() {

		c.So(
			Pattern{}.Lower().FindAllString("oNE৩", -1),
			c.ShouldEqual,
			[]string{"o"},
		)

	})

	c.Convey("Print", t, func() {

		c.So(
			Pattern{}.Print().FindAllString(" o\t", -1),
			c.ShouldEqual,
			[]string{" ", "o"},
		)

	})

	c.Convey("Punct", t, func() {

		c.So(
			Pattern{}.Punct().FindAllString("pun'c\t", -1),
			c.ShouldEqual,
			[]string{"'"},
		)

	})

	c.Convey("Space", t, func() {

		c.So(
			Pattern{}.Space().FindAllString("tab\tnl\nvtab\vftab\fcr\rspace ", -1),
			c.ShouldEqual,
			[]string{"\t", "\n", "\v", "\f", "\r", " "},
		)

	})

	c.Convey("Upper", t, func() {

		c.So(
			Pattern{}.Upper().FindAllString("oNE৩", -1),
			c.ShouldEqual,
			[]string{"N", "E"},
		)

	})

	c.Convey("Word", t, func() {

		c.So(
			Pattern{}.Word().FindAllString("_w'o,r.d_", -1),
			c.ShouldEqual,
			[]string{"_", "w", "o", "r", "d", "_"},
		)

	})

	c.Convey("Xdigit", t, func() {

		c.So(
			Pattern{}.Xdigit().FindAllString("_a-fvgl1", -1),
			c.ShouldEqual,
			[]string{"a", "f", "1"},
		)

	})

}
