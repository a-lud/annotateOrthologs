package tab

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

// Parse the 'idmapping_selected.tab.gz' file for key columns that match a set of UniProtKB accessions.
func ParseTabParallel(fpath string, idtype []string, mapIdx map[string]int, accessions []string) []string {

	start := time.Now()

	log.Printf("[ParseTabParallel]\tStarted: %s\n", start.Format("01-02-2006 15:04:05"))

	var (
		matches []string = []string{}
		wg               = sync.WaitGroup{}
		mutex            = &sync.Mutex{}
	)

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
		kp := make([][]string, 0, linesChunkLen)
		atomic.AddInt64(&keepPoolAllocated, 1)
		return kp
	}}

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

	// Increase default size of maxScanTokenSize or else long lines will crash the scanner
	cbuffer := make([]byte, 0, bufio.MaxScanTokenSize)
	scanner.Buffer(cbuffer, bufio.MaxScanTokenSize*50)

	// Skip the first entry to prevent error for len(0)
	scanner.Scan()

	log.Printf("[ParseTabParallel]\tFiltering '%s' for '%s' entries\n", fpath, idtype)
	for {
		lines = append(lines, scanner.Text()) // Append current line to 'lines' variable
		willScan := scanner.Scan()            // Boolean

		if len(lines) == linesChunkLen || !willScan {

			linesToProcess := lines
			wg.Add(len(linesToProcess))

			// Go routine - process current batch of lines
			go func() {

				keep := keepPool.Get().([][]string)[:0]

				for _, text := range linesToProcess {

					// Split row into slice
					str := strings.Split(text, "\t")

					// Continue if UniProtKB accession not in 'accessions' slice
					if !utility.Contains(strings.TrimSpace(str[0]), accessions) {
						continue
					}

					// User requested fields
					tmp := []string{}
					tmp = append(tmp, str[0]) // accession

					for _, i := range idtype {
						tmp = append(tmp, str[mapIdx[i]]) // user columns
					}

					// return object
					keep = append(keep, tmp)
					// matches = append(matches, strings.Join(tmp, ","))
				}

				mutex.Lock()
				for _, i := range keep {
					matches = append(matches, strings.Join(i, ","))
				}

				linesPool.Put(linesToProcess)
				keepPool.Put(keep)

				// decrement waitgroup counter by length of input
				wg.Add(-len(linesToProcess))

				mutex.Unlock()
			}()

			lines = linesPool.Get().([]string)[:0] // Make empty again
		}

		// Exit loop when can't scan any more / we find data for all the accessions of interest
		if !willScan || len(matches) == len(accessions) {
			break
		}
	}

	wg.Wait() // Wait for above to finish - counter needs to be 0. If negative, it'll error

	log.Printf("\t\t\t\t- Duration: %v\n", time.Since(start))

	return matches
}
