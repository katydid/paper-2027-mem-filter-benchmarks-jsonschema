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
	"strings"
	"time"

	jsonschema "github.com/json-schema-spec/json-schema-go"
)

// Does not support: "dependentRequired", "format", "json-schema-go/issues/2"

const WarmupIterations = 1000
const MaxWarmupTime = 10_000_000_000

func validateAll(instances []any, sch jsonschema.Validator, want bool) error {
	var result error
	for i := range instances {
		_, err := sch.Validate(instances[i])
		if want && err != nil {
			result = err
		} else if !want && err == nil {
			result = fmt.Errorf("expected invalid, but got valid at %d", i)
		}
	}
	return result
}

func main() {
	log.SetFlags(log.Lshortfile)
	if len(os.Args) < 2 {
		log.Fatal("Please provide the example folder path as an argument")
	}

	exampleFolder := os.Args[1]
	want := !strings.Contains(exampleFolder, "-invalid")
	log.Printf("folder %q with base %s expect %v", exampleFolder, filepath.Base(exampleFolder), want)

	// Construct and canonicalize file paths
	schemaFile, err := filepath.Abs(filepath.Join(exampleFolder, "schema.json"))
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

	schema, err := unmarshalToAny(schemaData)
	if err != nil {
		log.Fatalf("Error unmarshaling schema: %v", err)
	}
	validator, err := jsonschema.NewValidator([]any{schema})
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
	err = validateAll(instances, validator, want)
	if err != nil {
		log.Fatalf("Validation failed: %v", err)
	}
	coldDuration := time.Since(coldStart)

	// Warmup
	iterations := math.Ceil(float64(MaxWarmupTime) / float64(coldDuration.Nanoseconds()))
	for _ = range int64(min(iterations, WarmupIterations)) {
		validateAll(instances, validator, want)
	}

	warmStart := time.Now()
	validateAll(instances, validator, want)
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
