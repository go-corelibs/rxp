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
	runes := []rune(text)
	needLen := len(runes)

	return MakeMatcher(func(scope Flags, reps Reps, input []rune, index int, sm SubMatches) (consumed int, captured bool, negated bool, proceed bool) {

		// scan ahead without consuming runes

		// using the for-loop approach and unicode.ToLower reduces the actual
		// number of times unicode.ToLower is called compared to
		// strings.ToLower which has to scan the entire string each time

		if !IndexReady(input, index) {
			proceed = scope.Negated()
			return
		}

		for idx := 0; idx < needLen; idx++ {
			forward := index + idx // forward position
			if proceed = IndexReady(input, forward); !proceed {
				// forward is past EOF, OOB is not negated
				return
			}

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

		consumed = needLen

		return
	}, flags...)
}

// Dot creates a Matcher equivalent to the regexp dot (.)
func Dot(flags ...string) Matcher {
	return MakeMatcher(func(scope Flags, reps Reps, input []rune, index int, sm SubMatches) (consumed int, captured bool, negated bool, proceed bool) {

		if r, ok := IndexGet(input, index); ok {
			proceed = r != '\n' || scope.DotNL()
			if proceed {
				consumed = 1
			}
		} else {
			proceed = scope.Negated()
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
	panic(fmt.Errorf("invalid AsciiNames: %q", name))
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
	return MakeMatcher(func(scope Flags, reps Reps, input []rune, index int, sm SubMatches) (consumed int, captured bool, negated bool, proceed bool) {

		if r, ok := IndexGet(input, index); ok {
			if proceed = unicode.Is(table, r); scope.Negated() {
				proceed = !proceed
			}
			if proceed {
				consumed += 1
			}
		} else if scope.Negated() {
			proceed = true
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
	charLen := len(chars)

	for idx, this := range chars {
		if idx == 0 && this == '-' {
			// first dash is literal dash
			runes = append(runes, this)
			continue
		}
		if idx+2 < charLen {
			// range requires a total of three, the low and high runes and
			// a dash separator
			if chars[idx+1] == '-' {
				// next is a dash
				ranges = append(ranges, []rune{chars[idx+1], chars[idx+2]})
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
