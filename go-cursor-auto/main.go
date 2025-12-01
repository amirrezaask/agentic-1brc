package main

import (
	"bytes"
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

type entry struct {
	key   string
	used  bool
	stats stats
}

type table struct {
	entries []entry
	mask    int
}

func nextPow2(n int) int {
	if n <= 1 {
		return 1
	}
	p := 1
	for p < n {
		p <<= 1
	}
	return p
}

func newTable(expected int) *table {
	size := nextPow2(expected * 2)
	return &table{
		entries: make([]entry, size),
		mask:    size - 1,
	}
}

func hashBytes(b []byte) uint64 {
	var h uint64 = 1469598103934665603
	for _, c := range b {
		h ^= uint64(c)
		h *= 1099511628211
	}
	return h
}

func (t *table) add(nameBytes []byte, temp int32) {
	h := hashBytes(nameBytes)
	idx := int(h) & t.mask

	for {
		e := &t.entries[idx]
		if !e.used {
			e.used = true
			e.key = string(nameBytes)
			e.stats = stats{
				min:   temp,
				max:   temp,
				sum:   int64(temp),
				count: 1,
			}
			return
		}
		if e.key == string(nameBytes) {
			s := e.stats
			if temp < s.min {
				s.min = temp
			}
			if temp > s.max {
				s.max = temp
			}
			s.sum += int64(temp)
			s.count++
			e.stats = s
			return
		}
		idx = (idx + 1) & t.mask
	}
}

func (t *table) toSortedKeysAndStats() ([]string, map[string]stats) {
	out := make(map[string]stats, len(t.entries))
	for i := range t.entries {
		e := &t.entries[i]
		if !e.used {
			continue
		}
		out[e.key] = e.stats
	}
	names := make([]string, 0, len(out))
	for k := range out {
		names = append(names, k)
	}
	sort.Strings(names)
	return names, out
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

func processChunk(data []byte, start, end int, tbl *table) {
	i := start
	for i < end {
		lineStart := i

		relSemi := bytes.IndexByte(data[i:end], ';')
		if relSemi < 0 {
			break
		}
		semi := i + relSemi

		valStart := semi + 1
		if valStart >= end {
			break
		}

		relNL := bytes.IndexByte(data[valStart:end], '\n')
		var valEnd int
		if relNL < 0 {
			valEnd = end
			i = end
		} else {
			valEnd = valStart + relNL
			i = valEnd + 1
		}

		nameBytes := data[lineStart:semi]
		valBytes := data[valStart:valEnd]

		temp := parseTemperature(valBytes)

		tbl.add(nameBytes, temp)
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

	tbl := newTable(1 << 18)
	processChunk(data, 0, int(size), tbl)

	names, global := tbl.toSortedKeysAndStats()

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
