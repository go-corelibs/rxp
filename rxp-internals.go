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
	sync "github.com/go-corelibs/x-sync"
)

// spStringBuilder from go-corelibs/x-sync does not seem to impact performance
// to the same extent as is the case of sync.Append
var spStringBuilder = sync.NewStringBuilderPool(1)

func pushString(slice []string, data string) []string {
	m := len(slice)
	n := m + 1
	if n > cap(slice) {
		grown := make([]string, ((m+1)*2)+n) // double the existing space
		copy(grown, slice)                   // transfer to new slice
		slice = grown                        // grown becomes slice
	}
	slice = slice[0:n] // truncate in case too many present
	slice[n-1] = data  // populate with data
	return slice
}

func pushStrings(slice [][]string, data ...[]string) [][]string {
	m := len(slice)
	n := m + len(data)
	if n > cap(slice) {
		grown := make([][]string, ((m+1)*2)+n) // double the existing space
		copy(grown, slice)                     // transfer to new slice
		slice = grown                          // grown becomes slice
	}
	slice = slice[0:n]     // truncate in case too many present
	copy(slice[m:n], data) // populate with data
	return slice
}

func pushRunes(slice []rune, data ...rune) []rune {
	m := len(slice)
	n := m + len(data)
	if n > cap(slice) {
		grown := make([]rune, ((m+1)*2)+n) // double the existing space
		copy(grown, slice)                 // transfer to new slice
		slice = grown                      // grown becomes slice
	}
	slice = slice[0:n]     // truncate in case too many present
	copy(slice[m:n], data) // populate with data
	return slice
}

func pushMatch(slice [][][2]int, data [][2]int) [][][2]int {
	m := len(slice)
	n := m + 1
	if n > cap(slice) {
		grown := make([][][2]int, ((m+1)*2)+n) // double the existing space
		copy(grown, slice)                     // transfer to new slice
		slice = grown                          // grown becomes slice
	}
	slice = slice[0:n] // truncate in case too many present
	slice[n-1] = data
	return slice
}

func pushSMSlice(slice [][2]int, start, end int) [][2]int {
	m := len(slice)
	n := m + 1
	if n > cap(slice) {
		grown := make([][2]int, ((m+1)*2)+n) // double the existing space
		copy(grown, slice)                   // transfer to new slice
		slice = grown                        // grown becomes slice
	}
	slice = slice[0:n]              // truncate in case too many present
	slice[n-1] = [2]int{start, end} // populate with data
	return slice
}

func pushSegments(slice Segments, data ...Segment) Segments {
	m := len(slice)
	n := m + len(data)
	if n > cap(slice) {
		grown := make(Segments, ((m+1)*2)+n) // double the existing space
		copy(grown, slice)                   // transfer to new slice
		slice = grown                        // grown becomes slice
	}
	slice = slice[0:n]     // truncate in case too many present
	copy(slice[m:n], data) // populate with data
	return slice
}

func mapKeys[K comparable, V interface{}](m map[K]V) []K {
	var slice []K
	for k := range m {
		// no need for appendSlice because mapKeys is only used by NamedClass
		// when given a gid < 0, pattern-build-time issue
		slice = append(slice, k)
	}
	return slice
}
