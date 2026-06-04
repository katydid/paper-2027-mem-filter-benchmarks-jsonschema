#include <jsoncons/json.hpp>
#include <jsoncons_ext/jsonschema/jsonschema.hpp>

#include <chrono>
#include <filesystem>
#include <fstream>
#include <iostream>
#include <vector>
#include <filesystem>
#include <sstream>
#include <string>

#define WARMUP_ITERATIONS 100L
#define MAX_WARMUP_TIME 10000000000

namespace fs = std::filesystem;

using jsoncons::json;
namespace jsonschema = jsoncons::jsonschema;

template <typename Json>
bool validate_all(const jsonschema::json_schema<Json> &compiled, const std::vector<json> &instances, bool want) {
  for (const auto &instance : instances) {
    try
    {
      compiled.validate(instance);
      if (!want) {
        return false;
      }
    }
    catch (const std::exception& e)
    {
      if (want) {
        return false;
      }
    }
  }
  return true;
}

std::string read_file(const fs::path &path) {
  std::ifstream f(path, std::ios::binary);
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

int validate(const std::filesystem::path &example) {
  std::ifstream input_schema((example / "schema.json").string());
  const auto schema = json::parse(input_schema);
  const bool want = !(example.string().find("-invalid") != std::string::npos);
  std::cerr << std::boolalpha;
  std::cerr << "want:" << want << "\n";

  const auto compile_start{std::chrono::high_resolution_clock::now()};
  const auto compiled = jsonschema::make_json_schema(schema, jsonschema::evaluation_options{}.require_format_validation(true));
  const auto compile_end{std::chrono::high_resolution_clock::now()};
  const auto compile_duration{std::chrono::duration_cast<std::chrono::nanoseconds>(
      compile_end - compile_start)};

  std::string instance_str = read_file(example / "instances.jsonl");
  std::vector<std::string> instance_lines = split_string_by_newline(instance_str);

  const auto parse_start{std::chrono::high_resolution_clock::now()};
  std::vector<json> instances;
  for (const auto &instance_line : instance_lines) {
    const auto instance = json::parse(instance_line);
    instances.push_back(instance);
  }
  const auto parse_end{std::chrono::high_resolution_clock::now()};
  const auto parse_duration{std::chrono::duration_cast<std::chrono::nanoseconds>(
      parse_end - parse_start)};

  const auto cold_start{std::chrono::high_resolution_clock::now()};
  if (!validate_all(compiled, instances, want)) {
    return EXIT_FAILURE;
  }
  const auto cold_end{std::chrono::high_resolution_clock::now()};
  const auto cold_duration{std::chrono::duration_cast<std::chrono::nanoseconds>(
      cold_end - cold_start)};

  const auto iterations = 1 + ((MAX_WARMUP_TIME - 1) / cold_duration.count());
  for (int i = 0; i < std::min(iterations, WARMUP_ITERATIONS); i++) {
    validate_all(compiled, instances, want);
  }

  const auto warm_start{std::chrono::high_resolution_clock::now()};
  validate_all(compiled, instances, want);
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
    std::cerr << "Error during jsoncons benchmark: " << e.what() << "\n";
    return EXIT_FAILURE;
  }
}
