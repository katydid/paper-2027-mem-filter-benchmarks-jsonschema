package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/katydid/validator-jsonschema-benchmarks/analytics/analytics"
)

func main() {
	log.SetFlags(log.Lshortfile)
	format := flag.String("format", "html", "output format (html|latex|md)")
	rmUniqueItems := flag.Bool("rmUniqueItems", false, "if there is an rmUniqueItems version replace the original with it")
	latexPifont := flag.Bool("latex.pifont", false, "replace yes/no in with \\cmark/\\xmark, requires adding the following to your latex: \\usepackage{pifont}\\newcommand{\\cmark}{\\ding{51}}\\newcommand{\\xmark}{\\ding{55}}")
	rmSource := flag.Bool("rmSource", false, "remove prefix source from schema name, for example example-x becomes x")
	groupGen := flag.Bool("groupGen", false, "group generated schemas together")
	onlyCurated := flag.Bool("curated", true, "only show curated schemas")
	flag.Parse()
	rootFolder := "."
	if len(flag.Args()) > 0 {
		rootFolder = flag.Args()[0]
	}
	curated := analytics.GetCuratedSchemas(rootFolder)
	schemasFolder := filepath.Join(rootFolder, "schemas")
	log.Printf("analysing schemas folder at %s\n", schemasFolder)
	schemas, err := analytics.CollectSchemas(schemasFolder, curated)
	if err != nil {
		log.Fatal(err)
	}
	if *rmUniqueItems {
		schemas = analytics.RemoveUniqueItems(schemas)
	}
	if *rmSource {
		schemas = analytics.RemoveSourcePrefixFromName(schemas)
	}
	if *groupGen {
		schemas = analytics.GroupGenerated(schemas)
	}
	if *onlyCurated {
		schemas = slices.DeleteFunc(schemas, func(s *analytics.Schema) bool {
			return !s.Curated
		})
	}
	switch *format {
	case "html":
		fprintHTML(os.Stdout, schemas)
	case "latex":
		fprintLatex(os.Stdout, schemas, *latexPifont)
	case "md":
		fprintMarkdown(os.Stdout, schemas)
	}
}

func fprintHTML(w io.Writer, schemas []*analytics.Schema) {
	tabs := ""
	in := func() {
		tabs += "\t"
	}
	out := func() {
		tabs = tabs[:len(tabs)-1]
	}
	p := func(format string, a ...any) {
		fmt.Fprintf(w, tabs+format+"\n", a...)
	}
	p("<table>")
	in()
	p("<tr>")
	in()
	p("<td>%s</td>", "name")
	p("<td>%s</td>", "schema size (bytes)")
	p("<td>%s</td>", "number of docs")
	p("<td>%s</td>", "mean doc size (bytes)")
	out()
	p("</tr>")
	for _, schema := range schemas {
		p("<tr>")
		in()
		p("<!-- name --><td>%s</td>", schema.Name)
		p("<!-- schema size (bytes) --><td>%d</td>", schema.SchemaSizeBytes)
		p("<!-- number of docs --><td>%d</td>", schema.NumInstances)
		p("<!-- mean doc size (bytes) --><td>%.0f</td>", schema.MeanInstanceSizeBytes)
		out()
		p("</tr>")
	}
	out()
	p("</table>")
}

func fprintMarkdown(w io.Writer, schemas []*analytics.Schema) {
	p := func(format string, a ...any) {
		fmt.Fprintf(w, format, a...)
	}
	p("|#|Dataset name|Source|Valid/InValid|uniqueItems|# Docs|Schema Size (KB)|Mean Doc. Size (B)|\n")
	p("|---|---|---|---|---|---|---|---|\n")
	for i, schema := range schemas {
		p("| ")
		p("%d", i+1)
		p(" | ")
		p("%s", schema.Name)
		p(" | ")
		if schema.Generated {
			p("generated")
		} else if schema.Mutated {
			p("mutated")
		} else {
			p("collected")
		}
		p(" | ")
		if schema.ValidationKind == "invalid" {
			p("invalid")
		} else {
			p("valid")
		}
		p(" | ")
		if strings.Contains(schema.Name, "-rmUniqueItems") {
			p("removed")
		} else {
			p("none")
		}
		p(" | ")
		p("%d", schema.NumInstances)
		p(" | ")
		p("%d", schema.SchemaSizeBytes)
		p(" | ")
		p("%.0f", schema.MeanInstanceSizeBytes)
		p(" |")
		p("\n")
	}
}

func fprintLatex(w io.Writer, schemas []*analytics.Schema, pifont bool) {
	p := func(format string, a ...any) {
		fmt.Fprintf(w, format, a...)
	}
	//escape latex reserved characters
	esc := func(s string) string {
		return strings.Replace(s, "_", "\\_", -1)
	}
	sprintBool := func(b bool) string {
		return analytics.SprintLatexBool(pifont, b)
	}
	p("%% BEGIN Generated tabular\n")
	p("\\begin{tabular}{l|lllll}\n")
	p("name & schema size (B) & \\# Docs & mean doc size (B) & gen & rm \\\\\n")
	p("\\hline\n")
	for _, schema := range schemas {
		p("%s", esc(schema.Name))
		p(" & ")
		p("%d", schema.SchemaSizeBytes)
		p(" & ")
		p("%d", schema.NumInstances)
		p(" & ")
		p("%.0f", schema.MeanInstanceSizeBytes)
		p(" & ")
		p("%s", sprintBool(schema.Generated))
		p(" & ")
		p("%s", sprintBool(schema.RmUniqueItems))
		p(" \\\\\n")

	}
	p("\\end{tabular}\n")
	p("%% END Generated tabular\n")
}
