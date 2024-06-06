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

// FieldWord creates a Matcher equivalent to:
//
//	(?:\b[a-zA-Z0-9]+?['a-zA-Z0-9]*[a-zA-Z0-9]+\b|\b[a-zA-Z0-9]+\b)
func FieldWord(flags ...string) Matcher {
	_, cfg := ParseFlags(flags...)
	return func(scope Flags, _ Reps, input []rune, index int, matches SubMatches) (consumed int, captured bool, negated bool, proceed bool) {
		if IndexInvalid(input, index) {
			return
		}
		scope = scope.Merge(cfg)

		var this rune
		if this, proceed = IndexGet(input, index); !proceed {
			if scope.Negated() {
				proceed = true
			}
			return
		}

		if proceed = RuneIsALNUM(this); proceed {
			// first rune matched the first range [a-zA-Z]
			consumed += 1

			total := len(input)

			// scan for second range runes
			for idx := index + consumed; idx < total; idx += 1 {
				r := input[idx]

				if r == '\'' || RuneIsALNUM(r) {
					consumed += 1
					continue
				}

				// stop scanning, not within second range anymore
				break
			}

			if consumed > 1 && index+consumed-1 < total {
				for consumed > 1 && !RuneIsALNUM(input[index+consumed-1]) {
					consumed -= 1
				}
			}

			if scope.Capture() {
				captured = true
			}
		}

		if scope.Negated() {
			proceed = !proceed
		}

		return
	}
}

// FieldKey creates a Matcher equivalent to:
//
//	(?:\b[a-zA-Z][-_a-zA-Z0-9]+?[a-zA-Z0-9]\b)
func FieldKey(flags ...string) Matcher {
	_, cfg := ParseFlags(flags...)
	return func(scope Flags, reps Reps, input []rune, index int, sm SubMatches) (consumed int, captured bool, negated bool, proceed bool) {
		if IndexInvalid(input, index) {
			return
		}
		scope = scope.Merge(cfg)

		var this rune
		if this, proceed = IndexGet(input, index); !proceed {
			if scope.Negated() {
				proceed = true
			}
			return
		}

		if proceed = RuneIsALPHA(this); proceed {
			// matched first range [a-zA-Z]
			consumed += 1
			if scope.Capture() {
				captured = true
			}

			total := len(input)

			// scan for second range
			for idx := index + consumed; idx < total; idx += 1 {
				r := input[idx]

				if RuneIsDashUnderALNUM(r) {
					consumed += 1
					continue
				}

				// not a valid rune for this Matcher
				break
			}

			// if the last rune is [-_], rewind to last [^-_]
			if consumed > 1 && index+consumed-1 < total {
				for consumed > 1 && RuneIsDashUnder(input[index+consumed-1]) {
					consumed -= 1
				}
			}

		}

		if scope.Negated() {
			proceed = !proceed
		}

		return
	}
}

// Keyword is intended for Go-Enjin parsing of simple search keywords from
// user input and creates a Matcher equivalent to:
//
//	(?:\b[-+]?[a-zA-Z][-_a-zA-Z0-9]+?[a-zA-Z0-9]\b)
func Keyword(flags ...string) Matcher {
	_, cfg := ParseFlags(flags...)
	return func(scope Flags, reps Reps, input []rune, index int, sm SubMatches) (consumed int, captured bool, negated bool, proceed bool) {
		if IndexInvalid(input, index) {
			return
		}
		scope = scope.Merge(cfg)

		var this rune
		if this, proceed = IndexGet(input, index); !proceed {
			if scope.Negated() {
				proceed = true
			}
			return
		}

		var plusOrMinus rune
		if RuneIsPlusMinus(this) {
			plusOrMinus = this
			if this, proceed = IndexGet(input, index+1); !proceed {
				return
			}
		}

		if proceed = RuneIsALPHA(this); proceed {

			if plusOrMinus > 0 {
				// consume the previously detected keyword modifier rune
				consumed += 1
			}

			// this matched first range [a-zA-Z]
			consumed += 1
			if scope.Capture() {
				captured = true
			}

			total := len(input)

			// scan for second range
			for idx := index + consumed; idx < total; idx += 1 {
				r := input[idx]

				if RuneIsDashUnderALNUM(r) {
					consumed += 1
					continue
				}

				// not a valid rune for this Matcher
				break
			}

			// if the last rune is [-_], rewind to last [^-_]
			if consumed > 1 && index+consumed-1 < total {
				for consumed > 1 && RuneIsDashUnder(input[index+consumed-1]) {
					consumed -= 1
				}
			}

		}

		if scope.Negated() {
			proceed = !proceed
		}

		return
	}
}
