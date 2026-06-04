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

type SchemaName struct {
	Name          string
	PrefixName    string
	RmUniqueItems bool
	GeneratedKind string
	ShortName     string
}

type Schema struct {
	SchemaName
	Features               []string
	Source                 string
	Generated              bool
	SchemaSizeBytes        int
	NumInstances           int
	MeanInstanceSizeBytes  float64
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

func ParseSchemaName(name string) (*SchemaName, error) {
	s := &SchemaName{Name: name}
	shortName := name
	if strings.Contains(name, "-invalid") {
		s.GeneratedKind = "invalid"
		shortName = strings.Replace(shortName, "-invalid", "", 1)
	} else if strings.Contains(name, "-valid") {
		s.GeneratedKind = "valid"
		shortName = strings.Replace(shortName, "-valid", "", 1)
	}
	if strings.Contains(name, "-rmUniqueItems") {
		s.RmUniqueItems = true
		shortName = strings.Replace(shortName, "-rmUniqueItems", "", 1)
	}
	s.PrefixName, s.ShortName = RemovePrefix(shortName)
	return s, nil
}

func collectSchema(folder string) (*Schema, error) {
	s := &Schema{}
	name := filepath.Base(folder)
	schemaName, err := ParseSchemaName(name)
	if err != nil {
		return nil, err
	}
	s.SchemaName = *schemaName
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
			s.MeanInstanceSizeBytes = float64(len(data)) / float64(len(instances))
		}
	}
	return s, nil
}

func collectFeatures(data []byte) []string {
	features := []string{}
	if bytes.Contains(data, []byte(`"uniqueItems": true`)) || bytes.Contains(data, []byte(`"uniqueItems":true`)) {
		features = append(features, "uniqueItems")
	}
	if bytes.Contains(data, []byte("$dynamicRef")) {
		// cql2 and openapi
		features = append(features, "dynamicRef")
	}
	return features
}

func ContainsRmUniqueItems(schemas []*Schema, name string) bool {
	if strings.HasSuffix(name, "-invalid") {
		name = name[:len(name)-6] + "-rmUniqueItems-invalid"
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
		_, schemas[i].Name = RemovePrefix(schemas[i].Name)
	}
	return schemas
}

func RemovePrefix(name string) (string, string) {
	if strings.HasPrefix(name, "ajv-") {
		return "ajv", strings.Replace(name, "ajv-", "", 1)
	}
	if strings.HasPrefix(name, "jsck-") {
		return "jsck", strings.Replace(name, "jsck-", "", 1)
	}
	if strings.HasPrefix(name, "example-") {
		return "example", strings.Replace(name, "example-", "", 1)
	}
	if strings.HasPrefix(name, "zschema-") {
		return "zschema", strings.Replace(name, "zschema-", "", 1)
	}
	if strings.HasPrefix(name, "katydid-") {
		return "katydid", strings.Replace(name, "katydid-", "", 1)
	}
	return "", name
}

func GroupGenerated(schemas []*Schema) []*Schema {
	res := []*Schema{}
	for i := range schemas {
		if strings.HasSuffix(schemas[i].Name, "-valid") {
			validSchema := schemas[i]
			name := validSchema.Name[:len(validSchema.Name)-6]
			invalidName := name + "-invalid"
			invalidSchema := FindSchema(schemas, invalidName)
			groupSchema := &Schema{
				SchemaName: SchemaName{
					Name:          name,
					PrefixName:    validSchema.PrefixName,
					RmUniqueItems: validSchema.RmUniqueItems,
					GeneratedKind: "",
					ShortName:     validSchema.ShortName,
				},
				Features:               validSchema.Features,
				Source:                 validSchema.Source,
				SchemaSizeBytes:        validSchema.SchemaSizeBytes,
				NumInstances:           validSchema.NumInstances,
				MeanInstanceSizeBytes:  (validSchema.MeanInstanceSizeBytes + invalidSchema.MeanInstanceSizeBytes) / 2,
				HasExistingReplacement: validSchema.HasExistingReplacement,
				Generated:              true,
			}
			res = append(res, groupSchema)
		} else if strings.HasSuffix(schemas[i].Name, "-invalid") {
			// ignore, grouping happens in valid
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
