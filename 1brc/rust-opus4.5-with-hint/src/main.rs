use memmap2::Mmap;
use rustc_hash::FxHashMap;
use std::env;
use std::fs::File;
use std::io::{self, Write};
use std::thread;

#[derive(Clone, Copy)]
struct Stats {
    min: i32,
    max: i32,
    sum: i64,
    count: u32,
}

impl Stats {
    #[inline(always)]
    fn new(value: i32) -> Self {
        Stats {
            min: value,
            max: value,
            sum: value as i64,
            count: 1,
        }
    }

    #[inline(always)]
    fn update(&mut self, value: i32) {
        self.min = self.min.min(value);
        self.max = self.max.max(value);
        self.sum += value as i64;
        self.count += 1;
    }

    #[inline(always)]
    fn merge(&mut self, other: &Stats) {
        self.min = self.min.min(other.min);
        self.max = self.max.max(other.max);
        self.sum += other.sum;
        self.count += other.count;
    }
}

#[inline(always)]
fn parse_temperature(bytes: &[u8]) -> i32 {
    let len = bytes.len();
    
    match len {
        3 => {
            let d1 = (bytes[0] - b'0') as i32;
            let d2 = (bytes[2] - b'0') as i32;
            d1 * 10 + d2
        }
        4 => {
            if bytes[0] == b'-' {
                let d1 = (bytes[1] - b'0') as i32;
                let d2 = (bytes[3] - b'0') as i32;
                -(d1 * 10 + d2)
            } else {
                let d1 = (bytes[0] - b'0') as i32;
                let d2 = (bytes[1] - b'0') as i32;
                let d3 = (bytes[3] - b'0') as i32;
                d1 * 100 + d2 * 10 + d3
            }
        }
        5 => {
            let d1 = (bytes[1] - b'0') as i32;
            let d2 = (bytes[2] - b'0') as i32;
            let d3 = (bytes[4] - b'0') as i32;
            -(d1 * 100 + d2 * 10 + d3)
        }
        _ => 0,
    }
}

#[inline(always)]
fn find_semicolon(data: &[u8], start: usize) -> usize {
    let mut i = start;
    while i < data.len() && unsafe { *data.get_unchecked(i) } != b';' {
        i += 1;
    }
    i
}

#[inline(always)]
fn find_newline(data: &[u8], start: usize) -> usize {
    let mut i = start;
    while i < data.len() && unsafe { *data.get_unchecked(i) } != b'\n' {
        i += 1;
    }
    i
}

fn process_chunk(data: &[u8]) -> FxHashMap<Vec<u8>, Stats> {
    let mut map: FxHashMap<Vec<u8>, Stats> = FxHashMap::with_capacity_and_hasher(10000, Default::default());
    
    let mut i = 0;
    let len = data.len();
    
    while i < len {
        let semi_pos = find_semicolon(data, i);
        if semi_pos >= len {
            break;
        }
        
        let station = unsafe { data.get_unchecked(i..semi_pos) };
        
        let temp_start = semi_pos + 1;
        let newline_pos = find_newline(data, temp_start);
        
        let temp_bytes = unsafe { data.get_unchecked(temp_start..newline_pos) };
        let temp = parse_temperature(temp_bytes);
        
        match map.get_mut(station) {
            Some(stats) => stats.update(temp),
            None => {
                map.insert(station.to_vec(), Stats::new(temp));
            }
        }
        
        i = newline_pos + 1;
    }
    
    map
}

fn find_line_boundary(data: &[u8], pos: usize) -> usize {
    if pos == 0 || pos >= data.len() {
        return pos.min(data.len());
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

fn format_temp(value: i32) -> String {
    if value < 0 {
        let abs = -value;
        format!("-{}.{}", abs / 10, abs % 10)
    } else {
        format!("{}.{}", value / 10, value % 10)
    }
}

fn main() -> io::Result<()> {
    let args: Vec<String> = env::args().collect();
    let file_path = if args.len() > 1 {
        &args[1]
    } else {
        "../data/medium.txt"
    };

    let file = File::open(file_path)?;
    let mmap = unsafe { Mmap::map(&file)? };
    let data: &[u8] = &mmap[..];
    
    if data.is_empty() {
        println!("{{}}");
        return Ok(());
    }

    let num_threads = thread::available_parallelism()
        .map(|n| n.get())
        .unwrap_or(4);

    let chunk_size = data.len() / num_threads;

    let mut boundaries = Vec::with_capacity(num_threads + 1);
    boundaries.push(0);
    
    for i in 1..num_threads {
        boundaries.push(find_line_boundary(data, i * chunk_size));
    }
    boundaries.push(data.len());

    let results: Vec<FxHashMap<Vec<u8>, Stats>> = thread::scope(|s| {
        let handles: Vec<_> = (0..num_threads)
            .map(|i| {
                let start = boundaries[i];
                let end = boundaries[i + 1];
                let chunk = &data[start..end];
                
                s.spawn(move || process_chunk(chunk))
            })
            .collect();

        handles.into_iter().map(|h| h.join().unwrap()).collect()
    });

    let mut final_map: FxHashMap<Vec<u8>, Stats> = FxHashMap::with_capacity_and_hasher(10000, Default::default());
    
    for partial in results {
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
    stations.sort_unstable_by(|a, b| a.0.cmp(&b.0));

    let mut output = Vec::with_capacity(stations.len() * 50);
    output.push(b'{');

    for (i, (station, stats)) in stations.iter().enumerate() {
        if i > 0 {
            output.push(b',');
        }
        
        output.extend_from_slice(station);
        output.push(b'=');
        
        let mean = ((stats.sum as f64 / stats.count as f64).round()) as i32;
        
        let min_str = format_temp(stats.min);
        let mean_str = format_temp(mean);
        let max_str = format_temp(stats.max);
        
        output.extend_from_slice(min_str.as_bytes());
        output.push(b'/');
        output.extend_from_slice(mean_str.as_bytes());
        output.push(b'/');
        output.extend_from_slice(max_str.as_bytes());
    }
    
    output.push(b'}');
    output.push(b'\n');

    let stdout = io::stdout();
    let mut handle = stdout.lock();
    handle.write_all(&output)?;

    Ok(())
}
