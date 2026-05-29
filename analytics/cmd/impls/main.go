// Copyright 2026 Walter Schulze
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/katydid/validator-jsonschema-benchmarks/analytics/analytics"
)

func main() {
	log.SetFlags(log.Lshortfile)
	format := flag.String("format", "latex", "output format (md|latex)")
	rmSchema1 := flag.Bool("filterSchema1", false, "filter all schemas where one implementation had an non zero exit code")
	schemasFolder := flag.String("schemasFolder", "./schemas", "location of schemas folder")
	impl := flag.String("impls", "", "space separated list of implementations to filter")
	flag.Parse()
	reportFilename := "./dist/report.csv"
	if len(flag.Args()) == 0 {
		log.Printf("expected location of report.csv as first argument to command, defaulting to ./dist/report.csv")
	} else {
		reportFilename = flag.Args()[0]
	}

	log.Printf("analysing schemas folder at %s\n", *schemasFolder)
	schemas, err := analytics.CollectSchemas(*schemasFolder)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("analysing report file %s\n", reportFilename)
	lines, err := analytics.ReadLines(schemas, reportFilename)
	if err != nil {
		log.Fatal(err)
	}

	if impl != nil && *impl != "" {
		impls := strings.Split(*impl, " ")
		log.Printf("filtering implementations: %#v", impls)
		lines = analytics.FilterImplementations(lines, impls)
	}

	if *rmSchema1 {
		lines = analytics.FilterSchemasExitStatus0(lines)
	} else {
		lines = analytics.FilterExitStatus0(lines)
	}

	lines = analytics.AverageRuns(lines)

	groups := analytics.GroupBySchema(lines)
	scores := analytics.ScoreGroups(groups)
	analysed := analytics.AnalyseImplementations(scores)

	switch *format {
	case "latex":
		fprintLatex(os.Stdout, analysed)
	case "md":
		fprintMarkdown(os.Stdout, analysed)
	}
}

func fprintMarkdown(
	w io.Writer,
	impls []*analytics.Implementation,
) {
	p := func(format string, a ...any) {
		fmt.Fprintf(w, format, a...)
	}

	p(`| impl | warm avg/doc | warm mean/doc | warm avg rank | warm mean rank | cold avg/doc | cold mean/doc | cold avg rank | cold mean rank |`)
	p("\n")
	p(`| ---  | ---          | ---           | ---           | ---            | ---          | ---           | ---           | --- | `)
	p("\n")
	for _, impl := range impls {
		p("| ")
		p("%s", impl.Name)
		p(" | ")
		p("%.0f", impl.AverageWarmNsPerDoc)
		p(" | ")
		p("%.0f", impl.MedianWarmNsPerDoc)
		p(" | ")
		p("%.0f", impl.AverageWarmRank)
		p(" | ")
		p("%.0f", impl.MedianWarmRank)
		p(" | ")
		p("%.0f", impl.AverageColdNsPerDoc)
		p(" | ")
		p("%.0f", impl.MedianColdNsPerDoc)
		p(" | ")
		p("%.0f", impl.AverageColdRank)
		p(" | ")
		p("%.0f", impl.MedianColdRank)
		p(" |")
		p("\n")
	}
}

func fprintLatex(
	w io.Writer,
	impls []*analytics.Implementation,
) {
	p := func(format string, a ...any) {
		fmt.Fprintf(w, format, a...)
	}

	p("%% BEGIN Generated tabular\n")
	p("\\begin{tabular}{l|llll}\n")
	p(`impl & warm avg/doc & warm mean/doc & cold avg/doc & cold mean/doc \\`)
	p("\n")
	p(`\hline`)
	p("\n")
	for _, impl := range impls {
		p("%s", impl.Name)
		p(" & ")
		p("%.0f", impl.AverageWarmNsPerDoc)
		p(" & ")
		p("%.0f", impl.MedianWarmNsPerDoc)
		p(" & ")
		p("%.0f", impl.AverageColdNsPerDoc)
		p(" & ")
		p("%.0f", impl.MedianColdNsPerDoc)
		p(" \\\\\n")

	}
	p("\\end{tabular}\n")
	p("%% END Generated tabular\n")
}
