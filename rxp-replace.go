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
	"strings"
	"unicode"
)

// Transform is the function signature for non-rxp string transformation
// pipeline stages
type Transform[V []rune | []byte | string] func(input V) (output V)

// Replacer is a Replace processor function
//
// The captured argument is the result of the Pattern match process and is
// composed of the entire matched text as the first item in the captured list,
// and any Pattern capture groups following
//
// The modified string argument is the output of the previous Replacer in the
// Replace process, or the original matched input text if this is the first
// Replacer in the process
type Replacer[V []rune | []byte | string] func(input *InputReader, captured [][2]int, modified V) (replaced V)

// Replace is a Replacer pipeline
type Replace[V []rune | []byte | string] []Replacer[V]

func (r Replace[V]) process(input *InputReader, captured [][2]int, data V) (replaced V) {
	replaced = data
	for _, rpl := range r {
		replaced = rpl(input, captured, replaced)
	}
	return
}

// WithReplace replaces matches using a Replacer function
func (r Replace[V]) WithReplace(replacer Replacer[V]) Replace[V] {
	return append(r, replacer)
}

// WithTransform replaces matches with a Transform function
func (r Replace[V]) WithTransform(transform Transform[V]) Replace[V] {
	return append(r, func(input *InputReader, captured [][2]int, text V) (replaced V) {
		return transform(text)
	})
}

// WithLiteral replaces matches with the literal value given
func (r Replace[V]) WithLiteral(data V) Replace[V] {
	return append(r, func(input *InputReader, captured [][2]int, text V) (replaced V) {
		return data
	})
}

// ToLower is a convenience method which, based on Replace type, will use one of
// the following methods to lower-case all the lower-case-able values
//
//	| Type   | Function        |
//	|--------|-----------------|
//	| []rune | unicode.ToLower |
//	| []byte | bytes.ToLower   |
//	| string | strings.ToLower |
func (r Replace[V]) ToLower() Replace[V] {
	return append(r, func(input *InputReader, captured [][2]int, data V) (replaced V) {
		switch t := interface{}(&data).(type) {
		case *[]rune:
			slice := make([]rune, len(*t))
			for idx, ch := range *t {
				slice[idx] = unicode.ToLower(ch)
			}
			return interface{}(slice).(V)
		case *[]byte:
			slice := bytes.ToLower(*t)
			return interface{}(slice).(V)
		case *string:
			slice := strings.ToLower(*t)
			return interface{}(slice).(V)
		default:
			panic("the universe is broken")
		}
	})
}

// ToUpper is a convenience method which, based on Replace type, will use one of
// the following methods to upper-case all the upper-case-able values
//
//	| Type   | Function        |
//	|--------|-----------------|
//	| []rune | unicode.ToUpper |
//	| []byte | bytes.ToUpper   |
//	| string | strings.ToUpper |
func (r Replace[V]) ToUpper() Replace[V] {
	return append(r, func(input *InputReader, captured [][2]int, data V) (replaced V) {
		switch t := interface{}(&data).(type) {
		case *[]rune:
			slice := make([]rune, len(*t))
			for idx, ch := range *t {
				slice[idx] = unicode.ToUpper(ch)
			}
			return interface{}(slice).(V)
		case *[]byte:
			slice := bytes.ToUpper(*t)
			return interface{}(slice).(V)
		case *string:
			slice := strings.ToUpper(*t)
			return interface{}(slice).(V)
		default:
			panic("the universe is broken")
		}
	})
}
