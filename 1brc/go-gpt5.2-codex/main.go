package main

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"sort"
	"sync"
	"syscall"
)

type agg struct {
	min   int32
	max   int32
	sum   int64
	count int64
}

type segment struct {
	start int
	end   int
}

func main() {
	path := "../data/medium.txt"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	data, err := mmapFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read file: %v\n", err)
		os.Exit(1)
	}
	defer syscall.Munmap(data)

	segments := splitWork(data, runtime.NumCPU())
	results := make([]map[string]agg, len(segments))

	var wg sync.WaitGroup
	wg.Add(len(segments))
	for i, seg := range segments {
		i, seg := i, seg
		go func() {
			defer wg.Done()
			results[i] = processChunk(data, seg.start, seg.end)
		}()
	}
	wg.Wait()

	final := merge(results)
	printResult(final)
}

func mmapFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	size := info.Size()
	if size == 0 {
		return []byte{}, nil
	}
	if size > int64(^uint(0)>>1) {
		return nil, fmt.Errorf("file too large for mmap: %d bytes", size)
	}

	return syscall.Mmap(int(f.Fd()), 0, int(size), syscall.PROT_READ, syscall.MAP_SHARED)
}

func splitWork(data []byte, workers int) []segment {
	if workers < 1 {
		workers = 1
	}
	total := len(data)
	if total == 0 || workers == 1 {
		return []segment{{start: 0, end: total}}
	}

	chunk := total / workers
	if chunk < 1<<20 { // small file: avoid too many tiny chunks
		return []segment{{start: 0, end: total}}
	}

	segs := make([]segment, 0, workers)
	start := 0
	for w := 0; w < workers; w++ {
		end := start + chunk
		if w == workers-1 || end >= total {
			segs = append(segs, segment{start: start, end: total})
			break
		}
		for end < total && data[end] != '\n' {
			end++
		}
		if end < total {
			end++
		}
		segs = append(segs, segment{start: start, end: end})
		start = end
	}

	return segs
}

func processChunk(data []byte, start, end int) map[string]agg {
	result := make(map[string]agg, 1024)

	i := start
	for i < end {
		nameStart := i
		for i < end && data[i] != ';' {
			i++
		}
		if i >= end {
			break
		}
		nameEnd := i
		i++ // skip ';'

		if i >= end {
			break
		}

		sign := int32(1)
		if data[i] == '-' {
			sign = -1
			i++
		}

		var val int32
		for i < end && data[i] != '.' {
			val = val*10 + int32(data[i]-'0')
			i++
		}
		i++ // skip '.'
		if i < end {
			val = val*10 + int32(data[i]-'0')
			i++
		}

		temp := sign * val

		for i < end && data[i] != '\n' {
			i++
		}
		if i < end && data[i] == '\n' {
			i++
		}

		key := string(data[nameStart:nameEnd])
		if st, ok := result[key]; ok {
			if temp < st.min {
				st.min = temp
			}
			if temp > st.max {
				st.max = temp
			}
			st.sum += int64(temp)
			st.count++
			result[key] = st
		} else {
			result[key] = agg{
				min:   temp,
				max:   temp,
				sum:   int64(temp),
				count: 1,
			}
		}
	}

	return result
}

func merge(maps []map[string]agg) map[string]agg {
	final := make(map[string]agg, 16384)
	for _, m := range maps {
		for k, st := range m {
			if existing, ok := final[k]; ok {
				if st.min < existing.min {
					existing.min = st.min
				}
				if st.max > existing.max {
					existing.max = st.max
				}
				existing.sum += st.sum
				existing.count += st.count
				final[k] = existing
			} else {
				final[k] = st
			}
		}
	}
	return final
}

func printResult(m map[string]agg) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	buf.WriteByte('{')
	for idx, k := range keys {
		if idx > 0 {
			buf.WriteByte(',')
		}
		st := m[k]
		mean := float64(st.sum) / float64(st.count) / 10.0
		fmt.Fprintf(&buf, "%s=%.1f/%.1f/%.1f", k, float64(st.min)/10.0, mean, float64(st.max)/10.0)
	}
	buf.WriteByte('}')
	fmt.Print(buf.String())
}
