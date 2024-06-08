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
	"testing"

	c "github.com/smartystreets/goconvey/convey"
)

func TestRuneBuffer(t *testing.T) {
	c.Convey("RuneBuffer", t, func() {

		rb := NewRuneBuffer("stuff")
		// len
		c.So(rb.Len(), c.ShouldEqual, 5)
		// get
		r, size, ok := rb.Get(2)
		c.So(ok, c.ShouldBeTrue)
		c.So(r, c.ShouldEqual, 'u')
		c.So(size, c.ShouldEqual, 1)
		r, size, ok = rb.Get(6)
		c.So(ok, c.ShouldBeFalse)
		c.So(r, c.ShouldEqual, rune(0))
		c.So(size, c.ShouldEqual, 0)
		// ready
		c.So(rb.Ready(0), c.ShouldBeTrue)
		c.So(rb.Ready(5), c.ShouldBeFalse)
		// valid
		c.So(rb.Valid(0), c.ShouldBeTrue)
		c.So(rb.Valid(5), c.ShouldBeTrue)
		c.So(rb.Valid(6), c.ShouldBeFalse)
		// invalid
		c.So(rb.Invalid(0), c.ShouldBeFalse)
		c.So(rb.Invalid(5), c.ShouldBeTrue)
		c.So(rb.Invalid(6), c.ShouldBeTrue)
		// end
		c.So(rb.End(0), c.ShouldBeFalse)
		c.So(rb.End(5), c.ShouldBeTrue)
		c.So(rb.End(6), c.ShouldBeFalse)
		// slice
		slice, total := rb.Slice(1, 4)
		c.So(slice, c.ShouldEqual, []rune("tuff"))
		c.So(total, c.ShouldEqual, 4)
		slice, total = rb.Slice(5, 1)
		c.So(slice, c.ShouldEqual, []rune(nil))
		c.So(total, c.ShouldEqual, 0)
		// string
		c.So(rb.String(1, 4), c.ShouldEqual, "tuff")
		c.So(rb.String(5, 6), c.ShouldEqual, "")
		// prev
		r, size, ok = rb.Prev(3)
		c.So(ok, c.ShouldBeTrue)
		c.So(r, c.ShouldEqual, 'u')
		c.So(size, c.ShouldEqual, 1)
		r, size, ok = rb.Prev(0)
		c.So(ok, c.ShouldBeFalse)
		c.So(r, c.ShouldEqual, 0)
		c.So(size, c.ShouldEqual, 0)
		// next
		r, size, ok = rb.Next(1)
		c.So(ok, c.ShouldBeTrue)
		c.So(r, c.ShouldEqual, 'u')
		c.So(size, c.ShouldEqual, 1)
		r, size, ok = rb.Next(5)
		c.So(ok, c.ShouldBeFalse)
		c.So(r, c.ShouldEqual, 0)
		c.So(size, c.ShouldEqual, 0)

	})
}
