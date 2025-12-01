package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Stats struct {
	min, max, sum float64
	count         int
}

func main() {
	filePath := "../data/medium.txt"
	if len(os.Args) > 1 {
		filePath = os.Args[1]
	}

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	stations := make(map[string]*Stats)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ";")
		if len(parts) != 2 {
			continue
		}
		temp, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			continue
		}
		station := parts[0]
		if stats, ok := stations[station]; ok {
			if temp < stats.min {
				stats.min = temp
			}
			if temp > stats.max {
				stats.max = temp
			}
			stats.sum += temp
			stats.count++
		} else {
			stations[station] = &Stats{min: temp, max: temp, sum: temp, count: 1}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	// Collect keys
	keys := make([]string, 0, len(stations))
	for k := range stations {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build output
	results := make([]string, len(keys))
	for i, key := range keys {
		stats := stations[key]
		mean := stats.sum / float64(stats.count)
		mean_rounded := math.Round(mean*10) / 10
		results[i] = fmt.Sprintf("%s=%.1f/%.1f/%.1f", key, stats.min, mean_rounded, stats.max)
	}

	fmt.Printf("{%s}\n", strings.Join(results, ","))
}
