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

// Pipeline is a list of stages for transforming strings in a single procedure
type Pipeline []Stage

// Process returns the output of a complete Pipeline transformation of the
// input string, and is obviously not a buildable method as it returns a
// string instead of an updated Pipeline
func (p Pipeline) Process(input string) (output string) {
	output = input
	for _, stage := range p {
		output = stage.Process(output)
	}
	return
}

// Transform appends a Transform[string] function to the Pipeline
func (p Pipeline) Transform(transform Transform[string]) Pipeline {
	return append(p, Stage{Transform: transform})
}

// Replace appends a search Pattern and Replace[string] operation to the
// Pipeline
//
// The search argument can be a single Matcher function or a single Pattern
func (p Pipeline) Replace(search interface{}, replace Replace[string]) Pipeline {
	pattern, _, _ := ParseOptions(search)
	return append(p, Stage{Search: pattern, Replace: replace})
}

// Literal appends a search Pattern and literal string replace operation to the
// Pipeline
//
// The search argument can be a single Matcher function or a single Pattern
func (p Pipeline) Literal(search interface{}, text string) Pipeline {
	pattern, _, _ := ParseOptions(search)
	return append(p, Stage{Search: pattern, Transform: func(input string) (output string) {
		return text
	}})
}

// Substitute appends a search Pattern and Transform[string] operation to the
// Pipeline
//
// The search argument can be a single Matcher function or a single Pattern
func (p Pipeline) Substitute(search interface{}, transform Transform[string]) Pipeline {
	pattern, _, _ := ParseOptions(search)
	return append(p, Stage{Search: pattern, Transform: transform})
}
