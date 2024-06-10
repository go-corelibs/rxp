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

func TestFlags(t *testing.T) {
	c.Convey("Set/Unset/Has", t, func() {
		c.So(DefaultFlags.Set(DotNewlineFlag).String(), c.ShouldEqual, "s")
		c.So((MultilineFlag | DotNewlineFlag).Unset(DotNewlineFlag).String(), c.ShouldEqual, "m")
		c.So((MultilineFlag | DotNewlineFlag).Has(DotNewlineFlag), c.ShouldBeTrue)
		c.So((MultilineFlag | DotNewlineFlag).Has(CaptureFlag), c.ShouldBeFalse)
	})
	c.Convey("ParseFlags", t, func() {

		for idx, test := range []struct {
			input  []string
			reps   Reps
			output string
			panic  c.Assertion
		}{
			{[]string{}, Reps(nil), ``, c.ShouldNotPanic},
			{[]string{"^"}, Reps(nil), `^`, c.ShouldNotPanic},
			{[]string{"m"}, Reps(nil), `m`, c.ShouldNotPanic},
			{[]string{"s"}, Reps(nil), `s`, c.ShouldNotPanic},
			{[]string{"i"}, Reps(nil), `i`, c.ShouldNotPanic},
			{[]string{"c"}, Reps(nil), `c`, c.ShouldNotPanic},
			{[]string{"*"}, Reps{-1, -1}, `*`, c.ShouldNotPanic},
			{[]string{"+"}, Reps{1, -1}, `+`, c.ShouldNotPanic},
			{[]string{"?"}, Reps{0, 1}, `?`, c.ShouldNotPanic},
			{[]string{"*?"}, Reps{-1, -1}, `*?`, c.ShouldNotPanic},
			{[]string{"+?"}, Reps{1, -1}, `+?`, c.ShouldNotPanic},
			{[]string{"??"}, Reps{0, 1}, `??`, c.ShouldNotPanic},
			{[]string{"^msic*"}, Reps{-1, -1}, `*^msic`, c.ShouldNotPanic},
			{[]string{"^msic+"}, Reps{1, -1}, `+^msic`, c.ShouldNotPanic},
			{[]string{"^msic?"}, Reps{0, 1}, `?^msic`, c.ShouldNotPanic},
			{[]string{"^msic*?"}, Reps{-1, -1}, `*?^msic`, c.ShouldNotPanic},
			{[]string{"^msic+?"}, Reps{1, -1}, `+?^msic`, c.ShouldNotPanic},
			{[]string{"^msic??"}, Reps{0, 1}, `??^msic`, c.ShouldNotPanic},
			{[]string{"{10}"}, Reps{10, 10}, ``, c.ShouldNotPanic},
			{[]string{"{1,1}"}, Reps{1, 1}, ``, c.ShouldNotPanic},
			{[]string{"{1,}"}, Reps{1, -1}, ``, c.ShouldNotPanic},
			{[]string{"{1,}?"}, Reps{1, -1}, `?`, c.ShouldNotPanic},
			{[]string{"{1,-1}"}, Reps(nil), ``, c.ShouldPanic},
			{[]string{"{1,0}"}, Reps(nil), ``, c.ShouldPanic},
			{[]string{"NOPE"}, Reps(nil), ``, c.ShouldPanic},
			{[]string{"F"}, Reps(nil), ``, c.ShouldPanic},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d (panic)", idx),
				func() {
					reps, flag := ParseFlags(test.input...)
					c.SoMsg(
						fmt.Sprintf("test #%d (reps)", idx),
						reps, c.ShouldEqual, test.reps,
					)
					c.SoMsg(
						fmt.Sprintf("test #%d (flag)", idx),
						flag.String(), c.ShouldEqual, test.output,
					)
				},
				test.panic,
			)
		}

	})

	c.Convey("Mutable", t, func() {

		c.Convey("ParseOptions", func() {
			c.So(func() { ParseOptions(struct{ string }{"derp"}) }, c.ShouldPanic)
		})

		_, f := ParseFlags(`i`)
		c.So(f.String(), c.ShouldEqual, `i`)
		f |= CaptureFlag
		c.So(f.String(), c.ShouldEqual, `ic`)

		clone := f
		//c.So(clone.String(), c.ShouldEqual, `^ic`)

		_, other := ParseFlags(`^ms`)
		c.So(other.String(), c.ShouldEqual, `^ms`)
		clone = clone.Merge(other)
		c.So(clone.String(), c.ShouldEqual, `^msic`)
		nop := clone.Merge(0)
		c.So(nop.String(), c.ShouldEqual, `^msic`)

	})

	c.Convey("Reps", t, func() {

		for idx, test := range []struct {
			input  Reps
			valid  bool
			min    int
			max    int
			count  int
			minHit bool
			maxHit bool
		}{
			{Reps{}, false, 1, 1, 0, false, false},
			{Reps{2, 1}, false, 2, 1, 0, false, false},
			{Reps{1, 1}, true, 1, 1, 1, true, true},
			{Reps{1, 2}, true, 1, 2, 2, true, true},
			{Reps{1, 1}, true, 1, 1, 0, false, false},
			{Reps{1, -1}, true, 1, -1, 2, true, false},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d (valid)", idx),
				test.input.Valid(), c.ShouldEqual, test.valid,
			)

			c.SoMsg(
				fmt.Sprintf("test #%d (min)", idx),
				test.input.Min(), c.ShouldEqual, test.min,
			)

			c.SoMsg(
				fmt.Sprintf("test #%d (max)", idx),
				test.input.Max(), c.ShouldEqual, test.max,
			)

			if test.count >= 0 {
				minHit, maxHit := test.input.Satisfied(test.count)
				c.SoMsg(
					fmt.Sprintf("test #%d (hit min)", idx),
					minHit, c.ShouldEqual, test.minHit,
				)
				c.SoMsg(
					fmt.Sprintf("test #%d (hit max)", idx),
					maxHit, c.ShouldEqual, test.maxHit,
				)
			}
		}

	})
}
