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

const (
	defaultPath = "../data/medium.txt"
	// FNV-1a 64-bit constants
	offset64 = 14695981039346656037
	prime64  = 1099511628211
	// Hash table size - power of 2.
	// Max 10k stations, 64k slots gives load factor ~0.15
	tableSize = 1 << 16
)

// Stats holds the running statistics for a station
// We store temp as int64 (val * 10) to avoid float ops
type Stats struct {
	Min   int64
	Max   int64
	Sum   int64
	Count int64
	Name  []byte
}

// Entry for the open-addressing hash table
type Entry struct {
	Key   []byte
	Stats *Stats
}

// HashTable is a simple open addressing hash table
type HashTable struct {
	entries []Entry
	storage []Stats // Pre-allocated storage for stats
	count   int
}

func NewHashTable() *HashTable {
	return &HashTable{
		entries: make([]Entry, tableSize),
		storage: make([]Stats, 10000), // Reserve for max expected stations
		count:   0,
	}
}

// hash returns the FNV-1a hash of the key
// Inline-able simple hash
func fnv1a(key []byte) uint64 {
	var h uint64 = offset64
	for _, b := range key {
		h ^= uint64(b)
		h *= prime64
	}
	return h
}

// Get returns the stats for the given key, creating it if necessary
func (ht *HashTable) Get(key []byte, hashVal uint64) *Stats {
	idx := hashVal & (tableSize - 1)
	for {
		entryKey := ht.entries[idx].Key
		if entryKey == nil {
			// Found empty slot, insert new
			if ht.count >= len(ht.storage) {
				// Grow storage if needed (unlikely given constraints)
				ht.storage = append(ht.storage, Stats{})
			}

			// Allocate new key copy only on first sight
			keyCopy := make([]byte, len(key))
			copy(keyCopy, key)

			stats := &ht.storage[ht.count]
			stats.Name = keyCopy
			stats.Min = 9223372036854775807  // MaxInt64
			stats.Max = -9223372036854775808 // MinInt64
			ht.count++

			ht.entries[idx].Key = keyCopy
			ht.entries[idx].Stats = stats
			return stats
		}
		if bytes.Equal(entryKey, key) {
			return ht.entries[idx].Stats
		}
		idx = (idx + 1) & (tableSize - 1)
	}
}

// parseTemp parses a temperature string into an int64 (multiplied by 10)
// e.g., "12.3" -> 123, "-1.2" -> -12
func parseTemp(b []byte) int64 {
	var val int64
	neg := false
	i := 0
	if b[i] == '-' {
		neg = true
		i++
	}

	// Parse integer part
	for ; b[i] != '.'; i++ {
		val = val*10 + int64(b[i]-'0')
	}
	i++ // Skip dot

	// Parse decimal part (always 1 digit)
	if i < len(b) {
		val = val*10 + int64(b[i]-'0')
	}

	if neg {
		return -val
	}
	return val
}

func processChunk(data []byte, results chan<- *HashTable, wg *sync.WaitGroup) {
	defer wg.Done()

	ht := NewHashTable()

	var start int
	lenData := len(data)

	for start < lenData {
		// Find semicolon
		// Simple linear scan is fast enough for short strings
		end := start
		for data[end] != ';' {
			end++
		}

		name := data[start:end]
		h := fnv1a(name)

		end++ // Skip semicolon

		// Find newline
		tempStart := end
		for end < lenData && data[end] != '\n' {
			end++
		}

		temp := parseTemp(data[tempStart:end])

		stats := ht.Get(name, h)
		if temp < stats.Min {
			stats.Min = temp
		}
		if temp > stats.Max {
			stats.Max = temp
		}
		stats.Sum += temp
		stats.Count++

		start = end + 1
	}
	results <- ht
}

func main() {
	path := defaultPath
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}
	size := fi.Size()

	// mmap the file
	data, err := syscall.Mmap(int(f.Fd()), 0, int(size), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		panic(err)
	}
	defer syscall.Munmap(data)

	numWorkers := runtime.NumCPU()
	// Ensure we have at least one worker
	if numWorkers < 1 {
		numWorkers = 1
	}

	chunkSize := size / int64(numWorkers)

	var wg sync.WaitGroup
	results := make(chan *HashTable, numWorkers)

	start := int64(0)
	for i := 0; i < numWorkers; i++ {
		end := start + chunkSize
		if i == numWorkers-1 {
			end = size
		} else {
			// Align to newline
			// Be careful not to go out of bounds
			for end < size && data[end] != '\n' {
				end++
			}
			if end < size {
				end++ // Include the newline
			} else {
				end = size
			}
		}

		if start >= size {
			break
		}

		wg.Add(1)
		// Pass slice to worker
		go processChunk(data[start:end], results, &wg)
		start = end
	}

	wg.Wait()
	close(results)

	// Merge results
	// Use a map for merging since collisions are handled by key string
	finalStats := make(map[string]*Stats)

	for ht := range results {
		for i := 0; i < ht.count; i++ {
			s := &ht.storage[i]
			name := string(s.Name)
			if existing, ok := finalStats[name]; ok {
				if s.Min < existing.Min {
					existing.Min = s.Min
				}
				if s.Max > existing.Max {
					existing.Max = s.Max
				}
				existing.Sum += s.Sum
				existing.Count += s.Count
			} else {
				// Clone stats for final map
				finalStats[name] = &Stats{
					Min:   s.Min,
					Max:   s.Max,
					Sum:   s.Sum,
					Count: s.Count,
					Name:  s.Name,
				}
			}
		}
	}

	// Sort and print
	names := make([]string, 0, len(finalStats))
	for name := range finalStats {
		names = append(names, name)
	}
	sort.Strings(names)

	fmt.Print("{")
	for i, name := range names {
		s := finalStats[name]
		// Convert back to float representation
		min := float64(s.Min) / 10.0
		mean := float64(s.Sum) / float64(s.Count) / 10.0
		max := float64(s.Max) / 10.0

		if i > 0 {
			fmt.Print(",")
		}
		fmt.Printf("%s=%.1f/%.1f/%.1f", name, min, mean, max)
	}
	fmt.Println("}")
}
