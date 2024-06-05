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

// Matcher is a single string matching function
//
//	next  indicates that this Matcher function does not match the current rune
//	      and to progress the Pattern to the next Matcher instance
//	keep  indicates that this Matcher function was in fact satisfied, even
//	      though stop may also be true
type Matcher func(m MatchState) (next, keep bool)

// RuneMatcher is the signature for the MakeRuneMatcher matching function
//
// The MatchState provided to RuneMatcher functions is a MatchState.Clone and
// any changes to that MatchState by the RuneMatcher are discarded. In order
// to have the RuneMatcher consume any runes, return the consumed number of
// runes and a true proceed
type RuneMatcher func(cfg Flags, m MatchState, pos int, r rune) (consumed int, proceed bool)

// RuneMatchFn is the signature for the basic character matching functions
// such as RuneIsWord
//
// Implementations are expected to operate using the least amount of CPU
// instructions possible
type RuneMatchFn func(r rune) bool

// WrapFn wraps a RuneMatchFn with MakeRuneMatcher
func WrapFn(matcher RuneMatchFn, flags ...string) Matcher {
	return MakeRuneMatcher(func(scope Flags, m MatchState, start int, r rune) (consumed int, proceed bool) {
		if m.Has(start) {
			// within bounds
			if proceed = matcher(r); scope.Negated() {
				proceed = !proceed
			}
			if proceed {
				consumed += 1
			}
		}
		return
	}, flags...)

}

// MakeRuneMatcher creates a rxp standard Matcher implementation wrapped
// around a given RuneMatcher
func MakeRuneMatcher(match RuneMatcher, flags ...string) Matcher {
	reps, cfg := ParseFlags(flags...)
	// TODO: investigate optimizing MakeRuneMatcher further, likely has something to do with MatchState handling
	return func(m MatchState) (next, keep bool) {
		lh, scope := m.Scope(reps, cfg)
		minimum, maximum := lh[0], lh[1]

		if scope.Capture() {
			m.Capture()
		}

		var completed bool
		var count, queue int
		for start, this := m.Index()+m.Len(), 0; start+this <= m.InputLen(); {
			idx := start + this
			if idx <= m.InputLen() {
				// one past last is necessary for \z and $

				r, _ := m.Get(idx)
				clone := m.CloneWith(this, lh)
				consumed, proceed := match(scope, clone, idx, r)
				clone.Recycle()

				if proceed {
					count += 1

					if consumed == 0 {
						this += 1
					} else if consumed > 0 && idx <= m.InputLen() {
						// guard on progressing beyond EOF
						this += consumed
						queue += consumed
					}

					if count >= minimum {
						// met the min req
						completed = true
						if scope.Less() {
							// don't need more than the min?
							break
						} else if maximum > 0 && count >= maximum {
							// there is a limit and this is it
							break
						}
					}

					// passed less/min/max checks, see if there's any more
					continue
				}

				if count >= minimum {
					completed = true
				}

			}

			// did not pass match check
			break
		}

		// don't negate this, only negate actual RuneMatchers
		//if next = completed && scope.IsValidCount(count); next {
		if next = completed; /*&& scope.IsValidCount(count)*/ next {
			m.Consume(queue)
			keep = queue > 0
		}

		return
	}
}
