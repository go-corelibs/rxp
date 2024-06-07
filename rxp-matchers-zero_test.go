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

func TestMatchersClassZero(t *testing.T) {

	c.Convey("Caret", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  [][]string
		}{

			{
				input:   "",
				pattern: Pattern{}.Caret(),
				output:  [][]string(nil),
			},

			{
				input:   "a",
				pattern: Pattern{}.Caret("^", "m").Text("a"),
				output:  [][]string(nil),
			},

			{
				input:   "a",
				pattern: Pattern{}.Caret("^").Text("a"),
				output:  [][]string(nil),
			},

			{
				input:   "a\na",
				pattern: Pattern{}.Caret().Text("a", "c"),
				output:  [][]string{{"a", "a"}},
			},

			{
				input:   "a\na",
				pattern: Pattern{}.Caret("m").Text("a", "c"),
				output:  [][]string{{"a", "a"}, {"a", "a"}},
			},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.FindAllStringSubmatch(test.input, -1),
				c.ShouldEqual,
				test.output)
		}
	})

	c.Convey("Dollar", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  [][]string
		}{

			{
				input:   "",
				pattern: Pattern{}.Dollar(),
				output:  [][]string(nil),
			},

			{
				input:   "a",
				pattern: Pattern{}.Text("a").Dollar("^", "m"),
				output:  [][]string(nil),
			},

			{
				input:   "a",
				pattern: Pattern{}.Text("a").Dollar("^"),
				output:  [][]string(nil),
			},

			{
				input:   "a\na",
				pattern: Pattern{}.Text("a", "c").Dollar(),
				output:  [][]string{{"a", "a"}},
			},

			{
				input:   "a\na",
				pattern: Pattern{}.Text("a", "c").Dollar("m"),
				output:  [][]string{{"a", "a"}, {"a", "a"}},
			},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.FindAllStringSubmatch(test.input, -1),
				c.ShouldEqual,
				test.output)
		}
	})

	c.Convey("A", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  [][]string
		}{

			{
				input:   "",
				pattern: Pattern{}.A(),
				output:  [][]string(nil),
			},

			{
				input:   "",
				pattern: Pattern{}.A("^"),
				output:  [][]string(nil),
			},

			{
				input:   "a\nb",
				pattern: Pattern{}.A().Dot("c"),
				output:  [][]string{{"a", "a"}},
			},

			{
				input:   "\nb",
				pattern: Pattern{}.A().Dot("c"),
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

	c.Convey("B", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  [][]string
		}{

			{
				input:   "",
				pattern: Pattern{}.B(),
				output:  [][]string(nil),
			},

			{
				input:   "",
				pattern: Pattern{}.B("^"),
				output:  [][]string(nil),
			},

			{
				input:   "aa",
				pattern: Pattern{}.B().Text("a", "c"),
				output:  [][]string{{"a", "a"}},
			},

			{
				input:   "aa",
				pattern: Pattern{}.Text("a", "c").B(),
				output:  [][]string{{"a", "a"}},
			},

			{
				input:   "aa",
				pattern: Pattern{}.B().Text("a", "c").B(),
				output:  [][]string(nil),
			},

			{
				input:   "b.a.",
				pattern: Pattern{}.B().Text("a", "c").B(),
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

	c.Convey("Z", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  [][]string
		}{

			{
				input:   "",
				pattern: Pattern{}.Z(),
				output:  [][]string(nil),
			},

			{
				input:   "",
				pattern: Pattern{}.Z("^"),
				output:  [][]string(nil),
			},

			{
				input:   "a\na",
				pattern: Pattern{}.Dot("c").Z(),
				output:  [][]string{{"a", "a"}},
			},

			{
				input:   "a\n",
				pattern: Pattern{}.Dot("c").Z(),
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

	c.Convey("BackRef", t, func() {

		c.SoMsg(
			"should panic",
			func() {
				_ = Pattern{}.BackRef(0)
			},
			c.ShouldPanic,
		)

		for idx, test := range []struct {
			inputs  []string
			pattern Pattern
			output  [][][]string
		}{

			{
				inputs:  []string{"aBaBbb"},
				pattern: Pattern{}.Text("B", "c").BackRef(1, "i", "^"),
				output:  [][][]string{{{"Ba", "B"}}},
			},

			{ // this should not match because the backref is oob
				inputs:  []string{"aBabbb"},
				pattern: Pattern{}.Text("b", "i").BackRef(2),
				output:  [][][]string{nil},
			},

			{
				inputs:  []string{"aBabbb"},
				pattern: Pattern{}.Text("b", "{1}", "ic").BackRef(1, "^"),
				output: [][][]string{
					{{"Ba", "B"}, {"b", "b"}},
				},
			},

			{
				inputs: []string{"axaxa", "bxbxb", "cxcxc"}, // ([a-c])x\1x\1 matches axaxa, bxbxb and cxcxc
				pattern: Pattern{}.
					R("a-c", "c").
					Text("x").
					BackRef(1).
					Text("x").
					BackRef(1),
				output: [][][]string{
					{{"axaxa", "a"}},
					{{"bxbxb", "b"}},
					{{"cxcxc", "c"}},
				},
			},
		} {
			for jdx, text := range test.inputs {
				c.SoMsg(
					fmt.Sprintf("test #%d.%d - %q", idx, jdx, text),
					test.pattern.FindAllStringSubmatch(text, -1),
					c.ShouldEqual,
					test.output[jdx])

			}
		}

	})
}
