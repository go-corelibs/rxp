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
					Replace{}.WithReplace(func(s Segment) string {
						var buf strings.Builder
						buf.WriteString(s.String())
						buf.WriteString(s.String())
						if value, ok := s.Submatch(0); ok {
							buf.WriteString(value)
						} else if value, ok = s.Submatch(1); ok {
							buf.WriteString(value)
						}
						return buf.String()
					})).
				Process(`ONE TWO`),
			c.ShouldEqual,
			`ONE TWOONE TWOONE TWO`,
		)

		c.So(
			Pipeline{}.
				Replace(Pattern{}.Dot(), Replace{}.ToLower()).
				Replace(Pattern{}.Text("two"), Replace{func(s Segment) (replaced string) {
					return "2"
				}}).
				Process(`ONE TWO MANY MORE`),
			c.ShouldEqual,
			`one 2 many more`,
		)

		c.So(
			Pipeline{}.
				Replace(Pattern{}.Dot("*"), Replace{}.ToUpper()).
				Replace(Pattern{}.Text("two", "i"), Replace{func(s Segment) (replaced string) {
					return "2"
				}}).
				Process(`one two many more`),
			c.ShouldEqual,
			`ONE 2 MANY MORE`,
		)

		c.So(
			Pipeline{}.
				Transform(strings.ToUpper).
				Replace(Pattern{}.Text("two", "i"), Replace{func(s Segment) (replaced string) {
					return "2"
				}}).
				Process(`one two many more`),
			c.ShouldEqual,
			`ONE 2 MANY MORE`,
		)

		c.So(
			Pipeline{}.
				ReplaceWith(Pattern{}.Text("ONE", "i"), strings.ToUpper).
				Replace(Pattern{}.Text("two", "i"), Replace{func(s Segment) (replaced string) {
					return "2"
				}}).
				Process(`one two many more`),
			c.ShouldEqual,
			`ONE 2 many more`,
		)

		c.So(
			Pipeline{}.
				ReplaceWith(Pattern{}.Text("ONE", "i"), strings.ToUpper).
				Replace(Pattern{}.Text("two", "i"), Replace{}.WithText("2")).
				Process(`one two many more`),
			c.ShouldEqual,
			`ONE 2 many more`,
		)

		c.So(
			Pipeline{}.
				ReplaceWith(Pattern{}.Text("ONE", "i"), strings.ToUpper).
				ReplaceText(Pattern{}.Text("two", "i"), "2").
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
			ReplaceText(S("+"), " ").
			ReplaceText(Text("'"), "_").
			ReplaceText(Not(W(), S(), "+"), "").
			Process(`Isn't  this  neat?`),
			c.ShouldEqual, `isn_t this neat`)

		c.So(Pipeline{
			{Transform: strings.TrimSpace},
			//{ // this works too, though the Or introduces overhead
			//	Search:  Pattern{Or(S("+"), Not(Alnum("+")))},
			//	Replace: Replace{}.WithText("_"),
			//},
			//// these two work best it seems
			//{Search: Pattern{S("+", "c")}, Replace: Replace{}.WithText("_")},
			//{Search: Pattern{Not(Alnum("+"))}, Replace: Replace{}.WithText("_")},
			{Search: Pattern{S("+")}, Replace: Replace{}.WithText("_")},
			{Search: Pattern{Alnum("^")}, Replace: Replace{}.WithText("_")},
			//// caret is not what it seems?
			//{Search: Pattern{Alnum("+", "^")}, Replace: Replace{}.WithText("_")},
			//{ // so fast
			//	Transform: func(input string) string {
			//		var buf strings.Builder
			//		var spaces bool
			//		for _, r := range input {
			//			if r == ' ' {
			//				spaces = true
			//				continue
			//			} else if spaces {
			//				buf.WriteRune('_')
			//				spaces = false
			//			}
			//			if ('a' <= r && r <= 'z') ||
			//				('A' <= r && r <= 'Z') ||
			//				('0' <= r && r <= '9') {
			//				buf.WriteRune(r)
			//				continue
			//			} else {
			//				// not a space, nor alnum
			//				buf.WriteRune('_')
			//			}
			//		}
			//		return buf.String()
			//	},
			//},
			{Transform: strings.ToLower},
		}.Process(`Isn't this neat?`),
			c.ShouldEqual, `isn_t_this_neat_`)

	})
}
