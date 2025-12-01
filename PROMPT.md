# One Billion Row Challenge (1BRC)

## Challenge Overview

Your task is to implement a solution for the **One Billion Row Challenge** - a performance-focused programming challenge that tests how quickly you can process and aggregate data from a massive text file containing one billion rows of temperature measurements.

**Reference:** https://github.com/gunnarmorling/1brc

## Input Format

The input is a text file (`measurements.txt`) containing temperature measurements from various weather stations. The file format is:

```
<station_name>;<temperature>
```

**Example:**
```
Hamburg;12.0
Bulawayo;8.9
Palembang;38.8
St. John's;15.2
Cracow;12.6
...
```

### Data Specifications

- **File size:** ~13 GB (one billion rows)
- **Station names:** UTF-8 encoded strings (1-100 bytes), no `;` or `\n` characters
- **Temperature values:** Floating point numbers with exactly one decimal place, ranging from -99.9 to 99.9
- **Number of unique stations:** Up to 10,000 different weather stations
- **Line format:** `<station>;<temperature>\n`

## Task

For each unique weather station, calculate and output:
1. **Minimum temperature**
2. **Mean (average) temperature**
3. **Maximum temperature**

## Output Format

Print the results to standard output in the following format:

```
{<station1>=<min>/<mean>/<max>,<station2>=<min>/<mean>/<max>,...}
```

**Requirements:**
- Station names must be sorted **alphabetically**
- All temperature values must be formatted with **exactly one decimal place**
- Results are enclosed in curly braces `{}`
- Stations are separated by commas `,` (no spaces)
- Each station's stats are separated by `/`

**Example Output:**
```
{Abha=5.0/18.0/27.4,Abidjan=15.7/26.0/34.1,Accra=14.7/26.4/33.1,...}
```

## Implementation Requirements

1. **Create a new directory** named `$Language-$Model` (e.g., `rust-opus4.5`, `python-gpt4`, `cpp-gemini3`)
2. **Accept file path as command-line argument** - the program should accept the path to the measurements file as the first argument
3. **Default path:** If no argument provided, default to `../data/medium.txt`
4. **Focus on performance** - the goal is to process one billion rows as fast as possible
5. **Correctness first** - ensure your solution produces correct output before optimizing

## Testing

Test data files are available in the `data/` directory:
- `small.txt` - Small test file for quick validation
- `medium.txt` - Medium-sized file for development testing
- `measurements.txt` - Full one billion row file (generate using `create_measurements.py`)

## Evaluation Criteria

1. **Correctness:** Output must match the expected format exactly
2. **Performance:** Execution time on the full 1 billion row file
3. **Code quality:** Clean, readable implementation

## Getting Started

1. Create your solution directory: `mkdir $Language-$Model`
2. Implement your solution
3. Test with small/medium files first
4. Optimize for the full dataset
5. Run against `../data/measurements.txt`

# RULES
* In case you are generating for golang use go1.24 in go.mod
* Add the new implementation in run_all.py script as an implementation to be runned in next benchmark.
* Don't look into other agents implementation at all.

**Good luck! Make it fast! 🚀**

