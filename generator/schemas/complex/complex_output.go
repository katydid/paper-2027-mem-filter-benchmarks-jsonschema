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

package complex

import (
	"github.com/katydid/validator-jsonschema-benchmarks/generator/rand"
	. "github.com/katydid/validator-jsonschema-benchmarks/generator/rand/randjsonschema"
)

func Output() Rand {
	return Object(WithFields(
		Field("hash", TxID(), IsRequired()),
		Field("index", Integer(WithMinimum(0)), IsRequired()),
		Field("value", Integer(), IsRequired()),
		Field("script", Script(), IsRequired()),
		Field("address", Address()),
		Field("metadata", OutputMetadata()),
	))
}

func Script() Rand {
	return Object(WithFields(
		Field("type", ScriptType()),
		Field("asm", String()),
	))
}

func OutputMetadata() Rand {
	return Or(
		WithRight(
			Object(WithAdditionalFields(), WithFields(
				Field("wallet_path", String(), IsRequired()),
				Field("public_seeds", PublicSeeds(), IsRequired()),
			)),
			Object(WithAdditionalFields(), WithFields(
				Field("public_seeds", PublicSeeds(), IsRequired()),
			)),
		),
		WithWrong(
			Object(WithAdditionalFields(), WithFields(
				Field("wallet_path", String(), IsRequired()),
			)),
			Object(WithAdditionalFields(), WithFields(
				Field("wallet_path", Integer(), IsRequired()),
				Field("public_seeds", PublicSeeds(), IsRequired()),
			)),
		),
	)
}

//	"public_seeds": {
//		"type": "object",
//		"minProperties": 1,
//		"maxProperties": 3,
//		"additionalProperties": {
//			"anyOf": [{"$ref": "#base58"}, {"$ref": "#hex"}]
//		}
//	}
//
// test cases include:
//
//	"public_seeds": {
//	     "primary": "xpub661MyMwAqRbcFtXgS5sYJABqqG9YLmC4Q1Rdap9gSE8NqtwybGhePY2gZ29ESFjqJoCu1Rupje8YtGqsefD265TMg7usUDFdp6W1EGMcet8",
//	     "cosigner": "xpub661MyMwAqRbcFtXgS5sYJABqqG9YLmC4Q1Rdap9gSE8NqtwybGhePY2gZ29ESFjqJoCu1Rupje8YtGqsefD265TMg7usUDFdp6W1EGMcet8"
//	}
func PublicSeeds() Rand {
	return &randPublicSeeds{}
}

type randPublicSeeds struct{}

func (o *randPublicSeeds) Right(r rand.Rand) string {
	return Object(WithFields(
		Field("primary", Base58(), IsRequired()),
		Field("cosigner", Hex(), IsRequired()),
	)).Right(r)
}

func (o *randPublicSeeds) Wrong(r rand.Rand) string {
	switch r.Intn(2) {
	case 0: // too many properties
		return Object(WithFields(
			Field("primary", Base58(), IsRequired()),
			Field("cosigner", Hex(), IsRequired()),
			Field("primary2", Base58(), IsRequired()),
			Field("cosigner2", Hex(), IsRequired()),
		)).Right(r)
	case 1: // too few properties
		return "{}"
	}
	panic("unreachable")
}
