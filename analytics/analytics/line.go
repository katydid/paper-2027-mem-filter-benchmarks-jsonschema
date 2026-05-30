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
	"strconv"
)

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

func average(indexes []int, get func(int) int) int {
	sum := 0
	for _, index := range indexes {
		sum += get(index)
	}
	return sum / len(indexes)
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
