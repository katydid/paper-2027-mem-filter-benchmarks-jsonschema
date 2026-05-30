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
	Name                string
	Scores              []*ScoredLine
	AverageWarmNsPerDoc float64
	MedianWarmNsPerDoc  float64
	AverageWarmRank     float64
	MedianWarmRank      float64
	AverageColdNsPerDoc float64
	MedianColdNsPerDoc  float64
	AverageColdRank     float64
	MedianColdRank      float64
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

func Average(num int, get func(int) float64) float64 {
	sum := float64(0)
	for i := range num {
		sum += get(i)
	}
	return sum / float64(num)
}

func AnalyseImplementations(scores []*ScoredLine) []*Implementation {
	groups := GroupByImplementation(scores)
	impls := make([]*Implementation, len(groups))
	for i := range groups {
		impls[i] = AnalyseImplementation(groups[i])
	}
	return impls
}

func AnalyseImplementation(scores []*ScoredLine) *Implementation {
	res := &Implementation{}
	res.Name = scores[0].Line.Implementation
	res.MedianWarmNsPerDoc = MedianFloat(len(scores), func(index int) float64 {
		return scores[index].WarmNsPerDoc
	})
	res.AverageWarmNsPerDoc = Average(len(scores), func(index int) float64 {
		return scores[index].WarmNsPerDoc
	})
	res.MedianColdNsPerDoc = MedianFloat(len(scores), func(index int) float64 {
		return scores[index].ColdNsPerDoc
	})
	res.AverageColdNsPerDoc = Average(len(scores), func(index int) float64 {
		return scores[index].ColdNsPerDoc
	})
	res.MedianWarmRank = MedianFloat(len(scores), func(index int) float64 {
		return float64(scores[index].WarmRank)
	})
	res.AverageWarmRank = Average(len(scores), func(index int) float64 {
		return float64(scores[index].WarmRank)
	})
	res.MedianColdRank = MedianFloat(len(scores), func(index int) float64 {
		return float64(scores[index].ColdRank)
	})
	res.AverageColdRank = Average(len(scores), func(index int) float64 {
		return float64(scores[index].ColdRank)
	})
	return res
}
