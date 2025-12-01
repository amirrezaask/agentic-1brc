package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sort"
	"sync"
)

type Stats struct {
	Min   int64
	Max   int64
	Sum   int64
	Count int64
}

func main() {
	filePath := "../data/medium.txt"
	if len(os.Args) > 1 {
		filePath = os.Args[1]
	}

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	fileInfo, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fileSize := fileInfo.Size()
	f.Close()

	numWorkers := runtime.NumCPU()
	if numWorkers < 1 {
		numWorkers = 1
	}

	// Chunk size calculation
	chunkSize := fileSize / int64(numWorkers)
	var wg sync.WaitGroup

	results := make([]map[string]*Stats, numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		start := int64(i) * chunkSize
		end := start + chunkSize
		if i == numWorkers-1 {
			end = fileSize
		}

		go func(workerID int, startOffset, endOffset int64) {
			defer wg.Done()
			results[workerID] = processChunk(filePath, startOffset, endOffset)
		}(i, start, end)
	}

	wg.Wait()

	// Merge results
	finalStats := make(map[string]*Stats)
	for _, partial := range results {
		for station, stat := range partial {
			if current, exists := finalStats[station]; exists {
				if stat.Min < current.Min {
					current.Min = stat.Min
				}
				if stat.Max > current.Max {
					current.Max = stat.Max
				}
				current.Sum += stat.Sum
				current.Count += stat.Count
			} else {
				finalStats[station] = stat
			}
		}
	}

	// Sort stations
	stations := make([]string, 0, len(finalStats))
	for s := range finalStats {
		stations = append(stations, s)
	}
	sort.Strings(stations)

	// Print output
	fmt.Print("{")
	for i, station := range stations {
		s := finalStats[station]
		mean := float64(s.Sum) / float64(s.Count) / 10.0
		min := float64(s.Min) / 10.0
		max := float64(s.Max) / 10.0

		fmt.Printf("%s=%.1f/%.1f/%.1f", station, min, mean, max)
		if i < len(stations)-1 {
			fmt.Print(",")
		}
	}
	fmt.Println("}")
}

func processChunk(filePath string, start, end int64) map[string]*Stats {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	stats := make(map[string]*Stats)
	currentOffset := start

	skipFirst := false
	if start > 0 {
		// Check previous byte to see if we are at start of line
		f.Seek(start-1, 0)
		b := make([]byte, 1)
		f.Read(b)
		if b[0] != '\n' {
			skipFirst = true
		}
	}

	f.Seek(start, 0)
	r := bufio.NewReader(f)

	if skipFirst {
		skipped, err := r.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return stats
			}
			log.Fatal(err)
		}
		currentOffset += int64(len(skipped))
	}

	for currentOffset < end {
		line, err := r.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				if len(line) > 0 {
					processLine(line, stats)
				}
				break
			}
			log.Fatal(err)
		}

		currentOffset += int64(len(line))
		processLine(line, stats)
	}

	return stats
}

func processLine(line []byte, stats map[string]*Stats) {
	// Remove newline from end if present
	if len(line) > 0 && line[len(line)-1] == '\n' {
		line = line[:len(line)-1]
	}
	if len(line) == 0 {
		return
	}

	sepIndex := bytes.IndexByte(line, ';')
	if sepIndex == -1 {
		return
	}

	// Convert name to string for map key
	name := string(line[:sepIndex])

	tempBytes := line[sepIndex+1:]

	// Parse temp manually to int64 (x10)
	var val int64
	neg := false
	idx := 0
	if len(tempBytes) > 0 && tempBytes[0] == '-' {
		neg = true
		idx++
	}

	for i := idx; i < len(tempBytes); i++ {
		b := tempBytes[i]
		if b == '.' {
			continue
		}
		if b >= '0' && b <= '9' {
			val = val*10 + int64(b-'0')
		}
	}

	if neg {
		val = -val
	}

	s, ok := stats[name]
	if !ok {
		stats[name] = &Stats{
			Min:   val,
			Max:   val,
			Sum:   val,
			Count: 1,
		}
	} else {
		if val < s.Min {
			s.Min = val
		}
		if val > s.Max {
			s.Max = val
		}
		s.Sum += val
		s.Count++
	}
}
