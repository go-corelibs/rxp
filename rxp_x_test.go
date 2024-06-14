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
	"strings"
)

//func Benchmark_Text_Nope(b *testing.B) {
//	_ = Pattern{xTextNoLoop("lorem", "i")}.FindAllString(gTestDataRandomString, -1)
//}
//
//func Benchmark_Text_Loop(b *testing.B) {
//	_ = Pattern{Text("lorem", "i")}.FindAllString(gTestDataRandomString, -1)
//}

// xTextNoLoop is a benchmark for comparing the actual Text matcher performance
// without looping, this proved to be slower than the accepted Text version
// which ranges over characters instead of retrieving a string from the
// underlying runes.RuneReader
//
// this code is left here for others to play with and experiment with optimizing
// rxp further
func xTextNoLoop(text string, flags ...string) Matcher {
	runes := []rune(text)
	needLen := len(runes)

	return MakeMatcher(func(scope Flags, reps Reps, input *InputReader, index int, sm [][2]int) (scoped Flags, consumed int, proceed bool) {

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
