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

// https://json-schema.org/learn/json-schema-examples#movie
const SchemaJSONSchemaExampleMovie = `
{
  "$id": "https://example.com/movie.schema.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "description": "A representation of a movie",
  "type": "object",
  "required": ["title", "director", "releaseDate"],
  "properties": {
    "title": {
      "type": "string"
    },
    "director": {
      "type": "string"
    },
    "releaseDate": {
      "type": "string",
      "format": "date"
    },
    "genre": {
      "type": "string",
      "enum": ["Action", "Comedy", "Drama", "Science Fiction"]
    },
    "duration": {
      "type": "string"
    },
    "cast": {
      "type": "array",
      "items": {
        "type": "string"
      },
      "additionalItems": false
    }
  }
}
`

func RandomMovie() Rand {
	return Object(WithAdditionalFields(), WithFields(
		Field("title", String(), IsRequired()),
		Field("director", String(), IsRequired()),
		Field("releaseDate", Date(), IsRequired()),
		Field("genre", Or(WithAnyOf(Const(`"Action"`), Const(`"Comedy"`), Const(`"Drama"`), Const(`"Science Fiction"`)))),
		Field("duration", String()),
		Field("cast", ArrayOf(String())),
	))
}
