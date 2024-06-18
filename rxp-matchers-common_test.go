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

	c.Convey("IsFieldWord", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  []string
		}{
			{
				input:   "a' ",
				pattern: Pattern{}.Add(IsFieldWord("c", "^")),
				output:  []string{"'", " "},
			},

			{
				input:   "a' ",
				pattern: Pattern{}.Add(IsFieldWord("c")),
				output:  []string{"a"},
			},

			{
				input:   "a",
				pattern: Pattern{}.Add(IsFieldWord("c")),
				output:  []string{"a"},
			},

			{
				input:   "aa",
				pattern: Pattern{}.Add(IsFieldWord("c")),
				output:  []string{"aa"},
			},

			{
				input:   "one two won't do",
				pattern: Pattern{}.Add(IsFieldWord("c")),
				output:  []string{"one", "two", "won't", "do"},
			},

			{
				input:   "-one two won't do",
				pattern: Pattern{}.Add(IsFieldWord("c")),
				output:  []string{"one", "two", "won't", "do"},
			},

			{
				input:   "-one-two  won't do",
				pattern: Pattern{}.Add(IsFieldWord("c")),
				output:  []string{"one-two", "won't", "do"},
			},

			{
				input:   "-one-two'  won't do",
				pattern: Pattern{}.Add(IsFieldWord("c")),
				output:  []string{"one-two", "won't", "do"},
			},

			{
				input:   "aa-_-aa-",
				pattern: Pattern{}.Add(IsFieldWord("c")),
				output:  []string{"aa-_-aa"},
			},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.FindAllString(test.input, -1),
				c.ShouldEqual,
				test.output)
		}

	})

	c.Convey("IsFieldKey", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  [][]string
		}{
			{
				input:   "a",
				pattern: Pattern{}.Add(IsFieldKey("c")),
				output:  [][]string{{"a", "a"}},
			},

			{
				input:   "aa",
				pattern: Pattern{}.Add(IsFieldKey("c")),
				output:  [][]string{{"aa", "aa"}},
			},

			{
				input:   "a-a",
				pattern: Pattern{}.Add(IsFieldKey("c")),
				output:  [][]string{{"a-a", "a-a"}},
			},

			{
				input:   "a-a_a",
				pattern: Pattern{}.Add(IsFieldKey("c")),
				output:  [][]string{{"a-a_a", "a-a_a"}},
			},

			{
				input:   "0-a_a",
				pattern: Pattern{}.Add(IsFieldKey("c")),
				output:  [][]string{{"a_a", "a_a"}},
			},

			{
				input:   "-a-a",
				pattern: Pattern{}.Add(IsFieldKey("c")),
				output:  [][]string{{"a-a", "a-a"}},
			},

			{
				input:   "a-a-",
				pattern: Pattern{}.Add(IsFieldKey("c")),
				output:  [][]string{{"a-a", "a-a"}},
			},

			{
				input:   "a-a-",
				pattern: Pattern{}.Add(IsFieldKey("c", "^")),
				output:  [][]string(nil), // one may think that this should
				// match the trailing dash character because it is not a field
				// key, however, perl has the same results, ie:
				// $ perl -e '"a-a-" =~ m@((?!\b[a-zA-Z][-_a-zA-Z0-9]+?[a-zA-Z0-9]\b))@; print "matched: [${1}]\n";'
				// matched: []
			},

			{
				input:   "+a-a",
				pattern: Pattern{}.Add(IsFieldKey("c")),
				output:  [][]string{{"a-a", "a-a"}},
			},

			{ // IsFieldKey does not support single-quotes
				input:   "+a'a",
				pattern: Pattern{}.Add(IsFieldKey("c")),
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

	c.Convey("IsKeyword", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  [][]string
		}{
			{
				input:   "a",
				pattern: Pattern{}.Add(IsKeyword("c")),
				output:  [][]string{{"a", "a"}},
			},

			{
				input:   "a",
				pattern: Pattern{}.Add(IsKeyword("c", "^")),
				output:  [][]string(nil),
			},

			{
				input:   "aa",
				pattern: Pattern{}.Add(IsKeyword("c")),
				output:  [][]string{{"aa", "aa"}},
			},

			{
				input:   "a-a",
				pattern: Pattern{}.Add(IsKeyword("c")),
				output:  [][]string{{"a-a", "a-a"}},
			},

			{
				input:   "a-a_a",
				pattern: Pattern{}.Add(IsKeyword("c")),
				output:  [][]string{{"a-a_a", "a-a_a"}},
			},

			{
				input:   "0-a_a",
				pattern: Pattern{}.Add(IsKeyword("c")),
				output:  [][]string{{"-a_a", "-a_a"}},
			},

			{
				input:   "-a-a",
				pattern: Pattern{}.Add(IsKeyword("c")),
				output:  [][]string{{"-a-a", "-a-a"}},
			},

			{
				input:   "a-a-",
				pattern: Pattern{}.Add(IsKeyword("c")),
				output:  [][]string{{"a-a", "a-a"}},
			},

			{
				input:   "+a-a",
				pattern: Pattern{}.Add(IsKeyword("c")),
				output:  [][]string{{"+a-a", "+a-a"}},
			},

			{
				input:   "+ +a-a",
				pattern: Pattern{}.Add(IsKeyword("c")),
				output:  [][]string{{"+a-a", "+a-a"}},
			},

			{
				input:   "+ +a-a -",
				pattern: Pattern{}.Add(IsKeyword("c")),
				output:  [][]string{{"+a-a", "+a-a"}},
			},

			{
				input:   "aa-_-aa-",
				pattern: Pattern{}.Add(IsKeyword("c")),
				output:  [][]string{{"aa-_-aa", "aa-_-aa"}},
			},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				test.pattern.FindAllStringSubmatch(test.input, -1),
				c.ShouldEqual,
				test.output)
		}

	})

	c.Convey("IsHash10", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  [][]string
		}{

			{
				input:   "000a000",
				pattern: Pattern{}.Add(IsHash10()),
				output:  [][]string(nil),
			},

			{
				input:   "0000000000",
				pattern: Pattern{}.Add(IsHash10("^")),
				output:  [][]string(nil),
			},

			{
				input:   "0000000000",
				pattern: Pattern{}.Add(IsHash10()),
				output:  [][]string{{"0000000000"}},
			},

			{
				input:   "0000000000",
				pattern: Pattern{}.Add(IsHash10("^")),
				output:  [][]string(nil),
			},

			{
				input:   "000000000z",
				pattern: Pattern{}.Caret().Add(IsHash10()).Dollar(),
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

	c.Convey("IsAtLeastSixDigits", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  [][]string
		}{

			{
				input:   "000a000",
				pattern: Pattern{}.Add(IsAtLeastSixDigits()),
				output:  [][]string(nil),
			},

			{
				input:   "00000",
				pattern: Pattern{}.Add(IsAtLeastSixDigits()),
				output:  [][]string(nil),
			},

			{
				input:   "0000000",
				pattern: Pattern{}.Add(IsAtLeastSixDigits()),
				output:  [][]string{{"000000"}},
			},

			{
				input:   "0000000",
				pattern: Pattern{}.Add(IsAtLeastSixDigits("^")),
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

	c.Convey("IsUUID", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  [][]string
		}{

			{
				input:   "000a000",
				pattern: Pattern{}.Add(IsUUID()),
				output:  [][]string(nil),
			},

			{
				input:   "12345678-1234-1234-1234-123456789012",
				pattern: Pattern{}.Add(IsUUID("c")),
				output:  [][]string{{"12345678-1234-1234-1234-123456789012", "12345678-1234-1234-1234-123456789012"}},
			},

			{
				input:   "/blah-blah-12345678-1234-1234-1234-123456789012",
				pattern: Pattern{}.Add(IsUUID("c")).Dollar(),
				output:  [][]string{{"12345678-1234-1234-1234-123456789012", "12345678-1234-1234-1234-123456789012"}},
			},

			{
				input:   "/blah-blah-12345678-1234-1234-1234-123456789012-nope",
				pattern: Pattern{}.Add(IsUUID("c")).Dollar(),
				output:  [][]string(nil),
			},

			{
				input:   "12345678!1234-1234-1234-123456789012",
				pattern: Pattern{}.Add(IsUUID("c")),
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
}
