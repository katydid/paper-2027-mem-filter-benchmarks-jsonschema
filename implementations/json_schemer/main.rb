require 'json_schemer'
require 'json'


WARMUP_ITERATIONS = 100
MAX_WARMUP_TIME = 1e9 * 10

def validate_all(instances, schemer, want)
  instances.each do |instance|
    res = schemer.valid?(instance)
    if (res != want) then exit! end
  end
end

path = ARGV[0]

# Load the schema and build a validator
schema = JSON.parse(File.read(File.join(path, "schema.json")))
want = !(path.include? "-invalid")

compile_start = Process.clock_gettime(Process::CLOCK_REALTIME, :nanosecond)
schemer = JSONSchemer.schema(schema, format: true)
compile_end = Process.clock_gettime(Process::CLOCK_REALTIME, :nanosecond)

lines = File.readlines(File.join(path, "instances.jsonl"))

# Read all instances into an array
parse_start = Process.clock_gettime(Process::CLOCK_REALTIME, :nanosecond)
instances = lines.map do |line|
  JSON.parse(line)
end
parse_end = Process.clock_gettime(Process::CLOCK_REALTIME, :nanosecond)

# Run the validation
cold_start = Process.clock_gettime(Process::CLOCK_REALTIME, :nanosecond)
validate_all(instances, schemer, want)
cold_end = Process.clock_gettime(Process::CLOCK_REALTIME, :nanosecond)

iterations = (MAX_WARMUP_TIME / (cold_end - cold_start)).ceil

[WARMUP_ITERATIONS, iterations].min.times do
  validate_all(instances, schemer, want)
end

warm_start = Process.clock_gettime(Process::CLOCK_REALTIME, :nanosecond)
validate_all(instances, schemer, want)
warm_end = Process.clock_gettime(Process::CLOCK_REALTIME, :nanosecond)

print (cold_end - cold_start), ",", (warm_end - warm_start), ",", (parse_end - parse_start), ",", (compile_end - compile_start), "\n"
