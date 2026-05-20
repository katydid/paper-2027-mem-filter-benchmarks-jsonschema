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
	"strings"

	"github.com/katydid/validator-jsonschema-benchmarks/generator/rand"
	"github.com/katydid/validator-jsonschema-benchmarks/generator/rand/randjson"
)

type ArrayOption func(r *randArray)

type randArray struct {
	Item   Rand
	minLen int
}

func WithMinItems(n int) func(r *randArray) {
	return func(r *randArray) {
		r.minLen = n
	}
}

func ArrayOf(item Rand, opts ...ArrayOption) Rand {
	res := &randArray{Item: item}
	for _, o := range opts {
		o(res)
	}
	return res
}

func (o *randArray) Right(r rand.Rand) string {
	n := r.Intn(10) + o.minLen
	items := []string{}
	for range n {
		items = append(items, o.Item.Right(r))
	}
	return "[" + strings.Join(items, ",") + "]"
}

func (o *randArray) Wrong(r rand.Rand) string {
	switch r.Intn(10) {
	case 0:
		return randjson.Value(r, randjson.NotArray())
	default:
		if o.minLen > 0 {
			if r.Intn(10) == 0 {
				return "[]"
			}
		}
		n := r.Intn(9) + 1
		wrongIndex := r.Intn(n)
		items := []string{}
		for i := range n {
			if i == wrongIndex {
				items = append(items, o.Item.Wrong(r))
			} else {
				items = append(items, o.Item.Right(r))
			}
		}
		return "[" + strings.Join(items, ",") + "]"
	}
}
