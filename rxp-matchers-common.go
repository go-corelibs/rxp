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
	return func(scope Flags, _ Reps, input *RuneBuffer, index int, matches SubMatches) (consumed int, captured bool, negated bool, proceed bool) {
		scope |= cfg
		if input.Invalid(index) {
			proceed = scope.Negated()
			return
		}

		this, _ := input.Get(index) // this will never fail due to previous IndexInvalid check

		if proceed = RuneIsALNUM(this); proceed {
			// first rune matched the first range [a-zA-Z]
			consumed += 1

			total := input.Len()

			// scan for second range runes
			for idx := index + consumed; idx < total; {
				r, _ := input.Get(idx)

				if RuneIsALNUM(r) {
					consumed += 1
					idx += 1
					continue
				}

				if r == '\'' {
					if next, ok := input.Get(idx + 1); ok {
						// found the next rune
						if RuneIsALNUM(next) {
							// accept this rune
							consumed += 1
							idx += 1
							continue
						}
					}
					// end of input reached, or next is not alnum
					// do not consume this \' rune
				}

				// stop scanning, not within second range anymore
				break
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
	return func(scope Flags, reps Reps, input *RuneBuffer, index int, sm SubMatches) (consumed int, captured bool, negated bool, proceed bool) {
		scope |= cfg
		if input.Invalid(index) {
			proceed = scope.Negated()
			return
		}

		this, _ := input.Get(index)

		if proceed = RuneIsALPHA(this); proceed {
			// matched first range [a-zA-Z]
			consumed += 1
			captured = scope.Capture()

			total := input.Len()

			var scanForward func(start int) (size int, ok bool)
			scanForward = func(start int) (size int, ok bool) {
				if nxt, present := input.Get(start); present {
					// found the next rune
					if RuneIsALNUM(nxt) {
						// accept this rune
						return 1, true
					} else if RuneIsDashUnder(nxt) {
						if size, ok = scanForward(start + 1); ok {
							size += 1
							return size, true
						}
					}
				}
				return 0, false
			}

			// scan for second range
			for idx := index + consumed; idx < total; {
				if r, rok := input.Get(idx); rok {

					if RuneIsALNUM(r) {
						consumed += 1
						idx += 1
						continue
					}

					if RuneIsDashUnder(r) {

						if sz, ok := scanForward(idx + 1); ok {
							consumed += 1 + sz
							idx += 1 + sz
							continue
						}

						// end of input reached, or next is not alnum
						// do not consume this - or _ rune
					}

				}

				// not a valid rune for this Matcher
				break
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
	return func(scope Flags, reps Reps, input *RuneBuffer, index int, sm SubMatches) (consumed int, captured bool, negated bool, proceed bool) {
		scope |= cfg
		if input.Invalid(index) {
			proceed = scope.Negated()
			return
		}

		this, _ := input.Get(index)

		var plusOrMinus rune
		if RuneIsPlusMinus(this) {
			plusOrMinus = this
			if this, proceed = input.Get(index + 1); !proceed {
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

			total := input.Len()

			// scan for second range
			for idx := index + consumed; idx < total; {
				if r, rok := input.Get(idx); rok {

					if RuneIsALNUM(r) {
						consumed += 1
						idx += 1
						continue
					}

					if RuneIsDashUnder(r) {
						if next, ok := input.Get(idx + 1); ok {
							// found the next rune
							if RuneIsALNUM(next) {
								// accept this rune
								consumed += 1
								idx += 1
								continue
							}
						}
						// end of input reached, or next is not alnum
						// do not consume this - or _ rune
					}

				}

				// not a valid rune for this Matcher
				break
			}

		}

		if scope.Negated() {
			proceed = !proceed
		}

		return
	}
}
