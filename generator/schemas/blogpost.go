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

// https://json-schema.org/learn/json-schema-examples#blog-post
const SchemaJSONSchemaExampleBlogPost = `{
  "$id": "https://example.com/blog-post.schema.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "description": "A representation of a blog post",
  "type": "object",
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
  },
  "required": ["title", "content", "author"],
  "properties": {
    "title": {
      "type": "string"
    },
    "content": {
      "type": "string"
    },
    "publishedDate": {
      "type": "string",
      "format": "date-time"
    },
    "author": {
      "$ref": "#/definitions/user-profile"
    },
    "tags": {
      "type": "array",
      "items": {
        "type": "string"
      }
    }
  }
}`

func RandomBlogPost() Rand {
	return Object(WithAdditionalFields(), WithFields(
		Field("title", String(), IsRequired()),
		Field("content", String(), IsRequired()),
		Field("publishedDate", DateTime()),
		Field("author", RandomUserProfile(), IsRequired()),
		Field("tags", ArrayOf(String())),
	))
}
