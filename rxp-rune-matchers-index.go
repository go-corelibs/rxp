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

// IndexEnd returns true if this index is exactly the input length, denoting the
// Dollar zero-width position
func IndexEnd(input []rune, index int) bool {
	return index >= len(input)
}

// IndexReady returns true if this index is less than the input length
func IndexReady(input []rune, index int) bool {
	return 0 <= index && index < len(input)
}

// IndexValid returns true if this index is less than or equal to the input length
func IndexValid(input []rune, index int) bool {
	return 0 <= index && index <= len(input)
}

// IndexInvalid returns true if this index is greater than or equal to the input length
func IndexInvalid(input []rune, index int) bool {
	return index < 0 || index >= len(input)
}

// IndexGet returns the input rune at the given index, if there is one
func IndexGet(input []rune, index int) (r rune, ok bool) {
	if ok = 0 <= index && index < len(input); ok {
		r = input[index]
	}
	return
}
