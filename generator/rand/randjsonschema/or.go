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

import "github.com/katydid/validator-jsonschema-benchmarks/generator/rand"

type randOr struct {
	rightArgs []Rand
	wrongArgs []Rand
	args      []Rand
}

type OrOption func(r *randOr)

func Or(opts ...OrOption) Rand {
	res := &randOr{}
	for _, o := range opts {
		o(res)
	}
	if len(res.rightArgs) == 0 && len(res.args) == 0 {
		panic("Or was not correctly initialized")
	}
	return res
}

func WithAnyOf(args ...Rand) func(r *randOr) {
	return func(r *randOr) {
		r.args = args
	}
}

func WithRight(args ...Rand) func(r *randOr) {
	return func(r *randOr) {
		r.rightArgs = args
	}
}

func WithWrong(args ...Rand) func(r *randOr) {
	return func(r *randOr) {
		r.wrongArgs = args
	}
}

func (o *randOr) Right(r rand.Rand) string {
	i := r.Intn(len(o.rightArgs) + len(o.args))
	if i < len(o.rightArgs) {
		return o.rightArgs[i].Right(r)
	}
	return o.args[i-len(o.rightArgs)].Right(r)
}

func (o *randOr) Wrong(r rand.Rand) string {
	i := r.Intn(len(o.wrongArgs) + len(o.args))
	if i < len(o.wrongArgs) {
		return o.wrongArgs[i].Right(r)
	}
	return o.args[i-len(o.wrongArgs)].Wrong(r)
}
