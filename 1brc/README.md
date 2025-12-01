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

| Implementation | Input | Time (avg) |
|---|---|---|
| go-opus4.5-with-hint | medium.txt | 144.4ms |
| go-gemini3-with-hint | medium.txt | 178.8ms |
| go-opus4.5 | medium.txt | 182.8ms |
| go-gemini3 | medium.txt | 272.4ms |
| go-haiku-4.5 | medium.txt | 391.4ms |
| rust-opus4.5-with-hint | medium.txt | 415.8ms |
| rust-opus-4.5 | medium.txt | 458.2ms |
| go-gpt5.1-with-hint | medium.txt | 892.8ms |
| go-haiku-4.5-with-hint | medium.txt | 1.016s |
| go-qwen-with-hint | medium.txt | 1.101s |
| go-opencode-grok-code-fast | medium.txt | 1.112s |
| go-gpt5.1 | medium.txt | 1.144s |
| go-qwen | medium.txt | 1.263s |
| go-opencode-grok-code-fast-with-hint | medium.txt | 1.406s |

---
*Based on the original [1BRC by Gunnar Morling](https://github.com/gunnarmorling/1brc).*

