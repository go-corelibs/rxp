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
	input   []rune  // original string input
	index   int     // current match position (total runes consumed)
	pattern Pattern // list of fragments to satisfy as a match
	capture []bool  // denotes corresponding matches are capture groups or not
	matches Matches // list of matches (with matched capture groups)
}

func newPatternState(p Pattern, input []rune) *cPatternState {
	return &cPatternState{
		input:   input,
		index:   0,
		pattern: p,
	}
}

func (s *cPatternState) findString(count int) (matched [][]string) {
	if s.match(count) {
		for _, match := range s.matches {
			if len(match) > 0 {
				var list []string
				for _, submatch := range match {
					// first index is the complete match
					list = appendSlice(
						list,
						string(s.input[submatch.Start():submatch.End()]),
					)
				}
				matched = appendSlice(matched, list)
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

	// while there is input to process
	for s.index < len(s.input) {

		start := s.index
		var subMatches SubMatches

		// for each matcher in the pattern
		for _, matcher := range s.pattern {
			// call each matcher once, expecting matcher to progress the index
			if cons, capt, _, proceed := matcher(gDefaultFlags, gDefaultReps, s.input, s.index); proceed {
				if capt {
					subMatches = appendSlice(subMatches, []int{s.index, s.index + cons})
				}
				s.index += cons
				completed += 1
				continue // continue the pattern
			}
			// pattern did not match correctly
			break
		}

		if completed >= required {
			if len(subMatches) > 0 {
				s.matches = appendSlice(s.matches, append(SubMatches{SubMatch{start, s.index}}, subMatches...))
			} else if start != s.index {
				s.matches = appendSlice(s.matches, SubMatches{SubMatch{start, s.index}})
			}
			start = s.index
			if count > 0 && count >= completed {
				// early out, optimized for Pattern.Segment() calls
				// count is the developer requested total number of subMatches
				return true
			}
		}

		if last == s.index {
			//	// pattern did not advance the index
			s.index += 1
		}
		last = s.index

	}

	return len(s.matches) > 0
}
