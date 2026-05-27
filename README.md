# JSON Schema Benchmark

This benchmark suite builds on the amazing work of [SourceMeta's JSONSchema Benchmarks](https://github.com/sourcemeta-research/jsonschema-benchmark/) for the purpose of benchmarking Katydid's JSONSchema implementation against other JSONSchema validators.

Each validator is run with multiple schemas and a collection of documents that might be valid or mixed, which is expressed in the suffix of the folder name.

## Results

The most recent results can be seen [via GitHub Actions](https://github.com/katydid/validator-jsonschema-benchmarks/actions/workflows/ci.yml).

## Setup

The benchmarks require:

* posix tools: make, sed, printf
* [Docker](https://www.docker.com/)
* [dts](https://github.com/martinohmann/dts/releases), which can be installed by downloading the binary and putting it in your PATH.

## Running

There are several ways to run the benchmarks:

* Run all benchmarks and produce a report: `make run` or `make dist/report.csv`
* Run only specific implementations: `make IMPLEMENTATIONS='blaze jsoncons' RUNS=5`
* Run only specific schemas: `make SCHEMAS='example-address-valid'`

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

## Shortcomings

* Lots of libraries do not support unicode strings properly, so random generation of documents has been limited to ascii strings.
* Since we using JSONL for documents, we do not generate newlines in strings.
* Number generation has been limited to 64 bit floats.

## Schemas

<table>
<tr>
  <td>Name</td>
  <td>Schema Source</td>
  <td>Documents</td>
  <td>Mixed/Valid Only</td>
  <td>Features</td>
</tr>
<tr>
  <td><a href="./schemas/ajv_cosmicrealms-mixed/">ajv_cosmicrealms-mixed</a></td>
  <td><a href="https://github.com/ajv-validator/ajv/blob/master/spec/tests/schemas/cosmicrealms.json">ajv</a></td>
  <td>generated</td>
  <td>mixed</td>
  <td>uniqueItems</td>
</tr>
<tr>
  <td><a href="./schemas/ajv_cosmicrealms-noUniqueItems-mixed/">ajv_cosmicrealms-noUniqueItems-mixed</a></td>
  <td><a href="https://github.com/ajv-validator/ajv/blob/master/spec/tests/schemas/cosmicrealms.json">ajv</a> with uniqueItems removed</td>
  <td>generated</td>
  <td>mixed</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/ajv_cosmicrealms-noUniqueItems-valid/">ajv_cosmicrealms-noUniqueItems-valid</a></td>
  <td><a href="https://github.com/ajv-validator/ajv/blob/master/spec/tests/schemas/cosmicrealms.json">ajv</a> with uniqueItems removed</td>
  <td>generated</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/ajv_cosmicrealms-valid/">ajv_cosmicrealms-valid</a></td>
  <td><a href="https://github.com/ajv-validator/ajv/blob/master/spec/tests/schemas/cosmicrealms.json">ajv</a></td>
  <td>generated</td>
  <td>valid only</td>
  <td>uniqueItems</td>
</tr>
<tr>
  <td><a href="./schemas/ansible-meta/">ansible-meta</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/aws-cdk/">aws-cdk</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/babelrc/">babelrc</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/clang-format/">clang-format</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/cmake-presets/">cmake-presets</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/code-climate/">code-climate</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/cql2/">cql2</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/cspell/">cspell</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td>uniqueItems</td>
</tr>
<tr>
  <td><a href="./schemas/cypress/">cypress</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/deno/">deno</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td>uniqueItems</td>
</tr>
<tr>
  <td><a href="./schemas/dependabot/">dependabot</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/draft-04/">draft-04</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td>uniqueItems</td>
</tr>
<tr>
  <td><a href="./schemas/example-address-mixed/">example-address-mixed</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>mixed</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/example-address-valid/">example-address-valid</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/example-blogpost-mixed/">example-blogpost-mixed</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>mixed</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/example-blogpost-valid/">example-blogpost-valid</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/example-calendar-mixed/">example-calendar-mixed</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>mixed</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/example-calendar-valid/">example-calendar-valid</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/example-devicetype-mixed/">example-devicetype-mixed</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>mixed</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/example-devicetype-valid/">example-devicetype-valid</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/example-health-record-mixed/">example-health-record-mixed</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>mixed</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/example-health-record-valid/">example-health-record-valid</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/example-job-posting-mixed/">example-job-posting-mixed</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>mixed</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/example-job-posting-valid/">example-job-posting-valid</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/example-movie-mixed/">example-movie-mixed</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>mixed</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/example-movie-valid/">example-movie-valid</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/example-userprofile-mixed/">example-userprofile-mixed</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>mixed</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/example-userprofile-valid/">example-userprofile-valid</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/fabric-mod/">fabric-mod</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/geojson/">geojson</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/gitpod-configuration/">gitpod-configuration</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/helm-chart-lock/">helm-chart-lock</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/importmap/">importmap</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/jasmine/">jasmine</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/jsck_complex-mixed/">jsck_complex-mixed</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>mixed</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/jsck_complex-valid/">jsck_complex-valid</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/jsck_medium-mixed/">jsck_medium-mixed</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>mixed</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/jsck_medium-valid/">jsck_medium-valid</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/jsconfig/">jsconfig</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td>uniqueItems</td>
</tr>
<tr>
  <td><a href="./schemas/jshintrc/">jshintrc</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/katydid-conf-mixed/">katydid-conf-mixed</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>mixed</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/katydid-conf-valid/">katydid-conf-valid</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/krakend/">krakend</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td>uniqueItems</td>
</tr>
<tr>
  <td><a href="./schemas/lazygit/">lazygit</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td>uniqueItems</td>
</tr>
<tr>
  <td><a href="./schemas/lerna/">lerna</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/nest-cli/">nest-cli</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/omnisharp/">omnisharp</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/openapi/">openapi</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/pre-commit-hooks/">pre-commit-hooks</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/pulumi/">pulumi</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/semantic-release/">semantic-release</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/stale/">stale</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/stylecop/">stylecop</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td>uniqueItems</td>
</tr>
<tr>
  <td><a href="./schemas/tmuxinator/">tmuxinator</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/ui5/">ui5</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/ui5-manifest/">ui5-manifest</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td>uniqueItems</td>
</tr>
<tr>
  <td><a href="./schemas/unreal-engine-uproject/">unreal-engine-uproject</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td>uniqueItems</td>
</tr>
<tr>
  <td><a href="./schemas/vercel/">vercel</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/yamllint/">yamllint</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/zschema_advanced-mixed/">zschema_advanced-mixed</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>mixed</td>
  <td>uniqueItems</td>
</tr>
<tr>
  <td><a href="./schemas/zschema_advanced-noUniqueItems-mixed/">zschema_advanced-noUniqueItems-mixed</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>mixed</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/zschema_advanced-noUniqueItems-valid/">zschema_advanced-noUniqueItems-valid</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/zschema_advanced-valid/">zschema_advanced-valid</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td>uniqueItems</td>
</tr>
<tr>
  <td><a href="./schemas/zschema_basic-mixed/">zschema_basic-mixed</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>mixed</td>
  <td>uniqueItems</td>
</tr>
<tr>
  <td><a href="./schemas/zschema_basic-noUniqueItems-mixed/">zschema_basic-noUniqueItems-mixed</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>mixed</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/zschema_basic-noUniqueItems-valid/">zschema_basic-noUniqueItems-valid</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td></td>
</tr>
<tr>
  <td><a href="./schemas/zschema_basic-valid/">zschema_basic-valid</a></td>
  <td>TODO</td>
  <td>TODO</td>
  <td>valid only</td>
  <td>uniqueItems</td>
</tr>
</table>