package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"time"

	jsonschema "github.com/kaptinlin/jsonschema"
)

// Does not support: "dependencies"

const WarmupIterations = 1000
const MaxWarmupTime = 10_000_000_000

func validateAll(instances []any, sch *jsonschema.Schema) error {
	for _, inst := range instances {
		result := sch.Validate(inst)
		_ = result
		// result.IsValid()
	}
	return nil
}

func main() {
	log.SetFlags(log.Lshortfile)
	if len(os.Args) < 2 {
		log.Fatal("Please provide the example folder path as an argument")
	}

	exampleFolder := os.Args[1]
	log.Printf("benchmarking schema in folder: %s", exampleFolder)

	// Construct and canonicalize file paths
	schemaFile, err := filepath.Abs(filepath.Join(exampleFolder, "schema-noformat.json"))
	if err != nil {
		log.Fatalf("Error constructing schema file path: %v", err)
	}
	schemaData, err := os.ReadFile(schemaFile)
	if err != nil {
		log.Fatalf("Error reading schema data: %v", err)
	}

	instanceFile, err := filepath.Abs(filepath.Join(exampleFolder, "instances.jsonl"))
	if err != nil {
		log.Fatalf("Error constructing instance file path: %v", err)
	}

	// Compile the JSON schema
	compile_start := time.Now()

	compiler := jsonschema.NewCompiler()
	compiler.AssertFormat = true
	validator, err := compiler.Compile(schemaData)
	if err != nil {
		log.Fatalf("Error creating new validator: %v", err)
	}

	compile_duration := time.Since(compile_start)

	if err != nil {
		log.Fatal(err)
	}

	// Open the JSONL file
	f, err := os.Open(instanceFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Decode and store JSON objects
	var instances []any
	reader := bufio.NewReader(f)
	decoder := json.NewDecoder(reader)

	for {
		var inst any
		if err := decoder.Decode(&inst); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("Error decoding JSON: %v", err)
		}
		instances = append(instances, inst)
	}

	// Cold start
	coldStart := time.Now()
	err = validateAll(instances, validator)
	if err != nil {
		log.Fatalf("Validation failed: %v", err)
	}
	coldDuration := time.Since(coldStart)

	// Warmup
	iterations := math.Ceil(float64(MaxWarmupTime) / float64(coldDuration.Nanoseconds()))
	for _ = range int64(min(iterations, WarmupIterations)) {
		validateAll(instances, validator)
	}

	warmStart := time.Now()
	validateAll(instances, validator)
	warmDuration := time.Since(warmStart)

	// Print timing
	fmt.Printf("%d,%d,TODO,%d\n", coldDuration.Nanoseconds(), warmDuration.Nanoseconds(), compile_duration.Nanoseconds())
}

func unmarshalToAny(data []byte) (any, error) {
	if len(data) == 0 {
		return nil, nil
	}
	var v any
	return v, json.Unmarshal(data, &v)
}
