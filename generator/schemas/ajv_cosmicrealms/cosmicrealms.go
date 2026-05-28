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

package ajv_cosmicrealms

import (
	. "github.com/katydid/validator-jsonschema-benchmarks/generator/rand/randjsonschema"
)

func CosmicRealms() Rand {
	return Object(WithFields(
		Field("fullName", String(), IsRequired()),
		Field("age", Integer(WithMinimum(0)), IsRequired()),
		Field("optionalItem", String()),
		Field("state", String()),
		Field("city", String()),
		Field("zip", Integer(WithMinimum(0), WithMaximum(99999)), IsRequired()),
		Field("married", Bool(), IsRequired()),
		Field("dozen", Const("12"), IsRequired()),
		Field("dozenOrBakersDozen", Or(WithAnyOf(Const("12"), Const("13"))), IsRequired()),
		Field("favoriteEvenNumber", Integer(WithMultipleOf(2)), IsRequired()),
		Field("topThreeFavoriteColors", ArrayOf(String(), WithMinItems(3), WithMaxItems(3), WithUniqueItems()), IsRequired()),
		Field("favoriteSingleDigitWholeNumbers", ArrayOf(Integer(WithMinimum(0), WithMaximum(9)), WithMinItems(1), WithMaxItems(10), WithUniqueItems()), IsRequired()),
		Field("favoriteFiveLetterWord", String(WithMinLength(5), WithMaxLength(5)), IsRequired()),
		Field("emailAddresses", ArrayOf(Email(), WithMinItems(1), WithUniqueItems()), IsRequired()),
		Field("ipAddresses", ArrayOf(IPv4(), WithUniqueItems()), IsRequired()),
	))
}

func CosmicRealmsrmUniqueItems() Rand {
	return Object(WithFields(
		Field("fullName", String(), IsRequired()),
		Field("age", Integer(WithMinimum(0)), IsRequired()),
		Field("optionalItem", String()),
		Field("state", String()),
		Field("city", String()),
		Field("zip", Integer(WithMinimum(0), WithMaximum(99999)), IsRequired()),
		Field("married", Bool(), IsRequired()),
		Field("dozen", Const("12"), IsRequired()),
		Field("dozenOrBakersDozen", Or(WithAnyOf(Const("12"), Const("13"))), IsRequired()),
		Field("favoriteEvenNumber", Integer(WithMultipleOf(2)), IsRequired()),
		Field("topThreeFavoriteColors", ArrayOf(String(), WithMinItems(3), WithMaxItems(3)), IsRequired()),
		Field("favoriteSingleDigitWholeNumbers", ArrayOf(Integer(WithMinimum(0), WithMaximum(9)), WithMinItems(1), WithMaxItems(10)), IsRequired()),
		Field("favoriteFiveLetterWord", String(WithMinLength(5), WithMaxLength(5)), IsRequired()),
		Field("emailAddresses", ArrayOf(Email(), WithMinItems(1)), IsRequired()),
		Field("ipAddresses", ArrayOf(IPv4()), IsRequired()),
	))
}
