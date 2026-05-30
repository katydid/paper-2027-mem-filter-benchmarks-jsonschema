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

import "slices"

func filter[A any](lines []A, pred func(A) bool) []A {
	res := []A{}
	for i := range lines {
		if pred(lines[i]) {
			res = append(res, lines[i])
		}
	}
	return res
}

// Filter only lines with the select implementations.
func FilterImplementations(lines []*Line, impls []string) []*Line {
	return filter(lines, func(line *Line) bool {
		return slices.Contains(impls, line.Implementation)
	})
}

// Remove all schemas where one implementation had an exit status 1
func FilterSchemasExitStatus0(lines []*Line) []*Line {
	groups := GroupBySchema(lines)
	filteredGroups := filter(groups, func(lines []*Line) bool {
		for _, line := range lines {
			if line.ExitStatus == 1 {
				return false
			}
		}
		return true
	})
	return Ungroup(filteredGroups)
}

// Remove any schema with uniqueItems
func FilterNoUniqueItems(lines []*Line) []*Line {
	return filter(lines, func(line *Line) bool {
		return !slices.Contains(line.Schema.Features, "uniqueItems")
	})
}

// Remove exit status 1
func FilterExitStatus0(lines []*Line) []*Line {
	return filter(lines, func(line *Line) bool {
		return line.ExitStatus == 0
	})
}
