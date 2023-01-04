package main

import (
	"log"
	"os"
	"parseFunannotate/helpers"

	"github.com/akamensky/argparse"
)

var pro string = "A simple tool to parse the annotation output from Funannotate and generate a topGO compatible TSV for GO enrichment."

func main() {

	parser := argparse.NewParser(
		"msaSummary",
		pro,
	)

	// Arguments
	in := parser.String(
		"i",
		"in",
		&argparse.Options{
			Required: true,
			Help:     "Filepath to Funannotate '<species>.annotations.txt' file.",
		},
	)

	out := parser.String(
		"o",
		"out",
		&argparse.Options{
			Required: true,
			Help:     "Output filepath for tab separated GO Term file.",
		},
	)

	// Parse arguments
	err := parser.Parse(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	goTerms := helpers.GetGO(*in)
	helpers.WriteGoMap(goTerms, *out)
}
