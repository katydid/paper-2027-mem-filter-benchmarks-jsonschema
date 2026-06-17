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
	"fmt"
	"strconv"
	"strings"

	"github.com/katydid/validator-jsonschema-benchmarks/generator/rand"
)

type randIPv4 struct{}

func IPv4() Rand {
	return &randIPv4{}
}

func (o *randIPv4) Right(r rand.Rand) string {
	return fmt.Sprintf(`"%d.%d.%d.%d"`, r.Intn(126)+1, r.Intn(127), r.Intn(127), r.Intn(127))
}

func (o *randIPv4) Wrong(r rand.Rand) string {
	s := strconv.Quote(randId(r, 20))
	return strings.Replace(s, ".", "@", -1)
}
