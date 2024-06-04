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

func TestRuneIsALNUM(t *testing.T) {
	c.Convey("RuneIsALNUM", t, func() {

		for r := 'a'; r <= 'z'; r += 1 {
			c.SoMsg(
				fmt.Sprintf("test %q", string(r)),
				RuneIsALNUM(r),
				c.ShouldBeTrue,
			)
		}

		for r := 'A'; r <= 'Z'; r += 1 {
			c.SoMsg(
				fmt.Sprintf("test %q", string(r)),
				RuneIsALNUM(r),
				c.ShouldBeTrue,
			)
		}

		for r := '0'; r <= '9'; r += 1 {
			c.SoMsg(
				fmt.Sprintf("test %q", string(r)),
				RuneIsALNUM(r),
				c.ShouldBeTrue,
			)
		}

		for idx, test := range []struct {
			input  rune
			output bool
		}{
			{'!', false},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d", idx),
				RuneIsALNUM(test.input),
				c.ShouldEqual,
				test.output,
			)
		}

	})
}

func TestRuneIsALPHA(t *testing.T) {

	c.Convey("RuneIsALPHA", t, func() {

		for r := 'a'; r <= 'z'; r += 1 {
			c.SoMsg(
				fmt.Sprintf("test %q", string(r)),
				RuneIsALPHA(r),
				c.ShouldBeTrue,
			)
		}

		for r := 'A'; r <= 'Z'; r += 1 {
			c.SoMsg(
				fmt.Sprintf("test %q", string(r)),
				RuneIsALPHA(r),
				c.ShouldBeTrue,
			)
		}

		for r := '0'; r <= '9'; r += 1 {
			c.SoMsg(
				fmt.Sprintf("test %q", string(r)),
				RuneIsALPHA(r),
				c.ShouldBeFalse,
			)
		}

		for idx, test := range []struct {
			input  rune
			output bool
		}{
			{'!', false},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d", idx),
				RuneIsALPHA(test.input),
				c.ShouldEqual,
				test.output,
			)
		}

	})
}

func TestRuneIsBLANK(t *testing.T) {
	c.Convey("RuneIsBLANK", t, func() {
		c.So(RuneIsBLANK('!'), c.ShouldBeFalse)
		c.So(RuneIsBLANK('\t'), c.ShouldBeTrue)
		c.So(RuneIsBLANK(' '), c.ShouldBeTrue)
		c.So(RuneIsBLANK('\v'), c.ShouldBeFalse)
		c.So(RuneIsBLANK('\n'), c.ShouldBeFalse)
	})
}

func TestRuneIsCNTRL(t *testing.T) {
	c.Convey("RuneIsCNTRL", t, func() {
		c.So(RuneIsCNTRL('!'), c.ShouldBeFalse)
		c.So(RuneIsCNTRL(rune(127)), c.ShouldBeTrue)
		for r := rune(0); r <= 31; r++ {
			c.SoMsg(
				fmt.Sprintf("test %q", string(r)),
				RuneIsCNTRL(r),
				c.ShouldBeTrue,
			)
		}
	})
}

func TestRuneIsDIGIT(t *testing.T) {
	c.Convey("RuneIsDIGIT", t, func() {
		c.So(RuneIsDIGIT('!'), c.ShouldBeFalse)
		for r := '0'; r <= '9'; r++ {
			c.SoMsg(
				fmt.Sprintf("test %q", string(r)),
				RuneIsDIGIT(r),
				c.ShouldBeTrue,
			)
		}
	})
}

func TestRuneIsGRAPH(t *testing.T) {
	c.Convey("RuneIsGRAPH", t, func() {
		c.So(RuneIsGRAPH('\n'), c.ShouldBeFalse)
		for r := 'a'; r <= 'z'; r++ {
			c.SoMsg(
				fmt.Sprintf("test %q", string(r)),
				RuneIsGRAPH(r),
				c.ShouldBeTrue,
			)
		}
		for r := 'A'; r <= 'Z'; r++ {
			c.SoMsg(
				fmt.Sprintf("test %q", string(r)),
				RuneIsGRAPH(r),
				c.ShouldBeTrue,
			)
		}
		for r := '0'; r <= '9'; r++ {
			c.SoMsg(
				fmt.Sprintf("test %q", string(r)),
				RuneIsGRAPH(r),
				c.ShouldBeTrue,
			)
		}

		for r := range gLookupGraphical {
			c.SoMsg(
				fmt.Sprintf("test %q", string(r)),
				RuneIsGRAPH(r),
				c.ShouldBeTrue,
			)
		}
	})
}

func TestRuneIsLOWER(t *testing.T) {
	c.Convey("RuneIsLOWER", t, func() {
		c.So(RuneIsLOWER('!'), c.ShouldBeFalse)
		for r := 'a'; r <= 'z'; r++ {
			c.SoMsg(
				fmt.Sprintf("test %q", string(r)),
				RuneIsLOWER(r),
				c.ShouldBeTrue,
			)
		}
		for r := 'A'; r <= 'Z'; r++ {
			c.SoMsg(
				fmt.Sprintf("test %q", string(r)),
				RuneIsLOWER(r),
				c.ShouldBeFalse,
			)
		}
	})
}

func TestRuneIsPRINT(t *testing.T) {
	c.Convey("RuneIsPRINT", t, func() {
		c.So(RuneIsPRINT('\n'), c.ShouldBeFalse)
		c.So(RuneIsPRINT(' '), c.ShouldBeTrue)
		for r := 'a'; r <= 'z'; r++ {
			c.SoMsg(
				fmt.Sprintf("test %q", string(r)),
				RuneIsPRINT(r),
				c.ShouldBeTrue,
			)
		}
		for r := 'A'; r <= 'Z'; r++ {
			c.SoMsg(
				fmt.Sprintf("test %q", string(r)),
				RuneIsPRINT(r),
				c.ShouldBeTrue,
			)
		}
		for r := '0'; r <= '9'; r++ {
			c.SoMsg(
				fmt.Sprintf("test %q", string(r)),
				RuneIsPRINT(r),
				c.ShouldBeTrue,
			)
		}

		for r := range gLookupGraphical {
			c.SoMsg(
				fmt.Sprintf("test %q", string(r)),
				RuneIsPRINT(r),
				c.ShouldBeTrue,
			)
		}
	})
}

func TestRuneIsPUNCT(t *testing.T) {
	c.Convey("RuneIsPUNCT", t, func() {
		c.So(RuneIsPUNCT('\n'), c.ShouldBeFalse)
		c.So(RuneIsPUNCT(' '), c.ShouldBeFalse)
		for r := '!'; r <= '/'; r++ {
			c.SoMsg(
				fmt.Sprintf("test %q", string(r)),
				RuneIsPUNCT(r),
				c.ShouldBeTrue,
			)
		}
		for r := ':'; r <= '@'; r++ {
			c.SoMsg(
				fmt.Sprintf("test %q", string(r)),
				RuneIsPUNCT(r),
				c.ShouldBeTrue,
			)
		}
		for r := '['; r <= '`'; r++ {
			c.SoMsg(
				fmt.Sprintf("test %q", string(r)),
				RuneIsPUNCT(r),
				c.ShouldBeTrue,
			)
		}
		for r := '{'; r <= '~'; r++ {
			c.SoMsg(
				fmt.Sprintf("test %q", string(r)),
				RuneIsPUNCT(r),
				c.ShouldBeTrue,
			)
		}

	})
}

func TestRuneIsSPACE(t *testing.T) {
	c.Convey("RuneIsSPACE", t, func() {
		c.So(RuneIsSPACE('!'), c.ShouldBeFalse)
		for _, r := range []rune{' ', '\t', '\n', '\r', '\v', '\f'} {
			c.SoMsg(
				fmt.Sprintf("test %q", string(r)),
				RuneIsSPACE(r),
				c.ShouldBeTrue,
			)
		}
	})
}

func TestRuneIsUPPER(t *testing.T) {
	c.Convey("RuneIsUPPER", t, func() {
		c.So(RuneIsUPPER('!'), c.ShouldBeFalse)
		for r := 'a'; r <= 'z'; r++ {
			c.SoMsg(
				fmt.Sprintf("test %q", string(r)),
				RuneIsUPPER(r),
				c.ShouldBeFalse,
			)
		}
		for r := 'A'; r <= 'Z'; r++ {
			c.SoMsg(
				fmt.Sprintf("test %q", string(r)),
				RuneIsUPPER(r),
				c.ShouldBeTrue,
			)
		}
	})
}

func TestRuneIsXDIGIT(t *testing.T) {
	c.Convey("RuneIsXDIGIT", t, func() {
		c.So(RuneIsXDIGIT('!'), c.ShouldBeFalse)
		for r := 'a'; r <= 'f'; r++ {
			c.SoMsg(
				fmt.Sprintf("test %q", string(r)),
				RuneIsXDIGIT(r),
				c.ShouldBeTrue,
			)
		}
		for r := 'g'; r <= 'z'; r++ {
			c.SoMsg(
				fmt.Sprintf("test %q", string(r)),
				RuneIsXDIGIT(r),
				c.ShouldBeFalse,
			)
		}
		for r := 'A'; r <= 'F'; r++ {
			c.SoMsg(
				fmt.Sprintf("test %q", string(r)),
				RuneIsXDIGIT(r),
				c.ShouldBeTrue,
			)
		}
		for r := 'G'; r <= 'Z'; r++ {
			c.SoMsg(
				fmt.Sprintf("test %q", string(r)),
				RuneIsXDIGIT(r),
				c.ShouldBeFalse,
			)
		}
		for r := '0'; r <= '9'; r++ {
			c.SoMsg(
				fmt.Sprintf("test %q", string(r)),
				RuneIsXDIGIT(r),
				c.ShouldBeTrue,
			)
		}
	})
}
