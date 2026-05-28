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
	goreflect "reflect"
	"time"

	"github.com/katydid/parser-go-reflect/reflect"
	"github.com/katydid/validator-go-jsonschema/jsonschema"
)

const WarmupIterations = 1000
const MaxWarmupTime = 10_000_000_000

func validateAll(parser reflect.Parser, matcher jsonschema.Matcher, instances []goreflect.Value) error {
	for _, inst := range instances {
		parser.Init(inst)
		if _, err := matcher.MatchParser(parser); err != nil {
			return err
		}
	}
	return nil
}

var notSupported = map[string]string{
	"ajv-cosmicrealms-mixed": "uniqueItems not supported",
	"ajv-cosmicrealms-valid": "uniqueItems not supported",
	"cspell":                 "uniqueItems not supported",
	"deno":                   "uniqueItems not supported",
	"draft-04":               "uniqueItems not supported",
	"jsconfig":               "uniqueItems not supported",
	"krakend":                "uniqueItems not supported",
	"lazygit":                "uniqueItems not supported",
	"stylecop":               "uniqueItems not supported",
	"ui5-manifest":           "uniqueItems not supported",
	"unreal-engine-uproject": "uniqueItems not supported",
	"zschema-advanced-mixed": "uniqueItems not supported",
	"zschema-advanced-valid": "uniqueItems not supported",
	"zschema-basic-mixed":    "uniqueItems not supported",
	"zschema-basic-valid":    "uniqueItems not supported",

	"ansible-meta":  "json: cannot unmarshal bool into Go struct field Schema.Object.properties of type schema.Schema",
	"cmake-presets": "just takes long",
	"cql2":          "could not find schema for #/$defs/andOrExpression",
	"geojson":       "timed out",
}

func main() {
	log.SetFlags(log.Lshortfile)
	if len(os.Args) < 2 {
		log.Fatal("Please provide the example folder path as an argument")
	}

	exampleFolder := os.Args[1]
	log.Printf("folder %q with base %s", exampleFolder, filepath.Base(exampleFolder))
	if reason, ok := notSupported[filepath.Base(exampleFolder)]; ok {
		log.Fatalf("%s is not supported, because %s", exampleFolder, reason)
	}

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
	matcher, err := jsonschema.Compile([]byte(schemaData))
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
	var instances []goreflect.Value
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
		instances = append(instances, goreflect.ValueOf(inst))
	}
	parser := reflect.NewParser()

	// Cold start
	coldStart := time.Now()
	err = validateAll(parser, matcher, instances)
	if err != nil {
		log.Fatalf("Validation failed: %v", err)
	}
	coldDuration := time.Since(coldStart)

	// Warmup
	iterations := math.Ceil(float64(MaxWarmupTime) / float64(coldDuration.Nanoseconds()))
	for _ = range int64(min(iterations, WarmupIterations)) {
		validateAll(parser, matcher, instances)
	}

	warmStart := time.Now()
	validateAll(parser, matcher, instances)
	warmDuration := time.Since(warmStart)

	// Print timing
	fmt.Printf("%d,%d,TODO,%d\n", coldDuration.Nanoseconds(), warmDuration.Nanoseconds(), compile_duration.Nanoseconds())
}
