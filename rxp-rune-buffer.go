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

// RuneBuffer is an efficient rune based buffer
type RuneBuffer struct {
	len int
	buf runes.RuneReader
}

// NewRuneBuffer creates a new RuneBuffer instance for the given input string
func NewRuneBuffer[V []rune | []byte | string](input V) *RuneBuffer {
	rb := &RuneBuffer{}
	rb.len = len(input)
	rb.buf = runes.NewRuneReader(input)
	return rb
}

// Len returns the total number of runes in the RuneBuffer
func (rb *RuneBuffer) Len() int {
	return rb.len
}

// Get returns the Ready rune at the given index position
func (rb *RuneBuffer) Get(index int) (r rune, size int, ok bool) {
	if ok = rb.Ready(index); ok {
		r, size, _ = rb.buf.ReadRuneAt(int64(index))
	}
	return
}

// Ready returns true if the given index position is greater than or equal to
// zero and less than the total length of the RuneBuffer
func (rb *RuneBuffer) Ready(index int) (ready bool) {
	return 0 <= index && index < rb.len
}

// Valid returns true if the given index position is greater than or equal to
// zero and less than or equal to the total length of the RuneBuffer
func (rb *RuneBuffer) Valid(index int) (valid bool) {
	return 0 <= index && index <= rb.len
}

// Invalid returns true if the given index position is less than zero or greater
// than or equal to the total length of the RuneBuffer
func (rb *RuneBuffer) Invalid(index int) (invalid bool) {
	return index < 0 || index >= rb.len
}

// End returns true if this index is exactly the input length, denoting the
// Dollar zero-width position
func (rb *RuneBuffer) End(index int) (end bool) {
	return index == rb.len
}

// Prev returns the Ready rune before the given index position, or \0 if not
// Ready
//
// The Prev is necessary because to find the previous rune to the given index,
// Prev must incrementally scan backwards up to four bytes, trying to read a
// rune without error with each iteration
func (rb *RuneBuffer) Prev(index int) (r rune, size int, ok bool) {
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
func (rb *RuneBuffer) Next(index int) (r rune, size int, ok bool) {
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
func (rb *RuneBuffer) Slice(index, count int) (slice []rune, size int) {
	if rb.Ready(index) {
		slice, size, _ = rb.buf.ReadRuneSlice(int64(index), int64(count))
	}
	return
}

// String returns the string of runes from start (inclusive) to end (exclusive)
// if the entire range is Ready
func (rb *RuneBuffer) String(index, count int) string {
	if rb.Ready(index) {
		c := int64(count)
		if c < 0 {
			c = rb.buf.Size() - int64(index+count) - 1
		}
		slice, _, _ := rb.buf.ReadRuneSlice(int64(index), c)
		return string(slice)
	}
	return ""
}
