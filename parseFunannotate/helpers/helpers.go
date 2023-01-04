package helpers

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// RowInfo Contains key information that is to be parsed from the Funannotate 'Annotation.txt' object
type RowInfo struct {
	GeneID  string   // Funannotate gene ID
	TransID string   // Transcript id
	Name    string   // Gene symbol (if annotated)
	GoTerms []string // GO terms associated with gene
}

// ReadRows Reads the 'annotation.txt' file row-by-row, parsing the key information
func GetGO(f string) []RowInfo {
	sri := []RowInfo{}
	file, err := os.Open(f)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*2024)
	scanner.Scan() // skip first line

	// Iterate over rows
	log.Printf("[GetGO] Obtaining all GO Terms affiliated with each gene.")
	for scanner.Scan() {
		ri := RowInfo{}
		row := strings.Split(scanner.Text(), "\t") // Using Split rather than Field as it retains empty column values

		// Get the information
		ri.GeneID = row[0]
		ri.TransID = row[1]
		ri.Name = row[7]
		for _, term := range strings.Split(row[16], ";") {
			if term != "" {
				s := strings.Split(term, ": ")[1]
				s = strings.Split(s, " - ")[0]
				ri.GoTerms = append(ri.GoTerms, s)
			}
		}

		sri = append(sri, ri)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return sri
}

// WriteGoMap Export a tab separated file in the form of 'GeneID\t<GO Terms>'
func WriteGoMap(goInfo []RowInfo, out string) {
	log.Printf("[WriteGoMap] Writing GO map to %s", out)
	file, err := os.Create(out)

	if err != nil {
		log.Fatalf("[WriteGoMap] Failed to open file '%s' (%v)", out, err)
	}

	defer file.Close()

	// Write rows to file
	for _, og := range goInfo {
		// Concatenate GO terms
		goTerms := strings.Join(og.GoTerms, " ")

		// Write to file
		fmt.Fprintf(file, "%s\t%s\t%s\t%s\n", og.GeneID, og.TransID, og.Name, goTerms)
	}
}
