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
	sync "github.com/go-corelibs/x-sync"
)

type Segments []Segment

func (m Segments) Indexes() (indexes [][]int) {
	for _, match := range m {
		indexes = append(indexes, match.Index())
	}
	return
}

func (m Segments) String() string {
	buf := spStringBuilder.Get()
	defer spStringBuilder.Put(buf)
	for _, match := range m {
		buf.WriteString(match.String())
	}
	return buf.String()
}

func (m Segments) Strings() (found []string) {
	for _, match := range m {
		found = sync.Append(found, match.String())
	}
	return
}
