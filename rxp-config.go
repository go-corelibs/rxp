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

var (
	spConfig = newSyncPool(10, func() *Config {
		return new(Config)
	})
)

// Config is the structure for pattern matching settings
//
//	| Repetition | Settings                  |
//	|            | Minimum | Maximum |  Less |
//	|------------|---------|---------|-------|
//	|     x*     | 0 or -1 |    -1   | false |
//	|     x+     |    1    |    -1   | false |
//	|     x?     |    0    |     1   | false |
//	|   x{n,m}   |    n    |     m   | false |
//	|   x{n,}    |    n    |    -1   | false |
//	|   x{n}     |    n    |     n   | false |
//	|     x*?    |    0    |    -1   | true  |
//	|     x+?    |    1    |    -1   | true  |
//	|     x??    |    0    |     1   | true  |
//	|   x{n,m}?  |    n    |     m   | true  |
//	|   x{n,}?   |    n    |    -1   | true  |
//	|   x{n}?    |    n    |     n   | true  |
type Config struct {
	Negated   bool // invert the meaning of this match group, ie: [^a-zA-Z] negates an IsAlphabet Matcher
	Multiline bool // ^ and $ match begin/end of line, not just start/end of input
	DotNL     bool // . matches \n
	AnyCase   bool // case-insensitive
	Capture   bool // is a capture group
	Less      bool // prefer less
	Minimum   int  // the first number  in {n1,n2}
	Maximum   int  // the second number in {n1,n2}
	recycled  bool
	private   *cMatchState
}

func DefaultConfig() *Config {
	return &Config{Minimum: 1, Maximum: 1}
}

func newConfig() *Config {
	cfg := spConfig.Get()
	if !cfg.recycled {
		// pool created a new one, let's add ten more so next batch isn't new
		spConfig.Seed(-1)
	}
	cfg.Negated = false
	cfg.Multiline = false
	cfg.DotNL = false
	cfg.AnyCase = false
	cfg.Capture = false
	cfg.Less = false
	cfg.Minimum = 1
	cfg.Maximum = 1
	cfg.private = nil
	cfg.recycled = false
	return cfg
}

func (cfg *Config) Recycle() {
	if !cfg.recycled {
		cfg.private = nil
		cfg.recycled = true
		spConfig.Put(cfg)
	}
}

// Valid returns true if the repetition range:
//
// - does not have a negative minimum and a positive maximum
// - either the minimum is less than or equal to the maximum, or
// - the minimum is positive and the maximum is negative (unlimited)
func (cfg *Config) Valid() (ok bool) {
	if cfg.Minimum < 0 && cfg.Maximum > 0 {
		// cannot have "no min" and yet a max
		// not even sure this scenario is possible
		return
	}
	return (cfg.Minimum <= cfg.Maximum) || (cfg.Minimum > 0 && cfg.Maximum < 0)
}

// IsValidCount returns true if the given count satisfies the minimum and
// maximum range constraints
func (cfg *Config) IsValidCount(count int) (ok bool) {
	if cfg.Minimum < 0 {
		// no min req also means no limit
		return true
	}
	if count >= cfg.Minimum {
		// met the min req
		if cfg.Maximum > 0 {
			// though there is a limit
			return count <= cfg.Maximum
		}
		// met min, no limit
		return true
	}
	// count did not meet the minimum
	// using MakeRuneMatcher, makes this case very difficult to reach
	return false
}

func (cfg *Config) Apply(other *Config) {

	other.Negated = cfg.Negated
	other.Multiline = cfg.Multiline
	other.DotNL = cfg.DotNL
	other.AnyCase = cfg.AnyCase
	other.Capture = cfg.Capture
	other.Less = cfg.Less
	other.Minimum = cfg.Minimum
	other.Maximum = cfg.Maximum

}

// Merge returns a clone of this Config with all the non-zero values of the
// other Config applied
func (cfg *Config) Merge(other *Config) (merged *Config) {
	v := *cfg // struct copy is faster than newConfig + .Apply
	merged = &v
	if other != nil {
		if other.Negated {
			merged.Negated = other.Negated
		}
		if other.Multiline {
			merged.Multiline = other.Multiline
		}
		if other.DotNL {
			merged.DotNL = other.DotNL
		}
		if other.AnyCase {
			merged.AnyCase = other.AnyCase
		}
		if other.Capture {
			merged.Capture = other.Capture
		}
		if other.Less {
			merged.Less = other.Less
		}
		if other.Minimum != 1 {
			merged.Minimum = other.Minimum
		}
		if other.Maximum != 1 {
			merged.Maximum = other.Maximum
		}
	}
	return merged
}
