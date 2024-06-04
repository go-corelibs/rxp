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

type Match interface {
	// Runes returns teh complete text for this Pattern match as a rune slice
	Runes() []rune
	// String returns the complete text for this Pattern match
	String() (complete string)
	// Submatch returns the specific value of the capture group at the given
	// index, if it exists
	Submatch(idx int) (value string, ok bool)

	private()
}

type cMatch struct {
	s        *cPatternState
	matches  matchStates
	submatch matchStates
	complete *string
}

// newMatch will panic on nil matches
func newMatch(s *cPatternState, matches matchStates) *cMatch {
	if len(matches) > 0 {
		submatch := make(matchStates, len(matches))
		for _, m := range matches {
			submatch[len(submatch)-1] = m
		}
		return &cMatch{
			s:        s,
			matches:  matches,
			submatch: submatch,
		}
	}
	panic("programmer error - nil matches received")
}

func (r *cMatch) private() {}

func (r *cMatch) Runes() (runes []rune) {
	//runes = make([]rune, len(r.s.input))
	//for _, m := range r.matches {
	//	for idx := m.start; idx < m.start+m.this; idx += 1 {
	//		runes[idx] = r.s.input[idx]
	//	}
	//}

	//return []rune(r.String())

	for _, m := range r.matches {
		runes = appendSlice(runes, m.Runes()...)
	}

	return
}

func (r *cMatch) String() string {
	if r.complete != nil {
		return *r.complete
	}

	var buf strings.Builder
	for _, m := range r.matches {
		// regardless of m.capture state
		buf.WriteString(m.String())
	}
	r.complete = values.Ref(buf.String())

	//r.complete = values.Ref(string(r.Runes()))
	return *r.complete
}

func (r *cMatch) Submatch(idx int) (value string, ok bool) {
	if ok = 0 <= idx && idx < len(r.submatch); ok {
		value = r.submatch[idx].String()
	}
	return
}
