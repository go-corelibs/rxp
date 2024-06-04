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

var (
	_ MatchState = (*cMatchState)(nil)
)

var (
	spMatchState = newSyncPool(1, func() *cMatchState {
		return new(cMatchState)
	})
)

type IMatchState interface {
	// Input returns the complete input rune slice
	Input() []rune
	// InputLen returns the total number of input runes (not bytes)
	InputLen() int
	// Index returns the start position for this Matcher, use Index + Len to
	// get the correct position within Input for This rune
	Index() int
	// Len is count of runes consumed so far by this Matcher
	Len() int
	// End returns true if this MatchState position (index+len) is exactly the
	// InputLen, denoting the Dollar zero-width position (End Of Input)
	End() bool
	// Ready returns true if this MatchState position (index+len) is less than
	// the InputLen (Ready is before End, End is Valid, past End is Invalid)
	Ready() bool
	// Valid returns true if this MatchState position (index+len) is less than
	// or equal to the InputLen (End is Valid, past the End is Invalid)
	Valid() bool
	// Invalid returns true if this MatchState position (index+len) is greater
	// than to the InputLen (End is Valid, past the End is Invalid)
	Invalid() bool
	// Has returns true if the given index is an input rune
	Has(idx int) (ok bool)
	// Get returns the input rune at the given index, if there is one
	Get(idx int) (r rune, ok bool)
	// Prev returns the previous rune, if there is one
	Prev() (r rune, ok bool)
	// This is the current rune being considered for this Matcher, if there is one
	This() (r rune, ok bool)
	// Next peeks at the next rune, if there is one
	Next() (r rune, ok bool)
	// Captured returns true if Capture was previously called on this MatchState
	Captured() bool
	// Negated returns true if the scoped Config is negated
	Negated() bool

	Runes() []rune
	String() (text string)

	// Equal returns true if the other MatchState is the same as this one
	Equal(other MatchState) bool
}

// MatchState is used by Matcher functions to progress the Pattern matching process
type MatchState interface {
	IMatchState

	// Capture includes this Matcher in any sub-matching results
	Capture()

	// Consume associates the current rune with this match and moves the MatchState
	// forward by count runes, Keep returns true if there are still more runes
	// remaining to process
	Consume(count int) (ok bool)

	// Clone is a wrapper around a zero offset call to CloneWith
	Clone() MatchState

	// CloneWith returns an exact copy of this MatchState with the start point
	// offset by the positive amount given
	CloneWith(offset int) MatchState

	// Apply updates the other MatchState with this MatchState consumed runes,
	// capture state and Config Scope
	Apply(other MatchState)

	// Scope append (if the cfg is not nil) a new Config scope within this
	// matcher and always returns a composite Config of all scoped configs
	Scope(cfg *Config) *Config

	Recycle()
	private(_ *cMatchState)
}

// match is the MatchState of a Pattern Matcher
type cMatchState struct {
	s        *cPatternState
	this     int  // current rune offset from the start
	start    int  // starting rune index within the input string
	capture  bool // include as a capture group
	complete bool // did the match complete successfully
	scoped   *Config
	recycled bool
}

func newMatchState(s *cPatternState, start int) *cMatchState {
	state := spMatchState.Get()
	if !state.recycled {
		// pool created a new one, let's add one more so next isn't new either
		spMatchState.Seed(-1)
	}
	state.s = s
	state.this = 0
	state.start = start
	state.capture = false
	state.complete = false
	state.scoped = nil
	state.recycled = false
	return state
}

func (m *cMatchState) Recycle() {
	if !m.recycled {
		m.s = nil
		if m.scoped != nil {
			m.scoped.Recycle()
		}
		m.recycled = true
		spMatchState.Put(m)
	}
}

func (m *cMatchState) private(_ *cMatchState) {}

func (m *cMatchState) Input() []rune {
	return m.s.input[:]
}

func (m *cMatchState) InputLen() int {
	return len(m.s.input)
}

func (m *cMatchState) Index() int {
	return m.start
}

func (m *cMatchState) Len() int {
	return m.this
}

func (m *cMatchState) End() bool {
	return m.start+m.this == len(m.s.input)
}

func (m *cMatchState) Ready() bool {
	return m.start+m.this >= 0 && m.start+m.this < len(m.s.input)
}

func (m *cMatchState) Valid() bool {
	return m.start+m.this >= 0 && m.start+m.this <= len(m.s.input)
}

func (m *cMatchState) Invalid() bool {
	return m.start+m.this < 0 || m.start+m.this > len(m.s.input)
}

func (m *cMatchState) Has(idx int) (ok bool) {
	ok = 0 <= idx && idx < len(m.s.input)
	return
}

func (m *cMatchState) Get(idx int) (r rune, ok bool) {
	if ok = m.Has(idx); ok {
		r = m.s.input[idx]
	}
	return
}

func (m *cMatchState) Prev() (r rune, ok bool) {
	r, ok = m.Get(m.start + m.this - 1)
	return
}

func (m *cMatchState) This() (r rune, ok bool) {
	r, ok = m.Get(m.start + m.this)
	return
}

func (m *cMatchState) Next() (r rune, ok bool) {
	r, ok = m.Get(m.start + m.this + 1)
	return
}

func (m *cMatchState) Capture() {
	m.capture = true
}

func (m *cMatchState) Captured() bool {
	return m.capture
}

func (m *cMatchState) Consume(count int) (ok bool) {
	if delta := m.InputLen() - m.start - m.this - count; delta < 0 {
		// clamp on ceil
		count += delta
	}
	m.this += count
	ok = m.start+m.this+1 < len(m.s.input)
	return
}

func (m *cMatchState) Runes() []rune {
	if m.start+m.this <= len(m.s.input) {
		return m.s.input[m.start : m.start+m.this]
	}
	return nil
}

func (m *cMatchState) String() (text string) {
	if m.start+m.this <= len(m.s.input) {
		text = string(m.s.input[m.start : m.start+m.this])
	}
	return
}

func (m *cMatchState) Equal(other MatchState) bool {
	if o, ok := other.(*cMatchState); ok {
		if m.s == o.s {
			if m.start == o.start {
				if m.this == o.this {
					if m.capture == o.capture {
						if m.complete == o.complete {
							//if m.configs.Len() == o.configs.Len() {
							//}
							return true
						}
					}
				}
			}
		}
	}
	return false
}

func (m *cMatchState) Clone() MatchState {
	return m.CloneWith(0)
}

func (m *cMatchState) CloneWith(offset int) MatchState {
	if offset < 0 {
		offset = 0
	}
	cloned := newMatchState(m.s, m.start+offset)
	cloned.this = m.this
	cloned.capture = m.capture
	cloned.complete = m.complete
	cloned.scoped = m.scoped
	return cloned
}

func (m *cMatchState) Apply(other MatchState) {
	o, _ := other.(*cMatchState)
	o.this += m.this
	if m.capture {
		o.capture = m.capture
	}
	if m.complete {
		o.complete = m.complete
	}
}

func (m *cMatchState) Negated() bool {
	return m.scoped != nil && m.scoped.Negated
}

func (m *cMatchState) Scope(cfg *Config) (scope *Config) {

	if m.scoped == nil {
		m.scoped = newConfig()
	}
	if cfg != nil {
		scope = m.scoped.Merge(cfg)
		m.scoped = scope
		return
	}
	// not update, just copied
	scope = m.scoped.Merge(nil)
	return
}
