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

func TestFlagsReps(t *testing.T) {
	c.Convey("Reps", t, func() {

		for idx, test := range []struct {
			reps   Reps
			min    int
			max    int
			isNil  bool
			valid  bool
			count  int
			minHit bool
			maxHit bool
		}{
			{
				Reps{},
				DefaultMinReps, DefaultMaxReps,
				true, false,
				1, false, false,
			},

			{
				Reps{0},
				0, DefaultMaxReps,
				true, false,
				1, false, false,
			},

			{
				Reps{0, 1},
				0, 1,
				false, true,
				1, true, true,
			},

			{
				Reps{1, -1},
				1, -1,
				false, true,
				1, true, false,
			},

			{
				Reps{2, -1},
				2, -1,
				false, true,
				1, false, false,
			},

			{
				Reps{2, 1},
				2, 1,
				false, false,
				1, false, true,
			},
		} {
			mn, mx := test.reps.Satisfied(test.count)
			c.SoMsg(fmt.Sprintf("test #%d (Min)", idx), test.reps.Min(), c.ShouldEqual, test.min)
			c.SoMsg(fmt.Sprintf("test #%d (Max)", idx), test.reps.Max(), c.ShouldEqual, test.max)
			c.SoMsg(fmt.Sprintf("test #%d (IsNil)", idx), test.reps.IsNil(), c.ShouldEqual, test.isNil)
			c.SoMsg(fmt.Sprintf("test #%d (Valid)", idx), test.reps.Valid(), c.ShouldEqual, test.valid)
			c.SoMsg(fmt.Sprintf("test #%d (minHit)", idx), mn, c.ShouldEqual, test.minHit)
			c.SoMsg(fmt.Sprintf("test #%d (maxHit)", idx), mx, c.ShouldEqual, test.maxHit)
		}

	})
}
