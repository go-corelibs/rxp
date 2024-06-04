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

// RuneIsALNUM returns true for alphanumeric characters [a-zA-Z0-9]
func RuneIsALNUM(r rune) bool {
	return ('a' <= r && r <= 'z') ||
		('A' <= r && r <= 'Z') ||
		('0' <= r && r <= '9')
}

// RuneIsALPHA returns true for alphanumeric characters [a-zA-Z]
func RuneIsALPHA(r rune) bool {
	return ('a' <= r && r <= 'z') ||
		('A' <= r && r <= 'Z')
}

// RuneIsASCII returns true for valid ASCII characters [\x00-\x7F]
func RuneIsASCII(r rune) bool {
	return 0x00 <= r && r <= 0x7f
}

// RuneIsBLANK returns true for tab and space characters [\t ]
func RuneIsBLANK(r rune) bool {
	return r == '\t' || r == ' '
}

// RuneIsCNTRL returns true for control characters [\x00-\x1F\x7F]
func RuneIsCNTRL(r rune) bool {
	return r == 0x7f || (0x00 <= r && r <= 0x1f)
}

// RuneIsDIGIT returns true for number digits [0-9]
func RuneIsDIGIT(r rune) bool {
	return '0' <= r && r <= '9'
}

// RuneIsGRAPH returns true for graphical characters
// [a-zA-Z0-9!"$%&'()*+,\-./:;<=>?@[\\\]^_`{|}~]
//
// Note: upon the first use of RuneIsGRAPH, a lookup map is cached in a global
// variable and used for detecting the specific runes supported by the regexp
// [:graph:] class
func RuneIsGRAPH(r rune) bool {
	if gLookupGraphical == nil {
		gLookupGraphical = map[rune]struct{}{'!': {}, '"': {}, '#': {},
			'$': {}, '%': {}, '&': {}, '\'': {}, '(': {}, ')': {}, '*': {},
			'+': {}, ',': {}, '-': {}, '.': {}, '/': {}, ':': {}, ';': {},
			'<': {}, '=': {}, '>': {}, '?': {}, '@': {}, '[': {}, '\\': {},
			']': {}, '^': {}, '_': {}, '`': {}, '{': {}, '|': {}, '}': {},
			'~': {}}
	}
	_, present := gLookupGraphical[r]
	return present ||
		('a' <= r && r <= 'z') ||
		('A' <= r && r <= 'Z') ||
		('0' <= r && r <= '9')
}

// RuneIsLOWER returns true for lowercase alphabetic characters [a-z]
func RuneIsLOWER(r rune) bool {
	return 'a' <= r && r <= 'z'
}

// RuneIsPRINT returns true for space and RuneIsGRAPH characters [ [:graph:]]
//
// Note: uses RuneIsGRAPH
func RuneIsPRINT(r rune) bool {
	return r == ' ' || RuneIsGRAPH(r)
}

// RuneIsPUNCT returns true for punctuation characters [!-/:-@[-`{-~]
func RuneIsPUNCT(r rune) bool {
	return ('!' <= r && r <= '/') ||
		(':' <= r && r <= '@') ||
		('[' <= r && r <= '`') ||
		('{' <= r && r <= '~')
}

// RuneIsSPACE returns true for empty space characters [\t\n\v\f\r ]
func RuneIsSPACE(r rune) bool {
	return r == ' ' ||
		r == '\t' || r == '\n' || r == '\r' ||
		r == '\v' || r == '\f'
}

// RuneIsUPPER returns true for lowercase alphabetic characters [A-Z]
func RuneIsUPPER(r rune) bool {
	return 'A' <= r && r <= 'Z'
}

// RuneIsXDIGIT returns true for hexadecimal digits [z-fA-F0-9]
func RuneIsXDIGIT(r rune) bool {
	return ('a' <= r && r <= 'f') ||
		('A' <= r && r <= 'F') ||
		('0' <= r && r <= '9')
}
