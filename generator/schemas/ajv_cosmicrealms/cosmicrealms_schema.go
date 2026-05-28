// Copied from https://github.com/ajv-validator/ajv/blob/master/spec/tests/schemas/cosmicrealms.json
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

package ajv_cosmicrealms

const SchemaCosmicRealms = `
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "name": "test",
  "type": "object",
  "additionalProperties": false,
  "required": [
    "fullName",
    "age",
    "zip",
    "married",
    "dozen",
    "dozenOrBakersDozen",
    "favoriteEvenNumber",
    "topThreeFavoriteColors",
    "favoriteSingleDigitWholeNumbers",
    "favoriteFiveLetterWord",
    "emailAddresses",
    "ipAddresses"
  ],
  "properties": {
    "fullName": {"type": "string"},
    "age": {"type": "integer", "minimum": 0},
    "optionalItem": {"type": "string"},
    "state": {"type": "string"},
    "city": {"type": "string"},
    "zip": {"type": "integer", "minimum": 0, "maximum": 99999},
    "married": {"type": "boolean"},
    "dozen": {"type": "integer", "minimum": 12, "maximum": 12},
    "dozenOrBakersDozen": {"type": "integer", "minimum": 12, "maximum": 13},
    "favoriteEvenNumber": {"type": "integer", "multipleOf": 2},
    "topThreeFavoriteColors": {
      "type": "array",
      "minItems": 3,
      "maxItems": 3,
      "uniqueItems": true,
      "items": {"type": "string"}
    },
    "favoriteSingleDigitWholeNumbers": {
      "type": "array",
      "minItems": 1,
      "maxItems": 10,
      "uniqueItems": true,
      "items": {"type": "integer", "minimum": 0, "maximum": 9}
    },
    "favoriteFiveLetterWord": {
      "type": "string",
      "minLength": 5,
      "maxLength": 5
    },
    "emailAddresses": {
      "type": "array",
      "minItems": 1,
      "uniqueItems": true,
      "items": {"type": "string", "format": "email"}
    },
    "ipAddresses": {
      "type": "array",
      "uniqueItems": true,
      "items": {"type": "string", "format": "ipv4"}
    }
  }
}   
`

const SchemaCosmicRealmsrmUniqueItems = `
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "name": "test",
  "type": "object",
  "additionalProperties": false,
  "required": [
    "fullName",
    "age",
    "zip",
    "married",
    "dozen",
    "dozenOrBakersDozen",
    "favoriteEvenNumber",
    "topThreeFavoriteColors",
    "favoriteSingleDigitWholeNumbers",
    "favoriteFiveLetterWord",
    "emailAddresses",
    "ipAddresses"
  ],
  "properties": {
    "fullName": {"type": "string"},
    "age": {"type": "integer", "minimum": 0},
    "optionalItem": {"type": "string"},
    "state": {"type": "string"},
    "city": {"type": "string"},
    "zip": {"type": "integer", "minimum": 0, "maximum": 99999},
    "married": {"type": "boolean"},
    "dozen": {"type": "integer", "minimum": 12, "maximum": 12},
    "dozenOrBakersDozen": {"type": "integer", "minimum": 12, "maximum": 13},
    "favoriteEvenNumber": {"type": "integer", "multipleOf": 2},
    "topThreeFavoriteColors": {
      "type": "array",
      "minItems": 3,
      "maxItems": 3,
      "items": {"type": "string"}
    },
    "favoriteSingleDigitWholeNumbers": {
      "type": "array",
      "minItems": 1,
      "maxItems": 10,
      "items": {"type": "integer", "minimum": 0, "maximum": 9}
    },
    "favoriteFiveLetterWord": {
      "type": "string",
      "minLength": 5,
      "maxLength": 5
    },
    "emailAddresses": {
      "type": "array",
      "minItems": 1,
      "items": {"type": "string", "format": "email"}
    },
    "ipAddresses": {
      "type": "array",
      "items": {"type": "string", "format": "ipv4"}
    }
  }
}   
`

// "tests": [
//       {
//         "description": "valid data from cosmicrealms benchmark",
//         "data": {
//           "fullName": "John Smith",
//           "state": "CA",
//           "city": "Los Angeles",
//           "favoriteFiveLetterWord": "hello",
//           "emailAddresses": [
//             "NRorsfCYtvB5bKAf1jZMu1GAJzAhhg5lEvh@inTqnn.net",
//             "6tjWtYxjaan2Ivm5QZVhKxImKawRCA6gcqtMEwV1@bB01pCtIBY0F.org",
//             "j68UnHfrHiKwpAm8iYokoMuRTpWUj8bfxspusNFK@COoWeMZL.edu",
//             "qlnrIsYSWCGUQW6f8HL@UBOqUYQQzugVL.uk"
//           ],
//           "dozen": 12,
//           "dozenOrBakersDozen": 13,
//           "favoriteEvenNumber": 24,
//           "married": true,
//           "age": 17,
//           "zip": 65794,
//           "topThreeFavoriteColors": ["blue", "black", "yellow"],
//           "favoriteSingleDigitWholeNumbers": [2, 1, 3, 9],
//           "ipAddresses": ["225.234.40.3", "96.216.243.54", "18.126.145.83", "196.17.191.239"]
//         },
//         "valid": true
//       },
//       {
//         "description": "invalid data",
//         "data": {
//           "state": null,
//           "city": 90912,
//           "zip": [null],
//           "married": "married",
//           "dozen": 90912,
//           "dozenOrBakersDozen": null,
//           "favoriteEvenNumber": -1294145,
//           "emailAddresses": [],
//           "topThreeFavoriteColors": [
//             null,
//             null,
//             0.7925170068027211,
//             1.2478632478632479,
//             1.173913043478261,
//             0.4472049689440994
//           ],
//           "favoriteSingleDigitWholeNumbers": [],
//           "favoriteFiveLetterWord": "more than five letters",
//           "ipAddresses": [
//             "55.335.74.758",
//             "191.266.92.805",
//             "193.388.390.250",
//             "269.375.318.49",
//             "120.268.59.140"
//           ]
//         },
//         "valid": false
//       }
//     ]
//   }
// ]
