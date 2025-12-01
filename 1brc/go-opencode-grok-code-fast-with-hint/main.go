package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Station struct {
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

	stations := make(map[string]*Station)
	var mu sync.Mutex

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ";")
		if len(parts) != 2 {
			continue
		}
		station := parts[0]
		temp, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			continue
		}
		mu.Lock()
		if s, ok := stations[station]; ok {
			if temp < s.min {
				s.min = temp
			}
			if temp > s.max {
				s.max = temp
			}
			s.sum += temp
			s.count++
		} else {
			stations[station] = &Station{min: temp, max: temp, sum: temp, count: 1}
		}
		mu.Unlock()
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	// Get sorted keys
	var keys []string
	for k := range stations {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build output
	var results []string
	for _, k := range keys {
		s := stations[k]
		mean := s.sum / float64(s.count)
		results = append(results, fmt.Sprintf("%s=%.1f/%.1f/%.1f", k, s.min, mean, s.max))
	}
	fmt.Printf("{%s}\n", strings.Join(results, ","))
}
