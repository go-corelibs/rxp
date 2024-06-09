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

func (p Pattern) MatchString(input string) (ok bool) {
	if len(p) > 0 {
		s := newPatternState(p, input)
		ok = s.match(1)
		s.matches = nil
	}
	return
}

func (p Pattern) FindAllString(input string, count int) (found []string) {
	if len(p) > 0 {
		s := newPatternState(p, input)
		for _, m := range s.findString(count) {
			found = append(found, m[0]) // always at least one present
		}
		s.matches = nil
	}
	return
}

func (p Pattern) FindAllStringIndex(input string, count int) (found [][]int) {
	if len(p) > 0 {
		s := newPatternState(p, input)
		if s.match(count) {
			for _, mm := range s.matches {
				found = append(found, []int{mm[0][0], mm[0][1]})
			}
		}
		s.matches = nil
	}
	return
}

func (p Pattern) FindAllStringSubmatch(input string, count int) (found [][]string) {
	if len(p) > 0 {
		s := newPatternState(p, input)
		found = s.findString(count)
		s.matches = nil
	}
	return
}

func (p Pattern) FindAllStringSubmatchIndex(input string, count int) (found [][][2]int) {
	if len(p) > 0 {
		s := newPatternState(p, input)
		if s.match(count) {
			found = s.matches
		}
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
		if mm := s.findString(1); len(mm) > 0 {
			found = mm[0]
		}
		s.matches = nil
	}
	return
}

func (p Pattern) FindStringSubmatchIndex(input string, count int) (found [][2]int) {
	if len(p) > 0 {
		s := newPatternState(p, input)
		if s.match(count) && len(s.matches) > 0 {
			found = s.matches[0]
		}
		s.matches = nil
	}
	return
}

func (p Pattern) ReplaceAllString(input string, repl Replace) string {
	if len(p) > 0 {
		s := newPatternState(p, input)
		if s.match(-1) {
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
		if s.match(-1) {
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