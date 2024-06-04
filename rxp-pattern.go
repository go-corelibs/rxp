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

func (p Pattern) Match(input []rune) (ok bool) {
	if len(p) > 0 {
		s := newPatternState(p, input)
		ok = s.match(1)
		s.matches.Recycle()
		s.matches = nil
	}
	return
}

func (p Pattern) MatchString(input string) (ok bool) {
	return p.Match([]rune(input))
}

//func (p Pattern) Find(input []rune) (found []rune) {
//	if mm := p.FindSubmatch(input); len(mm) > 0 {
//		found = mm[0]
//	}
//	return
//}

//func (p Pattern) FindAll(input []rune, count int) (found [][]rune) {
//	if len(p) > 0 {
//		s := newPatternState(p, input)
//		for _, m := range s.find(count) {
//			found = append(found, m[0]) // always at least one present
//		}
//	}
//	return
//}

//func (p Pattern) FindAllIndex(input []rune, count int) (found [][]int) {
//	if len(p) > 0 {
//		s := newPatternState(p, input)
//		if s.match(count) {
//			for _, mm := range s.matches {
//				found = append(found, []int{mm[0].start, mm[0].start + mm[0].this})
//			}
//		}
//	}
//	return
//}

func (p Pattern) FindAllString(input string, count int) (found []string) {
	if len(p) > 0 {
		s := newPatternState(p, []rune(input))
		for _, m := range s.findString(count) {
			found = append(found, m[0]) // always at least one present
		}
		s.matches.Recycle()
		s.matches = nil
	}
	return
}

func (p Pattern) FindAllStringIndex(input string, count int) (found [][]int) {
	if len(p) > 0 {
		s := newPatternState(p, []rune(input))
		if s.match(count) {
			for _, mm := range s.matches {
				found = append(found, []int{mm[0].start, mm[0].start + mm[0].this})
			}
		}
		s.matches.Recycle()
		s.matches = nil
	}
	return
}

func (p Pattern) FindAllStringSubmatch(input string, count int) (found [][]string) {
	if len(p) > 0 {
		s := newPatternState(p, []rune(input))
		found = s.findString(count)
		s.matches.Recycle()
		s.matches = nil
	}
	return
}

//func (p Pattern) FindSubmatch(input []rune) (found [][]rune) {
//	if len(p) > 0 {
//		s := newPatternState(p, input)
//		if mm := s.find(1); len(mm) > 0 {
//			found = mm[0]
//		}
//	}
//	return
//}

func (p Pattern) FindString(input string) string {
	if mm := p.FindStringSubmatch(input); len(mm) > 0 {
		return mm[0]
	}
	return ""
}

func (p Pattern) FindStringSubmatch(input string) (found []string) {
	if len(p) > 0 {
		s := newPatternState(p, []rune(input))
		if mm := s.findString(1); len(mm) > 0 {
			found = mm[0]
		}
		s.matches.Recycle()
		s.matches = nil
	}
	return
}

func (p Pattern) FindIndex(input string) (found []int) {
	if len(p) > 0 {
		s := newPatternState(p, []rune(input))
		if s.match(1) {
			mm := s.matches[0]
			found = []int{mm[0].start, mm[0].start + mm[0].this}
		}
		s.matches.Recycle()
		s.matches = nil
	}
	return
}

func (p Pattern) ReplaceAllString(input string, repl Replace) string {
	if len(p) == 0 {
		return input
	}
	s := newPatternState(p, []rune(input))
	fragments := s.scan(-1)
	buf := getStringsBuilder()
	defer putStringsBuilder(buf)
	for _, frag := range fragments {
		if m, ok := frag.Match(); ok {
			for _, r := range repl {
				buf.WriteString(r(m))
			}
		} else {
			buf.WriteString(frag.String())
		}
		frag.Recycle()
	}
	s.matches.Recycle()
	s.matches = nil
	return buf.String()
}

func (p Pattern) ReplaceAllStringFunc(input string, repl Transform) string {
	if len(p) == 0 {
		return input
	}
	s := newPatternState(p, []rune(input))
	fragments := s.scan(-1)
	buf := getStringsBuilder()
	defer putStringsBuilder(buf)
	for _, frag := range fragments {
		if result, ok := frag.Match(); ok {
			buf.WriteString(repl(result.String()))
		} else {
			buf.WriteString(frag.String())
		}
		frag.Recycle()
	}
	s.matches.Recycle()
	s.matches = nil
	return buf.String()
}

func (p Pattern) ScanString(input string) (fragments Fragments) {
	if len(p) > 0 {
		s := newPatternState(p, []rune(input))
		fragments = s.scan(-1)
		s.matches = nil
	}
	return
}
