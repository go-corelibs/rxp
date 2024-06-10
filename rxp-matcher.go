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
//	| scoped   | possibly modified sub-match scope  |
//	| consumed | number of runes matched from index |
//	| proceed  | success, keep matching for more    |
type Matcher func(scope Flags, reps Reps, input *RuneBuffer, index int, sm [][2]int) (scoped Flags, consumed int, proceed bool)

// RuneMatcher is the signature for the basic character matching functions
// such as RuneIsWord
//
// Implementations are expected to operate using the least amount of CPU
// instructions possible
type RuneMatcher func(r rune) bool

// WrapMatcher wraps a RuneMatcher with MakeMatcher with support for negations
func WrapMatcher(matcher RuneMatcher, flags ...string) Matcher {
	return MakeMatcher(func(scope Flags, reps Reps, input *RuneBuffer, index int, sm [][2]int) (scoped Flags, consumed int, proceed bool) {
		scoped = scope
		if 0 <= index && index < input.len {
			r, size, _ := input.Get(index)
			if proceed = matcher(r); scoped.Negated() {
				proceed = !proceed
			}
			if proceed {
				consumed += size
			}
		}
		return
	}, flags...)

}

// MakeMatcher creates a rxp standard Matcher implementation wrapped
// around a given RuneMatcher
func MakeMatcher(match Matcher, flags ...string) Matcher {
	cfgReps, cfg := ParseFlags(flags...)
	return func(scope Flags, reps Reps, input *RuneBuffer, index int, sm [][2]int) (scoped Flags, consumed int, proceed bool) {
		scoped = scope | cfg
		if cfgReps != nil {
			reps = cfgReps
		}

		var scoping Flags
		var matched, completed bool
		var keep, count, queue int
		for this := 0; 0 <= index+this && index+this <= input.len; {
			// one past last is necessary for \z and $

			if scoping, keep, matched = match(scoped, reps, input, index+this, sm); scoping.Capture() {
				scoped |= CaptureFlag
			}

			if matched {
				scoped |= MatchedFlag
				count += 1

				if keep == 0 {
					this += 1
				} else if keep > 0 {
					this += keep
					queue += keep
				}

				if minHit, maxHit := reps.Satisfied(count); minHit {
					completed = true
					// met the min req
					if scoping.Less() {
						// preferring less, no need to continue
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

			// not a matched case
			break
		}

		// don't negate this, only negate actual RuneMatchers
		if proceed = completed; proceed {
			consumed += queue
		}

		return
	}
}
