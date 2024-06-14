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

// Stage is one phase of a text replacement Pipeline and receives an input
// string from the previous stage (or the initial input text) and returns the
// output provided to the next stage (or is finally returned to the caller)
type Stage struct {
	// Search is a Pattern of text, used with Replace to modify the matching
	// text
	Search    Pattern
	Replace   Replace[string]
	Transform Transform[string]
}

// Process is the Pipeline processor method for a Stage
//
// If there is a Search Pattern present, and there is at least one Replacer
// functions present, then the process returns a Search
// Pattern.ReplaceAllString
//
// If there is a Search Pattern present, and there are no Replacer functions
// present, then the process returns a Search Pattern.ReplaceAllStringFunc
func (s Stage) Process(input string) (output string) {
	if s.Transform != nil {
		if len(s.Search) > 0 {
			return s.Search.ReplaceAllStringFunc(input, s.Transform)
		}
		return s.Transform(input)
	}
	return s.Search.ReplaceAllString(input, s.Replace)
}
