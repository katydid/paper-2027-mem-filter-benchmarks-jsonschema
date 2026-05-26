// Copied from https://github.com/ajv-validator/ajv/blob/master/spec/tests/schemas/advanced.json
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

package zschema_advanced

const SchemaAdvanced = `
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "/": {"$ref": "#/definitions/entry"}
  },
  "patternProperties": {
    "^(/[^/]+)+$": {"$ref": "#/definitions/entry"}
  },
  "additionalProperties": false,
  "required": ["/"],
  "definitions": {
    "entry": {
      "$schema": "http://json-schema.org/draft-07/schema#",
      "description": "schema for an fstab entry",
      "type": "object",
      "required": ["storage"],
      "properties": {
        "storage": {
          "type": "object",
          "oneOf": [
            {"$ref": "#/definitions/entry/definitions/diskDevice"},
            {"$ref": "#/definitions/entry/definitions/diskUUID"},
            {"$ref": "#/definitions/entry/definitions/nfs"},
            {"$ref": "#/definitions/entry/definitions/tmpfs"}
          ]
        },
        "fstype": {
          "enum": ["ext3", "ext4", "btrfs"]
        },
        "options": {
          "type": "array",
          "minItems": 1,
          "items": {"type": "string"},
          "uniqueItems": true
        },
        "readonly": {"type": "boolean"}
      },
      "definitions": {
        "diskDevice": {
          "properties": {
            "type": {"enum": ["disk"]},
            "device": {
              "type": "string",
              "pattern": "^/dev/[^/]+(/[^/]+)*$"
            }
          },
          "required": ["type", "device"],
          "additionalProperties": false
        },
        "diskUUID": {
          "properties": {
            "type": {"enum": ["disk"]},
            "label": {
              "type": "string",
              "pattern": "^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$"
            }
          },
          "required": ["type", "label"],
          "additionalProperties": false
        },
        "nfs": {
          "properties": {
            "type": {"enum": ["nfs"]},
            "remotePath": {
              "type": "string",
              "pattern": "^(/[^/]+)+$"
            },
            "server": {
              "type": "string",
              "anyOf": [{"format": "hostname"}, {"format": "ipv4"}, {"format": "ipv6"}]
            }
          },
          "required": ["type", "server", "remotePath"],
          "additionalProperties": false
        },
        "tmpfs": {
          "properties": {
            "type": {"enum": ["tmpfs"]},
            "sizeInMB": {
              "type": "integer",
              "minimum": 16,
              "maximum": 512
            }
          },
          "required": ["type", "sizeInMB"],
          "additionalProperties": false
        }
      }
    }
  }
}
`

const SchemaAdvancedNoUniqueItems = `
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "/": {"$ref": "#/definitions/entry"}
  },
  "patternProperties": {
    "^(/[^/]+)+$": {"$ref": "#/definitions/entry"}
  },
  "additionalProperties": false,
  "required": ["/"],
  "definitions": {
    "entry": {
      "$schema": "http://json-schema.org/draft-07/schema#",
      "description": "schema for an fstab entry",
      "type": "object",
      "required": ["storage"],
      "properties": {
        "storage": {
          "type": "object",
          "oneOf": [
            {"$ref": "#/definitions/entry/definitions/diskDevice"},
            {"$ref": "#/definitions/entry/definitions/diskUUID"},
            {"$ref": "#/definitions/entry/definitions/nfs"},
            {"$ref": "#/definitions/entry/definitions/tmpfs"}
          ]
        },
        "fstype": {
          "enum": ["ext3", "ext4", "btrfs"]
        },
        "options": {
          "type": "array",
          "minItems": 1,
          "items": {"type": "string"}
        },
        "readonly": {"type": "boolean"}
      },
      "definitions": {
        "diskDevice": {
          "properties": {
            "type": {"enum": ["disk"]},
            "device": {
              "type": "string",
              "pattern": "^/dev/[^/]+(/[^/]+)*$"
            }
          },
          "required": ["type", "device"],
          "additionalProperties": false
        },
        "diskUUID": {
          "properties": {
            "type": {"enum": ["disk"]},
            "label": {
              "type": "string",
              "pattern": "^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$"
            }
          },
          "required": ["type", "label"],
          "additionalProperties": false
        },
        "nfs": {
          "properties": {
            "type": {"enum": ["nfs"]},
            "remotePath": {
              "type": "string",
              "pattern": "^(/[^/]+)+$"
            },
            "server": {
              "type": "string",
              "anyOf": [{"format": "hostname"}, {"format": "ipv4"}, {"format": "ipv6"}]
            }
          },
          "required": ["type", "server", "remotePath"],
          "additionalProperties": false
        },
        "tmpfs": {
          "properties": {
            "type": {"enum": ["tmpfs"]},
            "sizeInMB": {
              "type": "integer",
              "minimum": 16,
              "maximum": 512
            }
          },
          "required": ["type", "sizeInMB"],
          "additionalProperties": false
        }
      }
    }
  }
}
`

//     "tests": [
//       {
//         "description": "valid object from z-schema benchmark",
//         "data": {
//           "/": {
//             "storage": {
//               "type": "disk",
//               "device": "/dev/sda1"
//             },
//             "fstype": "btrfs",
//             "readonly": true
//           },
//           "/var": {
//             "storage": {
//               "type": "disk",
//               "label": "8f3ba6f4-5c70-46ec-83af-0d5434953e5f"
//             },
//             "fstype": "ext4",
//             "options": ["nosuid"]
//           },
//           "/tmp": {
//             "storage": {
//               "type": "tmpfs",
//               "sizeInMB": 64
//             }
//           },
//           "/var/www": {
//             "storage": {
//               "type": "nfs",
//               "server": "my.nfs.server",
//               "remotePath": "/exports/mypath"
//             }
//           }
//         },
//         "valid": true
//       },
//       {
//         "description": "not object",
//         "data": 1,
//         "valid": false
//       },
//       {
//         "description": "root only is valid",
//         "data": {
//           "/": {
//             "storage": {
//               "type": "disk",
//               "device": "/dev/sda1"
//             },
//             "fstype": "btrfs",
//             "readonly": true
//           }
//         },
//         "valid": true
//       },
//       {
//         "description": "missing root entry",
//         "data": {
//           "no root/": {
//             "storage": {
//               "type": "disk",
//               "device": "/dev/sda1"
//             },
//             "fstype": "btrfs",
//             "readonly": true
//           }
//         },
//         "valid": false
//       },
//       {
//         "description": "invalid entry key",
//         "data": {
//           "/": {
//             "storage": {
//               "type": "disk",
//               "device": "/dev/sda1"
//             },
//             "fstype": "btrfs",
//             "readonly": true
//           },
//           "invalid/var": {
//             "storage": {
//               "type": "disk",
//               "label": "8f3ba6f4-5c70-46ec-83af-0d5434953e5f"
//             },
//             "fstype": "ext4",
//             "options": ["nosuid"]
//           }
//         },
//         "valid": false
//       },
//       {
//         "description": "missing storage in entry",
//         "data": {
//           "/": {
//             "fstype": "btrfs",
//             "readonly": true
//           }
//         },
//         "valid": false
//       },
//       {
//         "description": "missing storage type",
//         "data": {
//           "/": {
//             "storage": {
//               "device": "/dev/sda1"
//             },
//             "fstype": "btrfs",
//             "readonly": true
//           }
//         },
//         "valid": false
//       },
//       {
//         "description": "storage type should be a string",
//         "data": {
//           "/": {
//             "storage": {
//               "type": null,
//               "device": "/dev/sda1"
//             },
//             "fstype": "btrfs",
//             "readonly": true
//           }
//         },
//         "valid": false
//       },
//       {
//         "description": "storage device should match pattern",
//         "data": {
//           "/": {
//             "storage": {
//               "type": null,
//               "device": "invalid/dev/sda1"
//             },
//             "fstype": "btrfs",
//             "readonly": true
//           }
//         },
//         "valid": false
//       }
//     ]
//   }
// ]
