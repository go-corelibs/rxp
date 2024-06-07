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

func TestPatternMatches(t *testing.T) {

	c.Convey("SubMatch", t, func() {

		for idx, test := range []struct {
			input SubMatch
			valid bool
			len   int
			start int
			end   int
		}{
			{SubMatch{}, false, 0, 0, 0},
			{SubMatch{0}, false, 0, 0, 0},
			{SubMatch{0, 1}, true, 1, 0, 1},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d (valid)", idx),
				test.input.Valid(), c.ShouldEqual, test.valid,
			)
			c.SoMsg(
				fmt.Sprintf("test #%d (size)", idx),
				test.input.Len(), c.ShouldEqual, test.len,
			)
			c.SoMsg(
				fmt.Sprintf("test #%d (start)", idx),
				test.input.Start(), c.ShouldEqual, test.start,
			)
			c.SoMsg(
				fmt.Sprintf("test #%d (end)", idx),
				test.input.End(), c.ShouldEqual, test.end,
			)
		}

	})

	c.Convey("SubMatches", t, func() {

		sm, ok := SubMatches{}.Get(1)
		c.So(ok, c.ShouldBeFalse)
		c.So(sm, c.ShouldBeNil)

		for idx, test := range []struct {
			input SubMatches
			valid bool
			len   int
			start int
			end   int
		}{
			{SubMatches{}, true, 0, 0, 0},
			{SubMatches{{0}}, false, 0, 0, 0},
			{SubMatches{{0, 1}}, true, 1, 0, 1},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d (valid)", idx),
				test.input.Valid(), c.ShouldEqual, test.valid,
			)
			c.SoMsg(
				fmt.Sprintf("test #%d (size)", idx),
				test.input.Size(), c.ShouldEqual, test.len,
			)
			c.SoMsg(
				fmt.Sprintf("test #%d (start)", idx),
				test.input.Start(), c.ShouldEqual, test.start,
			)
			c.SoMsg(
				fmt.Sprintf("test #%d (end)", idx),
				test.input.End(), c.ShouldEqual, test.end,
			)
		}

	})

	c.Convey("Matches", t, func() {

		m := Matches{{{0, 1}}}
		c.So(m.Len(), c.ShouldEqual, 1)
		sm, ok := m.Get(1)
		c.So(ok, c.ShouldBeFalse)
		c.So(sm, c.ShouldBeNil)
		sm, ok = m.Get(0)
		c.So(ok, c.ShouldBeTrue)
		c.So(sm, c.ShouldNotBeNil)

		for idx, test := range []struct {
			input Matches
			valid bool
			len   int
			start int
			end   int
		}{
			{Matches{}, true, 0, 0, 0},
			{Matches{{{0}}}, false, 0, 0, 0},
			{Matches{{{0, 1}}}, true, 1, 0, 1},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d (valid)", idx),
				test.input.Valid(), c.ShouldEqual, test.valid,
			)
			c.SoMsg(
				fmt.Sprintf("test #%d (size)", idx),
				test.input.Size(), c.ShouldEqual, test.len,
			)
			c.SoMsg(
				fmt.Sprintf("test #%d (start)", idx),
				test.input.Start(), c.ShouldEqual, test.start,
			)
			c.SoMsg(
				fmt.Sprintf("test #%d (end)", idx),
				test.input.End(), c.ShouldEqual, test.end,
			)
		}

	})

}
