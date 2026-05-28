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
	name                 string
	features             []string
	source               string
	generated            bool
	schemaSizeBytes      int
	numInstances         int
	avgInstanceSizeBytes float64
}

func main() {
	log.SetFlags(log.Lshortfile)
	format := flag.String("format", "html", "output format (html|latex)")
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
	}
	switch *format {
	case "html":
		fprintHTML(os.Stdout, schemas)
	case "latex":
		fprintLatex(os.Stdout, schemas)
	}
}

func collectSchema(folder string) (*schema, error) {
	s := &schema{}
	s.name = filepath.Base(folder)
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

func fprintLatex(w io.Writer, schemas []*schema) {
	p := func(format string, a ...any) {
		fmt.Fprintf(w, format, a...)
	}
	//escape latex reserved characters
	esc := func(s string) string {
		return strings.Replace(s, "_", "\\_", -1)
	}
	p("%% BEGIN Generated tabular\n")
	p("\\begin{tabular}{l|llll}\n")
	p("name & schema size (bytes) & number of docs & avg doc size (bytes) \\\\\n")
	p("\\hline\n")
	for _, schema := range schemas {
		p("%s & ", esc(schema.name))
		p("%d & ", schema.schemaSizeBytes)
		p("%d & ", schema.numInstances)
		p("%.0f \\\\\n", schema.avgInstanceSizeBytes)
	}
	p("\\end{tabular}")
	p("%% END Generated tabular\n")
}
