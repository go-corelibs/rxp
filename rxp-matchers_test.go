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

func TestMatchers(t *testing.T) {
	c.Convey("Custom", t, func() {

		c.So(
			Pattern{}.
				Add(func(scope Flags, reps Reps, input *RuneBuffer, index int, sm [][2]int) (scoped Flags, consumed int, proceed bool) {
					scoped = scope | CaptureFlag
					if prev, _, ok := input.Prev(index); ok {
						if prev == 'o' {
							if this, size, ok := input.Get(index); ok {
								if this == 'n' {
									if next, _, ok := input.Get(index + size); ok {
										if next == 'e' {
											consumed += size
											proceed = true
											return
										}
									}
								}
							}
						}
					}
					return
				}).
				FindAllStringSubmatch("one", -1),
			c.ShouldEqual,
			[][]string{{"n", "n"}},
		)

		p := Pattern{}.
			Caret().
			Text("/w/").
			Text("/", "^", "+", "c").
			Text("/", "??").
			Dollar()
		for idx, test := range []struct {
			input  string
			output [][]string
		}{
			{"/w/core", [][]string{{"/w/core", "core"}}},
			{"/w/core/", [][]string{{"/w/core/", "core"}}},
			{"/w//", [][]string(nil)},
			{"/nope/core/", [][]string(nil)},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				p.FindAllStringSubmatch(test.input, -1),
				c.ShouldEqual, test.output,
			)
		}

		// ^/b/([^/]+)/??$
		p = Pattern{}.
			Caret().
			Text("/b/").
			Not(Text("/"), "+", "c").
			Text("/", "??").
			Dollar()
		for idx, test := range []struct {
			input  string
			output [][]string
		}{
			{"/b/core", [][]string{{"/b/core", "core"}}},
			{"/b/core/", [][]string{{"/b/core/", "core"}}},
			{"/b//", [][]string(nil)},
			{"/nope/core/", [][]string(nil)},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				p.FindAllStringSubmatch(test.input, -1),
				c.ShouldEqual, test.output,
			)
		}

		// ^/build/([a-zA-Z0-9])?/??
		p = Pattern{}.
			Caret().
			Text("/build/").
			Alnum("?", "c").
			Text("/", "??").
			Dollar()
		for _, test := range []struct {
			input  string
			output [][]string
		}{
			{"/build/c", [][]string{{"/build/c", "c"}}},
			{"/build/c/", [][]string{{"/build/c/", "c"}}},
			{"/build//", [][]string{{"/build//", ""}}},
			{"/nope/core/", [][]string(nil)},
		} {
			c.So(p.FindAllStringSubmatch(test.input, -1), c.ShouldEqual, test.output)
		}

	})
}
