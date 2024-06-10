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

// Split slices the input into substrings separated by the Pattern and returns a
// slice of the substrings between those Pattern matches.
//
// The slice returned by this method consists of all the substrings of input
// not contained in the slice returned by [Pattern.FindAllString].
//
// Example:
//
//	s := regexp.MustCompile("a*").Split("abaabaccadaaae", 5)
//	// s: ["", "b", "b", "c", "cadaaae"]
//
// The count determines the number of substrings to return:
//
//	| Value Case | Description                                                       |
//	|------------|-------------------------------------------------------------------|
//	| count > 0  | at most count substrings; the last will be the un-split remainder |
//	| count == 0 | the result is nil (zero substrings)                               |
//	| count < 0  | all substrings                                                    |
func (p Pattern) Split(input string, count int) (found []string) {
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

	}
	s.matches = nil
	return
}

func (p Pattern) MatchString(input string) (ok bool) {
	if len(p) > 0 {
		s := newPatternState(p, input)
		ok = p.match(s, 1)
		s.matches = nil
	}
	return
}

func (p Pattern) FindString(input string) string {
	if mm := p.FindStringSubmatch(input); len(mm) > 0 {
		return mm[0]
	}
	return ""
}

func (p Pattern) FindStringIndex(input string, count int) (found [2]int) {
	if mm := p.FindStringSubmatchIndex(input, count); len(mm) > 0 {
		return mm[0]
	}
	return
}

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

func (p Pattern) FindStringSubmatchIndex(input string, count int) (found [][2]int) {
	if len(p) > 0 {
		s := newPatternState(p, input)
		if p.match(s, count) && len(s.matches) > 0 {
			found = s.matches[0]
		}
		s.matches = nil
	}
	return
}

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

func (p Pattern) FindAllStringIndex(input string, count int) (found [][]int) {
	if len(p) == 0 {
		return [][]int{{0, 0}}
	}
	s := newPatternState(p, input)
	if p.match(s, count) && len(s.matches) > 0 {
		for _, groups := range s.matches {
			found = append(found, []int{groups[0][0], groups[0][1]})
		}
	}
	s.matches = nil
	return
}

func (p Pattern) FindAllStringSubmatch(input string, count int) (found [][]string) {
	if len(p) > 0 {
		s := newPatternState(p, input)
		found = p.findString(s, count)
		s.matches = nil
	}
	return
}

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

func (p Pattern) ReplaceAllString(input string, repl Replace) string {
	if len(p) > 0 {
		s := newPatternState(p, input)
		if p.match(s, -1) {
			buf := spStringBuilder.Get()
			defer spStringBuilder.Put(buf)
			var last int
			for _, match := range s.matches {
				if last < match[0][0] {
					buf.WriteString(s.input.String(last, match[0][0]-last))
				}
				last = match[0][1]
				for _, rpl := range repl {
					buf.WriteString(rpl(&cSegment{
						input:   s.input,
						matched: true,
						matches: match,
					}))
				}
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

func (p Pattern) ReplaceAllStringFunc(input string, repl Transform) string {
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
				buf.WriteString(repl(s.input.String(group[0][0], group[0][1]-group[0][0])))
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

func (p Pattern) ScanStrings(input string) (segments Segments) {
	s := newPatternState(p, input)
	return p.scanner(s)
}
