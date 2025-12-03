# 🤖 Agentic Benchmarks

> Pitting AI models against each other in real-world coding challenges.

This repository hosts a collection of benchmarks designed to evaluate how well different AI models perform on practical programming tasks.

## 📚 Benchmarks

| Benchmark | Category | Description |
|-----------|----------|-------------|
| [1 Billion Row Challenge](./1brc/) | Performance | Process 1B temperature readings as fast as possible |
| [Project Euler](./projecteuler/) | Reasoning/Algorithm | Solve mathematical and programming problems |

## 🏎️ 1BRC at a Glance

| Rank | Implementation | Time |
|------|----------------|------|
| 🥇 | go-gemini3-with-hint | 91.5ms |
| 🥈 | go-opus4.5-with-hint | 146.2ms |
| 🥉 | go-opus4.5 | 174.7ms |
| 4 | go-haiku-4.5 | 195.1ms |
| 5 | go-gemini3 | 220.6ms |

[View full results →](./1brc/)

## 📂 Repository Structure

```
agentic-benchmarks/
├── README.md           # This file
├── 1brc/               # 1 Billion Row Challenge
├── projecteuler/       # Project Euler Challenge
└── ...
```

Each benchmark has its own directory with setup instructions, prompts, implementations, and results.
