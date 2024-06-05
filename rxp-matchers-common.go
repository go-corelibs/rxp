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
	reps, cfg := ParseFlags(flags...)
	return func(m MatchState) (next, keep bool) {
		if m.Invalid() {
			return
		}
		_, scope := m.Scope(reps, cfg)
		//lh, scope := m.Scope(reps, cfg)
		//minimum, maximum := lh[0], lh[1]

		var proceed bool
		if this, ok := m.This(); ok {
			if proceed = RuneIsALNUM(this); proceed {
				// first rune matched the first range [a-zA-Z]
				m.Consume(1)

				var consume int
				input := m.Input()
				total := len(input)
				start := m.Index() + m.Len()

				// scan for second range runes
				for idx := start; idx < total; idx += 1 {
					r := input[idx]

					if r == '\'' || RuneIsALNUM(r) {
						consume += 1
						continue
					}

					// stop scanning, not within second range anymore
					break
				}

				if consume > 0 && start+consume < total {
					for consume-1 >= 0 {
						if RuneIsALNUM(input[start+consume-1]) {
							break
						}
						consume -= 1
					}
				}

				m.Consume(consume)
				if keep = scope.Capture(); keep {
					m.Capture()
				}
			}
		}

		if scope.Negated() {
			proceed = !proceed
		}

		next = proceed
		return
	}
}

// FieldKey creates a Matcher equivalent to:
//
//	(?:\b[a-zA-Z][-_a-zA-Z0-9]+?[a-zA-Z0-9]\b)
func FieldKey(flags ...string) Matcher {
	reps, cfg := ParseFlags(flags...)
	return func(m MatchState) (next, keep bool) {
		if m.Invalid() {
			return
		}
		_, scope := m.Scope(reps, cfg)

		if this, ok := m.This(); ok {
			if next = RuneIsALPHA(this); next {
				// matched first range [a-zA-Z]
				m.Consume(1)
				if keep = scope.Capture(); keep {
					m.Capture()
				}

				var consume int
				input := m.Input()
				total := len(input)
				start := m.Index() + m.Len()

				// scan for second range
				for idx := start; idx < total; idx += 1 {
					r := input[idx]

					if RuneIsDashUnderALNUM(r) {
						consume += 1
						continue
					}

					// not a valid rune for this Matcher
					break
				}

				// if the last rune is [-_], rewind to last [^-_]
				if consume > 0 && start+consume < total {
					for last := input[start+consume]; consume > 0 && (last == '-' || last == '_'); consume -= 1 {
					}
				}
				if consume > 0 {
					m.Consume(consume)
				}

			}
		}

		if scope.Negated() {
			next = !next
		}

		return
	}
}

// Keyword is intended for Go-Enjin parsing of simple search keywords from
// user input and creates a Matcher equivalent to:
//
//	(?:\b[-+]?[a-zA-Z][-_a-zA-Z0-9]+?[a-zA-Z0-9]\b)
func Keyword(flags ...string) Matcher {
	reps, cfg := ParseFlags(flags...)
	return func(m MatchState) (next, keep bool) {
		if m.Invalid() {
			return
		}
		_, scope := m.Scope(reps, cfg)

		if this, ok := m.This(); ok {
			var plusOrMinus rune
			if RuneIsPlusMinus(this) {
				plusOrMinus = this
				if this, ok = m.Get(m.Index() + m.Len() + 1); !ok {
					return false, false
				}
			}
			if next = RuneIsALPHA(this); next {

				if plusOrMinus > 0 {
					// consume the keyword modifier rune
					m.Consume(1)
				}

				// matched first range [a-zA-Z]
				m.Consume(1)
				if keep = scope.Capture(); keep {
					m.Capture()
				}

				var consume int
				input := m.Input()
				total := len(input)
				start := m.Index() + m.Len()

				// scan for second range
				for idx := start; idx < total; idx += 1 {
					r := input[idx]

					if RuneIsDashUnderALNUM(r) {
						consume += 1
						continue
					}

					// not a valid rune for this Matcher
					break
				}

				// if the last rune is [-_], rewind to last [^-_]
				if consume > 0 {
					if start+consume < total {
						for last := input[start+consume]; consume > 0 && RuneIsDashUnder(last); consume -= 1 {
						}
					}
					m.Consume(consume)
				}

			}
		}

		if scope.Negated() {
			next = !next
		}

		return
	}
}
