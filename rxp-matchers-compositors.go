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
	return MakeMatcher(func(scope Flags, reps Reps, input *RuneBuffer, index int, sm [][2]int) (scoped Flags, consumed int, proceed bool) {
		scoped = scope
		for _, matcher := range matchers {
			clone := scoped
			if scoping, cons, next := matcher(clone, reps, input, index, sm); next {
				return scoping, cons, !scoped.Negated()
			}
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
	return MakeMatcher(func(scope Flags, reps Reps, input *RuneBuffer, index int, sm [][2]int) (scoped Flags, consumed int, proceed bool) {
		scoped = scope
		if 0 > index || index >= input.len {
			return
		}

		// Not is an Or-like operation when there are multiple Matcher
		// instances, first one proceeds
		for _, matcher := range matchers {
			if _, consumed, proceed = matcher(scoped, reps, input, index, sm); proceed {
				break
			}
		}

		// negate the matcher results
		if proceed = !proceed; proceed {
			if consumed == 0 {
				// always consume at least one?
				consumed = 1
			}
		}

		return
	}, flags...)
}

// Group processes the list of Matcher instances, in the order they were
// given, and stops at the first one that does not match, discarding any
// consumed runes. If all Matcher calls succeed, all consumed runes are
// accepted together as this group (sub-sub-matches are not a thing)
func Group(options ...interface{}) Matcher {
	matchers, flags, _ := ParseOptions(options...)
	return MakeMatcher(func(scope Flags, reps Reps, input *RuneBuffer, index int, sm [][2]int) (scoped Flags, consumed int, proceed bool) {
		scoped = scope

		for _, matcher := range matchers {
			if _, cons, next := matcher(scoped, reps, input, index+consumed, sm); !next {
				consumed = 0
				return
			} else {
				consumed += cons
			}
		}

		// successful match of entire group
		proceed = !scoped.Negated()

		return
	}, flags...)
}
