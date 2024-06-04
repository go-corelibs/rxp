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

// Pipeline is a list of stages for transforming text in a single procedure
//
// Pipeline is also a pseudo-buildable thing using the Transform, Replace,
// ReplaceText and ReplaceWith methods which return the updated Pipeline
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

func (p Pipeline) Transform(transform Transform) Pipeline {
	return append(p, Stage{Transform: transform})
}

func (p Pipeline) Replace(search interface{}, replace Replace) Pipeline {
	pattern, _, _ := ParseOptions(search)
	return append(p, Stage{Search: pattern, Replace: replace})
}

func (p Pipeline) ReplaceText(search interface{}, text string) Pipeline {
	pattern, _, _ := ParseOptions(search)
	return append(p, Stage{Search: pattern, Transform: func(input string) (output string) {
		return text
	}})
}

func (p Pipeline) ReplaceWith(search interface{}, transform Transform) Pipeline {
	pattern, _, _ := ParseOptions(search)
	return append(p, Stage{Search: pattern, Transform: transform})
}
