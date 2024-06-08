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

func (p Pattern) scanner(s *cPatternState) (segments Segments) {

	if len(p) == 0 {
		return []Segment{&cSegment{
			input:   s.input,
			matched: false,
			matches: SubMatches{SubMatch{0, s.input.Len()}},
		}}
	}

	if !s.match(-1) {
		return []Segment{&cSegment{
			input:   s.input,
			matched: false,
			matches: SubMatches{SubMatch{0, s.input.Len()}},
		}}
	}

	var last int
	for _, match := range s.matches {
		if last < match.Start() {
			segments = appendSlice[Segment](segments, &cSegment{
				input:   s.input,
				matched: false,
				matches: SubMatches{SubMatch{last, match.Start()}},
			})
		}
		last = match.End()

		segments = appendSlice[Segment](segments, &cSegment{
			input:   s.input,
			matched: true,
			matches: match,
		})
	}

	if last < s.input.Len() {
		segments = appendSlice[Segment](segments, &cSegment{
			input:   s.input,
			matched: false,
			matches: SubMatches{SubMatch{last, s.input.Len()}},
		})
	}
	s.matches = nil
	return
}
