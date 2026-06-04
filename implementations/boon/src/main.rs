use std::{error::Error, fs::File, io::{BufReader, BufRead}, time::Instant};
use boon::{Compiler, Schemas, SchemaIndex};
use serde_json::Value;
use std::env;

const WARMUP_ITERATIONS: u128 = 100;
const MAX_WARMUP_TIME: u128 = 10_000_000_000; // 10 seconds

fn validate_all(schemas: &Schemas, sch_index: SchemaIndex, serde_lines: &std::vec::Vec<Value>, want: bool) -> bool {
  let mut failed: bool = false;
  for (i, line) in serde_lines.iter().enumerate() {
    let result = schemas.validate(&line, sch_index);
    if !result.is_ok() {
      failed = true;
    }
    if want {
      assert!(result.is_ok(), "Validation failed for line {} with err {:?}: {}", i, result, line);
    } else {
      assert!(result.is_err(), "Validation succeeded for line {}: {}", i, line);
    }
  }
  return !failed;
}

fn main() -> Result<(), Box<dyn Error>> {
  // Get arguments
  let args: Vec<String> = env::args().collect();
  let example_folder = &args[1];
  let want = !example_folder.contains("-invalid");

  // Get the schema and instance paths
  let schema_file =   std::fs::canonicalize(example_folder.to_owned() + "/schema.json")?;
  let instance_file = std::fs::canonicalize(example_folder.to_owned() + "/instances.jsonl")?;

  // Read the instance file
  let file = File::open(&instance_file)?;
  let reader = BufReader::new(file);

  // Compile the schema
  let mut schemas = Schemas::new();
  let mut compiler = Compiler::new();

  let compile_start = Instant::now();
  let sch_index = compiler.compile(schema_file.to_str().ok_or("NULL")?, &mut schemas)?;
  let compile_duration = compile_start.elapsed().as_nanos();

  // Serialize instance lines
  let mut serde_lines = std::vec::Vec::new();
  for line in reader.lines() {
      let line = line?;
      serde_lines.push(line);
  }

  let parse_start = Instant::now();
  let mut instances = std::vec::Vec::new();
  for line in serde_lines {
      let instance: Value = serde_json::from_str(&line)?;
      instances.push(instance);
  }
  let parse_duration = parse_start.elapsed().as_nanos();

  // Validate the instances
  let cold_start = Instant::now();
  validate_all(&schemas, sch_index, &instances, want);
  let cold_duration = cold_start.elapsed().as_nanos();

  // Warmup
  let iterations: u128 = MAX_WARMUP_TIME / cold_duration;
  for _ in 0..std::cmp::min(iterations, WARMUP_ITERATIONS) {
    validate_all(&schemas, sch_index, &instances, want);
  }

  let warm_start = Instant::now();
  validate_all(&schemas, sch_index, &instances, want);
  let warm_duration = warm_start.elapsed().as_nanos();

  println!("{:?},{:?},{:?},{:?}", cold_duration, warm_duration, parse_duration, compile_duration);

  Ok(())
}
