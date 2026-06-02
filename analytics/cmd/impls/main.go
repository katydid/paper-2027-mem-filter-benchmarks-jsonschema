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

	totalSchemas := analytics.CountTotalSchemas(lines)
	completedPerImpl := analytics.CountCompletedPerImplementation(lines)

	if *rmSchema1 {
		lines = analytics.FilterSchemasExitStatus0(lines)
	} else {
		lines = analytics.FilterExitStatus0(lines)
	}

	lines = analytics.AverageRuns(lines)

	kindGroups := analytics.GroupByKind(lines)
	analyseds := [][]*analytics.Implementation{}
	for i := range kindGroups {
		kind := kindGroups[i][0].Schema.GeneratedKind
		if kind == "" {
			kind = "valid"
		}
		log.Printf("analysing kind group %d = %v with %d lines", i, kind, len(kindGroups[i]))
		groups := analytics.GroupBySchema(kindGroups[i])
		log.Printf("got %d groups", len(groups))
		scores := analytics.ScoreGroups(groups)
		log.Printf("got %d scores", len(scores))
		analysed := analytics.AnalyseImplementations(scores, totalSchemas, completedPerImpl)
		if len(analysed[0].Scores) == 0 {
			log.Fatalf("no scores for kind %v", kind)
		}
		analyseds = append(analyseds, analysed)
	}

	switch *format {
	case "latex":
		fprintLatex(os.Stdout, analyseds)
	case "md":
		fprintMarkdown(os.Stdout, analyseds)
	}
}

func getKind(impls []*analytics.Implementation) string {
	for _, impl := range impls {
		if len(impl.Scores) == 0 {
			log.Fatalf("no scores for %s", impl.Name)
		}
		for _, score := range impl.Scores {
			kind := score.Line.Schema.GeneratedKind
			if kind == "" {
				return "valid"
			}
			return kind
		}
	}
	log.Fatalf("could not find kind in %v", impls)
	return ""
}

func fprintMarkdown(
	w io.Writer,
	implss [][]*analytics.Implementation,
) {
	p := func(format string, a ...any) {
		fmt.Fprintf(w, format, a...)
	}

	for i, impls := range implss {
		p(`## Exclude Parse Time (%s)`, getKind(implss[i]))
		p("\n")
		p("\n")
		p(`| impl | completed | avg/doc | mean/doc | avg rank | mean rank |`)
		p("\n")
		p(`| ---  | ---      | ---     | ---      | ---      | ---       | `)
		p("\n")
		for _, impl := range impls {
			p("| ")
			p("%s", impl.Name)
			p(" | ")
			p("%d/%d", impl.CompletedSchemas, impl.TotalSchemas)
			p(" | ")
			p("%.0f", impl.AverageWarmNsPerDoc)
			p(" | ")
			p("%.0f", impl.MedianWarmNsPerDoc)
			p(" | ")
			p("%.0f", impl.AverageWarmRank)
			p(" | ")
			p("%.0f", impl.MedianWarmRank)
			p(" |")
			p("\n")
		}

		p("\n")
		p(`## Include Parse Time (%s)`, getKind(implss[i]))
		p("\n")
		p("\n")
		p(`| impl | completed    | avg/doc  | mean/doc | avg rank | mean rank |`)
		p("\n")
		p(`| ---  | ---          | ---      | ---      | ---      | ---       |`)
		p("\n")
		for _, impl := range impls {
			if !impl.ParseTODO {
				p("| ")
				p("%s", impl.Name)
				p(" | ")
				p("%d/%d", impl.CompletedSchemas, impl.TotalSchemas)
				p(" | ")
				p("%.0f", impl.AverageParsePlusWarmNsPerDoc)
				p(" | ")
				p("%.0f", impl.MedianParsePlusWarmNsPerDoc)
				p(" | ")
				p("%.0f", impl.AverageParsePlusWarmRank)
				p(" | ")
				p("%.0f", impl.MedianParsePlusWarmRank)
				p(" |")
				p("\n")
			} else {
				p("| ")
				p("%s", impl.Name)
				p(" | ")
				p("TODO")
				p(" | ")
				p("TODO")
				p(" | ")
				p("TODO")
				p(" | ")
				p("TODO")
				p(" | ")
				p("TODO")
				p(" |")
				p("\n")
			}
		}
		p("\n")
	}
}

func fprintLatex(
	w io.Writer,
	implss [][]*analytics.Implementation,
) {
	p := func(format string, a ...any) {
		fmt.Fprintf(w, format, a...)
	}

	for i, impls := range implss {
		p("%% BEGIN Generated tabular for kind: %s\n", getKind(implss[i]))
		p("\\begin{tabular}{l|ll}\n")
		p(`impl & avg/doc & mean/doc \\`)
		p("\n")
		p(`\hline`)
		p("\n")
		for _, impl := range impls {
			p("%s", impl.Name)
			p(" & ")
			p("%.0f", impl.AverageWarmNsPerDoc)
			p(" & ")
			p("%.0f", impl.MedianWarmNsPerDoc)
			p(" \\\\\n")

		}
		p("\\end{tabular}\n")
		p("%% END Generated tabular\n")
	}
}
