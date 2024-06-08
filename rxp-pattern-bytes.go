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

func (p Pattern) Match(input []byte) (ok bool) {
	if len(p) > 0 {
		s := newPatternState(p, input)
		defer s.recycle()
		ok = s.match(1)
		s.matches = nil
	}
	return
}

func (p Pattern) FindIndex(input string) (found []int) {
	if len(p) > 0 {
		s := newPatternState(p, input)
		defer s.recycle()
		if s.match(1) {
			mm := s.matches[0]
			found = []int{mm[0][0], mm[0][1]}
		}
		s.matches = nil
	}
	return
}

func (p Pattern) Scan(input []byte) (segments Segments) {
	s := newPatternState(p, input)
	defer s.recycle()
	return p.scanner(s)
}
