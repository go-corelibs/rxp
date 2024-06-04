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
)

type Fragments []Fragment

func (f Fragments) Indexes() (indexes [][]int) {
	for _, frag := range f {
		indexes = append(indexes, frag.Index())
	}
	return
}

func (f Fragments) String() string {
	var buf strings.Builder
	for _, frag := range f {
		buf.WriteString(frag.String())
	}
	return buf.String()
}

func (f Fragments) Strings() (found []string) {
	for _, frag := range f {
		found = append(found, frag.String())
	}
	return
}

func (f Fragments) Recycle() {
	for _, frag := range f {
		frag.Recycle()
	}
}
