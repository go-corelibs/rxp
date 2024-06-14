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
	"bytes"

	sync "github.com/go-corelibs/x-sync"
)

// spStringBuilder from go-corelibs/x-sync does not seem to impact performance
// to the same extent as is the case of sync.Append
var spStringBuilder = sync.NewStringBuilderPool(1)
var spBytesBuffer = sync.NewPool[*bytes.Buffer](1, func() *bytes.Buffer {
	return new(bytes.Buffer)
}, func(v *bytes.Buffer) *bytes.Buffer {
	// getter
	v.Reset()
	return v
}, func(v *bytes.Buffer) *bytes.Buffer {
	// setter
	if v.Len() < 64000 {
		return v
	}
	return nil
})

func pushByte(slice [][]byte, data []byte) [][]byte {
	have := len(slice)
	size := len(data)
	need := have + size
	if need > cap(slice) {
		grown := make([][]byte, have+need) // double the existing space
		copy(grown, slice)                 // transfer to new slice
		slice = grown                      // grown becomes slice
	}
	slice = slice[0:need]                  // truncate in case too many present
	copy(slice[have:need], [][]byte{data}) // populate with data
	return slice
}

func pushBytes(slice [][][]byte, data ...[][]byte) [][][]byte {
	have := len(slice)
	size := len(data)
	need := have + size
	if need > cap(slice) {
		grown := make([][][]byte, have+need) // double the existing space
		copy(grown, slice)                   // transfer to new slice
		slice = grown                        // grown becomes slice
	}
	slice = slice[0:need]        // truncate in case too many present
	copy(slice[have:need], data) // populate with data
	return slice
}

func pushString(slice []string, data string) []string {
	have := len(slice)
	need := have + 1
	if need > cap(slice) {
		grown := make([]string, have+need) // double the existing space
		copy(grown, slice)                 // transfer to new slice
		slice = grown                      // grown becomes slice
	}
	slice = slice[0:need] // truncate in case too many present
	slice[need-1] = data  // populate with data
	return slice
}

func pushStrings(slice [][]string, data ...[]string) [][]string {
	have := len(slice)
	size := len(data)
	need := have + size
	if need > cap(slice) {
		grown := make([][]string, have+need) // double the existing space
		copy(grown, slice)                   // transfer to new slice
		slice = grown                        // grown becomes slice
	}
	slice = slice[0:need]        // truncate in case too many present
	copy(slice[have:need], data) // populate with data
	return slice
}

func pushRune(slice []rune, data []rune) []rune {
	have := len(slice)
	size := len(data)
	need := have + size
	if need > cap(slice) {
		grown := make([]rune, have+need) // double the existing space
		copy(grown, slice)               // transfer to new slice
		slice = grown                    // grown becomes slice
	}
	slice = slice[0:need]        // truncate in case too many present
	copy(slice[have:need], data) // populate with data
	return slice
}

func pushRunes(slice [][]rune, data ...[]rune) [][]rune {
	have := len(slice)
	size := len(data)
	need := have + size
	if need > cap(slice) {
		grown := make([][]rune, have+need) // double the existing space
		copy(grown, slice)                 // transfer to new slice
		slice = grown                      // grown becomes slice
	}
	slice = slice[0:need]        // truncate in case too many present
	copy(slice[have:need], data) // populate with data
	return slice
}

func pushRuneSlices(slice [][][]rune, data ...[][]rune) [][][]rune {
	have := len(slice)
	size := len(data)
	need := have + size
	if need > cap(slice) {
		grown := make([][][]rune, have+need) // double the existing space
		copy(grown, slice)                   // transfer to new slice
		slice = grown                        // grown becomes slice
	}
	slice = slice[0:need]        // truncate in case too many present
	copy(slice[have:need], data) // populate with data
	return slice
}

func pushMatch(slice [][][2]int, data [][2]int) [][][2]int {
	have := len(slice)
	need := have + 1
	if need > cap(slice) {
		grown := make([][][2]int, have+need) // double the existing space
		copy(grown, slice)                   // transfer to new slice
		slice = grown                        // grown becomes slice
	}
	slice = slice[0:need] // truncate in case too many present
	slice[need-1] = data
	return slice
}

func pushSubMatch(slice [][2]int, sm [2]int) [][2]int {
	have := len(slice)
	need := have + 1
	if need > cap(slice) {
		grown := make([][2]int, have+need) // double the existing space
		copy(grown, slice)                 // transfer to new slice
		slice = grown                      // grown becomes slice
	}
	slice = slice[0:need] // truncate in case too many present
	slice[need-1] = sm    // populate with data
	return slice
}

func mapKeys[K comparable, V interface{}](m map[K]V) []K {
	var slice []K
	for k := range m {
		// no need for appendSlice because mapKeys is only used by NamedClass
		// as a panic message when given a gid <= 0
		slice = append(slice, k)
	}
	return slice
}
