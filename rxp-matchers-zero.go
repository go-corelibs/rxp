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

// Caret creates a Matcher equivalent to the regexp caret [^]
func Caret(options ...string) Matcher {
	return MakeRuneMatcher(func(scope Flags, m MatchState, start int, r rune) (consumed int, proceed bool) {

		if scope.Multiline() {
			// start of input or start of line
			prev, ok := m.Prev()
			// if there is no previous character ~or~ the previous is a newline
			if proceed = !ok || prev == '\n'; scope.Negated() {
				proceed = !proceed
			}
			return
		}

		// start of input
		if proceed = start == 0; scope.Negated() {
			proceed = !proceed
		}
		return
	}, options...)
}

// Dollar creates a Matcher equivalent to the regexp [$]
func Dollar(options ...string) Matcher {
	return MakeRuneMatcher(func(scope Flags, m MatchState, start int, r rune) (consumed int, proceed bool) {

		if scope.Multiline() {
			// end of input or end of line
			// if there is no this character ~or~ this is a newline
			if proceed = r == '\n'; scope.Negated() {
				proceed = !proceed
			}
			return
		}

		// end of input
		if proceed = start >= m.InputLen(); scope.Negated() {
			proceed = !proceed
		}
		return
	}, options...)
}

// A creates a Matcher equivalent to the regexp [\A]
func A() Matcher {
	return func(m MatchState) (next, keep bool) {
		if next = m.Index() == 0; m.Flags().Negated() {
			next = !next
		}
		return
	}
}

// B creates a Matcher equivalent to the regexp [\b]
func B() Matcher {
	return func(m MatchState) (proceed, keep bool) {

		start := m.Index()
		next, _ := m.This()
		prev, _ := m.Prev()

		if proceed = start == 0 && RuneIsWord(next); proceed {
			// at the start of input, boundary is prev if next is word

		} else if proceed = start == m.InputLen() && RuneIsWord(prev); proceed {
			// at the end of input, boundary is next if prev is word

		} else if proceed = RuneIsWord(prev) && !RuneIsWord(next); proceed {
			// in the middle of input, boundary is next if prev is word and next is not word

		} else if proceed = !RuneIsWord(prev) && RuneIsWord(next); proceed {
			// in the middle of input, boundary is prev if next is word and prev is not word

		}

		if m.Flags().Negated() {
			proceed = !proceed
		}

		return
	}
}

// Z is a Matcher equivalent to the regexp [\z]
func Z() Matcher {
	return func(m MatchState) (next, keep bool) {
		if next = m.Invalid(); m.Flags().Negated() {
			next = !next
		}
		return
	}
}
