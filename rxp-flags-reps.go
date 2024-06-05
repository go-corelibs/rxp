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

const (
	DefaultMinReps = 1
	DefaultMaxReps = 1
)

var gDefaultReps = Reps{1, 1}

type Reps []int

// Min returns the minimum number of repetitions, or the DefaultMinReps
// for nil instances
func (r Reps) Min() int {
	if len(r) == 0 {
		return 1
	}
	return r[0]
}

// Max returns the maximum number of repetitions, or the DefaultMaxReps
// for nil instances
func (r Reps) Max() int {
	if len(r) <= 1 {
		return 1
	}
	return r[1]
}

// IsNil returns true if the length of the Reps int slice is not exactly two
func (r Reps) IsNil() bool {
	return len(r) < 2
}

// Valid returns true if the Reps is not nil and if the minimum is not greater
// than the maximum, unless the maximum is unlimited (zero or a negative)
func (r Reps) Valid() bool {
	if r.IsNil() {
		return false
	}
	minimum, maximum := r.Min(), r.Max()
	if minimum > 0 && maximum > 0 && minimum > maximum {
		return false
	}
	return true
}

// Satisfied returns true if the given count meets the required repetitions
func (r Reps) Satisfied(count int) (minHit, maxHit bool) {
	if r.IsNil() {
		return
	}

	minimum, maximum := r.Min(), r.Max()

	if minHit = minimum >= 0 && count >= minimum; minHit {
		// has a min req and this count is good
	} else if minHit = minimum < 0 && count >= 0; minHit {
		// has unlimited min req and this count is good
	}

	if maxHit = maximum > 0 && count >= maximum; maxHit {
		// has a max req and this count is good
	}

	return
}
