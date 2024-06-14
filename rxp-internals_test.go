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

func TestInternals(t *testing.T) {
	c.Convey("Sync Pools", t, func() {

		buf := spBytesBuffer.Get()
		c.So(buf, c.ShouldNotBeNil)
		c.So(buf.String(), c.ShouldEqual, "")
		c.So(_spBytesBufferSetter(buf), c.ShouldNotBeNil)
		buf.WriteString(strings.Repeat(".", 64001))
		c.So(_spBytesBufferSetter(buf), c.ShouldBeNil)
		buf.Reset()
		c.So(_spBytesBufferSetter(buf), c.ShouldNotBeNil)

	})
}
