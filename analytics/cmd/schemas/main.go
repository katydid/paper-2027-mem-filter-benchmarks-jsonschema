package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	analytics "github.com/katydid/validator-jsonschema-benchmarks/analytics/lib"
)

func main() {
	log.SetFlags(log.Lshortfile)
	format := flag.String("format", "html", "output format (html|latex)")
	rmUniqueItems := flag.Bool("rmUniqueItems", false, "if there is an rmUniqueItems version replace the original with it")
	latexPifont := flag.Bool("latex.pifont", false, "replace yes/no in with \\cmark/\\xmark, requires adding the following to your latex: \\usepackage{pifont}\\newcommand{\\cmark}{\\ding{51}}\\newcommand{\\xmark}{\\ding{55}}")
	rmSource := flag.Bool("rmSource", false, "remove prefix source, for example example-x becomes x")
	groupGen := flag.Bool("groupGen", false, "group generated schemas together")
	flag.Parse()
	if len(flag.Args()) == 0 {
		log.Fatalf("expected location of schemas folder as first argument to command")
	}
	schemasFolder := flag.Args()[0]
	log.Printf("analysing schemas folder at %s\n", schemasFolder)
	schemas, err := analytics.CollectSchemas(schemasFolder)
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
	switch *format {
	case "html":
		fprintHTML(os.Stdout, schemas)
	case "latex":
		fprintLatex(os.Stdout, schemas, *latexPifont)
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
	p("<td>%s</td>", "avg doc size (bytes)")
	out()
	p("</tr>")
	for _, schema := range schemas {
		p("<tr>")
		in()
		p("<!-- name --><td>%s</td>", schema.Name)
		p("<!-- schema size (bytes) --><td>%d</td>", schema.SchemaSizeBytes)
		p("<!-- number of docs --><td>%d</td>", schema.NumInstances)
		p("<!-- avg doc size (bytes) --><td>%.0f</td>", schema.AvgInstanceSizeBytes)
		out()
		p("</tr>")
	}
	out()
	p("</table>")
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
		if pifont {
			if b {
				return "\\cmark"
			} else {
				return "\\xmark"
			}
		} else {
			if b {
				return "yes"
			} else {
				return "no"
			}
		}
	}
	p("%% BEGIN Generated tabular\n")
	p("\\begin{tabular}{l|lllll}\n")
	p("name & schema size (B) & \\# Docs & avg doc size (B) & gen & rm \\\\\n")
	p("\\hline\n")
	for _, schema := range schemas {
		p("%s", esc(schema.Name))
		p(" & ")
		p("%d", schema.SchemaSizeBytes)
		p(" & ")
		p("%d", schema.NumInstances)
		p(" & ")
		p("%.0f", schema.AvgInstanceSizeBytes)
		p(" & ")
		p("%s", sprintBool(schema.Generated))
		p(" & ")
		p("%s", sprintBool(schema.RmUniqueItems))
		p(" \\\\\n")

	}
	p("\\end{tabular}\n")
	p("%% END Generated tabular\n")
}
