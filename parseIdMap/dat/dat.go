package dat

import (
	"bufio"
	"compress/gzip"
	"log"
	"os"
	"parseIdMap/utility"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Parses the UniProtKB 'IDmapping.dat.gz' file in parallel using go-routines. Given the file path,
// identifier type and accession ID's to filter on, this function will return a string-slice of
// matches to the users input. Currently, this function expects the file to be gzipped.
func ParseDatParallel(fpath string, idtype []string, accessions []string) []string {
	start := time.Now()

	log.Printf("[ParseDatParallel]\tStarted: %s\n", start.Format("01-02-2006 15:04:05"))

	var (
		matches = []string{} // currently this is nonsense while I figure the code out!
		wg      = sync.WaitGroup{}
		mutex   = &sync.Mutex{}
	)

	// Open gzipped dat file
	file, err := os.Open(fpath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Unzip contents into reader
	gzfile, err := gzip.NewReader(file)
	if err != nil {
		log.Fatal(err)
	}
	defer gzfile.Close()

	// Scan lines
	scanner := bufio.NewScanner(gzfile)
	cbuffer := make([]byte, 0, bufio.MaxScanTokenSize)
	scanner.Buffer(cbuffer, bufio.MaxScanTokenSize*50)
	scanner.Scan()

	// Chunk scanned lines ~64k at a time
	linesChunkLen := 1 * 1024
	linesChunkPoolAllocated := int64(0)
	linesPool := sync.Pool{New: func() interface{} {
		lines := make([]string, 0, linesChunkLen)
		atomic.AddInt64(&linesChunkPoolAllocated, 1)
		return lines
	}}
	lines := linesPool.Get().([]string)[:0]

	keepPoolAllocated := int64(0)
	keepPool := sync.Pool{New: func() interface{} {
		kp := make([]string, 0, linesChunkLen)
		atomic.AddInt64(&keepPoolAllocated, 1)
		return kp
	}}

	// Continue looping until can't
	log.Printf("[ParseDatParallel]\tFiltering '%s' for '%s' entries\n", fpath, idtype)
	for {
		lines = append(lines, scanner.Text()) // Append current line to 'lines' variable
		willScan := scanner.Scan()            // Boolean
		if len(lines) == linesChunkLen || !willScan {
			linesToProcess := lines
			wg.Add(len(linesToProcess))

			// Go routine - process current batch of lines
			go func() {
				keep := keepPool.Get().([]string)[:0]

				// Iterate over each line in chunked set and check if pattern matches
				for _, text := range linesToProcess {
					str := strings.Split(text, "\t")

					// Continue if unwanted ID-type
					if !utility.Contains(strings.TrimSpace(str[1]), idtype) {
						continue
					}

					// Continue if accession not in list
					if !utility.Contains(strings.TrimSpace(str[0]), accessions) {
						continue
					}

					// Correct type and accession in 'accessions'
					keep = append(keep, text)
				}

				mutex.Lock()
				for _, i := range keep {
					matches = append(matches, strings.Replace(i, "\t", ",", 2))
				}

				linesPool.Put(linesToProcess)
				keepPool.Put(keep)

				// decrement waitgroup counter by length of input
				wg.Add(-len(linesToProcess))
				mutex.Unlock()
			}()

			lines = linesPool.Get().([]string)[:0] // Make empty again
		}

		// Exit loop when can't scan no more
		if !willScan || len(matches)/2 == len(accessions) {
			break
		}
	}

	wg.Wait() // Wait for above to finish - counter needs to be 0. If negative, it'll error

	log.Printf("\t\t\t\t- Duration: %v\n", time.Since(start))

	return matches
}
