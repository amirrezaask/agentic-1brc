package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// StationData holds the min, max, and count for each station
type StationData struct {
	Min   float64
	Max   float64
	Sum   float64
	Count int64
}

func main() {
	// Use first command-line argument as file path, otherwise use default
	inputPath := "../data/medium.txt"
	if len(os.Args) > 1 {
		inputPath = os.Args[1]
	}

	// Open the input file
	file, err := os.Open(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Map to store station data
	stations := make(map[string]*StationData)

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		
		// Split the line by semicolon
		parts := strings.Split(line, ";")
		if len(parts) != 2 {
			continue // Skip malformed lines
		}
		
		station := parts[0]
		tempStr := parts[1]
		
		// Parse temperature
		temp, err := strconv.ParseFloat(tempStr, 64)
		if err != nil {
			continue // Skip invalid temperatures
		}
		
		// Update station data
		if data, exists := stations[station]; exists {
			if temp < data.Min {
				data.Min = temp
			}
			if temp > data.Max {
				data.Max = temp
			}
			data.Sum += temp
			data.Count++
		} else {
			stations[station] = &StationData{
				Min:   temp,
				Max:   temp,
				Sum:   temp,
				Count: 1,
			}
		}
	}
	
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Create sorted list of station names
	stationNames := make([]string, 0, len(stations))
	for name := range stations {
		stationNames = append(stationNames, name)
	}
	sort.Strings(stationNames)

	// Format and print the results
	fmt.Print("{")
	for i, name := range stationNames {
		data := stations[name]
		mean := data.Sum / float64(data.Count)
		
		if i > 0 {
			fmt.Print(",")
		}
		fmt.Printf("%s=%.1f/%.1f/%.1f", name, data.Min, mean, data.Max)
	}
	fmt.Println("}")
}