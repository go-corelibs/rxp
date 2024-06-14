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

func Benchmark_FindAllString_Regexp(b *testing.B) {
	_ = regexp.MustCompile(`(?ms)(\b[a-zA-Z0-9]+?['a-zA-Z0-9]*[a-zA-Z0-9]+\b|\b[a-zA-Z0-9]+\b)`).
		FindAllString(gTestDataRandomString, -1)
}

func Benchmark_FindAllString_Rxp(b *testing.B) {
	_ = Pattern{IsFieldWord("c")}.FindAllString(gTestDataRandomString, -1)
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
	// this is going to be slightly slower thant the regexp version due to the
	// usage of Not, a better version would be one with a custom matcher that
	// is equivalent to Not(W(), S(), "+")
	// this benchmark is not about sheer performance but about comparing the
	// convenience of the regexp pattern strings versus the Pattern convenience
	// methods
	_ = Pipeline{}.
		Transform(strings.ToLower).
		Literal(S("+"), " ").
		Literal(Text("'"), "_").
		Literal(Not(W(), S(), "+"), "").
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
			ReplaceAllString(gTestDataRandomString, Replace[string]{}.ToUpper())
	}
}
