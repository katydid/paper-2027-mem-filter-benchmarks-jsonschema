// Copied from https://github.com/ajv-validator/ajv/blob/master/spec/tests/schemas/basic.json
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

// They originally copied it from https://github.com/zaggino/z-schema
// The MIT License (MIT)

// Copyright (c) 2014 Martin Zagora and other contributors
// https://github.com/zaggino/z-schema/graphs/contributors

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

package zschema_basic

const SchemaBasic = `
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Product set",
  "type": "array",
  "items": {
    "title": "Product",
    "type": "object",
    "properties": {
      "id": {
        "description": "The unique identifier for a product",
        "type": "number"
      },
      "name": {
        "type": "string"
      },
      "price": {
        "type": "number",
        "exclusiveMinimum": 0
      },
      "tags": {
        "type": "array",
        "items": {
          "type": "string"
        },
        "minItems": 1,
        "uniqueItems": true
      },
      "dimensions": {
        "type": "object",
        "properties": {
          "length": {"type": "number"},
          "width": {"type": "number"},
          "height": {"type": "number"}
        },
        "required": ["length", "width", "height"]
      },
      "warehouseLocation": {
        "description": "Coordinates of the warehouse with the product"
      }
    },
    "required": ["id", "name", "price"]
  }
}
`

// exactly the same as the schema above, except we removed uniqueItems
const SchemaBasicrmUniqueItems = `
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Product set",
  "type": "array",
  "items": {
    "title": "Product",
    "type": "object",
    "properties": {
      "id": {
        "description": "The unique identifier for a product",
        "type": "number"
      },
      "name": {
        "type": "string"
      },
      "price": {
        "type": "number",
        "exclusiveMinimum": 0
      },
      "tags": {
        "type": "array",
        "items": {
          "type": "string"
        },
        "minItems": 1
      },
      "dimensions": {
        "type": "object",
        "properties": {
          "length": {"type": "number"},
          "width": {"type": "number"},
          "height": {"type": "number"}
        },
        "required": ["length", "width", "height"]
      },
      "warehouseLocation": {
        "description": "Coordinates of the warehouse with the product"
      }
    },
    "required": ["id", "name", "price"]
  }
}
`

//     "tests": [
//       {
//         "description": "valid array from z-schema benchmark",
//         "data": [
//           {
//             "id": 2,
//             "name": "An ice sculpture",
//             "price": 12.5,
//             "tags": ["cold", "ice"],
//             "dimensions": {
//               "length": 7.0,
//               "width": 12.0,
//               "height": 9.5
//             },
//             "warehouseLocation": {
//               "latitude": -78.75,
//               "longitude": 20.4
//             }
//           },
//           {
//             "id": 3,
//             "name": "A blue mouse",
//             "price": 25.5,
//             "dimensions": {
//               "length": 3.1,
//               "width": 1.0,
//               "height": 1.0
//             },
//             "warehouseLocation": {
//               "latitude": 54.4,
//               "longitude": -32.7
//             }
//           }
//         ],
//         "valid": true
//       },
//       {
//         "description": "not array",
//         "data": 1,
//         "valid": false
//       },
//       {
//         "description": "array of not onjects",
//         "data": [1, 2, 3],
//         "valid": false
//       },
//       {
//         "description": "missing required properties",
//         "data": [{}],
//         "valid": false
//       },
//       {
//         "description": "required property of wrong type",
//         "data": [{"id": 1, "name": "product", "price": "not valid"}],
//         "valid": false
//       },
//       {
//         "description": "smallest valid product",
//         "data": [{"id": 1, "name": "product", "price": 100}],
//         "valid": true
//       },
//       {
//         "description": "tags should be array",
//         "data": [{"tags": {}, "id": 1, "name": "product", "price": 100}],
//         "valid": false
//       },
//       {
//         "description": "dimensions should be object",
//         "data": [{"dimensions": [], "id": 1, "name": "product", "price": 100}],
//         "valid": false
//       },
//       {
//         "description": "valid product with tag",
//         "data": [{"tags": ["product"], "id": 1, "name": "product", "price": 100}],
//         "valid": true
//       },
//       {
//         "description": "dimensions miss required properties",
//         "data": [
//           {
//             "dimensions": {},
//             "tags": ["product"],
//             "id": 1,
//             "name": "product",
//             "price": 100
//           }
//         ],
//         "valid": false
//       },
//       {
//         "description": "valid product with tag and dimensions",
//         "data": [
//           {
//             "dimensions": {"length": 7, "width": 12, "height": 9.5},
//             "tags": ["product"],
//             "id": 1,
//             "name": "product",
//             "price": 100
//           }
//         ],
//         "valid": true
//       }
//     ]
//   }
// ]
