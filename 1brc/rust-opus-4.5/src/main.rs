use memmap2::Mmap;
use rustc_hash::FxHashMap;
use std::env;
use std::fs::File;
use std::io::Write;
use std::thread;

#[derive(Clone)]
struct Stats {
    min: i64,
    max: i64,
    sum: i64,
    count: u64,
}

impl Stats {
    fn new(value: i64) -> Self {
        Stats {
            min: value,
            max: value,
            sum: value,
            count: 1,
        }
    }

    fn update(&mut self, value: i64) {
        if value < self.min {
            self.min = value;
        }
        if value > self.max {
            self.max = value;
        }
        self.sum += value;
        self.count += 1;
    }

    fn merge(&mut self, other: &Stats) {
        if other.min < self.min {
            self.min = other.min;
        }
        if other.max > self.max {
            self.max = other.max;
        }
        self.sum += other.sum;
        self.count += other.count;
    }
}

fn parse_temperature(bytes: &[u8]) -> i64 {
    let mut idx = 0;
    let negative = bytes[0] == b'-';
    if negative {
        idx = 1;
    }

    let mut value: i64 = 0;
    while idx < bytes.len() {
        let b = bytes[idx];
        if b == b'.' {
            idx += 1;
            continue;
        }
        value = value * 10 + (b - b'0') as i64;
        idx += 1;
    }

    if negative {
        -value
    } else {
        value
    }
}

fn format_temp(value: i64) -> String {
    let negative = value < 0;
    let abs_value = value.abs();
    let integer = abs_value / 10;
    let decimal = abs_value % 10;
    if negative {
        format!("-{}.{}", integer, decimal)
    } else {
        format!("{}.{}", integer, decimal)
    }
}

fn process_chunk(data: &[u8]) -> FxHashMap<Vec<u8>, Stats> {
    let mut map: FxHashMap<Vec<u8>, Stats> = FxHashMap::default();

    let mut i = 0;
    while i < data.len() {
        let line_start = i;

        while i < data.len() && data[i] != b';' {
            i += 1;
        }
        let station = &data[line_start..i];
        i += 1;

        let temp_start = i;
        while i < data.len() && data[i] != b'\n' {
            i += 1;
        }
        let temp_bytes = &data[temp_start..i];
        i += 1;

        let temp = parse_temperature(temp_bytes);

        match map.get_mut(station) {
            Some(stats) => stats.update(temp),
            None => {
                map.insert(station.to_vec(), Stats::new(temp));
            }
        }
    }

    map
}

fn find_line_start(data: &[u8], pos: usize) -> usize {
    if pos == 0 {
        return 0;
    }
    let mut p = pos;
    while p < data.len() && data[p] != b'\n' {
        p += 1;
    }
    if p < data.len() {
        p + 1
    } else {
        data.len()
    }
}

fn main() {
    let args: Vec<String> = env::args().collect();
    let file_path = if args.len() > 1 {
        &args[1]
    } else {
        "../data/medium.txt"
    };

    let file = File::open(file_path).expect("Failed to open file");
    let mmap = unsafe { Mmap::map(&file).expect("Failed to mmap file") };
    let data = &mmap[..];

    let num_threads = thread::available_parallelism()
        .map(|n| n.get())
        .unwrap_or(4);

    let chunk_size = data.len() / num_threads;

    let handles: Vec<_> = (0..num_threads)
        .map(|i| {
            let start = if i == 0 {
                0
            } else {
                find_line_start(data, i * chunk_size)
            };
            let end = if i == num_threads - 1 {
                data.len()
            } else {
                find_line_start(data, (i + 1) * chunk_size)
            };

            let chunk = data[start..end].to_vec();

            thread::spawn(move || process_chunk(&chunk))
        })
        .collect();

    let mut final_map: FxHashMap<Vec<u8>, Stats> = FxHashMap::default();
    for handle in handles {
        let partial = handle.join().expect("Thread panicked");
        for (station, stats) in partial {
            match final_map.get_mut(&station) {
                Some(existing) => existing.merge(&stats),
                None => {
                    final_map.insert(station, stats);
                }
            }
        }
    }

    let mut stations: Vec<_> = final_map.into_iter().collect();
    stations.sort_by(|a, b| a.0.cmp(&b.0));

    let mut output = String::with_capacity(stations.len() * 50);
    output.push('{');

    for (i, (station, stats)) in stations.iter().enumerate() {
        if i > 0 {
            output.push(',');
        }
        let station_name = String::from_utf8_lossy(station);
        let mean = (stats.sum as f64 / stats.count as f64).round() as i64;
        output.push_str(&format!(
            "{}={}/{}/{}",
            station_name,
            format_temp(stats.min),
            format_temp(mean),
            format_temp(stats.max)
        ));
    }
    output.push('}');

    let stdout = std::io::stdout();
    let mut handle = stdout.lock();
    writeln!(handle, "{}", output).expect("Failed to write output");
}

