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

type AsciiNames string

const (
	ALNUM  AsciiNames = "alnum"
	ALPHA  AsciiNames = "alpha"
	ASCII  AsciiNames = "ascii"
	BLANK  AsciiNames = "blank"
	CNTRL  AsciiNames = "cntrl"
	DIGIT  AsciiNames = "digit"
	GRAPH  AsciiNames = "graph"
	LOWER  AsciiNames = "lower"
	PRINT  AsciiNames = "print"
	PUNCT  AsciiNames = "punct"
	SPACE  AsciiNames = "space"
	UPPER  AsciiNames = "upper"
	WORD   AsciiNames = "word"
	XDIGIT AsciiNames = "xdigit"
)

var LookupAsciiClass = map[AsciiNames]func(r rune) bool{
	ALNUM:  RuneIsALNUM,
	ALPHA:  RuneIsALPHA,
	ASCII:  RuneIsASCII,
	BLANK:  RuneIsBLANK,
	CNTRL:  RuneIsCNTRL,
	DIGIT:  RuneIsDIGIT,
	GRAPH:  RuneIsGRAPH,
	LOWER:  RuneIsLOWER,
	PRINT:  RuneIsPRINT,
	PUNCT:  RuneIsPUNCT,
	SPACE:  RuneIsSpace,
	UPPER:  RuneIsUPPER,
	WORD:   RuneIsWord,
	XDIGIT: RuneIsXDIGIT,
}

// cached on first use of RuneIsGRAPH
var gLookupGraphical map[rune]struct{}
