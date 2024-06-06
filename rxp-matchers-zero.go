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
func Caret(flags ...string) Matcher {
	return MakeMatcher(func(scope Flags, reps Reps, input []rune, index int, sm SubMatches) (consumed int, captured bool, negated bool, proceed bool) {

		if scope.Multiline() {
			// start of input or start of line
			prev, ok := IndexGet(input, index-1)
			// if there is no previous character ~or~ the previous is a newline
			if proceed = !ok || prev == '\n'; scope.Negated() {
				proceed = !proceed
			}
			return
		}

		// start of input
		if proceed = index == 0; scope.Negated() {
			proceed = !proceed
		}
		return
	}, flags...)
}

// Dollar creates a Matcher equivalent to the regexp [$]
func Dollar(flags ...string) Matcher {
	return MakeMatcher(func(scope Flags, reps Reps, input []rune, index int, sm SubMatches) (consumed int, captured bool, negated bool, proceed bool) {

		if scope.Multiline() {
			// end of input or end of line
			// if there is no this character ~or~ this is a newline
			if proceed = IndexEnd(input, index); proceed {
				return
			}
			r, ok := IndexGet(input, index)
			if proceed = ok && r == '\n'; scope.Negated() {
				proceed = !proceed
			}
			// matched on the newline
			return
		}

		// matching on end of input
		if proceed = IndexEnd(input, index); scope.Negated() {
			proceed = !proceed
		}
		return
	}, flags...)
}

// A creates a Matcher equivalent to the regexp [\A]
func A() Matcher {
	return func(scope Flags, reps Reps, input []rune, index int, sm SubMatches) (consumed int, captured bool, negated bool, proceed bool) {
		if proceed = index == 0; scope.Negated() {
			proceed = !proceed
		}
		return
	}
}

// B creates a Matcher equivalent to the regexp [\b]
func B() Matcher {
	return func(scope Flags, reps Reps, input []rune, index int, sm SubMatches) (consumed int, captured bool, negated bool, proceed bool) {

		this, _ := IndexGet(input, index)
		next, _ := IndexGet(input, index+1)
		prev, _ := IndexGet(input, index-1)

		if index == 0 {

			// at start of input, boundary is to the left if this is a word
			proceed = RuneIsWord(this)

		} else if index >= len(input) {

			// at the end of input, boundary is to the right if this is a word
			proceed = RuneIsWord(prev)

		} else {

			// somewhere in the middle of the string

			if RuneIsWord(this) {
				// this is a word, boundary is to the left
				proceed = !RuneIsWord(prev)
			} else {
				// this is not a word, boundary is to the right
				proceed = RuneIsWord(next)
			}

		}

		if scope.Negated() {
			proceed = !proceed
		}

		return
	}
}

// Z is a Matcher equivalent to the regexp [\z]
func Z() Matcher {
	return func(scope Flags, reps Reps, input []rune, index int, sm SubMatches) (consumed int, captured bool, negated bool, proceed bool) {
		if proceed = IndexInvalid(input, index); scope.Negated() {
			proceed = !proceed
		}
		return
	}
}
