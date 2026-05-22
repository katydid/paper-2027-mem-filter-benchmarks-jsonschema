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

// https://json-schema.org/learn/json-schema-examples#job-posting
const SchemaJSONSchemaExampleJobPosting = `
{
  "$id": "https://example.com/job-posting.schema.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "description": "A representation of a job posting",
  "type": "object",
  "required": ["title", "company", "location", "description"],
  "properties": {
    "title": {
      "type": "string"
    },
    "company": {
      "type": "string"
    },
    "location": {
      "type": "string"
    },
    "description": {
      "type": "string"
    },
    "employmentType": {
      "type": "string"
    },
    "salary": {
      "type": "number",
      "minimum": 0
    },
    "applicationDeadline": {
      "type": "string",
      "format": "date"
    }
  }
}
`

func RandomJobPosting() Rand {
	return Object(WithAdditionalFields(), WithFields(
		Field("title", String(), IsRequired()),
		Field("company", String(), IsRequired()),
		Field("location", String(), IsRequired()),
		Field("description", String(), IsRequired()),
		Field("employmentType", String()),
		Field("salary", Number(WithMinimum(0))),
		Field("applicationDeadline", Date()),
	))
}
