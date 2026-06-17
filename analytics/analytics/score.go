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

type ScoredLine struct {
	Line                  *Line
	WarmRank              int
	WarmNsPerDoc          float64
	ColdRank              int
	ColdNsPerDoc          float64
	ParseNsPerDoc         float64
	ParsePlusWarmNsPerDoc float64
	ParsePlusWarmRank     int
	ParsePlusColdNsPerDoc float64
	ParsePlusColdRank     int
}

func ScoreGroups(groups [][]*Line) []*ScoredLine {
	res := []*ScoredLine{}
	for i := range groups {
		res = append(res, Score(groups[i])...)
	}
	return res
}

func nums(n int) []int {
	ns := make([]int, n)
	for i := range ns {
		ns[i] = i
	}
	return ns
}

func Score(lines []*Line) []*ScoredLine {
	res := make([]*ScoredLine, len(lines))
	for i := range lines {
		res[i] = &ScoredLine{Line: lines[i]}
	}

	for i := range lines {
		res[i].WarmNsPerDoc = float64(lines[i].WarmNs) / float64(lines[i].Schema.NumInstances)
	}
	sortedWarmPerDoc := nums(len(res))
	slices.SortFunc(sortedWarmPerDoc, func(i, j int) int {
		return FloatCompare(res[i].WarmNsPerDoc, res[j].WarmNsPerDoc)
	})
	for i := range sortedWarmPerDoc {
		res[sortedWarmPerDoc[i]].WarmRank = i + 1
	}

	for i := range lines {
		res[i].ColdNsPerDoc = float64(lines[i].ColdNs) / float64(lines[i].Schema.NumInstances)
	}
	sortedColdPerDoc := nums(len(res))
	slices.SortFunc(sortedColdPerDoc, func(i, j int) int {
		return FloatCompare(res[i].ColdNsPerDoc, res[j].ColdNsPerDoc)
	})
	for i := range sortedColdPerDoc {
		res[sortedColdPerDoc[i]].ColdRank = i + 1
	}

	for i := range lines {
		res[i].ParseNsPerDoc = float64(lines[i].ParseNs) / float64(lines[i].Schema.NumInstances)
		res[i].ParsePlusWarmNsPerDoc = res[i].ParseNsPerDoc + res[i].WarmNsPerDoc
		res[i].ParsePlusColdNsPerDoc = res[i].ParseNsPerDoc + res[i].ColdNsPerDoc
	}

	sortedParsePlusWarm := nums(len(res))
	slices.SortFunc(sortedParsePlusWarm, func(i, j int) int {
		return FloatCompare(res[i].ParsePlusWarmNsPerDoc, res[j].ParsePlusWarmNsPerDoc)
	})
	for i := range sortedParsePlusWarm {
		res[sortedParsePlusWarm[i]].ParsePlusWarmRank = i + 1
	}

	sortedParsePlusCold := nums(len(res))
	slices.SortFunc(sortedParsePlusCold, func(i, j int) int {
		return FloatCompare(res[i].ParsePlusColdNsPerDoc, res[j].ParsePlusColdNsPerDoc)
	})
	for i := range sortedParsePlusCold {
		res[sortedParsePlusCold[i]].ParsePlusColdRank = i + 1
	}

	return res
}
