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

func (p Pattern) findBytes(s *cPatternState, count int) (matched [][][]byte) {
	if p.match(s, count) {
		for _, match := range s.matches {
			if len(match) > 0 {
				var groups [][]byte
				for _, submatch := range match {
					groups = pushByte(groups, s.input.Bytes(submatch[0], submatch[1]-submatch[0]))
					//groups = append(groups, s.input.Bytes(submatch[0], submatch[1]-submatch[0]))
				}
				matched = pushBytes(matched, groups)
				//matched = append(matched, groups)
			}
		}
	}
	return
}

func (p Pattern) MatchBytes(input []byte) (ok bool) {
	if len(p) > 0 {
		s := newPatternState(p, input)
		ok = p.match(s, 1)
		s.matches = nil
	}
	return
}

func (p Pattern) FindBytes(input []byte) []byte {
	if mm := p.FindBytesSubmatch(input); len(mm) > 0 {
		return mm[0]
	}
	return nil
}

func (p Pattern) FindBytesIndex(input []byte) (found [2]int) {
	if mm := p.FindBytesSubmatchIndex(input); len(mm) > 0 {
		return mm[0]
	}
	return
}

func (p Pattern) FindBytesSubmatch(input []byte) (found [][]byte) {
	if len(p) > 0 {
		s := newPatternState(p, input)
		if mm := p.findBytes(s, 1); len(mm) > 0 {
			found = mm[0]
		}
		s.matches = nil
	}
	return
}

func (p Pattern) FindBytesSubmatchIndex(input []byte) (found [][2]int) {
	if len(p) > 0 {
		s := newPatternState(p, input)
		if p.match(s, 1) && len(s.matches) > 0 {
			found = s.matches[0]
		}
		s.matches = nil
	}
	return
}

func (p Pattern) FindAllBytes(input []byte, count int) (found [][]byte) {
	if len(p) > 0 {
		s := newPatternState(p, input)
		for _, m := range p.findBytes(s, count) {
			found = pushByte(found, m[0]) // always at least one present
		}
		s.matches = nil
	}
	return
}

func (p Pattern) FindAllBytesIndex(input []byte, count int) (found [][2]int) {
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

func (p Pattern) FindAllBytesSubmatch(input []byte, count int) (found [][][]byte) {
	if len(p) > 0 {
		s := newPatternState(p, input)
		found = p.findBytes(s, count)
		s.matches = nil
	}
	return
}

func (p Pattern) FindAllBytesSubmatchIndex(input []byte, count int) (found [][][2]int) {
	if len(p) > 0 {
		s := newPatternState(p, input)
		if p.match(s, count) {
			found = s.matches
		}
		s.matches = nil
	}
	return
}

func (p Pattern) ReplaceAllBytes(input []byte, replacements Replace[[]byte]) []byte {
	if len(p) > 0 {
		s := newPatternState(p, input)
		if p.match(s, -1) {
			buf := spBytesBuffer.Get()
			defer spBytesBuffer.Put(buf)
			var last int
			for _, match := range s.matches {
				if last < match[0][0] {
					// keep outsiders
					buf.Write(s.input.Bytes(last, match[0][0]-last))
				}
				last = match[0][1] // track outsider space

				buf.Write(replacements.process(
					s.input,
					match,
					s.input.Bytes(match[0][0], match[0][1]-match[0][0]),
				))

			}
			if last < s.input.Len() {
				// keep outsiders
				buf.Write(s.input.Bytes(last, s.input.Len()-last))
			}
			s.matches = nil
			return buf.Bytes()
		}
	}
	return input
}

func (p Pattern) ReplaceAllBytesFunc(input []byte, transform Transform[[]byte]) []byte {
	if len(p) > 0 {
		s := newPatternState(p, input)
		if p.match(s, -1) {
			buf := spBytesBuffer.Get()
			defer spBytesBuffer.Put(buf)
			var last int
			for _, group := range s.matches {
				if last < group[0][0] {
					buf.Write(s.input.Bytes(last, group[0][0]-last))
				}
				last = group[0][1]
				buf.Write(transform(s.input.Bytes(group[0][0], group[0][1]-group[0][0])))
			}
			if last < s.input.Len() {
				buf.Write(s.input.Bytes(last, s.input.Len()-last))
			}
			s.matches = nil
			return buf.Bytes()
		}
	}
	return input
}

func (p Pattern) ReplaceAllLiteralBytes(input []byte, replacement []byte) []byte {
	if len(p) > 0 {
		s := newPatternState(p, input)
		if p.match(s, -1) {
			buf := spBytesBuffer.Get()
			defer spBytesBuffer.Put(buf)
			var last int
			for _, match := range s.matches {
				if last < match[0][0] {
					// keep outsiders
					buf.Write(s.input.Bytes(last, match[0][0]-last))
				}
				last = match[0][1] // track outsider space
				buf.Write(replacement)
			}
			if last < s.input.Len() {
				// keep outsiders
				buf.Write(s.input.Bytes(last, s.input.Len()-last))
			}
			s.matches = nil
			return buf.Bytes()
		}
	}
	return input
}

func (p Pattern) SplitBytes(input []byte, count int) (found [][]byte) {
	if len(p) > 0 && len(input) == 0 {
		return [][]byte(nil)
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
		found = [][]byte{input}
	}
	s.matches = nil
	return
}
