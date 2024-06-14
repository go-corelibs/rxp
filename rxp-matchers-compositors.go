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

// Or processes the list of Matcher instances, in the order they were given,
// and stops at the first one that returns a true next
//
// Or accepts Pattern, Matcher and string types and will panic on all others
func Or(options ...interface{}) Matcher {
	matchers, flags, _ := ParseOptions(options...)
	return MakeMatcher(func(scope Flags, reps Reps, input *InputReader, index int, sm [][2]int) (scoped Flags, consumed int, proceed bool) {
		scoped = scope
		if scoped&NegatedFlag == NegatedFlag {
			// can't be any of
			clone := scoped.Unset(NegatedFlag)
			for _, matcher := range matchers {
				if _, _, proceed = matcher(clone, reps, input, index, sm); proceed {
					break
				}
			}
			if proceed = !proceed; proceed {
				_, size, _ := input.Get(index)
				consumed += size
			}
			return
		}
		// stop at the first
		for _, matcher := range matchers {
			clone := scoped
			if scoping, cons, next := matcher(clone, reps, input, index, sm); next {
				return scoping, cons, next
			}
		}
		return
	}, flags...)
}

// Not processes all the matchers given, in the order they were given, stopping
// at the first one that succeeds and inverts the proceed return value
//
// Not is equivalent to a negated character class in traditional regular
// expressions, for example:
//
//	[^xyza-f]
//
// could be implemented as any of the following:
//
//	// slower due to four different matchers being present
//	Not(Text("x"),Text("y"),Text("z"),R("a-f"))
//	// better but still has two matchers
//	Not(R("xyza-f"))
//	// no significant difference from the previous
//	Or(R("xyza-f"), "^") //< negation (^) flag
//	// simplified to just one matcher present
//	R("xyza-f", "^") //< negation (^) flag
//
// here's the interesting bit about rxp though, if speed is really the goal,
// then the following would capture single characters matching [^xyz-af] with
// significant performance over MakeMatcher based matchers (use Pattern.Add to
// include the custom Matcher)
//
//	func(scope Flags, reps Reps, input *InputReader, index int, sm [][2]int) (scoped Flags, consumed int, proceed bool) {
//	    scoped = scope
//	    if r, size, ok := input.Get(index); ok {
//	        // test for [xyza-f]
//	        proceed = (r >= 'x' && r <= 'z') || (r >= 'a' && r <= 'f')
//	        // and invert the result
//	        proceed = !proceed
//	        if proceed { // true means the negation is a match
//	        // MatchedFlag is required, CaptureFlag optionally if needed
//	        scoped |= MatchedFlag | CaptureFlag
//	        // consume this rune's size if a capture group is needed
//	        // using size instead of just 1 will allow support for
//	        // accurate []rune, []byte and string processing
//	        consumed += size
//	    }
//	    return
//	}
func Not(options ...interface{}) Matcher {
	matchers, flags, _ := ParseOptions(options...)
	return MakeMatcher(func(scope Flags, reps Reps, input *InputReader, index int, sm [][2]int) (scoped Flags, consumed int, proceed bool) {
		scoped = scope
		if 0 > index || index >= input.len {
			return
		}

		// Not is an Or-like operation when there are multiple Matcher
		// instances, first one proceeds
		for _, matcher := range matchers {
			if _, consumed, proceed = matcher(scoped, reps, input, index, sm); proceed {
				break
			}
		}

		// always negate the matcher results
		if proceed = !proceed; proceed {
			if consumed == 0 {
				// always consume at least one?
				consumed = 1
			}
		}

		return
	}, flags...)
}

// Group processes the list of Matcher instances, in the order they were given,
// and stops at the first one that does not match, discarding any consumed
// runes. If all Matcher calls succeed, all consumed runes are accepted together
// as this group (sub-sub-matches are not a thing in rxp)
func Group(options ...interface{}) Matcher {
	matchers, flags, _ := ParseOptions(options...)
	return MakeMatcher(func(scope Flags, reps Reps, input *InputReader, index int, sm [][2]int) (scoped Flags, consumed int, proceed bool) {
		scoped = scope

		for _, matcher := range matchers {
			if _, cons, next := matcher(scoped, reps, input, index+consumed, sm); !next {
				consumed = 0
				return
			} else {
				consumed += cons
			}
		}

		// successful match of entire group
		proceed = !scoped.Negated()
		return
	}, flags...)
}
