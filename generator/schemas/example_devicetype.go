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

// https://json-schema.org/learn/json-schema-examples#device-type
const SchemaJSONSchemaExampleDevicetype = `
{
  "$id": "https://example.com/device.schema.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "properties": {
    "deviceType": {
      "type": "string"
    }
  },
  "required": ["deviceType"],
  "oneOf": [
    {
      "properties": {
        "deviceType": { "const": "smartphone" }
      },
      "$ref": "#/definitions/smartphone"
    },
    {
      "properties": {
        "deviceType": { "const": "laptop" }
      },
      "$ref": "#/definitions/laptop"
    }
  ],
  "definitions": {
    "smartphone": {
      "$id": "https://example.com/smartphone.schema.json",
      "$schema": "https://json-schema.org/draft/2020-12/schema",
      "type": "object",
      "properties": {
        "brand": {
          "type": "string"
        },
        "model": {
          "type": "string"
        },
        "screenSize": {
          "type": "number"
        }
      },
      "required": ["brand", "model", "screenSize"]
    },
	"laptop": {
      "$id": "https://example.com/laptop.schema.json",
      "$schema": "https://json-schema.org/draft/2020-12/schema",
      "type": "object",
      "properties": {
        "brand": {
          "type": "string"
        },
        "model": {
          "type": "string"
        },
        "processor": {
          "type": "string"
        },
        "ramSize": {
          "type": "number"
        }
      },
      "required": ["brand", "model", "processor", "ramSize"]
    }
  }
}
`

func RandomDevicetype() Rand {
	return Or(WithAnyOf(
		RandomLaptop(),
		RandomSmartphone(),
	))
}

func RandomLaptop() Rand {
	return Object(WithAdditionalFields(), WithFields(
		Field("deviceType", Const(`"laptop"`), IsRequired()),
		Field("brand", String(), IsRequired()),
		Field("model", String(), IsRequired()),
		Field("processor", String(), IsRequired()),
		Field("ramSize", Number(), IsRequired()),
	))
}

func RandomSmartphone() Rand {
	return Object(WithAdditionalFields(), WithFields(
		Field("deviceType", Const(`"smartphone"`), IsRequired()),
		Field("brand", String(), IsRequired()),
		Field("model", String(), IsRequired()),
		Field("screenSize", Number(), IsRequired()),
	))
}
