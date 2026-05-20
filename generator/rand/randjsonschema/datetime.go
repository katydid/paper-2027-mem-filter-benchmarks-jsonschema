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

	"github.com/katydid/validator-jsonschema-benchmarks/generator/rand"
)

type randDatetime struct{}

func DateTime() Rand {
	return &randDatetime{}
}

func (o *randDatetime) Right(r rand.Rand) string {
	return strconv.Quote(randValidDateTimeString(r))
}

func (o *randDatetime) Wrong(r rand.Rand) string {
	return strconv.Quote(randId(r, 20))
}

func randValidDateString(r rand.Rand) string {
	return fmt.Sprintf("%d%d%d%d-%d%d-%d%d", 2, 0, r.Intn(10), r.Intn(10), r.Intn(2), r.Intn(2)+1, r.Intn(2), r.Intn(7)+1)
}

func randValidTimeString(r rand.Rand) string {
	return fmt.Sprintf("%d%d:%d%d:%d%dZ", r.Intn(2), r.Intn(4), r.Intn(2), r.Intn(2)+1, r.Intn(2), r.Intn(7)+1)
}

func randValidDateTimeString(r rand.Rand) string {
	return randValidDateString(r) + "T" + randValidTimeString(r)
}
