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

type matchGroups []matchStates

func (g matchGroups) Recycle() {
	for _, m := range g {
		m.Destroy()
	}
}

type matchStates []*cMatchState

func (m matchStates) Destroy() {
	for _, state := range m {
		state.Recycle()
	}
}

func (m matchStates) start() int {
	if len(m) > 0 {
		return m[0].start
	}
	// though unlikely, this may happen and returning -1 may not be the best
	// thing to do, just haven't figured out a real-world case
	return -1
}

func (m matchStates) end() int {
	if end := len(m) - 1; end >= 0 {
		return m[end].start + m[end].this
	}
	// though unlikely, this may happen and returning -1 may not be the best
	// thing to do, just haven't figured out a real-world case
	return -1
}
