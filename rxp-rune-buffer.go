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

var (
	spRuneBuffer = sync.NewPool(1, func() *RuneBuffer {
		return &RuneBuffer{}
	}, func(v *RuneBuffer) *RuneBuffer {
		// getter
		return v
	}, func(v *RuneBuffer) *RuneBuffer {
		// setter
		v.buf = nil
		return v
	})
)

// RuneBuffer is an efficient rune based buffer
type RuneBuffer struct {
	len int
	buf []rune
}

// NewRuneBuffer creates a new RuneBuffer instance for the given input string
func NewRuneBuffer(input []rune) *RuneBuffer {
	rb := spRuneBuffer.Get()
	rb.len = len(input)
	rb.buf = input
	return rb
}

func (rb *RuneBuffer) recycle() {
	spRuneBuffer.Put(rb)
}

// Len returns the total number of runes in the RuneBuffer
func (rb *RuneBuffer) Len() int {
	return len(rb.buf)
}

// Get returns the Ready rune at the given index position
func (rb *RuneBuffer) Get(index int) (r rune, ok bool) {
	if ok = rb.Ready(index); ok {
		r = rb.buf[index]
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

// Invalid returns true if the given index position is less than zero or equal
// greater than or equal to the total length of the RuneBuffer
func (rb *RuneBuffer) Invalid(index int) (invalid bool) {
	return index < 0 || index >= rb.len
}

// End returns true if this index is exactly the input length, denoting the
// Dollar zero-width position
func (rb *RuneBuffer) End(index int) (end bool) {
	return index == rb.len
}

// Slice returns the range of runes from start (inclusive) to end (exclusive)
// if the entire range is Ready
func (rb *RuneBuffer) Slice(start, end int) (slice []rune) {
	//slice = rb.buf[start:end]
	if end >= 0 {
		if rb.Ready(start) && rb.Valid(end) {
			slice = rb.buf[start:end]
		}
	} else {
		if rb.Ready(start) {
			slice = rb.buf[start:]
		}
	}
	return
}

// String returns the string of runes from start (inclusive) to end (exclusive)
// if the entire range is Ready
func (rb *RuneBuffer) String(start, end int) (slice string) {
	buf := spStringBuilder.Get()
	defer spStringBuilder.Put(buf)
	if start <= end {
		if rb.Ready(start) && rb.Valid(end) {
			for _, r := range rb.buf[start:end] {
				buf.WriteRune(r)
			}
		}
	} else if start >= 0 && end <= 0 {
		if rb.Ready(start) {
			for _, r := range rb.buf[start:] {
				buf.WriteRune(r)
			}
		}
	}
	return buf.String()
}
