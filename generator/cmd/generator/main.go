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

package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v6"

	"github.com/katydid/validator-jsonschema-benchmarks/generator/rand"
	"github.com/katydid/validator-jsonschema-benchmarks/generator/rand/randjsonschema"
	"github.com/katydid/validator-jsonschema-benchmarks/generator/schemas"
	"github.com/katydid/validator-jsonschema-benchmarks/generator/schemas/ajv_cosmicrealms"
	complex "github.com/katydid/validator-jsonschema-benchmarks/generator/schemas/jsck_complex"
	medium "github.com/katydid/validator-jsonschema-benchmarks/generator/schemas/jsck_medium"
	advanced "github.com/katydid/validator-jsonschema-benchmarks/generator/schemas/zschema_advanced"
	basic "github.com/katydid/validator-jsonschema-benchmarks/generator/schemas/zschema_basic"
)

type schemaGenerator struct {
	name   string
	gen    randjsonschema.Rand
	schema string
	kind   string
	num    int
}

var generators = []schemaGenerator{
	{
		name:   "example-address-invalid",
		gen:    schemas.RandomAddress(),
		schema: schemas.SchemaJSONSchemaExampleAddress,
		kind:   "invalid",
	},
	{
		name:   "example-address-valid",
		gen:    schemas.RandomAddress(),
		schema: schemas.SchemaJSONSchemaExampleAddress,
		kind:   "valid",
	},
	{
		name:   "example-blogpost-invalid",
		gen:    schemas.RandomBlogPost(),
		schema: schemas.SchemaJSONSchemaExampleBlogPost,
		kind:   "invalid",
	},
	{
		name:   "example-blogpost-valid",
		gen:    schemas.RandomBlogPost(),
		schema: schemas.SchemaJSONSchemaExampleBlogPost,
		kind:   "valid",
	},
	{
		name:   "jsck-complex-invalid",
		gen:    complex.Complex(),
		schema: complex.SchemaComplexNew,
		kind:   "invalid",
		num:    100,
	},
	{
		name:   "jsck-complex-valid",
		gen:    complex.Complex(),
		schema: complex.SchemaComplexNew,
		kind:   "valid",
		num:    100,
	},
	{
		name:   "katydid-conf-invalid",
		gen:    schemas.RandomConfIsIn2026OrLate2025AndEU(),
		schema: schemas.SchemaConfIsIn2026OrLate2025AndEU,
		kind:   "invalid",
	},
	{
		name:   "katydid-conf-valid",
		gen:    schemas.RandomConfIsIn2026OrLate2025AndEU(),
		schema: schemas.SchemaConfIsIn2026OrLate2025AndEU,
		kind:   "valid",
	},
	{
		name:   "example-userprofile-invalid",
		gen:    schemas.RandomUserProfile(),
		schema: schemas.SchemaJSONSchemaExampleUserProfile,
		kind:   "invalid",
	},
	{
		name:   "example-userprofile-valid",
		gen:    schemas.RandomUserProfile(),
		schema: schemas.SchemaJSONSchemaExampleUserProfile,
		kind:   "valid",
	},
	{
		name:   "example-calendar-invalid",
		gen:    schemas.RandomCalendar(),
		schema: schemas.SchemaJSONSchemaExampleCalendar,
		kind:   "invalid",
	},
	{
		name:   "example-calendar-valid",
		gen:    schemas.RandomCalendar(),
		schema: schemas.SchemaJSONSchemaExampleCalendar,
		kind:   "valid",
	},
	{
		name:   "example-devicetype-invalid",
		gen:    schemas.RandomDevicetype(),
		schema: schemas.SchemaJSONSchemaExampleDevicetype,
		kind:   "invalid",
	},
	{
		name:   "example-devicetype-valid",
		gen:    schemas.RandomDevicetype(),
		schema: schemas.SchemaJSONSchemaExampleDevicetype,
		kind:   "valid",
	},
	{
		name:   "example-health-record-invalid",
		gen:    schemas.RandomHealthRecord(),
		schema: schemas.SchemaJSONSchemaExampleHealthRecord,
		kind:   "invalid",
	},
	{
		name:   "example-health-record-valid",
		gen:    schemas.RandomHealthRecord(),
		schema: schemas.SchemaJSONSchemaExampleHealthRecord,
		kind:   "valid",
	},
	{
		name:   "example-job-posting-invalid",
		gen:    schemas.RandomJobPosting(),
		schema: schemas.SchemaJSONSchemaExampleJobPosting,
		kind:   "invalid",
	},
	{
		name:   "example-job-posting-valid",
		gen:    schemas.RandomJobPosting(),
		schema: schemas.SchemaJSONSchemaExampleJobPosting,
		kind:   "valid",
	},
	{
		name:   "example-movie-invalid",
		gen:    schemas.RandomMovie(),
		schema: schemas.SchemaJSONSchemaExampleMovie,
		kind:   "invalid",
	},
	{
		name:   "example-movie-valid",
		gen:    schemas.RandomMovie(),
		schema: schemas.SchemaJSONSchemaExampleMovie,
		kind:   "valid",
	},
	{
		name:   "jsck-medium-invalid",
		gen:    medium.Medium(),
		schema: medium.SchemaMedium,
		kind:   "invalid",
	},
	{
		name:   "jsck-medium-valid",
		gen:    medium.Medium(),
		schema: medium.SchemaMedium,
		kind:   "valid",
	},
	{
		name:   "zschema-basic-invalid",
		gen:    basic.ProductSet(),
		schema: basic.SchemaBasic,
		kind:   "invalid",
		num:    200,
	},
	{
		name:   "zschema-basic-valid",
		gen:    basic.ProductSet(),
		schema: basic.SchemaBasic,
		kind:   "valid",
		num:    200,
	},
	{
		name:   "zschema-basic-rmUniqueItems-invalid",
		gen:    basic.ProductSetrmUniqueItems(),
		schema: basic.SchemaBasicrmUniqueItems,
		kind:   "invalid",
		num:    200,
	},
	{
		name:   "zschema-basic-rmUniqueItems-valid",
		gen:    basic.ProductSetrmUniqueItems(),
		schema: basic.SchemaBasicrmUniqueItems,
		kind:   "valid",
		num:    200,
	},
	{
		name:   "zschema-advanced-invalid",
		gen:    advanced.Advanced(),
		schema: advanced.SchemaAdvanced,
		kind:   "invalid",
	},
	{
		name:   "zschema-advanced-valid",
		gen:    advanced.Advanced(),
		schema: advanced.SchemaAdvanced,
		kind:   "valid",
	},
	{
		name:   "zschema-advanced-rmUniqueItems-invalid",
		gen:    advanced.AdvancedrmUniqueItems(),
		schema: advanced.SchemaAdvancedrmUniqueItems,
		kind:   "invalid",
	},
	{
		name:   "zschema-advanced-rmUniqueItems-valid",
		gen:    advanced.AdvancedrmUniqueItems(),
		schema: advanced.SchemaAdvancedrmUniqueItems,
		kind:   "valid",
	},
	{
		name:   "ajv-cosmicrealms-invalid",
		gen:    ajv_cosmicrealms.CosmicRealms(),
		schema: ajv_cosmicrealms.SchemaCosmicRealms,
		kind:   "invalid",
	},
	{
		name:   "ajv-cosmicrealms-valid",
		gen:    ajv_cosmicrealms.CosmicRealms(),
		schema: ajv_cosmicrealms.SchemaCosmicRealms,
		kind:   "valid",
	},
	{
		name:   "ajv-cosmicrealms-rmUniqueItems-invalid",
		gen:    ajv_cosmicrealms.CosmicRealmsrmUniqueItems(),
		schema: ajv_cosmicrealms.SchemaCosmicRealmsrmUniqueItems,
		kind:   "invalid",
	},
	{
		name:   "ajv-cosmicrealms-rmUniqueItems-valid",
		gen:    ajv_cosmicrealms.CosmicRealmsrmUniqueItems(),
		schema: ajv_cosmicrealms.SchemaCosmicRealmsrmUniqueItems,
		kind:   "valid",
	},
}

func main() {
	log.SetFlags(log.Lshortfile)
	seed := flag.Int64("seed", time.Now().UnixNano(), "seed for random generator (defaults to now)")
	num := flag.Int("num", 1000, "number of random json files to generate (defaults to 10)")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		panic("expected folder where to generate")
	}
	folder := args[0]
	fmt.Printf("generating at %s with seed %d\n", folder, *seed)

	for _, gen := range generators {
		fmt.Printf("generating schema %s\n", gen.name)
		subfolder := filepath.Join(folder, gen.name)
		if _, err := os.Stat(subfolder); err != nil {
			if err := os.Mkdir(subfolder, 0755); err != nil {
				panic(err)
			}
		}
		validator, err := newValidator(gen.schema)
		if err != nil {
			panic(fmt.Sprintf("given schema %s, error: %v", gen.schema, err))
		}

		// mark that this folder was generated
		generatedFilename := filepath.Join(subfolder, ".generated")
		if err := os.WriteFile(generatedFilename, []byte{}, 0644); err != nil {
			panic(err)
		}

		schemaFilename := filepath.Join(subfolder, "schema.json")
		if err := os.WriteFile(schemaFilename, []byte(gen.schema), 0644); err != nil {
			panic(err)
		}

		number := *num
		if gen.num > 0 {
			fmt.Printf("overriding num with %d for schema %s\n", gen.num, gen.name)
			number = gen.num
		}
		generateJSONL(*seed, gen.gen, validator, number, subfolder, gen.kind)
	}

}

func generateJSONL(seed int64, gen randjsonschema.Rand, validator *jsonschema.Schema, num int, folder string, kind string) {
	file, err := os.Create(filepath.Join(folder, "instances.jsonl"))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	r := rand.NewRandWithSeed(seed)
	for range num {
		var s string
		if kind == "valid" {
			s = genRight(r, gen, validator)
		} else if kind == "invalid" {
			s = genWrong(r, gen, validator)
		} else {
			panic(fmt.Sprintf("kind: %q not supported", kind))
		}
		file.WriteString(s + "\n")
	}
}

func genWrong(r rand.Rand, gen randjsonschema.Rand, validator *jsonschema.Schema) string {
	s := gen.Wrong(r)
	v, err := isValid(validator, []byte(s))
	if err != nil {
		panic(fmt.Sprintf("got err: %v for input: %s", err, s))
	}
	if v {
		log.Fatalf("regenerating, since we expected invalid for %s", s)
		return genWrong(r, gen, validator)
	}
	return s
}

func genRight(r rand.Rand, gen randjsonschema.Rand, validator *jsonschema.Schema) string {
	s := gen.Right(r)
	v, err := isValid(validator, []byte(s))
	if err != nil {
		panic(fmt.Sprintf("given input %s error: %v", s, err))
	}
	if !v {
		log.Fatalf("regenerating, since we expected valid for %s", s)
		return genRight(r, gen, validator)
	}
	return s
}

func newValidator(schemaJSON string) (*jsonschema.Schema, error) {
	c := jsonschema.NewCompiler()
	c.AssertFormat()
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader([]byte(schemaJSON)))
	if err != nil {
		return nil, err
	}
	if err := c.AddResource("schema.json", doc); err != nil {
		return nil, err
	}
	return c.Compile("schema.json")
}

func isValid(validator *jsonschema.Schema, jsonData []byte) (bool, error) {
	a, err := jsonschema.UnmarshalJSON(bytes.NewReader(jsonData))
	if err != nil {
		return false, err
	}
	return validator.Validate(a) == nil, nil
}
