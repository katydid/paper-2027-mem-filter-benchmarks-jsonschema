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

func Ungroup(groups [][]*Line) []*Line {
	res := []*Line{}
	for _, group := range groups {
		res = append(res, group...)
	}
	return res
}

func GroupBySchema(lines []*Line) [][]*Line {
	groupIndexes := make(map[string]int)
	groups := [][]*Line{}
	for i, line := range lines {
		if _, ok := groupIndexes[line.Schema.Name]; !ok {
			groupIndexes[line.Schema.Name] = len(groups)
			groups = append(groups, []*Line{})
		}
		index := groupIndexes[line.Schema.Name]
		groups[index] = append(groups[index], lines[i])
	}
	return groups
}

func GroupByImplementation(lines []*ScoredLine) [][]*ScoredLine {
	groupIndexes := make(map[string]int)
	groups := [][]*ScoredLine{}
	for i, line := range lines {
		if _, ok := groupIndexes[line.Line.Implementation]; !ok {
			groupIndexes[line.Line.Implementation] = len(groups)
			groups = append(groups, []*ScoredLine{})
		}
		index := groupIndexes[line.Line.Implementation]
		groups[index] = append(groups[index], lines[i])
	}
	return groups
}
