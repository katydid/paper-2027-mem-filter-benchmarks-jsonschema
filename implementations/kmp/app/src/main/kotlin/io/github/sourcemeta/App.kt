package io.github.sourcemeta


import io.github.optimumcode.json.schema.OutputCollector
import io.github.optimumcode.json.schema.JsonSchema
import io.github.optimumcode.json.schema.JsonSchemaLoader
import io.github.optimumcode.json.schema.SchemaOption
import io.github.optimumcode.json.schema.FormatBehavior.ANNOTATION_AND_ASSERTION
import io.github.optimumcode.json.schema.ValidationError
import java.io.File
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.JsonElement


val WARMUP_ITERATIONS: ULong = 1000.toULong()
val MAX_WARMUP_TIME: ULong = (1e9 * 10).toULong()


fun validateAll(schema: JsonSchema, docs: List<JsonElement>, want: Boolean): Boolean {
  for (doc in docs) {
    val got = schema.validate(doc, OutputCollector.flag()).valid
    if (got != want) {
      System.err.println("error for ${want} with ${doc}")
      return false
    }
  }
  return true
}

fun main(args: Array<String>) {
    val json = Json { ignoreUnknownKeys = true }

    // Prepare the schema
    val want = !args[0].contains("-invalid")
    val schemaDefinition = File(args[0]).readText()
    val compileStart = System.nanoTime()
    val schema = JsonSchemaLoader.create().withSchemaOption(SchemaOption.FORMAT_BEHAVIOR_OPTION, ANNOTATION_AND_ASSERTION).fromDefinition(schemaDefinition)
    val compileEnd = System.nanoTime()

    // Load all documents
    var lines = File(args[1]).readLines()
    val parseStart = System.nanoTime()
    val docs = lines.map { json.parseToJsonElement(it) }
    val parseEnd = System.nanoTime()

    val coldStart = System.nanoTime()
    val valid = validateAll(schema, docs, want)
    val coldEnd = System.nanoTime()

    if (!valid) {
        System.exit(1)
    }

    // Run some warmup iterations
    val iterations: ULong = kotlin.math.ceil(MAX_WARMUP_TIME.toDouble() / (coldEnd - coldStart)).toULong()
    repeat(kotlin.math.min(iterations, WARMUP_ITERATIONS).toInt()) {
        validateAll(schema, docs, want)
    }

    val warmStart = System.nanoTime()
    validateAll(schema, docs, want)
    val warmEnd = System.nanoTime()

    println("${coldEnd - coldStart},${warmEnd - warmStart},${parseEnd - parseStart},${compileEnd - compileStart}")
}
