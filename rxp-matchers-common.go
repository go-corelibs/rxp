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
	return func(scope Flags, reps Reps, input *RuneBuffer, index int, sm [][2]int) (scoped Flags, consumed int, proceed bool) {
		scoped = scope | cfg
		if 0 > index || index >= input.len {
			proceed = scoped.Negated()
			return
		}

		this, size, _ := input.Get(index) // this will never fail due to previous IndexInvalid check

		if proceed = RuneIsALNUM(this); proceed {
			// first rune matched the first range [a-zA-Z]
			consumed += size

			total := input.Len()

			// scan for second range runes
			for idx := index + consumed; idx < total; {
				if r, rs, rok := input.Get(idx); rok {

					if RuneIsALNUM(r) {
						consumed += rs
						idx += rs
						continue
					}

					if r == '\'' {
						if next, _, ok := input.Get(idx + rs); ok {
							// found the next rune
							if RuneIsALNUM(next) {
								// accept this rune
								consumed += rs
								idx += rs
								continue
							}
						}
						// end of input reached, or next is not alnum
						// do not consume this \' rune
					}

				}

				// stop scanning, not within second range anymore
				break
			}

		}

		if scoped.Negated() {
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
	return func(scope Flags, reps Reps, input *RuneBuffer, index int, sm [][2]int) (scoped Flags, consumed int, proceed bool) {
		scoped = scope | cfg
		if 0 > index || index >= input.len {
			proceed = scoped.Negated()
			return
		}

		this, size, _ := input.Get(index)

		if proceed = RuneIsALPHA(this); proceed {
			// matched first range [a-zA-Z]
			consumed += size

			total := input.Len()

			var scanForward func(start int) (size int, ok bool)
			scanForward = func(start int) (size int, ok bool) {
				if nxt, sz, present := input.Get(start); present {
					// found the next rune
					if RuneIsALNUM(nxt) {
						// accept this rune
						return sz, true
					} else if RuneIsDashUnder(nxt) {
						if size, ok = scanForward(start + sz); ok {
							size += sz
							return size, true
						}
					}
				}
				return 0, false
			}

			// scan for second range
			for idx := index + consumed; idx < total; {
				if r, rs, rok := input.Get(idx); rok {

					if RuneIsALNUM(r) {
						consumed += rs
						idx += rs
						continue
					}

					if RuneIsDashUnder(r) {

						if sz, ok := scanForward(idx + rs); ok {
							consumed += rs + sz
							idx += rs + sz
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

		if scoped.Negated() {
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
	return func(scope Flags, reps Reps, input *RuneBuffer, index int, sm [][2]int) (scoped Flags, consumed int, proceed bool) {
		scoped = scope | cfg
		if 0 > index || index >= input.len {
			proceed = scoped.Negated()
			return
		}

		this, size, _ := input.Get(index)

		var plusOrMinus rune
		if RuneIsPlusMinus(this) {
			plusOrMinus = this // is one byte, size overwritten by next
			if this, size, proceed = input.Get(index + 1); !proceed {
				return
			}
		}

		if proceed = RuneIsALPHA(this); proceed {

			if plusOrMinus > 0 {
				// consume the previously detected keyword modifier rune
				consumed += 1 // just one byte
			}

			// this matched first range [a-zA-Z]
			consumed += size

			total := input.Len()

			// scan for second range
			for idx := index + consumed; idx < total; {
				if r, rs, rok := input.Get(idx); rok {

					if RuneIsALNUM(r) {
						consumed += rs
						idx += rs
						continue
					}

					if RuneIsDashUnder(r) {
						if next, _, ok := input.Get(idx + rs); ok {
							// found the next rune
							if RuneIsALNUM(next) {
								// accept this rune
								consumed += rs
								idx += rs
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

		if scoped.Negated() {
			proceed = !proceed
		}

		return
	}
}
