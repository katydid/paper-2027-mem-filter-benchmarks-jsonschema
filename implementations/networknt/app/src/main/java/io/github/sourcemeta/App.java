package io.github.sourcemeta;

import com.networknt.schema.Schema;
import com.networknt.schema.SchemaRegistry;
import com.networknt.schema.dialect.Dialects;
import com.networknt.schema.SchemaRegistryConfig;
import com.networknt.schema.InputFormat;
import com.networknt.schema.OutputFormat;
import com.networknt.schema.regex.GraalJSRegularExpressionFactory;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.List;
import java.util.stream.Collectors;

public class App {
  static int WARMUP_ITERATIONS = 1000;
  static long MAX_WARMUP_TIME = (long) 1e9 * 10;

  public static void main(String[] args) throws IOException {
    org.apache.log4j.BasicConfigurator.configure();
    org.apache.log4j.Logger.getRootLogger().setLevel(org.apache.log4j.Level.ERROR);
    SchemaRegistryConfig schemaRegistryConfig = SchemaRegistryConfig.builder()
      .regularExpressionFactory(GraalJSRegularExpressionFactory.getInstance()).build();
    SchemaRegistry schemaRegistry = SchemaRegistry.withDefaultDialect(Dialects.getDraft202012(),
      builder -> builder.schemaRegistryConfig(schemaRegistryConfig));
    boolean want = !args[0].contains("-invalid");
    if ((args[0].contains("krakend")) || (args[0].contains("ui5-manifest"))) {
      // unable to handle these and throws an exception
      System.exit(1);
    }
    System.err.println("want" + want);
    String schemaString = new String(Files.readAllBytes(Paths.get(args[0])));

    // Register the schema
    Long compileStart = System.nanoTime();
    Schema schema = schemaRegistry.getSchema(schemaString);
    Long compileEnd = System.nanoTime();

    // Load all documents
    List<String> lines = Files.readAllLines(Paths.get(args[1]));

    Long coldStart = System.nanoTime();
    boolean valid = validateAll(schema, lines, want);
    Long coldEnd = System.nanoTime();

    if (!valid) {
      System.exit(1);
    }

    // Warmup
    long iterations = (long) Math.ceil(((double) MAX_WARMUP_TIME) / (coldEnd - coldStart));
    for (int i = 0; i < WARMUP_ITERATIONS; i++) {
      validateAll(schema, lines, want);
    }

    Long warmStart = System.nanoTime();
    validateAll(schema, lines, want);
    Long warmEnd = System.nanoTime();

    System.out.println(
        (coldEnd - coldStart) + "," + (warmEnd - warmStart) + ",0," + (compileEnd - compileStart));
  }

  public static boolean validateAll(Schema schema, List<String> docs, boolean want) {
    for (String doc : docs) {
      boolean got = schema.validate(doc, InputFormat.JSON, OutputFormat.BOOLEAN, executionContext -> {
        /*
        * By default since Draft 2019-09 the format keyword only generates annotations
        * and not assertions.
        */
        executionContext.executionConfig(executionConfig -> executionConfig.formatAssertionsEnabled(true));
      });
      if (want != got) {
        System.err.println("error:" + doc);
        return false;
      }
    }
    return true;
  }
}
