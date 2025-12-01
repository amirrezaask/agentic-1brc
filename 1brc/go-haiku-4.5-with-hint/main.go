package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Stats struct {
	min   float64
	max   float64
	sum   float64
	count int64
}

func main() {
	filePath := flag.String("file", "../data/medium.txt", "Path to the measurements file")
	flag.Parse()

	if flag.NArg() > 0 {
		*filePath = flag.Arg(0)
	}

	measurements, err := os.Open(*filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer measurements.Close()

	stats := make(map[string]*Stats)

	scanner := bufio.NewScanner(measurements)
	const maxBufferSize = 256 * 1024
	buf := make([]byte, maxBufferSize)
	scanner.Buffer(buf, maxBufferSize)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ";", 2)
		if len(parts) != 2 {
			continue
		}

		station := parts[0]
		tempStr := parts[1]

		temp, err := strconv.ParseFloat(tempStr, 64)
		if err != nil {
			continue
		}

		if s, exists := stats[station]; exists {
			s.sum += temp
			s.count++
			if temp < s.min {
				s.min = temp
			}
			if temp > s.max {
				s.max = temp
			}
		} else {
			stats[station] = &Stats{
				min:   temp,
				max:   temp,
				sum:   temp,
				count: 1,
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	stations := make([]string, 0, len(stats))
	for station := range stats {
		stations = append(stations, station)
	}
	sort.Strings(stations)

	var output bytes.Buffer
	output.WriteString("{")

	for i, station := range stations {
		if i > 0 {
			output.WriteString(",")
		}

		s := stats[station]
		mean := s.sum / float64(s.count)

		output.WriteString(fmt.Sprintf("%s=%.1f/%.1f/%.1f",
			station, s.min, mean, s.max))
	}

	output.WriteString("}")
	fmt.Println(output.String())
}

