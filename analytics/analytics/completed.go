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

package analytics

import "fmt"

func CountTotalSchemas(lines []*Line) int {
	schemas := map[string]struct{}{}
	for i := range lines {
		schemas[lines[i].Schema.Name] = struct{}{}
	}
	return len(schemas)
}

func CountCompletedPerImplementation(lines []*Line) map[string]int {
	res := make(map[string]map[string]bool)
	for i := range lines {
		impl := lines[i].Implementation
		if _, ok := res[impl]; !ok {
			res[impl] = make(map[string]bool)
		}
		schema := lines[i].Schema.Name
		success := lines[i].ExitStatus == 0
		if _, ok := res[impl][schema]; !ok {
			res[impl][schema] = success
		} else {
			prev := res[impl][schema]
			if prev != success {
				panic(fmt.Sprintf("some failed and some succeeded for implementation: %s and schema: %s", impl, schema))
			}
		}
	}
	result := make(map[string]int)
	for impl := range res {
		count := 0
		for _, succ := range res[impl] {
			if succ {
				count++
			}
		}
		result[impl] = count
	}
	return result
}
