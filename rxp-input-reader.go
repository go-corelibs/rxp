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
	"github.com/go-corelibs/runes"
)

// InputReader is an efficient rune based buffer
type InputReader struct {
	len int
	buf runes.RuneReader
}

// NewInputReader creates a new InputReader instance for the given input string
func NewInputReader[V []rune | []byte | string](input V) *InputReader {
	rb := &InputReader{}
	rb.len = len(input)
	rb.buf = runes.NewRuneReader(input)
	return rb
}

// Len returns the total number of runes in the InputReader
func (rb *InputReader) Len() int {
	return rb.len
}

// Get returns the Ready rune at the given index position
func (rb *InputReader) Get(index int) (r rune, size int, ok bool) {
	if ok = 0 <= index && index < rb.len; ok {
		r, size, _ = rb.buf.ReadRuneAt(int64(index))
	}
	return
}

// Ready returns true if the given index position is greater than or equal to
// zero and less than the total length of the InputReader
func (rb *InputReader) Ready(index int) (ready bool) {
	return 0 <= index && index < rb.len
}

// Valid returns true if the given index position is greater than or equal to
// zero and less than or equal to the total length of the InputReader
func (rb *InputReader) Valid(index int) (valid bool) {
	return 0 <= index && index <= rb.len
}

// Invalid returns true if the given index position is less than zero or greater
// than or equal to the total length of the InputReader
func (rb *InputReader) Invalid(index int) (invalid bool) {
	return index < 0 || index >= rb.len
}

// End returns true if this index is exactly the input length, denoting the
// Dollar zero-width position
func (rb *InputReader) End(index int) (end bool) {
	return index == rb.len
}

// Prev returns the Ready rune before the given index position, or \0 if not
// Ready
//
// The Prev is necessary because to find the previous rune to the given index,
// Prev must incrementally scan backwards up to four bytes, trying to read a
// rune without error with each iteration
func (rb *InputReader) Prev(index int) (r rune, size int, ok bool) {
	var err error
	if index > 0 {
		if r, size, err = rb.buf.ReadPrevRuneFrom(int64(index)); err == nil {
			ok = true
			return
		}
	}
	// this is no previous rune from the given index
	return 0, 0, false
}

// Next returns the Ready rune after the given index position, or \0 if not
// Ready
func (rb *InputReader) Next(index int) (r rune, size int, ok bool) {
	var err error
	if r, size, err = rb.buf.ReadNextRuneFrom(int64(index)); err == nil {
		ok = true
		return
	}
	// this is no previous rune from the given index
	return 0, 0, false
}

// Slice returns the range of runes from start (inclusive) to end (exclusive)
// if the entire range is Ready
func (rb *InputReader) Slice(index, count int) (slice []rune, size int) {
	if rb.Ready(index) {
		slice, size, _ = rb.buf.ReadRuneSlice(int64(index), int64(count))
	}
	return
}

// String returns the string of runes from start (inclusive) to end (exclusive)
// if the entire range is Ready
func (rb *InputReader) String(index, count int) string {
	if rb.Ready(index) {
		c := int64(count)
		if c < 0 {
			c = int64(rb.len) - int64(index+count) - 1
		}
		data, _ := rb.buf.ReadString(int64(index), c)
		return data
	}
	return ""
}

// Bytes returns the range of runes from start (inclusive) to end (exclusive)
// if the entire range is Ready
func (rb *InputReader) Bytes(index, count int) (slice []byte) {
	if rb.Ready(index) {
		c := int64(count)
		if c < 0 {
			c = int64(rb.len) - int64(index+count) - 1
		}
		slice, _ = rb.buf.ReadByteSlice(int64(index), c)
	}
	return
}
