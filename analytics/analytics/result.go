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
	"encoding/csv"
	"fmt"
	"os"
	"slices"
	"strconv"
)

type ScoredLine struct {
	Line                  *Line
	WarmRank              int
	WarmSlowDown          float64
	WarmNsPerDoc          float64
	ColdRank              int
	ColdSlowDown          float64
	ColdNsPerDoc          float64
	ParseTODO             bool
	ParseNsPerDoc         float64
	ParsePlusWarmNsPerDoc float64
	ParsePlusWarmRank     int
	ParsePlusWarmSlowDown float64
	ParsePlusColdNsPerDoc float64
	ParsePlusColdRank     int
	ParsePlusColdSlowDown float64
}

type Line struct {
	Schema         *Schema
	Implementation string
	Version        string
	ColdNs         int
	WarmNs         int
	ParseNs        int
	ParseTODO      bool
	CompileNs      int
	Memory         int
	ExitStatus     int
}

func ReadLines(schemas []*Schema, filename string) ([]*Line, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	lines, err := parseLines(schemas, records)
	if err != nil {
		return nil, err
	}
	return lines, nil
}

func parseLines(schemas []*Schema, records [][]string) ([]*Line, error) {
	// check if first line is header and remove it
	if records[0][0] == "implementation" {
		records = records[1:]
	}
	lines := []*Line{}
	for i := range records {
		line, err := parseLine(schemas, records[i])
		if err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}
	return lines, nil
}

func parseLine(schemas []*Schema, record []string) (*Line, error) {
	var err error
	line := &Line{}
	line.Implementation = record[0]
	line.Version = record[1]
	schemaName := record[2]
	line.Schema = FindSchema(schemas, schemaName)
	if line.Schema == nil {
		return nil, fmt.Errorf("could not find schema %s", schemaName)
	}
	line.ColdNs, err = strconv.Atoi(record[3])
	if err != nil {
		return nil, err
	}
	line.WarmNs, err = strconv.Atoi(record[4])
	if err != nil {
		return nil, err
	}
	if record[5] != "TODO" {
		line.ParseNs, err = strconv.Atoi(record[5])
		if err != nil {
			return nil, err
		}
	} else {
		line.ParseTODO = true
	}
	line.CompileNs, err = strconv.Atoi(record[6])
	if err != nil {
		return nil, err
	}
	line.Memory, err = strconv.Atoi(record[7])
	if err != nil {
		return nil, err
	}
	line.ExitStatus, err = strconv.Atoi(record[8])
	if err != nil {
		return nil, err
	}
	return line, nil
}

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

func Ungroup(groups [][]*Line) []*Line {
	res := []*Line{}
	for _, group := range groups {
		res = append(res, group...)
	}
	return res
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

func nums(n int) []int {
	ns := make([]int, n)
	for i := range ns {
		ns[i] = i
	}
	return ns
}

func AverageRuns(lines []*Line) []*Line {
	runs := map[string][]int{}
	runNames := []string{}
	for i, line := range lines {
		runName := line.Implementation + ":" + line.Schema.Name
		if _, ok := runs[runName]; !ok {
			runs[runName] = []int{}
			runNames = append(runNames, runName)
		}
		runs[runName] = append(runs[runName], i)
	}

	res := make([]*Line, len(runNames))
	for i, runName := range runNames {
		runIndexes := runs[runName]
		avgColdNs := average(runIndexes, func(index int) int {
			return lines[index].ColdNs
		})
		avgWarmNs := average(runIndexes, func(index int) int {
			return lines[index].WarmNs
		})
		avgParseNs := average(runIndexes, func(index int) int {
			return lines[index].ParseNs
		})
		avgCompileNs := average(runIndexes, func(index int) int {
			return lines[index].CompileNs
		})
		avgMemory := average(runIndexes, func(index int) int {
			return lines[index].Memory
		})
		line := lines[runIndexes[0]]
		res[i] = &Line{
			Schema:         line.Schema,
			Implementation: line.Implementation,
			Version:        line.Version,
			ColdNs:         avgColdNs,
			WarmNs:         avgWarmNs,
			ParseNs:        avgParseNs,
			ParseTODO:      line.ParseTODO,
			CompileNs:      avgCompileNs,
			Memory:         avgMemory,
			ExitStatus:     line.ExitStatus,
		}
	}
	return res
}

func average(indexes []int, get func(int) int) int {
	sum := 0
	for _, index := range indexes {
		sum += get(index)
	}
	return sum / len(indexes)
}

func ScoreGroups(groups [][]*Line) []*ScoredLine {
	res := []*ScoredLine{}
	for i := range groups {
		res = append(res, Score(groups[i])...)
	}
	return res
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
	fastestWarmNsPerDoc := res[sortedWarmPerDoc[0]].WarmNsPerDoc
	for i := range sortedWarmPerDoc {
		res[sortedWarmPerDoc[i]].WarmRank = i + 1
	}
	for i := range res {
		res[i].WarmSlowDown = float64(res[i].WarmNsPerDoc) / float64(fastestWarmNsPerDoc)
	}

	for i := range lines {
		res[i].ColdNsPerDoc = float64(lines[i].ColdNs) / float64(lines[i].Schema.NumInstances)
	}
	sortedColdPerDoc := nums(len(res))
	slices.SortFunc(sortedColdPerDoc, func(i, j int) int {
		return FloatCompare(res[i].ColdNsPerDoc, res[j].ColdNsPerDoc)
	})
	fastestColdNsPerDoc := res[sortedColdPerDoc[0]].ColdNsPerDoc
	for i := range sortedColdPerDoc {
		res[sortedColdPerDoc[i]].ColdRank = i + 1
	}
	for i := range lines {
		res[i].ColdSlowDown = float64(res[i].ColdNsPerDoc) / float64(fastestColdNsPerDoc)
	}

	const verySlow = float64(10e20)
	for i := range lines {
		res[i].ParseTODO = lines[i].ParseTODO
		if !res[i].ParseTODO {
			res[i].ParseNsPerDoc = float64(lines[i].ParseNs) / float64(lines[i].Schema.NumInstances)
			res[i].ParsePlusWarmNsPerDoc = res[i].ParseNsPerDoc + res[i].WarmNsPerDoc
			res[i].ParsePlusColdNsPerDoc = res[i].ParseNsPerDoc + res[i].ColdNsPerDoc
		} else {
			res[i].ParseNsPerDoc = verySlow
			res[i].ParsePlusWarmNsPerDoc = verySlow
			res[i].ParsePlusColdNsPerDoc = verySlow
		}
	}

	sortedParsePlusWarm := nums(len(res))
	slices.SortFunc(sortedParsePlusWarm, func(i, j int) int {
		return FloatCompare(res[i].ParsePlusWarmNsPerDoc, res[j].ParsePlusWarmNsPerDoc)
	})
	fastestParsePlusWarmNsPerDoc := res[sortedParsePlusWarm[0]].ParsePlusWarmNsPerDoc
	for i := range sortedParsePlusWarm {
		res[sortedParsePlusWarm[i]].ParsePlusWarmRank = i + 1
	}
	for i := range res {
		res[i].ParsePlusWarmSlowDown = float64(lines[i].WarmNs) / float64(fastestParsePlusWarmNsPerDoc)
	}

	sortedParsePlusCold := nums(len(res))
	slices.SortFunc(sortedParsePlusCold, func(i, j int) int {
		return FloatCompare(res[i].ParsePlusColdNsPerDoc, res[j].ParsePlusColdNsPerDoc)
	})
	fastestParsePlusColdNsPerDoc := res[sortedParsePlusCold[0]].ParsePlusColdNsPerDoc
	for i := range sortedParsePlusCold {
		res[sortedParsePlusCold[i]].ParsePlusColdRank = i + 1
	}
	for i := range res {
		res[i].ParsePlusColdSlowDown = float64(lines[i].ColdNs) / float64(fastestParsePlusColdNsPerDoc)
	}

	return res
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
