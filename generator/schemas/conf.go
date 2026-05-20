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

package schemas

import (
	"strconv"

	"github.com/katydid/validator-jsonschema-benchmarks/generator/rand"

	. "github.com/katydid/validator-jsonschema-benchmarks/generator/rand/randjsonschema"
)

const SchemaConfIsIn2026OrLate2025AndEU = `
{
	"definitions": {
		"due": {
			"type": "object",
			"anyOf": [
			{
				"properties": {
					"Year": {
						"$ref": "#/definitions/year2026"
					}
				},
				"required": ["Year"]
			},
			{
				"allOf": [                  
					{
						"properties": {
							"Year": {
								"$ref": "#/definitions/year2025"
							}
						},
						"required": ["Year"]
					},
					{
						"properties": {
							"Month": {
								"$ref": "#/definitions/month10"
							}
						},
						"required": ["Month"]
					}
				]
			}
			]
		},
		"loc": {
			"type": "object",
			"properties": {
				"Cont": {
					"$ref": "#/definitions/conteu"
				}
			},
			"required": ["Cont"]
		},
		"year2026": {
			"const": "2026"
		},
		"year2025": {
			"const": "2025"
		},
		"month10": {
			"minimum": 10
		},
		"conteu": {
			"const": "EU"
		}
	},
	"type": "object",
	"properties": {
		"Due": {
			"$ref": "#/definitions/due"
		},
		"Loc": {
			"$ref": "#/definitions/loc"
		}
	},
	"required": ["Due", "Loc"]
}`

func RandomConfIsIn2026OrLate2025AndEU() Rand {
	return Object(
		WithAdditionalFields(),
		WithAlwaysRightFields(
			Field("Name", String()),
			Field("Notify", Object(WithAdditionalFields(), WithFields(
				Field("Year", String(), IsRequired()),
				Field("Month", String(), IsRequired()),
				Field("Day", String(), IsRequired()),
			))),
			Field("Category", String()),
			Field("Tags", ArrayOf(String())),
		),
		WithFields(
			Field("Due", Due(), IsRequired()),
			Field("Loc", Object(
				WithAdditionalFields(),
				WithAlwaysRightFields(
					Field("Ctry", String()),
					Field("City", String()),
				),
				WithFields(
					Field("Cont", EU(), IsRequired()),
				),
			), IsRequired()),
		),
	)
}

type randDue struct{}

func Due() Rand {
	return &randDue{}
}

func (o *randDue) Right(r rand.Rand) string {
	switch r.Intn(3) {
	case 0:
		return Object(WithFields(
			Field("Year", Const(`"2026"`), IsRequired()),
			Field("Month", Integer(), IsRequired()),
			Field("Day", String(), IsRequired()),
		)).Right(r)
	case 1:
		return Object(WithFields(
			Field("Year", Const(`"2025"`), IsRequired()),
			Field("Month", Const(`11`), IsRequired()),
			Field("Day", String(), IsRequired()),
		)).Right(r)
	case 2:
		return Object(WithFields(
			Field("Year", Const(`"2025"`), IsRequired()),
			Field("Month", Const(`12`), IsRequired()),
			Field("Day", String(), IsRequired()),
		)).Right(r)
	}
	panic("unreachable")
}

func (o *randDue) Wrong(r rand.Rand) string {
	switch r.Intn(3) {
	case 0:
		return Object(WithFields(
			Field("Year", Const(`"2027"`), IsRequired()),
			Field("Month", Integer(), IsRequired()),
			Field("Day", String(), IsRequired()),
		)).Right(r)
	case 1:
		return Object(WithFields(
			Field("Year", Const(`"2025"`), IsRequired()),
			Field("Month", Const(`7`), IsRequired()),
			Field("Day", String(), IsRequired()),
		)).Right(r)
	case 2:
		return Object(WithFields(
			Field("Year", Const(`"1991"`), IsRequired()),
			Field("Month", Const(`7`), IsRequired()),
			Field("Day", String(), IsRequired()),
		)).Right(r)
	}
	panic("unreachable")
}

type randEU struct{}

func EU() Rand {
	return &randEU{}
}

var conts = []string{"AF", "AN", "AS", "EU", "NA", "SA", "OC"}

func (o *randEU) Right(r rand.Rand) string {
	return `"EU"`
}

func (o *randEU) Wrong(r rand.Rand) string {
	c := conts[r.Intn(len(conts))]
	for c == "EU" {
		c = conts[r.Intn(len(conts))]
	}
	return strconv.Quote(c)
}
