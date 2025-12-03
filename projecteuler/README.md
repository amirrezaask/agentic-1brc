# 🧮 Project Euler Benchmark

This benchmark evaluates the ability of AI models to solve Project Euler problems. The goal is to implement solutions for a set of mathematical/programming problems, testing logical reasoning, algorithm design, and correctness.

## 📂 Structure

- `problems/`: Contains the problem descriptions (e.g., `problem_0001.txt`).
- `solutions.txt`: Contains the expected answers for validation.
- `score.py`: Script to run implementations and verify their answers.
- `PROMPT.md`: The prompt provided to agents to generate the solutions.
- `get_problems.py`: Utility to fetch problems from Project Euler.

## 🚀 How to Run

Use the `score.py` script to run all implementations and generate a report.

```bash
# Run all implementations and show a table
python3 score.py

# Run and output CSV
python3 score.py --format csv

# Run and output JSON
python3 score.py --format json
```

## 📝 Implementations

Each implementation is in its own directory following the naming convention `<language>-<model>`.


## Latest Scores
```
Implementation                 Correct    Total      Accuracy     Time (s)  
------------------------------------------------------------------------
go-gemini-3                    12         13         92.31%       2.20
go-haiku4.5                    9          20         45.00%       0.01
go-opus4.5                     13         20         65.00%       5.52
go-sonnet4.5                   11         20         55.00%       0.66

```
