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

import "github.com/katydid/validator-jsonschema-benchmarks/generator/rand"

func randId(r rand.Rand, maxLen int) string {
	l := r.Intn(maxLen-1) + 1
	s := make([]rune, l)
	for i := range l {
		switch r.Intn(3) {
		case 0:
			s[i] = 'a' + rune(r.Intn(26))
		case 1:
			s[i] = 'A' + rune(r.Intn(26))
		case 2:
			s[i] = '0' + rune(r.Intn(10))
			if i == 0 {
				// first rune has to be a letter
				s[i] = 'a' + rune(r.Intn(26))
			}
		}
	}
	return string(s)
}
