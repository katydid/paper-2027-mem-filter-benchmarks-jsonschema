import math
import json
import os
import pathlib
import sys
import logging
import time

import jsonschema

WARMUP_ITERATIONS = 1000
MAX_WARMUP_TIME = 1e9 * 10

if __name__ == "__main__":

    logging.basicConfig()
    log = logging.getLogger("py-jsonschema")
    log.setLevel(logging.DEBUG)

    example_dir = pathlib.Path(sys.argv[1])
    log.debug(example_dir)
    schema = json.load(open(example_dir / "schema.json"))
    want = not ("-invalid" in example_dir.__str__())
    lines = open(example_dir / "instances.jsonl").readlines()

    parse_start = time.time_ns()
    instances = [json.loads(doc) for doc in lines]
    parse_end = time.time_ns()

    Validator = jsonschema.validators.validator_for(schema)
    compile_start = time.time_ns()
    validator = Validator(schema, format_checker=jsonschema.Draft202012Validator.FORMAT_CHECKER)
    compile_end = time.time_ns()

    cold_start = time.time_ns()
    for instance in instances:
        got = validator.is_valid(instance)
        if got != want:
            print(instance, file=sys.stderr)
            exit(1)
    cold_end = time.time_ns()

    iterations = math.ceil(MAX_WARMUP_TIME / (cold_end - cold_start))
    for _ in range(min(iterations, WARMUP_ITERATIONS)):
        for instance in instances:
            validator.is_valid(instance)

    warm_start = time.time_ns()
    for instance in instances:
        validator.is_valid(instance)
    warm_end = time.time_ns()

    print((cold_end - cold_start), ",", (warm_end - warm_start), ",", (parse_end - parse_start), ",", (compile_end - compile_start), sep='')
