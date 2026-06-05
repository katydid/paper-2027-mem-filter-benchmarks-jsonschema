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
	"strings"

	"github.com/katydid/validator-jsonschema-benchmarks/generator/rand"
)

type randEmail struct{}

func Email() Rand {
	return &randEmail{}
}

var extensions = []string{"com", "org", "co.uk", "co.za"}

func (o *randEmail) Right(r rand.Rand) string {
	return strconv.Quote(randId(r, 10) + "@" + randId(r, 10) + extensions[r.Intn(len(extensions))])
}

func (o *randEmail) Wrong(r rand.Rand) string {
	s := strconv.Quote(randId(r, 20))
	for strings.Contains(s, "@") {
		s = strconv.Quote(randId(r, 20))
	}
	return s
}
