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

func TestMatchersCommon(t *testing.T) {

	c.Convey("FieldWord", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  []string
		}{
			{
				input:   "a' ",
				pattern: Pattern{}.Add(FieldWord("c")),
				output:  []string{"a", "' "},
			},

			{
				input:   "a",
				pattern: Pattern{}.Add(FieldWord("c")),
				output:  []string{"a"},
			},

			{
				input:   "aa",
				pattern: Pattern{}.Add(FieldWord("c")),
				output:  []string{"aa"},
			},

			{
				input:   "one two won't do",
				pattern: Pattern{}.Add(FieldWord("c")),
				output:  []string{"one", " ", "two", " ", "won't", " ", "do"},
			},

			{
				input:   "-one two won't do",
				pattern: Pattern{}.Add(FieldWord("c")),
				output:  []string{"-", "one", " ", "two", " ", "won't", " ", "do"},
			},

			{
				input:   "-one-two  won't do",
				pattern: Pattern{}.Add(FieldWord("c")),
				output:  []string{"-", "one", "-", "two", "  ", "won't", " ", "do"},
			},

			{
				input:   "-one-two'  won't do",
				pattern: Pattern{}.Add(FieldWord("c")),
				output:  []string{"-", "one", "-", "two", "'  ", "won't", " ", "do"},
			},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d", idx),
				test.pattern.ScanString(test.input).Strings(),
				c.ShouldEqual,
				test.output)
		}

		for idx, test := range []struct {
			input   string
			pattern Pattern
		}{
			{
				input:   "a",
				pattern: Pattern{}.Add(FieldWord("c")),
			},

			{
				input:   "aa",
				pattern: Pattern{}.Add(FieldWord("c")),
			},

			{
				input:   "a-a",
				pattern: Pattern{}.Add(FieldWord("c")),
			},

			{
				input:   "a-a_a",
				pattern: Pattern{}.Add(FieldWord("c")),
			},

			{
				input:   "0-a_a",
				pattern: Pattern{}.Add(FieldWord("c")),
			},

			{
				input:   "-a-a",
				pattern: Pattern{}.Add(FieldWord("c")),
			},

			{
				input:   "+a-a",
				pattern: Pattern{}.Add(FieldWord("c")),
			},

			{
				input:   "one two won't do",
				pattern: Pattern{}.Add(FieldWord("c")),
			},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d", idx),
				test.pattern.ScanString(test.input).String(),
				c.ShouldEqual,
				test.input)
		}

	})

	c.Convey("FieldKey", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  [][]string
		}{
			{
				input:   "a",
				pattern: Pattern{}.Add(FieldKey("c")),
				output:  [][]string{{"a", "a"}},
			},

			{
				input:   "aa",
				pattern: Pattern{}.Add(FieldKey("c")),
				output:  [][]string{{"aa", "aa"}},
			},

			{
				input:   "a-a",
				pattern: Pattern{}.Add(FieldKey("c")),
				output:  [][]string{{"a-a", "a-a"}},
			},

			{
				input:   "a-a_a",
				pattern: Pattern{}.Add(FieldKey("c")),
				output:  [][]string{{"a-a_a", "a-a_a"}},
			},

			{
				input:   "0-a_a",
				pattern: Pattern{}.Add(FieldKey("c")),
				output:  [][]string{{"a_a", "a_a"}},
			},

			{
				input:   "-a-a",
				pattern: Pattern{}.Add(FieldKey("c")),
				output:  [][]string{{"a-a", "a-a"}},
			},

			{
				input:   "+a-a",
				pattern: Pattern{}.Add(FieldKey("c")),
				output:  [][]string{{"a-a", "a-a"}},
			},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d", idx),
				test.pattern.FindAllStringSubmatch(test.input, -1),
				c.ShouldEqual,
				test.output)
		}

	})

	c.Convey("Keyword", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  [][]string
		}{
			{
				input:   "a",
				pattern: Pattern{}.Add(Keyword("c")),
				output:  [][]string{{"a", "a"}},
			},

			{
				input:   "aa",
				pattern: Pattern{}.Add(Keyword("c")),
				output:  [][]string{{"aa", "aa"}},
			},

			{
				input:   "a-a",
				pattern: Pattern{}.Add(Keyword("c")),
				output:  [][]string{{"a-a", "a-a"}},
			},

			{
				input:   "a-a_a",
				pattern: Pattern{}.Add(Keyword("c")),
				output:  [][]string{{"a-a_a", "a-a_a"}},
			},

			{
				input:   "0-a_a",
				pattern: Pattern{}.Add(Keyword("c")),
				output:  [][]string{{"-a_a", "-a_a"}},
			},

			{
				input:   "-a-a",
				pattern: Pattern{}.Add(Keyword("c")),
				output:  [][]string{{"-a-a", "-a-a"}},
			},

			{
				input:   "+a-a",
				pattern: Pattern{}.Add(Keyword("c")),
				output:  [][]string{{"+a-a", "+a-a"}},
			},

			{
				input:   "+ +a-a",
				pattern: Pattern{}.Add(Keyword("c")),
				output:  [][]string{{"+a-a", "+a-a"}},
			},

			{
				input:   "+ +a-a -",
				pattern: Pattern{}.Add(Keyword("c")),
				output:  [][]string{{"+a-a", "+a-a"}},
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
