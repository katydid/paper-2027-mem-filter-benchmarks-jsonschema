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
	"strconv"

	"github.com/katydid/validator-jsonschema-benchmarks/generator/rand"
)

type FieldOption func(r *randField)

func Field(name string, value Rand, opts ...FieldOption) Rand {
	res := &randField{name: name, value: value}
	for _, o := range opts {
		o(res)
	}
	return res
}

type randField struct {
	name     string
	required bool
	value    Rand
}

func IsRequired() func(r *randField) {
	return func(r *randField) {
		r.required = true
	}
}

func (o *randField) Right(r rand.Rand) string {
	if !o.required && r.Intn(5) == 0 {
		return ""
	}
	return strconv.Quote(o.name) + ":" + o.value.Right(r)
}

func (o *randField) Wrong(r rand.Rand) string {
	if o.required && r.Intn(3) == 0 {
		return ""
	}
	return strconv.Quote(o.name) + ":" + o.value.Wrong(r)
}
