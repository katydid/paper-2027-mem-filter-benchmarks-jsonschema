import { registerSchema, validate } from "@hyperjump/json-schema/draft-2020-12";
import fs from 'fs';
import readline from 'readline';
import { performance } from 'perf_hooks';

const WARMUP_ITERATIONS = 100;
const MAX_WARMUP_TIME = 1e9 * 10; // 10 seconds

await Promise.all([
  import("@hyperjump/json-schema/draft-2019-09"),
  import("@hyperjump/json-schema/draft-07"),
]);

function readJSONFile(filePath) {
  try {
    const fileContent = fs.readFileSync(filePath, 'utf8');
    return JSON.parse(fileContent);
  } catch (error) {
    process.exit(1);
  }
}

async function* readJSONLines(filePath) {
  const rl = readline.createInterface({
    input: fs.createReadStream(filePath),
  });
  for await (const line of rl) {
    yield JSON.parse(line);
  }
}

async function validateAll(instances, schemaId, want) {
  let failed = false;
  for (const instance of instances) {
    const output = await validate(schemaId, instance);
    if (output.valid != want) {
      failed = true;
    }
  }
  return failed;
}

async function validateSchema(schemaPath, instancePath, want) {
  const schema = readJSONFile(schemaPath);

  const schemaId = schema["$id"] || "https://example.com" + schemaPath;

  const compileStart = performance.now();
  registerSchema(schema, schemaId);
  const compiled = await validate(schemaId);
  const compileEnd = performance.now();
  const compileDurationNs = (compileEnd - compileStart) * 1e6;

  const instances = [];
  for await (const instance of readJSONLines(instancePath)) {
    instances.push(instance);
  }

  const coldStartTime = performance.now();
  const failed = await validateAll(instances, schemaId, want);
  const coldEndTime = performance.now();
  const coldDurationNs = (coldEndTime - coldStartTime) * 1e6;

  const iterations = Math.ceil(MAX_WARMUP_TIME / coldDurationNs);
  for (let i = 0; i < Math.min(iterations, WARMUP_ITERATIONS); i++) {
    await validateAll(instances, schemaId, want);
  }

  const warmStartTime = performance.now();
  await validateAll(instances, schemaId, want);
  const warmEndTime = performance.now();
  const warmDurationNs = (warmEndTime - warmStartTime) * 1e6;

  console.log(coldDurationNs.toFixed(0) + ',' + warmDurationNs.toFixed(0) + ',' + 'TODO' + ',' + compileDurationNs.toFixed(0));

  // Exit with non-zero status on validation failure
  if (failed) {
    process.exit(1);
  }
}

if (process.argv.length !== 4) {
  process.exit(1);
}

const schemaPath = process.argv[2];
const instancePath = process.argv[3];

const want = !schemaPath.includes("-invalid");

await validateSchema(schemaPath, instancePath, want);
