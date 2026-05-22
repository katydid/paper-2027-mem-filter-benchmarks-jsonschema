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

// https://json-schema.org/learn/json-schema-examples#user-profile
const SchemaJSONSchemaExampleUserProfile = `{
  "$id": "https://example.com/user-profile.schema.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
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
}`

func RandomUserProfile() Rand {
	return Object(WithAdditionalFields(), WithFields(
		Field("username", String(), IsRequired()),
		Field("email", Email(), IsRequired()),
		Field("fullName", String()),
		Field("age", Integer(WithMinimum(0))),
		Field("location", String()),
		Field("interests", ArrayOf(String())),
	))
}
