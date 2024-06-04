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
				fmt.Sprintf("test #%d", idx),
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
				fmt.Sprintf("test #%d", idx),
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
				fmt.Sprintf("test #%d", idx),
				test.pattern.FindAllStringSubmatch(test.input, -1),
				c.ShouldEqual,
				test.output)
		}
	})

	c.Convey("Class", t, func() {

		c.So(func() {
			_ = Class("nope")
		}, c.ShouldPanic)

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  [][]string
		}{

			{
				input:   "",
				pattern: Pattern{}.Class(ASCII),
				output:  [][]string(nil),
			},

			{
				input:   "a",
				pattern: Pattern{}.Class(ASCII, "c"),
				output:  [][]string{{"a", "a"}},
			},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d", idx),
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
				pattern: Pattern{}.Range(unicode.Common),
				output:  [][]string(nil),
			},

			{
				input:   "a",
				pattern: Pattern{}.Range(unicode.Latin, "c"),
				output:  [][]string{{"a", "a"}},
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
