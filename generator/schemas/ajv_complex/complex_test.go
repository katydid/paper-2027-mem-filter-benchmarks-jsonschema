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
	"strconv"
	"testing"

	"github.com/dlclark/regexp2/v2"
	"github.com/katydid/validator-jsonschema-benchmarks/generator/rand"
)

func TestRandomAddress(t *testing.T) {
	r := rand.NewRand()
	for i := 0; i < 1000; i++ {
		address := Address().Right(r)
		uaddress, err := strconv.Unquote(address)
		if err != nil {
			t.Fatal(err)
		}
		if len(uaddress) != 34 {
			t.Fatalf("expected address length 34, but got %d", len(uaddress))
		}
		rx, err := regexp2.Compile("^[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]+$")
		if err != nil {
			t.Fatal(err)
		}
		if m, err := rx.MatchString(uaddress); err != nil || !m {
			t.Fatalf("address did not match %s: %v", uaddress, err)
		}
	}
}

func TestRandomBase58(t *testing.T) {
	r := rand.NewRand()
	for i := 0; i < 1000; i++ {
		base := Base58().Right(r)
		ubase, err := strconv.Unquote(base)
		if err != nil {
			t.Fatal(err)
		}
		rx, err := regexp2.Compile("^[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]+$")
		if err != nil {
			t.Fatal(err)
		}
		if m, err := rx.MatchString(ubase); err != nil || !m {
			t.Fatalf("base did not match %s: %v", ubase, err)
		}
	}
}
