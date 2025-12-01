#!/usr/bin/env python3
import requests
from requests.adapters import HTTPAdapter
from urllib3.util.retry import Retry
from bs4 import BeautifulSoup
import re
from pathlib import Path
from concurrent.futures import ThreadPoolExecutor, as_completed
import threading
import time

BASE_URL = "https://projecteuler.net"
ARCHIVES_URL = f"{BASE_URL}/archives"
PROBLEMS_DIR = Path(__file__).parent / "problems"
MAX_WORKERS = 10
TIMEOUT = 30
MAX_RETRIES = 3

session_local = threading.local()

def get_session():
    if not hasattr(session_local, 'session'):
        session = requests.Session()
        session.headers.update({
            'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36'
        })
        retry_strategy = Retry(
            total=MAX_RETRIES,
            backoff_factor=1,
            status_forcelist=[429, 500, 502, 503, 504],
        )
        adapter = HTTPAdapter(max_retries=retry_strategy)
        session.mount("https://", adapter)
        session.mount("http://", adapter)
        session_local.session = session
    return session_local.session

def fetch_with_retry(url, max_attempts=MAX_RETRIES):
    session = get_session()
    for attempt in range(max_attempts):
        try:
            response = session.get(url, timeout=TIMEOUT)
            response.raise_for_status()
            return response
        except requests.RequestException as e:
            if attempt < max_attempts - 1:
                wait_time = 2 ** attempt
                print(f"Retry {attempt + 1}/{max_attempts} for {url} after {wait_time}s...")
                time.sleep(wait_time)
            else:
                raise e
    return None

def get_max_page_number():
    response = fetch_with_retry(ARCHIVES_URL)
    soup = BeautifulSoup(response.content, 'html.parser')
    
    max_page = 1
    for link in soup.find_all('a', href=True):
        href = link.get('href', '')
        match = re.search(r'page=(\d+)', href)
        if match:
            page_num = int(match.group(1))
            max_page = max(max_page, page_num)
    
    return max_page

def extract_problems_from_page(soup):
    problems = []
    
    table = soup.find('table')
    if not table:
        return problems
    
    rows = table.find_all('tr')[1:]
    for row in rows:
        cells = row.find_all('td')
        if len(cells) >= 2:
            title_cell = cells[1]
            link = title_cell.find('a')
            if link:
                href = link.get('href', '')
                match = re.search(r'problem=(\d+)', href)
                if match:
                    problem_id = int(match.group(1))
                    title = link.get_text(strip=True)
                    problems.append((problem_id, title))
    
    return problems

def fetch_archive_page(page):
    if page == 1:
        url = ARCHIVES_URL
    else:
        url = f"{BASE_URL}/archives;page={page}"
    
    try:
        response = fetch_with_retry(url)
        soup = BeautifulSoup(response.content, 'html.parser')
        return extract_problems_from_page(soup)
    except requests.RequestException as e:
        print(f"Failed to fetch page {page}: {e}")
        return []

def get_all_problems():
    print("Finding total number of pages...")
    try:
        max_page = get_max_page_number()
    except requests.RequestException as e:
        print(f"Failed to connect to Project Euler: {e}")
        print("Please check your internet connection and try again.")
        return []
    
    print(f"Found {max_page} pages, fetching problem list...")
    
    all_problems = []
    
    with ThreadPoolExecutor(max_workers=MAX_WORKERS) as executor:
        futures = {executor.submit(fetch_archive_page, page): page for page in range(1, max_page + 1)}
        for future in as_completed(futures):
            page = futures[future]
            try:
                problems = future.result()
                all_problems.extend(problems)
                print(f"Fetched page {page}/{max_page} ({len(problems)} problems)")
            except Exception as e:
                print(f"Error fetching page {page}: {e}")
    
    all_problems.sort(key=lambda x: x[0])
    print(f"Found {len(all_problems)} total problems")
    return all_problems

def clean_math_notation(text):
    text = re.sub(r'\$([^$]+)\$', r'\1', text)
    return text

def extract_problem_text(html_content):
    soup = BeautifulSoup(html_content, 'html.parser')
    
    problem_content = soup.find('div', class_='problem_content')
    if not problem_content:
        return None
    
    text = problem_content.get_text(separator='\n', strip=True)
    return clean_math_notation(text)

def download_and_save_problem(problem_id, title):
    problem_file = PROBLEMS_DIR / f"problem_{problem_id:04d}.txt"
    
    if problem_file.exists():
        return (problem_id, "skipped")
    
    url = f"{BASE_URL}/problem={problem_id}"
    
    try:
        response = fetch_with_retry(url)
        problem_text = extract_problem_text(response.text)
        
        if not problem_text:
            return (problem_id, "failed")
        
        content = f"""================================================================================
PROBLEM {problem_id}: {title}
================================================================================

{problem_text}

================================================================================
END OF PROBLEM {problem_id}
================================================================================
"""
        problem_file.write_text(content, encoding='utf-8')
        return (problem_id, "downloaded")
    except requests.RequestException:
        return (problem_id, "failed")

def main():
    PROBLEMS_DIR.mkdir(exist_ok=True)
    
    all_problems = get_all_problems()
    
    if not all_problems:
        print("No problems found. Exiting.")
        return
    
    print(f"\nDownloading {len(all_problems)} problems with {MAX_WORKERS} workers...")
    
    downloaded = 0
    skipped = 0
    failed = 0
    
    with ThreadPoolExecutor(max_workers=MAX_WORKERS) as executor:
        futures = {
            executor.submit(download_and_save_problem, pid, title): pid 
            for pid, title in all_problems
        }
        
        for future in as_completed(futures):
            problem_id, status = future.result()
            if status == "downloaded":
                downloaded += 1
                print(f"Downloaded problem {problem_id} ({downloaded + skipped + failed}/{len(all_problems)})")
            elif status == "skipped":
                skipped += 1
            else:
                failed += 1
                print(f"Failed problem {problem_id}")
    
    print(f"\nDownload complete!")
    print(f"Downloaded: {downloaded}")
    print(f"Skipped: {skipped}")
    print(f"Failed: {failed}")
    print(f"Total: {len(all_problems)}")

if __name__ == "__main__":
    main()
