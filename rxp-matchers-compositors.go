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

// Or processes the list of Matcher instances, in the order they were given,
// and stops at the first one that returns a true next
//
// Or accepts Pattern, Matcher and string types and will panic on all others
func Or(options ...interface{}) Matcher {
	matchers, flags, _ := ParseOptions(options...)
	return MakeRuneMatcher(func(scope Flags, m MatchState, start int, r rune) (consumed int, proceed bool) {
		for _, matcher := range matchers {
			clone := m.Clone()
			if next, keep := matcher(clone); next {
				if keep {
					// apply cloned MatchState to actual
					clone.Recycle()
					return clone.Len(), true
				}
				// or consumes at least one rune, even if not keeping? nope?
				clone.Recycle()
				return 0, true
			}
			clone.Recycle()
		}
		return
	}, flags...)
}

// Not processes the given Matcher and inverts the next response without
// consuming any runes (zero-width Matcher)
//
// Not accepts Pattern, Matcher and string types and will panic on all others
func Not(options ...interface{}) Matcher {
	matchers, flags, _ := ParseOptions(options...)
	return MakeRuneMatcher(func(scope Flags, m MatchState, start int, r rune) (consumed int, proceed bool) {
		scope.SetNegated() // adding a ^ to the Not Matcher would cancel out the whole point of an explicit Not

		clone := m.Clone()
		defer clone.Recycle()
		if len(matchers) == 1 {
			proceed, _ = matchers[0](clone)
		} else {
			// Not is an Or-like operation when there are multiple Matcher
			// instances, first on proceeds
			for _, matcher := range matchers {
				iter := m.Clone()
				if proceed, _ = matcher(clone); proceed {
					iter.Apply(clone)
					iter.Recycle()
					break
				}
				iter.Recycle()
			}
		}

		// invert the positives
		if proceed = !proceed; proceed {
			if consumed = clone.Len(); consumed == 0 {
				consumed = 1
			}
		}
		return
	}, flags...)
}

// Group processes the list of Matcher instances, in the order they were
// given, and stops at the first one that does not match, discarding any
// consumed runes. If all Matcher calls succeed, all consumed runes are
// applied
func Group(options ...interface{}) Matcher {
	matchers, flags, _ := ParseOptions(options...)
	return MakeRuneMatcher(func(scope Flags, m MatchState, start int, r rune) (consumed int, proceed bool) {
		clone := m.Clone() // accumulate matched runes
		defer clone.Recycle()
		for _, matcher := range matchers {
			iter := m.CloneWith(clone.Len(), m.Reps())
			if next, _ := matcher(iter); !next {
				// entire group did not in fact match
				iter.Recycle()
				return 0, false
			}
			iter.Apply(clone)
			iter.Recycle()
		}
		proceed = true
		consumed = clone.Len()
		return
	}, flags...)
}
