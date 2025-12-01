package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// StationStats holds the min, mean, and max temperature values for a station
type StationStats struct {
	min float64
	max float64
	sum float64
	cnt int64
}

func main() {
	var filepath string
	if len(os.Args) > 1 {
		filepath = os.Args[1]
	} else {
		filepath = "../data/medium.txt"
	}

	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	stations := make(map[string]*StationStats)

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

		if stats, exists := stations[station]; exists {
			if temp < stats.min {
				stats.min = temp
			}
			if temp > stats.max {
				stats.max = temp
			}
			stats.sum += temp
			stats.cnt++
		} else {
			stations[station] = &StationStats{
				min: temp,
				max: temp,
				sum: temp,
				cnt: 1,
			}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	// Get sorted station names
	sortedStations := make([]string, 0, len(stations))
	for station := range stations {
		sortedStations = append(sortedStations, station)
	}
	sort.Strings(sortedStations)

	// Format output
	fmt.Print("{")
	for i, station := range sortedStations {
		stats := stations[station]
		mean := stats.sum / float64(stats.cnt)
		if i > 0 {
			fmt.Print(",")
		}
		fmt.Printf("%s=%.1f/%.1f/%.1f", station, stats.min, mean, stats.max)
	}
	fmt.Println("}")
}