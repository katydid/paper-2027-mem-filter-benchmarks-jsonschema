package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

const WarmupIterations = 1000
const MaxWarmupTime = 10_000_000_000

func validateAll(instances []any, sch *jsonschema.Schema, want bool) error {
	var result error
	for i := range instances {
		err := sch.Validate(instances[i])
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
	c := jsonschema.NewCompiler()
	c.AssertFormat()

	compile_start := time.Now()

	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(schemaData))
	if err != nil {
		log.Fatalf("Error unmarshaling schema: %v", err)
	}
	if err := c.AddResource("schema.json", doc); err != nil {
		log.Fatalf("Error adding resource: %v", err)
	}
	sch, err := c.Compile("schema.json")
	if err != nil {
		log.Fatalf("Error compiling schema: %v", err)
	}
	compile_duration := time.Since(compile_start)

	if err != nil {
		log.Fatal(err)
	}

	// Open the JSONL file
	data, err := os.ReadFile(instanceFile)
	if err != nil {
		log.Fatal(err)
	}
	f := bytes.NewBuffer(data)
	reader := bufio.NewReader(f)

	// Decode and store JSON objects
	parsingStart := time.Now()
	var instances []any
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
	parsingDuration := time.Since(parsingStart)

	// Cold start
	coldStart := time.Now()
	err = validateAll(instances, sch, want)
	if err != nil {
		log.Fatalf("Validation failed: %v for schema in folder %s", err, exampleFolder)
	}
	coldDuration := time.Since(coldStart)

	// Warmup
	iterations := math.Ceil(float64(MaxWarmupTime) / float64(coldDuration.Nanoseconds()))
	for _ = range int64(min(iterations, WarmupIterations)) {
		validateAll(instances, sch, want)
	}

	warmStart := time.Now()
	validateAll(instances, sch, want)
	warmDuration := time.Since(warmStart)

	// Print timing
	fmt.Printf("%d,%d,%d,%d\n", coldDuration.Nanoseconds(), warmDuration.Nanoseconds(), parsingDuration.Nanoseconds(), compile_duration.Nanoseconds())
}
