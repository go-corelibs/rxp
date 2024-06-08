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
	_ "embed"
	"regexp"
	"strings"
	"testing"
	"unicode"
)

//go:embed testdata/random.txt
var gTestDataRandomString string

var gTestDataRandomRunes []rune

func init() {
	gTestDataRandomRunes = []rune(gTestDataRandomString)
}

func xTextNoLoop(text string, flags ...string) Matcher {
	runes := []rune(text)
	needLen := len(runes)

	return MakeMatcher(func(scope Flags, reps Reps, input *RuneBuffer, index int, sm [][2]int) (consumed int, captured bool, negated bool, proceed bool) {

		// scan ahead without consuming runes
		// without any for looping
		// this is marginally slower than the for-loop

		if input.Ready(index) {

			inputLen := input.Len()

			end := index + needLen
			if proceed = end <= inputLen; !proceed {
				return
			}

			if maybe := input.String(index, end-index); scope.AnyCase() {
				proceed = strings.ToLower(text) == strings.ToLower(maybe)
			} else {
				proceed = text == maybe
			}

			if scope.Negated() {
				if proceed = !proceed; proceed {
					// negations only move the needle by one
					consumed = 1
				}
			} else if proceed {
				// positives move the needle however much is needed
				consumed = needLen
			}

			return // bypass working code for wip
		}

		return
	}, flags...)
}

func Benchmark_Text_Nope(b *testing.B) {
	_ = Pattern{xTextNoLoop("lorem", "i")}.FindAllString(gTestDataRandomString, -1)
}

func Benchmark_Text_Loop(b *testing.B) {
	_ = Pattern{Text("lorem", "i")}.FindAllString(gTestDataRandomString, -1)
}

func Benchmark_ScanString_Regexp(b *testing.B) {
	// There is no ScanString within the regexp package, this benchmark test
	// is emulating what would be done (without using rxp).

	// ScanString returns a special slice of Segment types that include both the
	// matched and unmatched input, and indicates which is which. This cannot
	// be done with FindAllStringSubmatch because there's no way to know where
	// one fragment begins and another ends

	// FieldWord equivalent regexp pattern:
	pattern := `(?ms)(\b[a-zA-Z0-9]+?['a-zA-Z0-9]*[a-zA-Z0-9]+\b|\b[a-zA-Z0-9]+\b)`
	m := regexp.MustCompile(pattern).FindAllStringIndex(gTestDataRandomString, -1)
	var results []string
	var last int
	for _, mm := range m {
		if last < mm[0] {
			results = append(results, gTestDataRandomString[last:mm[0]])
			last = mm[0]
		}
		results = append(results, gTestDataRandomString[mm[0]:mm[1]])
	}
}

func Benchmark_ScanString_Rxp(b *testing.B) {
	_ = Pattern{FieldWord("c")}.ScanRunes(gTestDataRandomRunes).Strings()
}

func Benchmark_FindAllString_Regexp(b *testing.B) {
	_ = regexp.MustCompile(`(?ms)(\b[a-zA-Z0-9]+?['a-zA-Z0-9]*[a-zA-Z0-9]+\b|\b[a-zA-Z0-9]+\b)`).
		FindAllString(gTestDataRandomString, -1)
}

func Benchmark_FindAllString_Rxp(b *testing.B) {
	_ = Pattern{FieldWord("c")}.FindAllString(gTestDataRandomString, -1)
}

func Benchmark_Pipeline_Combo_Regexp(b *testing.B) {
	o := strings.TrimSpace(gTestDataRandomString)
	o = regexp.MustCompile(`(?ms)([^a-zA-Z0-9]+)`).ReplaceAllString(o, "_") // [^a-zA-Z0-9]
	o = regexp.MustCompile(`(?ms)(\s+)`).ReplaceAllString(o, "_")           // \s+
	_ = strings.ToLower(o)
}

func Benchmark_Pipeline_Combo_Rxp(b *testing.B) {
	_ = Pipeline{
		{Transform: strings.TrimSpace},
		{ // so nice
			Transform: func(input string) string {
				var buf strings.Builder
				var spaces bool
				for _, r := range input {
					if RuneIsSPACE(r) {
						spaces = true
						continue
					} else if spaces {
						buf.WriteRune('_')
						spaces = false
					}
					if RuneIsALNUM(r) {
						buf.WriteRune(r)
						continue
					} else {
						// not a space, nor alnum
						buf.WriteRune('_')
					}
				}
				return buf.String()
			},
		},
		{Transform: strings.ToLower},
	}.Process(gTestDataRandomString)
}

func Benchmark_Pipeline_Readme_Regexp(b *testing.B) {
	// using regexp:
	output := strings.ToLower(gTestDataRandomString)
	output = regexp.MustCompile(`\s+`).ReplaceAllString(output, " ")
	output = regexp.MustCompile(`[']`).ReplaceAllString(output, "_")
	_ = regexp.MustCompile(`[^\w\s]+`).ReplaceAllString(output, "")
}

func Benchmark_Pipeline_Readme_Rxp(b *testing.B) {
	_ = Pipeline{}.
		Transform(strings.ToLower).
		ReplaceText(S("+"), " ").
		ReplaceText(Text("'"), "_").
		ReplaceText(Not(W(), S(), "+"), "").
		Process(gTestDataRandomString)
}

func Benchmark_Replace_ToUpper_Regexp(b *testing.B) {
	for i := 1; i < 10; i += 1 {
		_ = regexp.MustCompile(`(?ms)(\pL+)`).ReplaceAllStringFunc(gTestDataRandomString, func(s string) string {
			return strings.ToUpper(s)
		})
	}
}

func Benchmark_Replace_ToUpper_Rxp(b *testing.B) {
	for i := 1; i < 10; i += 1 {
		_ = Pattern{}.RangeTable(unicode.L, "+", "m", "s", "c").
			ReplaceAllString(gTestDataRandomString, Replace{}.ToUpper())
	}
}

