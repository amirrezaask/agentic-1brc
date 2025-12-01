package main

import (
	"fmt"
	"os"
	"sort"
	"syscall"

	"golang.org/x/sys/unix"
)

type stats struct {
	min   int32
	max   int32
	sum   int64
	count int64
}

func parseTemperature(b []byte) int32 {
	n := len(b)
	if n == 0 {
		return 0
	}

	sign := int32(1)
	i := 0
	if b[0] == '-' {
		sign = -1
		i = 1
	}

	var v int32
	for ; i < n; i++ {
		c := b[i]
		if c == '.' {
			continue
		}
		v = v*10 + int32(c-'0')
	}
	return sign * v
}

func processChunk(data []byte, start, end int, out map[string]stats) {
	local := out

	i := start
	for i < end {
		lineStart := i

		// find ';'
		semi := -1
		for i < end {
			if data[i] == ';' {
				semi = i
				i++
				break
			}
			i++
		}
		if semi == -1 {
			// no complete line
			break
		}

		// find end of line '\n'
		valStart := i
		for i < end && data[i] != '\n' {
			i++
		}
		valEnd := i
		if i < end && data[i] == '\n' {
			i++
		}

		nameBytes := data[lineStart:semi]
		valBytes := data[valStart:valEnd]

		temp := parseTemperature(valBytes)

		name := string(nameBytes)
		s, ok := local[name]
		if !ok {
			local[name] = stats{
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
			local[name] = s
		}
	}
}

func main() {
	path := "../data/medium.txt"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to stat file: %v\n", err)
		os.Exit(1)
	}

	size := info.Size()
	if size == 0 {
		fmt.Println("{}")
		return
	}

	data, err := syscall.Mmap(int(f.Fd()), 0, int(size), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to mmap file: %v\n", err)
		os.Exit(1)
	}

	_ = unix.Madvise(data, unix.MADV_SEQUENTIAL)

	defer syscall.Munmap(data)

	global := make(map[string]stats, 1<<18)
	processChunk(data, 0, int(size), global)

	names := make([]string, 0, len(global))
	for name := range global {
		names = append(names, name)
	}
	sort.Strings(names)

	fmt.Print("{")
	for i, name := range names {
		st := global[name]
		minVal := float64(st.min) / 10.0
		maxVal := float64(st.max) / 10.0
		avgVal := float64(st.sum) / float64(st.count) / 10.0

		if i > 0 {
			fmt.Print(",")
		}
		fmt.Printf("%s=%.1f/%.1f/%.1f", name, minVal, avgVal, maxVal)
	}
	fmt.Println("}")
}
