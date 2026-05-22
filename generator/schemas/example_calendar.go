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

// https://json-schema.org/learn/json-schema-examples#calendar
const SchemaJSONSchemaExampleCalendar = `
{
  "$id": "https://example.com/calendar.schema.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "description": "A representation of an event",
  "type": "object",
  "required": [ "startDate", "summary" ],
  "properties": {
    "startDate": {
      "type": "string",
      "description": "Event starting time"
    },
    "endDate": {
      "type": "string",
      "description": "Event ending time"
    },
    "summary": {
      "type": "string"
    },
    "location": {
      "type": "string"
    },
    "url": {
      "type": "string"
    },
    "duration": {
      "type": "string",
      "description": "Event duration"
    },
    "recurrenceDate": {
      "type": "string",
      "description": "Recurrence date"
    },
    "recurrenceRule": {
      "type": "string",
      "description": "Recurrence rule"
    },
    "category": {
      "type": "string"
    },
    "description": {
      "type": "string"
    },
    "geo": {
      "$ref": "#/definitions/geographical-location"
    }
  },
  "definitions": {
    "geographical-location": {
      "title": "Longitude and Latitude Values",
	  "description": "A geographical coordinate.",
      "required": [ "latitude", "longitude" ],
      "type": "object",
      "properties": {
        "latitude": {
          "type": "number",
          "minimum": -90,
          "maximum": 90
        },
        "longitude": {
          "type": "number",
          "minimum": -180,
          "maximum": 180
        }
      }
    }
  }
}
`

func RandomCalendar() Rand {
	return Object(WithAdditionalFields(), WithFields(
		Field("startDate", String(), IsRequired()),
		Field("endDate", String()),
		Field("summary", String(), IsRequired()),
		Field("location", String()),
		Field("url", String()),
		Field("duration", String()),
		Field("recurrenceDate", String()),
		Field("recurrenceRule", String()),
		Field("category", String()),
		Field("description", String()),
		Field("geo", RandomGeographicalLocation()),
	))
}
