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

type cPatternState struct {
	input   *RuneBuffer // input rune buffer
	index   int         // current match position (total runes consumed)
	pattern Pattern     // list of fragments to satisfy as a match
	capture []bool      // denotes corresponding matches are capture groups or not
	matches [][][2]int  // list of matches (with matched capture groups)
}

func newPatternState[V []rune | []byte | string](p Pattern, input V) *cPatternState {
	return &cPatternState{
		input:   NewRuneBuffer(input),
		index:   0,
		pattern: p,
		matches: [][][2]int{},
	}
}

func (s *cPatternState) findString(count int) (matched [][]string) {
	if s.match(count) {
		for _, match := range s.matches {
			if len(match) > 0 {
				var groups []string
				for _, submatch := range match {
					groups = pushString(groups, s.input.String(submatch[0], submatch[1]-submatch[0]))
				}
				matched = pushStrings(matched, groups)
			}
		}
	}
	return
}

// match returns true if the state can process the Pattern at least count times
func (s *cPatternState) match(count int) (matched bool) {

	var last int               // track Matcher progress
	var completed int          // track completed Matcher count
	required := len(s.pattern) // completed requirement

	set := [][2]int{{}}

	// while there is input to process
	for s.input.Valid(s.index) {

		start := s.index

		// for each matcher in the pattern
		for _, matcher := range s.pattern {
			// call each matcher once, expecting matcher to progress the index
			if cons, capt, _, proceed := matcher(DefaultFlags, gDefaultReps, s.input, s.index, set); proceed {
				if capt {
					set = pushSMSlice(set, s.index, s.index+cons)
				}
				s.index += cons
				completed += 1
				continue // continue the pattern
			}
			// pattern did not match correctly
			completed = 0
			if len(set) > 1 {
				set = [][2]int{{}}
			}
			break
		}

		if completed >= required {
			if len(set) > 1 || start != s.index {
				set[0][0] = start
				set[0][1] = s.index
				s.matches = pushMatch(s.matches, set)
				set = [][2]int{{}}
			}
			if count > 0 && count >= completed {
				// early out, optimized for Pattern.Segment() calls
				// count is the developer requested total number of subMatches
				return true
			}
		}

		if last == s.index {
			// pattern did not advance the index
			if _, size, ok := s.input.Get(s.index); ok && size > 0 {
				s.index += size // move the needle correctly
			} else {
				s.index += 1 // must move the needle to progress
			}
		}
		last = s.index

	}

	return len(s.matches) > 0
}
