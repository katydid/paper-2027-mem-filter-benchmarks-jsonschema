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
	"github.com/katydid/validator-jsonschema-benchmarks/generator/rand/randjson"
)

type NumberOption func(r *randNumber)

func Number(opts ...NumberOption) Rand {
	res := &randNumber{}
	for _, o := range opts {
		o(res)
	}
	return res
}

func Integer(opts ...NumberOption) Rand {
	res := &randNumber{}
	for _, o := range opts {
		o(res)
	}
	res.isInteger = true
	return res
}

func WithMinimum(n uint) func(r *randNumber) {
	return func(r *randNumber) {
		r.minimum = &n
	}
}

func WithMultipleOf(n uint) func(r *randNumber) {
	return func(r *randNumber) {
		r.multipleOf = &n
	}
}

type randNumber struct {
	minimum    *uint
	multipleOf *uint
	isInteger  bool
}

func (o *randNumber) Right(r rand.Rand) string {
	if o.minimum != nil && r.Intn(10) == 0 {
		return strconv.Itoa(int(*o.minimum))
	}
	for {
		num := randjson.Number(r)
		if o.isInteger {
			num = randjson.Integer(r)
		}
		i, err := strconv.Atoi(num)
		if err != nil {
			continue
		}
		if o.minimum != nil && i < int(*o.minimum) {
			continue
		}
		if o.multipleOf != nil {
			min := 0
			if o.minimum != nil {
				min = int(*o.minimum)
			}
			i = (r.Intn(1000) + min) * int(*o.multipleOf)
			num = strconv.Itoa(i)
		}
		return num
	}
}

func (o *randNumber) Wrong(r rand.Rand) string {
	if r.Intn(2) == 0 {
		// generate not a number
		return randjson.String(r)
	}
	var num string
	for {
		num = randjson.Number(r)
		i, err := strconv.Atoi(num)
		if err != nil {
			continue
		}
		if o.minimum != nil && i >= int(*o.minimum) {
			continue
		}
		if o.multipleOf != nil && ((i/int(*o.multipleOf))*int(*o.multipleOf)) == i {
			continue
		}
		if o.minimum == nil && o.multipleOf == nil {
			return randjson.Value(r, randjson.NotNumber())
		}
		return num
	}
}
