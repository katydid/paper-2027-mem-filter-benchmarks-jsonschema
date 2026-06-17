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

import (
	"slices"
)

type Implementation struct {
	Name   string
	Scores []*ScoredLine

	CompletedSchemas    []string
	NotCompletedSchemas []string
	AllSchemas          []string

	MeanWarmNsPerDoc   float64
	MedianWarmNsPerDoc float64
	MeanWarmRank       float64
	MedianWarmRank     float64
	MeanColdNsPerDoc   float64
	MedianColdNsPerDoc float64
	MeanColdRank       float64
	MedianColdRank     float64

	MeanParsePlusWarmNsPerDoc   float64
	MedianParsePlusWarmNsPerDoc float64
	MeanParsePlusWarmRank       float64
	MedianParsePlusWarmRank     float64
	MeanParsePlusColdNsPerDoc   float64
	MedianParsePlusColdNsPerDoc float64
	MeanParsePlusColdRank       float64
	MedianParsePlusColdRank     float64
}

func MedianInt(num int, get func(int) int) int {
	nums := make([]int, num)
	for i := range num {
		nums[i] = get(i)
	}
	slices.Sort(nums)
	// 1 / 2 = 0 => 0
	// 2 / 2 = 1 => 0/1
	// 3 / 2 = 1 => 1
	// 4 / 2 = 2 => 1/2
	if len(nums)%2 == 1 {
		medianIndex := len(nums) / 2
		return nums[medianIndex]
	}
	medianIndex1 := (len(nums) / 2) - 1
	medianIndex2 := len(nums) / 2
	return (nums[medianIndex1] + nums[medianIndex2]) / 2
}

func MedianFloat(num int, get func(int) float64) float64 {
	nums := make([]float64, num)
	for i := range num {
		nums[i] = get(i)
	}
	slices.Sort(nums)
	// 1 / 2 = 0 => 0
	// 2 / 2 = 1 => 0/1
	// 3 / 2 = 1 => 1
	// 4 / 2 = 2 => 1/2
	if len(nums)%2 == 1 {
		medianIndex := len(nums) / 2
		return nums[medianIndex]
	}
	medianIndex1 := (len(nums) / 2) - 1
	medianIndex2 := len(nums) / 2
	return (nums[medianIndex1] + nums[medianIndex2]) / 2
}

func Mean(num int, get func(int) float64) float64 {
	sum := float64(0)
	for i := range num {
		sum += get(i)
	}
	return sum / float64(num)
}

func AnalyseImplementations(scores []*ScoredLine, allSchemas []string, completedPerImpl map[string][]string) []*Implementation {
	groups := GroupByImplementation(scores)
	impls := make([]*Implementation, len(groups))
	for i := range groups {
		name := groups[i][0].Line.Implementation
		impls[i] = AnalyseImplementation(groups[i], allSchemas, completedPerImpl[name])
	}
	return impls
}

func AnalyseImplementation(scores []*ScoredLine, allSchemas []string, completedSchemas []string) *Implementation {
	res := &Implementation{}
	res.Scores = scores
	res.Name = scores[0].Line.Implementation
	res.CompletedSchemas = completedSchemas
	res.AllSchemas = allSchemas
	notCompleted := []string{}
	for i := range allSchemas {
		schema := allSchemas[i]
		if !slices.Contains(res.CompletedSchemas, schema) {
			notCompleted = append(notCompleted, schema)
		}
	}
	res.NotCompletedSchemas = notCompleted

	res.MedianWarmNsPerDoc = MedianFloat(len(scores), func(index int) float64 {
		return scores[index].WarmNsPerDoc
	})
	res.MeanWarmNsPerDoc = Mean(len(scores), func(index int) float64 {
		return scores[index].WarmNsPerDoc
	})
	res.MedianWarmRank = MedianFloat(len(scores), func(index int) float64 {
		return float64(scores[index].WarmRank)
	})
	res.MeanWarmRank = Mean(len(scores), func(index int) float64 {
		return float64(scores[index].WarmRank)
	})

	res.MedianColdNsPerDoc = MedianFloat(len(scores), func(index int) float64 {
		return scores[index].ColdNsPerDoc
	})
	res.MeanColdNsPerDoc = Mean(len(scores), func(index int) float64 {
		return scores[index].ColdNsPerDoc
	})
	res.MedianColdRank = MedianFloat(len(scores), func(index int) float64 {
		return float64(scores[index].ColdRank)
	})
	res.MeanColdRank = Mean(len(scores), func(index int) float64 {
		return float64(scores[index].ColdRank)
	})

	res.MedianParsePlusWarmNsPerDoc = MedianFloat(len(scores), func(index int) float64 {
		return scores[index].ParsePlusWarmNsPerDoc
	})
	res.MeanParsePlusWarmNsPerDoc = Mean(len(scores), func(index int) float64 {
		return scores[index].ParsePlusWarmNsPerDoc
	})
	res.MedianParsePlusWarmRank = MedianFloat(len(scores), func(index int) float64 {
		return float64(scores[index].ParsePlusWarmRank)
	})
	res.MeanParsePlusWarmRank = Mean(len(scores), func(index int) float64 {
		return float64(scores[index].ParsePlusWarmRank)
	})

	res.MedianParsePlusColdNsPerDoc = MedianFloat(len(scores), func(index int) float64 {
		return scores[index].ParsePlusColdNsPerDoc
	})
	res.MeanParsePlusColdNsPerDoc = Mean(len(scores), func(index int) float64 {
		return scores[index].ParsePlusColdNsPerDoc
	})
	res.MedianParsePlusColdRank = MedianFloat(len(scores), func(index int) float64 {
		return float64(scores[index].ParsePlusColdRank)
	})
	res.MeanParsePlusColdRank = Mean(len(scores), func(index int) float64 {
		return float64(scores[index].ParsePlusColdRank)
	})

	return res
}
