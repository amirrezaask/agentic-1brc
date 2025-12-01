package main

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"sort"
	"sync"
	"syscall"
)

type Stats struct {
	min   int64
	max   int64
	sum   int64
	count int64
}

func parseTemp(b []byte) int64 {
	negative := false
	idx := 0
	if b[0] == '-' {
		negative = true
		idx = 1
	}

	var val int64
	for ; idx < len(b); idx++ {
		c := b[idx]
		if c == '.' {
			continue
		}
		val = val*10 + int64(c-'0')
	}

	if negative {
		return -val
	}
	return val
}

func processChunk(data []byte, results map[string]*Stats) {
	for len(data) > 0 {
		newlineIdx := bytes.IndexByte(data, '\n')
		if newlineIdx == -1 {
			if len(data) > 0 {
				semiIdx := bytes.IndexByte(data, ';')
				if semiIdx != -1 && semiIdx < len(data)-1 {
					station := string(data[:semiIdx])
					temp := parseTemp(data[semiIdx+1:])
					if s, ok := results[station]; ok {
						if temp < s.min {
							s.min = temp
						}
						if temp > s.max {
							s.max = temp
						}
						s.sum += temp
						s.count++
					} else {
						results[station] = &Stats{min: temp, max: temp, sum: temp, count: 1}
					}
				}
			}
			break
		}

		line := data[:newlineIdx]
		data = data[newlineIdx+1:]

		if len(line) == 0 {
			continue
		}

		semiIdx := bytes.IndexByte(line, ';')
		if semiIdx == -1 {
			continue
		}

		station := string(line[:semiIdx])
		temp := parseTemp(line[semiIdx+1:])

		if s, ok := results[station]; ok {
			if temp < s.min {
				s.min = temp
			}
			if temp > s.max {
				s.max = temp
			}
			s.sum += temp
			s.count++
		} else {
			results[station] = &Stats{min: temp, max: temp, sum: temp, count: 1}
		}
	}
}

func mergeResults(dst, src map[string]*Stats) {
	for station, s := range src {
		if d, ok := dst[station]; ok {
			if s.min < d.min {
				d.min = s.min
			}
			if s.max > d.max {
				d.max = s.max
			}
			d.sum += s.sum
			d.count += s.count
		} else {
			dst[station] = s
		}
	}
}

func formatTemp(val int64) string {
	negative := val < 0
	if negative {
		val = -val
	}
	intPart := val / 10
	decPart := val % 10
	if negative {
		return fmt.Sprintf("-%d.%d", intPart, decPart)
	}
	return fmt.Sprintf("%d.%d", intPart, decPart)
}

func main() {
	filePath := "../data/medium.txt"
	if len(os.Args) > 1 {
		filePath = os.Args[1]
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting file info: %v\n", err)
		os.Exit(1)
	}
	fileSize := int(fileInfo.Size())

	if fileSize == 0 {
		fmt.Println("{}")
		return
	}

	data, err := syscall.Mmap(int(file.Fd()), 0, fileSize, syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error mmapping file: %v\n", err)
		os.Exit(1)
	}
	defer syscall.Munmap(data)

	numWorkers := runtime.NumCPU()
	chunkSize := fileSize / numWorkers
	if chunkSize < 1024*1024 {
		chunkSize = 1024 * 1024
	}

	type chunkRange struct {
		start int
		end   int
	}

	chunks := make([]chunkRange, 0, numWorkers)
	start := 0
	for i := 0; i < numWorkers && start < fileSize; i++ {
		end := start + chunkSize
		if end >= fileSize {
			end = fileSize
		} else {
			for end < fileSize && data[end] != '\n' {
				end++
			}
			if end < fileSize {
				end++
			}
		}
		chunks = append(chunks, chunkRange{start: start, end: end})
		start = end
	}

	var wg sync.WaitGroup
	resultsChan := make(chan map[string]*Stats, len(chunks))

	for _, chunk := range chunks {
		wg.Add(1)
		go func(c chunkRange) {
			defer wg.Done()
			localResults := make(map[string]*Stats, 10000)
			processChunk(data[c.start:c.end], localResults)
			resultsChan <- localResults
		}(chunk)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	finalResults := make(map[string]*Stats, 10000)
	for r := range resultsChan {
		mergeResults(finalResults, r)
	}

	stations := make([]string, 0, len(finalResults))
	for station := range finalResults {
		stations = append(stations, station)
	}
	sort.Strings(stations)

	var output bytes.Buffer
	output.Grow(len(stations) * 50)
	output.WriteByte('{')
	for i, station := range stations {
		s := finalResults[station]
		mean := s.sum / s.count
		if i > 0 {
			output.WriteByte(',')
		}
		output.WriteString(station)
		output.WriteByte('=')
		output.WriteString(formatTemp(s.min))
		output.WriteByte('/')
		output.WriteString(formatTemp(mean))
		output.WriteByte('/')
		output.WriteString(formatTemp(s.max))
	}
	output.WriteByte('}')
	fmt.Println(output.String())
}
