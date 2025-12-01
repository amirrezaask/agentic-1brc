package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"sort"
)

type stats struct {
	min   int32
	max   int32
	sum   int64
	count int64
}

func parseTemp(b []byte) int32 {
	if len(b) == 0 {
		return 0
	}

	var sign int32 = 1
	i := 0
	if b[0] == '-' {
		sign = -1
		i = 1
	}

	var v int32
	for ; i < len(b); i++ {
		c := b[i]
		if c == '.' {
			continue
		}
		v = v*10 + int32(c-'0')
	}
	return sign * v
}

func processFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	r := bufio.NewReaderSize(f, 1<<20)

	m := make(map[string]*stats, 1024)

	for {
		line, err := r.ReadBytes('\n')
		if len(line) == 0 && err != nil {
			break
		}

		line = bytes.TrimRight(line, "\r\n")
		if len(line) == 0 {
			if err != nil {
				break
			}
			continue
		}

		sep := bytes.IndexByte(line, ';')
		if sep <= 0 || sep+2 >= len(line) {
			if err != nil {
				break
			}
			continue
		}

		station := string(line[:sep])
		tempBytes := line[sep+1:]
		temp := parseTemp(tempBytes)

		s, ok := m[station]
		if !ok {
			m[station] = &stats{
				min:   temp,
				max:   temp,
				sum:   int64(temp),
				count: 1,
			}
		} else {
			if temp < s.min {
				s.min = temp
			}
			if temp > s.max {
				s.max = temp
			}
			s.sum += int64(temp)
			s.count++
		}

		if err != nil {
			break
		}
	}

	if len(m) == 0 {
		return "{}", nil
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	buf.Grow(len(m) * 32)
	buf.WriteByte('{')

	for i, k := range keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(k)
		buf.WriteByte('=')
		s := m[k]

		min := float64(s.min) / 10.0
		mean := float64(s.sum) / float64(s.count) / 10.0
		max := float64(s.max) / 10.0

		fmt.Fprintf(&buf, "%.1f/%.1f/%.1f", min, mean, max)
	}

	buf.WriteByte('}')
	return buf.String(), nil
}

func main() {
	path := "../data/medium.txt"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	out, err := processFile(path)
	if err != nil {
		log.Fatalf("failed to process file: %v", err)
	}

	fmt.Println(out)
}


