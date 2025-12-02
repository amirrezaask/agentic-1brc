use std::collections::HashMap;
use std::env;
use std::fs::File;
use std::hash::{BuildHasherDefault, Hasher};
use std::io::{self, Read};

type FastMap = HashMap<Vec<u8>, Stats, BuildHasherDefault<FnvHasher>>;

#[derive(Clone, Copy, Debug)]
struct Stats {
    min: i32,
    max: i32,
    sum: i64,
    count: u32,
}

impl Stats {
    fn new(value: i32) -> Self {
        Self {
            min: value,
            max: value,
            sum: value as i64,
            count: 1,
        }
    }

    fn update(&mut self, value: i32) {
        if value < self.min {
            self.min = value;
        }
        if value > self.max {
            self.max = value;
        }
        self.sum += value as i64;
        self.count += 1;
    }
}

#[derive(Clone, Copy)]
struct FnvHasher(u64);

impl Default for FnvHasher {
    fn default() -> Self {
        Self(0xcbf29ce484222325)
    }
}

impl Hasher for FnvHasher {
    fn write(&mut self, bytes: &[u8]) {
        const PRIME: u64 = 0x00000100000001b3;
        for &byte in bytes {
            self.0 ^= byte as u64;
            self.0 = self.0.wrapping_mul(PRIME);
        }
    }

    fn finish(&self) -> u64 {
        self.0
    }
}

fn parse_temperature(bytes: &[u8]) -> i32 {
    let mut idx = 0;
    let mut negative = false;
    if bytes[idx] == b'-' {
        negative = true;
        idx += 1;
    }

    let mut value: i32 = 0;
    while idx < bytes.len() {
        let b = bytes[idx];
        if b == b'.' {
            idx += 1;
            let tenth = (bytes[idx] - b'0') as i32;
            value = value * 10 + tenth;
            break;
        } else {
            value = value * 10 + (b - b'0') as i32;
            idx += 1;
        }
    }

    if negative {
        -value
    } else {
        value
    }
}

fn process_line(line: &[u8], map: &mut FastMap) {
    let mut split = 0;
    while split < line.len() && line[split] != b';' {
        split += 1;
    }

    let station = &line[..split];
    let temp_value = parse_temperature(&line[split + 1..]);

    if let Some(stats) = map.get_mut(station) {
        stats.update(temp_value);
    } else {
        map.insert(station.to_vec(), Stats::new(temp_value));
    }
}

fn process_file(path: &str) -> io::Result<FastMap> {
    let mut file = File::open(path)?;
    let mut buffer = vec![0u8; 8 * 1024 * 1024];
    let mut carry = Vec::with_capacity(256);
    let mut map: FastMap =
        HashMap::with_capacity_and_hasher(16_384, BuildHasherDefault::<FnvHasher>::default());

    loop {
        let read = file.read(&mut buffer)?;
        if read == 0 {
            break;
        }

        let mut start = 0;
        for idx in 0..read {
            if buffer[idx] == b'\n' {
                if carry.is_empty() {
                    process_line(&buffer[start..idx], &mut map);
                } else {
                    carry.extend_from_slice(&buffer[start..idx]);
                    process_line(&carry, &mut map);
                    carry.clear();
                }
                start = idx + 1;
            }
        }

        if start < read {
            carry.extend_from_slice(&buffer[start..read]);
        }
    }

    if !carry.is_empty() {
        process_line(&carry, &mut map);
    }

    Ok(map)
}

fn main() {
    let path = env::args()
        .nth(1)
        .unwrap_or_else(|| "../data/medium.txt".to_string());

    let map = match process_file(&path) {
        Ok(m) => m,
        Err(err) => {
            eprintln!("Failed to process file {}: {}", path, err);
            std::process::exit(1);
        }
    };

    let mut entries: Vec<(String, Stats)> = map
        .into_iter()
        .map(|(name, stats)| {
            let station = String::from_utf8(name).expect("station names must be valid UTF-8");
            (station, stats)
        })
        .collect();

    entries.sort_by(|a, b| a.0.cmp(&b.0));

    print!("{{");
    for (idx, (name, stats)) in entries.iter().enumerate() {
        if idx > 0 {
            print!(",");
        }
        let min = stats.min as f64 / 10.0;
        let mean = (stats.sum as f64 / stats.count as f64) / 10.0;
        let max = stats.max as f64 / 10.0;
        print!("{}={:.1}/{:.1}/{:.1}", name, min, mean, max);
    }
    println!("}}");
}
