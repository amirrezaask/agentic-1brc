package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type stats struct {
	min  int16
	max  int16
	sum  int64
	cnt  int64
	init bool
}

func parseTemp(b []byte) int16 {
	n := len(b)
	if n == 0 {
		return 0
	}

	i := 0
	sign := int16(1)
	if b[i] == '-' {
		sign = -1
		i++
	}

	var v int16
	for ; i < n && b[i] != '.'; i++ {
		v = v*10 + int16(b[i]-'0')
	}

	if i+2 < n && b[i] == '.' {
		v = v*10 + int16(b[i+1]-'0')
	}

	return v * sign
}

func formatTemp(value int64) string {
	sign := ""
	if value < 0 {
		sign = "-"
		value = -value
	}
	intPart := value / 10
	dec := value % 10
	return fmt.Sprintf("%s%d.%d", sign, intPart, dec)
}

func main() {
	path := "../data/medium.txt"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open file", err)
		os.Exit(1)
	}
	defer f.Close()

	reader := bufio.NewReaderSize(f, 1<<20)
	statsMap := make(map[string]*stats, 10000)

	for {
		line, err := reader.ReadBytes('\n')
		if len(line) == 0 && err != nil {
			break
		}

		if len(line) > 0 && line[len(line)-1] == '\n' {
			line = line[:len(line)-1]
		}

		sep := -1
		for i, c := range line {
			if c == ';' {
				sep = i
				break
			}
		}
		if sep <= 0 || sep+1 >= len(line) {
			if err != nil {
				break
			}
			continue
		}

		name := string(line[:sep])
		tempVal := parseTemp(line[sep+1:])

		s, ok := statsMap[name]
		if !ok {
			s = &stats{
				min:  tempVal,
				max:  tempVal,
				sum:  int64(tempVal),
				cnt:  1,
				init: true,
			}
			statsMap[name] = s
		} else {
			if !s.init {
				s.min = tempVal
				s.max = tempVal
				s.sum = int64(tempVal)
				s.cnt = 1
				s.init = true
			} else {
				if tempVal < s.min {
					s.min = tempVal
				}
				if tempVal > s.max {
					s.max = tempVal
				}
				s.sum += int64(tempVal)
				s.cnt++
			}
		}

		if err != nil {
			break
		}
	}

	names := make([]string, 0, len(statsMap))
	for name := range statsMap {
		names = append(names, name)
	}

	for i := 1; i < len(names); i++ {
		j := i
		for j > 0 && names[j-1] > names[j] {
			names[j-1], names[j] = names[j], names[j-1]
			j--
		}
	}

	out := make([]byte, 0, len(statsMap)*32)
	out = append(out, '{')
	for i, name := range names {
		if i > 0 {
			out = append(out, ',')
		}
		s := statsMap[name]
		out = append(out, name...)
		out = append(out, '=')

		out = append(out, formatTemp(int64(s.min))...)
		out = append(out, '/')

		meanTenths := int64(0)
		if s.cnt > 0 {
			avg := float64(s.sum) / float64(s.cnt)
			meanTenths = int64(math.Round(avg))
		}
		out = append(out, formatTemp(meanTenths)...)
		out = append(out, '/')

		out = append(out, formatTemp(int64(s.max))...)
	}
	out = append(out, '}')

	os.Stdout.Write(out)
}


