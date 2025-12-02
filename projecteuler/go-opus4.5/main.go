package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

type Result struct {
	Solution int     `json:"solution"`
	Result   string  `json:"result"`
	TimeMs   float64 `json:"time_ms"`
}

var solutions = map[int]func() int64{
	1:   solution_0001,
	2:   solution_0002,
	3:   solution_0003,
	4:   solution_0004,
	5:   solution_0005,
	6:   solution_0006,
	7:   solution_0007,
	8:   solution_0008,
	9:   solution_0009,
	10:  solution_0010,
	350: solution_0350,
	387: solution_0387,
	407: solution_0407,
	416: solution_0416,
	428: solution_0428,
	434: solution_0434,
	447: solution_0447,
	458: solution_0458,
	510: solution_0510,
	709: solution_0709,
}

func main() {
	problemNums := make([]int, 0, len(solutions))
	for num := range solutions {
		problemNums = append(problemNums, num)
	}
	sort.Ints(problemNums)

	for _, problemNum := range problemNums {
		solutionFn := solutions[problemNum]
		start := time.Now()
		answer := solutionFn()
		elapsedMs := float64(time.Since(start).Microseconds()) / 1000.0
		result := Result{
			Solution: problemNum,
			Result:   fmt.Sprintf("%d", answer),
			TimeMs:   elapsedMs,
		}
		jsonBytes, _ := json.Marshal(result)
		fmt.Println(string(jsonBytes))
	}
}
