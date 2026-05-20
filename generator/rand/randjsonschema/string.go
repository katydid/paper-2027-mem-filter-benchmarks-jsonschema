// Copyright 2026 Walter Schulze
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package randjsonschema

import (
	"github.com/katydid/validator-jsonschema-benchmarks/generator/rand"
	"github.com/katydid/validator-jsonschema-benchmarks/generator/rand/randjson"
)

type StringOption func(r *randString)

func String(opts ...StringOption) Rand {
	res := &randString{defaultMaxLen: 10}
	for _, o := range opts {
		o(res)
	}
	return res
}

func WithMaxLength(n uint) func(r *randString) {
	return func(r *randString) {
		r.maxLen = &n
	}
}

func WithMinLength(n uint) func(r *randString) {
	return func(r *randString) {
		r.minLen = &n
	}
}

func WithNonEmpty() func(r *randString) {
	return func(r *randString) {
		n := uint(1)
		r.minLen = &n
	}
}

type randString struct {
	minLen        *uint
	maxLen        *uint
	defaultMaxLen uint
}

func (o *randString) Right(r rand.Rand) string {
	maxLen := int(o.defaultMaxLen)
	if o.maxLen != nil {
		maxLen = int(*o.maxLen)
	}
	minLen := int(0)
	if o.minLen != nil {
		minLen = int(*o.minLen)
	}
	return randjson.String(r, randjson.WithMinStringLength(minLen), randjson.WithMaxStringLength(maxLen))
}

func (o *randString) Wrong(r rand.Rand) string {
	if o.maxLen != nil {
		if r.Intn(2) == 0 {
			// generate a too long string
			return randjson.String(r, randjson.WithMinStringLength(int(*o.maxLen+1)))
		}
	}
	if o.minLen != nil && *o.minLen > 1 {
		if r.Intn(2) == 0 {
			// generate a too short string
			return randjson.String(r, randjson.WithMaxStringLength(int(*o.minLen-1)))
		}
	}
	// generate not a string
	return randjson.Value(r, randjson.NotString())
}
