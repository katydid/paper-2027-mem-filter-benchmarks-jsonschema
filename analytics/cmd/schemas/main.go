package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type schema struct {
	name                   string
	features               []string
	source                 string
	generated              bool
	schemaSizeBytes        int
	numInstances           int
	avgInstanceSizeBytes   float64
	rmUniqueItems          bool
	hasExistingReplacement bool
}

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
	dirs, err := os.ReadDir(schemasFolder)
	if err != nil {
		log.Fatalf("problem reading folder %s got error: %v", schemasFolder, err)
	}
	schemas := []*schema{}
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		path := filepath.Join(schemasFolder, dir.Name())
		log.Printf("collecting analytics at %s", path)
		s, err := collectSchema(path)
		if err != nil {
			log.Fatal(err)
		}
		schemas = append(schemas, s)
		for _, schema := range schemas {
			if containsRmUniqueItems(schemas, schema.name) {
				schema.hasExistingReplacement = true
			}
		}
	}
	if *rmUniqueItems {
		schemas = removeUniqueItems(schemas)
	}
	if *rmSource {
		schemas = removeSourcePrefix(schemas)
	}
	if *groupGen {
		schemas = groupGenerated(schemas)
	}
	switch *format {
	case "html":
		fprintHTML(os.Stdout, schemas)
	case "latex":
		fprintLatex(os.Stdout, schemas, *latexPifont)
	}
}

func collectSchema(folder string) (*schema, error) {
	s := &schema{}
	s.name = filepath.Base(folder)
	if strings.Contains(s.name, "rmUniqueItems") {
		s.rmUniqueItems = true
	}
	files, err := os.ReadDir(folder)
	if err != nil {
		log.Printf("error reading folder: %s", folder)
		return nil, err
	}
	for _, file := range files {
		basename := file.Name()
		filename := filepath.Join(folder, basename)
		data, err := os.ReadFile(filename)
		if err != nil {
			log.Printf("error reading file: %s", filename)
			return nil, err
		}
		switch basename {
		case "schema.json":
			s.schemaSizeBytes = len(data)
			s.features = collectFeatures(data)
		case "source.txt":
			s.source = string(data)
		case ".generated":
			s.generated = true
		case "instances.jsonl":
			instances := bytes.Split(data, []byte("\n"))
			s.numInstances = len(instances)
			// last line might be empty, remember to not count that
			if len(bytes.TrimSpace(instances[len(instances)-1])) == 0 {
				s.numInstances = s.numInstances - 1
			}
			s.avgInstanceSizeBytes = float64(len(data)) / float64(len(instances))
		}
	}
	return s, nil
}

func collectFeatures(data []byte) []string {
	features := []string{}
	if bytes.Contains(data, []byte(`"uniqueItems": true`)) || bytes.Contains(data, []byte(`"uniqueItems":true`)) {
		features = append(features, "uniqueItems")
	}
	return features
}

func fprintHTML(w io.Writer, schemas []*schema) {
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
		p("<!-- name --><td>%s</td>", schema.name)
		p("<!-- schema size (bytes) --><td>%d</td>", schema.schemaSizeBytes)
		p("<!-- number of docs --><td>%d</td>", schema.numInstances)
		p("<!-- avg doc size (bytes) --><td>%.0f</td>", schema.avgInstanceSizeBytes)
		out()
		p("</tr>")
	}
	out()
	p("</table>")
}

func fprintLatex(w io.Writer, schemas []*schema, pifont bool) {
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
		p("%s", esc(schema.name))
		p(" & ")
		p("%d", schema.schemaSizeBytes)
		p(" & ")
		p("%d", schema.numInstances)
		p(" & ")
		p("%.0f", schema.avgInstanceSizeBytes)
		p(" & ")
		p("%s", sprintBool(schema.generated))
		p(" & ")
		p("%s", sprintBool(schema.rmUniqueItems))
		p(" \\\\\n")

	}
	p("\\end{tabular}\n")
	p("%% END Generated tabular\n")
}

func removeUniqueItems(schemas []*schema) []*schema {
	res := []*schema{}
	for i, schema := range schemas {
		if schema.hasExistingReplacement {
			continue
		}
		if strings.Contains(schema.name, "rmUniqueItems") {
			s := schema
			s.name = strings.Replace(s.name, "-rmUniqueItems", "", 1)
			res = append(res, s)
			continue
		}
		res = append(res, schemas[i])
	}
	return res
}

func containsRmUniqueItems(schemas []*schema, name string) bool {
	if strings.HasSuffix(name, "-mixed") {
		name = name[:len(name)-6] + "-rmUniqueItems-mixed"
	} else if strings.HasSuffix(name, "-valid") {
		name = name[:len(name)-6] + "-rmUniqueItems-valid"
	} else {
		name = name + "-rmUniqueItems"
	}
	for _, schema := range schemas {
		if schema.name == name {
			return true
		}
	}
	return false
}

func removeSourcePrefix(schemas []*schema) []*schema {
	for i := range schemas {
		if strings.HasPrefix(schemas[i].name, "ajv-") {
			schemas[i].name = strings.Replace(schemas[i].name, "ajv-", "", 1)
		}
		if strings.HasPrefix(schemas[i].name, "jsck-") {
			schemas[i].name = strings.Replace(schemas[i].name, "jsck-", "", 1)
		}
		if strings.HasPrefix(schemas[i].name, "example-") {
			schemas[i].name = strings.Replace(schemas[i].name, "example-", "", 1)
		}
		if strings.HasPrefix(schemas[i].name, "zschema-") {
			schemas[i].name = strings.Replace(schemas[i].name, "zschema-", "", 1)
		}
		if strings.HasPrefix(schemas[i].name, "katydid-") {
			schemas[i].name = strings.Replace(schemas[i].name, "katydid-", "", 1)
		}
	}
	return schemas
}

func groupGenerated(schemas []*schema) []*schema {
	res := []*schema{}
	for i := range schemas {
		if strings.HasSuffix(schemas[i].name, "-mixed") {
			mixedSchema := schemas[i]
			name := mixedSchema.name[:len(mixedSchema.name)-6]
			validName := name + "-valid"
			validSchema := findSchema(schemas, validName)
			groupSchema := &schema{
				name:                   name,
				features:               mixedSchema.features,
				source:                 mixedSchema.source,
				schemaSizeBytes:        mixedSchema.schemaSizeBytes,
				numInstances:           mixedSchema.numInstances,
				avgInstanceSizeBytes:   (mixedSchema.avgInstanceSizeBytes + validSchema.avgInstanceSizeBytes) / 2,
				rmUniqueItems:          mixedSchema.rmUniqueItems,
				hasExistingReplacement: mixedSchema.hasExistingReplacement,
				generated:              true,
			}
			res = append(res, groupSchema)
		} else if strings.HasSuffix(schemas[i].name, "-valid") {
			// ignore, grouping happens in mixed
			continue
		} else {
			res = append(res, schemas[i])
		}
	}
	return res
}

func findSchema(schemas []*schema, name string) *schema {
	for i := range schemas {
		if schemas[i].name == name {
			return schemas[i]
		}
	}
	return nil
}
