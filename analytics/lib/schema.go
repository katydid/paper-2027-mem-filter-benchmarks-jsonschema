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
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Schema struct {
	Name                   string
	Features               []string
	Source                 string
	Generated              bool
	SchemaSizeBytes        int
	NumInstances           int
	AvgInstanceSizeBytes   float64
	RmUniqueItems          bool
	HasExistingReplacement bool
}

func CollectSchemas(folder string) ([]*Schema, error) {
	dirs, err := os.ReadDir(folder)
	if err != nil {
		log.Fatalf("problem reading folder %s got error: %v", folder, err)
	}
	schemas := []*Schema{}
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		path := filepath.Join(folder, dir.Name())
		log.Printf("collecting analytics at %s", path)
		s, err := collectSchema(path)
		if err != nil {
			log.Fatal(err)
		}
		schemas = append(schemas, s)
		for _, schema := range schemas {
			if ContainsRmUniqueItems(schemas, schema.Name) {
				schema.HasExistingReplacement = true
			}
		}
	}
	return schemas, nil
}

func collectSchema(folder string) (*Schema, error) {
	s := &Schema{}
	s.Name = filepath.Base(folder)
	if strings.Contains(s.Name, "rmUniqueItems") {
		s.RmUniqueItems = true
	}
	files, err := os.ReadDir(folder)
	if err != nil {
		log.Printf("error reading folder: %s", folder)
		return nil, err
	}
	for _, file := range files {
		basename := file.Name()
		filename := filepath.Join(folder, basename)
		data, err := os.ReadFile(filename)
		if err != nil {
			log.Printf("error reading file: %s", filename)
			return nil, err
		}
		switch basename {
		case "schema.json":
			s.SchemaSizeBytes = len(data)
			s.Features = collectFeatures(data)
		case "source.txt":
			s.Source = string(data)
		case ".generated":
			s.Generated = true
		case "instances.jsonl":
			instances := bytes.Split(data, []byte("\n"))
			s.NumInstances = len(instances)
			// last line might be empty, remember to not count that
			if len(bytes.TrimSpace(instances[len(instances)-1])) == 0 {
				s.NumInstances = s.NumInstances - 1
			}
			s.AvgInstanceSizeBytes = float64(len(data)) / float64(len(instances))
		}
	}
	return s, nil
}

func collectFeatures(data []byte) []string {
	features := []string{}
	if bytes.Contains(data, []byte(`"uniqueItems": true`)) || bytes.Contains(data, []byte(`"uniqueItems":true`)) {
		features = append(features, "uniqueItems")
	}
	return features
}

func ContainsRmUniqueItems(schemas []*Schema, name string) bool {
	if strings.HasSuffix(name, "-mixed") {
		name = name[:len(name)-6] + "-rmUniqueItems-mixed"
	} else if strings.HasSuffix(name, "-valid") {
		name = name[:len(name)-6] + "-rmUniqueItems-valid"
	} else {
		name = name + "-rmUniqueItems"
	}
	for _, schema := range schemas {
		if schema.Name == name {
			return true
		}
	}
	return false
}

func RemoveUniqueItems(schemas []*Schema) []*Schema {
	res := []*Schema{}
	for i, schema := range schemas {
		if schema.HasExistingReplacement {
			continue
		}
		if strings.Contains(schema.Name, "rmUniqueItems") {
			s := schema
			s.Name = strings.Replace(s.Name, "-rmUniqueItems", "", 1)
			res = append(res, s)
			continue
		}
		res = append(res, schemas[i])
	}
	return res
}

func RemoveSourcePrefixFromName(schemas []*Schema) []*Schema {
	for i := range schemas {
		if strings.HasPrefix(schemas[i].Name, "ajv-") {
			schemas[i].Name = strings.Replace(schemas[i].Name, "ajv-", "", 1)
		}
		if strings.HasPrefix(schemas[i].Name, "jsck-") {
			schemas[i].Name = strings.Replace(schemas[i].Name, "jsck-", "", 1)
		}
		if strings.HasPrefix(schemas[i].Name, "example-") {
			schemas[i].Name = strings.Replace(schemas[i].Name, "example-", "", 1)
		}
		if strings.HasPrefix(schemas[i].Name, "zschema-") {
			schemas[i].Name = strings.Replace(schemas[i].Name, "zschema-", "", 1)
		}
		if strings.HasPrefix(schemas[i].Name, "katydid-") {
			schemas[i].Name = strings.Replace(schemas[i].Name, "katydid-", "", 1)
		}
	}
	return schemas
}

func GroupGenerated(schemas []*Schema) []*Schema {
	res := []*Schema{}
	for i := range schemas {
		if strings.HasSuffix(schemas[i].Name, "-mixed") {
			mixedSchema := schemas[i]
			name := mixedSchema.Name[:len(mixedSchema.Name)-6]
			validName := name + "-valid"
			validSchema := FindSchema(schemas, validName)
			groupSchema := &Schema{
				Name:                   name,
				Features:               mixedSchema.Features,
				Source:                 mixedSchema.Source,
				SchemaSizeBytes:        mixedSchema.SchemaSizeBytes,
				NumInstances:           mixedSchema.NumInstances,
				AvgInstanceSizeBytes:   (mixedSchema.AvgInstanceSizeBytes + validSchema.AvgInstanceSizeBytes) / 2,
				RmUniqueItems:          mixedSchema.RmUniqueItems,
				HasExistingReplacement: mixedSchema.HasExistingReplacement,
				Generated:              true,
			}
			res = append(res, groupSchema)
		} else if strings.HasSuffix(schemas[i].Name, "-valid") {
			// ignore, grouping happens in mixed
			continue
		} else {
			res = append(res, schemas[i])
		}
	}
	return res
}

func FindSchema(schemas []*Schema, name string) *Schema {
	for i := range schemas {
		if schemas[i].Name == name {
			return schemas[i]
		}
	}
	return nil
}
