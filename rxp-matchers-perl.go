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
	"fmt"
	"unicode"
)

// Text creates a Matcher for the plain text given
func Text(text string, flags ...string) Matcher {
	content := []rune(text)

	return MakeMatcher(func(scope Flags, reps Reps, input *InputReader, index int, sm [][2]int) (scoped Flags, consumed int, proceed bool) {
		scoped = scope

		if scoped&NegatedFlag == NegatedFlag {
			// the meaning of "proceed" is inverted in a negation context

			if proceed = 0 > index || index >= input.len; proceed {
				// out-of-bounds is a negative negated making it positive with
				// zero consumed
				return
			} else if proceed = len(content) > input.len-index; proceed {
				// content matching is impossible due to insufficient input, and
				// negated consumes this index
				_, consumed, _ = input.Get(index)
				return
			}

			// track this index size to increment []byte and string readers correctly
			_, size, _ := input.Get(index)

			var matched bool
			for idx := 0; idx < len(content); idx++ {
				forward := index + idx

				r, _, _ := input.Get(forward)

				if scoped.AnyCase() {
					matched = unicode.ToLower(content[idx]) == unicode.ToLower(r)
				} else {
					matched = content[idx] == r
				}

				if !matched {
					break
				}

			}

			if !matched {
				// not matched is true proceed, consuming the size of this index
				proceed = true
				consumed = size
			}

			return
		}

		if 0 > index || index >= input.len || len(content) > input.len-index {
			// negative index, or index oob, or insufficient input
			return
		}

		// using the for-loop approach and unicode.ToLower reduces the actual
		// number of times unicode.ToLower is called compared to
		// strings.ToLower which has to scan the entire string each time
		// see the commented benchmarks in rxp_x_test.go

		var size int
		for idx := 0; idx < len(content); idx++ {
			forward := index + idx // forward position

			r, rs, _ := input.Get(forward)

			if scoped.AnyCase() {
				proceed = unicode.ToLower(content[idx]) == unicode.ToLower(r)
			} else {
				proceed = content[idx] == r
			}

			if !proceed {
				// definitely not a match, early out
				return
			}

			size += rs
		}

		consumed = size
		return
	}, flags...)
}

// Dot creates a Matcher equivalent to the regexp dot (.)
func Dot(flags ...string) Matcher {
	return MakeMatcher(func(scope Flags, reps Reps, input *InputReader, index int, sm [][2]int) (scoped Flags, consumed int, proceed bool) {
		scoped = scope
		if r, rs, ok := input.Get(index); ok {
			proceed = r != '\n' || scoped&DotNewlineFlag == DotNewlineFlag
			if proceed {
				consumed = rs
			}
		} else {
			proceed = scoped.Negated()
		}

		return
	}, flags...)
}

// D creates a Matcher equivalent to the regexp \d
func D(flags ...string) Matcher {
	return WrapMatcher(RuneIsDIGIT, flags...)
}

// S creates a Matcher equivalent to the regexp \s
func S(flags ...string) Matcher {
	return WrapMatcher(RuneIsSpace, flags...)
}

// W creates a Matcher equivalent to the regexp \w
func W(flags ...string) Matcher {
	return WrapMatcher(RuneIsWord, flags...)
}

// Alnum creates a Matcher equivalent to [:alnum:]
func Alnum(flags ...string) Matcher {
	return WrapMatcher(RuneIsALNUM, flags...)
}

// Alpha creates a Matcher equivalent to [:alpha:]
func Alpha(flags ...string) Matcher {
	return WrapMatcher(RuneIsALPHA, flags...)
}

// Ascii creates a Matcher equivalent to [:ascii:]
func Ascii(flags ...string) Matcher {
	return WrapMatcher(RuneIsASCII, flags...)
}

// Blank creates a Matcher equivalent to [:blank:]
func Blank(flags ...string) Matcher {
	return WrapMatcher(RuneIsBLANK, flags...)
}

// Cntrl creates a Matcher equivalent to [:cntrl:]
func Cntrl(flags ...string) Matcher {
	return WrapMatcher(RuneIsCNTRL, flags...)
}

// Digit creates a Matcher equivalent to [:digit:]
func Digit(flags ...string) Matcher {
	return WrapMatcher(RuneIsDIGIT, flags...)
}

// Graph creates a Matcher equivalent to [:graph:]
func Graph(flags ...string) Matcher {
	return WrapMatcher(RuneIsGRAPH, flags...)
}

// Lower creates a Matcher equivalent to [:lower:]
func Lower(flags ...string) Matcher {
	return WrapMatcher(RuneIsLOWER, flags...)
}

// Print creates a Matcher equivalent to [:print:]
func Print(flags ...string) Matcher {
	return WrapMatcher(RuneIsPRINT, flags...)
}

// Punct creates a Matcher equivalent to [:punct:]
func Punct(flags ...string) Matcher {
	return WrapMatcher(RuneIsPUNCT, flags...)
}

// Space creates a Matcher equivalent to [:space:]
func Space(flags ...string) Matcher {
	return WrapMatcher(RuneIsSPACE, flags...)
}

// Upper creates a Matcher equivalent to [:upper:]
func Upper(flags ...string) Matcher {
	return WrapMatcher(RuneIsUPPER, flags...)
}

// Word creates a Matcher equivalent to [:word:]
func Word(flags ...string) Matcher {
	return WrapMatcher(RuneIsWord, flags...)
}

// Xdigit creates a Matcher equivalent to [:xdigit:]
func Xdigit(flags ...string) Matcher {
	return WrapMatcher(RuneIsXDIGIT, flags...)
}

// NamedClass creates a Matcher equivalent to the regexp [:AsciiNames:],
// see the AsciiNames constants for the list of supported ASCII class
// names
//
// NamedClass will panic if given an invalid class name
func NamedClass(name AsciiNames, flags ...string) Matcher {
	if matcher, ok := LookupAsciiClass[name]; ok {
		return WrapMatcher(matcher, flags...)
	}
	panic(fmt.Errorf("invalid ASCII name: %q, valid names are: %q", name, mapKeys(LookupAsciiClass)))
}

// IsUnicodeRange creates a Matcher equivalent to the regexp \pN where N
// is a unicode character class, passed to IsUnicodeRange as a
// unicode.RangeTable instance
//
// For example, creating a Matcher for a single braille character:
//
//	IsUnicodeRange(unicode.Braille)
func IsUnicodeRange(table *unicode.RangeTable, flags ...string) Matcher {
	_ = unicode.Is(table, 'a') // compile-time test for panic cases
	return MakeMatcher(func(scope Flags, reps Reps, input *InputReader, index int, sm [][2]int) (scoped Flags, consumed int, proceed bool) {
		scoped = scope
		if r, rs, ok := input.Get(index); ok {
			if proceed = unicode.Is(table, r); scoped&NegatedFlag == NegatedFlag {
				proceed = !proceed
			}
			if proceed {
				consumed += rs
			}
		}

		return
	}, flags...)
}

// R creates a Matcher equivalent to regexp character class ranges such as:
// [xyza-f] where x, y and z are individual runes to accept and a-f is the
// inclusive range of letters from lowercase a to lowercase f to accept
//
// Note: do not include the [] brackets unless the intent is to actually accept
// those characters
func R(characters string, flags ...string) Matcher {
	var runes []rune
	var ranges [][]rune
	chars := []rune(characters)
	charsLen := len(chars)

	// this is all happening at compile time so no need to use a pushRune or
	// other append() optimizations

	// parse the character class rune pattern [xyza-f]
	for idx := 0; idx < charsLen; idx++ {
		this := chars[idx]
		if idx == 0 && this == '-' {
			// first dash is literal dash
			runes = append(runes, this)
			continue
		}
		if idx+2 < charsLen {
			// range requires a total of three, the low and high runes and
			// a dash separator
			if chars[idx+1] == '-' {
				// next is a dash
				if chars[idx] > chars[idx+2] {
					// the low rune is greater than the high rune
					// allow these mistakes? hmm...
					ranges = append(ranges, []rune{chars[idx+2], chars[idx]})
				} else {
					ranges = append(ranges, []rune{chars[idx], chars[idx+2]})
				}
				idx += 2 // one for the dash and one for the other rune
				continue
			}
		}
		runes = append(runes, this)
	}

	hasRunes, hasRanges := len(runes) > 0, len(ranges) > 0
	return WrapMatcher(func(r rune) bool {

		if hasRanges {
			for _, check := range ranges {
				if check[0] <= r && r <= check[1] {
					return true
				}
			}
		}

		if hasRunes {
			for _, check := range runes {
				if check == r {
					return true
				}
			}
		}

		return false
	}, flags...)
}
