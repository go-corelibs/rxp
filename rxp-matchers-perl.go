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

	return MakeRuneMatcher(func(scope *Config, m MatchState, start int, r rune) (consumed int, proceed bool) {

		// scan ahead without consuming runes

		// using the for-loop approach and unicode.ToLower reduces the actual
		// number of times unicode.ToLower is called compared to
		// strings.ToLower which has to scan the entire string each time

		if m.Has(start) {
			// with for looping
			// this version works

			input := m.Input()
			inputLen := len(input)

			for idx := 0; idx < needLen; idx++ {
				// match this rune with corresponding index value

				forward := start + m.Len() + idx // forward position
				if proceed = forward < inputLen; !proceed {
					// forward is past EOF, OOB is not negated?
					return
				}

				if scope.AnyCase {
					proceed = unicode.ToLower(runes[idx]) == unicode.ToLower(input[forward])
				} else {
					proceed = runes[idx] == input[forward]
				}

				if scope.Negated {
					proceed = !proceed
				}

				if !proceed {
					return
				}

			}

			consumed = needLen

		}

		return
	}, flags...)
}

// Dot creates a Matcher equivalent to the regexp dot (.)
func Dot(flags ...string) Matcher {
	return MakeRuneMatcher(func(scope *Config, m MatchState, start int, r rune) (consumed int, proceed bool) {
		if m.Has(start) {
			if proceed = r != '\n' || scope.DotNL; scope.Negated {
				proceed = !proceed
			}
			if proceed {
				consumed = 1
			}
		} else if scope.Negated {
			proceed = true
		}
		return
	}, flags...)
}

// D creates a Matcher equivalent to the regexp \d
func D(flags ...string) Matcher {
	return WrapFn(RuneIsDIGIT, flags...)
}

// S creates a Matcher equivalent to the regexp \s
func S(flags ...string) Matcher {
	return WrapFn(RuneIsSpace, flags...)
}

// W creates a Matcher equivalent to the regexp \w
func W(flags ...string) Matcher {
	return WrapFn(RuneIsWord, flags...)
}

// Alnum creates a Matcher equivalent to [:alnum:]
func Alnum(flags ...string) Matcher {
	return WrapFn(RuneIsALNUM, flags...)
}

// Alpha creates a Matcher equivalent to [:alpha:]
func Alpha(flags ...string) Matcher {
	return WrapFn(RuneIsALPHA, flags...)
}

// Ascii creates a Matcher equivalent to [:ascii:]
func Ascii(flags ...string) Matcher {
	return WrapFn(RuneIsASCII, flags...)
}

// Blank creates a Matcher equivalent to [:blank:]
func Blank(flags ...string) Matcher {
	return WrapFn(RuneIsBLANK, flags...)
}

// Cntrl creates a Matcher equivalent to [:cntrl:]
func Cntrl(flags ...string) Matcher {
	return WrapFn(RuneIsCNTRL, flags...)
}

// Digit creates a Matcher equivalent to [:digit:]
func Digit(flags ...string) Matcher {
	return WrapFn(RuneIsDIGIT, flags...)
}

// Graph creates a Matcher equivalent to [:graph:]
func Graph(flags ...string) Matcher {
	return WrapFn(RuneIsGRAPH, flags...)
}

// Lower creates a Matcher equivalent to [:lower:]
func Lower(flags ...string) Matcher {
	return WrapFn(RuneIsLOWER, flags...)
}

// Print creates a Matcher equivalent to [:print:]
func Print(flags ...string) Matcher {
	return WrapFn(RuneIsPRINT, flags...)
}

// Punct creates a Matcher equivalent to [:punct:]
func Punct(flags ...string) Matcher {
	return WrapFn(RuneIsPUNCT, flags...)
}

// Space creates a Matcher equivalent to [:space:]
func Space(flags ...string) Matcher {
	return WrapFn(RuneIsSpace, flags...)
}

// Upper creates a Matcher equivalent to [:upper:]
func Upper(flags ...string) Matcher {
	return WrapFn(RuneIsUPPER, flags...)
}

// Word creates a Matcher equivalent to [:word:]
func Word(flags ...string) Matcher {
	return WrapFn(RuneIsWord, flags...)
}

// Xdigit creates a Matcher equivalent to [:xdigit:]
func Xdigit(flags ...string) Matcher {
	return WrapFn(RuneIsXDIGIT, flags...)
}

// Class creates a Matcher equivalent to the regexp [:AsciiNames:],
// see the AsciiNames constants for the list of supported ASCII class
// names
//
// Class will panic if given an invalid class name
func Class(name AsciiNames, flags ...string) Matcher {
	if matcher, ok := LookupAsciiClass[name]; ok {
		return WrapFn(matcher, flags...)
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
	return MakeRuneMatcher(func(scope *Config, m MatchState, start int, r rune) (consumed int, proceed bool) {
		if m.Has(start) {
			if proceed = unicode.Is(table, r); scope.Negated {
				proceed = !proceed
			}
			if proceed {
				consumed += 1
			}
		} else if scope.Negated {
			proceed = true
		}
		return
	}, flags...)
}
