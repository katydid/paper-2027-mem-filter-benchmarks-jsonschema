// Copied from https://github.com/ajv-validator/ajv/blob/master/spec/tests/schemas/medium.json
// The MIT License (MIT)

// Copyright (c) 2015-2021 Evgeny Poberezkin

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// They originally copied it from https://github.com/pandastrike/jsck
// The MIT License (MIT)

// Copyright (c) 2013 Matthew King

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package jsck_medium

const SchemaMedium = `
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "description": "A moderately complex schema with some nesting and value constraints",
  "type": "object",
  "additionalProperties": false,
  "required": ["api_server", "transport", "storage", "chain"],
  "properties": {
    "api_server": {
      "description": "Settings for the HTTP API server",
      "type": "object",
      "additionalProperties": false,
      "required": ["url", "host", "port"],
      "properties": {
        "url": {
          "type": "string",
          "format": "uri"
        },
        "host": {
          "type": "string"
        },
        "port": {
          "type": "integer",
          "minimum": 1000
        }
      }
    },
    "transport": {
      "description": "Settings for the Redis tranport",
      "additionalProperties": false,
      "required": ["server"],
      "properties": {
        "server": {
          "type": "string"
        },
        "options": {
          "type": "object"
        },
        "queues": {
          "properties": {
            "blocking_timeout": {
              "type": "integer",
              "minimum": 0
            }
          }
        }
      }
    },
    "storage": {
      "description": "Settings for the PostgreSQL storage",
      "required": ["server", "database", "user"],
      "properties": {
        "server": {
          "type": "string"
        },
        "database": {
          "type": "string"
        },
        "user": {
          "type": "string"
        },
        "options": {
          "type": "object"
        }
      }
    },
    "chain": {
      "description": "Settings for the Chain.com client",
      "required": ["api_key_id", "api_key_secret"],
      "properties": {
        "api_key_id": {
          "type": "string"
        },
        "api_key_secret": {
          "type": "string"
        }
      }
    }
  }
}
`

// "tests": [
//   {
//     "description": "valid object from jsck benchmark",
//     "data": {
//       "api_server": {
//         "url": "http://example.com:8998",
//         "host": "example.com",
//         "port": 8998
//       },
//       "transport": {
//         "server": "127.0.0.1:6381",
//         "queues": {
//           "blocking_timeout": 0
//         }
//       },
//       "storage": {
//         "server": "127.0.0.1:5432",
//         "database": "thingy-test",
//         "user": "thingy-test",
//         "password": "password"
//       },
//       "chain": {
//         "api_key_id": "cafebabe",
//         "api_key_secret": "babecafe"
//       }
//     },
//     "valid": true
//   },
//   {
//     "description": "not object",
//     "data": 1,
//     "valid": false
//   }
// ]
