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

func (p Pattern) findString(s *cPatternState, count int) (matched [][]string) {
	if p.match(s, count) {
		for _, match := range s.matches {
			if len(match) > 0 {
				var groups []string
				for _, submatch := range match {
					groups = pushString(groups, s.input.String(submatch[0], submatch[1]-submatch[0]))
					//groups = append(groups, s.input.String(submatch[0], submatch[1]-submatch[0]))
				}
				matched = pushStrings(matched, groups)
				//matched = append(matched, groups)
			}
		}
	}
	return
}

// MatchString returns true if the input contains at least one match of this
// Pattern
func (p Pattern) MatchString(input string) (ok bool) {
	if len(p) > 0 {
		s := newPatternState(p, input)
		ok = p.match(s, 1)
		s.matches = nil
	}
	return
}

// FindString returns the leftmost Pattern match within the input given
func (p Pattern) FindString(input string) string {
	if mm := p.FindStringSubmatch(input); len(mm) > 0 {
		return mm[0]
	}
	return ""
}

// FindStringIndex returns the leftmost Pattern match starting and ending
// indexes within the string given
func (p Pattern) FindStringIndex(input string) (found [2]int) {
	if mm := p.FindStringSubmatchIndex(input); len(mm) > 0 {
		return mm[0]
	}
	return
}

// FindStringSubmatch returns a slice of strings holding the leftmost match of
// this Pattern and any of its sub-matches. FindStringSubmatch returns nil if
// there was no match of this Pattern within the input given
func (p Pattern) FindStringSubmatch(input string) (found []string) {
	if len(p) > 0 {
		s := newPatternState(p, input)
		if mm := p.findString(s, 1); len(mm) > 0 {
			found = mm[0]
		}
		s.matches = nil
	}
	return
}

// FindStringSubmatchIndex returns a slice of starting and ending indices
// denoting the leftmost match of this Pattern and any of its sub-matches
func (p Pattern) FindStringSubmatchIndex(input string) (found [][2]int) {
	if len(p) > 0 {
		s := newPatternState(p, input)
		if p.match(s, 1) && len(s.matches) > 0 {
			found = s.matches[0]
		}
		s.matches = nil
	}
	return
}

// FindAllString returns a slice of strings containing all of the Pattern
// matches present in the input given, in the order the matches are found
func (p Pattern) FindAllString(input string, count int) (found []string) {
	if len(p) > 0 {
		s := newPatternState(p, input)
		for _, m := range p.findString(s, count) {
			found = append(found, m[0]) // always at least one present
		}
		s.matches = nil
	}
	return
}

// FindAllStringIndex returns a slice of starting and ending indices denoting
// each of the Pattern matches present in the input given
func (p Pattern) FindAllStringIndex(input string, count int) (found [][2]int) {
	if len(p) == 0 {
		return [][2]int{{0, 0}}
	}
	s := newPatternState(p, input)
	if p.match(s, count) && len(s.matches) > 0 {
		for _, groups := range s.matches {
			found = append(found, groups[0])
		}
	}
	s.matches = nil
	return
}

// FindAllStringSubmatch returns a slice of all Pattern matches (and any
// sub-matches) present in the input given
func (p Pattern) FindAllStringSubmatch(input string, count int) (found [][]string) {
	if len(p) > 0 {
		s := newPatternState(p, input)
		found = p.findString(s, count)
		s.matches = nil
	}
	return
}

// FindAllStringSubmatchIndex returns a slice of starting and ending points for
// all Pattern matches (and any sub-matches) present in the input given
func (p Pattern) FindAllStringSubmatchIndex(input string, count int) (found [][][2]int) {
	if len(p) > 0 {
		s := newPatternState(p, input)
		if p.match(s, count) {
			found = s.matches
		}
		s.matches = nil
	}
	return
}

// ReplaceAllString returns a copy of the input string with all Pattern matches
// replaced with text returned by the given Replace process
func (p Pattern) ReplaceAllString(input string, replacements Replace[string]) string {
	if len(p) > 0 {
		s := newPatternState(p, input)
		if p.match(s, -1) {
			buf := spStringBuilder.Get()
			defer spStringBuilder.Put(buf)
			var last int
			for _, match := range s.matches {
				if last < match[0][0] {
					// keep outsiders
					buf.WriteString(s.input.String(last, match[0][0]-last))
				}
				last = match[0][1] // track outsider space

				buf.WriteString(replacements.process(
					s.input,
					match,
					s.input.String(match[0][0], match[0][1]-match[0][0]),
				))

			}
			if last < s.input.Len() {
				// keep outsiders
				buf.WriteString(s.input.String(last, s.input.Len()-last))
			}
			s.matches = nil
			return buf.String()
		}
	}
	return input
}

// ReplaceAllStringFunc returns a copy of the input string with all Pattern
// matches replaced with the text returned by the given Transform function
func (p Pattern) ReplaceAllStringFunc(input string, transform Transform[string]) string {
	if len(p) > 0 {
		s := newPatternState(p, input)
		if p.match(s, -1) {
			buf := spStringBuilder.Get()
			defer spStringBuilder.Put(buf)
			var last int
			for _, group := range s.matches {
				if last < group[0][0] {
					buf.WriteString(s.input.String(last, group[0][0]-last))
				}
				last = group[0][1]
				buf.WriteString(transform(s.input.String(group[0][0], group[0][1]-group[0][0])))
			}
			if last < s.input.Len() {
				buf.WriteString(s.input.String(last, s.input.Len()-last))
			}
			s.matches = nil
			return buf.String()
		}
	}
	return input
}

// ReplaceAllLiteralString returns a copy of the input string with all Pattern
// matches replaced with the unmodified replacement text
func (p Pattern) ReplaceAllLiteralString(input string, replacement string) string {
	if len(p) > 0 {
		s := newPatternState(p, input)
		if p.match(s, -1) {
			buf := spStringBuilder.Get()
			defer spStringBuilder.Put(buf)
			var last int
			for _, match := range s.matches {
				if last < match[0][0] {
					// keep outsiders
					buf.WriteString(s.input.String(last, match[0][0]-last))
				}
				last = match[0][1] // track outsider space
				buf.WriteString(replacement)
			}
			if last < s.input.Len() {
				// keep outsiders
				buf.WriteString(s.input.String(last, s.input.Len()-last))
			}
			s.matches = nil
			return buf.String()
		}
	}
	return input
}

// SplitString slices the input into substrings separated by the Pattern and
// returns a slice of the substrings between those Pattern matches
//
// The slice returned by this method consists of all the substrings of input
// not contained in the slice returned by [Pattern.FindAllString]
//
// Example:
//
//	s := regexp.MustCompile("a*").SplitString("abaabaccadaaae", 5)
//	// s: ["", "b", "b", "c", "cadaaae"]
//
// The count determines the number of substrings to return:
//
//	| Value Case | Description                                                       |
//	|------------|-------------------------------------------------------------------|
//	| count > 0  | at most count substrings; the last will be the un-split remainder |
//	| count == 0 | the result is nil (zero substrings)                               |
//	| count < 0  | all substrings                                                    |
func (p Pattern) SplitString(input string, count int) (found []string) {
	if len(p) > 0 && input == "" {
		return []string{""}
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
		found = []string{input}
	}
	s.matches = nil
	return
}
