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
	. "github.com/katydid/validator-jsonschema-benchmarks/generator/rand/randjsonschema"
)

func Complex() Rand {
	return ArrayOf(Transaction(), WithMinItems(1))
}

func Transaction() Rand {
	return Object(WithFields(
		Field("metadata", TransactionMetadata(), IsRequired()),
		Field("version", Integer()),
		Field("lock_time", Integer()),
		Field("hash", TxID(), IsRequired()),
		Field("inputs", ArrayOf(Input(), WithMinItems(1)), IsRequired()),
		Field("outputs", ArrayOf(Output(), WithMinItems(1)), IsRequired()),
	))
}

func TransactionMetadata() Rand {
	return Object(WithFields(
		Field("amount", Integer(), IsRequired()),
		Field("fee", Integer(WithMultipleOf(10000)), IsRequired()),
		Field("status", Status()),
		Field("confirmations", Integer(WithMinimum(0))),
		Field("block_time", Integer()),
	))
}
