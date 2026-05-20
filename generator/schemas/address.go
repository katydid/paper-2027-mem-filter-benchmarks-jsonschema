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
	. "github.com/katydid/validator-jsonschema-benchmarks/generator/rand/randjsonschema"
)

// https://json-schema.org/learn/json-schema-examples#address
// A schema representing an address, with optional properties for different address components
// which enforces that locality, region, and countryName are required,
// and if postOfficeBox or extendedAddress is provided, streetAddress must also be provided.
const SchemaJSONSchemaExampleAddress = `{
  "$id": "https://example.com/address.schema.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "description": "An address similar to http://microformats.org/wiki/h-card",
  "type": "object",
  "properties": {
    "postOfficeBox": {
      "type": "string"
    },
    "extendedAddress": {
      "type": "string"
    },
    "streetAddress": {
      "type": "string"
    },
    "locality": {
      "type": "string"
    },
    "region": {
      "type": "string"
    },
    "postalCode": {
      "type": "string"
    },
    "countryName": {
      "type": "string"
    }
  },
  "required": [ "locality", "region", "countryName" ],
  "dependentRequired": {
    "postOfficeBox": [ "streetAddress" ],
    "extendedAddress": [ "streetAddress" ]
  }
}`

func RandomAddress() Rand {
	otherFields := []Rand{
		Field("locality", String(WithNonEmpty()), IsRequired()),
		Field("region", String(WithNonEmpty()), IsRequired()),
		Field("postalCode", String(WithNonEmpty())),
		Field("countryName", String(WithNonEmpty()), IsRequired()),
	}
	obj := func(fields ...Rand) Rand {
		return Object(WithAdditionalFields(),
			WithFields(append(fields, otherFields...)...),
		)
	}
	a := Or(
		WithRight(
			obj(
				Field("postOfficeBox", String(WithNonEmpty()), IsRequired()),
				Field("extendedAddress", String(WithNonEmpty())),
				Field("streetAddress", String(WithNonEmpty()), IsRequired()),
			),
			obj(
				Field("postOfficeBox", String(WithNonEmpty())),
				Field("extendedAddress", String(WithNonEmpty()), IsRequired()),
				Field("streetAddress", String(WithNonEmpty()), IsRequired()),
			),
			obj(
				Field("streetAddress", String(WithNonEmpty())),
			),
		),
		WithWrong(
			obj(
				Field("postOfficeBox", String(WithNonEmpty()), IsRequired()),
				Field("extendedAddress", String(WithNonEmpty())),
			),
			obj(
				Field("postOfficeBox", String(WithNonEmpty())),
				Field("extendedAddress", String(WithNonEmpty()), IsRequired()),
			),
		),
	)
	return a
}
