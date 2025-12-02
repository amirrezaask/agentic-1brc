# Project Euler Challenge

## Challenge Overview

Your task is to implement solutions for following problems in problem set:
Problem 1
Problem 2
Problem 3
Problem 4
Problem 5
Problem 6
Problem 7
Problem 8
Problem 9
Problem 10
Problem 350
Problem 387
Problem 407
Problem 416
Problem 428
Problem 434
Problem 447
Problem 458
Problem 510
Problem 709

**Reference:** https://projecteuler.net/

## Problem Files

Problem descriptions are available in the `problems/` directory. Each file follows this naming convention:
- `problem_XXXX.txt` where `XXXX` is the zero-padded problem number (e.g., `problem_0001.txt`, `problem_0042.txt`)

## Task

Implement solutions for first 100 Project Euler problems as you can. Each solution should:
1. Read and understand the problem from the corresponding `problems/problem_XXXX.txt` file
2. Implement a correct solution
3. Return the answer as the function's return value

## Implementation Requirements

### Directory Structure

Create a new directory named `${Laguage}-${Model}` (e.g., `python-opus4.5`, `go-gpt4`, `rust-gemini3`)

```
projecteuler/
├── ${Laguage}-${Model}/
│   ├── solutions.{ext}      # All solution functions
│   └── main.{ext}           # Entry point that runs all solutions
```

### Function Naming Convention

Each solution must be implemented as a function with this naming convention:

```
solution_XXXX
```

Where `XXXX` is the zero-padded problem number.

**Examples:**
- `solution_0001` - Solution for Problem 1
- `solution_0042` - Solution for Problem 42
- `solution_0100` - Solution for Problem 100

### Function Signature

Each solution function should:
- Take no arguments (problem parameters are hardcoded per the problem description)
- Return the answer (as a string or integer depending on the language)

**Example (Python):**
```python
def solution_0001():
    # Find the sum of all multiples of 3 or 5 below 1000
    return sum(x for x in range(1000) if x % 3 == 0 or x % 5 == 0)
```

**Example (Go):**
```go
func solution_0001() int {
    sum := 0
    for i := 0; i < 1000; i++ {
        if i%3 == 0 || i%5 == 0 {
            sum += i
        }
    }
    return sum
}
```

**Example (Rust):**
```rust
fn solution_0001() -> i64 {
    (0..1000).filter(|x| x % 3 == 0 || x % 5 == 0).sum()
}
```

### Main Entry Point

Create a main entry point that:
1. Runs all implemented solutions
2. Prints results in **machine-readable JSON format**

**Output Format:**
```json
{"results":[{"problem":1,"answer":"233168","time_ms":0.5},{"problem":2,"answer":"4613732","time_ms":0.2},...]}
```

Each result object contains:
- `problem`: The problem number (integer)
- `answer`: The computed answer (string)
- `time_ms`: Execution time in milliseconds (float)

**Example Main (Python):**
```python
import json
import time

solutions = {
    1: solution_0001,
    2: solution_0002,
    # ... add more as implemented
}

def main():
    results = []
    for problem_num, solution_fn in sorted(solutions.items()):
        start = time.perf_counter()
        answer = solution_fn()
        elapsed_ms = (time.perf_counter() - start) * 1000
        results.append({
            "problem": problem_num,
            "answer": str(answer),
            "time_ms": round(elapsed_ms, 3)
        })
    print(json.dumps({"results": results}))

if __name__ == "__main__":
    main()
```

**Example Main (Go):**
```go
package main

import (
    "encoding/json"
    "fmt"
    "time"
)

type Result struct {
    Problem int     `json:"problem"`
    Answer  string  `json:"answer"`
    TimeMs  float64 `json:"time_ms"`
}

type Output struct {
    Results []Result `json:"results"`
}

var solutions = map[int]func() int{
    1: solution_0001,
    2: solution_0002,
    // ... add more as implemented
}

func main() {
    var results []Result
    for problemNum, solutionFn := range solutions {
        start := time.Now()
        answer := solutionFn()
        elapsedMs := float64(time.Since(start).Microseconds()) / 1000.0
        results = append(results, Result{
            Problem: problemNum,
            Answer:  fmt.Sprintf("%d", answer),
            TimeMs:  elapsedMs,
        })
    }
    output := Output{Results: results}
    jsonBytes, _ := json.Marshal(output)
    fmt.Println(string(jsonBytes))
}
```

## Evaluation Criteria

1. **Correctness:** Solutions must produce the correct answer for each problem
2. **Coverage:** Number of problems solved
3. **Performance:** Execution time for each solution (aim for under 1 minute per problem)
4. **Code quality:** Clean, readable implementation

## Getting Started

1. Create your solution directory: `mkdir ${Laguage}-${Model}`
2. Read problem descriptions from `problems/` directory
3. Implement solutions starting with Problem 1
4. Test each solution before moving to the next
5. Create the main entry point to run all solutions

## Rules
- Always compile the binary with -o runner
- Do NOT look at other agents' implementations
- Do NOT use external problem-solving resources or lookup answers
- Solutions must compute the answer, not hardcode known answers
- Each solution should complete in under 60 seconds
- Use only standard library features (no external dependencies unless absolutely necessary)
- Implement solutions one by one and write each one to the file before going to the next one.
- make sure to follow out score.py script output expectation.
- make sure to create a new directory following convention ${Language}-{MODEL}
- DONT LOOK INTO OTHER IMPLEMENTATIONS AND INTERNET YOU NEED TO SOLVE THEM YOUSELF.
- YOU DON"T HAVE PERMISSION TO READ OTHER AGENTS IMPLEMENTATIONS

## Language-Specific Notes

### Python
- Use Python 3.10+
- File: `solutions.py` and `main.py` (or combine into single file)

### Go
- Use Go 1.24 in go.mod
- File: `solutions.go` and `main.go`

### Rust
- Use latest stable Rust
- File: `src/main.rs` (solutions can be in same file or separate modules)

### JavaScript/TypeScript
- Use Node.js 20+
- File: `solutions.js` and `main.js`

**Good luck! Solve as many as you can! 🧮**

