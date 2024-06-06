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

// appendSlice is a slightly better option than stock append and using the one
// from go-corelibs/x-sync crosses the package boundary and actually makes a
// difference in must-fast-path scenarios
func appendSlice[V interface{}](slice []V, data ...V) []V {
	m := len(slice)     // current length
	n := m + len(data)  // needed length
	if n > cap(slice) { // current cap size
		grown := make([]V, (m+1)*2) // double the existing space
		copy(grown, slice)          // transfer to new slice
		slice = grown               // grown becomes slice
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