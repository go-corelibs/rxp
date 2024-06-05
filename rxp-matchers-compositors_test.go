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

func TestMatchersCompositors(t *testing.T) {

	c.Convey("Or", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  [][]string
		}{

			//{
			//	input: "aBb Aa!",
			//	//pattern: Pattern{}.Or(Pattern{}.W("^").S(), "c"),
			//	//pattern: Pattern{}.Or(Pattern{}.W("^"), "c"),
			//	pattern: Pattern{}.W("^", "c"),
			//	output:  [][]string{{" ", " "}, {"!", "!"}},
			//},

			{
				input: "aBbAa",
				pattern: Pattern{}.Or(
					Text("b", "{1,2}", "i"),
					Text("a", "{1}"),
					"c"),
				output: [][]string{{"a", "a"}, {"Bb", "Bb"}, {"a", "a"}},
			},

			{
				input: "aBb\nAa",
				pattern: Pattern{}.Or(
					Text("b", "{1,2}", "i"),
					Dot("{1}"),
					"c"),
				output: [][]string{{"a", "a"}, {"Bb", "Bb"}, {"A", "A"}, {"a", "a"}},
			},

			{ // testing for true next and false keep within Or
				input: "a . b",
				pattern: Pattern{}.Or("+",
					A(), Z(), B(),
				),
				output: [][]string(nil),
			},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d", idx),
				test.pattern.FindAllStringSubmatch(test.input, -1),
				c.ShouldEqual,
				test.output)
		}

	})

	c.Convey("Not", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  [][]string
		}{

			{
				input: "aBbAa",
				pattern: Pattern{}.
					Text("b", "{1,2}", "i", "c").
					Text("a", "^"), // not lower, no capture
				output: [][]string{{"BbA", "Bb"}},
			},

			{
				input: "aBbaaBB",
				pattern: Pattern{}.
					Text("b", "{1,2}", "i", "c").
					Text("a", "^"),
				output: [][]string{{"BB", "BB"}},
			},

			{
				input: "aBbaaBB",
				pattern: Pattern{}.Not("+",
					Text("b", "{1,2}", "i", "c"),
				),
				output: [][]string{{"a"}, {"aa"}},
			},

			{
				input: "aBbaaBB",
				pattern: Pattern{}.Not("+",
					Text("b"),
					Text("B"),
				),
				output: [][]string{{"a"}, {"aa"}},
			},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d", idx),
				test.pattern.FindAllStringSubmatch(test.input, -1),
				c.ShouldEqual,
				test.output)
		}

	})

	c.Convey("Group", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  [][]string
		}{

			{ // ([^\w\s]+)
				input: "aBb Aa!",
				//pattern: Pattern{}.Group(Pattern{}.Or(Pattern{}.W("^").S("^")), "+", "c"),
				//pattern: Pattern{}.Not(Group(Pattern{}.Or(Pattern{}.W().S()), "+"), "c"),
				//pattern: Pattern{}.Group(Pattern{}.Not(Or(Pattern{}.W().S())), "+", "c"),
				pattern: Pattern{}.Not(Or(W(), S()), "+", "c"),
				output:  [][]string{{"!", "!"}},
			},

			{
				input:   "abaaa",
				pattern: Pattern{}.Group(Text("a"), Text("b"), "+", "c"),
				output:  [][]string{{"ab", "ab"}},
			},

			//{ // (\w+\s)
			//	input:   "aBb Aa!",
			//	pattern: Pattern{}.Group(Pattern{}.W("+", "c")).S(),
			//	output:  [][]string{{"aBb ", "aBb"}, {"Aa", "Aa"}},
			//},

			{ // ([^aB])
				input: "aBbAa",
				//pattern: Pattern{}.Group(Pattern{}.Text("a").Text("B"), "^", "c"),
				pattern: Pattern{}.Not(Or(Text("a"), Text("B")), "c"),
				output:  [][]string{{"b", "b"}, {"A", "A"}},
			},

			// BEGIN oddities
			////
			//// not sure what this combination of Not and Group is intended
			//// to produce, however it is producing strange results
			//// indicating something is not entirely correct
			//{ // this one is not supposed to match the first "a"
			//	input:   "aBbAa",
			//	pattern: Pattern{}.Not(Group(Text("a"), Text("B")), "+", "c"),
			//	output:  [][]string{{"Bb", "Bb"}, {"Aa", "Aa"}},
			//	//actual:  [][]string{{"aBbAa", "aBbAa"}},
			//},
			//{ // this one should in fact return pairs of letters that are not "aB"
			//	input:   "aBbAa",
			//	pattern: Pattern{}.Not(Group(Text("a"), Text("B"))),
			//	output:  [][]string{{"Bb", "Bb"}, {"Aa", "Aa"}},
			//	// actual: [][]string(nil)
			//},
			// END oddities
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d", idx),
				test.pattern.FindAllStringSubmatch(test.input, -1),
				c.ShouldEqual,
				test.output)
		}

	})

	c.Convey("Complicated", t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  [][]string
		}{

			//{
			//	input: "kebab-thing[10]",
			//	pattern: Pattern{}.Caret("m").Group(
			//		WrapMatcher(RuneIsALPHA),
			//		Or(Text("-", "+?"), W("+?")),
			//		WrapMatcher(RuneIsALNUM, "*")).
			//		Text("[").
			//		D("+", "c").
			//		Text("]").
			//		Dollar(),
			//	output: [][]string{{"kebab-thing[10]", "kebab-thing", "10"}},
			//},

			{
				input: "thing[10]",
				pattern: Pattern{}.
					A().
					Add(FieldKey("+", "c")).
					Text("[").
					D("+", "c").
					Text("]").
					Dollar(),
				output: [][]string{{"thing[10]", "thing", "10"}},
			},

			{
				input: "other-thing[10]",
				pattern: Pattern{}.
					A().
					Add(FieldKey("+", "c")).
					Text("[").
					D("+", "c").
					Text("]").
					Dollar(),
				output: [][]string{{"other-thing[10]", "other-thing", "10"}},
			},

			{
				input: "other_-thing[10]",
				pattern: Pattern{}.
					A().
					Add(FieldKey("+", "c")).
					Text("[").
					D("+", "c").
					Text("]").
					Dollar(),
				output: [][]string{{"other_-thing[10]", "other_-thing", "10"}},
			},

			{
				input: "other_-_thing[10]",
				pattern: Pattern{}.
					A().
					Add(FieldKey("+", "c")).
					Text("[").
					D("+", "c").
					Text("]").
					Dollar(),
				output: [][]string{{"other_-_thing[10]", "other_-_thing", "10"}},
			},

			{
				input: "other1thing[10]",
				pattern: Pattern{}.
					A().
					Add(FieldKey("+", "c")).
					Text("[").
					D("+", "c").
					Text("]").
					Dollar(),
				output: [][]string{{"other1thing[10]", "other1thing", "10"}},
			},

			{
				input: "other-1thing[10]",
				pattern: Pattern{}.
					A().
					Add(FieldKey("+", "c")).
					Text("[").
					D("+", "c").
					Text("]").
					Dollar(),
				output: [][]string{{"other-1thing[10]", "other-1thing", "10"}},
			},

			{
				input: "0other_-_thing[10]",
				pattern: Pattern{}.
					A().
					Add(FieldKey("+", "c")).
					Text("[").
					D("+", "c").
					Text("]").
					Dollar(),
				output: [][]string(nil),
			},

			{
				input: "0-other_-_thing[10]",
				pattern: Pattern{}.
					A().
					Add(FieldKey("+", "c")).
					Text("[").
					D("+", "c").
					Text("]").
					Dollar(),
				output: [][]string(nil),
			},

			{
				input: "-other_-_thing[10]",
				pattern: Pattern{}.
					A().
					Add(FieldKey("+", "c")).
					Text("[").
					D("+", "c").
					Text("]").
					Dollar(),
				output: [][]string(nil),
			},

			//{
			//	input: "thing[10]",
			//	pattern: Pattern{}.
			//		A().
			//		Group(Pattern{}.
			//			// this AsciiNames matches
			//			AsciiNames(ALPHA, "{1}").
			//			// this Or consumes everything
			//			Or(Pattern{}.Text("-").W(), "+?"). // this isn't respecting the less
			//			// this does not match because previous got everything
			//			AsciiNames(ALNUM, "{1}"),
			//			"c").
			//		Text("[").
			//		D("+", "c").
			//		Text("]").
			//		Dollar(),
			//	output: [][]string{{"thing[10]", "thing", "10"}},
			//},

			{
				input: " func1  ",
				pattern: Pattern{
					Caret(),
					S("*"),
					Or(
						D("+"),
						Group(
							Text("func"),
							D("+"),
						),
						"c"),
					S("*"),
					Dollar(),
				},
				output: [][]string{{" func1  ", "func1"}},
			},

			{
				input: " fun1  ",
				pattern: Pattern{
					Caret(),
					S("*"),
					Or(
						D("+"),
						Group(
							Text("func"),
							D("+"),
						),
						"c"),
					S("*"),
					Dollar(),
				},
				output: [][]string(nil),
			},

			{
				input: " 1  ",
				pattern: Pattern{
					Caret(),
					S("*"),
					Or(
						D("+"),
						Group(
							Text("func"),
							D("+"),
						),
						"c"),
					S("*"),
					Dollar(),
				},
				output: [][]string{{" 1  ", "1"}},
			},
		} {
			c.SoMsg(fmt.Sprintf("test #%d", idx),
				test.pattern.FindAllStringSubmatch(test.input, -1),
				c.ShouldEqual,
				test.output)
		}
	})

}
