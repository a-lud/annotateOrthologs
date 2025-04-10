package utility

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// Check if a value is in a slice. Returns 'true' if found, 'false' otherwise
func Contains(qry string, vals []string) bool {
	for _, i := range vals {
		if i == qry {
			return true
		}
	}

	return false
}

// Checks that the user provided ID-Type variable is valid by comparing it to a pre-defined slice of entries.
func CheckIdType(id []string, valid []string) {
	for _, i := range id {
		if !Contains(i, valid) {
			log.Fatalf("Invalid ID-Type provided. Choose from one of the following: %s", strings.Join(valid, " "))
		}
	}
}

// Removes duplicate strings from a slice.
func RemoveDuplicates(strSlice []string) []string {
	log.Printf("[removeDuplicates]\tRemoving duplicate UniProt accessions\n")
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

// Reads in a best-BLAST-hit CSV file generated by 'bestBlast.py' and returns all unique
// UniProt IDs that were present as a string slice.
func GetUniprotIDs(path string) []string {
	log.Printf("[getUniprotIDs]\tExtracting UniProt accessions\n")
	upid := []string{}

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan() // skip header line

	// Iterate over each line -> store UniprotID in slice
	for scanner.Scan() {
		upid = append(upid, strings.Split(scanner.Text(), ",")[2])
	}

	return RemoveDuplicates(upid)
}

// Write a slice of strings to file in CSV format.
func WriteSliceToCsv(ftype string, idtype []string, matches []string, fout string) {
	// Determine column names depending on type of file parsed
	cnames := ""
	switch ftype {
	case "idmapping":
		cnames = "accession,idtype,id"
	case "idmapping_selected":
		cnames = strings.Join(append([]string{"accession"}, idtype...), ",")
	}

	// Output file (temporary CSV used by python script)
	o, err := os.Create(fout)
	if err != nil {
		log.Fatal(err)
	}

	// Create a writer - so we can write the output lines
	w := bufio.NewWriter(o)

	// Add column names to matches slice
	outslice := append([]string{cnames}, matches...)

	for _, current := range outslice {
		w.WriteString(current + "\n")
		w.Flush()
	}

	o.Close()
	log.Printf("[writeSliceToCsv]\tMatches written to '%s'\n", fout)
}
