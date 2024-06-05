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

var _ Fragment = (*cFragment)(nil)

var (
	spFragment = newSyncPool(10, func() *cFragment {
		return new(cFragment)
	})
)

// Fragment is the Pattern.ScanString return value and can be either a Result
// of a captured Matcher or unmatched plain text. The Fragment slice returned
// by Pattern.ScanString can be used to rebuild the original input
type Fragment interface {
	// Match returns the associated Match for this Fragment, if this
	// Fragment is for matched text
	Match() (result Match, ok bool)
	// Index returns the start and end indexes of this Fragment
	Index() []int
	// Runes returns the input runes for this Index range
	Runes() []rune
	// String is a convenience wrapper around Runes
	String() string
	// Bytes is a convenience wrapper around Runes
	Bytes() []byte

	Recycle()

	private()
}

type cFragment struct {
	s          *cPatternState
	start, end int
	match      bool
	result     *cMatch
	recycled   bool
}

func newFragment(s *cPatternState, start, end int, result *cMatch) Fragment {
	frag := spFragment.Get()
	//if !frag.recycled {
	//	// pool created a new one, let's add ten more so next batch aren't new
	//	spFragment.Seed(-1)
	//}
	frag.s = s
	frag.start = start
	frag.end = end
	frag.match = result != nil && len(result.matches) > 0
	frag.result = result
	frag.recycled = false
	return frag
}

func (f *cFragment) Recycle() {
	if !f.recycled {
		f.s = nil
		f.result = nil
		f.recycled = true
		spFragment.Put(f)
		//TODO: recycle f.result?
	}
}

func (f *cFragment) private() {}

func (f *cFragment) Match() (result Match, ok bool) {
	if ok = f.match && f.result != nil; ok {
		result = f.result
	}
	return
}

func (f *cFragment) Index() []int {
	if size := len(f.s.input); f.start < f.end && f.end <= size {
		return []int{f.start, f.end}
	}
	return nil
}

func (f *cFragment) Runes() []rune {
	if size := len(f.s.input); f.start < f.end && f.end <= size {
		return f.s.input[f.start:f.end]
	}
	return nil
}

func (f *cFragment) String() string {
	return string(f.Runes())
}

func (f *cFragment) Bytes() (data []byte) {
	return []byte(f.String())
}
