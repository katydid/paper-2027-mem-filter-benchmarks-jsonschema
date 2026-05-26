# JSON Schema Benchmark

This benchmark suite builds on the amazing work of [SourceMeta's JSONSchema Benchmarks](https://github.com/sourcemeta-research/jsonschema-benchmark/) for the purpose of benchmarking Katydid's JSONSchema implementation against other JSONSchema validators.

Each validator is run with multiple schemas and a collection of documents that might be valid or mixed, which is expressed in the suffix of the folder name.

## Results

The most recent results can be seen (via GitHub Actions)[https://github.com/katydid/validator-jsonschema-benchmarks/actions/workflows/ci.yml].

## Setup

The benchmark requires make, [Docker](https://www.docker.com/), [dts](https://github.com/martinohmann/dts/releases) and sed.

## Running

To run all the benchmarks, a report can be produced via `make dist/report.csv`.
The `Makefile` accepts parameters `IMPLEMENTATIONS` to specify which implementations to run and `RUNS` for the number of runs per implementation.
For example, `make IMPLEMENTATIONS='blaze jsoncons' RUNS=5` will run the Blaze and jsoncons implementations 5 times each.

## Implementations

All implementations can be found in the `implementations/` subdirectory.
A summary of these implementations is given below:

- [ajv](./implementations/ajv/) (JS)
- [ajv-bun](./implementations/ajv-bun/) (JS with BUN runtime)
- [blaze](./implementations/blaze/) (C++)
- [boon](./implementations/boon/) (Rust)
- [corvus](./implementations/corvus/) (C#)
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
- [networknt](./implementations/networknt/) (Java) (Ignored)
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

### Adding a new implementation

First, each implementation must have a `Dockerfile` that copies in any necessary scripts and installs dependencies.
There is also a `version.sh` script that must output the version of the implementation (often extracted from whatever dependency management tool is used).
Finally, appropriate targets must be added to the `Makefile` to build the Docker container and run the benchmark.
