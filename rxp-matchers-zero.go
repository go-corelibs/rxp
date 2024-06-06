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

import (
	"unicode"
)

// Caret creates a Matcher equivalent to the regexp caret [^]
func Caret(flags ...string) Matcher {
	return MakeMatcher(func(scope Flags, reps Reps, input []rune, index int, sm SubMatches) (consumed int, captured bool, negated bool, proceed bool) {

		if scope.Multiline() {
			// start of input or start of line
			prev, ok := IndexGet(input, index-1)
			// if there is no previous character ~or~ the previous is a newline
			if proceed = !ok || prev == '\n'; scope.Negated() {
				// check negation before return
				proceed = !proceed
			}
			return
		}

		// check only for the start of input
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
			// look for: end of input or end of line
			r, ok := IndexGet(input, index)
			if proceed = !ok || r == '\n'; scope.Negated() {
				// check negation before return
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
func A(flags ...string) Matcher {
	_, cfg := ParseFlags(flags...)
	return func(scope Flags, reps Reps, input []rune, index int, sm SubMatches) (consumed int, captured bool, negated bool, proceed bool) {
		scope = scope.Merge(cfg)
		if proceed = index == 0; scope.Negated() {
			proceed = !proceed
		}
		return
	}
}

// B creates a Matcher equivalent to the regexp [\b]
func B(flags ...string) Matcher {
	_, cfg := ParseFlags(flags...)
	return func(scope Flags, reps Reps, input []rune, index int, sm SubMatches) (consumed int, captured bool, negated bool, proceed bool) {
		scope = scope.Merge(cfg)

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
			} else if next > 0 {
				// next is not null, boundary is to the right
				proceed = RuneIsWord(next)
			} else if prev > 0 {
				// prev is not null, boundary is to the left
				proceed = RuneIsWord(prev)
			}

		}

		if scope.Negated() {
			proceed = !proceed
		}

		return
	}
}

// Z is a Matcher equivalent to the regexp [\z]
func Z(flags ...string) Matcher {
	_, cfg := ParseFlags(flags...)
	return func(scope Flags, reps Reps, input []rune, index int, sm SubMatches) (consumed int, captured bool, negated bool, proceed bool) {
		scope = scope.Merge(cfg)
		if proceed = IndexInvalid(input, index); scope.Negated() {
			proceed = !proceed
		}
		return
	}
}

// BackRef is a Matcher equivalent to Perl backreferences where the gid
// argument is the match group to use
//
// BackRef will panic if the gid argument is less than one
func BackRef(gid int, flags ...string) Matcher {
	if gid < 1 {
		panic("BackRef requires a positive non-zero gid argument")
	}
	_, cfg := ParseFlags(flags...)
	return func(scope Flags, reps Reps, input []rune, index int, sm SubMatches) (consumed int, captured, negated, proceed bool) {
		scope = scope.Merge(cfg)
		captured = scope.Capture()

		if count := len(sm); count == 0 || gid > count { // gid > count is correct because gid is 1-indexed
			// id is out of range (non-zero index) or no matches present
			proceed = scope.Negated()
			return
		}

		smidx := gid - 1 // convert the gid to the submatch 0-indexed position

		group := sm[smidx]
		runes := input[group.Start():group.End()]
		groupLen := group.Len()

		if proceed = len(input) >= index+groupLen; !proceed {
			// forward is past EOF, OOB is not negated
			proceed = scope.Negated()
			return
		}

		for idx := 0; idx < groupLen; idx++ {
			forward := index + idx // forward position

			if scope.AnyCase() {
				proceed = unicode.ToLower(runes[idx]) == unicode.ToLower(input[forward])
			} else {
				proceed = runes[idx] == input[forward]
			}

			if scope.Negated() {
				proceed = !proceed
			}

			if !proceed {
				// early out
				return
			}

		}

		consumed = groupLen
		return
	}
}
