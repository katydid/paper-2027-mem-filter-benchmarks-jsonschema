import fs from 'fs';
import readline from 'readline';
import { performance } from 'perf_hooks';

const DRAFTS = {
  "https://json-schema.org/draft/2020-12/schema": (await import("ajv/dist/2020.js")).Ajv2020,
  "https://json-schema.org/draft/2019-09/schema": (await import("ajv/dist/2019.js")).Ajv2019,
  "http://json-schema.org/draft-07/schema": (await import("ajv")).Ajv,
};
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

function readFile(filePath) {
  try {
    const fileContent = fs.readFileSync(filePath, 'utf8');
    return fileContent;
  } catch (error) {
    process.exit(1);
  }
}

function validateAll(instances, validator, want) {
  for (const instance of instances) {
    const result = validator(instance);
    if (result != want) {
      return true;
    }
  }
  return false;
}

async function validateSchema(schemaPath, instancePath, want) {
  const schema = readJSONFile(schemaPath);

  const ajv = new DRAFTS[schema['$schema'].replace(/#$/, '')]({strict: false});

  const compileStart = performance.now();
  const validate = ajv.compile(schema);
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
