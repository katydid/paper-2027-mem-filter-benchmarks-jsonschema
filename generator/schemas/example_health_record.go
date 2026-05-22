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

// https://json-schema.org/learn/json-schema-examples#health-record
const SchemaJSONSchemaExampleHealthRecord = `
{
  "$id": "https://example.com/health-record.schema.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "description": "Schema for representing a health record",
  "type": "object",
  "required": ["patientName", "dateOfBirth", "bloodType"],
  "properties": {
    "patientName": {
      "type": "string"
    },
    "dateOfBirth": {
      "type": "string",
      "format": "date"
    },
    "bloodType": {
      "type": "string"
    },
    "allergies": {
      "type": "array",
      "items": {
        "type": "string"
      }
    },
    "conditions": {
      "type": "array",
      "items": {
        "type": "string"
      }
    },
    "medications": {
      "type": "array",
      "items": {
        "type": "string"
      }
    },
    "emergencyContact": {
      "$ref": "#/definitions/user-profile"
    }
  },
  "definitions": {
    "user-profile": {
      "description": "A representation of a user profile",
      "type": "object",
      "required": ["username", "email"],
      "properties": {
        "username": {
          "type": "string"
        },
        "email": {
          "type": "string",
          "format": "email"
        },
        "fullName": {
          "type": "string"
        },
        "age": {
          "type": "integer",
          "minimum": 0
        },
        "location": {
          "type": "string"
        },
        "interests": {
          "type": "array",
          "items": {
            "type": "string"
          }
        }
      }
    }
  }
}
`

func RandomHealthRecord() Rand {
	return Object(WithAdditionalFields(), WithFields(
		Field("patientName", String(), IsRequired()),
		Field("dateOfBirth", Date(), IsRequired()),
		Field("bloodType", String(), IsRequired()),
		Field("allergies", ArrayOf(String())),
		Field("conditions", ArrayOf(String())),
		Field("medications", ArrayOf(String())),
		Field("emergencyContact", RandomUserProfile()),
	))
}
