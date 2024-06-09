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

func TestVersusRegexp(t *testing.T) {
	c.Convey(`a*`, t, func() {

		for idx, test := range []struct {
			input   string
			pattern Pattern
			output  []string
		}{
			{
				input:   `bb`,
				pattern: Pattern{}.Text("a", "*"),
				output: []string{
					"", "", "",
				},
			},
			{
				input:   `abaabaccadaaae`,
				pattern: Pattern{}.Text("a", "*"),
				output: []string{
					"a", "aa", "a", "", "a", "aaa", "",
				},
			},
		} {
			c.SoMsg(
				fmt.Sprintf("test #%d - %q", idx, test.input),
				Pattern{}.
					Text("a", "*").
					FindAllString(test.input, -1),
				c.ShouldEqual,
				test.output,
			)
		}

	})
}
