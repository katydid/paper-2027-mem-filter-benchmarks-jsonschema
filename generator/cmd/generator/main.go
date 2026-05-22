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
	complex "github.com/katydid/validator-jsonschema-benchmarks/generator/schemas/ajv_complex"

	"github.com/katydid/validator-jsonschema-benchmarks/generator/std"
)

type schemaGenerator struct {
	name      string
	gen       randjsonschema.Rand
	schema    string
	validOnly bool
	num       int
}

var generators = []schemaGenerator{
	{
		name:   "example-address-mixed",
		gen:    schemas.RandomAddress(),
		schema: schemas.SchemaJSONSchemaExampleAddress,
	},
	{
		name:      "example-address-valid",
		gen:       schemas.RandomAddress(),
		schema:    schemas.SchemaJSONSchemaExampleAddress,
		validOnly: true,
	},
	{
		name:   "example-blogpost-mixed",
		gen:    schemas.RandomBlogPost(),
		schema: schemas.SchemaJSONSchemaExampleBlogPost,
	},
	{
		name:      "example-blogpost-valid",
		gen:       schemas.RandomBlogPost(),
		schema:    schemas.SchemaJSONSchemaExampleBlogPost,
		validOnly: true,
	},
	{
		name:   "ajv-complex-mixed",
		gen:    complex.Complex(),
		schema: complex.SchemaComplexNew,
		num:    100,
	},
	{
		name:      "ajv-complex-valid",
		gen:       complex.Complex(),
		schema:    complex.SchemaComplexNew,
		validOnly: true,
		num:       100,
	},
	{
		name:   "katydid-conf-mixed",
		gen:    schemas.RandomConfIsIn2026OrLate2025AndEU(),
		schema: schemas.SchemaConfIsIn2026OrLate2025AndEU,
	},
	{
		name:      "katydid-conf-valid",
		gen:       schemas.RandomConfIsIn2026OrLate2025AndEU(),
		schema:    schemas.SchemaConfIsIn2026OrLate2025AndEU,
		validOnly: true,
	},
	{
		name:   "example-userprofile-mixed",
		gen:    schemas.RandomUserProfile(),
		schema: schemas.SchemaJSONSchemaExampleUserProfile,
	},
	{
		name:      "example-userprofile-valid",
		gen:       schemas.RandomUserProfile(),
		schema:    schemas.SchemaJSONSchemaExampleUserProfile,
		validOnly: true,
	},
	{
		name:   "example-calendar-mixed",
		gen:    schemas.RandomCalendar(),
		schema: schemas.SchemaJSONSchemaExampleCalendar,
	},
	{
		name:      "example-calendar-valid",
		gen:       schemas.RandomCalendar(),
		schema:    schemas.SchemaJSONSchemaExampleCalendar,
		validOnly: true,
	},
	{
		name:   "example-devicetype-mixed",
		gen:    schemas.RandomDevicetype(),
		schema: schemas.SchemaJSONSchemaExampleDevicetype,
	},
	{
		name:      "example-devicetype-valid",
		gen:       schemas.RandomDevicetype(),
		schema:    schemas.SchemaJSONSchemaExampleDevicetype,
		validOnly: true,
	},
	{
		name:   "example-health-record-mixed",
		gen:    schemas.RandomHealthRecord(),
		schema: schemas.SchemaJSONSchemaExampleHealthRecord,
	},
	{
		name:      "example-health-record-valid",
		gen:       schemas.RandomHealthRecord(),
		schema:    schemas.SchemaJSONSchemaExampleHealthRecord,
		validOnly: true,
	},
	{
		name:   "example-job-posting-mixed",
		gen:    schemas.RandomJobPosting(),
		schema: schemas.SchemaJSONSchemaExampleJobPosting,
	},
	{
		name:      "example-job-posting-valid",
		gen:       schemas.RandomJobPosting(),
		schema:    schemas.SchemaJSONSchemaExampleJobPosting,
		validOnly: true,
	},
	{
		name:   "example-movie-mixed",
		gen:    schemas.RandomMovie(),
		schema: schemas.SchemaJSONSchemaExampleMovie,
	},
	{
		name:      "example-movie-valid",
		gen:       schemas.RandomMovie(),
		schema:    schemas.SchemaJSONSchemaExampleMovie,
		validOnly: true,
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
	r := rand.NewRandWithSeed(*seed)
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
		generateJSONL(r, gen.gen, validator, number, subfolder, !gen.validOnly)
	}

}

func generateJSONL(r rand.Rand, gen randjsonschema.Rand, validator *jsonschema.Schema, num int, folder string, mixed bool) {
	file, err := os.Create(filepath.Join(folder, "instances.jsonl"))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	for i := range num {
		var s string
		if mixed && i%2 == 0 && i > 10 {
			s = genWrong(r, gen, validator)
		} else {
			s = genRight(r, gen, validator)
		}
		file.WriteString(s + "\n")
	}
}

func genWrong(r rand.Rand, gen randjsonschema.Rand, validator *jsonschema.Schema) string {
	s := gen.Wrong(r)
	v, err := isValid(validator, []byte(s))
	if err != nil {
		panic(err)
	}
	if v {
		log.Printf("regenerating, since we expected invalid for %s", s)
		return genWrong(r, gen, validator)
	}
	return s
}

func genRight(r rand.Rand, gen randjsonschema.Rand, validator *jsonschema.Schema) string {
	s := gen.Right(r)
	v, err := isValid(validator, []byte(s))
	if err != nil {
		panic(err)
	}
	if !v {
		log.Printf("regenerating, since we expected valid for %s", s)
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
	a, err := std.UnmarshalToAny(jsonData)
	if err != nil {
		return false, err
	}
	err = validator.Validate(a)
	return err == nil, nil
}
