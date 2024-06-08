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

	"github.com/go-corelibs/values"
)

var _ Segment = (*cSegment)(nil)

type Segment interface {
	// Match returns true if this is a matched segment
	Match() bool

	// Runes returns the complete rune slice for this Pattern Segment
	Runes() []rune

	// String returns the complete string for this Pattern match
	String() (complete string)

	// Index returns the start and end input rune indexes
	//
	// Note: rune slice indexes, not byte or string indexes
	Index() []int

	// Submatch returns the specific submatch value, if it exists
	Submatch(idx int) (value string, ok bool)

	private(v *cSegment)
}

type cSegment struct {
	input    *RuneBuffer
	matched  bool
	matches  SubMatches
	complete *string
}

func (r *cSegment) private(_ *cSegment) {}

func (r *cSegment) Match() bool {
	return r.matched
}

func (r *cSegment) Runes() (runes []rune) {

	for _, m := range r.matches {
		slice, _ := r.input.Slice(m.Start(), m.End()-m.Start())
		runes = appendSlice(runes, slice...)
	}

	return
}

func (r *cSegment) String() string {
	if r.complete != nil {
		return *r.complete
	}

	var buf strings.Builder
	if len(r.matches) == 1 {
		buf.WriteString(r.input.String(r.matches[0].Start(), r.matches[0].End()-r.matches[0].Start()))
	} else {
		for _, m := range r.matches[1:] {
			// skip first which is the entire matched text
			// regardless of m.capture state
			buf.WriteString(r.input.String(m.Start(), m.End()-m.Start()))
		}
	}

	r.complete = values.Ref(buf.String())
	return *r.complete
}

func (r *cSegment) Index() []int {
	return []int{r.matches.Start(), r.matches.End()}
}

func (r *cSegment) Submatch(idx int) (value string, ok bool) {
	if ok = 0 <= idx+1 && idx+1 < len(r.matches); ok {
		m := r.matches[idx+1]
		value = r.input.String(m.Start(), m.End()-m.Start())
	}
	return
}
