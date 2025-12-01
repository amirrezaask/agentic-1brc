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
    args = parser.parse_args(argv)

    file_path = os.path.abspath(args.file)
    if not os.path.isfile(file_path):
        print(f"Error: file not found: {file_path}", file=sys.stderr)
        return 1

    implementations: List[Implementation] = [
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
        # Qwen
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
        # Opus4.5
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
        # Gemini 3
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
        
    ]

    names: List[str] = []
    times: List[float] = []

    runs = max(1, args.runs)

    for impl in implementations:
        names.append(impl.name)

        per_run: List[float] = []
        last_code: int = 0

        for _ in range(runs):
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

    print()
    header = f"{'Implementation':20} {'Input':10} {'Time (avg)':15}"
    sep = f"{'-' * 20:20} {'-' * 10:10} {'-' * 15:15}"
    print(header)
    print(sep)

    input_label = os.path.basename(file_path)
    for name, t in zip(names, times):
        pretty = format_duration(t) if t > 0 else "-"
        print(f"{name:20} {input_label:10} {pretty:15}")

    print()
    return 0


if __name__ == "__main__":
    raise SystemExit(main())


