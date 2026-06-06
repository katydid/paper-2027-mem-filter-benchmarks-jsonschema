import { validator } from '@exodus/schemasafe';
import fs from 'fs';
import readline from 'readline';
import { performance } from 'perf_hooks';

const WARMUP_ITERATIONS = 100;
const MAX_WARMUP_TIME = 1e9 * 10; // 10 seconds

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

function readFile(filePath) {
  try {
    const fileContent = fs.readFileSync(filePath, 'utf8');
    return fileContent;
  } catch (error) {
    process.exit(1);
  }
}

function validateAll(instances, validate, want) {
  for (const instance of instances) {
    const result = validate(instance);
    if (result != want) {
      console.error(instance);
      return true;
    }
  }

  return false;
}

async function validateSchema(schemaPath, instancePath, want) {
  const schema = readJSONFile(schemaPath);

  const compileStart = performance.now();
  let validate;
  try {
    validate = validator(schema, {
      mode: 'lax', // enables formats
      isJSON: true
    });
  } catch (error) {
    console.error(error);
    process.exit(1);
  }
  const compileEnd = performance.now();
  const compileDurationNs = (compileEnd - compileStart) * 1e6;

  const instanceData = readFile(instancePath);
  const instanceDatas = instanceData.split(/\r?\n/);

  const parseStartTime = performance.now();
  const instances = [];
  instanceDatas.forEach(function (item) {
    if (item.length > 0) {
      instances.push(JSON.parse(item));
    }
  })
  const parseEndTime = performance.now();
  const parseDurationNs = (parseEndTime - parseStartTime) * 1e6;

  const coldStartTime = performance.now();
  const failed = validateAll(instances, validate, want);
  const coldEndTime = performance.now();
  const coldDurationNs = (coldEndTime - coldStartTime) * 1e6;

  const iterations = Math.ceil(MAX_WARMUP_TIME / coldDurationNs);
  for (let i = 0; i < Math.min(iterations, WARMUP_ITERATIONS); i++) {
    validateAll(instances, validate, want);
  }

  const warmStartTime = performance.now();
  validateAll(instances, validate, want);
  const warmEndTime = performance.now();
  const warmDurationNs = (warmEndTime - warmStartTime) * 1e6;

  console.log(coldDurationNs.toFixed(0) + ',' + warmDurationNs.toFixed(0) + ',' + parseDurationNs.toFixed(0) + ',' + compileDurationNs.toFixed(0));

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
