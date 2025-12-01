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

// Constants for the challenge
const (
	defaultPath = "../data/medium.txt"
	// Table size needs to be power of 2 and larger than max stations (10,000).
	// 1<<15 = 32768 should be enough to reduce collisions.
	tableSize = 1 << 15
	tableMask = tableSize - 1
)

// Stats holds the aggregate data for a station
type Stats struct {
	Min   int32
	Max   int32
	Sum   int64
	Count int32
}

// Entry is a bucket in our custom hash table
type Entry struct {
	Key   []byte // We'll store a reference to the mmapped data to avoid allocation
	Stats Stats
}

// HashTable is a simple linear-probing hash table
type HashTable struct {
	Entries []Entry
}

func NewHashTable() *HashTable {
	return &HashTable{
		Entries: make([]Entry, tableSize),
	}
}

// hash computes the FNV-1a hash of the key
func hash(key []byte) uint64 {
	var h uint64 = 14695981039346656037
	for _, b := range key {
		h ^= uint64(b)
		h *= 1099511628211
	}
	return h
}

// Put updates the stats for a given station
func (ht *HashTable) Put(key []byte, temp int32) {
	h := hash(key)
	idx := h & tableMask

	for {
		entry := &ht.Entries[idx]
		if entry.Key == nil {
			// Found empty slot, insert new
			// Note: we need to copy the key if we were not mmapped and the buffer was reused,
			// but here 'key' is a slice of the mmapped file which is immutable for our purpose.
			// actually, 'key' is a slice of the chunk. It's safe to keep it if the chunk lives on.
			// But wait, we want to use the key in the final merge.
			// The mmap region is valid until munmap.
			entry.Key = key
			entry.Stats.Min = temp
			entry.Stats.Max = temp
			entry.Stats.Sum = int64(temp)
			entry.Stats.Count = 1
			return
		}
		if bytes.Equal(entry.Key, key) {
			// Found existing, update
			if temp < entry.Stats.Min {
				entry.Stats.Min = temp
			}
			if temp > entry.Stats.Max {
				entry.Stats.Max = temp
			}
			entry.Stats.Sum += int64(temp)
			entry.Stats.Count++
			return
		}
		// Collision, linear probe
		idx = (idx + 1) & tableMask
	}
}

// parseInt parses a temperature string (e.g. "12.3", "-8.9") into an int32 * 10
// It assumes the format is always X.Y or -X.Y or XX.Y etc.
// Returns value * 10
func parseInt(b []byte) int32 {
	var val int32
	var neg bool
	i := 0
	if b[0] == '-' {
		neg = true
		i++
	}
	// Parse integer part
	for ; i < len(b); i++ {
		if b[i] == '.' {
			i++
			break
		}
		val = val*10 + int32(b[i]-'0')
	}
	// Parse decimal part (assumed 1 digit)
	if i < len(b) {
		val = val*10 + int32(b[i]-'0')
	}
	if neg {
		return -val
	}
	return val
}

func worker(data []byte, wg *sync.WaitGroup, resultChan chan *HashTable) {
	defer wg.Done()
	ht := NewHashTable()

	// Iterate over the data
	// Format: station;temp\n
	start := 0
	for start < len(data) {
		// Find semi-colon
		endName := start
		for endName < len(data) && data[endName] != ';' {
			endName++
		}
		if endName >= len(data) {
			break
		}

		// Find newline
		endLine := endName + 1
		for endLine < len(data) && data[endLine] != '\n' {
			endLine++
		}

		name := data[start:endName]
		tempBytes := data[endName+1 : endLine]
		temp := parseInt(tempBytes)

		ht.Put(name, temp)

		start = endLine + 1
	}
	resultChan <- ht
}

func main() {
	// 1. Parse args
	path := defaultPath
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	// 2. Open file
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

	if size == 0 {
		fmt.Print("{}")
		return
	}

	// 3. Mmap
	data, err := syscall.Mmap(int(f.Fd()), 0, int(size), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		panic(err)
	}
	defer syscall.Munmap(data)

	// 4. Launch workers
	numWorkers := runtime.NumCPU()
	chunkSize := int(size) / numWorkers
	var wg sync.WaitGroup
	resultChan := make(chan *HashTable, numWorkers)

	start := 0
	for i := 0; i < numWorkers; i++ {
		end := start + chunkSize
		if i == numWorkers-1 {
			end = int(size)
		} else {
			// Adjust end to next newline
			for end < int(size) && data[end] != '\n' {
				end++
			}
			end++ // Include the newline
		}

		if start >= int(size) {
			break
		}

		wg.Add(1)
		// Run worker on slice
		go worker(data[start:end], &wg, resultChan)
		start = end
	}

	wg.Wait()
	close(resultChan)

	// 5. Merge results
	// We'll use a final map to merge results.
	// We can reuse the same structure or just a simple Go map since collision handling is done.
	// Actually, let's just use a Go map for the final merge to make sorting easier.
	finalStats := make(map[string]*Stats)

	for ht := range resultChan {
		for _, entry := range ht.Entries {
			if entry.Key == nil {
				continue
			}
			// Convert key to string for the final map.
			// In a super optimized version we might avoid this, but for merge it's fine.
			// However, string(entry.Key) allocates.
			// Since we only have max 10k stations, this is negligible compared to 1B rows.
			name := string(entry.Key)
			if s, ok := finalStats[name]; ok {
				if entry.Stats.Min < s.Min {
					s.Min = entry.Stats.Min
				}
				if entry.Stats.Max > s.Max {
					s.Max = entry.Stats.Max
				}
				s.Sum += entry.Stats.Sum
				s.Count += entry.Stats.Count
			} else {
				// Copy stats
				finalStats[name] = &Stats{
					Min:   entry.Stats.Min,
					Max:   entry.Stats.Max,
					Sum:   entry.Stats.Sum,
					Count: entry.Stats.Count,
				}
			}
		}
	}

	// 6. Sort and print
	stations := make([]string, 0, len(finalStats))
	for name := range finalStats {
		stations = append(stations, name)
	}
	sort.Strings(stations)

	fmt.Print("{")
	for i, name := range stations {
		s := finalStats[name]
		mean := float64(s.Sum) / float64(s.Count) / 10.0
		min := float64(s.Min) / 10.0
		max := float64(s.Max) / 10.0
		fmt.Printf("%s=%.1f/%.1f/%.1f", name, min, mean, max)
		if i < len(stations)-1 {
			fmt.Print(",")
		}
	}
	fmt.Println("}")
}

