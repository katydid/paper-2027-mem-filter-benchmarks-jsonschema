using Corvus.Json;
using System.Diagnostics;
using System.Linq;

const int WarmupIterations = 1000;
const long MaxWarmupTime = 10_000_000_000;

bool ValidateAll(JSB.Schema[] docs, bool want) {
  var valid = true;
  foreach (var doc in docs) {
    var result = doc.Validate(ValidationContext.ValidContext, ValidationLevel.Flag);
    if (result.IsValid != want) {
      valid = false;  
    }
  }

  return valid;
}


// Read and parse all instances
var lines = File.ReadLines(args[0]);
var want = !args[0].Contains("-invalid");
Stopwatch stopWatch = new Stopwatch();

stopWatch.Start();
JSB.Schema[] docs = Array.Empty<JSB.Schema>();
try  {
  docs = lines.Select(l => JSB.Schema.Parse(l)).ToArray();
} catch (System.Text.Json.JsonException e) {
  Environment.Exit(1);
}
stopWatch.Stop();
TimeSpan parseTs = stopWatch.Elapsed;

// Loop and validate all instances
stopWatch.Start();
var valid = ValidateAll(docs, want);
stopWatch.Stop();
TimeSpan coldTs = stopWatch.Elapsed;

var iterations = (int) Math.Ceiling(((double) MaxWarmupTime) / coldTs.TotalNanoseconds);
for (int i = 0; i < Math.Min(iterations, WarmupIterations); i++) {
  ValidateAll(docs, want);
}

stopWatch.Restart();
ValidateAll(docs, want);
stopWatch.Stop();
TimeSpan warmTs = stopWatch.Elapsed;

// Output file time and exit
Console.WriteLine(coldTs.TotalNanoseconds + "," + warmTs.TotalNanoseconds + "," + parseTs.TotalNanoseconds);
Environment.Exit(valid ? 0 : 1);
