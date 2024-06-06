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
	return MakeMatcher(func(scope Flags, reps Reps, input []rune, index int, sm SubMatches) (consumed int, captured bool, negated bool, proceed bool) {
		if scope.Capture() {
			captured = true
		}
		for _, matcher := range matchers {
			clone := scope.Clone()
			if cons, capt, nega, next := matcher(clone, reps, input, index, sm); next {
				return cons, capt || captured, nega, !scope.Negated()
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
	return MakeMatcher(func(scope Flags, reps Reps, input []rune, index int, sm SubMatches) (consumed int, captured bool, negated bool, proceed bool) {
		if IndexInvalid(input, index) {
			proceed = scope.Negated()
			return
		}
		negated = true

		// Not is an Or-like operation when there are multiple Matcher
		// instances, first one proceeds
		for _, matcher := range matchers {
			if consumed, _, _, proceed = matcher(scope.Clone(), reps, input, index, sm); proceed {
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
	return MakeMatcher(func(scope Flags, reps Reps, input []rune, index int, sm SubMatches) (consumed int, captured bool, negated bool, proceed bool) {

		if scope.Capture() {
			captured = true
		}

		for _, matcher := range matchers {
			if cons, _, _, next := matcher(scope, reps, input, index+consumed, sm); !next {
				consumed = 0
				return
			} else {
				consumed += cons
			}
		}

		// successful match of entire group
		proceed = !scope.Negated()

		return
	}, flags...)
}
