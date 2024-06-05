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
	input   []rune      // original string input
	index   int         // current match position (total runes consumed)
	pattern Pattern     // list of fragments to satisfy as a match
	matches matchGroups // total list of completed Pattern matches
}

func newPatternState(p Pattern, input []rune) *cPatternState {
	//spConfig.Seed(spConfig.Scale() * 10)
	//spFragment.Seed(spFragment.Scale() * 10)
	//spMatchState.Seed(spMatchState.Scale() * 1)
	return &cPatternState{
		input:   input,
		index:   0,
		pattern: p,
	}
}

//func (s *cPatternState) find(count int) (matched [][][]rune) {
//	if s.match(count) {
//		for _, mm := range s.matches {
//			list := [][]rune{{}} // first index is the complete match
//			for _, m := range mm {
//				this := s.input[m.start : m.start+m.this]
//				list[0] = appendSlice(list[0], this...)
//				if m.capture {
//					list = appendSlice(list, this)
//				}
//			}
//			matched = appendSlice(matched, list)
//		}
//	}
//	return
//}

func (s *cPatternState) findString(count int) (matched [][]string) {
	if s.match(count) {
		for _, mm := range s.matches {
			list := []string{""} // first index is the complete match
			for _, m := range mm {
				this := string(s.input[m.start : m.start+m.this])
				list[0] += this
				if m.flags.Capture() {
					list = appendSlice(list, this)
				}
			}
			matched = appendSlice(matched, list)
		}
	}
	return
}

// match returns true if the state can process the Pattern at least count times
func (s *cPatternState) match(count int) (matched bool) {

	var lastIndex int          // track Matcher progress
	var completed int          // track completed Matcher count
	required := len(s.pattern) // completed requirement

	// while there is input to process
	for s.index < len(s.input) {

		var matches matchStates

		// for each matcher in the pattern
		for _, matcher := range s.pattern {
			m := newMatchState(s, s.index)
			// call each matcher once, expecting matcher to progress the index
			if next, keep := matcher(m); next && keep {
				completed += 1
				s.index += m.this
				m.complete = true
				matches = appendSlice(matches, m)
				continue // continue the pattern
			} else if next {
				completed += 1
				continue // did not progress though go to next
			}
			// pattern did not match correctly
			break
		}

		if completed >= required {
			if len(matches) > 0 {
				s.matches = appendSlice(s.matches, matches)
			}
			if count > 0 && count >= completed {
				// early out, optimized for Pattern.Match() calls
				// count is the developer requested total number of matches
				return true
			}
		}

		if lastIndex == s.index {
			//	// pattern did not advance the index
			s.index += 1
		}
		lastIndex = s.index

	}

	return len(s.matches) > 0
}

func (s *cPatternState) scan(count int) (fragments Fragments) {

	if !s.match(count) {
		// no matches
		return Fragments{
			newFragment(s, 0, len(s.input), nil),
		}
	}

	fragments = make(Fragments, 0, (len(s.matches)+1)*4)
	total := len(s.input)
	var lastEnd int
	for midx, m := range s.matches {
		start, end := m.start(), m.end()

		// include preceeding non-match fragments

		if midx == 0 {
			// this is the first match
			if start > 0 {
				// and it starts after the first input rune
				fragments = appendSlice(fragments, newFragment(s, 0, start, nil))
			}
		} else {
			// this is not the first match
			if start-1 >= lastEnd {
				// and there's a gap between the previous and this match
				fragments = appendSlice(fragments, newFragment(s, lastEnd, start, nil))
			}
		}

		fragments = appendSlice(fragments, newFragment(s, start, end, newMatch(s, m)))

		lastEnd = end
	}

	// make sure any leftovers aren't actually left behind
	if lastEnd < total {
		fragments = appendSlice(fragments, newFragment(s, lastEnd, total, nil))
	}
	return
}
