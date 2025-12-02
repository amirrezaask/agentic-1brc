package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
)

type stats struct {
	min int32
	max int32
	sum int64
	cnt int64
}

func parseTemp(field []byte) int32 {
	var neg bool
	var i int
	if len(field) > 0 {
		if field[0] == '-' {
			neg = true
			i++
		} else if field[0] == '+' {
			i++
		}
	}
	var val int32
	for i < len(field) && field[i] != '.' {
		val = val*10 + int32(field[i]-'0')
		i++
	}
	if i < len(field) && field[i] == '.' {
		i++
		if i < len(field) {
			val = val*10 + int32(field[i]-'0')
		}
	}
	if neg {
		val = -val
	}
	return val
}

func process(path string) (map[string]stats, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReaderSize(file, 1<<20)
	result := make(map[string]stats, 16_384)

	for {
		line, err := reader.ReadBytes('\n')
		if len(line) == 0 {
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
			continue
		}
		line = bytes.TrimRight(line, "\n")
		sep := bytes.IndexByte(line, ';')
		if sep < 0 {
			if err == io.EOF {
				break
			}
			continue
		}
		station := string(line[:sep])
		temp := parseTemp(line[sep+1:])

		if s, ok := result[station]; ok {
			if temp < s.min {
				s.min = temp
			}
			if temp > s.max {
				s.max = temp
			}
			s.sum += int64(temp)
			s.cnt++
			result[station] = s
		} else {
			result[station] = stats{min: temp, max: temp, sum: int64(temp), cnt: 1}
		}

		if err == io.EOF {
			break
		}
	}

	return result, nil
}

func formatTenths(v int32) string {
	return fmt.Sprintf("%.1f", float64(v)/10.0)
}

func formatMean(sum int64, cnt int64) string {
	return fmt.Sprintf("%.1f", (float64(sum)/float64(cnt))/10.0)
}

func main() {
	path := "../data/medium.txt"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	statsMap, err := process(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	names := make([]string, 0, len(statsMap))
	for k := range statsMap {
		names = append(names, k)
	}
	sort.Strings(names)

	out := make([]byte, 0, len(names)*32)
	out = append(out, '{')
	for i, name := range names {
		if i > 0 {
			out = append(out, ',')
		}
		s := statsMap[name]
		out = append(out, name...)
		out = append(out, '=')
		out = append(out, formatTenths(s.min)...)
		out = append(out, '/')
		out = append(out, formatMean(s.sum, s.cnt)...)
		out = append(out, '/')
		out = append(out, formatTenths(s.max)...)
	}
	out = append(out, '}')
	os.Stdout.Write(out)
}
