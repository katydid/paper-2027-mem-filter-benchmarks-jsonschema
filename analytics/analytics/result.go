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
	Line         *Line
	WarmRank     int
	WarmSlowDown float64
	WarmNsPerDoc float64
	ColdRank     int
	ColdSlowDown float64
	ColdNsPerDoc float64
}

type Line struct {
	Schema         *Schema
	Implementation string
	Version        string
	ColdNs         int
	WarmNs         int
	ParseNs        int
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

func filter(lines []*Line, pred func(*Line) bool) []*Line {
	res := []*Line{}
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

// Remove any schema with uniqueItems
func FilterNoUniqueItems(lines []*Line) []*Line {
	return filter(lines, func(line *Line) bool {
		return slices.Contains(line.Schema.Features, "uniqueItems")
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

	sortedWarm := nums(len(lines))
	slices.SortFunc(sortedWarm, func(i, j int) int {
		return IntCompare(lines[i].WarmNs, lines[j].WarmNs)
	})
	fastestWarmNs := lines[sortedWarm[0]].WarmNs
	for i := range sortedWarm {
		res[sortedWarm[i]].WarmRank = i
	}
	for i := range lines {
		res[i].WarmSlowDown = float64(lines[i].WarmNs) / float64(fastestWarmNs)
	}
	for i := range lines {
		res[i].WarmNsPerDoc = float64(lines[i].WarmNs) / float64(lines[i].Schema.NumInstances)
	}

	sortedCold := nums(len(lines))
	slices.SortFunc(sortedCold, func(i, j int) int {
		return IntCompare(lines[i].ColdNs, lines[j].ColdNs)
	})
	fastestColdNs := lines[sortedCold[0]].ColdNs
	for i := range sortedCold {
		res[sortedCold[i]].ColdRank = i
	}
	for i := range lines {
		res[i].ColdSlowDown = float64(lines[i].ColdNs) / float64(fastestColdNs)
	}
	for i := range lines {
		res[i].ColdNsPerDoc = float64(lines[i].ColdNs) / float64(lines[i].Schema.NumInstances)
	}
	return res
}
