#!/usr/bin/env python3

import argparse
import os
import subprocess
import sys
import time
from dataclasses import dataclass
from typing import List, Optional


@dataclass
class Implementation:
    name: str
    cwd: str
    build: Optional[List[str]]
    cmd: List[str]


def get_root_dir() -> str:
    return os.path.dirname(os.path.abspath(__file__))


def run_command(args: List[str], cwd: str) -> int:
    try:
        result = subprocess.run(
            args,
            cwd=cwd,
            check=False,
            stdout=subprocess.DEVNULL,
            stderr=subprocess.DEVNULL,
        )
        return result.returncode
    except OSError as exc:
        print(f"failed to run {' '.join(args)}: {exc}", file=sys.stderr)
        return 1


def run_impl(impl: Implementation, file_path: str) -> tuple[float, int, str]:
    root = get_root_dir()
    cwd = os.path.join(root, impl.cwd)

    if impl.build:
        code = run_command(impl.build, cwd=cwd)
        if code != 0:
            return 0.0, code, ""

    cmd = [arg.format(file=file_path) for arg in impl.cmd]

    start = time.perf_counter()
    try:
        proc = subprocess.run(
            cmd,
            cwd=cwd,
            check=False,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
        )
    except OSError as exc:
        print(f"{impl.name}: failed to execute: {exc}", file=sys.stderr)
        return 0.0, 1, ""

    duration = time.perf_counter() - start
    output = proc.stdout.replace("\n", "")
    return duration, proc.returncode, output


def format_duration(seconds: float) -> str:
    if seconds >= 1.0:
        return f"{seconds:.3f}s"
    millis = seconds * 1000.0
    return f"{millis:.1f}ms"


def main(argv: Optional[List[str]] = None) -> int:
    root = get_root_dir()

    parser = argparse.ArgumentParser(
        description="Run all 1BRC implementations and compare results."
    )
    parser.add_argument(
        "file",
        nargs="?",
        default=os.path.join(root, "data", "medium.txt"),
        help="input measurements file (default: data/medium.txt)",
    )
    parser.add_argument(
        "--runs",
        type=int,
        default=1,
        help="number of timed runs per implementation (default: 5; first is treated as warmup)",
    )
    parser.add_argument(
        "--format",
        choices=["table", "json", "csv"],
        default="table",
        help="output format (default: table)",
    )
    args = parser.parse_args(argv)

    file_path = os.path.abspath(args.file)
    if not os.path.isfile(file_path):
        print(f"Error: file not found: {file_path}", file=sys.stderr)
        return 1

    implementations: List[Implementation] = [
        Implementation(
            name="go-haiku-4.5",
            cwd="go-haiku-4.5",
            build=["go", "build", "-o", "1brc-go", "main.go"],
            cmd=["./1brc-go", "{file}"],
        ),
        Implementation(
            name="go-gpt5.1-with-hint",
            cwd="go-gpt5.1-with-hint",
            build=["go", "build", "-o", "1brc-go", "main.go"],
            cmd=["./1brc-go", "{file}"],
        ),
        Implementation(
            name="go-gpt5.1",
            cwd="go-gpt5.1",
            build=["go", "build", "-o", "1brc-go", "main.go"],
            cmd=["./1brc-go", "{file}"],
        ),
        Implementation(
            name="go-qwen",
            cwd="go-qwen",
            build=["go", "build", "-o", "1brc-go", "main.go"],
            cmd=["./1brc-go", "{file}"],
        ),
        Implementation(
            name="go-qwen-with-hint",
            cwd="go-qwen-with-hint",
            build=["go", "build", "-o", "1brc-go", "main.go"],
            cmd=["./1brc-go", "{file}"],
        ),
        Implementation(
            name="go-opus4.5",
            cwd="go-opus4.5",
            build=["go", "build", "-o", "1brc-go", "main.go"],
            cmd=["./1brc-go", "{file}"],
        ),
        Implementation(
            name="go-opus4.5-with-hint",
            cwd="go-opus4.5-with-hint",
            build=["go", "build", "-o", "1brc-go", "main.go"],
            cmd=["./1brc-go", "{file}"],
        ),
        Implementation(
            name="rust-opus-4.5",
            cwd="rust-opus-4.5",
            build=["cargo", "build", "--release"],
            cmd=["./target/release/onebrc", "{file}"],
        ),
        Implementation(
            name="go-gemini3",
            cwd="go-gemini3",
            build=["go", "build", "-o", "1brc-go", "main.go"],
            cmd=["./1brc-go", "{file}"],
        ),
        Implementation(
            name="go-gemini3-with-hint",
            cwd="go-gemini3-with-hint",
            build=["go", "build", "-o", "1brc-go", "main.go"],
            cmd=["./1brc-go", "{file}"],
        ),
        Implementation(
            name="go-haiku-4.5-with-hint",
            cwd="go-haiku-4.5-with-hint",
            build=["go", "build", "-o", "1brc-go", "main.go"],
            cmd=["./1brc-go", "{file}"],
        ),
        Implementation(
            name="go-opencode-grok-code-fast",
            cwd="go-opencode-grok-code-fast",
            build=["go", "build", "-o", "1brc-go", "main.go"],
            cmd=["./1brc-go", "{file}"],
        ),
        Implementation(
            name="go-opencode-grok-code-fast-with-hint",
            cwd="go-opencode-grok-code-fast-with-hint",
            build=["go", "build", "-o", "1brc-go", "main.go"],
            cmd=["./1brc-go", "{file}"],
        ),
        Implementation(
            name="rust-opus4.5-with-hint",
            cwd="rust-opus4.5-with-hint",
            build=["cargo", "build", "--release"],
            cmd=["./target/release/onebrc", "{file}"],
        ),
        
    ]

    names: List[str] = []
    times: List[float] = []

    runs = max(1, args.runs)

    impl_num = 0
    for impl in implementations:
        impl_num += 1
        names.append(impl.name)

        per_run: List[float] = []
        last_code: int = 0

        for run_num in range(1, runs + 1):
            if args.format == "table":
                print(f"Running {impl_num}/{len(implementations)} run {run_num}/{runs}", end='\r')
                
            duration, code, _output = run_impl(impl, file_path)
            per_run.append(duration)
            last_code = code
            if code != 0:
                break

        if last_code != 0:
            print(f"{impl.name}: exited with code {last_code}", file=sys.stderr)
            times.append(0.0)
        else:
            if len(per_run) > 1:
                warmup_excluded = per_run[1:]
            else:
                warmup_excluded = per_run
            avg = sum(warmup_excluded) / len(warmup_excluded)
            times.append(avg)

    if args.format == "table":
        print()  # Clear progress line
    input_label = os.path.basename(file_path)
    sorted_items = sorted(zip(names, times), key=lambda x: (x[1] == 0.0, x[1]))

    results = []
    for name, t in sorted_items:
        results.append({"name": name, "input": input_label, "time": t})

    if args.format == "table":
        print()
        header = f"{'Implementation':25} {'Input':15} {'Time (avg)':15}"
        sep = f"{'-' * 25:25} {'-' * 15:15} {'-' * 15:15}"
        print(header)
        print(sep)

        for item in results:
            pretty = format_duration(item["time"]) if item["time"] > 0 else "-"
            print(f"{item['name']:25} {item['input']:15} {pretty:15}")

        print()
    elif args.format == "json":
        import json
        print(json.dumps(results))
    elif args.format == "csv":
        import csv
        writer = csv.DictWriter(sys.stdout, fieldnames=["name", "input", "time"])
        writer.writeheader()
        writer.writerows(results)

    return 0


if __name__ == "__main__":
    raise SystemExit(main())


