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

type randConst struct {
	c string
}

func Const(c string) Rand {
	return &randConst{c}
}

func (o *randConst) Right(r rand.Rand) string {
	return o.c
}

func (o *randConst) Wrong(r rand.Rand) string {
	// not the constant
	res := randjson.Value(r)
	for res == o.c {
		res = randjson.Value(r)
	}
	return res
}
