# JSON Schema Benchmark

This benchmark suite builds on the amazing work of [SourceMeta's JSONSchema Benchmarks](https://github.com/sourcemeta-research/jsonschema-benchmark/) for the purpose of benchmarking Katydid's JSONSchema implementation against other JSONSchema validators.

Each validator is run with multiple schemas and a collection of documents that might be valid or invalid, which is expressed in the suffix of the folder name.

## Results

The most recent results can be seen [via GitHub Actions](https://github.com/katydid/validator-jsonschema-benchmarks/actions/workflows/ci.yml).

## Setup

The benchmarks require:

* posix tools: make, sed, printf
* [Docker](https://www.docker.com/)

## Running

There are several ways to run the benchmarks:

* Run all benchmarks and produce a report: `make run` or `make dist/report.csv`
* Run only specific implementations: `make IMPLEMENTATIONS='blaze jsoncons' RUNS=5`
* Run only specific schemas: `make SCHEMAS='example-address-valid'`

## Analytics

Analytics currently requires Go to be installed:

* `make analytics-schemas-latex` creates a latex file `schemas.latex` with analytics on the schemas.
* `make analytics-results-latex` creates a latex file `results.latex` with analytics on the results after running per schema and implementation.
* `make analytics-impls-md` creates a markdown file `impls.md` with analytics on the results after running per implementation.

## Implementations

All implementations can be found in the `implementations/` subdirectory.
A summary of these implementations is given below:

- [ajv](./implementations/ajv/) (JS)
- [ajv-bun](./implementations/ajv-bun/) (JS with BUN runtime)
- [blaze](./implementations/blaze/) (C++)
- [boon](./implementations/boon/) (Rust)
- [corvus](./implementations/corvus/) (generated C#)
- [go-google](./implementations/go-google/) (Go)
- [go-json-schema-spec](./implementations/go-json-schema-spec/) (Go)
- [go-kaptinlin](./implementations/go-kaptinlin/) (Go)
- [go-katydid-auto](./implementations/go-katydid-auto/) (Go)
- [go-katydid-mem](./implementations/go-katydid-mem/) (Go)
- [go-santhosh-tekuri](./implementations/go-santhosh-tekuri/) (Go)
- [hyperjump](./implementations/hyperjump/) (JS)
- [jsdotnet](./implementations/jsdotnet/) (C#)
- [json_schemer](./implementations/json_schemer/) (Ruby)
- [jsoncons](./implementations/jsoncons/) (C++)
- [jsu-c](./implementations/jsu-c/) (generated C)
- [jsu-java](./implementations/jsu-java/) (generated Java)
- [jsu-js](./implementations/jsu-js/) (generated JS)
- [jsu-pl](./implementations/jsu-pl/) (generated Perl)
- [jsu-py](./implementations/jsu-py/) (generated Python)
- [JSV](./implementations/jsv) (Elixir)
- [kmp](./implementations/kmp) (Kotlin)
- [networknt](./implementations/networknt/) (Java)
- [opis](./implementations/opis/) (PHP)
- [py-jsonschema](./implementations/py-jsonschema/) (Python)
- [rapidjson](./implementations/rapidjson/) (C++)
- [schemasafe](./implementations/schemasafe/) (JS)

Compared to the original [SourceMeta's JSONSchema Benchmarks](https://github.com/sourcemeta-research/jsonschema-benchmark/) the following libraries were added: 

- [go-google](./implementations/go-google/) (Go)
- [go-json-schema-spec](./implementations/go-json-schema-spec/) (Go)
- [go-kaptinlin](./implementations/go-kaptinlin/) (Go)
- [go-katydid-auto](./implementations/go-katydid-auto/) (Go)
- [go-katydid-mem](./implementations/go-katydid-mem/) (Go)
- [rapidjson](./implementations/rapidjson/) (C++)

Also note that [go-santhosh-tekuri](./implementations/go-santhosh-tekuri/) was renamed from [go-jsonschema](https://github.com/sourcemeta-research/jsonschema-benchmark/tree/main/implementations/go-jsonschema).

Each implementation is run via Docker.
First, a Docker container is built with all the necessary dependencies.
Then, at runtime, a folder containing the schema and the necessary dependencies is mounted and the time to validate all documents is measured.

Implementations can be ignored by adding a `.benchmark-ignore` file in the implementation subdirectory.
It also worth noting that some implementations compile schemas ahead of time into a more efficient representation, while others interpret the entire schema at runtime.

### Curated implementations

We also have a list of [curated_implemenatations.txt](./curated_implemenatations.txt).
So far we have filtered out the following implementations for not supporting `format`: 
* go-google
* go-json-schema-spec
* rapidjson

### TODO these to handle invalid inputs

* json_schemer
* jsu-java, jsu-js, jsu-pl and jsu-py.
* jsv
* kmp
* opis
* schemasafe

### TODO these need to handle format or be removed from the curated list

* json_schemer
* jsv
* kmp
* networknt
* opis
* schemasafe

### TODO these need to add a measurement for parsing time

* hyperjump
* json_schemer
* jsu-java, jsu-js, jsu-pl and jsu-py
* jsv
* kmp
* networknt
* opis
* schemasafe

### Adding a new implementation

First, each implementation must have a `Dockerfile` that copies in any necessary scripts and installs dependencies.
There is also a `version.sh` script that must output the version of the implementation (often extracted from whatever dependency management tool is used).
Finally, appropriate targets must be added to the `Makefile` to build the Docker container and run the benchmark.

## Schemas

All schemas are found in the schemas folder
We run a curated list of [curated_schemas.txt](./curated_schemas.txt) where:
* `uniqueItems` have been removed,
* we have removed `cspell` and `ui5-manifest` for using Perl syntax regexes `(?=`
* we have deleted instances of `helm-chart-lock` that had empty strings for `repository` fields, since they are not valid uri's.
* we have totally excluded schemas that use `dynamicRef`.

Some schemas had a collection of instances gathered from github.
The rest we can regenerate `instances.jsonl` files for by running: `make generate`.

### Shortcomings of document generator

* Lots of libraries do not support unicode strings properly, so random generation of documents has been limited to ascii strings.
* Since we using JSONL for documents, we do not generate newlines in strings.
* Number generation has been limited to 64 bit floats.