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

package randjsonschema

import (
	"strconv"

	"github.com/katydid/validator-jsonschema-benchmarks/generator/rand"
)

type randURI struct{}

// [scheme:][//[userinfo@]host][/]path[?query][#fragment]
func URI() Rand {
	return &randURI{}
}

func (o *randURI) Right(r rand.Rand) string {
	scheme := "http://"
	if r.Intn(2) == 0 {
		scheme = "https://"
	}
	path := randID(r) + "." + randExtension(r)
	fragment := ""
	if r.Intn(4) == 0 {
		fragment = "#" + randID(r)
	}
	return strconv.Quote(scheme + path + fragment)
}

func (o *randURI) Wrong(r rand.Rand) string {
	return strconv.Quote(randId(r, 20))
}

func randID(r rand.Rand) string {
	l := r.Intn(20) + 1
	rs := make([]rune, l)
	for i := 0; i < l; i++ {
		switch r.Intn(2) {
		case 0:
			rs[i] = 'A' + rune(r.Intn(26))
		case 1:
			rs[i] = 'a' + rune(r.Intn(26))
		}
	}
	return string(rs)
}
