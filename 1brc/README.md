# 🏎️ 1 Billion Row Challenge

> **Who writes the fastest code?** AI models tackle the [One Billion Row Challenge](https://github.com/gunnarmorling/1brc).

## 📋 The Challenge

Process a massive text file containing **1,000,000,000 rows** of temperature measurements and calculate the min, mean, and max temperature for each weather station.

**Input:** `Hamburg;12.0` (x 1 Billion)  
**Output:** `{Abha=5.0/18.0/27.4, Abidjan=15.7/26.0/34.1, ...}`

## 🤖 Methodology

Each model is tested with two variations:
1. **Raw Prompt** — Standard challenge description only ([PROMPT.md](./PROMPT.md))
2. **With Hints** — Nudged to use advanced techniques like mmap, SIMD, concurrency ([PROMPT_WITH_HINTS.md](./PROMPT_WITH_HINTS.md))

## 📊 Results

**Latest Benchmark Run:**
```
Rank  Implementation            Input           Time (avg)
----- ------------------------- --------------- ---------------
    1 go-gemini3-with-hint      medium.txt      335.8ms
    2 go-opus4.5-with-hint      medium.txt      339.7ms
    3 go-opus4.5                medium.txt      340.2ms
    4 rust-opus4.5-with-hint    medium.txt      353.6ms
    5 rust-opus-4.5             medium.txt      358.2ms
    6 go-gpt5.2-codex           medium.txt      360.5ms
    7 go-gemini3                medium.txt      477.4ms
    8 rust-gpt5.2-codex         medium.txt      595.9ms
    9 go-gpt5.1                 medium.txt      809.0ms
   10 go-gpt5.1-with-hint       medium.txt      840.4ms
   11 rust-gpt5.2-codex-with-hint medium.txt      918.7ms
   12 go-haiku-4.5              medium.txt      933.9ms
   13 go-gpt5.2-codex-with-hint medium.txt      952.1ms
   14 go-haiku-4.5-with-hint    medium.txt      1.030s
   15 go-qwen                   medium.txt      1.084s
   16 go-opencode-grok-code-fast-with-hint medium.txt      1.107s
   17 go-opencode-grok-code-fast medium.txt      1.178s
   18 go-qwen-with-hint         medium.txt      1.390s
```

*Based on the original [1BRC by Gunnar Morling](https://github.com/gunnarmorling/1brc).*
