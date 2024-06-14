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

var regexpFindTests = []struct {
	regexp  string
	pattern Pattern
	input   string
	output  [][2]int
}{
	{
		regexp:  ``,
		pattern: Pattern{},
		input:   ``,
		output:  build(0, 0),
	},
	{
		regexp:  `^abcdefg`,
		pattern: Pattern{}.Caret().Text(`abcdefg`),
		input:   "abcdefg",
		output:  build(0, 7),
	},
	{
		regexp:  `a+`,
		pattern: Pattern{}.Text(`a`, "+"),
		input:   "baaab",
		output:  build(1, 4),
	},
	{
		regexp:  `abcd..`,
		pattern: Pattern{}.Text("abcd").Dot().Dot(),
		input:   "abcdef",
		output:  build(0, 6),
	},
	{
		regexp:  `a`,
		pattern: Pattern{}.Text(`a`),
		input:   "a",
		output:  build(0, 1),
	},
	{
		regexp:  `x`,
		pattern: Pattern{}.Text(`x`),
		input:   "y",
		output:  nil,
	},
	{
		regexp:  `b`,
		pattern: Pattern{}.Text(`b`),
		input:   "abc",
		output:  build(1, 2),
	},
	{
		regexp:  `.`,
		pattern: Pattern{}.Dot(),
		input:   "a",
		output:  build(0, 1),
	},
	{
		regexp:  `.*`,
		pattern: Pattern{}.Dot("*"),
		input:   "abcdef",
		output:  build(0, 6),
	},
	{
		regexp:  `^`,
		pattern: Pattern{}.Caret(),
		input:   "abcde",
		output:  build(0, 0),
	},
	{
		regexp:  `$`,
		pattern: Pattern{}.Dollar(),
		input:   "abcde",
		output:  build(5, 5),
	},
	{
		regexp:  `^abcd$`,
		pattern: Pattern{}.Caret().Text(`abcd`).Dollar(),
		input:   "abcd",
		output:  build(0, 4),
	},
	{
		regexp:  `^bcd'`,
		pattern: Pattern{}.Caret().Text(`bcd'`),
		input:   "abcdef",
		output:  nil,
	},
	{
		regexp:  `^abcd$`,
		pattern: Pattern{}.Caret().Text(`abcd`).Dollar(),
		input:   "abcde",
		output:  nil,
	},
	{
		regexp:  `a+`,
		pattern: Pattern{}.Text(`a`, "+"),
		input:   "baaab",
		output:  build(1, 4),
	},
	{
		regexp:  `a*`,
		pattern: Pattern{}.Text(`a`, "*"),
		input:   "baaab",
		output:  build(0, 0, 1, 4, 5, 5),
	},
	{
		regexp:  `[a-z]+`,
		pattern: Pattern{}.R(`a-z`, "+"),
		input:   "abcd",
		output:  build(0, 4),
	},
	{
		regexp:  `[^a-z]+`,
		pattern: Pattern{}.R(`a-z`, "^", "+"),
		input:   "ab1234cd",
		output:  build(2, 6),
	},
	{
		regexp:  `[a\-\]z]+`,
		pattern: Pattern{}.R(`-a]z`, "+"),
		input:   "az]-bcz",
		output:  build(0, 4, 6, 7),
	},
	{
		regexp:  `[^\n]+`,
		pattern: Pattern{}.R("\n", "^", "+"),
		input:   "abcd\n",
		output:  build(0, 4),
	},
	{
		regexp:  `[日本語]+`,
		pattern: Pattern{}.R(`日本語`, "+"),
		input:   "日本語日本語",
		output:  build(0, 18),
	},
	{
		regexp:  `日本語+`,
		pattern: Pattern{}.Text(`日`).Text(`本`).Text(`語`, "+"),
		input:   "日本語",
		output:  build(0, 9),
	},
	{
		regexp:  `日本語+`,
		pattern: Pattern{}.Text(`日`).Text(`本`).Text(`語`, "+"),
		input:   "日本語語語語",
		output:  build(0, 18),
	},

	//{
	//	regexp:  `()`,
	//	pattern: Pattern{}.Group("c"),
	//	input:   "",
	//	output:  build(1, 0, 0, 0, 0),
	//},
	//{
	//	regexp:  `(a)`,
	//	pattern: Pattern{}.Group("c", Text(`a`)),
	//	input:   "a",
	//	output:  build(1, 0, 1, 0, 1),
	//},
	//{
	//	regexp:  `(.)(.)`,
	//	pattern: Pattern{}.Group(Dot()).Group(Dot()),
	//	input:   "日a",
	//	output:  build(1, 0, 4, 0, 3, 3, 4),
	//},
	//{
	//	regexp:  `(.*)`,
	//	pattern: Pattern{}.Group("c", Dot("*")),
	//	input:   "",
	//	output:  build(1, 0, 0, 0, 0),
	//},
	//{
	//	regexp:  `(.*)`,
	//	pattern: Pattern{}.Group("c", Dot("*")),
	//	input:   "abcd",
	//	output:  build(1, 0, 4, 0, 4),
	//},
	//{
	//	regexp:  `(..)(..)`,
	//	pattern: Pattern{}.Group("c", Dot(), Dot()).Group("c", Dot(), Dot()),
	//	input:   "abcd",
	//	output:  build(1, 0, 4, 0, 2, 2, 4),
	//},

	//{ // (([^xyz]*)(d))
	//	// sub-sub groups not supported yet!
	//	pattern: Pattern{}.Text(`(([^xyz]*)(d))`),
	//	input:   "abcd",
	//	output:  build(1, 0, 4, 0, 4, 0, 3, 3, 4),
	//},
	//{
	//	pattern: Pattern{}.Text(`((a|b|c)*(d))`),
	//	input:   "abcd",
	//	output:  build(1, 0, 4, 0, 4, 2, 3, 3, 4),
	//},
	//{
	//	pattern: Pattern{}.Text(`(((a|b|c)*)(d))`),
	//	input:   "abcd",
	//	output:  build(1, 0, 4, 0, 4, 0, 3, 2, 3, 3, 4),
	//},

	//{
	//	pattern: Pattern{}.Text("\a\f\n\r\t\v"),
	//	input:   "\a\f\n\r\t\v",
	//	output:  build(1, 0, 6),
	//},
	//{
	//	pattern: Pattern{}.Text(`[\a\f\n\r\t\v]+`),
	//	input:   "\a\f\n\r\t\v",
	//	output:  build(1, 0, 6),
	//},

	//{
	//	pattern: Pattern{}.Text(`a*(|(b))c*`),
	//	input:   "aacc",
	//	output:  build(1, 0, 4, 2, 2, -1, -1),
	//},
	//{
	//	pattern: Pattern{}.Text(`(.*).*`),
	//	input:   "ab",
	//	output:  build(1, 0, 2, 0, 2),
	//},

	{
		regexp:  `[.]`, // 23
		pattern: Pattern{}.R("."),
		input:   ".",
		output:  build(0, 1),
	},
	{
		regexp:  `/$`,
		pattern: Pattern{}.Text(`/`).Dollar(),
		input:   "/abc/",
		output:  build(4, 5),
	},
	{
		regexp:  `/$`,
		pattern: Pattern{}.Text(`/`).Dollar(),
		input:   "/abc",
		output:  nil,
	},

	// multiple matches
	{
		regexp:  `.`,
		pattern: Pattern{}.Dot(),
		input:   "abc",
		output:  build(0, 1, 1, 2, 2, 3),
	},
	//{
	//	pattern: Pattern{}.Text(`(.)`),
	//	input:   "abc",
	//	output:  build(3, 0, 1, 0, 1, 1, 2, 1, 2, 2, 3, 2, 3),
	//},
	//{
	//	pattern: Pattern{}.Text(`.(.)`),
	//	input:   "abcd",
	//	output:  build(2, 0, 2, 1, 2, 2, 4, 3, 4),
	//},
	{
		regexp:  `ab*`, //27
		pattern: Pattern{}.Text(`a`).Text(`b`, "*"),
		input:   "abbaab",
		output:  build(0, 3, 3, 4, 4, 6),
	},
	//{
	//	pattern: Pattern{}.Text(`a(b*)`),
	//	input:   "abbaab",
	//	output:  build(3, 0, 3, 1, 3, 3, 4, 4, 4, 4, 6, 5, 6),
	//},
	//
	// fixed bugs
	{
		regexp:  `ab$`,
		pattern: Pattern{}.Text(`ab`).Dollar(),
		input:   "cab",
		output:  build(1, 3),
	},
	{
		regexp:  `axxb$`,
		pattern: Pattern{}.Text(`axxb`).Dollar(),
		input:   "axxcb",
		output:  nil,
	},
	{
		regexp:  `data`,
		pattern: Pattern{}.Text(`data`),
		input:   "daXY data",
		output:  build(5, 9),
	},
	//{
	//	pattern: Pattern{}.Text(`da(.)a$`),
	//	input:   "daXY data",
	//	output:  build(1, 5, 9, 7, 8),
	//},
	{
		regexp:  `zx+`,
		pattern: Pattern{}.Text(`z`).Text(`x`, "+"),
		input:   "zzx",
		output:  build(1, 3),
	},
	{
		regexp:  `ab$`,
		pattern: Pattern{}.Text(`ab`).Dollar(),
		input:   "abcab",
		output:  build(3, 5),
	},
	//{
	//	pattern: Pattern{}.Text(`(aa)*$`),
	//	input:   "a",
	//	output:  build(1, 1, 1, -1, -1),
	//},
	//{
	//	pattern: Pattern{}.Text(`(?:.|(?:.a))`),
	//	input:   "",
	//	output:  nil,
	//},
	//{
	//	pattern: Pattern{}.Text(`(?:A(?:A|a))`),
	//	input:   "Aa",
	//	output:  build(1, 0, 2),
	//},
	//{
	//	pattern: Pattern{}.Text(`(?:A|(?:A|a))`),
	//	input:   "a",
	//	output:  build(1, 0, 1),
	//},
	//{
	//	pattern: Pattern{}.Text(`(a){0}`),
	//	input:   "",
	//	output:  build(1, 0, 0, -1, -1),
	//},
	//{
	//	pattern: Pattern{}.Text(`(?-s)(?:(?:^).)`),
	//	input:   "\n",
	//	output:  nil,
	//},
	//{
	//	pattern: Pattern{}.Text(`(?s)(?:(?:^).)`),
	//	input:   "\n",
	//	output:  build(1, 0, 1),
	//},
	//{
	//	pattern: Pattern{}.Text(`(?:(?:^).)`),
	//	input:   "\n",
	//	output:  nil,
	//},
	{
		regexp:  `\b`, // 33
		pattern: Pattern{}.B(),
		input:   "x",
		output:  build(0, 0, 1, 1),
	},
	{
		regexp:  `\b`,
		pattern: Pattern{}.B(),
		input:   "xx",
		output:  build(0, 0, 2, 2),
	},
	{
		regexp:  `\b`,
		pattern: Pattern{}.B(),
		input:   "x y",
		output:  build(0, 0, 1, 1, 2, 2, 3, 3),
	},
	{
		regexp:  `\b`,
		pattern: Pattern{}.B(),
		input:   "xx yy",
		output:  build(0, 0, 2, 2, 3, 3, 5, 5),
	},
	{
		regexp:  `\b`,
		pattern: Pattern{}.B("^"),
		input:   "x",
		output:  nil,
	},
	{
		regexp:  `\b`,
		pattern: Pattern{}.B("^"),
		input:   "xx",
		output:  build(1, 1),
	},
	{
		regexp:  `\b`,
		pattern: Pattern{}.B("^"),
		input:   "x y",
		output:  nil,
	},
	{
		regexp:  `\b`,
		pattern: Pattern{}.B("^"),
		input:   "xx yy",
		output:  build(1, 1, 4, 4),
	},
	//{
	//	pattern: Pattern{}.Text(`(|a)*`),
	//	input:   "aa",
	//	output:  build(3, 0, 0, 0, 0, 1, 1, 1, 1, 2, 2, 2, 2),
	//},

	// RE2 tests
	{
		regexp:  `[^\S\s]`,
		pattern: Pattern{}.Group(S("^"), S(), "^"),
		input:   "abcd",
		output:  nil,
	},
	{
		regexp:  `[^\S[:space:]]`,
		pattern: Pattern{}.Text(`[^\S[:space:]]`),
		input:   "abcd",
		output:  nil,
	},
	{
		regexp:  `[^\D\d]`,
		pattern: Pattern{}.Text(`[^\D\d]`),
		input:   "abcd",
		output:  nil,
	},
	{
		regexp:  `[^\D[:digit:]]`,
		pattern: Pattern{}.Text(`[^\D[:digit:]]`),
		input:   "abcd",
		output:  nil,
	},
	{
		regexp:  `(?i)\W`,
		pattern: Pattern{}.Text(`(?i)\W`),
		input:   "x",
		output:  nil,
	},
	{
		regexp:  `(?i)\W`,
		pattern: Pattern{}.Text(`(?i)\W`),
		input:   "k",
		output:  nil,
	},
	{
		regexp:  `(?i)\W`,
		pattern: Pattern{}.Text(`(?i)\W`),
		input:   "s",
		output:  nil,
	},

	//// can backslash-escape any punctuation
	//{
	//	pattern: Pattern{}.Text(`\!\"\#\$\%\&\'\(\)\*\+\,\-\.\/\:\;\<\=\>\?\@\[\\\]\^\_\{\|\}\~`),
	//	input:   `!"#$%&'()*+,-./:;<=>?@[\]^_{|}~`,
	//	output:  build(1, 0, 31),
	//},
	//{
	//	pattern: Pattern{}.Text(`[\!\"\#\$\%\&\'\(\)\*\+\,\-\.\/\:\;\<\=\>\?\@\[\\\]\^\_\{\|\}\~]+`),
	//	input:   `!"#$%&'()*+,-./:;<=>?@[\]^_{|}~`,
	//	output:  build(1, 0, 31),
	//},
	{
		regexp:  "`",
		pattern: Pattern{}.Text("`"),
		input:   "`",
		output:  build(0, 1),
	},
	{
		regexp:  "[`]+",
		pattern: Pattern{}.R("`", "+"),
		input:   "`",
		output:  build(0, 1),
	},

	//{
	//	pattern: Pattern{}.Text("\ufffd"),
	//	input:   "\xff",
	//	output:  build(1, 0, 1),
	//},
	//{
	//	pattern: Pattern{}.Text("\ufffd"),
	//	input:   "hello\xffworld",
	//	output:  build(1, 5, 6),
	//},
	{
		regexp:  `.*`,
		pattern: Pattern{}.Dot("*"),
		input:   "hello\xffworld",
		output:  build(0, 11),
	},
	//{
	//	pattern: Pattern{}.Text(`\x{fffd}`),
	//	input:   "\xc2\x00",
	//	output:  build(1, 0, 1),
	//},
	//{
	//	pattern: Pattern{}.Text("[\ufffd]"),
	//	input:   "\xff",
	//	output:  build(1, 0, 1),
	//},
	//{
	//	pattern: Pattern{}.Text(`[\x{fffd}]`),
	//	input:   "\xc2\x00",
	//	output:  build(1, 0, 1),
	//},

	// long set of matches (longer than startSize)
	{ // .
		regexp:  `.`,
		pattern: Pattern{}.Dot(),
		input:   "qwertyuiopasdfghjklzxcvbnm1234567890",
		output:  build(0, 1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7, 8, 8, 9, 9, 10, 10, 11, 11, 12, 12, 13, 13, 14, 14, 15, 15, 16, 16, 17, 17, 18, 18, 19, 19, 20, 20, 21, 21, 22, 22, 23, 23, 24, 24, 25, 25, 26, 26, 27, 27, 28, 28, 29, 29, 30, 30, 31, 31, 32, 32, 33, 33, 34, 34, 35, 35, 36),
	},
}

// build is a helper to construct a [][2]int by extracting n sequences from x.
// This represents n matches with len(x)/n sub-matches each.
//
// build is borrowed from regexp v1.22.4
func build(pairs ...int) [][2]int {
	size := len(pairs)
	if size%2 != 0 {
		panic("qa developer error, build requires an even number of start/end index arguments")
	}
	count := size / 2
	ret := make([][2]int, count)
	jdx := 0
	for idx := 0; idx < size; idx += 2 {
		ret[jdx] = [2]int{pairs[idx], pairs[idx+1]}
		jdx += 1
	}
	return ret
}

func TestRegexp(t *testing.T) {
	c.Convey("Find Tests", t, func() {
		for idx, test := range regexpFindTests {

			c.SoMsg(
				fmt.Sprintf("test #%d | %q | %q", idx, test.regexp, test.input),
				test.pattern.FindAllStringIndex(test.input, -1),
				c.ShouldEqual,
				test.output,
			)

		}
	})
}

func TestVersusRegexp(t *testing.T) {
	c.Convey(`a*`, t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  []string
		}{
			{
				input:   `bb`,
				pattern: Pattern{}.Text("a", "*"),
				output: []string{
					"", "", "",
				},
			},
			{
				input:   `abaabaccadaaae`,
				pattern: Pattern{}.Text("a", "*"),
				output: []string{
					"a", "aa", "a", "", "a", "aaa", "",
				},
			},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				Pattern{}.
					Text("a", "*").
					FindAllString(test.input, -1),
				c.ShouldEqual,
				test.output,
			)
		}

	})
}
