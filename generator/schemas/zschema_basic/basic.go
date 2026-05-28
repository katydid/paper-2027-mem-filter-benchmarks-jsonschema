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

package zschema_basic

import (
	. "github.com/katydid/validator-jsonschema-benchmarks/generator/rand/randjsonschema"
)

func ProductSet() Rand {
	return ArrayOf(Product())
}

func Product() Rand {
	return Object(WithAdditionalFields(),
		WithAlwaysRightFields(
			Field("warhouseLocation", Any()),
		),
		WithFields(
			Field("id", Number(), IsRequired()),
			Field("name", String(), IsRequired()),
			Field("price", Number(WithExclusiveMinimum(0)), IsRequired()),
			Field("tags", ArrayOf(String(), WithMinItems(1), WithUniqueItems())),
			Field("dimensions", Object(WithAdditionalFields(), WithFields(
				Field("length", Number(), IsRequired()),
				Field("width", Number(), IsRequired()),
				Field("height", Number(), IsRequired()),
			))),
		),
	)
}

// exactly the same as the schema above, except we removed uniqueItems
func ProductSetrmUniqueItems() Rand {
	return ArrayOf(ProductrmUniqueItems())
}

func ProductrmUniqueItems() Rand {
	return Object(WithAdditionalFields(),
		WithAlwaysRightFields(
			Field("warhouseLocation", Any()),
		),
		WithFields(
			Field("id", Number(), IsRequired()),
			Field("name", String(), IsRequired()),
			Field("price", Number(WithExclusiveMinimum(0)), IsRequired()),
			Field("tags", ArrayOf(String(), WithMinItems(1))),
			Field("dimensions", Object(WithAdditionalFields(), WithFields(
				Field("length", Number(), IsRequired()),
				Field("width", Number(), IsRequired()),
				Field("height", Number(), IsRequired()),
			))),
		),
	)
}
