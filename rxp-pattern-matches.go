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

// Matches is the list of all SubMatches for each of the Pattern matches
type Matches []SubMatches

func (m Matches) Valid() bool {
	for _, sm := range m {
		if !sm.Valid() {
			return false
		}
	}
	return true
}

func (m Matches) Get(index int) (sm SubMatches, ok bool) {
	if ok = 0 <= index && index < len(m); ok {
		sm = m[index]
	}
	return
}

func (m Matches) Len() int {
	return len(m)
}

func (m Matches) Size() int {
	if lm := len(m); lm > 0 {
		// this may look odd because it's the start of the first and the end
		// of the last and yet it also catches the end of the first if there's
		// only one SubMatches present
		return m[lm-1].End() - m[0].Start()
	}
	return 0
}

func (m Matches) Start() int {
	if len(m) > 0 {
		return m[0].Start()
	}
	return 0
}

func (m Matches) End() int {
	if lm := len(m); lm > 0 {
		return m[lm-1].End()
	}
	return 0
}

// SubMatches is the list of SubMatch indices for a single matched Pattern
type SubMatches []SubMatch

func (m SubMatches) Valid() bool {
	for _, sm := range m {
		if !sm.Valid() {
			return false
		}
	}
	return true
}

func (m SubMatches) Get(index int) (sm SubMatch, ok bool) {
	if 0 <= index && index < len(m) {
		return m[index], true
	}
	return
}

func (m SubMatches) Len() int {
	return len(m)
}

func (m SubMatches) Size() int {
	if lm := len(m); lm > 0 {
		if lmm := len(m[lm-1]); lmm > 0 {
			return m[lm-1].End() - m[0].Start()
		}
	}
	return 0
}

func (m SubMatches) Start() int {
	if len(m) > 0 && len(m[0]) > 0 {
		return m[0].Start()
	}
	return 0
}

func (m SubMatches) End() int {
	if lm := len(m); lm > 0 {
		if lmm := len(m[lm-1]); lmm > 0 {
			return m[lm-1].End()
		}
	}
	return 0
}

// SubMatch is a single pair of start and end input rune indices
type SubMatch []int

func (m SubMatch) Valid() bool {
	return len(m) > 1
}

func (m SubMatch) Len() int {
	if len(m) > 1 {
		return m[1] - m[0]
	}
	return 0
}

func (m SubMatch) Start() int {
	if len(m) > 0 {
		return m[0]
	}
	return 0
}

func (m SubMatch) End() int {
	if lm := len(m); lm > 1 {
		return m[lm-1]
	}
	return 0
}
