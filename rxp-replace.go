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

// Replacement is a Segment.String replacement function
type Replacement func(s Segment) (replaced string)

// Replace is a Replacement pipeline
type Replace []Replacement

// WithReplace replaces matches using a Replacement function
func (r Replace) WithReplace(replacer Replacement) Replace {
	return append(r, replacer)
}

// WithTransform replaces matches with a Transform function
func (r Replace) WithTransform(transform Transform) Replace {
	return append(r, func(m Segment) (replaced string) {
		return transform(m.String())
	})
}

// WithText replaces matches with the plain text given
func (r Replace) WithText(text string) Replace {
	return append(r, func(s Segment) (replaced string) {
		return text
	})
}

// ToLower is a convenience wrapper around WithTransform and strings.ToLower
func (r Replace) ToLower() Replace {
	return r.WithTransform(strings.ToLower)
}

// ToUpper is a convenience wrapper around WithTransform and strings.ToUpper
func (r Replace) ToUpper() Replace {
	return r.WithTransform(strings.ToUpper)
}
