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

// Pattern is a list of Matcher functions, all of which must match, in the
// order present, in order to consider the Pattern to match
type Pattern []Matcher

type cPatternState struct {
	input   *InputReader // input rune buffer
	index   int          // current match position (total runes consumed)
	pattern Pattern      // list of fragments to satisfy as a match
	capture []bool       // denotes corresponding matches are capture groups or not
	matches [][][2]int   // list of matches (with matched capture groups)
}

func newPatternState[V []rune | []byte | string](p Pattern, input V) *cPatternState {
	return &cPatternState{
		input:   NewInputReader(input),
		index:   0,
		pattern: p,
		matches: [][][2]int{},
	}
}

// match returns true if the state can process the Pattern at least count times
//
//gocyclo:ignore
func (p Pattern) match(s *cPatternState, count int) (matched bool) {

	var lastInputIndex int     // track Matcher progress
	var totalCompleted int     // track completed Matcher count
	required := len(s.pattern) // completed requirement

	set := [][2]int{{}}

	lastMatcherIdx := -1
	lastMatchedIdx := -1

	// while there is input to process
	for 0 <= s.index && s.index <= s.input.len {

		start := s.index
		var consumed int

		// for each matcher in the pattern
		var atLeastZero int
		for idx, matcher := range s.pattern {
			// call each matcher once, expecting matcher to progress the index
			if scoping, keep, proceed := matcher(DefaultFlags, gDefaultReps, s.input, s.index, set); proceed {
				if scoping.Capture() {
					set = pushSubMatch(set, [2]int{s.index, clamp(s.index+keep, s.input.len)})
				}
				if keep > 0 {
					consumed += keep
					s.index += keep
				} else if scoping.ZeroOrMore() {
					// if the previous index and this index are the same
					if lastMatcherIdx != idx || lastMatchedIdx != s.index || s.index == s.input.len {
						atLeastZero += 1
					}
				} else if scoping.Matched() {
					atLeastZero += 1
				}
				totalCompleted += 1
				lastMatcherIdx = idx
				lastMatchedIdx = s.index
				continue // continue the pattern
			}
			// pattern did not match correctly
			totalCompleted = 0
			if len(set) > 1 {
				set = [][2]int{{}}
			}
			break
		}

		if totalCompleted >= required {
			if start < s.index {
				set[0][0], set[0][1] = start, clamp(s.index, s.input.len)
				s.matches = pushMatch(s.matches, set)
				set = [][2]int{{}}
			} else if atLeastZero > 0 && 0 <= s.index && s.index <= s.input.len {
				set[0][0], set[0][1] = start, start
				s.matches = pushMatch(s.matches, set)
				set = [][2]int{{}}
			}
			if count > 0 {
				if count >= totalCompleted {
					// early out, count is the requested total number of subMatches
					return true
				}
			} else if 0 > s.index || s.index >= s.input.len {
				return len(s.matches) > 0
			}
		}

		if lastInputIndex == s.index {
			// pattern did not advance the index
			if _, size, ok := s.input.Get(s.index); ok && size > 0 {
				s.index += size // move the needle correctly
			} else {
				s.index += 1 // must move the needle to progress
			}
		}
		lastInputIndex = s.index

	}

	return len(s.matches) > 0
}
