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

package jsck_medium

import (
	. "github.com/katydid/validator-jsonschema-benchmarks/generator/rand/randjsonschema"
)

func Medium() Rand {
	return Object(WithFields(
		Field("api_server", APIServer(), IsRequired()),
		Field("transport", Transport(), IsRequired()),
		Field("storage", Storage(), IsRequired()),
		Field("chain", Chain(), IsRequired()),
	))
}

func APIServer() Rand {
	return Object(WithFields(
		Field("url", URI(), IsRequired()),
		Field("host", String(), IsRequired()),
		Field("port", Integer(WithMinimum(1000)), IsRequired()),
	))
}

func Transport() Rand {
	return Object(WithFields(
		Field("server", String(), IsRequired()),
		Field("options", Object(WithAdditionalFields())),
		Field("queues", Object(WithAdditionalFields(), WithFields(
			Field("blocking_timeout", Integer(WithMinimum(0))),
		))),
	))
}

func Storage() Rand {
	return Object(WithAdditionalFields(), WithFields(
		Field("server", String(), IsRequired()),
		Field("database", String(), IsRequired()),
		Field("user", String(), IsRequired()),
		Field("options", Object(WithAdditionalFields())),
	))
}

func Chain() Rand {
	return Object(WithFields(
		Field("api_key_id", String(), IsRequired()),
		Field("api_key_secret", String(), IsRequired()),
	))
}
