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

func isFieldWordScanForward(input *InputReader, start int) (size int, ok bool) {
	if nxt, sz, present := input.Get(start); present {
		// found the next rune
		if RuneIsALNUM(nxt) {
			// accept this rune
			return sz, true
		} else if nxt == '-' || nxt == '_' || nxt == '\'' {
			if size, ok = isFieldWordScanForward(input, start+sz); ok {
				size += sz
				return size, true
			}
		}
	}
	return 0, false
}

// IsFieldWord creates a Matcher equivalent to:
//
//	(?:\b[a-zA-Z0-9]+?[-_a-zA-Z0-9']*[a-zA-Z0-9]+\b|\b[a-zA-Z0-9]+\b)
func IsFieldWord(flags ...string) Matcher {
	_, cfg := ParseFlags(flags...)
	return func(scope Flags, reps Reps, input *InputReader, index int, sm [][2]int) (scoped Flags, consumed int, proceed bool) {
		scoped = scope | cfg
		if 0 > index || index >= input.len {
			proceed = scoped.Negated()
			return
		}

		this, size, _ := input.Get(index) // this will never fail due to previous IndexInvalid check

		if proceed = RuneIsALNUM(this); proceed {
			// first rune matched the first range [a-zA-Z]
			consumed += size

			// scan for second range runes
			for idx := index + consumed; idx < input.len; {
				if r, rs, rok := input.Get(idx); rok {
					if RuneIsALNUM(r) {
						consumed += rs
						idx += rs
						continue
					} else if r == '-' || r == '_' || r == '\'' {
						if sz, ok := isFieldWordScanForward(input, idx+rs); ok {
							consumed += rs + sz
							idx += rs + sz
							continue
						}
					}
				}
				// stop scanning, not within second range anymore
				break
			}

		}

		if scoped.Negated() {
			proceed = !proceed
		}
		if proceed {
			scoped |= MatchedFlag
			if consumed == 0 {
				consumed += 1
			}
		}

		return
	}
}

func isFieldKeyScanForward(input *InputReader, start int) (size int, ok bool) {
	if nxt, sz, present := input.Get(start); present {
		// found the next rune
		if RuneIsALNUM(nxt) {
			// accept this rune
			return sz, true
		} else if nxt == '-' || nxt == '_' {
			if size, ok = isFieldKeyScanForward(input, start+sz); ok {
				size += sz
				return size, true
			}
		}
	}
	return 0, false
}

// IsFieldKey creates a Matcher equivalent to:
//
//	(?:\b[a-zA-Z][-_a-zA-Z0-9]+?[a-zA-Z0-9]\b)
//
// IsFieldKey is intended to validate CSS and HTML attribute key names such as
// "data-thing" or "some_value"
func IsFieldKey(flags ...string) Matcher {
	_, cfg := ParseFlags(flags...)
	return func(scope Flags, reps Reps, input *InputReader, index int, sm [][2]int) (scoped Flags, consumed int, proceed bool) {
		scoped = scope | cfg
		if 0 > index || index >= input.len {
			proceed = scoped.Negated()
			return
		}

		this, size, _ := input.Get(index)

		if proceed = RuneIsALPHA(this); proceed {
			// matched first range [a-zA-Z]
			consumed += size

			// scan for second range
			for idx := index + consumed; idx < input.len; {
				if r, rs, rok := input.Get(idx); rok {
					if RuneIsALNUM(r) {
						consumed += rs
						idx += rs
						continue
					} else if r == '-' || r == '_' {
						if sz, ok := isFieldKeyScanForward(input, idx+rs); ok {
							consumed += rs + sz
							idx += rs + sz
							continue
						}
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

func isKeywordScanForward(input *InputReader, start int) (size int, ok bool) {
	if nxt, sz, present := input.Get(start); present {
		// found the next rune
		if RuneIsALNUM(nxt) {
			// accept this rune
			return sz, true
		} else if nxt == '-' || nxt == '_' || nxt == '\'' {
			if size, ok = isKeywordScanForward(input, start+sz); ok {
				size += sz
				return size, true
			}
		}
	}
	return 0, false
}

// IsKeyword is intended for Go-Enjin parsing of simple search keywords from
// user input and creates a Matcher equivalent to:
//
//	(?:\b[-+]?[a-zA-Z][-_a-zA-Z0-9']+?[a-zA-Z0-9]\b)
func IsKeyword(flags ...string) Matcher {
	_, cfg := ParseFlags(flags...)
	return func(scope Flags, reps Reps, input *InputReader, index int, sm [][2]int) (scoped Flags, consumed int, proceed bool) {
		scoped = scope | cfg
		if 0 > index || index >= input.len {
			proceed = scoped.Negated()
			return
		}

		this, size, _ := input.Get(index)

		var plusOrMinus rune
		if this == '+' || this == '-' {
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

			// scan for second range
			for idx := index + consumed; idx < input.len; {
				if r, rs, rok := input.Get(idx); rok {

					if RuneIsALNUM(r) {
						consumed += rs
						idx += rs
						continue
					} else if r == '-' || r == '_' || r == '\'' {
						if sz, ok := isKeywordScanForward(input, idx+rs); ok {
							consumed += rs + sz
							idx += rs + sz
							continue
						}
					}

				}

				// not a valid rune for this Matcher
				break
			}

			proceed = !scoped.Negated()
		}

		return
	}
}

// IsHash10 creates a Matcher equivalent to:
//
//	(?:[a-fA-F0-9]{10})
func IsHash10(flags ...string) Matcher {
	_, cfg := ParseFlags(flags...)
	return func(scope Flags, reps Reps, input *InputReader, index int, sm [][2]int) (scoped Flags, consumed int, proceed bool) {
		scoped = scope | cfg

		// exactly 10 characters required
		if 0 > index || index >= input.len || input.len-index < 10 {
			proceed = scoped.Negated()
			return
		}

		for idx := 0; idx < 10; idx++ {
			r, _, _ := input.Get(index + consumed + idx)
			if !RuneIsXDIGIT(r) {
				proceed = scoped.Negated()
				return
			}
		}

		if proceed = !scoped.Negated(); proceed {
			consumed = 10
		}

		return
	}
}

// IsAtLeastSixDigits creates a Matcher equivalent to:
//
//	(?:\A[0-9]{6,}\z)
func IsAtLeastSixDigits(flags ...string) Matcher {
	_, cfg := ParseFlags(flags...)
	return func(scope Flags, reps Reps, input *InputReader, index int, sm [][2]int) (scoped Flags, consumed int, proceed bool) {
		scoped = scope | cfg

		// exactly 10 characters required
		if 0 > index || index >= input.len || input.len-index < 6 {
			proceed = scoped.Negated()
			return
		}

		for idx := 0; idx < 6; idx++ {
			r, _, _ := input.Get(index + consumed + idx)
			if !RuneIsDIGIT(r) {
				proceed = scoped.Negated()
				return
			}
		}

		if proceed = !scoped.Negated(); proceed {
			consumed = 6
		}

		return
	}
}

// IsUUID creates a Matcher equivalent to:
//
//	(?:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})
func IsUUID(flags ...string) Matcher {
	return MakeMatcher(func(scope Flags, reps Reps, input *InputReader, index int, sm [][2]int) (scoped Flags, consumed int, proceed bool) {
		scoped = scope

		// exactly 36 (8+4+4+4+12+4) characters long required
		if 0 > index || index >= input.len || input.len-index < 36 {
			proceed = scoped.Negated()
			return
		}

		for i, need := range []int{8, 4, 4, 4, 12} {
			for idx := 0; idx < need; idx++ {
				r, _, _ := input.Get(index + consumed + idx)
				if !RuneIsXDIGIT(r) {
					proceed = scoped.Negated()
					return
				}
			}
			if i < 4 {
				// need dash too
				r, _, _ := input.Get(index + consumed + need)
				if r != '-' {
					proceed = scoped.Negated()
					return
				}
				consumed += need + 1
			}
		}

		if proceed = !scoped.Negated(); proceed {
			consumed = 36
		}

		return
	}, flags...)
}
