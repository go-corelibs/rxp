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

func (p Pattern) findRunes(s *cPatternState, count int) (matched [][][]rune) {
	if p.match(s, count) {
		for _, match := range s.matches {
			if len(match) > 0 {
				var groups [][]rune
				for _, submatch := range match {
					slice, _ := s.input.Slice(submatch[0], submatch[1]-submatch[0])
					groups = pushRunes(groups, slice)
					//groups = append(groups, slice)
				}
				matched = pushRuneSlices(matched, groups)
				//matched = append(matched, groups)
			}
		}
	}
	return
}

// MatchRunes returns true if the input contains at least one match of this
// Pattern
func (p Pattern) MatchRunes(input []rune) (ok bool) {
	if len(p) > 0 {
		s := newPatternState(p, input)
		ok = p.match(s, 1)
		s.matches = nil
	}
	return
}

// FindRunes returns the leftmost Pattern match within the input given
func (p Pattern) FindRunes(input []rune) []rune {
	if mm := p.FindRunesSubmatch(input); len(mm) > 0 {
		return mm[0]
	}
	return nil
}

// FindRunesIndex returns the leftmost Pattern match starting and ending
// indexes within the string given
func (p Pattern) FindRunesIndex(input []rune) (found [2]int) {
	if mm := p.FindRunesSubmatchIndex(input); len(mm) > 0 {
		return mm[0]
	}
	return
}

// FindRunesSubmatch returns a slice of strings holding the leftmost match of
// this Pattern and any of its sub-matches. FindRunesSubmatch returns nil if
// there was no match of this Pattern within the input given
func (p Pattern) FindRunesSubmatch(input []rune) (found [][]rune) {
	if len(p) > 0 {
		s := newPatternState(p, input)
		if mm := p.findRunes(s, 1); len(mm) > 0 {
			found = mm[0]
		}
		s.matches = nil
	}
	return
}

// FindRunesSubmatchIndex returns a slice of starting and ending indices
// denoting the leftmost match of this Pattern and any of its sub-matches
func (p Pattern) FindRunesSubmatchIndex(input []rune) (found [][2]int) {
	if len(p) > 0 {
		s := newPatternState(p, input)
		if p.match(s, 1) && len(s.matches) > 0 {
			found = s.matches[0]
		}
		s.matches = nil
	}
	return
}

// FindAllRunes returns a slice of strings containing all of the Pattern
// matches present in the input given, in the order the matches are found
func (p Pattern) FindAllRunes(input []rune, count int) (found [][]rune) {
	if len(p) > 0 {
		s := newPatternState(p, input)
		for _, m := range p.findRunes(s, count) {
			found = pushRunes(found, m[0]) // always at least one present
		}
		s.matches = nil
	}
	return
}

// FindAllRunesIndex returns a slice of starting and ending indices denoting
// each of the Pattern matches present in the input given
func (p Pattern) FindAllRunesIndex(input []rune, count int) (found [][2]int) {
	if len(p) == 0 {
		return [][2]int{{0, 0}}
	}
	s := newPatternState(p, input)
	if p.match(s, count) && len(s.matches) > 0 {
		for _, groups := range s.matches {
			found = pushSubMatch(found, groups[0])
		}
	}
	s.matches = nil
	return
}

// FindAllRunesSubmatch returns a slice of all Pattern matches (and any
// sub-matches) present in the input given
func (p Pattern) FindAllRunesSubmatch(input []rune, count int) (found [][][]rune) {
	if len(p) > 0 {
		s := newPatternState(p, input)
		found = p.findRunes(s, count)
		s.matches = nil
	}
	return
}

// FindAllRunesSubmatchIndex returns a slice of starting and ending points for
// all Pattern matches (and any sub-matches) present in the input given
func (p Pattern) FindAllRunesSubmatchIndex(input []rune, count int) (found [][][2]int) {
	if len(p) > 0 {
		s := newPatternState(p, input)
		if p.match(s, count) {
			found = s.matches
		}
		s.matches = nil
	}
	return
}

// ReplaceAllRunes returns a copy of the input []rune with all Pattern matches
// replaced with text returned by the given Replace process
func (p Pattern) ReplaceAllRunes(input []rune, replacements Replace[[]rune]) (replaced []rune) {
	if len(p) > 0 {
		s := newPatternState(p, input)
		if p.match(s, -1) {
			var last int
			for _, group := range s.matches {
				if last < group[0][0] {
					// keep outsiders
					slice, _ := s.input.Slice(last, group[0][0]-last)
					replaced = pushRune(replaced, slice)
				}
				last = group[0][1] // track outsider space

				slice, _ := s.input.Slice(group[0][0], group[0][1]-group[0][0])
				replaced = pushRune(replaced, replacements.process(s.input, group, slice))

			}
			if last < s.input.Len() {
				// keep outsiders
				slice, _ := s.input.Slice(last, s.input.len-last)
				replaced = pushRune(replaced, slice)
			}
			s.matches = nil
			return
		}
	}
	return input
}

// ReplaceAllRunesFunc returns a copy of the input []rune with all Pattern
// matches replaced with the text returned by the given Transform function
func (p Pattern) ReplaceAllRunesFunc(input []rune, transform Transform[[]rune]) (replaced []rune) {
	if len(p) > 0 {
		s := newPatternState(p, input)
		if p.match(s, -1) {
			var last int
			for _, group := range s.matches {
				if last < group[0][0] {
					slice, _ := s.input.Slice(last, group[0][0]-last)
					replaced = pushRune(replaced, slice)
				}
				last = group[0][1]

				slice, _ := s.input.Slice(group[0][0], group[0][1]-group[0][0])
				replaced = pushRune(replaced, transform(slice))
			}
			if last < s.input.Len() {
				slice, _ := s.input.Slice(last, s.input.len-last)
				replaced = pushRune(replaced, slice)
			}
			s.matches = nil
			return
		}
	}
	return input
}

// ReplaceAllLiteralRunes returns a copy of the input []rune with all Pattern
// matches replaced with the unmodified replacement text
func (p Pattern) ReplaceAllLiteralRunes(input []rune, replacement []rune) (replaced []rune) {
	if len(p) > 0 {
		s := newPatternState(p, input)
		if p.match(s, -1) {
			var last int
			for _, group := range s.matches {
				if last < group[0][0] {
					// keep outsiders
					slice, _ := s.input.Slice(last, group[0][0]-last)
					replaced = pushRune(replaced, slice)
				}
				last = group[0][1] // track outsider space
				replaced = pushRune(replaced, replacement)
			}
			if last < s.input.Len() {
				// keep outsiders
				slice, _ := s.input.Slice(last, s.input.Len()-last)
				replaced = pushRune(replaced, slice)
			}
			s.matches = nil
			return
		}
	}
	return input
}

func (p Pattern) SplitRunes(input []rune, count int) (found [][]rune) {
	if len(p) > 0 && len(input) == 0 {
		return [][]rune(nil)
	} else if count == 0 || len(p) == 0 {
		return // zero substrings
	}

	s := newPatternState(p, input)
	// get the complete match set first
	if p.match(s, -1) {

		beg := 0
		end := 0
		for _, match := range s.matches {
			if count > 0 && len(found) >= count-1 {
				break
			}

			end = match[0][0]
			if match[0][1] != 0 {
				found = append(found, input[beg:end])
			}
			beg = match[0][1]
		}

		if end != len(input) {
			found = append(found, input[beg:])
		}

	} else {
		found = [][]rune{input}
	}
	s.matches = nil
	return
}
