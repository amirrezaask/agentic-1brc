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

func main() {
	solutions := map[int]func() interface{}{
		1:   func() interface{} { return solution_0001() },
		2:   func() interface{} { return solution_0002() },
		3:   func() interface{} { return solution_0003() },
		4:   func() interface{} { return solution_0004() },
		5:   func() interface{} { return solution_0005() },
		6:   func() interface{} { return solution_0006() },
		7:   func() interface{} { return solution_0007() },
		8:   func() interface{} { return solution_0008() },
		9:   func() interface{} { return solution_0009() },
		10:  func() interface{} { return solution_0010() },
		350: func() interface{} { return solution_0350() },
		387: func() interface{} { return solution_0387() },
		407: func() interface{} { return solution_0407() },
		416: func() interface{} { return solution_0416() },
		428: func() interface{} { return solution_0428() },
		434: func() interface{} { return solution_0434() },
		447: func() interface{} { return solution_0447() },
		458: func() interface{} { return solution_0458() },
		510: func() interface{} { return solution_0510() },
		709: func() interface{} { return solution_0709() },
	}

	var problemNums []int
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
			Result:   fmt.Sprintf("%v", answer),
			TimeMs:   elapsedMs,
		}
		jsonBytes, _ := json.Marshal(result)
		fmt.Println(string(jsonBytes))
	}
}
