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

func TestMatchState(t *testing.T) {

	c.Convey("Corner Cases", t, func() {

		m := &cMatchState{
			s:        &cPatternState{input: []rune("1234567890")},
			this:     10,
			capture:  false,
			complete: false,
		}
		c.So(m.Len(), c.ShouldEqual, 10)
		c.So(m.Captured(), c.ShouldBeFalse)
		m.Capture()
		c.So(m.Captured(), c.ShouldBeTrue)
		c.So(m.Ready(), c.ShouldBeFalse)
		c.So(m.Consume(0), c.ShouldBeFalse)
		m.this = 1
		c.So(m.Ready(), c.ShouldBeTrue)
		c.So(m.Consume(0), c.ShouldBeTrue)
		c.So(m.Len(), c.ShouldEqual, 1)
		c.So(m.Consume(1), c.ShouldBeTrue)
		c.So(m.Len(), c.ShouldEqual, 2)
		c.So(m.Consume(20), c.ShouldBeFalse)
		c.So(m.Len(), c.ShouldEqual, 10)
		m.this = 0
		c.So(m.String(), c.ShouldEqual, "")

		m1 := &cMatchState{
			s:        m.s,
			this:     7,
			capture:  true,
			complete: true,
		}
		m1.Apply(m)
		c.So(m.this, c.ShouldEqual, m1.this)
		c.So(m.capture, c.ShouldEqual, m1.capture)
		c.So(m.complete, c.ShouldEqual, m1.complete)
		m2 := m.CloneWith(-1)
		c.So(m2.Equal(m), c.ShouldBeTrue)
	})

}
