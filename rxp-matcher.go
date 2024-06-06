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

// Matcher is a single string matching function
//
//	| Argument | Description                        |
//	|----------|------------------------------------|
//	|  scope   | current Flags for this iteration   |
//	|  reps    | min and max repetition settings    |
//	|  input   | input rune slice (do not modify!)  |
//	|  index   | current input rune index to match  |
//
//	| Return   | Description                        |
//	|----------|------------------------------------|
//	| consumed | number of runes matched from index |
//	| captured | indicate this is a capture group   |
//	| negated  | indicate this group is negated     |
//	| proceed  | matched success, match for more    |
type Matcher func(scope Flags, reps Reps, input []rune, index int, sm SubMatches) (consumed int, captured, negated, proceed bool)

// RuneMatcher is the signature for the basic character matching functions
// such as RuneIsWord
//
// Implementations are expected to operate using the least amount of CPU
// instructions possible
type RuneMatcher func(r rune) bool

// WrapMatcher wraps a RuneMatcher with MakeMatcher with support for negations
func WrapMatcher(matcher RuneMatcher, flags ...string) Matcher {
	return MakeMatcher(func(scope Flags, reps Reps, input []rune, index int, sm SubMatches) (consumed int, captured, negated, proceed bool) {
		if IndexReady(input, index) {
			if proceed = matcher(input[index]); scope.Negated() {
				proceed = !proceed
			}
			if proceed {
				consumed += 1
				captured = scope.Capture()
			}
		}
		return
	}, flags...)

}

// MakeMatcher creates a rxp standard Matcher implementation wrapped
// around a given RuneMatcher
func MakeMatcher(match Matcher, flags ...string) Matcher {
	cfgReps, cfg := ParseFlags(flags...)
	return func(scope Flags, reps Reps, input []rune, index int, sm SubMatches) (consumed int, captured, negated, proceed bool) {
		scope |= cfg
		if cfgReps != nil {
			reps = cfgReps
		}

		if scope.Capture() {
			captured = true
		}

		inputLen := len(input)
		var next, capt, completed bool
		var keep, count, queue int
		for this := 0; index+this <= inputLen; {
			idx := index + this
			if IndexValid(input, idx) {
				// one past last is necessary for \z and $

				keep, capt, _, next = match(scope, reps, input, idx, sm)
				if capt {
					captured = true
				}

				if next {
					count += 1

					if keep == 0 {
						this += 1
					} else if keep > 0 && idx <= inputLen {
						// guard on progressing beyond EOF
						this += keep
						queue += keep
					}

					if minHit, maxHit := reps.Satisfied(count); minHit {
						// met the min req
						completed = true
						if scope.Less() {
							// don't need more than the min?
							break
						} else if maxHit {
							// there is a limit and this is it
							break
						}
					}

					// passed less/min/max checks, see if there's any more
					continue
				}

				if minOk, _ := reps.Satisfied(count); minOk {
					completed = true
				}

			}

			// did not pass match check
			break
		}

		// don't negate this, only negate actual RuneMatchers
		if proceed = completed; proceed {
			consumed += queue
		}

		return
	}
}
