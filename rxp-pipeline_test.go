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
	"strings"
	"testing"

	c "github.com/smartystreets/goconvey/convey"
)

func TestPipeline(t *testing.T) {

	c.Convey("Mixed", t, func() {

		c.So(
			Pipeline{}.
				Replace(nil, nil).
				Replace(
					Pattern{}.Dot("+sc"),
					Replace[string]{}.WithReplace(func(input *InputReader, captured [][2]int, text string) (replaced string) {
						return text + text + text
					})).
				Process(`ONE TWO`),
			c.ShouldEqual,
			`ONE TWOONE TWOONE TWO`,
		)

		c.So(
			Pipeline{}.
				Replace(Pattern{}.Dot("+"), Replace[string]{}.ToLower()).
				Replace(Pattern{}.Text("two"), Replace[string]{func(input *InputReader, captured [][2]int, text string) (replaced string) {
					return "2"
				}}).
				Process(`ONE TWO MANY MORE`),
			c.ShouldEqual,
			`one 2 many more`,
		)

		c.So(
			Pipeline{}.
				Replace(Pattern{}.Dot("*"), Replace[string]{}.ToUpper()).
				Replace(Pattern{}.Text("two", "i"), Replace[string]{func(input *InputReader, captured [][2]int, text string) (replaced string) {
					return "2"
				}}).
				Process(`one two many more`),
			c.ShouldEqual,
			`ONE 2 MANY MORE`,
		)

		c.So(
			Pipeline{}.
				Transform(strings.ToUpper).
				Replace(Pattern{}.Text("two", "i"), Replace[string]{func(input *InputReader, captured [][2]int, text string) (replaced string) {
					return "2"
				}}).
				Process(`one two many more`),
			c.ShouldEqual,
			`ONE 2 MANY MORE`,
		)

		c.So(
			Pipeline{}.
				Substitute(Pattern{}.Text("ONE", "i"), strings.ToUpper).
				Replace(Pattern{}.Text("two", "i"), Replace[string]{func(input *InputReader, captured [][2]int, text string) (replaced string) {
					return "2"
				}}).
				Process(`one two many more`),
			c.ShouldEqual,
			`ONE 2 many more`,
		)

		c.So(
			Pipeline{}.
				Substitute(Pattern{}.Text("ONE", "i"), strings.ToUpper).
				Replace(Pattern{}.Text("two", "i"), Replace[string]{}.WithLiteral("2")).
				Process(`one two many more`),
			c.ShouldEqual,
			`ONE 2 many more`,
		)

		c.So(
			Pipeline{}.
				Substitute(Pattern{}.Text("ONE", "i"), strings.ToUpper).
				Literal(Pattern{}.Text("two", "i"), "2").
				Process(`one two many more`),
			c.ShouldEqual,
			`ONE 2 many more`,
		)

		// from the README.md:
		// output := strings.ToLower(`Isn't  this  neat?`)
		// output = regexp.MustCompile(`\s+`).ReplaceAllString(output, " ")
		// output = regexp.MustCompile(`[']`).ReplaceAllString(output, "_")
		// output = regexp.MustCompile(`[^\w ]`).ReplaceAllString(output, "")

		c.So(Pipeline{}.
			Transform(strings.ToLower).
			Literal(S("+"), " ").
			Literal(Text("'"), "_").
			Literal(Not(W(), S(), "+"), "").
			//Literal(Or(W(), S(), "+", "^"), ""). // this is equivalent
			Process(`Isn't  this  neat?`),
			c.ShouldEqual, `isn_t this neat`)

		c.So(Pipeline{
			{Transform: strings.TrimSpace},
			{Search: Pattern{S("+")}, Replace: Replace[string]{}.WithLiteral("_")},
			{Search: Pattern{Alnum("^")}, Replace: Replace[string]{}.WithLiteral("_")},
			{Transform: strings.ToLower},
		}.Process(`Isn't this neat?`),
			c.ShouldEqual, `isn_t_this_neat_`)

		c.So(Pipeline{
			{Transform: strings.TrimSpace},
			{Search: Pattern{S("+")}, Replace: Replace[string]{}.WithTransform(func(input string) (output string) {
				return "_"
			})},
			{Search: Pattern{Alnum("^")}, Replace: Replace[string]{}.WithLiteral("_")},
			{Transform: strings.ToLower},
		}.Process(`Isn't this neat?`),
			c.ShouldEqual, `isn_t_this_neat_`)

	})
}
