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

package ajv_complex

import (
	"github.com/katydid/validator-jsonschema-benchmarks/generator/rand"
	. "github.com/katydid/validator-jsonschema-benchmarks/generator/rand/randjsonschema"
)

func Input() Rand {
	return Object(WithFields(
		Field("index", Integer(WithMinimum(0)), IsRequired()),
		Field("output", Output(), IsRequired()),
		Field("sig_hash", Hex()),
		Field("script_sig", Hex(), IsRequired()),
		Field("signatures", Signatures()),
	))
}

//	"signatures": {
//		"type": "object",
//		"description": "A dictionary of signatures.  Keys represent keypair names",
//		"minProperties": 1,
//		"maxProperties": 3,
//		"additionalProperties": {"$ref": "#signature"}
//	}
//
//	"signature": {
//		"$anchor": "signature",
//		"allOf": [
//			{"$ref": "#hex"},
//			{
//				"minLength": 128,
//				"maxLength": 128
//			}
//		]
//	}
//
// The tests include an example:
//
//	"signatures": {
//	   "primary": "3046022100be69797cf5d784412b1258256eb657c191a04893479dfa2ae5c7f2088c7adbe0022100e6b000bd633b286ed1b9bc7682fe753d9fdad61fbe5da2a7",
//	   "cosigner": "a2ad5ebf16dadf9d357ef2867cb9b1de682b336db000b6e0012200ebda7c8802f7c5ea2afd97439840a191c756be6528521b214487d5fc79796eb00122064037"
//	},
func Signatures() Rand {
	return &randSignatures{}
}

type randSignatures struct{}

func (o *randSignatures) Right(r rand.Rand) string {
	return Object(WithFields(
		Field("primary", Signature(), IsRequired()),
		Field("cosigner", Signature(), IsRequired()),
	)).Right(r)
}

func (o *randSignatures) Wrong(r rand.Rand) string {
	switch r.Intn(2) {
	case 0: // too many properties
		return Object(WithFields(
			Field("primary", Signature(), IsRequired()),
			Field("cosigner", Signature(), IsRequired()),
			Field("primary2", Signature(), IsRequired()),
			Field("cosigner2", Signature(), IsRequired()),
		)).Right(r)
	case 1: // too few properties
		return "{}"
	}
	panic("unreachable")
}
