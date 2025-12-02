use std::collections::hash_map::Entry;
use std::collections::HashMap;
use std::env;
use std::error::Error;
use std::fs::File;
use std::io::{BufRead, BufReader};
use std::path::Path;

#[derive(Debug, Clone)]
struct Stats {
    min: i32,   // tenths of degrees
    max: i32,   // tenths of degrees
    sum: i64,   // tenths of degrees
    count: u32, // number of measurements
}

fn parse_temp(raw: &[u8]) -> i32 {
    let mut idx = 0;
    let mut negative = false;
    if let Some(&first) = raw.first() {
        if first == b'-' {
            negative = true;
            idx += 1;
        } else if first == b'+' {
            idx += 1;
        }
    }

    let mut int_part: i32 = 0;
    while idx < raw.len() && raw[idx] != b'.' {
        int_part = int_part * 10 + (raw[idx] - b'0') as i32;
        idx += 1;
    }

    // Skip the dot.
    idx += 1;
    let frac = if idx < raw.len() {
        (raw[idx] - b'0') as i32
    } else {
        0
    };

    let mut value = int_part * 10 + frac;
    if negative {
        value = -value;
    }
    value
}

fn process_file(path: &Path) -> Result<HashMap<String, Stats>, Box<dyn Error>> {
    let file = File::open(path)?;
    let mut reader = BufReader::with_capacity(1 << 20, file);
    let mut buffer: Vec<u8> = Vec::with_capacity(256);
    let mut stations: HashMap<String, Stats> = HashMap::with_capacity(16_384);

    loop {
        let bytes = reader.read_until(b'\n', &mut buffer)?;
        if bytes == 0 {
            break;
        }
        if let Some(&b'\n') = buffer.last() {
            buffer.pop();
        }
        if buffer.is_empty() {
            continue;
        }

        let mut sep_idx = None;
        for (idx, b) in buffer.iter().enumerate() {
            if *b == b';' {
                sep_idx = Some(idx);
                break;
            }
        }

        if let Some(pos) = sep_idx {
            let station = &buffer[..pos];
            let temp_raw = &buffer[pos + 1..];
            let temp = parse_temp(temp_raw);
            let key = String::from_utf8(station.to_vec())?;

            match stations.entry(key) {
                Entry::Vacant(v) => {
                    v.insert(Stats {
                        min: temp,
                        max: temp,
                        sum: temp as i64,
                        count: 1,
                    });
                }
                Entry::Occupied(mut o) => {
                    let s = o.get_mut();
                    if temp < s.min {
                        s.min = temp;
                    }
                    if temp > s.max {
                        s.max = temp;
                    }
                    s.sum += temp as i64;
                    s.count += 1;
                }
            }
        }

        buffer.clear();
    }

    Ok(stations)
}

fn format_tenths(value: i32) -> String {
    format!("{:.1}", value as f64 / 10.0)
}

fn format_mean(sum_tenths: i64, count: u32) -> String {
    let average = (sum_tenths as f64) / (count as f64) / 10.0;
    format!("{:.1}", average)
}

fn main() {
    let args: Vec<String> = env::args().collect();
    let path = args
        .get(1)
        .map(|s| s.as_str())
        .unwrap_or("../data/medium.txt");

    let file_path = Path::new(path);
    match process_file(file_path) {
        Ok(mut stats) => {
            let mut names: Vec<String> = stats.keys().cloned().collect();
            names.sort_unstable();

            let mut output = String::with_capacity(names.len() * 32);
            output.push('{');
            for (idx, name) in names.iter().enumerate() {
                if idx > 0 {
                    output.push(',');
                }
                if let Some(s) = stats.remove(name) {
                    output.push_str(name);
                    output.push('=');
                    output.push_str(&format_tenths(s.min));
                    output.push('/');
                    output.push_str(&format_mean(s.sum, s.count));
                    output.push('/');
                    output.push_str(&format_tenths(s.max));
                }
            }
            output.push('}');
            println!("{output}");
        }
        Err(err) => {
            eprintln!("Error processing file: {err}");
            std::process::exit(1);
        }
    }
}
