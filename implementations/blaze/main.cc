#include <sourcemeta/core/io.h>
#include <sourcemeta/core/json.h>
#include <sourcemeta/core/jsonl.h>
#include <sourcemeta/core/jsonschema.h>

#include <sourcemeta/blaze/compiler.h>
#include <sourcemeta/blaze/evaluator.h>

#include <chrono>
#include <filesystem>
#include <iostream>
#include <vector>
#include <sstream>
#include <string>

#define WARMUP_ITERATIONS 100L
#define MAX_WARMUP_TIME 10000000000

namespace fs = std::filesystem;

std::string read_file(const fs::path &path) {
  std::ifstream f(path, std::ios::binary);
  if (!f) {
    throw std::runtime_error(std::format("Cannot open file: {}", path.string()));
  }
  return {std::istreambuf_iterator<char>{f}, {}};
}

std::vector<std::string> split_string_by_newline(const std::string& str)
{
    auto result = std::vector<std::string>{};
    auto ss = std::stringstream{str};
    for (std::string line; std::getline(ss, line, '\n');) {
        result.push_back(line);
    }
    return result;
}

bool validate_all(auto &evaluator, const auto &instances, const auto &schema_template, bool want) {
  for (std::size_t num = 0; num < instances.size(); num++) {
    const bool result = evaluator.validate(schema_template, instances[num]);
    if (result != want) {
      std::cerr << "Error validating instance " << num << "\n";
      return false;
    }
  }

  return true;
}

int validate(const std::filesystem::path &example) {
  const auto schema{
      sourcemeta::core::read_json(example / "schema-noformat.json")};
  std::string instance_str = read_file(example / "instances.jsonl");
  const bool want = !example.string().contains("-invalid");
  std::cerr << std::boolalpha;
  std::cerr << "want:" << want << "\n";
  std::vector<std::string> instance_lines = split_string_by_newline(instance_str);

  const auto parse_start{std::chrono::high_resolution_clock::now()};
  std::vector<sourcemeta::core::JSON> instances;
  for (const auto &instance_line : instance_lines) {
    sourcemeta::core::JSON instance = sourcemeta::core::parse_json(instance_line);
    instances.push_back(instance);
  }
  const auto parse_end{std::chrono::high_resolution_clock::now()};
  const auto parse_duration{std::chrono::duration_cast<std::chrono::nanoseconds>(
      parse_end - parse_start)};

  const auto compile_start{std::chrono::high_resolution_clock::now()};
  const auto schema_template{sourcemeta::blaze::compile(
      schema, sourcemeta::core::schema_walker,
      sourcemeta::core::schema_resolver,
      sourcemeta::blaze::default_schema_compiler,
      sourcemeta::blaze::Mode::FastValidation)};

  const auto compile_end{std::chrono::high_resolution_clock::now()};
  const auto compile_duration{std::chrono::duration_cast<std::chrono::nanoseconds>(
      compile_end - compile_start)};

  sourcemeta::blaze::Evaluator evaluator;

  const auto cold_start{std::chrono::high_resolution_clock::now()};
  if (!validate_all(evaluator, instances, schema_template, want)) {
    return EXIT_FAILURE;
  }
  const auto cold_end{std::chrono::high_resolution_clock::now()};
  const auto cold_duration{std::chrono::duration_cast<std::chrono::nanoseconds>(
      cold_end - cold_start)};

  const auto iterations = 1 + ((MAX_WARMUP_TIME - 1) / cold_duration.count());
  for (int i = 0; i < std::min(iterations, WARMUP_ITERATIONS); i++) {
    validate_all(evaluator, instances, schema_template, want);
  }

  const auto warm_start{std::chrono::high_resolution_clock::now()};
  validate_all(evaluator, instances, schema_template, want);
  const auto warm_end{std::chrono::high_resolution_clock::now()};
  const auto warm_duration{std::chrono::duration_cast<std::chrono::nanoseconds>(
      warm_end - warm_start)};

  std::cout << cold_duration.count() << "," << warm_duration.count() << "," << parse_duration.count() << "," << compile_duration.count() << "\n";

  return EXIT_SUCCESS;
}

int main(int argc, char **argv) {
  if (argc < 2) {
    std::cerr << "Usage: " << argv[0] << " <schema>\n";
    return EXIT_FAILURE;
  }

  try {
    const std::filesystem::path example{argv[1]};
    return validate(example);
  } catch (const std::exception &e) {
    std::cerr << "Error during Blaze benchmark: " << e.what() << " in folder " << argv[1] << "\n";
    return EXIT_FAILURE;
  }
}
