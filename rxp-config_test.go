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

func TestConfig(t *testing.T) {

	c.Convey("Panic", t, func() {

		c.Convey("ParseFlags", func() {

			c.So(func() { _ = ParseFlags("!") }, c.ShouldPanic)
			c.So(func() { _ = ParseFlags("#!") }, c.ShouldPanic)
			c.So(func() { _ = ParseFlags("{1,0}") }, c.ShouldPanic)
			c.So(func() { _ = ParseFlags("{0,1}z?") }, c.ShouldPanic)
			c.So(func() { _ = ParseFlags("{0,1}?z") }, c.ShouldPanic)
			c.So(func() { _ = ParseFlags("{nope}?") }, c.ShouldPanic)

		})

		c.Convey("ParseOptions", func() {
			c.So(func() { ParseOptions(struct{ string }{"derp"}) }, c.ShouldPanic)
		})

	})

	c.Convey("Corner Cases", t, func() {

		c.So((&Config{Minimum: -1, Maximum: 1}).Valid(), c.ShouldBeFalse)
		c.So((&Config{Minimum: 3, Maximum: 4}).IsValidCount(1), c.ShouldBeFalse)

	})

	c.Convey("One Rune Argument", t, func() {

		for idx, test := range []struct {
			input  []string
			output *Config
		}{

			{input: nil, output: &Config{Minimum: 1, Maximum: 1}},
			{input: []string{"m"}, output: &Config{Minimum: 1, Maximum: 1, Multiline: true}},
			{input: []string{"i"}, output: &Config{Minimum: 1, Maximum: 1, AnyCase: true}},
			{input: []string{"s"}, output: &Config{Minimum: 1, Maximum: 1, DotNL: true}},
			{input: []string{"c"}, output: &Config{Minimum: 1, Maximum: 1, Capture: true}},
			{input: []string{"*"}, output: &Config{Minimum: -1, Maximum: -1, Less: false}},
			{input: []string{"+"}, output: &Config{Minimum: 1, Maximum: -1}},
			{input: []string{"?"}, output: &Config{Minimum: 0, Maximum: 1}},
			{input: []string{"^"}, output: &Config{Minimum: 1, Maximum: 1, Negated: true}},
			{input: []string{""}, output: &Config{Minimum: 1, Maximum: 1}},
		} {

			c.SoMsg(
				fmt.Sprintf("test #%d", idx),
				ParseFlags(test.input...),
				c.ShouldEqual,
				test.output,
			)

		}

	})

	c.Convey("Mixed Arguments", t, func() {

		for idx, test := range []struct {
			input  []string
			output *Config
		}{

			{input: nil, output: &Config{Minimum: 1, Maximum: 1}},
			{input: []string{"mi"}, output: &Config{Minimum: 1, Maximum: 1, Multiline: true, AnyCase: true}},
			{input: []string{"is"}, output: &Config{Minimum: 1, Maximum: 1, AnyCase: true, DotNL: true}},
			{input: []string{"sc"}, output: &Config{Minimum: 1, Maximum: 1, DotNL: true, Capture: true}},
			{input: []string{"c?*"}, output: &Config{Minimum: -1, Maximum: -1, Capture: true, Less: false}},
			{input: []string{"*+?"}, output: &Config{Minimum: 1, Maximum: -1, Less: true}},
			{input: []string{"*??"}, output: &Config{Minimum: 0, Maximum: 1, Less: true}},
			{input: []string{"??i"}, output: &Config{Minimum: 0, Maximum: 1, Less: true, AnyCase: true}},
			{input: []string{"{1}"}, output: &Config{Minimum: 1, Maximum: 1, Less: false}},
			{input: []string{"{1,}"}, output: &Config{Minimum: 1, Maximum: -1, Less: false}},
			{input: []string{"{1,10}"}, output: &Config{Minimum: 1, Maximum: 10, Less: false}},
			{input: []string{"{1,10}?"}, output: &Config{Minimum: 1, Maximum: 10, Less: true}},
			{input: []string{"?^"}, output: &Config{Minimum: 0, Maximum: 1, Negated: true}},
			{input: []string{"  ", ""}, output: &Config{Minimum: 1, Maximum: 1}},
		} {

			c.SoMsg(
				fmt.Sprintf("test #%d", idx),
				ParseFlags(test.input...),
				c.ShouldEqual,
				test.output,
			)

		}

	})

}
