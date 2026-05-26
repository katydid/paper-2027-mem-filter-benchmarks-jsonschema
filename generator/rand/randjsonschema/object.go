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

type ObjectOption func(r *randObject)

type randObject struct {
	fields            []Rand
	additionalFields  bool
	alwaysRightFields []Rand
}

func Object(opts ...ObjectOption) Rand {
	res := &randObject{}
	for _, o := range opts {
		o(res)
	}
	return res
}

func WithFields(fields ...Rand) func(r *randObject) {
	return func(r *randObject) {
		r.fields = fields
	}
}

func WithAdditionalFields() func(r *randObject) {
	return func(r *randObject) {
		r.additionalFields = true
	}
}

// Always Right Fields are always generated as Right
func WithAlwaysRightFields(fields ...Rand) func(r *randObject) {
	return func(r *randObject) {
		r.alwaysRightFields = fields
	}
}

func (o *randObject) Right(r rand.Rand) string {
	if len(o.fields) == 0 && !o.additionalFields {
		return "{}"
	}

	fieldmap := map[Rand]bool{}
	fields := []Rand{}
	for i := range o.fields {
		fieldmap[o.fields[i]] = true
		fields = append(fields, o.fields[i])
	}
	for i := range o.alwaysRightFields {
		fieldmap[o.alwaysRightFields[i]] = true
		fields = append(fields, o.alwaysRightFields[i])
	}
	if o.additionalFields && r.Intn(2) == 0 {
		additionalField := Field("AdditionalField", String(WithMinLength(1)), IsRequired())
		fieldmap[additionalField] = true
		fields = append(fields, additionalField)
	}
	rand.Shuffle(r, fields)

	fieldstrs := []string{}
	for i, field := range fields {
		right := fieldmap[fields[i]]
		fstr := field.Right(r)
		if !right {
			fstr = field.Wrong(r)
		}
		if len(fstr) > 0 {
			fieldstrs = append(fieldstrs, fstr)
		}
	}

	return "{" + strings.Join(fieldstrs, ",") + "}"
}

func (o *randObject) Wrong(r rand.Rand) string {
	if len(o.fields) == 0 {
		return randjson.Value(r, randjson.NotObject())
	}

	wrongIndex := r.Intn(len(o.fields))

	fieldmap := map[Rand]bool{}
	fields := []Rand{}
	for i := range o.fields {
		fieldmap[o.fields[i]] = i != wrongIndex
		fields = append(fields, o.fields[i])
	}
	for i := range o.alwaysRightFields {
		fieldmap[o.alwaysRightFields[i]] = true
		fields = append(fields, o.alwaysRightFields[i])
	}
	if !o.additionalFields && r.Intn(2) == 0 {
		additionalField := Field("AdditionalField", String(WithMinLength(1)), IsRequired())
		fieldmap[additionalField] = true
		fields = append(fields, additionalField)
	}
	rand.Shuffle(r, fields)

	fieldstrs := []string{}
	for i, field := range fields {
		right := fieldmap[fields[i]]
		fstr := field.Right(r)
		if !right {
			fstr = field.Wrong(r)
		}
		if len(fstr) > 0 {
			fieldstrs = append(fieldstrs, fstr)
		}
	}

	return "{" + strings.Join(fieldstrs, ",") + "}"
}
